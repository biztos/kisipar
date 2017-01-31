// datasource_examples_test.go

package kisipar_test

import (
	"fmt"

	"github.com/biztos/kisipar"
)

func ExampleStandardDataSource() {

	ds, err := kisipar.StandardDataSourceFromYaml("# nothing yet")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v", ds)

}
