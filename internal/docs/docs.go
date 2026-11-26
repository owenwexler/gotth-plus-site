// Package docs is the documentation site's content. It is Go data rather than Markdown so that every
// page is rendered by Templ, escaped like the rest of the site, with no raw HTML.
package docs

import "strings"

// Block is one piece of a section. Kind is "p" (Text), "code" (Text), "ul" (Items), "table" (Head, Rows)
// or "links" (Links).
type Block struct {
	Kind  string
	Text  string
	Items []string
	Head  []string
	Rows  [][]string
	Links []Link
}

// Link is an external link with a short note on what is there. URL must be an https:// address.
type Link struct {
	Label string
	URL   string
	Note  string
}

// P is a paragraph.
func P(text string) Block { return Block{Kind: "p", Text: text} }

// Code is a preformatted block of code or a terminal session.
func Code(text string) Block { return Block{Kind: "code", Text: text} }

// UL is a bulleted list.
func UL(items ...string) Block { return Block{Kind: "ul", Items: items} }

// Links is a bulleted list of external links.
func Links(links ...Link) Block { return Block{Kind: "links", Links: links} }

// Table is a table with a header row.
func Table(head []string, rows ...[]string) Block {
	return Block{Kind: "table", Head: head, Rows: rows}
}

// Section is a titled part of a page.
type Section struct {
	Title  string
	Blocks []Block
}

// Anchor is the section's id on the page, so it can be linked to (#section-title).
func (s Section) Anchor() string { return anchor(s.Title) }

// Doc is one documentation page.
type Doc struct {
	Slug     string
	Title    string
	Summary  string
	Sections []Section
}

// Group is a heading in the navigation with the pages under it.
type Group struct {
	Title string
	Docs  []Doc
}

func anchor(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ', r == '-':
			b.WriteRune('-')
		}
	}
	return b.String()
}

// Groups returns the navigation, in order.
func Groups() []Group { return groups }

// All returns every page in navigation order.
func All() []Doc {
	var all []Doc
	for _, g := range groups {
		all = append(all, g.Docs...)
	}
	return all
}

// Find returns the page with the given slug.
func Find(slug string) (Doc, bool) {
	for _, d := range All() {
		if d.Slug == slug {
			return d, true
		}
	}
	return Doc{}, false
}

// Group returns the navigation group a page belongs to.
func GroupOf(slug string) (Group, bool) {
	for _, g := range groups {
		for _, d := range g.Docs {
			if d.Slug == slug {
				return g, true
			}
		}
	}
	return Group{}, false
}
