package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/owenwexler/gotth-plus-site/internal/model"
)

// VoteRepository stores the helpful/not-helpful counts of the docs pages.
type VoteRepository interface {
	// Counts returns the totals for a page. A page nobody has voted on has zero counts and no error.
	Counts(ctx context.Context, slug string) (model.PageVote, error)
	// Record adds one vote and returns the new totals.
	Record(ctx context.Context, slug string, helpful bool) (model.PageVote, error)
}

type voteRepository struct{ db *gorm.DB }

// NewVoteRepository creates a VoteRepository backed by db.
func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) Counts(ctx context.Context, slug string) (model.PageVote, error) {
	var vote model.PageVote

	// Find, not First: a page with no votes yet is normal, and First would log it as an error
	result := r.db.WithContext(ctx).Where("slug = ?", slug).Limit(1).Find(&vote)
	if result.Error != nil {
		return model.PageVote{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.PageVote{Slug: slug}, nil
	}

	return vote, nil
}

// Record is a single upsert: the database adds one to the counter, so two votes at the same moment both count.
func (r *voteRepository) Record(ctx context.Context, slug string, helpful bool) (model.PageVote, error) {
	now := time.Now()

	row := model.PageVote{Slug: slug, UpdatedAt: now}
	var bump clause.Set

	if helpful {
		row.Helpful = 1
		bump = clause.Assignments(map[string]any{"helpful": gorm.Expr("helpful + 1"), "updated_at": now})
	} else {
		row.NotHelpful = 1
		bump = clause.Assignments(map[string]any{"not_helpful": gorm.Expr("not_helpful + 1"), "updated_at": now})
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoUpdates: bump}).
		Create(&row).Error
	if err != nil {
		return model.PageVote{}, err
	}

	return r.Counts(ctx, slug)
}
