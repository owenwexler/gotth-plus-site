// Command migrate brings the database schema up to date. Run it on every deploy, before the new server
// starts. The server never changes the schema itself.
package main

import (
	"log/slog"
	"os"

	"github.com/owenwexler/gotth-plus-site/internal/config"
	"github.com/owenwexler/gotth-plus-site/internal/database"
)

func main() {
	env := config.GetEnv()

	db, err := database.ConnectDb(env.DatabasePath)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("getting database handle", "err", err)
		os.Exit(1)
	}
	defer func() { _ = sqlDB.Close() }()

	if err := database.Migrate(db); err != nil {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}

	slog.Info("schema is up to date", "database", env.DatabasePath)
}
