package docs

import (
	"strings"
	"testing"
)

func TestCatalogIsWellFormed(t *testing.T) {
	slugs := map[string]bool{}

	for _, g := range Groups() {
		if g.Title == "" || len(g.Docs) == 0 {
			t.Errorf("group %q is empty or untitled", g.Title)
		}

		for _, d := range g.Docs {
			if slugs[d.Slug] {
				t.Errorf("duplicate slug %q", d.Slug)
			}
			slugs[d.Slug] = true

			if d.Title == "" || d.Summary == "" || len(d.Sections) == 0 {
				t.Errorf("%s: needs a title, a summary and at least one section", d.Slug)
			}

			anchors := map[string]bool{}
			for _, s := range d.Sections {
				if s.Anchor() == "" || anchors[s.Anchor()] {
					t.Errorf("%s: section %q has an empty or duplicate anchor", d.Slug, s.Title)
				}
				anchors[s.Anchor()] = true

				for _, b := range s.Blocks {
					switch b.Kind {
					case "p", "code":
						if b.Text == "" {
							t.Errorf("%s/%s: empty %s block", d.Slug, s.Title, b.Kind)
						}
					case "ul":
						if len(b.Items) == 0 {
							t.Errorf("%s/%s: empty list", d.Slug, s.Title)
						}
					case "links":
						if len(b.Links) == 0 {
							t.Errorf("%s/%s: empty links list", d.Slug, s.Title)
						}
						for _, l := range b.Links {
							if l.Label == "" || !strings.HasPrefix(l.URL, "https://") {
								t.Errorf("%s/%s: link %q needs a label and an https:// URL", d.Slug, s.Title, l.URL)
							}
						}
					case "table":
						for _, row := range b.Rows {
							if len(row) != len(b.Head) {
								t.Errorf("%s/%s: a row has %d cells, the header has %d", d.Slug, s.Title, len(row), len(b.Head))
							}
						}
					default:
						t.Errorf("%s/%s: unknown block kind %q", d.Slug, s.Title, b.Kind)
					}
				}
			}
		}
	}

	if !slugs["deployment"] {
		t.Error("the deployment page must exist: the site links to it")
	}
}

func TestFindAndGroupOf(t *testing.T) {
	d, ok := Find("cli")
	if !ok || d.Title == "" {
		t.Fatal("Find(cli) failed")
	}
	if _, ok := Find("../etc/passwd"); ok {
		t.Error("Find must only know our slugs")
	}
	if g, ok := GroupOf("cli"); !ok || g.Title != "Get started" {
		t.Errorf("GroupOf(cli) = %q, %v", g.Title, ok)
	}
}

func TestAnchor(t *testing.T) {
	for in, want := range map[string]string{"What you need": "what-you-need", "Pages & fragments": "pages--fragments", "CSRF (double-submit)": "csrf-double-submit"} {
		if got := (Section{Title: in}).Anchor(); got != want {
			t.Errorf("Anchor(%q) = %q, want %q", in, got, want)
		}
	}
}
