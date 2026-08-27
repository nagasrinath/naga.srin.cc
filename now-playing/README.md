# now-playing

Small Go service backing the "now playing" line on naga.srin.cc.
Polls the site owner's Spotify account and exposes one endpoint:

`GET /now-playing` → `{"playing":true,"track":"...","artist":"..."}` when
something's currently playing; otherwise falls back to Spotify's play
history, returning `{"playing":false,"track":"...","artist":"..."}` for
the last-played track, or bare `{"playing":false}` if even that's
unavailable.

- `go build .` / `go run .` — needs `SPOTIFY_CLIENT_ID`,
  `SPOTIFY_CLIENT_SECRET`, `SPOTIFY_REFRESH_TOKEN` set; `ALLOWED_ORIGIN`
  and `PORT` are optional (default to `https://naga.srin.cc` and `8080`).
- `docker build -t now-playing:test .` — same env vars via `-e`/`--env-file`.

Deployed as a second service in the repo root's `compose.yaml`, built
and pushed to GHCR by `.github/workflows/deploy.yml` alongside the main
site. See that workflow's header comment for the required secrets and
one-time VPS/DNS setup (routing `api.srin.cc` to this container).

## Getting a Spotify refresh token (one-time)

1. Create an app at the [Spotify developer dashboard](https://developer.spotify.com/dashboard).
   Development Mode is fine — this only ever authorizes your own
   account, no quota extension needed. Add a redirect URI of
   `https://127.0.0.1:8888/callback` (it doesn't need to actually be
   running anything — you're just reading the `code` param back out of
   the browser's address bar after Spotify redirects there).
2. Note the app's **Client ID** and **Client Secret**.
3. Open this in a browser, filling in your client ID (needs both
   `user-read-currently-playing` and `user-read-recently-played`, the
   latter for the "was listening" fallback when nothing's playing):

   ```
   https://accounts.spotify.com/authorize?client_id=YOUR_CLIENT_ID&response_type=code&redirect_uri=https://127.0.0.1:8888/callback&scope=user-read-currently-playing%20user-read-recently-played
   ```

   Log in, approve, and the browser will try to load
   `https://127.0.0.1:8888/callback` — since nothing's listening there
   (and there's no valid cert for it), you'll get a connection or
   certificate warning page. That's expected; you don't need it to
   load. Copy the `code=...` value out of the address bar.
4. Exchange that code for a refresh token:

   ```sh
   curl -X POST https://accounts.spotify.com/api/token \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d grant_type=authorization_code \
     -d code=PASTE_THE_CODE_HERE \
     -d redirect_uri=https://127.0.0.1:8888/callback \
     -u YOUR_CLIENT_ID:YOUR_CLIENT_SECRET
   ```

   The JSON response's `refresh_token` field is what goes in the
   `SPOTIFY_REFRESH_TOKEN` GitHub Environment secret. It doesn't expire
   on its own — this is a one-time step unless you revoke the app's
   access from your Spotify account.
