package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestEndpointsLifecycle(t *testing.T) {
	srv := &Server{store: NewEventStore()}
	handler := srv.routes()

	reqHealth := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recHealth := httptest.NewRecorder()
	handler.ServeHTTP(recHealth, reqHealth)
	if recHealth.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /healthz, got %d", recHealth.Code)
	}

	validBody := []byte(`{"session_id": "call-123", "event_type": "call.started"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(validBody))
	recPost := httptest.NewRecorder()
	handler.ServeHTTP(recPost, reqPost)
	if recPost.Code != http.StatusAccepted {
		t.Fatalf("expected 202 Accepted, got %d", recPost.Code)
	}

	invalidBody := []byte(`{"session_id": "call-123"}`)
	reqBad := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(invalidBody))
	recBad := httptest.NewRecorder()
	handler.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", recBad.Code)
	}

	reqStats := httptest.NewRequest(http.MethodGet, "/stats", nil)
	recStats := httptest.NewRecorder()
	handler.ServeHTTP(recStats, reqStats)
	if recStats.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /stats, got %d", recStats.Code)
	}

	var stats map[string]int
	if err := json.NewDecoder(recStats.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode stats response: %v", err)
	}
	if stats["call.started"] != 1 {
		t.Errorf("expected call.started count to be 1, got %d", stats["call.started"])
	}
}

func TestConcurrentEvents(t *testing.T) {
	srv := &Server{store: NewEventStore()}
	handler := srv.routes()

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			payload := []byte(`{"session_id": "call-999", "event_type": "call.ended"}`)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(payload))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusAccepted {
				t.Errorf("expected 202 Accepted, got %d", rec.Code)
			}
		}()
	}

	wg.Wait()

	stats := srv.store.GetStats()
	if stats["call.ended"] != workers {
		t.Errorf("expected %d events, got %d", workers, stats["call.ended"])
	}
}