package factory

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/owenwexler/gotth-plus-site/internal/docs"
	"github.com/owenwexler/gotth-plus-site/internal/dto"
	"github.com/owenwexler/gotth-plus-site/internal/repository"
	"github.com/owenwexler/gotth-plus-site/internal/view"
)

// Vote records a "Was this page helpful?" answer and returns the vote box with the new totals.
// POST /docs/{slug}/vote?vote=yes|no
func Vote(votes repository.VoteRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// the slug must be one of our pages; it is never trusted as a database key
		page, ok := docs.Find(r.PathValue("slug"))
		if !ok {
			respondWithToast(w, r, http.StatusNotFound, view.ErrorToast("Page not found"))
			return
		}

		req := dto.VoteRequest{Vote: strings.TrimSpace(r.URL.Query().Get("vote"))}
		if err := validate.Struct(&req); err != nil {
			logRejected(r, "vote validation failed", err)
			respondWithToast(w, r, http.StatusUnprocessableEntity, view.ErrorToast("Could not record your vote"))
			return
		}

		counts, err := votes.Record(r.Context(), page.Slug, req.Vote == "yes")
		if err != nil {
			respondWithToast(w, r, http.StatusInternalServerError, view.ErrorToast("Could not record your vote"))
			logError(r, "record vote", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := view.VoteBox(page.Slug, counts, req.Vote).Render(r.Context(), w); err != nil {
			logError(r, "render vote box", err)
		}
	}
}
