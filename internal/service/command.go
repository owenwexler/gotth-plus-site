// Package service holds the site's business rules. It has no HTTP and no SQL.
package service

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrInvalidProjectName is returned for names the project generator would not accept as a directory name.
var ErrInvalidProjectName = errors.New("invalid project name")

// projectName: lowercase letters, digits and hyphens, starting with a letter, at most 40 characters.
var projectName = regexp.MustCompile(`^[a-z][a-z0-9-]{0,39}$`)

// ProjectCommand returns the shell command that creates a project with the given name and database,
// with every prompt answered by a flag (--yes), as create-gotth-plus-app documents it for scripts.
func ProjectCommand(name, database string) (string, error) {
	if !projectName.MatchString(name) {
		return "", ErrInvalidProjectName
	}

	return fmt.Sprintf("go run github.com/owenwexler/create-gotth-plus-app@latest %s --db=%s --yes", name, database), nil
}
