// app_examples_test.go - examples, and testing trickery, for kisipar app!

package app_test

import (

	// Standard Library:
	"log"
	"os"
	"path/filepath"

	// Kisipar:
	"github.com/biztos/kisipar"
	"github.com/biztos/kisipar/app"
)

func Example() {

	// For testing purposes, let's not really launch the servers:
	kisipar.LAUNCH_SERVERS = false

	// ...and let's not do tricky logging:
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// We conveniently have some test sites on disk we can use for this:
	path1 := filepath.Join("test_data", "site_1")
	path2 := filepath.Join("test_data", "site_2")

	// Imagining this usage:
	os.Args = []string{
		"kisipar", // arg zero: the app
		path1,     // arg one: the first site to launch
		path2,     // arg two: the second site to launch
	}

	// This, our simple launcher, will launch those:
	err := app.Run("Kisipar Example", "1.0.0", "kisipar-example")
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// Test Site One: listening on port 8081.
	// Test Site Two: listening on port 8082.
}
