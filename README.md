# GOTTH+ site

The marketing and documentation site for [GOTTH+](https://github.com/owenwexler/gotth-plus), built with GOTTH+
and deployed the way its documentation describes. It is started from
[create-gotth-plus-app](https://github.com/owenwexler/create-gotth-plus-app).

Go serves every page, Templ renders them, htmx and a little Alpine.js make them interactive, Tailwind styles
them, and SQLite stores the feedback on the docs pages. There is no Node, no npm and no bundler.

## Run it

```bash
make setup    # .env, Go modules, the browser for the tests, CSS, and the database schema
make dev      # live reload; open http://localhost:7331
```

`make run` starts the server without live reload, on the `Port` in `.env` (3000 by default).
`make help` lists every target.

Requirements: Go 1.27.3 or newer, templ v0.3.1070 (the version in `go.mod`), the Tailwind CSS v4 standalone CLI
on your `PATH` as `tailwindcss`, `git`, `curl` and `make`.

| Command | What it does |
|---|---|
| `make dev` | Tailwind and templ watchers; templ starts the server and reloads the browser |
| `make run` | Build the CSS and the templ code, then run the server |
| `make migrate` | Bring the database schema up to date (run on every deploy) |
| `make test` | Unit, integration and browser tests |
| `make check` | Generate, format, vet and test |
| `make build` | Build `tmp/server` and `tmp/migrate` |
| `make docker` | Build the production image |
| `make js` | Re-download the pinned browser libraries (see below) |

## Configuration

Configuration comes from the environment. `.env` is for local development only and is git-ignored.

| Variable | Default | Meaning |
|---|---|---|
| `GoEnv` | required | `development` or `production`. Production turns on HSTS, `Secure` cookies, JSON logs and the strict rate limits. |
| `Port` | required | The port the server listens on. |
| `DatabasePath` | `app.db` | The SQLite file. |
| `TrustedProxies` | empty | Comma-separated networks of the reverse proxies in front of the app, such as `127.0.0.1/32`. Leave it empty when nothing is in front of the app. |

[`DEPLOYMENT.md`](DEPLOYMENT.md) explains `TrustedProxies` and why it matters behind a proxy.

## Pages and routes

| Method | Path | Purpose |
|---|---|---|
| GET | `/` | Landing page, with the project-command builder |
| POST | `/try/command` | htmx: validates a project name and database and returns the command |
| GET | `/docs` | Documentation index |
| GET | `/docs/{slug}` | A documentation page, with its vote counts |
| POST | `/docs/{slug}/vote?vote=yes\|no` | htmx: records a "was this page helpful?" answer |
| GET | `/healthz` | Health check: `ok` when the database answers, 503 otherwise |
| GET | `/static/*` | Static assets (no directory listings) |

The documentation is Go data in `internal/docs/content.go`, rendered by Templ, so every page is escaped like
the rest of the site and no raw HTML is written.

## Structure

```
cmd/server/         entry point: config, database, server start, graceful shutdown
cmd/migrate/        schema migrations (run on every deploy)
internal/
  app/              routes (newMux) and the middleware chain (buildHandler)
  config/           environment loading and validation, TrustedProxies
  csrf/             double-submit-cookie CSRF protection
  database/         SQLite connection (WAL, busy timeout) and Migrate
  docs/             the documentation content
  dto/              validated request shapes
  factory/          handlers (closures over their dependencies)
  middleware/       security headers, rate limit, trusted proxies, compression, body limit
  model/            PageVote
  repository/       VoteRepository: counts and atomic upserts
  service/          the project-command rule
  styles/           Tailwind class strings
  view/             Templ components and their generated *_templ.go
static/
  input.css         Tailwind entry point (and the docs menu rules)
  js/app.js         Alpine components and htmx hooks
  js/*.min.js       the pinned browser libraries
  icon/favicon.svg  the favicon
deploy/Caddyfile    TLS and the reverse proxy for the compose stack
docker-compose.yml  migrate, app and caddy, for one server
Dockerfile          the production image (distroless, non-root)
```

## Browser libraries

htmx, its `preload` extension, the Alpine.js CSP build and `hx-optimistic` are committed in `static/js`, so
production serves the copies reviewed here. Each file's SHA-384 hash:

| File | Library | Version | SHA-384 |
|---|---|---|---|
| `htmx.min.js` | htmx | 2.0.11 | `sha384-2OatzQy1H+Zd/IIrjr1TcuDGqLXeHhbooAyJY1KdQMKnr4LZ22k31GBLdYKHmVjg` |
| `htmx-ext-preload.min.js` | htmx preload extension | 2.1.2 | `sha384-PRIcY6hH1Y5784C76/Y8SqLyTanY9rnI3B8F3+hKZFNED55hsEqMJyqWhp95lgfk` |
| `alpine-csp.min.js` | Alpine.js (CSP build) | 3.17.5 | `sha384-Q4UOujiwqAb/jTt3ki4pdkETk4I3rIFQX8h271PLI41okj/YIKY7IQQFMw73wkrX` |
| `hx-optimistic.min.js` | hx-optimistic | 1.0.12 | `sha384-SJqi53iK1nRHjsBTLO3JMiFcQEoyxpnW+MH3kP6R44JMCgvGToCQ6kEkL/q6UvMo` |

`make js` downloads exactly these versions. To upgrade, read each library's release notes, change the pins in
the Makefile (`JS_PINS`), run `make js`, compare the printed hashes with the upstream release, and update this
table in the same commit.

## Testing

```bash
make e2e-install   # once: the Chromium that the browser tests drive (make setup does this too)
make test          # unit, integration and browser tests
```

- Unit tests sit next to their packages.
- `cmd/server/integration_test.go` calls the real routes and middleware on a throwaway SQLite file.
- `cmd/server/e2e_test.go` drives Chromium: the command builder, boosted navigation, the vote, and the phone menu.
- `cmd/server/e2e_generic_test.go` holds checks for any app: focus rings, the toast live region, no third-party
  requests, no horizontal scroll, loading without JavaScript. An axe accessibility scan runs when an
  `axe.min.js` is at `cmd/server/testdata/axe.min.js` or `AXE_CORE_JS` points at one, and is skipped otherwise.

## Deploying

[`DEPLOYMENT.md`](DEPLOYMENT.md) covers the deploy on one server with Docker Compose and Caddy (HTTPS), updates,
backups, rollbacks, the requirements for other platforms, and troubleshooting.

## Guidance for AI agents

[`AGENTS.md`](AGENTS.md) holds the conventions for this project. `CLAUDE.md` points to it.
