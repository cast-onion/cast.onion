package stream

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/cast-onion/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	hub       *Hub
	streamSvc *service.StreamService
}

func NewHandler(hub *Hub, streamSvc *service.StreamService) *Handler {
	return &Handler{hub: hub, streamSvc: streamSvc}
}

func (h *Handler) Broadcast(w http.ResponseWriter, r *http.Request) {
	stationID := chi.URLParam(r, "station_id")
	rawKey := r.Header.Get("X-Station-Key")

	if err := h.streamSvc.ValidateStationKey(r.Context(), rawKey); err != nil {
		log.Printf("broadcast auth failed for station %s: %v", stationID, err)
		http.Error(w, "invalid station key", http.StatusUnauthorized)
		return
	}

	log.Printf("broadcast started: station=%s", stationID)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<40)

	if err := h.hub.StartBroadcast(stationID, r.Body); err != nil && err != io.EOF {
		log.Printf("broadcast ended: station=%s err=%v", stationID, err)
	} else {
		log.Printf("broadcast ended: station=%s", stationID)
	}
}

func (h *Handler) Listen(w http.ResponseWriter, r *http.Request) {
	stationID := chi.URLParam(r, "station_id")

	listener, ok := h.hub.AddListener(stationID)
	if !ok {
		http.Error(w, "station not live", http.StatusServiceUnavailable)
		return
	}
	defer h.hub.RemoveListener(stationID, listener)

	w.Header().Set("Content-Type", "audio/ogg")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	flusher, canFlush := w.(http.Flusher)
	if canFlush {
		flusher.Flush()
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case data, ok := <-listener.ch:
			if !ok {
				return
			}
			if _, err := w.Write(data); err != nil {
				return
			}
			if canFlush {
				flusher.Flush()
			}
		}
	}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	stationID := chi.URLParam(r, "station_id")
	live := h.hub.IsLive(stationID)
	count := h.hub.ViewerCount(stationID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"station_id": stationID,
		"live":       live,
		"viewers":    count,
	})
}

func (h *Handler) ViewerCount(w http.ResponseWriter, r *http.Request) {
	stationID := chi.URLParam(r, "station_id")
	count := h.hub.ViewerCount(stationID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"station_id": stationID,
		"count":      count,
	})
}
