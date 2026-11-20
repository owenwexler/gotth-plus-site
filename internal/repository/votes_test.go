package repository

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/owenwexler/gotth-plus-site/internal/database"
)

func newVoteRepo(t *testing.T) VoteRepository {
	t.Helper()

	db, err := database.ConnectDb(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}

	return NewVoteRepository(db)
}

func TestVoteCounts(t *testing.T) {
	repo := newVoteRepo(t)
	ctx := context.Background()

	got, err := repo.Counts(ctx, "cli")
	if err != nil {
		t.Fatal(err)
	}
	if got.Slug != "cli" || got.Helpful != 0 || got.NotHelpful != 0 {
		t.Errorf("a page with no votes = %+v, want zero counts", got)
	}

	for _, helpful := range []bool{true, true, false} {
		if _, err := repo.Record(ctx, "cli", helpful); err != nil {
			t.Fatal(err)
		}
	}

	got, err = repo.Counts(ctx, "cli")
	if err != nil {
		t.Fatal(err)
	}
	if got.Helpful != 2 || got.NotHelpful != 1 {
		t.Errorf("after two yes and one no: helpful=%d not=%d, want 2 and 1", got.Helpful, got.NotHelpful)
	}

	// other pages are counted separately
	if other, _ := repo.Counts(ctx, "security"); other.Helpful != 0 {
		t.Errorf("votes leaked to another page: %+v", other)
	}
}

// Votes that arrive at the same moment must all be counted: the increment happens in the database.
func TestConcurrentVotesAreAllCounted(t *testing.T) {
	repo := newVoteRepo(t)
	ctx := context.Background()

	const votes = 25
	var wg sync.WaitGroup
	errs := make(chan error, votes)
	for i := 0; i < votes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repo.Record(ctx, "testing", true); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("record vote: %v", err)
	}

	got, err := repo.Counts(ctx, "testing")
	if err != nil {
		t.Fatal(err)
	}
	if got.Helpful != votes {
		t.Errorf("helpful = %d, want %d", got.Helpful, votes)
	}
}
