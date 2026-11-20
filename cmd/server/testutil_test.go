package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"

	"github.com/owenwexler/gotth-plus-site/internal/app"
	"github.com/owenwexler/gotth-plus-site/internal/config"
	"github.com/owenwexler/gotth-plus-site/internal/database"
)

// TestMain moves to the project root, because the server reads static/ relative to its working directory.
func TestMain(m *testing.M) {
	if err := os.Chdir(filepath.Join("..", "..")); err != nil {
		fmt.Fprintln(os.Stderr, "chdir to project root:", err)
		os.Exit(1)
	}
	code := m.Run()
	for _, cleanup := range afterTests {
		cleanup()
	}
	os.Exit(code)
}

// afterTests run once, after every test in the package.
var afterTests []func()

// newTestApp builds the real handler (every route behind every middleware) on a throwaway database.
func newTestApp(t *testing.T, env config.Env) (http.Handler, *gorm.DB) {
	t.Helper()

	db, err := database.ConnectDb(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { closeDB(db) })

	// the server never migrates on its own (cmd/migrate does), so a fresh test database is migrated here
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	return app.NewHandler(env, db), db
}

var devEnv = config.Env{GoEnv: "development", Port: "0"}

func mustNets(t *testing.T, cidrs string) []*net.IPNet {
	t.Helper()
	nets, err := config.ParseTrustedProxies(cidrs)
	if err != nil {
		t.Fatal(err)
	}
	return nets
}
