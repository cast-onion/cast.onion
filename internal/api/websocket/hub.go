package websocket

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type      string `json:"type"`
	StationID string `json:"station_id,omitempty"`
	Count     int    `json:"count,omitempty"`
	Payload   any    `json:"payload,omitempty"`
}

type Client struct {
	conn interface {
		WriteMessage(messageType int, data []byte) error
		ReadMessage() (messageType int, p []byte, err error)
		Close() error
	}
	SessionID string
	send      chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(eventType, stationID string) {
	b, err := json.Marshal(Event{Type: eventType, StationID: stationID})
	if err != nil {
		return
	}
	h.broadcast <- b
}

func (h *Hub) BroadcastViewerCount(stationID string, count int) {
	b, _ := json.Marshal(Event{
		Type:      "viewer_count",
		StationID: stationID,
		Count:     count,
	})
	h.broadcast <- b
}

func (h *Hub) BroadcastRaw(data []byte) {
	h.broadcast <- data
}

func ViewerCountMessage(stationID string, count int) []byte {
	b, _ := json.Marshal(fmt.Sprintf(`{"type":"viewer_count","station_id":"%s","count":%d}`, stationID, count))
	return b
}
