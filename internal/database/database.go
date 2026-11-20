package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// ConnectDb opens the SQLite database at path (created if it does not exist). Every connection uses WAL
// journaling, so readers never block the writer, and waits up to five seconds for a lock before failing.
// Foreign keys are enforced.
func ConnectDb(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", path)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
