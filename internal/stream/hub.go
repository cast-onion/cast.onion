package stream

import (
	"io"
	"sync"
)

type ViewerCountFn func(stationID string, count int)

type Listener struct {
	ch   chan []byte
	done chan struct{}
}

type Station struct {
	mu        sync.RWMutex
	listeners map[*Listener]struct{}
	live      bool
}

func newStation() *Station {
	return &Station{listeners: make(map[*Listener]struct{})}
}

func (s *Station) addListener() *Listener {
	l := &Listener{ch: make(chan []byte, 256), done: make(chan struct{})}
	s.mu.Lock()
	s.listeners[l] = struct{}{}
	s.mu.Unlock()
	return l
}

func (s *Station) removeListener(l *Listener) {
	s.mu.Lock()
	delete(s.listeners, l)
	s.mu.Unlock()
	close(l.done)
}

func (s *Station) listenerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.listeners)
}

func (s *Station) broadcast(buf []byte) {
	data := make([]byte, len(buf))
	copy(data, buf)
	s.mu.RLock()
	for l := range s.listeners {
		select {
		case l.ch <- data:
		default:
		}
	}
	s.mu.RUnlock()
}

type Hub struct {
	mu             sync.RWMutex
	stations       map[string]*Station
	onViewerChange ViewerCountFn
}

func NewHub() *Hub {
	return &Hub{stations: make(map[string]*Station)}
}

func (h *Hub) SetViewerCountCallback(fn ViewerCountFn) {
	h.onViewerChange = fn
}

func (h *Hub) getOrCreate(stationID string) *Station {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s, ok := h.stations[stationID]; ok {
		return s
	}
	s := newStation()
	h.stations[stationID] = s
	return s
}

func (h *Hub) get(stationID string) (*Station, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	s, ok := h.stations[stationID]
	return s, ok
}

func (h *Hub) IsLive(stationID string) bool {
	s, ok := h.get(stationID)
	if !ok {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.live
}

func (h *Hub) ViewerCount(stationID string) int {
	s, ok := h.get(stationID)
	if !ok {
		return 0
	}
	return s.listenerCount()
}

func (h *Hub) StartBroadcast(stationID string, r io.Reader) error {
	station := h.getOrCreate(stationID)
	station.mu.Lock()
	station.live = true
	station.mu.Unlock()

	defer func() {
		station.mu.Lock()
		station.live = false
		station.mu.Unlock()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			station.broadcast(buf[:n])
		}
		if err != nil {
			return err
		}
	}
}

func (h *Hub) AddListener(stationID string) (*Listener, bool) {
	s, ok := h.get(stationID)
	if !ok || !s.live {
		return nil, false
	}
	l := s.addListener()
	if h.onViewerChange != nil {
		h.onViewerChange(stationID, s.listenerCount())
	}
	return l, true
}

func (h *Hub) RemoveListener(stationID string, l *Listener) {
	if s, ok := h.get(stationID); ok {
		s.removeListener(l)
		if h.onViewerChange != nil {
			h.onViewerChange(stationID, s.listenerCount())
		}
	}
}

func (h *Hub) Broadcast(stationID string, data []byte) {
	if s, ok := h.get(stationID); ok {
		s.broadcast(data)
	}
}
