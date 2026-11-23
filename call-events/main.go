package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

type CallEvent struct {
	SessionID string `json:"session_id"`
	EventType string `json:"event_type"`
}

type EventStore struct {
	mu     sync.RWMutex
	counts map[string]int
}

func NewEventStore() *EventStore {
	return &EventStore{
		counts: make(map[string]int),
	}
}

func (s *EventStore) Increment(eventType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[eventType]++
}

func (s *EventStore) GetStats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statsCopy := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		statsCopy[k] = v
	}
	return statsCopy
}

type Server struct {
	store *EventStore
}

func (srv *Server) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	var event CallEvent
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&event); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(event.SessionID) == "" || strings.TrimSpace(event.EventType) == "" {
		http.Error(w, "session_id and event_type are required", http.StatusBadRequest)
		return
	}

	srv.store.Increment(event.EventType)
	w.WriteHeader(http.StatusAccepted)
}

func (srv *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stats := srv.store.GetStats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func (srv *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (srv *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", srv.handlePostEvent)
	mux.HandleFunc("GET /stats", srv.handleGetStats)
	mux.HandleFunc("GET /healthz", srv.handleHealthz)
	return mux
}

func main() {
	srv := &Server{store: NewEventStore()}
	log.Println("call-events service listening on :8080")
	if err := http.ListenAndServe(":8080", srv.routes()); err != nil {
		log.Fatal(err)
	}
}
