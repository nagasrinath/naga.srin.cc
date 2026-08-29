package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// testService wires a service at the given fake Spotify base URL, with a
// token already cached so tests that don't care about auth skip that hop.
func testService(base string) *service {
	s := &service{
		clientID:     "id",
		clientSecret: "secret",
		refreshToken: "refresh",
		allowedOrig:  "https://naga.srin.cc",
		tokenURL:     base + "/token",
		nowURL:       base + "/now",
		recentURL:    base + "/recent",
		httpClient:   &http.Client{Timeout: 2 * time.Second},
	}
	s.tokens.accessToken = "cached-token"
	s.tokens.expiresAt = time.Now().Add(time.Hour)
	return s
}

// expireCache forces the next read to miss, standing in for waiting out
// nowPlayingTTL.
func (s *service) expireCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	s.cache.attemptedAt = time.Time{}
}

func TestTrackResponse(t *testing.T) {
	var track spotifyTrack
	if err := json.Unmarshal([]byte(`{
		"name": "Svefn-g-englar",
		"artists": [{"name": "Sigur Rós"}, {"name": "Guest"}],
		"external_urls": {"spotify": "https://open.spotify.com/track/x"}
	}`), &track); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	got := track.response(true)
	want := nowPlayingResponse{
		Playing: true,
		Track:   "Svefn-g-englar",
		Artist:  "Sigur Rós, Guest",
		URL:     "https://open.spotify.com/track/x",
	}
	if got != want {
		t.Errorf("response(true) = %+v, want %+v", got, want)
	}

	if got := track.response(false); got.Playing {
		t.Error("response(false).Playing = true, want false")
	}
}

// A cache miss must produce exactly one upstream call no matter how many
// requests arrive at once.
func TestGetNowPlayingSingleFlight(t *testing.T) {
	var calls atomic.Int64
	release := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		<-release // hold the request open so every caller piles up
		w.Write([]byte(`{"is_playing":true,"item":{"name":"t","artists":[{"name":"a"}]}}`))
	}))
	defer srv.Close()

	s := testService(srv.URL)

	const n = 20
	var wg sync.WaitGroup
	results := make([]nowPlayingResponse, n)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = s.getNowPlaying()
		}()
	}

	// Give the goroutines time to contend, then let the upstream answer.
	time.Sleep(50 * time.Millisecond)
	close(release)
	wg.Wait()

	if got := calls.Load(); got != 1 {
		t.Errorf("upstream calls = %d, want 1", got)
	}
	for i, got := range results {
		if got.Track != "t" || !got.Playing {
			t.Errorf("result[%d] = %+v, want the fetched track", i, got)
		}
	}
}

func TestGetNowPlayingUsesCacheWithinTTL(t *testing.T) {
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"is_playing":true,"item":{"name":"t","artists":[]}}`))
	}))
	defer srv.Close()

	s := testService(srv.URL)
	s.getNowPlaying()
	s.getNowPlaying()

	if got := calls.Load(); got != 1 {
		t.Errorf("upstream calls = %d, want 1 (second read should hit the cache)", got)
	}
}

// A failing upstream must not blank out a good cached value.
func TestGetNowPlayingServesStaleOnError(t *testing.T) {
	var fail atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"is_playing":true,"item":{"name":"good","artists":[{"name":"a"}]}}`))
	}))
	defer srv.Close()

	s := testService(srv.URL)
	if got := s.getNowPlaying(); got.Track != "good" {
		t.Fatalf("warm-up track = %q, want %q", got.Track, "good")
	}

	fail.Store(true)
	s.expireCache()

	got := s.getNowPlaying()
	if got.Track != "good" || !got.Playing {
		t.Errorf("after upstream failure = %+v, want the previous good value", got)
	}
}

// With nothing cached yet, a failure degrades to "no signal" rather than
// surfacing an error.
func TestGetNowPlayingColdCacheFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	got := testService(srv.URL).getNowPlaying()
	if got != (nowPlayingResponse{Playing: false}) {
		t.Errorf("cold-cache failure = %+v, want {Playing:false}", got)
	}
}

// A failed refresh must not retry on every request; the TTL applies to
// attempts, not just successes.
func TestFailedRefreshIsRateLimited(t *testing.T) {
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	s := testService(srv.URL)
	s.getNowPlaying()
	s.getNowPlaying()
	s.getNowPlaying()

	if got := calls.Load(); got != 1 {
		t.Errorf("upstream calls = %d, want 1 (failures should be rate-limited too)", got)
	}
}

// 204 means nothing is playing, which should fall back to play history.
func TestFallsBackToRecentlyPlayed(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/now", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/recent", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"items":[{"track":{"name":"past","artists":[{"name":"b"}]}}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	got := testService(srv.URL).getNowPlaying()
	if got.Track != "past" || got.Artist != "b" || got.Playing {
		t.Errorf("got %+v, want the last-played track with Playing=false", got)
	}
}

// A token failure must not trigger the recently-played call, which would be
// guaranteed to fail for the same reason.
func TestTokenFailureSkipsRecentlyPlayed(t *testing.T) {
	var recentCalls atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	mux.HandleFunc("/recent", func(w http.ResponseWriter, r *http.Request) {
		recentCalls.Add(1)
		w.Write([]byte(`{"items":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	s := testService(srv.URL)
	s.tokens.accessToken = "" // force a refresh
	s.getNowPlaying()

	if got := recentCalls.Load(); got != 0 {
		t.Errorf("recently-played calls = %d, want 0 when the token refresh failed", got)
	}
}

func TestAccessTokenIsReusedUntilExpiry(t *testing.T) {
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"access_token":"fresh","expires_in":3600}`))
	}))
	defer srv.Close()

	s := testService(srv.URL)
	s.tokens.accessToken = ""

	for range 3 {
		tok, err := s.accessToken()
		if err != nil {
			t.Fatalf("accessToken: %v", err)
		}
		if tok != "fresh" {
			t.Fatalf("token = %q, want %q", tok, "fresh")
		}
	}

	if got := calls.Load(); got != 1 {
		t.Errorf("token endpoint calls = %d, want 1", got)
	}
}

func TestHandlerServesJSONWithCORS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"is_playing":true,"item":{"name":"t","artists":[{"name":"a"}]}}`))
	}))
	defer srv.Close()

	rec := httptest.NewRecorder()
	testService(srv.URL).handleNowPlaying(rec, httptest.NewRequest(http.MethodGet, "/now-playing", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://naga.srin.cc" {
		t.Errorf("CORS header = %q, want the configured origin", got)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var body nowPlayingResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Track != "t" || !body.Playing {
		t.Errorf("body = %+v, want the fetched track", body)
	}
}

func TestHandlerRejectsNonGET(t *testing.T) {
	rec := httptest.NewRecorder()
	testService("http://unused.invalid").handleNowPlaying(
		rec, httptest.NewRequest(http.MethodPost, "/now-playing", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got == "" {
		t.Error("CORS header missing on the 405 response")
	}
}
