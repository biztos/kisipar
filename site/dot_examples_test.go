// site/dot_examples_test.go - examples for the Dot.
// -------------------------

package site_test

import (
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/site"
	"log"
	"net/http"
	"os"
	"text/template"
)

func ExampleSite_NewDot() {

	// Let's imagine a Dot for a single non-index page, say "/foo/bar.md"
	p, err := page.LoadVirtualString("foo/bar.md", "# Bar!")
	if err != nil {
		log.Fatal(err)
	}

	// We've requested it from somewhere:
	r, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}

	// We have a Site:
	s, err := site.LoadVirtual(nil, []*page.Page{p}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Now we can have a Dot!
	dot := &site.Dot{
		Site:    s,
		Request: r,
		Page:    p,
		Pageset: nil, // nil because we're not in a directory
	}

	// Well enough, let's render it now.
	tmplStr := `{{ .Site.Name }} - {{ .Page.Title }} : {{ .URL }}`
	tmpl, err := template.New("example").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(os.Stdout, dot); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Anonymous Kisipar Site - Bar! : http://localhost/foo/bar

}
