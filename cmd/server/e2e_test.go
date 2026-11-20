//go:build e2e

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mxschmitt/playwright-go"

	"github.com/owenwexler/gotth-plus-site/internal/docs"
)

// Browser tests for the site's own behaviour, in a real browser against the real handler. The harness is in
// e2e_helpers_test.go; checks that hold for any app are in e2e_generic_test.go.

// TestE2ECommandBuilder: the server answers the form with the command, and the result replaces the placeholder.
func TestE2ECommandBuilder(t *testing.T) {
	a := newE2EApp(t)
	a.open(t, "/")

	expect(t, assert.Locator(a.Page.Locator("#command-result")).ToContainText("Your command will appear here."))

	expect(t, a.Page.Locator("#project-name").Fill("my-app"))
	_, err := a.Page.Locator("#project-database").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice("postgres")})
	expect(t, err)
	expect(t, a.Page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Build command"}).Click())

	expect(t, assert.Locator(a.Page.Locator("#command-result code")).ToContainText("create-gotth-plus-app@latest my-app --db=postgres --yes"))
	expect(t, assert.Locator(a.Page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Copy command"})).ToBeVisible())
	a.expectNoConsoleErrors(t)
}

// TestE2EDocsNavigationIsBoosted: moving between pages swaps the body in place, so the page is not reloaded.
func TestE2EDocsNavigationIsBoosted(t *testing.T) {
	a := newE2EApp(t)
	a.open(t, "/docs/getting-started")
	a.markPage(t)

	expect(t, a.Page.Locator("#docs-nav").GetByRole("link", playwright.LocatorGetByRoleOptions{Name: "Security"}).Click())
	expect(t, assert.Page(a.Page).ToHaveURL(a.URL+"/docs/security"))
	expect(t, assert.Locator(a.Page.GetByRole("heading", playwright.PageGetByRoleOptions{Name: "Security", Level: playwright.Int(1)})).ToBeVisible())
	a.expectNotReloaded(t)
	a.expectNoConsoleErrors(t)
}

// TestE2EVoteRecordsAnswer: a visitor answers the question, and the totals come back from the server.
func TestE2EVoteRecordsAnswer(t *testing.T) {
	a := newE2EApp(t)
	a.open(t, "/docs/cli")

	expect(t, assert.Locator(a.Page.Locator("#vote")).ToContainText("Be the first to answer."))
	expect(t, a.Page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Yes", Exact: playwright.Bool(true)}).Click())

	expect(t, assert.Locator(a.Page.Locator("#vote")).ToContainText("Thanks for your feedback."))
	expect(t, assert.Locator(a.Page.Locator("#vote")).ToContainText("1 of 1 readers found this page helpful."))
	a.expectNoConsoleErrors(t)
}

// TestE2EDocsMenuOpensOnPhones: on a phone the sidebar is behind a button, and the button opens it.
func TestE2EDocsMenuOpensOnPhones(t *testing.T) {
	a := newE2EApp(t)
	a.setViewport(t, 375, 800)
	a.open(t, "/docs/cli")

	nav := a.Page.Locator("#docs-nav")
	expect(t, assert.Locator(nav).ToBeHidden())
	expect(t, a.Page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Pages"}).Click())
	expect(t, assert.Locator(nav).ToBeVisible())
}

// TestE2EAxeOnEveryPage runs the axe scan on the home page, the docs index and each documentation page, the
// pages with tables and code. It needs axe-core: set AXE_CORE_JS to an axe.min.js, and it skips otherwise.
func TestE2EAxeOnEveryPage(t *testing.T) {
	axePath := os.Getenv("AXE_CORE_JS")
	if axePath == "" {
		axePath = "cmd/server/testdata/axe.min.js"
	}
	if _, err := os.Stat(axePath); err != nil {
		t.Skipf("axe-core not found at %s: set AXE_CORE_JS or copy axe.min.js there", axePath)
	}

	paths := []string{"/", "/docs"}
	for _, d := range docs.All() {
		paths = append(paths, "/docs/"+d.Slug)
	}

	a := newE2EApp(t)
	// the CSP only allows the app's own scripts, so this context bypasses it to inject axe
	scan := a.inContext(t, playwright.BrowserNewContextOptions{BypassCSP: playwright.Bool(true)})

	for _, path := range paths {
		scan.open(t, path)
		expect(t, scan.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle}))

		_, err := scan.Page.AddScriptTag(playwright.PageAddScriptTagOptions{Path: playwright.String(axePath)})
		expect(t, err)
		result, err := scan.Page.Evaluate(`axe.run().then(r => r.violations
			.filter(v => v.impact === 'serious' || v.impact === 'critical')
			.map(v => v.impact + ' ' + v.id + ': ' + v.help + ' (' + v.nodes.length + ' elements)'))`)
		expect(t, err)

		if violations, _ := result.([]any); len(violations) > 0 {
			lines := make([]string, len(violations))
			for i, v := range violations {
				lines[i] = fmt.Sprint(v)
			}
			t.Errorf("%s: axe found serious or critical violations:\n%s", path, strings.Join(lines, "\n"))
		}
	}
}
