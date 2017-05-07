// kisipar-fsp -- a demo site running off the file system, sans frills.
// -----------
// TODO: simple opts for file, port (use docopt?)

// Kisipar demo filesystem-based server supporting Frosted Markdown and YAML.
package main

import (
	"log"
	"os"

	"github.com/biztos/kisipar"
)

func main() {

	usage := "Usage: kisipar-fsp CONTENT-DIR [TEMPLATE-DIR]"
	cdir := ""
	tdir := ""
	if len(os.Args) == 3 {
		cdir = os.Args[1]
		tdir = os.Args[2]
	} else if len(os.Args) == 2 {
		cdir = os.Args[1]
	} else {
		log.Fatalf("Wrong number of args.\n%s\n", usage)
	}
	log.Println("Loading...")
	site, err := kisipar.NewSite(&kisipar.Config{
		Port:     8080,
		Provider: "filesystem",
		ProviderConfig: map[string]interface{}{
			"ContentDir":  cdir,
			"TemplateDir": tdir,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(site.Serve())
}
