// kisipar.go
// ----------

// Package kisipar provides an opinionated web server for small(ish)
// Markdown-based sites.  The API presented here allows for a standard
// server to be run via the command subpackage cmd/kisipar, or through a
// custom application.  Much deeper customization is possible using the
// subpackages.
//
// It is STRONGLY recommended that you run any public-facing kisipar servers
// behind a reverse-proxy web server such as Nginx: https://www.nginx.com
//
// For more information see https://github.com/biztos/kisipar
//
// Building the Server
//
// The standard server should be sufficient for most intended purposes.
//
//  go get github.com/biztos/kisipar
//  go build github.com/biztos/kisipar/cmd/kisipar
//  ./kisipar --help
//
// You may of course build a custom server in order to expand -- or contract
// -- the Kisipar functionality.
//
// Site Configuration
//
// The configuration file is in YAML format, must be named "config.yaml" and
// lives at the top of the site directory.
//
// Site Layout
//
// A standard Kisipar site is contained in a directory, with a YAML
// configuration file and three subdirectories:
//
//   config.yaml
//   pages/
//   static/
//   templates/
//
// Static assets override pages.  Templates are go-style (html/template).
package kisipar

import (
	"errors"
	"fmt"
	"log"

	"github.com/biztos/kisipar/site"
)

// LAUNCH_SERVERS controls whether to actually launch the site servers; set to
// false for testing without spawning listeners.
var LAUNCH_SERVERS = true

// Kisipar represents a set of one or more Sites to serve.
type Kisipar struct {
	Sites []*site.Site
}

// Load initializes a Kisipar struct with sites loaded from the  directories
// located at the given paths.  Each site must have its own config file.
func Load(paths ...string) (*Kisipar, error) {
	if len(paths) == 0 {
		return nil, errors.New("kisipar.Load requires at least one site path.")
	}
	portPaths := map[int]string{}
	sites := make([]*site.Site, len(paths))
	for i, path := range paths {
		site, err := site.Load(path)
		if err != nil {
			return nil, fmt.Errorf("Site error at %s: %s", path, err)
		}
		if pp := portPaths[site.Port]; pp != "" {
			return nil, fmt.Errorf("Duplicate Port %d: %s vs. %s",
				site.Port, pp, path)
		}
		portPaths[site.Port] = path
		sites[i] = site
	}
	return &Kisipar{Sites: sites}, nil
}

// Serve launches listen-and-serve routines for all Sites, with or without
// TLS as per the configuration.  The last site in the list will block until
// it is finished.
func (k *Kisipar) Serve() {

	final := len(k.Sites) - 1
	for i, s := range k.Sites {
		log.Printf("%s: listening on port %d.", s.Name, s.Port)
		if LAUNCH_SERVERS {
			if i == final {
				log.Printf("%s (port %d): %s\n", s.Name, s.Port, s.Serve())
			} else {
				go log.Printf("%s (port %d): %s\n", s.Name, s.Port, s.Serve())
			}
		}
	}

}
