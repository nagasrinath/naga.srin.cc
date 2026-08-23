package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	spotifyTokenURL = "https://accounts.spotify.com/api/token"
	spotifyNowURL   = "https://api.spotify.com/v1/me/player/currently-playing"
	nowPlayingTTL   = 20 * time.Second
)

type nowPlayingResponse struct {
	Playing bool   `json:"playing"`
	Track   string `json:"track,omitempty"`
	Artist  string `json:"artist,omitempty"`
}

type tokenCache struct {
	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

type nowPlayingCache struct {
	mu        sync.Mutex
	resp      nowPlayingResponse
	fetchedAt time.Time
}

type service struct {
	clientID     string
	clientSecret string
	refreshToken string
	allowedOrig  string
	httpClient   *http.Client
	tokens       tokenCache
	cache        nowPlayingCache
}

func main() {
	svc := &service{
		clientID:     mustEnv("SPOTIFY_CLIENT_ID"),
		clientSecret: mustEnv("SPOTIFY_CLIENT_SECRET"),
		refreshToken: mustEnv("SPOTIFY_REFRESH_TOKEN"),
		allowedOrig:  envOr("ALLOWED_ORIGIN", "https://naga.srin.cc"),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/now-playing", svc.handleNowPlaying)

	addr := ":" + envOr("PORT", "8080")
	log.Printf("now-playing listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var %s", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (s *service) handleNowPlaying(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Access-Control-Allow-Origin", s.allowedOrig)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(nowPlayingResponse{Playing: false})
		return
	}

	resp := s.getNowPlaying()
	json.NewEncoder(w).Encode(resp)
}

func (s *service) getNowPlaying() nowPlayingResponse {
	s.cache.mu.Lock()
	if time.Since(s.cache.fetchedAt) < nowPlayingTTL {
		resp := s.cache.resp
		s.cache.mu.Unlock()
		return resp
	}
	s.cache.mu.Unlock()

	resp := s.fetchNowPlaying()

	s.cache.mu.Lock()
	s.cache.resp = resp
	s.cache.fetchedAt = time.Now()
	s.cache.mu.Unlock()

	return resp
}

func (s *service) fetchNowPlaying() nowPlayingResponse {
	fallback := nowPlayingResponse{Playing: false}

	token, err := s.accessToken()
	if err != nil {
		log.Printf("access token error: %v", err)
		return fallback
	}

	req, err := http.NewRequest(http.MethodGet, spotifyNowURL, nil)
	if err != nil {
		log.Printf("build request error: %v", err)
		return fallback
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("spotify request error: %v", err)
		return fallback
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return fallback
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Printf("spotify unexpected status %d: %s", res.StatusCode, body)
		return fallback
	}

	var payload struct {
		IsPlaying bool `json:"is_playing"`
		Item      *struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"item"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		log.Printf("decode error: %v", err)
		return fallback
	}

	if !payload.IsPlaying || payload.Item == nil {
		return fallback
	}

	names := make([]string, 0, len(payload.Item.Artists))
	for _, a := range payload.Item.Artists {
		names = append(names, a.Name)
	}

	return nowPlayingResponse{
		Playing: true,
		Track:   payload.Item.Name,
		Artist:  strings.Join(names, ", "),
	}
}

func (s *service) accessToken() (string, error) {
	s.tokens.mu.Lock()
	defer s.tokens.mu.Unlock()

	if s.tokens.accessToken != "" && time.Now().Before(s.tokens.expiresAt) {
		return s.tokens.accessToken, nil
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", s.refreshToken)

	req, err := http.NewRequest(http.MethodPost, spotifyTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	basic := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+basic)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", &tokenError{status: res.StatusCode, body: string(body)}
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}

	if payload.RefreshToken != "" && payload.RefreshToken != s.refreshToken {
		log.Printf("spotify issued a new refresh token; update SPOTIFY_REFRESH_TOKEN to keep working long-term")
	}

	s.tokens.accessToken = payload.AccessToken
	// refresh a little early to avoid edge-of-expiry races
	s.tokens.expiresAt = time.Now().Add(time.Duration(payload.ExpiresIn-30) * time.Second)

	return s.tokens.accessToken, nil
}

type tokenError struct {
	status int
	body   string
}

func (e *tokenError) Error() string {
	return "spotify token refresh failed: status " + http.StatusText(e.status) + ": " + e.body
}
