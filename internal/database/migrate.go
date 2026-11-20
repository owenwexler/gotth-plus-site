package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/owenwexler/gotth-plus-site/internal/model"
)

// Migrate brings the schema up to date. Every model the app stores belongs in this list. The server never
// calls it: run cmd/migrate on each deploy, before the new server starts.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.PageVote{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
