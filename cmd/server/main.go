package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/cast-onion/internal/api/rest"
	webadmin "github.com/cast-onion/internal/api/web/admin"
	ws "github.com/cast-onion/internal/api/websocket"
	"github.com/cast-onion/internal/cache"
	"github.com/cast-onion/internal/cdn"
	"github.com/cast-onion/internal/config"
	"github.com/cast-onion/internal/db"
	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/email"
	"github.com/cast-onion/internal/room"
	"github.com/cast-onion/internal/service"
	"github.com/cast-onion/internal/stream"
)

func startSvelteDevServer(siteDir string, port string) *exec.Cmd {
	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = siteDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PORT="+port)
	if err := cmd.Start(); err != nil {
		log.Printf("web: could not start npm run dev: %v", err)
		return nil
	}
	log.Printf("cast.onion web  → http://localhost:%s (svelte dev)", port)
	return cmd
}

func findSiteDir() string {
	p := "web"
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(filepath.Join(abs, "package.json")); err == nil {
		return abs
	}
	return ""
}

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(cfg.DatabaseDSN); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	redisClient := cache.NewRedis(cfg.RedisAddr)
	q := queries.New(database)

	hub := ws.NewHub()
	go hub.Run()

	streamHub := stream.NewHub()
	streamHub.SetViewerCountCallback(func(stationID string, count int) {
		hub.BroadcastViewerCount(stationID, count)
	})

	roomHub := room.NewHub()
	mailer := email.NewDispatcher()

	stationSvc := service.NewStationService(q)
	authSvc := service.NewAuthService(q)
	adminSvc := service.NewAdminService(q, stationSvc, authSvc, hub, mailer)

	// =========================
	// ADMIN HANDLER (Go templates — /admin/*)
	// =========================
	adminCfg, err := webadmin.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("admin config: %v", err)
	}

	adminH, err := webadmin.NewHandler("web/templates/admin", adminCfg, q, adminSvc)
	if err != nil {
		log.Fatalf("admin templates: %v", err)
	}

	// =========================
	// SVELTE DEV SERVER
	// =========================
	siteDir := findSiteDir()
	var svelteCmd *exec.Cmd
	if siteDir != "" {
		svelteCmd = startSvelteDevServer(siteDir, cfg.WebPort)
	} else {
		log.Println("web: cast-onion-site not found, skipping npm run dev")
	}

	// =========================
	// ADMIN-ONLY WEB SERVER (for /admin/* Go template routes)
	// =========================
	adminMux := http.NewServeMux()
	adminRoutes := adminH.Routes()
	adminMux.Handle("/admin", adminRoutes)
	adminMux.Handle("/admin/", adminRoutes)

	adminSrv := &http.Server{
		Addr:         ":2051",
		Handler:      adminMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// =========================
	// API ROUTER
	// =========================
	apiRouter := rest.NewRouter(database, redisClient, hub, streamHub, roomHub, mailer, cfg)

	apiSrv := &http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     apiRouter,
		IdleTimeout: 120 * time.Second,
	}

	// =========================
	// CDN SERVER
	// =========================
	cdnMux := http.NewServeMux()
	cdnMux.HandleFunc("/", cdn.Handler)

	cdnSrv := &http.Server{
		Addr:         ":" + cfg.CDNPort,
		Handler:      cdnMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// =========================
	// START SERVERS
	// =========================
	go func() {
		log.Printf("cast.onion api   → http://localhost:%s/v1/docs", cfg.Port)
		if err := apiSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("api server: %v", err)
		}
	}()

	go func() {
		log.Printf("cast.onion admin → http://localhost:2051/admin")
		if err := adminSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("admin server: %v", err)
		}
	}()

	go func() {
		log.Printf("cast.onion cdn   → http://localhost:%s", cfg.CDNPort)
		if err := cdnSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cdn server: %v", err)
		}
	}()

	// =========================
	// SHUTDOWN HANDLER
	// =========================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	if svelteCmd != nil && svelteCmd.Process != nil {
		svelteCmd.Process.Kill()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	apiSrv.Shutdown(ctx)
	adminSrv.Shutdown(ctx)
	cdnSrv.Shutdown(ctx)

	log.Println("servers stopped")
}
