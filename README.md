# naga.srin.cc

Zola site.

- `zola serve` — dev server
- `zola build` — outputs to `public/`
- `docker build -t naga-site:test . && docker run --rm -p 8080:80 naga-site:test` — test the deploy image

## Deploy

`.github/workflows/deploy.yml` builds/pushes to GHCR on push to `main`,
then deploys `compose.yaml` over SSH. Secrets and VPS setup are documented
in the workflow file's header comment. Two images are built: the site
itself, and `now-playing/` (a separate Go service backing the "now
playing" line — see `now-playing/README.md`).

## Posts

Add `content/blog/<slug>.md` with `title`, `date`, `description`,
`[taxonomies] tags`. Toggle the blog with `extra.blog_enabled` in
`config.toml`.

## Palette / uptime

Two colours per theme (a ground and an ink) as tokens at the top of
`sass/style.scss`, along with the per-theme mark and portraits in
`static/img/`. The theme follows `prefers-color-scheme` until the nav
toggle (`static/theme.js`) saves a choice. Uptime counts from `extra.dob`
in `config.toml`, recomputed client-side by `static/uptime.js`.

## Entrance

Pages come into focus with CSS animations (the "entrances" block in
`sass/style.scss`). They wait until the web fonts are in and, on the
homepage, the song line has its text, so nothing moves mid-animation:
an inline script in `templates/base.html` holds them before first paint
and `static/entrance.js` lets go, after 0.8s at the latest. The song
line reads `extra.now_playing_api` in `config.toml`; remove it to drop
the line.
