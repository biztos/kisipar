// datasource_examples_test.go

package kisipar_test

import (
	"fmt"

	"github.com/biztos/kisipar"
)

func ExampleVirtualDataSource() {

	ds, err := kisipar.VirtualDataSourceFromYaml("# nothing yet")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v", ds)

	// Output:
	// <VirtualDataSource: 0 pages, 0 data>

}
