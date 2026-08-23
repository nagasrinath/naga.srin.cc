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

Light/dark tokens in `sass/style.scss`, derived from `static/avatar.png`'s
colors, switching automatically via `prefers-color-scheme`. Uptime counts
from `extra.dob` in
`config.toml`, recomputed client-side by `static/uptime.js`.
