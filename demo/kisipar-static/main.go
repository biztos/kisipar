// kisipar-static -- a demo site serving just static assets.
// --------------

// Kisipar demo static-asset server on port 8080.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/biztos/kisipar"
)

func main() {

	usage := "Usage: kisipar-static STATIC-DIR"
	dir := ""
	if len(os.Args) == 2 {
		dir = os.Args[1]
	} else {
		log.Fatalf("Wrong number of args.\n%s\n", usage)
	}

	// In real life you wouldn't want to do expose the top static dir as
	// the site root of course, because then you'd be going to all this
	// trouble just to make a worse version of $ANY_EXISTING_SERVER.
	// But for the sake of this demo we will do exactly that.
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Could not read %s: %s\n%s\n", dir, err.Error(), usage)
	}

	// Mimic the rather spartan default style from net/http.
	// The wonky indent here is so we can stick it in the yaml... less code.
	html := "        <pre>\n"
	for _, info := range infos {
		name := info.Name()
		link := name
		if info.IsDir() {
			link += "/"
		}
		html += fmt.Sprintf("        <a href=\"%s\">%s</a>\n", link, name)

	}
	html += "        </pre>\n"

	// Set up an empty provider and manually initialize the site.
	yaml := "templates:\n    index.html: |\n" + html
	provider, err := kisipar.StandardProviderFromYAML(yaml)
	if err != nil {
		// Programmer error, oops:
		panic(err)
	}
	config := &kisipar.Config{
		Name:       "kisipar-static",
		Port:       8080,
		StaticDir:  dir,
		ListStatic: true,
	}
	site := &kisipar.Site{
		Config:   config,
		Provider: provider,
	}
	if err := site.InitMux(); err != nil {
		log.Fatalf("InitMux failed: %s", err.Error())
	}
	if err := site.InitServer(); err != nil {
		log.Fatalf("InitServer failed: %s", err.Error())
	}

	log.Fatal(site.Serve())
}
