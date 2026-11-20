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
						P("When you're ready for your own data, read Data: SQLite and GORM, and remove the example pieces the README lists under \"Removing the example app\"."),
					}},
				},
			},
			{
				Slug:    "cli",
				Title:   "The project generator",
				Summary: "What create-gotth-plus-app asks, what it writes, and every flag it accepts.",
				Sections: []Section{
					{Title: "How it works", Blocks: []Block{
						P("The gotth-plus template is embedded in the create-gotth-plus-app binary, so the generator works offline and each version of the CLI always writes the same project. Database and sign-in choices add or replace files from fixed variants, and the dependency files come from pre-computed go.mod and go.sum sets, so a new project is tidy without network access."),
						P("Every edit to the template text must match exactly once. When the template changes in a way the generator doesn't expect, generation fails and the tests say which edit broke, instead of quietly writing a wrong project."),
					}},
					{Title: "Flags", Blocks: []Block{
						Table([]string{"Flag", "Default", "What it does"},
							[]string{"[name]", "asked", "Project directory. . writes into the current directory."},
							[]string{"--module", "github.com/<user>/<name>", "Go module path."},
							[]string{"--port", "3000", "Development port, written to .env."},
							[]string{"--auth", "none", "none, supabase, oidc or builtin (advanced, experimental). All but none use the app's own server-side sessions."},
							[]string{"--oidc-preset", "other", "With --auth=oidc: workos, logto, kinde, auth0, zitadel, keycloak or other."},
							[]string{"--auth-methods", "password", "With --auth=builtin: any of password, google, github, apple."},
							[]string{"--db", "sqlite", "sqlite, postgres or mysql. Postgres and MySQL need a server; the tests create a throwaway database on it."},
							[]string{"--sqlite-driver", "purego", "purego (no C compiler, CGO_ENABLED=0 works) or cgo."},
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
						P("The default is no sign-in, which suits learning and prototypes. For production the generator recommends a hosted service: Supabase Auth (--auth=supabase, which moves the database to Postgres) or any OpenID Connect provider (--auth=oidc). The built-in option (--auth=builtin) is labelled advanced and experimental, and the generated SECURITY.md says so. This site uses no sign-in."),
					}},
				},
			},
			{
				Slug:    "project-structure",
				Title:   "Project structure",
				Summary: "Where each part of a GOTTH+ project lives, and what it is for.",
				Sections: []Section{
					{Title: "Layout", Blocks: []Block{
						Code(`cmd/server/         entry point: config, database, server start, graceful shutdown
cmd/migrate/        schema migrations, run on each deploy
internal/
  app/              routes (newMux) and middleware chain (buildHandler)
  config/           environment loading and validation
  csrf/             double-submit-cookie CSRF protection
  database/         SQLite connection (GORM) and migrations
  docs/             this site's documentation content
  dto/              validated request shapes
  factory/          handler constructors (closures over dependencies)
  middleware/       security headers, rate limit, compression, body limit, trusted proxies
  model/            domain types
  repository/       data access behind interfaces
  requestid/        request IDs in headers and logs
  service/          business rules
  styles/           shared Tailwind class strings
  view/             .templ components and their generated *_templ.go
static/
  input.css         Tailwind entry point
  js/app.js         Alpine components and htmx hooks
  js/*.min.js       htmx, its preload extension, Alpine (CSP) and hx-optimistic
  icon/  font/      your icons and fonts
Makefile            every command, run make help
Dockerfile          the production image`),
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
							"Trusted proxies (only when TrustedProxies is set): replaces the client address when the peer is your proxy.",
							"Request ID: reads or creates X-Request-ID and puts it on every log line.",
							"Security headers: Content-Security-Policy, nosniff, Referrer-Policy, Permissions-Policy, Cross-Origin-Opener-Policy, and HSTS in production.",
							"Compression: gzip and Brotli.",
							"Rate limit: per client address, with a separate bucket for mutations.",
							"Cross-origin protection: Go's http.CrossOriginProtection, using Sec-Fetch-Site and Origin.",
							"Body limit: 1 MiB on every request.",
							"CSRF: issues the token cookie and checks the header on every non-safe method.",
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
						P("The template self-hosts four pinned libraries in static/js. Production serves those committed copies and never downloads anything at build time."),
						Table([]string{"File", "Library", "Version"},
							[]string{"htmx.min.js", "htmx", "2.0.11"},
							[]string{"htmx-ext-preload.min.js", "htmx preload extension", "2.1.2"},
							[]string{"alpine-csp.min.js", "Alpine.js (CSP build)", "3.17.5"},
							[]string{"hx-optimistic.min.js", "hx-optimistic", "1.0.12"},
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
						P("Add preload to a link or an hx-get element to fetch it on mousedown, or preload=\"mouseover\" to fetch it after a short hover. A preload is a real request, so only use it for cheap GETs that change nothing. This site's navigation uses hx-boost with preload=\"mouseover\"."),
						Code("<nav hx-boost=\"true\" preload=\"mouseover\">\n  <a href=\"/docs/cli\">The project generator</a>\n</nav>"),
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
						P("The CSS is built in the Docker image as well, so the committed output.css is not used in production builds."),
					}},
					{Title: "Finding classes", Blocks: []Block{
						P("Tailwind only generates classes it can find as complete literal strings. input.css lists the view and styles directories as sources, so a class written in a .templ file or in internal/styles is included. Never build a class by concatenation such as \"bg-\" + color."),
						P("Repeated class lists belong in internal/styles as functions, which keeps a component's markup short and the class strings in one place."),
					}},
					{Title: "Fonts and focus", Blocks: []Block{
						P("To self-host a font, put the .woff2 file in static/font, declare the @font-face in input.css and register it in @theme. Preload it in the layout. The site uses no third-party font CDN, which keeps requests private and fast."),
						P("Keep a visible focus ring. Use the focus-visible: classes alone: adding focus:outline-none beside a focus-visible:outline-* ring hides the ring in Tailwind v4."),
					}},
				},
			},
			{
				Slug:    "data",
				Title:   "Data: SQLite and GORM",
				Summary: "The database layer, the migrations that create its tables, and the safe patterns for queries.",
				Sections: []Section{
					{Title: "Connecting", Blocks: []Block{
						P("The database is SQLite through GORM, using the pure-Go glebarez/sqlite driver, so the build needs no C compiler. DatabasePath sets the file, and it defaults to app.db. In production, put the file on a volume that survives redeploys."),
						P("This site opens the database in WAL mode with a busy timeout, so a read never blocks a write and a short lock wait is retried instead of failing."),
					}},
					{Title: "Migrations", Blocks: []Block{
						P("The server never changes the schema. cmd/migrate does, and you run it before each new server starts. database.Migrate lists every model the app stores, so adding a model means adding it there and running the migrator."),
						Code("make migrate    # brings the schema up to date (the same binary runs in the image)"),
					}},
					{Title: "Safe queries", Blocks: []Block{
						UL(
							"Use bound parameters: Where(\"slug = ?\", slug). Never concatenate user input into SQL.",
							"Map user-chosen sort columns through an allow-list before they reach Order.",
							"Map a form into a DTO, validate it, then copy the allowed fields. Never bind a request straight into a model: that lets a client set its own id or timestamps.",
							"Escape LIKE wildcards with escaper.Like, and cap the length of any search.",
							"Paginate every list. repository.Paginate clamps the page size to 100.",
						),
					}},
					{Title: "Counters", Blocks: []Block{
						P("The vote on each docs page is a single upsert: the database adds one to the counter, so two votes at the same moment both count. Read the counts back in the same request, rather than trusting the number the client already has."),
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
						P("The rate limiter counts per TCP peer address, and it ignores X-Forwarded-For, because a client can forge that header. Behind a reverse proxy every visitor would then share the proxy's one bucket."),
						P("Set TrustedProxies to the proxy's address or network, in CIDR form. The server then takes the client address from X-Forwarded-For, but only when the request really comes from that proxy, and it uses the right-most entry the proxy appended, never the left-most one the client sent."),
						Code("TrustedProxies=127.0.0.1/32"),
					}},
					{Title: "Known limits", Blocks: []Block{
						UL(
							"There is no sign-in. Every route is public, and nothing is tied to a user.",
							"Rate limits live in memory, per process. Several instances each keep their own counters.",
							"Votes are not deduplicated by visitor. The rate limit slows repeated voting, but it does not stop it.",
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
							"Integration tests call the real routes and middleware through httptest on a throwaway SQLite file.",
							"Browser tests (build tag e2e) drive Chromium with Playwright-Go, one browser for the run and a fresh context per test. They check focus rings, the toast region, that no third-party request is made, that nothing scrolls sideways, that the page loads without JavaScript, and an axe accessibility scan when axe-core is available.",
						),
					}},
				},
			},
			{
				Slug:    "deployment",
				Title:   "Deployment",
				Summary: "How this site is deployed, and where the step-by-step guide lives.",
				Sections: []Section{
					{Title: "The short version", Blocks: []Block{
						P("The image is built by the Dockerfile in the repository. It is distroless and runs as a non-root user with GoEnv=production on port 8080, and the SQLite file lives on a volume at /data. Each deploy runs the migrator, then starts the server. Put a TLS-terminating proxy or platform in front of it, point its health check at /healthz, and set TrustedProxies to that proxy."),
						P("The full guide, with commands for a server running Docker behind Caddy, and notes for hosted platforms, is DEPLOYMENT.md in the repository."),
					}},
					{Title: "Before the first deploy", Blocks: []Block{
						UL(
							"GoEnv=production, so HSTS, Secure cookies and JSON logs are on and the strict limits apply.",
							"Configuration in the platform's environment or secret store, never a committed .env.",
							"The JavaScript libraries committed and pinned. Production never runs download-js.sh.",
							"Backups for the database volume, tested by restoring them.",
						),
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
