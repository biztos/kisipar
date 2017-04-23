// standardprovider_examples_test.go

package kisipar_test

import (
	"fmt"

	"github.com/biztos/kisipar"
)

func ExampleStandardProvider() {

	ds, err := kisipar.StandardProviderFromYAML("# nothing yet")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v", ds)

}
