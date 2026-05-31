package room

import (
	"encoding/json"
	"net/http"

	"github.com/cast-onion/internal/service"
	"github.com/cast-onion/internal/stream"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	hub       *Hub
	streamHub *stream.Hub
	authSvc   *service.AuthService
}

func NewHandler(hub *Hub, streamHub *stream.Hub, authSvc *service.AuthService) *Handler {
	return &Handler{hub: hub, streamHub: streamHub, authSvc: authSvc}
}

type createRoomRequest struct {
	StationID string `json:"station_id"`
}

type createRoomResponse struct {
	RoomID string `json:"room_id"`
	Code   string `json:"code"`
	Link   string `json:"link"`
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	rawToken := r.Header.Get("X-Access-Token")
	stationID, err := h.authSvc.ValidateAccessToken(r.Context(), rawToken)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	room, err := h.hub.CreateRoom(stationID, rawToken)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createRoomResponse{
		RoomID: room.ID,
		Code:   room.Code,
		Link:   "cast://join/" + room.Code,
	})
}

type joinRoomRequest struct {
	Name string `json:"name"`
}

type joinRoomResponse struct {
	RoomID    string      `json:"room_id"`
	GuestID   string      `json:"guest_id"`
	StationID string      `json:"station_id"`
	Guests    []GuestInfo `json:"guests"`
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var body joinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		body.Name = "guest"
	}

	room, ok := h.hub.GetByCode(code)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	guest, err := room.Join(body.Name)
	if err == ErrRoomFull {
		http.Error(w, "room is full", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(joinRoomResponse{
		RoomID:    room.ID,
		GuestID:   guest.ID,
		StationID: room.StationID,
		Guests:    room.ListGuests(),
	})
}

func (h *Handler) GuestStream(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	guestID := r.Header.Get("X-Guest-ID")

	room, ok := h.hub.GetByID(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	room.StreamGuestAudio(guestID, r.Body, func(data []byte) {
		h.streamHub.Broadcast(room.StationID, data)
	})
}

func (h *Handler) MuteGuest(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	guestID := chi.URLParam(r, "guest_id")
	rawToken := r.Header.Get("X-Access-Token")

	room, ok := h.hub.GetByID(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	_, err := h.authSvc.ValidateAccessToken(r.Context(), rawToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Muted bool `json:"muted"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := room.MuteGuest(guestID, body.Muted); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SelfMute(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	guestID := r.Header.Get("X-Guest-ID")

	room, ok := h.hub.GetByID(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	var body struct {
		Muted bool `json:"muted"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	room.GuestMuteSelf(guestID, body.Muted)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RoomInfo(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	room, ok := h.hub.GetByID(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"room_id":    room.ID,
		"code":       room.Code,
		"station_id": room.StationID,
		"guests":     room.ListGuests(),
	})
}
