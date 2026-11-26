package docs

// groups is the whole documentation, in navigation order. Facts come from the GOTTH+ template
// (github.com/owenwexler/gotth-plus) and the project generator (github.com/owenwexler/create-gotth-plus-app).
var groups = []Group{
	{
		Title: "Get started",
		Docs: []Doc{
			{
				Slug:    "getting-started",
				Title:   "Getting started",
				Summary: "Create a GOTTH+ project, run it locally and make your first change.",
				Sections: []Section{
					{Title: "What you need", Blocks: []Block{
						UL(
							"Go 1.27.3 or newer (the template's go.mod sets the version).",
							"templ, at the same version as the template's go.mod. Install it with go install github.com/a-h/templ/cmd/templ@<version>.",
							"The Tailwind CSS v4 standalone CLI on your PATH as tailwindcss. No Node and no npm are needed.",
							"git, curl and make. Optional: air, for live reload instead of templ's proxy.",
						),
						P("create-gotth-plus-app checks these tools before it writes anything. A missing or wrong version is a warning with install steps, and the project is still created."),
					}},
					{Title: "Create a project", Blocks: []Block{
						P("The generator asks for anything you don't pass as a flag, shows a summary, and writes nothing if you cancel."),
						Code("go run github.com/owenwexler/create-gotth-plus-app@latest my-app\ncd my-app"),
						P("It then runs make setup for you: Go modules, the JavaScript libraries, Chromium for the browser tests, and a first build. Pass --setup=false to skip that."),
					}},
					{Title: "Run it", Blocks: []Block{
						Code("make dev    # live reload, open http://localhost:7331\nmake run    # no live reload, on the Port in .env (3000 by default)"),
						P("make dev runs the Tailwind watcher and templ's watcher, which starts the server and reloads the browser. make help lists every target."),
					}},
					{Title: "Make your first change", Blocks: []Block{
						P("The splash page is internal/view/pages.templ. Edit it, save, and the browser reloads. Components live beside it in internal/view, and Tailwind picks up any class you write there, as long as it is a complete literal string."),
						P("When you're ready for your own data, read Data: SQLite, Postgres or MySQL, and remove the example pieces the README lists under \"Removing the example app\"."),
					}},
					{Title: "Official sites and documentation", Blocks: []Block{
						P("These docs cover how GOTTH+ uses each part of the stack. For everything else, go to the source:"),
						Links(
							Link{"Go", "https://go.dev/", "the language, downloads and the standard library"},
							Link{"Go documentation", "https://go.dev/doc/", "the tour, Effective Go, the language spec and tutorials"},
							Link{"templ", "https://templ.guide/", "the templ guide: syntax, components and the CLI"},
							Link{"htmx", "https://htmx.org/", "the htmx home page, essays and examples"},
							Link{"htmx documentation", "https://htmx.org/docs/", "attributes, events, configuration and the reference"},
							Link{"htmx preload extension", "https://htmx.org/extensions/preload/", "preloading links and fragments on mousedown or hover"},
							Link{"hx-optimistic", "https://github.com/lorenseanstewart/hx-optimistic", "the optimistic-update extension: its README documents the data-optimistic config"},
							Link{"Alpine.js", "https://alpinejs.dev/", "directives, magics and the global store"},
							Link{"Alpine.js CSP build", "https://alpinejs.dev/advanced/csp", "the build GOTTH+ uses, and what it can't evaluate"},
							Link{"Tailwind CSS", "https://tailwindcss.com/", "the Tailwind home page and the full class reference"},
							Link{"Tailwind CSS standalone CLI", "https://tailwindcss.com/blog/standalone-cli", "the single binary GOTTH+ uses instead of Node and npm"},
							Link{"Tailwind CSS: styling with utility classes", "https://tailwindcss.com/docs/styling-with-utility-classes", "the core concepts: utilities, variants such as hover: and data-*, and theme values"},
							Link{"GORM", "https://gorm.io/", "the Go ORM GOTTH+ uses for every database"},
							Link{"GORM documentation", "https://gorm.io/docs/", "models, queries, associations, migrations and the database drivers"},
							Link{"SQLite", "https://www.sqlite.org/", "the default database: a single file, no server"},
							Link{"SQLite documentation", "https://www.sqlite.org/docs.html", "SQL syntax, pragmas such as WAL mode, and limits"},
							Link{"PostgreSQL", "https://www.postgresql.org/", "the Postgres project, downloads and news"},
							Link{"PostgreSQL documentation", "https://www.postgresql.org/docs/", "the manual for every supported version"},
							Link{"MySQL", "https://www.mysql.com/", "the MySQL home page and downloads"},
							Link{"MySQL documentation", "https://dev.mysql.com/doc/", "the reference manual and connector guides"},
						),
					}},
				},
			},
			{
				Slug:    "cli",
				Title:   "The project generator",
				Summary: "What create-gotth-plus-app asks, what it writes, and every flag it accepts.",
				Sections: []Section{
					{Title: "How it works", Blocks: []Block{
						P("The gotth-plus template is embedded in the create-gotth-plus-app binary, so the generator works offline and each version of the CLI always writes the same project. Database, sign-in and cache choices add or replace files from fixed variants, and the dependency files come from pre-computed go.mod and go.sum sets, so a new project is tidy without network access."),
						P("Every edit to the template text must match exactly once. When the template changes in a way the generator doesn't expect, generation fails and the tests say which edit broke, instead of quietly writing a wrong project."),
					}},
					{Title: "Flags", Blocks: []Block{
						Table([]string{"Flag", "Default", "What it does"},
							[]string{"[name]", "asked", "Project directory. . writes into the current directory."},
							[]string{"--module", "github.com/<user>/<name>", "Go module path."},
							[]string{"--port", "3000", "Development port, written to .env."},
							[]string{"--auth", "none", "none, supabase, oidc or builtin (advanced, experimental). All but none use the app's own server-side sessions. See Sign-in."},
							[]string{"--oidc-preset", "other", "With --auth=oidc: workos, logto, kinde, auth0, zitadel, keycloak or other."},
							[]string{"--auth-methods", "password", "With --auth=builtin: any of password, google, github, apple."},
							[]string{"--db", "sqlite", "sqlite, postgres or mysql. Postgres and MySQL need a server; the tests create a throwaway database on it."},
							[]string{"--sqlite-driver", "purego", "purego (no C compiler, CGO_ENABLED=0 works) or cgo."},
							[]string{"--cache", "none", "Server-side cache: none, redis (Redis or Redis-compatible: Valkey, KeyDB, Dragonfly, managed Redis) or memcached (Memcached or Memcached-compatible). See Server-side cache."},
							[]string{"--e2e", "true", "Playwright browser tests. --e2e=false removes them."},
							[]string{"--live-reload", "templ", "templ (its proxy) or air. Both reload the browser on port 7331."},
							[]string{"--docker", "true", "Production Dockerfile (distroless, non-root) and .dockerignore."},
							[]string{"--ci", "true", "GitHub Actions workflow."},
							[]string{"--js", "true", "Download and commit htmx, its preload extension, Alpine (CSP) and hx-optimistic, with their SHA-384 hashes."},
							[]string{"--axe", "true", "Download axe-core (sha512-verified) for the accessibility test."},
							[]string{"--git", "true", "git init and a first commit."},
							[]string{"--setup", "true", "Run make setup."},
							[]string{"--yes, -y", "", "Ask nothing: flags, then defaults. For scripts and AI agents."},
							[]string{"--force", "", "Write into a directory that is not empty."},
						),
					}},
					{Title: "Sign-in", Blocks: []Block{
						P("The default is no sign-in, which suits learning and prototypes. For production the generator recommends a hosted service: Supabase Auth (--auth=supabase, which moves the database to Postgres) or any OpenID Connect provider (--auth=oidc). The built-in option (--auth=builtin) is labelled experimental and only for advanced use, and the generated SECURITY.md says so. Sign-in covers what each option sets up."),
					}},
					{Title: "Live reload", Blocks: []Block{
						P("make dev reloads the browser at http://localhost:7331 after every change, and runs the Tailwind watcher beside the reloader. --live-reload picks the reloader."),
						UL(
							"templ (the default): templ's own watcher regenerates the views, rebuilds and restarts the server, and its proxy reloads the browser. Nothing else to install.",
							"air: air rebuilds on every .go or .templ change, running templ generate first so the build never uses stale generated files, and its proxy reloads the browser on the same port. Install it with go install github.com/air-verse/air@latest. Its settings are in air.toml.",
						),
					}},
					{Title: "Browser tests", Blocks: []Block{
						P("--e2e (on by default) keeps the Playwright-Go browser tests, and make setup downloads the Chromium they drive. Besides the example app, they check things any app needs: visible focus rings, accessible toasts, no third-party requests, no sideways scrolling on a phone, and pages that load without JavaScript."),
						P("--axe (on by default with --e2e) downloads axe-core from the npm registry, checked against its sha512 hash, for an accessibility scan. Without it, that scan is skipped."),
						P("--e2e=false removes the browser tests, their make targets and the Chromium download, so make test needs no browser. Nothing then catches a CDN script or font sneaking past the Content-Security-Policy, so the README suggests adding the generic checks back from the gotth-plus repository before production."),
					}},
					{Title: "JavaScript libraries", Blocks: []Block{
						P("htmx, its preload extension, Alpine's CSP build and hx-optimistic are plain files in static/js, downloaded by make js. --js (on by default) downloads them during setup, records each version and SHA-384 hash in the README, and commits them, so production serves copies you have pinned and reviewed."),
						P("With --js=false, the files are git-ignored and downloaded later with make js. Commit them before you deploy: the Dockerfile refuses to build without them, and a production build must never download them, because make js fetches the latest versions."),
					}},
					{Title: "Docker and CI", Blocks: []Block{
						P("--docker (on by default) adds a Dockerfile, a .dockerignore and make docker. The image builds the CSS with a pinned, checksum-verified Tailwind CLI and builds the server, then copies both into a distroless image that runs as a non-root user with GoEnv=production on port 8080. On SQLite, the database file lives on a volume at /data. On Postgres and MySQL, pass DatabaseURL at runtime from the platform's secret store: it is never baked into the image. The CGO SQLite driver needs glibc, so that one choice uses a slightly larger distroless base."),
						P("--ci (on by default) adds a GitHub Actions workflow, with actions pinned by commit. It runs go mod verify, checks that templ generate, templ fmt and gofmt change nothing, then runs make vet, make test and make vuln. On Postgres or MySQL it starts a database server for the tests. With --docker, it also checks that the production image builds."),
					}},
					{Title: "Git and setup", Blocks: []Block{
						UL(
							"--git (on by default) runs git init and makes a first commit, so the project starts from a clean, committed state.",
							"--setup (on by default) runs make setup: the Go modules, the JavaScript libraries, Chromium for the browser tests, and a first build. With --setup=false, run make setup yourself.",
							"--yes asks nothing: flags first, then defaults. With --yes, the project name is required.",
						),
					}},
				},
			},
			{
				Slug:    "project-structure",
				Title:   "Project structure",
				Summary: "Where each part of a GOTTH+ project lives, and what it is for.",
				Sections: []Section{
					{Title: "Layout", Blocks: []Block{
						Code(`cmd/server/         entry point: config, database, cache, server start, graceful shutdown
internal/
  app/              routes (newMux) and middleware chain (buildHandler), shared by main and the tests
  auth/             with sign-in: sessions, the signed-in user, providers, per-account limits
  cache/            with --cache: the Redis or Memcached cache and CheckCache
  config/           environment loading and validation
  csrf/             double-submit-cookie CSRF protection
  database/         database connection (GORM: SQLite, Postgres or MySQL)
  dto/              validated request shapes
  escaper/          SQL LIKE wildcard escaping
  factory/          handler constructors (closures over dependencies)
  handler/          struct-based handlers (the health check)
  helper/           generic helpers
  mail/             with password sign-in: emails, by SMTP in production and to the log in development
  middleware/       security headers, rate limit, compression, body limit, admin key
  model/            domain types
  optimistic/       hx-optimistic config builders
  repository/       data access behind interfaces, sentinel errors, pagination
  requestid/        request IDs in headers and logs
  seed/             the dev-only /seed route: creates the tables and loads sample data
  service/          business rules
  styles/           shared Tailwind class strings
  view/             .templ components and their generated *_templ.go
static/
  input.css         Tailwind entry point
  js/app.js         Alpine components and htmx hooks
  js/*.min.js       htmx, its preload extension, Alpine (CSP) and hx-optimistic
  font/  icon/      your self-hosted fonts and icons
download-js.sh      downloads the JavaScript libraries (make js)
seed.sh             runs POST /seed with a CSRF token (make seed)
seedData.json       sample data
Makefile            every command, run make help
Dockerfile          with --docker: the production image`),
					}},
					{Title: "Where to start", Blocks: []Block{
						UL(
							"Pages and components: internal/view. Full pages wrap themselves in Layout; fragments never do.",
							"Routes: newMux in internal/app/app.go. Register new routes there.",
							"Handlers: internal/factory. They parse, validate, call a service or repository, then render.",
							"Client behaviour: static/js/app.js, registered by name, because Alpine runs its CSP build.",
						),
					}},
				},
			},
		},
	},
	{
		Title: "Building",
		Docs: []Doc{
			{
				Slug:    "request-lifecycle",
				Title:   "Request lifecycle",
				Summary: "The middleware chain every request passes through, and the layers a handler uses.",
				Sections: []Section{
					{Title: "Middleware order", Blocks: []Block{
						P("buildHandler in internal/app/app.go wraps the router. Outermost runs first, so cheap rejections happen before any handler work:"),
						UL(
							"Request ID: reads or creates X-Request-ID and puts it on every log line.",
							"Security headers: Content-Security-Policy, nosniff, Referrer-Policy, Permissions-Policy, Cross-Origin-Opener-Policy, and HSTS in production.",
							"Compression: gzip and Brotli.",
							"Rate limit: per client address, with a separate bucket for mutations.",
							"Cross-origin protection: Go's http.CrossOriginProtection, using Sec-Fetch-Site and Origin.",
							"Body limit: 1 MiB on every request.",
							"CSRF: issues the token cookie and checks the header on every non-safe method.",
							"With sign-in: the session and the signed-in user (auth.Middleware), skipped for /static/.",
						),
					}},
					{Title: "Layers", Blocks: []Block{
						P("Dependencies point inward: model, then repository, service, and dto, then factory and handler, then view. A handler parses the request, validates it, calls a service or repository, and renders a fragment or a page. It contains no SQL and no business rules."),
						UL(
							"Repositories are interfaces defined next to their implementation, so handlers can be tested with fakes.",
							"Services contain rules and have no HTTP and no SQL.",
							"Factories are closures over their dependencies, so routes are wired in one place.",
						),
					}},
					{Title: "Errors", Blocks: []Block{
						P("Clients never see internal errors. The handler logs the detail with its request ID and sends a generic message. For htmx requests that is a toast (respondWithToast) with HX-Reswap: none, which leaves the page as it was. Full page loads get the error page."),
					}},
				},
			},
			{
				Slug:    "views",
				Title:   "Views and fragments",
				Summary: "Templ components, the difference between pages and fragments, and how escaping works.",
				Sections: []Section{
					{Title: "Components", Blocks: []Block{
						P("A component is a Go function that returns a templ.Component. Write it in a .templ file, then run templ generate, which writes a *_templ.go file. Never edit the generated file: when a page looks stale, the generated file is almost always out of date."),
						Code("templ Greeting(text string) {\n    <h2 id=\"greeting\">{ text }</h2>\n}"),
						P("Templ escapes text, attributes and URLs by their context. Hostile input such as <script> is shown as text. Avoid templ.Raw, templ.SafeURL and inline scripts that carry user data."),
					}},
					{Title: "Pages and fragments", Blocks: []Block{
						P("A page wraps itself in Layout, which holds the document head, the scripts, the CSRF token for htmx, and the toast container. A fragment is what htmx swaps into the page, so it never renders the layout. Keep a stable id on every swap target."),
						P("Out-of-band toasts are the one place a response changes something other than its target: a toast carries hx-swap-oob=\"beforeend:#toasts\", and the layout's container receives it."),
					}},
				},
			},
			{
				Slug:    "htmx",
				Title:   "htmx",
				Summary: "The htmx version, its extensions, and the error handling this template commits to.",
				Sections: []Section{
					{Title: "Versions", Blocks: []Block{
						P("A project self-hosts four libraries in static/js. make js downloads the latest release of each, and htmx stays on its 2.x line. Pin a version with the variable beside it, for example HTMX_VERSION=2.0.11 make js. Production serves the committed copies and never downloads anything at build time."),
						Table([]string{"File", "Library", "Pin with"},
							[]string{"htmx.min.js", "htmx", "HTMX_VERSION"},
							[]string{"htmx-ext-preload.min.js", "htmx preload extension", "PRELOAD_VERSION"},
							[]string{"alpine-csp.min.js", "Alpine.js (CSP build)", "ALPINE_VERSION"},
							[]string{"hx-optimistic.min.js", "hx-optimistic", "HX_OPTIMISTIC_VERSION"},
						),
						P("download-js.sh prints each file's SHA-384 hash. Compare it with the upstream release when you upgrade, and record it in the README."),
					}},
					{Title: "Configuration", Blocks: []Block{
						P("The layout sets htmx's configuration with a meta tag: selfRequestsOnly on, allowEval off, and includeIndicatorStyles off, so the strict CSP needs no inline styles. The body sets hx-ext=\"optimistic, preload\" and the CSRF token in hx-headers."),
					}},
					{Title: "Errors", Blocks: []Block{
						P("htmx ignores 4xx and 5xx responses by default, so a failure would be silent. app.js changes that: a response with an HX-Reswap header is swapped as the server asked, and network failures or unlabelled errors show a generic toast. Check xhr.status before you reset a form or close a dialog on success."),
					}},
					{Title: "Preloading and boosting", Blocks: []Block{
						P("Add preload to a link or an hx-get element to fetch it on mousedown, or preload=\"mouseover\" to fetch it after a short hover. A preload is a real request, sent with an HX-Preloaded: true header, so only use it for cheap GETs that change nothing. The extension is loaded on every page, but a new project doesn't use it anywhere yet. Good places are navigation, tabs and pagination."),
						Code("<nav hx-boost=\"true\" preload=\"mouseover\">\n  <a href=\"/docs/cli\">The project generator</a>\n</nav>"),
					}},
					{Title: "Optimistic updates", Blocks: []Block{
						P("hx-optimistic makes the page change the moment someone clicks, before the server answers, and undoes the change if the request fails. A toggle such as a done checkbox or a like button feels instant even on a slow connection, while the server stays the source of truth: its response still replaces the target when it arrives. The extension is loaded on every page, and an element opts in with a data-optimistic attribute that holds a JSON config saying what to change on its hx-target."),
						P("For toggles, optimistic.OptimisticFlip builds that config. The target carries its state as a data attribute, and flipping that attribute is the whole optimistic update, so style the target from it with Tailwind's data-[done=true]: variants. If the request fails on the network, the extension rolls the change back. If the server answers with an error toast, app.js restores the old value instead, because htmx then treats the response as a swap rather than an error, and the extension's own rollback never runs. A native checkbox flips itself on click, so app.js also sets a checkbox inside the target back from its data-done attribute. Use done as the key for checkbox toggles, or extend that line in app.js."),
						Code("<li id=\"task-1\" data-done=\"false\" class=\"data-[done=true]:line-through\">\n  <input type=\"checkbox\" hx-post=\"/tasks/1/toggle\" hx-target=\"#task-1\" hx-swap=\"outerHTML\"\n         data-optimistic={ optimistic.OptimisticFlip(\"done\", false) }/>\n  Write the docs\n</li>"),
						P("Only use it where the server almost always agrees, and where a brief flicker back on failure is acceptable. Never use it for anything the user must not see succeed before it really has, such as a payment."),
					}},
				},
			},
			{
				Slug:    "alpine",
				Title:   "Alpine.js without eval",
				Summary: "Client-only state with Alpine's CSP build, and why components are registered by name.",
				Sections: []Section{
					{Title: "Why the CSP build", Blocks: []Block{
						P("The strict Content-Security-Policy forbids 'unsafe-eval', and Alpine's standard build evaluates expressions with it. The CSP build can't evaluate inline objects or statements, so every component is registered in static/js/app.js with Alpine.data and referenced by name."),
						Code("Alpine.data('copyText', () => ({\n    copied: false,\n    async copy() { /* ... */ },\n}))\n\n<div x-data=\"copyText\">\n  <button type=\"button\" @click=\"copy\">Copy</button>\n</div>"),
						P("The CSP build reads property names and method references, not arbitrary expressions. Use a getter for a negation, such as idle, instead of writing !copied in the template."),
					}},
					{Title: "Shared state", Blocks: []Block{
						P("The globalState store is shared by every component, and the htmx hooks in app.js keep its pendingRequests count and online flag current. Read it with $store.globalState."),
					}},
					{Title: "What Alpine is for", Blocks: []Block{
						P("Keep Alpine for state that never needs the server: a copy button, a disclosure, a menu. Anything that changes data goes through htmx, so the server stays the source of truth."),
					}},
				},
			},
			{
				Slug:    "tailwind",
				Title:   "Tailwind CSS v4",
				Summary: "The standalone Tailwind CLI, how it finds classes, and how to add your own font.",
				Sections: []Section{
					{Title: "Building the CSS", Blocks: []Block{
						P("There is no Node and no bundler. The Tailwind v4 standalone CLI builds static/output.css from static/input.css. make css builds it once, and make dev keeps it watching."),
						Code("make css        # minified static/output.css\nmake dev-css    # watch mode (also run by make dev)"),
						P("The production image (--docker) builds the CSS itself, with a pinned, checksum-verified Tailwind CLI, so a stale local output.css never reaches production."),
					}},
					{Title: "Finding classes", Blocks: []Block{
						P("Tailwind only generates classes it can find as complete literal strings. input.css lists the view and styles directories as sources, so a class written in a .templ file or in internal/styles is included. Never build a class by concatenation such as \"bg-\" + color."),
						P("Repeated class lists belong in internal/styles as functions, which keeps a component's markup short and the class strings in one place."),
					}},
					{Title: "Fonts and focus", Blocks: []Block{
						P("To self-host a font, put the .woff2 file in static/font, declare the @font-face in input.css and register it in @theme. Preload it in the layout. Don't use a third-party font CDN: the Content-Security-Policy only allows fonts from the app itself, and self-hosting keeps requests private and fast."),
						P("Keep a visible focus ring. Use the focus-visible: classes alone: adding focus:outline-none beside a focus-visible:outline-* ring hides the ring in Tailwind v4."),
					}},
				},
			},
			{
				Slug:    "data",
				Title:   "Data: SQLite, Postgres or MySQL",
				Summary: "Choosing a database, connecting to it, the migrations that create its tables, and the safe patterns for queries.",
				Sections: []Section{
					{Title: "Choosing a database", Blocks: []Block{
						P("Every project talks to its database through GORM, so models, repositories and queries are the same code whichever database you pick. Choose one with --db when you create the project. Only the connection code, the configuration and the test setup change. With Postgres or MySQL, the generator names the database after the project, with each character that isn't a letter, digit or underscore turned into an underscore, and prints the command that creates it."),
						Table([]string{"--db", "Driver", "Connection", "Good for"},
							[]string{"sqlite (default)", "glebarez/sqlite, pure Go (or mattn/go-sqlite3 with --sqlite-driver=cgo)", "DatabasePath, a file. Defaults to app.db.", "Learning, prototypes and single-server apps. Nothing else to run."},
							[]string{"postgres", "gorm.io/driver/postgres, pure Go", "DatabaseURL, required. For example postgres://localhost:5432/my_app?sslmode=disable", "Most production apps, managed cloud databases, and Supabase Auth, which always uses Postgres."},
							[]string{"mysql", "gorm.io/driver/mysql, pure Go", "DatabaseURL, required, in go-sql-driver/mysql's DSN format. For example root@tcp(localhost:3306)/my_app?parseTime=true&charset=utf8mb4&loc=UTC", "Teams and hosts already running MySQL or a MySQL-compatible server."},
						),
						Code("go run github.com/owenwexler/create-gotth-plus-app@latest my-app --db=postgres\ncreatedb my_app    # Postgres: create the database before the first run\n\ngo run github.com/owenwexler/create-gotth-plus-app@latest my-app --db=mysql\nmysql -u root -e 'CREATE DATABASE my_app CHARACTER SET utf8mb4'"),
					}},
					{Title: "Connecting", Blocks: []Block{
						P("SQLite needs no server. DatabasePath sets the file, and the file is created on first use. In production, put it on a volume that survives redeploys and keep it readable only by the app's user. A new project opens SQLite with GORM's defaults. If reads and writes start waiting on each other, turn on WAL mode and a busy timeout, so a read never blocks a write and a short lock wait is retried instead of failing. With the pure-Go driver, add them to the path in ConnectDb:"),
						Code("dsn := fmt.Sprintf(\"file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)\", path)"),
						P("Postgres and MySQL read DatabaseURL, and the app refuses to start without it. Point it at a local server or a cloud database. internal/database opens a pool of up to 25 connections, well below the server's default connection limit so more than one instance fits, and recycles each connection after 30 minutes, or 5 idle, so failovers and server-side timeouts are picked up."),
						UL(
							"Connect over TLS to any server that isn't localhost: sslmode=verify-full on Postgres, tls=true on MySQL.",
							"Log in as the app's own user with only the privileges it needs, never as a superuser or root.",
							"MySQL needs parseTime=true in the DSN, because the models use time.Time.",
							"IDs are stored as a uuid column on Postgres and char(36) on MySQL. SQLite takes the value as it is.",
						),
					}},
					{Title: "Migrations", Blocks: []Block{
						P("A new project creates its tables in development. make seed calls the /seed route, which runs GORM's AutoMigrate for the example model and then loads seedData.json. It is safe to run more than once. The route only exists in development: GoEnv=production unregisters it, so nothing changes the schema in production until you decide how it should."),
						Code("make seed    # development: create the tables, then load seedData.json"),
						P("Before the first deploy, add an explicit migration step. Run it before each new server starts, never from inside the server, so a schema change happens once and not on every instance. AutoMigrate creates missing tables, columns and indexes, and can change a column's type, but it never drops or renames a column. Review what it will do, or switch to a versioned migration tool such as goose, golang-migrate or Atlas. All three work with SQLite, Postgres and MySQL."),
						P("If AutoMigrate is enough for your schema, a small command is all it takes: a cmd/migrate whose main connects with database.ConnectDb and calls AutoMigrate with every model the app stores, built into the image beside the server and run by each deploy before the new server starts. With sign-in, model.All() already lists every model."),
					}},
					{Title: "Safe queries", Blocks: []Block{
						UL(
							"Use bound parameters: Where(\"slug = ?\", slug). Never concatenate user input into SQL.",
							"Map user-chosen sort columns through an allow-list before they reach Order.",
							"Map a form into a DTO, validate it, then copy the allowed fields. Never bind a request straight into a model: that lets a client set its own id or timestamps.",
							"Escape LIKE wildcards with escaper.Like, and cap the length of any search. On SQLite and Postgres, add ESCAPE '\\' to the query. MySQL already uses \\ as its escape character, so MySQL queries need no ESCAPE clause.",
							"Paginate every list. repository.Paginate clamps the page size to 100.",
						),
					}},
					{Title: "Testing", Blocks: []Block{
						P("On SQLite, each test opens its own temporary file, so the tests need nothing installed. On Postgres and MySQL, the tests need a running server: internal/database/dbtest creates a fresh database for each test and drops it when the test ends. It uses TestDatabaseURL, or DatabaseURL when that is empty, so that user needs CREATEDB on Postgres, or CREATE and DROP on MySQL. Never point either at a production server."),
					}},
					{Title: "Concurrent updates", Blocks: []Block{
						P("Let the database do the arithmetic. Reading a value, changing it in Go and writing it back loses one of two updates that arrive at the same moment. Use a single UPDATE with gorm.Expr, or an upsert with clause.OnConflict, which all three databases support, then read the result back in the same request rather than trusting the number the client already has."),
						Code("db.Model(&model.Post{}).Where(\"id = ?\", id).Update(\"likes\", gorm.Expr(\"likes + ?\", 1))"),
					}},
				},
			},
			{
				Slug:    "auth",
				Title:   "Sign-in",
				Summary: "The four sign-in options, what they set up, and the settings each one needs.",
				Sections: []Section{
					{Title: "Choosing a sign-in option", Blocks: []Block{
						P("Pick one with --auth when you create the project. Every option except none adds the same server-side sessions, and only the code for the methods you choose is written into the project."),
						Table([]string{"--auth", "Where users live", "Good for"},
							[]string{"none (default)", "Nowhere: no users and no sessions. Every route is public.", "Learning, prototypes and public sites."},
							[]string{"supabase", "Supabase Auth: email and password, plus any provider you enable in Supabase. The database is Supabase's Postgres, so --db is set to postgres.", "Production apps that want accounts, MFA and recovery handled by a hosted service."},
							[]string{"oidc", "A hosted OpenID Connect provider: WorkOS, Logto, Kinde, Auth0, Zitadel, Keycloak or any other with discovery.", "Production apps, single sign-on, and teams that already have a provider."},
							[]string{"builtin", "The app's own database: email and password, Google, GitHub and Apple. Experimental and only for advanced use.", "Teams ready to own password storage, email delivery and provider sign-in themselves."},
						),
						Code("go run github.com/owenwexler/create-gotth-plus-app@latest my-app --auth=supabase\ngo run github.com/owenwexler/create-gotth-plus-app@latest my-app --auth=oidc --oidc-preset=workos\ngo run github.com/owenwexler/create-gotth-plus-app@latest my-app --auth=builtin --auth-methods=password,github"),
					}},
					{Title: "What every option sets up", Blocks: []Block{
						UL(
							"Server-side sessions (alexedwards/scs) in a sessions table, so they can be revoked. The cookie is __Host-session in production: HttpOnly, Secure and SameSite=Lax, lasting 7 days at most and 24 hours idle.",
							"Every sign-in issues a new session token and a new CSRF token (csrf.Rotate), so a token planted or seen earlier stops working.",
							"auth.CurrentUser(ctx) is the signed-in user. Wrap routes that need one in auth.RequireUser, and authorise every object by fetching it with both its id and its owner's id.",
							"Sign out, and Sign out everywhere, which raises User.SessionVersion and so ends every session the user has.",
							"?next= after sign-in only accepts paths on the same site (auth.SafeNext), so it can't send anyone elsewhere.",
							"Routes: GET /login, GET /account, POST /logout and POST /logout-everywhere, plus the routes of the chosen methods.",
							"BaseURL, the app's public address. Emailed links and the providers' redirect URLs are built from it, never from the request's Host header, which the client controls. It must be https in production, and defaults to http://localhost:7331 in development.",
						),
						P("In development, a method whose keys are empty is left off the sign-in page, so the project runs before you register anything. In production, the app refuses to start without them."),
					}},
					{Title: "Supabase Auth", Blocks: []Block{
						Table([]string{"Setting", "What it is"},
							[]string{"SupabaseURL", "The project URL, from Project Settings > API."},
							[]string{"SupabasePublishableKey", "The publishable (anon) key. Never the secret or service_role key: it bypasses Supabase's row-level security."},
							[]string{"SupabaseProviders", "Optional. The providers enabled in Supabase that get a Continue with button, such as google,github. Empty means email and password only."},
							[]string{"DatabaseURL", "The project's Postgres, from Connect > Session pooler, with sslmode=require."},
						),
						P("In Supabase, under Authentication > URL Configuration, set the Site URL to BaseURL and allow <BaseURL>/auth/supabase/callback and <BaseURL>/reset-password as redirect URLs."),
						UL(
							"Users and passwords live in Supabase. The app keeps its own users row, linked by Supabase's user ID, and its own session. Supabase's tokens are checked with Supabase, used once and never stored.",
							"Every flow, including provider sign-in and the confirmation and reset links, uses PKCE. The verifier stays in this browser's __Host-supabase-flow cookie, so a reset link only works in the browser that asked for it.",
							"Passwords are 15 to 72 characters. Set Supabase's own minimum no higher.",
							"A user banned or deleted in Supabase stays signed in to the app until their session ends, or until you raise their SessionVersion.",
						),
					}},
					{Title: "OpenID Connect", Blocks: []Block{
						P("--oidc-preset fills .env.example with the right hints for a known provider. The app reads everything else from the issuer's discovery document."),
						Table([]string{"--oidc-preset", "Create the client as", "Issuer looks like"},
							[]string{"workos", "An OAuth application under Connect", "https://<your-subdomain>.authkit.app"},
							[]string{"logto", "A Traditional web application", "https://<tenant-id>.logto.app/oidc"},
							[]string{"kinde", "A back-end web application", "https://<your-business>.kinde.com"},
							[]string{"auth0", "A Regular Web Application", "https://<tenant>.<region>.auth0.com/"},
							[]string{"zitadel", "A web application with code flow", "https://<instance>.zitadel.cloud"},
							[]string{"keycloak", "A confidential OpenID Connect client in the realm", "https://<host>/realms/<realm>"},
							[]string{"other (default)", "Any provider with /.well-known/openid-configuration", "https://issuer.example.com"},
						),
						P("Set OIDCIssuerURL, OIDCClientID and OIDCClientSecret, and optionally OIDCLabel, the button's text after Continue with. At the provider, allow <BaseURL>/auth/oidc/callback as the redirect URI and <BaseURL>/ as the post-logout redirect URI."),
						UL(
							"Sign-in checks a random state, uses PKCE, and checks the ID token's signature, audience, expiry, nonce and issuer, which must be exactly OIDCIssuerURL.",
							"The provider must send email_verified. Accounts are only made for verified emails, and users are matched by the provider's subject, never by email alone.",
							"Signing out also ends the session at the provider when its discovery document has an end_session_endpoint. Some providers, such as WorkOS, don't have one.",
						),
					}},
					{Title: "Built-in sign-in", Blocks: []Block{
						P("Built-in sign-in follows current guidance (OWASP, NIST SP 800-63B-4) and is covered by tests, but it hasn't been independently audited, and keeping passwords, sessions and provider sign-in safe becomes your job. A hosted service also takes care of MFA, breach protection and account recovery. Choose the methods with --auth-methods: any of password, google, github and apple."),
						Table([]string{"Method", "Settings", "Notes"},
							[]string{"password", "SMTPHost, SMTPPort (587), SMTPUsername, SMTPPassword, MailFrom", "Confirmation and reset emails. In development, with SMTPHost empty, emails and their links go to the server log."},
							[]string{"google", "GoogleClientID, GoogleClientSecret", "An OAuth client in the Google Cloud console."},
							[]string{"github", "GitHubClientID, GitHubClientSecret", "An OAuth app in GitHub's developer settings."},
							[]string{"apple", "AppleClientID, AppleTeamID, AppleKeyID, ApplePrivateKeyFile", "A Services ID and a .p8 key. Apple only allows https redirect URLs, so trying it locally needs an https tunnel."},
						),
						P("Each provider's redirect URL is <BaseURL>/auth/<provider>/callback."),
						UL(
							"Passwords are hashed with argon2id, 15 to 256 characters, with no rules about kinds of characters. An unknown email costs the same time as a wrong password, and sign-in, sign-up and reset answer the same whether or not an account exists.",
							"Emailed links carry single-use, expiring tokens stored only as SHA-256. They open a page that changes nothing until its button posts.",
							"Sign-in, sign-up and reset emails are rate-limited per email address as well as per client address.",
							"Provider sign-in checks a random state, uses PKCE where the provider supports it (Apple doesn't), and checks the ID token's nonce for Google and Apple. An email links to an existing account only when the provider has verified it and the account's email is confirmed too, so nobody can take over an account by signing in elsewhere with someone else's address.",
							"Apple comes back with a cross-site POST, so POST /auth/apple/callback alone skips the CSRF and cross-origin checks. Its state and nonce protect it. Keep the .p8 file out of the repository and the image.",
							"make seed also creates dev@example.com, confirmed, with the password in SeedUserPassword, which the generator fills with a random value in .env.",
						),
					}},
					{Title: "Known limits", Blocks: []Block{
						UL(
							"There is no second factor (TOTP or passkeys) yet.",
							"Authorisation is yours: auth.RequireUser only checks that someone is signed in.",
							"The CSRF token is replaced at sign-in and sign-out but not bound to the session.",
							"Password sign-in, built-in or Supabase, needs JavaScript: the forms rely on htmx sending the CSRF header.",
						),
					}},
				},
			},
			{
				Slug:    "caching",
				Title:   "Server-side cache",
				Summary: "The optional Redis or Memcached cache, CheckCache, and how to key, invalidate and secure what you cache.",
				Sections: []Section{
					{Title: "Choosing a cache", Blocks: []Block{
						P("A GOTTH+ project can keep the results of slow or frequent database queries in a cache server, so most requests skip the database. Pick one when you create the project:"),
						Table([]string{"--cache", "Server", "CacheURL"},
							[]string{"none (default)", "No cache. The project has no internal/cache package and no CacheURL.", ""},
							[]string{"redis", "Redis, or anything that speaks its protocol: Valkey, KeyDB, Dragonfly, a managed Redis. Uses go-redis.", "redis://localhost:6379/0, or rediss://user:password@host:6380/0 for TLS"},
							[]string{"memcached", "Memcached, or anything that speaks its text protocol. Uses gomemcache.", "memcached://localhost:11211, or several servers: memcached://host1:11211,host2:11211"},
						),
						Code("go run github.com/owenwexler/create-gotth-plus-app@latest my-app --cache=redis"),
						P("The gotth-plus template itself ships the Redis version. A project whose pages read little from the database doesn't need a cache at all."),
					}},
					{Title: "Turning it on", Blocks: []Block{
						P("The cache is off until CacheURL is set in .env, so a new project runs, and its tests pass, without a cache server. CacheKeyPrefix is optional: every key becomes <prefix>:<key>, which keeps apps or environments apart on a shared server, and changing it drops everything cached at once."),
						Code("CacheURL=redis://localhost:6379/0\nCacheKeyPrefix=myapp"),
						P("main.go connects at startup with cache.Connect and passes the cache to app.NewHandler(env, db, c). Hand it from newMux to the repositories or services that need it. A CacheURL that can't be parsed stops the app at startup; a cache server that doesn't answer is only a warning in the log. With no CacheURL, the cache is cache.Disabled(), which stores nothing."),
						P("For a local server, install Redis, Valkey or Memcached from your package manager, or run one in Docker: docker run --rm -p 6379:6379 valkey/valkey, or docker run --rm -p 127.0.0.1:11211:11211 memcached."),
					}},
					{Title: "CheckCache", Blocks: []Block{
						P("cache.CheckCache is a read-through helper. It returns the value cached under a key. On a miss it calls your function, usually a repository query, stores the result as JSON for the given time, and returns it."),
						Code(`cities, err := cache.CheckCache(ctx, c, "cities", 7*24*time.Hour,
	func(ctx context.Context) ([]model.City, error) { return repo.Cities(ctx) })`),
						UL(
							"An error from your function is returned as it is, so errors.Is still works, and it is never cached.",
							"The cache failing never fails the request. A read error counts as a miss, and a write error is only logged, so with the cache server down every call goes to the database. Short timeouts (500 ms) keep that fallback quick.",
							"Values go through encoding/json, so cache types with exported fields. An entry that no longer decodes, because the type changed, is fetched again and overwritten.",
							"The time to live must be positive. Nothing is cached forever.",
							"Hits and misses are logged at debug level with the key.",
						),
					}},
					{Title: "Keys and invalidation", Blocks: []Block{
						UL(
							"Put everything the value depends on in the key: the page, the filter and, above all, the user or owner for per-user data, such as tasks:owner:<id>:page:2. A key without the owner serves one user's data to another.",
							"Check authorisation before reading the cache, exactly as you would before querying the database.",
							"When the data changes, delete the keys it affects with c.Delete(ctx, keys...), or keep the time to live short enough that stale data doesn't matter.",
							"Read through the cache in repositories or services, never in handlers.",
							"Memcached keys must be at most 250 bytes, prefix included, with no spaces or control characters, and values at most 1 MB by default.",
						),
					}},
					{Title: "Security", Blocks: []Block{
						UL(
							"Never expose the cache server to the internet. Anyone who reaches it can read every entry, plant poisoned ones, or flush it all. Keep it on localhost or a private network, behind a firewall.",
							"Redis: set a password or an ACL user, and use rediss:// (TLS) for anything that isn't on the same host. CacheURL then carries a password, so it belongs in the secret store, and the app never repeats it in an error.",
							"Memcached: the client has no passwords and no TLS, so the network is the only protection. Keep UDP off on the servers (-U 0).",
							"Don't cache secrets such as password hashes, session data or tokens, and don't put secrets in keys, which are logged.",
							"/healthz doesn't check the cache, because the app works without it. Watch the log for cache read failed warnings instead.",
						),
					}},
					{Title: "Testing", Blocks: []Block{
						P("The cache tests need no server: the Redis client is tested against an in-memory Redis (miniredis), and the Memcached client against a small fake Memcached server in memcached_test.go. Both pass the same contract, testContract in internal/cache/cache_test.go, which covers misses, overwrites, expiry, deletes and CheckCache. Integration tests build the app with cache.Disabled()."),
					}},
				},
			},
		},
	},
	{
		Title: "Operating",
		Docs: []Doc{
			{
				Slug:    "security",
				Title:   "Security",
				Summary: "What the template protects against, and the checklist to run before production.",
				Sections: []Section{
					{Title: "What the template does", Blocks: []Block{
						Table([]string{"Protection", "How"},
							[]string{"Content-Security-Policy", "script-src 'self' with no inline script and no 'unsafe-eval'; frame-ancestors 'none'; form-action 'self'; object-src 'none'."},
							[]string{"Other headers", "nosniff, Referrer-Policy, Permissions-Policy, Cross-Origin-Opener-Policy. HSTS with includeSubDomains in production only."},
							[]string{"CSRF", "A 32-byte random token in an HttpOnly, SameSite=Lax cookie, echoed in X-CSRF-Token and compared in constant time on every non-safe request."},
							[]string{"Cross-origin requests", "Go's CrossOriginProtection rejects cross-site browser requests, as defence in depth."},
							[]string{"Rate limits", "Per client, strict in production: 20 requests a second with a burst of 60, and 5 mutations a second with a burst of 20."},
							[]string{"Body size", "1 MiB on every request."},
							[]string{"Timeouts", "Read, write, idle and header timeouts on the server, and graceful shutdown."},
							[]string{"Errors", "Clients get a generic message. The details go to the server log with the request ID."},
							[]string{"Config", "An unknown GoEnv refuses to start instead of acting like development."},
						),
					}},
					{Title: "Behind a proxy", Blocks: []Block{
						P("The rate limiter counts per TCP peer address, and it ignores X-Forwarded-For, because a client can forge that header. Behind a reverse proxy or load balancer every visitor would then share the proxy's one bucket, and one busy user could lock everyone out. Choose one fix:"),
						UL(
							"Rate-limit at the proxy (nginx limit_req, Caddy, your CDN or load balancer), and loosen or remove the in-app limits.",
							"Or trust the forwarded address only from your proxy: add middleware in front of the rate limiter that replaces r.RemoteAddr with the client address from X-Forwarded-For, only when the request really comes from your proxy. Take the right-most address your proxy appended, never the left-most one the client sent.",
						),
					}},
					{Title: "Known limits", Blocks: []Block{
						UL(
							"Without --auth there is no sign-in. Every route is public, and nothing is tied to a user. With sign-in, authorisation is still yours: auth.RequireUser only checks that someone is signed in.",
							"There is no second factor (TOTP or passkeys) yet.",
							"The CSRF token is a double-submit cookie. With sign-in it is replaced at sign-in and sign-out, but it is not bound to the session.",
							"Rate limits live in memory, per process. Several instances each keep their own counters, so use the proxy or a shared store for a global limit.",
							"The app does not terminate TLS. Put a TLS-terminating proxy or platform in front of it.",
						),
					}},
				},
			},
			{
				Slug:    "testing",
				Title:   "Testing",
				Summary: "Unit, integration and browser tests, and the commands that run them.",
				Sections: []Section{
					{Title: "Commands", Blocks: []Block{
						Code("make test     # everything: unit, integration and browser tests\nmake check    # generate, format, vet and test\nmake e2e-install   # once: the Chromium that the browser tests drive"),
					}},
					{Title: "Layers", Blocks: []Block{
						UL(
							"Unit tests sit next to their packages: services, validators, helpers, the CSRF middleware.",
							"Integration tests call the real routes and middleware through httptest on a throwaway database: a temporary SQLite file, or a fresh database on your Postgres or MySQL server.",
							"Browser tests (build tag e2e) drive Chromium with Playwright-Go, one browser for the run and a fresh context per test. They check focus rings, the toast region, that no third-party request is made, that nothing scrolls sideways, that the page loads without JavaScript, and an axe accessibility scan when axe-core is available.",
						),
						P("A plain go test ./... runs only the unit and integration tests and needs no browser. Projects created with --e2e=false have no browser tests at all. Point the generic checks in e2e_generic_test.go at your own pages as the app grows: they only load the home page at first."),
					}},
				},
			},
			{
				Slug:    "deployment",
				Title:   "Deployment",
				Summary: "How to deploy a GOTTH+ project: the image, the database, and what to set before the first deploy.",
				Sections: []Section{
					{Title: "The short version", Blocks: []Block{
						P("make docker builds the production image from the Dockerfile (--docker). It is distroless and runs as a non-root user with GoEnv=production on port 8080. Fly.io, Railway, Render and most container platforms build that Dockerfile as it is."),
						Code("make docker\ndocker run --rm -p 8080:8080 -v my-app-data:/data my-app              # SQLite\ndocker run --rm -p 8080:8080 -e DatabaseURL='...' my-app             # Postgres or MySQL"),
						P("Put a TLS-terminating proxy or platform in front of it, and point its health check at /healthz, which pings the database. Behind a proxy, read Security: Behind a proxy, or every visitor shares one rate-limit bucket."),
					}},
					{Title: "The database", Blocks: []Block{
						P("The image creates no tables: AutoMigrate only runs in the development seed route. Run your reviewed migrations before the first start, and before each deploy that changes the schema. Data: SQLite, Postgres or MySQL describes the options."),
						UL(
							"SQLite: the file is on the /data volume (DatabasePath=/data/app.db). Keep the volume across deploys and back it up.",
							"Postgres or MySQL: pass DatabaseURL from the platform's secret store, connect over TLS, and log in as the app's own user with only the privileges it needs.",
						),
					}},
					{Title: "Before the first deploy", Blocks: []Block{
						UL(
							"GoEnv=production, so HSTS, Secure cookies and JSON logs are on, the strict rate limits apply, and /seed is gone.",
							"Configuration in the platform's environment or secret store, never a committed .env.",
							"The JavaScript libraries committed and pinned. Production never runs download-js.sh.",
							"With sign-in: BaseURL set to the public https:// address, and every setting the chosen option needs. Production refuses to start without them.",
							"With a cache: CacheURL pointing at a server only the app can reach.",
							"Backups of the database, tested by restoring them.",
						),
						P("The generated SECURITY.md has the full checklist."),
					}},
				},
			},
			{
				Slug:    "ai-agents",
				Title:   "Working with AI agents",
				Summary: "The AGENTS.md and CLAUDE.md files, and the checklist an agent follows before it finishes.",
				Sections: []Section{
					{Title: "Guidance files", Blocks: []Block{
						P("Every GOTTH+ project has an AGENTS.md. It holds the conventions for the stack and the security checklist, and its last section lists this project's specifics. CLAUDE.md points to it, so Claude Code and any other agent read the same rules."),
					}},
					{Title: "The checklist", Blocks: []Block{
						UL(
							"templ generate has run, and the generated files are up to date.",
							"gofmt and templ fmt are clean, and go vet and go test -race pass.",
							"New input is validated and length-capped on the server.",
							"No user data reaches templ.Raw, SQL built by concatenation, or an inline script.",
							"Every mutating route uses POST, PUT, PATCH or DELETE, with CSRF and authorisation.",
							"Errors are logged on the server and generic to the client.",
							"No secrets are committed, and .env.example covers any new configuration.",
						),
					}},
				},
			},
		},
	},
}
