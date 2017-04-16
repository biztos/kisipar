// fsp -- a demo site running off the file system, sans frills.
// ---
// TODO: use kisipar.serve etc when available.
// TODO: simple opts for file, port (use docopt?)

// Kisipar demo filesystem-based server supporting Frosted Markdown and YAML.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/biztos/kisipar"
)

func main() {

	usage := "Usage: fsp CONTENT-DIR [TEMPLATE-DIR]"
	cfg := kisipar.FileSystemProviderConfig{}
	if len(os.Args) == 3 {
		cfg.ContentDir = os.Args[1]
		cfg.TemplateDir = os.Args[2]
	} else if len(os.Args) == 2 {
		cfg.ContentDir = os.Args[1]
	} else {
		log.Fatalf("Wrong number of args.\n%s\n", usage)
	}
	log.Println("Loading...")
	provider, err := kisipar.LoadFileSystemProvider(cfg)
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

	log.Println("Listening on port 8080.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}