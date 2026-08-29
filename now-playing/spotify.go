package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	spotifyTokenURL  = "https://accounts.spotify.com/api/token"
	spotifyNowURL    = "https://api.spotify.com/v1/me/player/currently-playing"
	spotifyRecentURL = "https://api.spotify.com/v1/me/player/recently-played?limit=1"

	nowPlayingTTL  = 20 * time.Second
	spotifyTimeout = 5 * time.Second
)

type nowPlayingResponse struct {
	Playing bool   `json:"playing"`
	Track   string `json:"track,omitempty"`
	Artist  string `json:"artist,omitempty"`
	URL     string `json:"url,omitempty"`
}

type spotifyTrack struct {
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func (t spotifyTrack) response(playing bool) nowPlayingResponse {
	names := make([]string, 0, len(t.Artists))
	for _, a := range t.Artists {
		names = append(names, a.Name)
	}
	return nowPlayingResponse{
		Playing: playing,
		Track:   t.Name,
		Artist:  strings.Join(names, ", "),
		URL:     t.ExternalURLs.Spotify,
	}
}

type tokenCache struct {
	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

type nowPlayingCache struct {
	mu   sync.Mutex
	resp nowPlayingResponse
	// ok reports whether resp was ever populated by a successful call.
	ok bool
	// attemptedAt is the last upstream attempt, successful or not. Failures
	// bump it without touching resp, so an outage keeps serving the last
	// known track instead of retrying on every single request.
	attemptedAt time.Time
}

type service struct {
	clientID     string
	clientSecret string
	refreshToken string
	allowedOrig  string

	// Endpoints are fields rather than consts so tests can point them at an
	// httptest server.
	tokenURL  string
	nowURL    string
	recentURL string

	httpClient *http.Client
	tokens     tokenCache
	cache      nowPlayingCache
	// fetchMu admits one goroutine at a time to the Spotify refresh path.
	fetchMu sync.Mutex
}

func newService() *service {
	return &service{
		clientID:     mustEnv("SPOTIFY_CLIENT_ID"),
		clientSecret: mustEnv("SPOTIFY_CLIENT_SECRET"),
		refreshToken: mustEnv("SPOTIFY_REFRESH_TOKEN"),
		allowedOrig:  envOr("ALLOWED_ORIGIN", "https://naga.srin.cc"),
		tokenURL:     spotifyTokenURL,
		nowURL:       spotifyNowURL,
		recentURL:    spotifyRecentURL,
		httpClient:   &http.Client{Timeout: spotifyTimeout},
	}
}

// cached returns what we'd serve right now, plus whether it is still within
// the TTL. The response is safe to serve either way; the bool only decides
// whether a refresh is due.
func (s *service) cached() (nowPlayingResponse, bool) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	resp := s.cache.resp
	if !s.cache.ok {
		resp = nowPlayingResponse{Playing: false}
	}
	return resp, time.Since(s.cache.attemptedAt) < nowPlayingTTL
}

func (s *service) getNowPlaying() nowPlayingResponse {
	if resp, fresh := s.cached(); fresh {
		return resp
	}

	// One goroutine refreshes; the rest queue here and then read what it
	// stored. Without this, every request arriving after the TTL expired
	// would issue its own Spotify call.
	s.fetchMu.Lock()
	defer s.fetchMu.Unlock()

	// The goroutine ahead of us may have just refreshed.
	if resp, fresh := s.cached(); fresh {
		return resp
	}

	resp, err := s.fetchNowPlaying()

	s.cache.mu.Lock()
	s.cache.attemptedAt = time.Now()
	if err == nil {
		s.cache.resp, s.cache.ok = resp, true
	}
	current, ok := s.cache.resp, s.cache.ok
	s.cache.mu.Unlock()

	if err != nil {
		log.Printf("now-playing refresh failed, serving last known value: %v", err)
	}
	if !ok {
		return nowPlayingResponse{Playing: false}
	}
	return current
}

// fetchNowPlaying distinguishes "Spotify says nothing is playing" (a real
// answer, worth caching) from "the call failed" (an error, so the caller
// keeps whatever it had).
func (s *service) fetchNowPlaying() (nowPlayingResponse, error) {
	track, err := s.currentlyPlaying()
	if err != nil {
		return nowPlayingResponse{}, err
	}
	if track != nil {
		return track.response(true), nil
	}
	// Nothing playing right now, so fall back to play history and let the
	// site say "was listening" rather than going blank.
	return s.recentlyPlayed()
}

// currentlyPlaying returns nil, nil when Spotify reports nothing playing.
func (s *service) currentlyPlaying() (*spotifyTrack, error) {
	res, err := s.spotifyGet(s.nowURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, statusError("currently-playing", res)
	}

	var payload struct {
		IsPlaying bool          `json:"is_playing"`
		Item      *spotifyTrack `json:"item"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if !payload.IsPlaying || payload.Item == nil {
		return nil, nil
	}
	return payload.Item, nil
}

// recentlyPlayed needs the user-read-recently-played scope. A refresh token
// predating that scope makes this fail, which surfaces to the client as the
// last known value — or bare {"playing":false} on a cold cache, the same as
// before this fallback existed.
func (s *service) recentlyPlayed() (nowPlayingResponse, error) {
	res, err := s.spotifyGet(s.recentURL)
	if err != nil {
		return nowPlayingResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nowPlayingResponse{}, statusError("recently-played", res)
	}

	var payload struct {
		Items []struct {
			Track spotifyTrack `json:"track"`
		} `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nowPlayingResponse{}, err
	}
	if len(payload.Items) == 0 {
		// Genuinely nothing to show: a real answer, not a failure.
		return nowPlayingResponse{Playing: false}, nil
	}
	return payload.Items[0].Track.response(false), nil
}

func (s *service) spotifyGet(endpoint string) (*http.Response, error) {
	token, err := s.accessToken()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return s.httpClient.Do(req)
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

	req, err := http.NewRequest(http.MethodPost, s.tokenURL, strings.NewReader(form.Encode()))
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
		return "", statusError("token refresh", res)
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

type apiError struct {
	call   string
	status int
	body   string
}

func (e *apiError) Error() string {
	return "spotify " + e.call + ": status " + strconv.Itoa(e.status) + ": " + e.body
}

func statusError(call string, res *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
	return &apiError{call: call, status: res.StatusCode, body: string(body)}
}
