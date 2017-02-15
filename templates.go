// templates.go -- general template logic

package kisipar

import (
	"errors"
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/biztos/vebben"
)

// DefaultFuncMap is the starting point for FuncMap.
var DefaultFuncMap = vebben.NewFuncMap()

// FuncMap returns a function map for use in templates.  It starts with the
// map in DefaultFuncMap and adds the following Provider-related functions:
//
// XXXXX TODO
func FuncMap() template.FuncMap {
	return DefaultFuncMap // TODO: add functions
}

// TemplatesFromData returns a set of templates under a master template.
// The master template is empty and has the empty string as its name.
// All other templates are parsed from the input string value and given
// the input key name.  The first error encountered is returned.
//
// The templates have functions available as defined in DefaultFuncMap,
// which is designed to cover may be overridden.
func TemplatesFromData(in map[string]string) (*template.Template, error) {

	if in == nil {
		return nil, errors.New("TemplatesFromData input may not be nil.")
	}
	tmpl, err := template.New("").Funcs(FuncMap()).Parse("")
	if err != nil {
		return nil, fmt.Errorf("Master template failed: %s", err.Error())
	}

	for path, src := range in {
		if _, err := tmpl.New(path).Parse(src); err != nil {
			return nil, fmt.Errorf("Template %s failed: %s", path, err.Error())
		}
	}
	return tmpl, nil
}

// PageTemplate returns the template recommended for rendering the given
// page, as follows:
//
// Nil Template
//
// If the provided template is nil, nil is returned.  In this case we assume
// the HTTP handler knows what to do with the site's templates.
//
// Meta Template
//
// If the Page has a MetaString named "Template" (or "template" or
// "TEMPLATE") and there is a template of that exact name available, it is
// returned. No guesses are made regarding the template in this case: the
// match must be exact.
//
// Path Match
//
// If the Page's Path exactly matches a template path, that template is
// returned. Variations are checked in order to reflect common practices:
// the path with and without a leading slash; and with and without a suffix
// of ".html" (leading-slash practice is easy to get wrong, and templates
// are usually stored as files with an .html extension).
//
//  Path: /foo --> Template: /foo OR foo OR /foo.html OR foo.html
//
// Best Guess
//
// If by walking up the path we find a template matching the dirname of the
// path, that template is returned:
//
//  Path: /foo/bar/baz --> Template: /foo/bar OR foo/bar
//                                     OR /foo/bar.html OR foo/bar.html
//                                     OR /foo OR foo
//                                     OR /foo.html OR foo.html
//                                     (but not "" nor ".html")
//
// The same variations are checked as in the Path Match.
//
// Default
//
// A template named "/default" is returned if present; again this is subject
// to the variations above.
//
// No Match
//
// If there is no match at all, nil is returned and the caller is responsible
// for finding a template.
//
// Note
//
// There is usually more to template selection than simply matching the page.
// In most cases there will be some templates that are served based on URL
// path matching before the Provider is consulted; however for single-item
// pages the PageTemplate function is recommended.
func PageTemplate(tmpl *template.Template, p Page) *template.Template {

	// Nil in, nil out.
	if tmpl == nil {
		return nil
	}

	// Take the page's word where possible.
	tname := p.MetaString("Template")
	if tname == "" {
		tname = p.MetaString("template")
	}
	if tname == "" {
		tname = p.MetaString("TEMPLATE")
	}
	if tname != "" {
		if match := tmpl.Lookup(tname); match != nil {
			return match
		}
	}

	if p.Path() != "/" {

		// An exact-ish match makes things easy.
		if match := lookupPathVariations(tmpl, p.Path()); match != nil {
			return match
		}

		// Otherwise, up we go!
		for d := path.Dir(p.Path()); d != "/" && d != "" && d != "."; d = path.Dir(d) {
			if match := lookupPathVariations(tmpl, d); match != nil {
				return match
			}
		}
	}

	// Fallback is the default.html or similar:
	if match := lookupPathVariations(tmpl, "/default"); match != nil {
		return match
	}

	// No match and no default.
	return nil

}

func lookupPathVariations(t *template.Template, p string) *template.Template {

	if tmpl := t.Lookup(p); tmpl != nil {
		return tmpl
	}
	if strings.HasPrefix(p, "/") {
		if tmpl := t.Lookup(strings.TrimPrefix(p, "/")); tmpl != nil {
			return tmpl
		}
	}

	if !strings.HasSuffix(p, ".html") {
		p = p + ".html"
		if tmpl := t.Lookup(p); tmpl != nil {
			return tmpl
		}
		if strings.HasPrefix(p, "/") {
			if tmpl := t.Lookup(strings.TrimPrefix(p, "/")); tmpl != nil {
				return tmpl
			}
		}
	}

	return nil
}
