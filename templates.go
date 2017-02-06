// templates.go -- general template logic

package kisipar

import (
	"html/template"
)

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
// If the Page has a MetaString named Template and there is a template of that
// exact name available, it is returned.  No guesses are made regarding the
// template in this case: the match must be perfectly exact.
//
// Exact Match
//
// If the Page's Path exactly matches a template path, that template is
// returned. In most cases the root directory of the templates will be
// stripped before assigning template paths, but this is up to the
// implementation. A suffix of ".html" will be checked as a fallback for the
// template, in keeping with the common practice of using that extension for
// Go templates.
//
//  Path: /foo/bar --> Template: /foo/bar OR /foo/bar.html
//
// Best Guess
//
// If by walking up the path we find a template matching the dirname of the
// path, that template is returned:
//
//  Path: /foo/bar/baz --> Template: /foo/bar OR /foo/bar.html
//                                      OR /foo OR /foo.html
//                                      (but not "" nor ".html")
//
// Default
//
// The master template is returned, i.e. tmpl.
//
// Note
//
// There is usually more to template selection than simply matching the page.
// In most cases there will be templates which are served based on URL path
// matching; however for single-item pages the PageTemplate function is
// recommended.
func PageTemplate(tmpl *template.Template, p Page) *template.Template {
	if tmpl == nil {
		return nil
	}
	return tmpl // TODO!
	// var tpath string
	// ...
	return nil

}
