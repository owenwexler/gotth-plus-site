package model

import "time"

// PageVote counts the "Was this page helpful?" answers for one docs page. Slug is the page's identifier
// from internal/docs (a fixed list, never user input), so it is the primary key.
type PageVote struct {
	Slug       string `gorm:"primaryKey;size:64"`
	Helpful    int64  `gorm:"not null;default:0"`
	NotHelpful int64  `gorm:"not null;default:0"`
	UpdatedAt  time.Time
}
