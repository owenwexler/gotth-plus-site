package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/owenwexler/gotth-plus-site/internal/config"
	"github.com/owenwexler/gotth-plus-site/internal/csrf"
	"github.com/owenwexler/gotth-plus-site/internal/docs"
)

// client talks to the app the way a browser would: it keeps the CSRF cookie from its first GET and
// echoes it in the X-CSRF-Token header, like the hx-headers attribute on <body> does.
type client struct {
	t       *testing.T
	handler http.Handler
	token   string
}

func newClient(t *testing.T, handler http.Handler) *client {
	t.Helper()
	c := &client{t: t, handler: handler}

	rec := c.do("GET", "/", nil, nil)
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == csrf.CookieName {
			c.token = cookie.Value
		}
	}
	if c.token == "" {
		t.Fatal("GET / did not set the CSRF cookie")
	}
	return c
}

func (c *client) do(method, path string, form url.Values, headers map[string]string) *httptest.ResponseRecorder {
	c.t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if c.token != "" {
		req.AddCookie(&http.Cookie{Name: csrf.CookieName, Value: c.token})
		req.Header.Set(csrf.HeaderName, c.token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	c.handler.ServeHTTP(rec, req)
	return rec
}

// htmx issues the request with HX-Request set; these helpers do the same.
func (c *client) htmx(method, path string, form url.Values) *httptest.ResponseRecorder {
	return c.do(method, path, form, map[string]string{"HX-Request": "true"})
}

func TestHome(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Errorf("Content-Type = %q", got)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"<!doctype html>",
		"<title>GOTTH+ · Go, Templ, Tailwind and HTMX</title>",
		`name="description"`,
		`href="/static/icon/favicon.svg"`,
		`id="command-result"`,
		`hx-post="/try/command"`,
		`src="/static/js/htmx-ext-preload.min.js"`,
		`hx-ext="optimistic, preload"`,
		`href="/docs/getting-started"`,
		`id="toasts"`,
		csrf.HeaderName,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page is missing %q", want)
		}
	}

	for _, h := range []string{"Content-Security-Policy", "X-Content-Type-Options", "Referrer-Policy", "X-Request-ID"} {
		if rec.Header().Get(h) == "" {
			t.Errorf("missing %s header", h)
		}
	}
	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS must not be sent in development")
	}

	// the CSRF token rendered into the page is the one in the cookie
	var token string
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == csrf.CookieName {
			token = cookie.Value
			if !cookie.HttpOnly {
				t.Error("CSRF cookie must be HttpOnly")
			}
		}
	}
	if token == "" || !strings.Contains(body, token) {
		t.Error("the CSRF token in the cookie must be rendered into hx-headers")
	}
}

func TestDocsPages(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	index := c.do("GET", "/docs", nil, nil)
	if index.Code != http.StatusOK {
		t.Fatalf("GET /docs: status %d", index.Code)
	}

	for _, d := range docs.All() {
		if !strings.Contains(index.Body.String(), `href="/docs/`+d.Slug+`"`) {
			t.Errorf("/docs does not link to %s", d.Slug)
		}

		rec := c.do("GET", "/docs/"+d.Slug, nil, nil)
		if rec.Code != http.StatusOK {
			t.Errorf("GET /docs/%s: status %d", d.Slug, rec.Code)
			continue
		}
		body := rec.Body.String()
		if !strings.Contains(body, "<h1") || !strings.Contains(body, d.Title) {
			t.Errorf("/docs/%s does not show its title %q", d.Slug, d.Title)
		}
		if !strings.Contains(body, `id="docs-nav"`) {
			t.Errorf("/docs/%s has no navigation", d.Slug)
		}
		if !strings.Contains(body, `hx-boost="true"`) || !strings.Contains(body, `preload="mouseover"`) {
			t.Errorf("/docs/%s navigation is not boosted and preloaded", d.Slug)
		}
	}
}

func TestUnknownPagesAre404(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	for _, path := range []string{"/nope", "/docs/nope", "/docs/", "/.env", "/app.db", "/go.mod"} {
		if got := c.do("GET", path, nil, nil).Code; got != http.StatusNotFound && got != http.StatusMethodNotAllowed {
			t.Errorf("GET %s: status %d, want 404 or 405", path, got)
		}
	}

	// the mux cleans traversal paths (a redirect), so the check is that no file is served
	for _, path := range []string{"/static/", "/static/../go.mod", "/static/%2e%2e/go.mod"} {
		if got := c.do("GET", path, nil, nil).Code; got == http.StatusOK {
			t.Errorf("GET %s: status 200, want it refused", path)
		}
	}

	// the 404 page is a page, not an empty response
	rec := c.do("GET", "/docs/nope", nil, nil)
	if !strings.Contains(rec.Body.String(), "This page doesn’t exist.") {
		t.Errorf("404 page body: %q", rec.Body.String())
	}
}

// Every page's script tags load from /static/js. A bare <script> would be blocked by the CSP.
var inlineScript = regexp.MustCompile(`<script(\s[^>]*)?>`)

func TestNoInlineScripts(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	paths := []string{"/", "/docs", "/docs/deployment"}
	for _, d := range docs.All() {
		paths = append(paths, "/docs/"+d.Slug)
	}

	for _, path := range paths {
		body := c.do("GET", path, nil, nil).Body.String()
		for _, tag := range inlineScript.FindAllString(body, -1) {
			if !strings.Contains(tag, "src=") {
				t.Errorf("%s has an inline script: %s", path, tag)
			}
		}
		if strings.Contains(body, "<style") {
			t.Errorf("%s has an inline <style>, which the CSP forbids", path)
		}
	}
}

func TestStaticFiles(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	rec := c.do("GET", "/static/js/app.js", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /static/js/app.js: status %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "javascript") {
		t.Errorf("app.js Content-Type = %q", got)
	}

	if rec := c.do("GET", "/static/icon/favicon.svg", nil, nil); rec.Code != http.StatusOK {
		t.Errorf("favicon: status %d", rec.Code)
	}

	// no directory listings, and nothing outside static/
	for _, path := range []string{"/static/js/", "/static/%2e%2e/go.mod"} {
		if got := c.do("GET", path, nil, nil).Code; got == http.StatusOK {
			t.Errorf("GET %s: status 200, want it refused", path)
		}
	}
}

func TestHealth(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	rec := c.do("GET", "/healthz", nil, nil)
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Errorf("GET /healthz: status %d, body %q", rec.Code, rec.Body.String())
	}
}

func TestVote(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)
	path := "/docs/cli"

	if body := c.do("GET", path, nil, nil).Body.String(); !strings.Contains(body, "Be the first to answer.") {
		t.Fatalf("a page nobody has voted on should say so; body has no such text")
	}

	t.Run("a helpful vote is counted and shown", func(t *testing.T) {
		rec := c.htmx("POST", "/docs/cli/vote?vote=yes", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, `id="vote"`) || !strings.Contains(body, "1 of 1 readers found this page helpful.") {
			t.Errorf("unexpected fragment: %s", body)
		}
		if strings.Contains(strings.ToLower(body), "<html") {
			t.Error("a fragment must not render the layout")
		}
	})

	t.Run("a second answer is added, and the page shows the totals", func(t *testing.T) {
		if body := c.htmx("POST", "/docs/cli/vote?vote=no", nil).Body.String(); !strings.Contains(body, "1 of 2 readers found this page helpful.") {
			t.Errorf("after a no vote: %s", body)
		}
		if body := c.do("GET", path, nil, nil).Body.String(); !strings.Contains(body, "1 of 2 readers found this page helpful.") {
			t.Error("the totals did not survive a reload")
		}
	})

	t.Run("a vote must be yes or no", func(t *testing.T) {
		rec := c.htmx("POST", "/docs/cli/vote?vote=maybe", nil)
		if rec.Code != http.StatusUnprocessableEntity || rec.Header().Get("HX-Reswap") != "none" {
			t.Errorf("status %d, HX-Reswap %q; want a 422 toast that leaves the page alone", rec.Code, rec.Header().Get("HX-Reswap"))
		}
		if !strings.Contains(rec.Body.String(), "Could not record your vote") {
			t.Errorf("no error toast: %s", rec.Body.String())
		}
	})

	t.Run("only our pages can be voted on", func(t *testing.T) {
		if got := c.htmx("POST", "/docs/made-up/vote?vote=yes", nil).Code; got != http.StatusNotFound {
			t.Errorf("status %d, want 404", got)
		}
	})

	t.Run("a vote needs the CSRF token", func(t *testing.T) {
		anonymous := &client{t: t, handler: handler}
		if got := anonymous.do("POST", "/docs/cli/vote?vote=yes", nil, nil).Code; got != http.StatusForbidden {
			t.Errorf("vote without CSRF token: status %d, want 403", got)
		}
	})
}

func TestCommandBuilder(t *testing.T) {
	handler, _ := newTestApp(t, devEnv)
	c := newClient(t, handler)

	t.Run("builds the command", func(t *testing.T) {
		rec := c.htmx("POST", "/try/command", url.Values{"name": {" my-app "}, "database": {"postgres"}})
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		want := "go run github.com/owenwexler/create-gotth-plus-app@latest my-app --db=postgres --yes"
		if !strings.Contains(body, want) {
			t.Errorf("command missing from %s", body)
		}
		if !strings.Contains(body, `x-data="copyText"`) {
			t.Error("the result must offer the copy button")
		}
	})

	t.Run("rejects a database we do not generate", func(t *testing.T) {
		rec := c.htmx("POST", "/try/command", url.Values{"name": {"my-app"}, "database": {"oracle"}})
		if rec.Code != http.StatusUnprocessableEntity || rec.Header().Get("HX-Reswap") != "none" {
			t.Errorf("status %d; want a 422 toast", rec.Code)
		}
	})

	t.Run("rejects names that are not directory-safe", func(t *testing.T) {
		for _, name := range []string{"", "My App", "../etc", "<script>alert(1)</script>", "a;rm"} {
			rec := c.htmx("POST", "/try/command", url.Values{"name": {name}, "database": {"sqlite"}})
			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("name %q: status %d, want 422", name, rec.Code)
			}
			if strings.Contains(rec.Body.String(), "<script>") {
				t.Errorf("name %q was echoed unescaped", name)
			}
		}
	})

	t.Run("needs the CSRF token", func(t *testing.T) {
		anonymous := &client{t: t, handler: handler}
		form := url.Values{"name": {"my-app"}, "database": {"sqlite"}}
		if got := anonymous.do("POST", "/try/command", form, nil).Code; got != http.StatusForbidden {
			t.Errorf("status %d, want 403", got)
		}
	})
}

func TestProductionHeadersAndLimits(t *testing.T) {
	prodEnv := config.Env{GoEnv: "production", Port: "0"}
	handler, _ := newTestApp(t, prodEnv)
	c := newClient(t, handler)

	if rec := c.do("GET", "/", nil, nil); rec.Header().Get("Strict-Transport-Security") == "" {
		t.Error("HSTS must be sent in production")
	}
	if got := c.do("GET", "/docs/deployment", nil, nil).Code; got != http.StatusOK {
		t.Errorf("docs in production: status %d", got)
	}
}

// Behind a proxy, every visitor arrives from the proxy's address. With TrustedProxies set, each visitor
// is limited on their own, so one busy visitor cannot lock out another.
func TestTrustedProxyGivesEachVisitorTheirOwnLimit(t *testing.T) {
	prodEnv := config.Env{GoEnv: "production", Port: "0"}
	prodEnv.TrustedProxies = mustNets(t, "10.0.0.0/8")
	handler, _ := newTestApp(t, prodEnv)

	fromProxy := func(client string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("GET", "/healthz", nil)
		req.RemoteAddr = "10.0.0.2:5000"
		req.Header.Set("X-Forwarded-For", client)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	// one visitor exhausts their bucket (burst 60)
	limited := false
	for i := 0; i < 200; i++ {
		if fromProxy("198.51.100.7").Code == http.StatusTooManyRequests {
			limited = true
			break
		}
	}
	if !limited {
		t.Fatal("a busy visitor was never rate limited")
	}

	if got := fromProxy("198.51.100.8").Code; got != http.StatusOK {
		t.Errorf("a second visitor behind the same proxy: status %d, want 200", got)
	}
}
