package rest

import (
	"net/http"
	"time"

	"github.com/cast-onion/internal/api/rest/handlers"
	"github.com/cast-onion/internal/api/rest/middleware"
	ws "github.com/cast-onion/internal/api/websocket"
	"github.com/cast-onion/internal/cache"
	"github.com/cast-onion/internal/config"
	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/email"
	"github.com/cast-onion/internal/room"
	"github.com/cast-onion/internal/service"
	"github.com/cast-onion/internal/stream"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB, redis *cache.Redis, hub *ws.Hub, streamHub *stream.Hub, roomHub *room.Hub, mailer *email.Dispatcher, cfg *config.Config) http.Handler {
	q := queries.New(db)

	sessionSvc := service.NewSessionService(q, redis)
	stationSvc := service.NewStationService(q)
	appSvc := service.NewApplicationService(q)
	authSvc := service.NewAuthService(q)
	adminSvc := service.NewAdminService(q, stationSvc, authSvc, hub, mailer)
	streamSvc := service.NewStreamService(q)

	stationH := handlers.NewStationHandler(stationSvc)
	appH := handlers.NewApplicationHandler(appSvc, mailer)
	adminH := handlers.NewAdminHandler(adminSvc)
	ownerH := handlers.NewOwnerHandler(authSvc, stationSvc)
	streamH := stream.NewHandler(streamHub, streamSvc)
	roomH := room.NewHandler(roomHub, streamHub, authSvc)

	r := chi.NewRouter()
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RateLimit(redis, 120, time.Minute))

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/v1/docs", http.StatusFound)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/docs", handlers.Docs)

		r.Get("/ws", func(w http.ResponseWriter, req *http.Request) {
			ws.Handler(hub, sessionSvc.Create)(w, req)
		})

		r.Post("/broadcast/{station_id}", streamH.Broadcast)
		r.Get("/listen/{station_id}", streamH.Listen)
		r.Get("/status/{station_id}", streamH.Status)
		r.Get("/viewers/{station_id}", streamH.ViewerCount)

		r.Post("/room/create", roomH.CreateRoom)
		r.Post("/room/join/{code}", roomH.JoinRoom)
		r.Get("/room/{room_id}", roomH.RoomInfo)
		r.Post("/room/{room_id}/stream", roomH.GuestStream)
		r.Post("/room/{room_id}/mute/{guest_id}", roomH.MuteGuest)
		r.Post("/room/{room_id}/selfmute", roomH.SelfMute)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireSession(sessionSvc))

			r.Get("/stations", stationH.Directory)
			r.Get("/stations/{id}", stationH.Get)

			r.Post("/apply", appH.Submit)
			r.Get("/apply/{id}", appH.Get)

			r.Get("/owner/dashboard", ownerH.Dashboard)
			r.Patch("/owner/station", ownerH.UpdateStation)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdmin(q))

			r.Get("/admin/applications", adminH.ListPending)
			r.Post("/admin/applications/{id}/approve", adminH.Approve)
			r.Post("/admin/applications/{id}/deny", adminH.Deny)
			r.Post("/admin/stations/{id}/suspend", adminH.Suspend)
			r.Post("/admin/stations/{id}/revoke", adminH.Revoke)
			r.Post("/admin/stations/{id}/unsuspend", adminH.Unsuspend)
		})
	})

	return r
}
