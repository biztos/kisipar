// standardprovider_examples_test.go

package provider_test

import (
	"fmt"

	"github.com/biztos/kisipar/provider"
)

func ExampleStandardProvider() {

	ds, err := provider.StandardProviderFromYAML("# nothing yet")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v", ds)

}
