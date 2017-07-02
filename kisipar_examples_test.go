// kisipar_examples_test.go

package kisipar_test

import (
	"fmt"
	"path/filepath"

	"github.com/biztos/kisipar"
)

func Example() {

	// TODO: a real example, probably with Serve, meaning we'd have to come
	// up with a way for that to be testable.
	site, err := kisipar.NewSite(filepath.Join("testdata", "config.yaml"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(site.Config.Name)

	// Output:
	// Kisipar Test Site

}
