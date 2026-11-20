package database

import (
	"path/filepath"
	"testing"
)

func TestConnectUsesWALAndMigrates(t *testing.T) {
	db, err := ConnectDb(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}

	var mode string
	if err := db.Raw("PRAGMA journal_mode").Scan(&mode).Error; err != nil {
		t.Fatal(err)
	}
	if mode != "wal" {
		t.Errorf("journal_mode = %q, want wal", mode)
	}

	var timeout int
	if err := db.Raw("PRAGMA busy_timeout").Scan(&timeout).Error; err != nil {
		t.Fatal(err)
	}
	if timeout != 5000 {
		t.Errorf("busy_timeout = %d, want 5000", timeout)
	}

	// migrating twice must be harmless: the migrator runs on every deploy
	for i := 0; i < 2; i++ {
		if err := Migrate(db); err != nil {
			t.Fatalf("migrate run %d: %v", i+1, err)
		}
	}
	if !db.Migrator().HasTable("page_votes") {
		t.Error("page_votes table was not created")
	}
}
