// main.go - The Kisipar web server application.

// The kisipar program loads and serves Kisipar web sites.  For more
// information use the --help option, or consult the project page:
// https://github.com/biztos/kisipar
package main

import (
	"log"
	"os"

	"github.com/biztos/kisipar/app"
)

const VERSION = "0.1.0"

var EXIT_FUNCTION = os.Exit

func main() {

	err := app.Run(
		"Kisipar", // heading name
		VERSION,   // current version
		"kisipar", // binary name
	)
	if err != nil {
		log.Println(err)
		EXIT_FUNCTION(1)
	} else {
		EXIT_FUNCTION(0)
	}
}
