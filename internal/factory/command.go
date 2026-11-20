package factory

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/owenwexler/gotth-plus-site/internal/dto"
	"github.com/owenwexler/gotth-plus-site/internal/service"
	"github.com/owenwexler/gotth-plus-site/internal/view"
)

// Command answers the project command builder on the home page with the command for the chosen name and database.
func Command(validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logRejected(r, "parse command form", err)
			respondWithToast(w, r, http.StatusBadRequest, view.ErrorToast("Could not build the command"))
			return
		}

		req := dto.CommandRequest{
			Name:     strings.TrimSpace(r.FormValue("name")),
			Database: strings.TrimSpace(r.FormValue("database")),
		}

		if err := validate.Struct(&req); err != nil {
			logRejected(r, "command validation failed", err)
			respondWithToast(w, r, http.StatusUnprocessableEntity, view.ToastArgs{Kind: view.ToastError, Header: "Could not build the command", Body: "Choose a database and enter a name"})
			return
		}

		command, err := service.ProjectCommand(req.Name, req.Database)
		if errors.Is(err, service.ErrInvalidProjectName) {
			logRejected(r, "command name rejected", err)
			respondWithToast(w, r, http.StatusUnprocessableEntity, view.ToastArgs{Kind: view.ToastError, Header: "Use a valid project name", Body: "Lowercase letters, digits and hyphens, starting with a letter"})
			return
		}
		if err != nil {
			respondWithErrorPage(w, r, http.StatusInternalServerError, "build command", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := view.CommandResult(command).Render(r.Context(), w); err != nil {
			logError(r, "render command", err)
		}
	}
}
