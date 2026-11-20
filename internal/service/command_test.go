package service

import (
	"errors"
	"strings"
	"testing"
)

func TestProjectCommand(t *testing.T) {
	got, err := ProjectCommand("my-app", "postgres")
	if err != nil {
		t.Fatal(err)
	}
	want := "go run github.com/owenwexler/create-gotth-plus-app@latest my-app --db=postgres --yes"
	if got != want {
		t.Errorf("ProjectCommand = %q, want %q", got, want)
	}
}

func TestProjectCommandRejectsNames(t *testing.T) {
	for _, name := range []string{"", "1app", "My App", "app;rm -rf /", "../etc", "-leading", "UPPER", strings.Repeat("a", 41)} {
		if _, err := ProjectCommand(name, "sqlite"); !errors.Is(err, ErrInvalidProjectName) {
			t.Errorf("ProjectCommand(%q): err = %v, want ErrInvalidProjectName", name, err)
		}
	}

	if _, err := ProjectCommand(strings.Repeat("a", 40), "sqlite"); err != nil {
		t.Errorf("a 40 character name must be accepted: %v", err)
	}
}
