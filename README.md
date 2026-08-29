# naga.srin.cc

Zola site.

- `nix develop` — dev shell with the pinned `zola` and `go`
- `zola serve` — dev server
- `zola build` — outputs to `public/`
- `docker build -t naga-site:test . && docker run --rm -p 8080:80 naga-site:test` — test the deploy image

## Toolchain

`flake.nix` pins the toolchain; **the `zola` version in the dev shell must
match the `zola` tag in `Dockerfile`.** Zola makes breaking `config.toml`
schema changes across minor versions (0.22 moved `[markdown]`'s
`highlight_code`/`highlight_theme` to `[markdown.highlighting]`'s
`code`/`theme`), so a drift between the two means the site builds in one
place and not the other. Bumping either means bumping both, then running
`zola check`.

`now-playing/go.mod` declares a minimum of Go 1.24 to match
`now-playing/Dockerfile`'s builder image; the dev shell's Go may be newer.
Keep the service's own sources buildable by that minimum — tests aren't
copied into the image, so they're free to assume the dev shell's version.

## Deploy

`.github/workflows/deploy.yml` builds/pushes to GHCR on push to `main`,
then deploys `compose.yaml` over SSH. Secrets and VPS setup are documented
in the workflow file's header comment. Two images are built: the site
itself, and `now-playing/` (a separate Go service backing the "now
playing" line — see `now-playing/README.md`).

## Posts

Add `content/blog/<slug>.md` with `title`, `date`, `description`,
`[taxonomies] tags`. `extra.blog_enabled` in `config.toml` only controls
whether the homepage links to the writing list — Zola still builds
`/blog/` and `/tags/` and lists them in `sitemap.xml` and `atom.xml`
regardless, so posts are public either way.

## Palette / uptime

Light/dark tokens in `sass/style.scss`, derived from the avatar's colors,
switching automatically via `prefers-color-scheme`. The avatar itself is
served from `extra.avatar_url`, not the repo. Uptime counts from
`extra.dob` in `config.toml`, recomputed client-side by `static/uptime.js`.

## Client-side JS

All of it is index-only and loaded from `templates/index.html`'s
`{% block scripts %}` — blog, tag, and 404 pages ship none.
`static/term.js` owns the footer prompt (`#typed`): the typewriter intro,
the caret, and the command REPL all touch that one element, so they live
together. `static/uptime.js` and `static/spotify.js` are independent.
