package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type SessionCreator func(ctx context.Context, ip, userAgent string) (string, error)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler(hub *Hub, createSession SessionCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}

		sessionID, err := createSession(r.Context(), r.RemoteAddr, r.UserAgent())
		if err != nil {
			log.Printf("ws createSession error: %v", err)
			conn.Close()
			return
		}

		log.Printf("ws session created: %s", sessionID)

		conn.WriteJSON(map[string]string{
			"type":       "session",
			"session_id": sessionID,
		})

		client := &Client{
			conn:      conn,
			SessionID: sessionID,
			send:      make(chan []byte, 64),
		}
		hub.register <- client

		go writePump(client)
		go readPump(client, hub)
	}
}

func writePump(c *Client) {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func readPump(c *Client, hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func PingLoop(c *Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}
