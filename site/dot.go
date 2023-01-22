// site/dot.go - the Kisipar Dot!
// -----------

package site

import (
	// Standard library:
	"html/template"
	"net/http"
	"strings"
	"time"

	// Kisipar packages:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
)

// A Dot is the "dot" available in a template for rendering a page.
type Dot struct {

	// The full Request used to retrieve the page.  Use with caution when
	// caching rendered pages.
	Request *http.Request

	// The Page being rendered; nil if rendering an anonymous index.
	Page *page.Page

	// The Pageset if we are rendering an index; nil if we are rendering a
	// final page.
	Pageset *pageset.Pageset

	// The Site, which of course contains its own Pageset and so on.
	Site *Site

	// Now hold the timestamp of the Dot's creation.
	Now time.Time

	// Register is useful in templates when one needs, say, to keep track
	// of indent levels.  It is set via - ta-da! - SetRegister.
	Register int
}

// SetRegister sets the Register to the provided value and returns the
// previous value.  It is mostly useful in templates.
func (d *Dot) SetRegister(v int) int {
	old := d.Register
	d.Register = v
	return old
}

// URL returns the full URL string of the request.
func (d *Dot) URL() string {
	if d.Request == nil {
		return ""
	}
	return d.Request.URL.String()
}

// Template returns the template in which to render the Dot, based on the
// following logic.
//
// TODO: possibly start with a site-level configured override set in case
// you want to say /oranges/* uses /fruit.html or whatever.
//
// If the Page defines a Template property in its Meta, and that template
// exists as an exact match, then it is used.  This allows per-page template
// overrides.
//
// If no override is present, a match is sought based on the Request URL.
// An exact match of the normalized request path to the template name is
// preferred; failing that, a template named "index" (if the Dot has a
// Pageset) or "single" is sought at the same level.  This is repeated
// up the path chain until we hit or miss at the top-level "index" or
// "single" template.  The fallback is the top-level (default) template.
//
// Thus a request for /foo/bar with a Pageset present matches:
//  foo/bar
//  foo/bar/index
//  foo
//  foo/index
//  index
//  <default>
func (d *Dot) Template() *template.Template {

	// Sanity checks:
	if d.Site == nil {
		panic("Site is nil")
	}
	if d.Site.Template == nil {
		panic("Site.Template is nil")
	}

	// Page override:
	if d.Page != nil {
		if name := d.Page.MetaString("Template"); name != "" {
			if tmpl := d.Site.Template.Lookup(name); tmpl != nil {
				return tmpl
			}
		}
	}

	// Request based:
	alt := "single"
	if d.Pageset != nil {
		alt = "index"
	}
	if d.Request != nil {
		path := strings.TrimPrefix(strings.ToLower(d.Request.URL.Path), "/")
		if tmpl := d.Site.Template.Lookup(path + "/" + alt); tmpl != nil {
			return tmpl
		}
		parts := strings.Split(path, "/")
		for len(parts) > 0 {
			name := strings.Join(parts, "/")
			if name != "" {
				if tmpl := d.Site.Template.Lookup(name); tmpl != nil {
					return tmpl
				}
				if tmpl := d.Site.Template.Lookup(name + "/" + alt); tmpl != nil {
					return tmpl
				}
			}
			parts = parts[:len(parts)-1]

		}
	}

	// Top-level special templates:
	if tmpl := d.Site.Template.Lookup(alt); tmpl != nil {
		return tmpl
	}

	// Final fallback: the top-level (default) template.
	return d.Site.Template

}
