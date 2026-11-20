# GOTTH+ site - run `make help` for the list of targets.

PLAYWRIGHT := go run github.com/mxschmitt/playwright-go/cmd/playwright
PORT = $(shell sed -n 's/^Port=//p' .env | tr -d '"')

# The browser libraries committed in static/js. `make js` re-downloads exactly these versions (see README.md);
# never run it as part of a production build.
JS_PINS := HTMX_VERSION=2.0.11 PRELOAD_VERSION=2.1.2 ALPINE_VERSION=3.17.5 HX_OPTIMISTIC_VERSION=1.0.12

.DEFAULT_GOAL := help
.PHONY: help setup js generate css build run dev dev-templ dev-css air migrate fmt vet test e2e-install vuln check clean docker

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  %-14s %s\n", $$1, $$2}'

setup: ## First-time setup: .env, Go modules, the test browser, CSS and a first build
	@test -f .env || (cp .env.example .env && echo "created .env from .env.example")
	go mod download
	$(MAKE) generate css e2e-install
	$(MAKE) migrate

js: ## Re-download the pinned htmx, preload extension, Alpine (CSP) and hx-optimistic into static/js
	$(JS_PINS) ./download-js.sh

generate: ## Regenerate *_templ.go from the .templ files (templ v0.3.1070, the version in go.mod)
	templ generate

css: ## Build the minified Tailwind CSS
	tailwindcss -i static/input.css -o static/output.css --minify

build: generate css ## Build the server and the migrator into ./tmp
	go build -o ./tmp/server ./cmd/server
	go build -o ./tmp/migrate ./cmd/migrate

run: generate css ## Run the server without live reload
	go run ./cmd/server

migrate: ## Bring the database schema up to date (run on every deploy, before the server starts)
	go run ./cmd/migrate

dev: ## Live reload: Tailwind and templ watchers; open the proxy at http://localhost:7331
	$(MAKE) -j2 dev-css dev-templ

dev-templ:
	templ generate --watch --cmd="go run ./cmd/server" --proxy="http://localhost:$(PORT)"

dev-css:
	tailwindcss -i static/input.css -o static/output.css --watch

air: ## Alternative live reload with air (config in air.toml); run `make dev-css` beside it
	air -c air.toml

fmt: ## Format Go and templ files
	gofmt -w cmd internal
	templ fmt .

vet: ## Vet everything, including the browser tests
	go vet -tags e2e ./...

# The e2e build tag adds the browser tests to the unit and integration tests. They need Chromium (`make e2e-install`,
# once); set AXE_CORE_JS to an axe.min.js to include the accessibility scan.
test: generate css ## Unit, integration and browser tests
	go test -race -count=1 -tags e2e ./...

e2e-install: ## Download the browser the browser tests drive (once)
	$(PLAYWRIGHT) install chromium

vuln: ## Check dependencies and the standard library for known vulnerabilities
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: generate fmt vet test ## Everything that must pass before a commit

clean: ## Remove build output
	rm -rf tmp static/output.css

docker: ## Build the production image (Dockerfile)
	docker build -t gotth-plus-site .
