# Deploying the GOTTH+ site

This site is a single Go binary with a SQLite database, built into a small image and run behind a
TLS-terminating proxy. The steps below deploy it to one Linux server with Docker. The same image runs on
other platforms; the requirements for those are in [Other platforms](#other-platforms).

**Status of these steps.** The Go parts are checked by the test suite, including browser tests, and the
compose file parses. The image build and the running stack have not been run on a server yet; the first
deploy is the test of this guide. The [troubleshooting](#troubleshooting) section lists what to look for.

## What runs

| Part | Where | Notes |
|---|---|---|
| `migrate` | a one-off container | Brings the SQLite schema up to date. It runs before every start. |
| `app` | the image's `/app/server`, port 8080 | No published port. Only Caddy can reach it. Runs read-only, as a non-root user. |
| `caddy` | `caddy:2`, ports 80 and 443 | Issues and renews the TLS certificate, and forwards requests to `app`. |
| data | the Docker volume `data`, mounted at `/data` | The SQLite database `app.db` and its WAL files. Back this up. |

The server sets `GoEnv=production` in the image. That turns on HSTS, `Secure` cookies, JSON logs and the
strict rate limits, and it has no development-only routes.

## Before you start

- [ ] **A domain name** whose DNS `A` (and `AAAA`, if you use IPv6) record points at the server.
- [ ] **Ports 80 and 443 open** on the server's firewall. Caddy needs port 80 to get the certificate.
- [ ] **The repository committed to git.** The image is built from the files in the repository, and
      `.dockerignore` leaves out `.env` and `app.db`. Commit `static/js/*.min.js` as well: production
      serves those committed copies and never downloads anything.
- [ ] **Tests pass** on the commit you deploy: `make check`.
- [ ] **A backup plan** (see [Backups](#backups)).

## Deploy with Docker and Compose

These steps run on a Linux server with Docker Engine and the Compose plugin installed
(`docker compose version` should print a version). Run them as a user that can use Docker.

**1. Get the code.**

```bash
git clone <your repository URL> gotth-plus-site
cd gotth-plus-site
```

**2. Start the stack.** `SITE_ADDRESS` is the domain name only, without `https://`.

```bash
SITE_ADDRESS=example.com docker compose up -d --build
```

Compose does these things in order: it builds the image, runs `migrate` and waits for it to succeed, then
starts `app`, then starts `caddy`. If `migrate` fails, `app` does not start. Caddy obtains the certificate
when it starts, which can take a few seconds to a minute.

**3. Check it.**

```bash
docker compose ps                                    # app and caddy are running; migrate has exited 0
docker compose logs --tail=50 app                    # JSON lines, "server listening"
curl -fsS https://example.com/healthz                # prints: ok
curl -sI https://example.com/ | grep -i -E 'strict-transport|content-security'
```

The response must carry `Strict-Transport-Security` and a `Content-Security-Policy` header. Open the site in a
browser and try the command builder and a docs vote, which both use the database.

**4. Keep the settings in one place.** The only value you choose is `SITE_ADDRESS`. Put it in a file the
server keeps, so you don't retype it, and export it before each compose command:

```bash
echo 'SITE_ADDRESS=example.com' > /etc/gotth-plus-site.env   # any path you like; keep it out of git
set -a; . /etc/gotth-plus-site.env; set +a
docker compose up -d --build
```

Do not copy the development `.env` to the server. It is for `make dev` only; the image never reads it.

## Updating

```bash
git pull
docker compose build
docker compose run --rm migrate     # the schema, before the new app starts
docker compose up -d                # restarts app and caddy on the new image
```

Note the commit you are running first (`git rev-parse --short HEAD`), so you can return to it.

Schema changes must be additive (new tables, new nullable columns) so the previous version still works
during a rollback. Take a [backup](#backups) before any change that rewrites data.

## Backups

The database is one file in the `data` volume. Back up while the app is stopped, so the WAL file is
checkpointed into `app.db`:

```bash
mkdir -p ~/backups
docker compose stop app
docker run --rm -v gotth-plus-site_data:/data -v ~/backups:/backup alpine \
  sh -c 'cp -a /data/. /backup/app-$(date +%Y%m%d-%H%M%S)/'
docker compose start app
```

The volume name is the project name (`gotth-plus-site`) with `_data`; check it with `docker volume ls`.
The site is small and the stop lasts a few seconds. For a copy taken while the app keeps running, use
SQLite's online backup (`sqlite3 app.db "VACUUM INTO '/backup/app.db'"`, with the `sqlite3` command-line tool
run against the volume's contents). Store backups off the server, and test a restore before you need one.

**Restore:** stop the app, copy the backup's `app.db` (and any `-wal` and `-shm` files) into the volume, then
start the app again:

```bash
docker compose stop app
docker run --rm -v gotth-plus-site_data:/data -v ~/backups/app-20260101-000000:/restore alpine \
  sh -c 'cp -a /restore/. /data/'
docker compose start app
```

## Behind a proxy

The rate limiter counts requests per client address. Behind a proxy, every request comes from the proxy's
address, so every visitor would share one limit. The stack avoids that: `app` gets
`TrustedProxies=172.28.0.0/16`, the subnet of the compose network that only Caddy can reach. For those
requests, the app takes the client address from the `X-Forwarded-For` header that Caddy sets.

If you change the proxy or the network, set `TrustedProxies` to the proxy's address or network, in CIDR form.
Never trust a network that clients can reach directly: a client on a trusted network can forge the header.
An empty `TrustedProxies` is correct when nothing is in front of the app.

## Operating it

- **Logs.** `docker compose logs -f app`. Production logs are JSON, one object per line. Lines written while
  handling a request carry a `request_id`, which matches the `X-Request-Id` response header.
- **Health.** `GET /healthz` returns `ok` when the database answers, and 503 when it does not.
- **Restarts.** Both services restart unless stopped (`restart: unless-stopped`).
- **Secrets.** The site has no secrets and no sign-in. If you add either, set them as environment variables
  or in a secret store. Never bake them into the image.

## Other platforms

The image runs on any platform that can run a container with a persistent volume. The stack has one
requirement that a stateless platform does not: **SQLite needs one instance and a volume that survives
redeploys.** Two instances would each keep their own rate-limit counters and write to separate copies of
the database.

Before you deploy, check your platform's documentation for:

1. **A persistent volume or disk**, mounted at `/data`. The image's default is `DatabasePath=/data/app.db`.
2. **One running instance**, with the volume attached to it.
3. **A migration before the new version starts.** Run `/app/migrate` in the image before the new server
   starts. Confirm how your platform runs a one-off command against the same volume; some platforms
   do not attach volumes to their pre-deploy steps.
4. **HTTPS** from the platform's proxy, and `TrustedProxies` set to the address range of that proxy, so the
   rate limiter sees the real client.
5. **A health check** on `GET /healthz`.
6. **Port 8080** (the image's `Port`), or the platform's port setting.

Read [Behind a proxy](#behind-a-proxy) before you rely on a platform's proxy: find out which address
the proxy connects from, and trust only that range.

## Troubleshooting

| Symptom | Likely cause |
|---|---|
| Every visitor is rate limited, `429 too many requests` | `TrustedProxies` is not set, so the proxy's address is the client for everyone. |
| A form or the vote returns `403 invalid CSRF token` | The site is being opened over plain `http://`. Production cookies are `Secure`, so the browser does not send them. Use HTTPS. |
| `no such table: page_votes` | `migrate` did not run or failed. Check `docker compose logs migrate`. |
| The site loads with no styles | `static/output.css` is built by the image. If you build locally, run `make css`. |
| `app` exits at start with `GoEnv must be "development" or "production"` | The environment has an unknown `GoEnv`. The image sets `production`; check overrides. |
| `TrustedProxies: ... is not an IP address or CIDR network` | A value in `TrustedProxies` is not valid. Use comma-separated IPs or CIDRs. |
| The certificate is not issued | DNS does not point at the server yet, or ports 80 and 443 are blocked. `docker compose logs caddy` says which. |
| Browser shows an HSTS error after a test on another subdomain | HSTS is sent with `includeSubDomains`. Every subdomain must serve HTTPS before you deploy. |

## Rolling back

Check out the commit you noted before the update (`git checkout <commit>`), then run the same three
commands: build, `docker compose run --rm migrate`, and `docker compose up -d`. Because schema changes are
additive, the older version can still read the database. Restore a backup only if data was changed in a way
the older version cannot read.
