// yamly -- a demo site running off YAML
// -----
// TODO: use kisipar.serve etc when available.
// TODO: simple opts for file, port (use docopt?)

// Kisipar demo YAML-based server.
package main

import (
	"log"
	"net/http"

	"github.com/biztos/kisipar"
)

func main() {

	provider, err := kisipar.StandardProviderFromYaml(DefaultYaml)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		item, err := provider.Get(r.URL.Path)
		if kisipar.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			log.Fatalf("Error from Get for path %s: %s",
				r.URL.Path, err.Error())
		}

		// Our YAML-based provider only knows about two item types.
		if c, ok := item.(kisipar.Content); ok {
			if ct := c.ContentType(); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
			http.ServeContent(w, r, "", c.ModTime(), c.ReadSeeker())
			return
		}
		if p, ok := item.(kisipar.Page); ok {
			tmpl := provider.TemplateFor(p)
			if tmpl == nil {
				log.Fatal("No template returned for " + p.Path())
			}
			tmpl.Execute(w, p)
			return
		}
		// Any other type means we forgot to keep the code up to date. :-(
		log.Fatalf("Unexpected type for %s: %T", item.Path(), item)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

var DefaultYaml = `# I am YAML! YAMLIAM!
pages:
    /:
        title: Kisipar from YAML
        html: |
            Hello world!
    /foo/bar:
        title: I am the Foo Bar!
        tags: [foo,bar]
        created: 2016-01-02T15:04:05Z
        updated: 2017-02-02T15:04:05Z
        content: |
            This is the foo, the bar, the baz and
            the bat if you like.  For sanity's sake
            let's not let it be Markdown.
    /baz/bat:
        title: The BazzerBat
        tags: [foo,bazzers,badgers]
content:
    /js/goober.js:
        type: application/javascript
        content: |
            window.alert('hello world');
templates:
    any/random/tmpl.html: |
        <!doctype html>
        <script src="/js/goober.js"></script>
        <h1>Hello {{ .Title }}</h1>
`
