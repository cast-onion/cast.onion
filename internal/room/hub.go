package room

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"
)

const MaxGuests = 4

var (
	ErrRoomFull      = errors.New("room is full")
	ErrRoomNotFound  = errors.New("room not found")
	ErrNotHost       = errors.New("not the host")
	ErrGuestNotFound = errors.New("guest not found")
)

type Guest struct {
	ID          string
	Name        string
	MutedByHost bool
	MutedSelf   bool
	audio       chan []byte
	done        chan struct{}
}

type Room struct {
	ID        string
	Code      string
	StationID string
	HostKey   string
	CreatedAt time.Time

	mu     sync.RWMutex
	guests map[string]*Guest
	mixer  *Mixer
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	codes map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		codes: make(map[string]*Room),
	}
}

func (h *Hub) CreateRoom(stationID, hostKey string) (*Room, error) {
	id, err := randomHex(16)
	if err != nil {
		return nil, err
	}
	code, err := randomHex(4)
	if err != nil {
		return nil, err
	}

	r := &Room{
		ID:        id,
		Code:      code,
		StationID: stationID,
		HostKey:   hostKey,
		CreatedAt: time.Now(),
		guests:    make(map[string]*Guest),
		mixer:     newMixer(),
	}

	h.mu.Lock()
	h.rooms[id] = r
	h.codes[code] = r
	h.mu.Unlock()

	return r, nil
}

func (h *Hub) GetByCode(code string) (*Room, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.codes[code]
	return r, ok
}

func (h *Hub) GetByID(id string) (*Room, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.rooms[id]
	return r, ok
}

func (h *Hub) DeleteRoom(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[id]; ok {
		delete(h.codes, r.Code)
		delete(h.rooms, id)
	}
}

func (r *Room) Join(name string) (*Guest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.guests) >= MaxGuests {
		return nil, ErrRoomFull
	}

	id, _ := randomHex(8)
	g := &Guest{
		ID:    id,
		Name:  name,
		audio: make(chan []byte, 64),
		done:  make(chan struct{}),
	}
	r.guests[id] = g
	return g, nil
}

func (r *Room) Leave(guestID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.guests[guestID]; ok {
		close(g.done)
		delete(r.guests, guestID)
	}
}

func (r *Room) MuteGuest(guestID string, muted bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.guests[guestID]
	if !ok {
		return ErrGuestNotFound
	}
	g.MutedByHost = muted
	return nil
}

func (r *Room) GuestMuteSelf(guestID string, muted bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.guests[guestID]
	if !ok {
		return ErrGuestNotFound
	}
	g.MutedSelf = muted
	return nil
}

func (r *Room) ListGuests() []GuestInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]GuestInfo, 0, len(r.guests))
	for _, g := range r.guests {
		out = append(out, GuestInfo{
			ID:          g.ID,
			Name:        g.Name,
			MutedByHost: g.MutedByHost,
			MutedSelf:   g.MutedSelf,
		})
	}
	return out
}

func (r *Room) StreamGuestAudio(guestID string, reader io.Reader, broadcastFn func([]byte)) error {
	r.mu.RLock()
	g, ok := r.guests[guestID]
	r.mu.RUnlock()
	if !ok {
		return ErrGuestNotFound
	}

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			r.mu.RLock()
			muted := g.MutedByHost || g.MutedSelf
			r.mu.RUnlock()
			if !muted {
				data := make([]byte, n)
				copy(data, buf[:n])
				broadcastFn(data)
			}
		}
		if err != nil {
			return err
		}
	}
}

type GuestInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MutedByHost bool   `json:"muted_by_host"`
	MutedSelf   bool   `json:"muted_self"`
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
