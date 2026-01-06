package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestSleepBackoff_Cancelled(t *testing.T) {
	s := NewStockService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if s.sleepBackoff(ctx, 0) {
		t.Fatalf("expected sleepBackoff to return false when ctx is canceled")
	}
}

func TestSSEClient_NoTimeout(t *testing.T) {
	s := NewStockService()
	if s.sseClient == nil {
		t.Fatalf("expected sseClient to be initialized")
	}
	if s.sseClient.Timeout != 0 {
		t.Fatalf("expected sseClient.Timeout=0 for long-lived SSE, got %v", s.sseClient.Timeout)
	}
}

func TestStream_StopReplacesOld(t *testing.T) {
	s := NewStockService()
	s.Startup(context.Background())

	// Override emitter so tests don't depend on Wails runtime context.
	s.emitIntraday = func(ctx context.Context, code string, trends []string) {
		// no-op
	}

	// avoid real network calls by intercepting with a test server transport
	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		// immediately close
		_, _ = io.WriteString(w, "data: {\"data\": {\"trends\": [\"t,0,0,0,0,0,0,0\"]}}\n\n")
	}))
	defer ts.Close()

	// Replace sseClient transport to route any request to our server.
	s.sseClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = ts.Listener.Addr().String()
		return http.DefaultTransport.RoundTrip(req)
	})

	code := "600519"
	// start twice, second should cancel first
	s.StreamIntradayData(code)
	s.StreamIntradayData(code)

	// wait a bit to allow goroutines to run
	time.Sleep(200 * time.Millisecond)

	s.StopIntradayStream(code)

	// we should have at least one attempt
	if atomic.LoadInt32(&hits) == 0 {
		t.Fatalf("expected at least one SSE connect attempt")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
