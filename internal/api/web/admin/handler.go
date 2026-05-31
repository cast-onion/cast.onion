package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/model"
	"github.com/cast-onion/internal/service"
)

type Handler struct {
	templates *template.Template
	cfg       *AdminConfig
	q         *queries.Queries
	adminSvc  *service.AdminService
}

type DashboardData struct {
	Username     string
	Applications []*model.Application
	Stations     []*model.Station
	Actions      []*model.AdminAction
	Keys         map[string]string
	Tokens       map[string]string
}

func NewHandler(templateDir string, cfg *AdminConfig, q *queries.Queries, adminSvc *service.AdminService) (*Handler, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"timeAgo": timeAgo,
	}).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}
	return &Handler{templates: tmpl, cfg: cfg, q: q, adminSvc: adminSvc}, nil
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin", h.Login)
	mux.HandleFunc("/admin/login", h.DoLogin)
	mux.HandleFunc("/admin/logout", h.DoLogout)
	mux.Handle("/admin/dashboard", RequireAdminSession(http.HandlerFunc(h.Dashboard)))
	mux.Handle("/admin/approve/", RequireAdminSession(http.HandlerFunc(h.Approve)))
	mux.Handle("/admin/deny/", RequireAdminSession(http.HandlerFunc(h.Deny)))
	mux.Handle("/admin/suspend/", RequireAdminSession(http.HandlerFunc(h.Suspend)))
	mux.Handle("/admin/revoke/", RequireAdminSession(http.HandlerFunc(h.Revoke)))
	mux.Handle("/admin/unsuspend/", RequireAdminSession(http.HandlerFunc(h.Unsuspend)))
	mux.Handle("/admin/restart", RequireAdminSession(http.HandlerFunc(h.Restart)))
	return mux
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if _, ok := ValidateSession(r); ok {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
		return
	}
	h.templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) DoLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !h.cfg.Authenticate(username, password) {
		h.templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "invalid credentials"})
		return
	}

	SetSessionCookie(w, username)
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

func (h *Handler) DoLogout(w http.ResponseWriter, r *http.Request) {
	ClearSessionCookie(w)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	username, _ := ValidateSession(r)
	apps, _ := h.q.ListAllApplications(r.Context())
	stationsList, _ := h.q.ListAllStations(r.Context())
	actions, _ := h.q.ListAdminActions(r.Context())

	keys := make(map[string]string)
	tokens := make(map[string]string)
	for _, st := range stationsList {
		key, _ := h.q.GetLatestStationKey(r.Context(), st.ID)
		token, _ := h.q.GetLatestAccessToken(r.Context(), st.ID)
		keys[st.ID] = key
		tokens[st.ID] = token
	}

	// support JSON for Svelte frontend
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Applications": apps,
			"Stations":     stationsList,
			"Actions":      actions,
			"Keys":         keys,
			"Tokens":       tokens,
		})
		return
	}

	h.templates.ExecuteTemplate(w, "dashboard.html", DashboardData{
		Username:     username,
		Applications: apps,
		Stations:     stationsList,
		Actions:      actions,
		Keys:         keys,
		Tokens:       tokens,
	})
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	username, _ := ValidateSession(r)

	result, err := h.adminSvc.ApproveApplication(r.Context(), username, id, "")
	if err != nil {
		log.Printf("approve error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"station_id":   result.StationID,
		"station_key":  result.RawStationKey,
		"access_token": result.RawAccessToken,
	})
}

func (h *Handler) Deny(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	username, _ := ValidateSession(r)

	if err := h.adminSvc.DenyApplication(r.Context(), username, id, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Suspend(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	username, _ := ValidateSession(r)
	h.adminSvc.SuspendStation(r.Context(), username, id, "")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	username, _ := ValidateSession(r)
	h.adminSvc.RevokeStation(r.Context(), username, id, "")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Unsuspend(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	username, _ := ValidateSession(r)
	h.adminSvc.UnsuspendStation(r.Context(), username, id, "")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("restarting..."))
	go func() {
		time.Sleep(500 * time.Millisecond)
		p, err := os.FindProcess(os.Getpid())
		if err == nil {
			p.Signal(syscall.SIGTERM)
		}
	}()
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return t.Format("2 Jan 2006")
	}
}
