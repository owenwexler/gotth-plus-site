package factory

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/owenwexler/gotth-plus-site/internal/docs"
	"github.com/owenwexler/gotth-plus-site/internal/repository"
	"github.com/owenwexler/gotth-plus-site/internal/view"
)

// renderPage writes a full page. A render error after the first byte is logged; the client has already
// received the start of the page, so there is nothing better to send.
func renderPage(w http.ResponseWriter, r *http.Request, status int, page templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := page.Render(r.Context(), w); err != nil {
		logError(r, "render page", err)
	}
}

// RenderDocsIndex serves /docs.
func RenderDocsIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, http.StatusOK, view.DocsIndex(docs.Groups()))
	}
}

// RenderDoc serves /docs/{slug}, with the page's vote counts read from the database.
func RenderDoc(votes repository.VoteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, ok := docs.Find(r.PathValue("slug"))
		if !ok {
			renderPage(w, r, http.StatusNotFound, view.NotFoundPage())
			return
		}

		group, _ := docs.GroupOf(page.Slug)

		counts, err := votes.Counts(r.Context(), page.Slug)
		if err != nil {
			respondWithErrorPage(w, r, http.StatusInternalServerError, "read vote counts", err)
			return
		}

		renderPage(w, r, http.StatusOK, view.DocPage(page, group, counts))
	}
}
