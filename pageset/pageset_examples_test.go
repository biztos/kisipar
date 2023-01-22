// pageset/pageset_examples_test.go - examples for the Kisipar Pageset
// --------------------------------

package pageset_test

import (
	"fmt"
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
)

func Example() {

	// Given a set of pages:
	p1, _ := page.LoadVirtualString("/here/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/here/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/there/d.md", "# Third!")

	// ...we create a Pageset:
	ps, err := pageset.New([]*page.Page{p1, p2, p3})
	if err != nil {
		panic(err)
	}

	// ...which then can give us back its Pages in various ways:
	fmt.Println(ps.Page("/here/a").Title())
	for i, p := range ps.ByModTime() {
		fmt.Println(i, p.Title())
	}
	// Output:
	// First!
	// 0 Third!
	// 1 Second!
	// 2 First!
}

func ExamplePageset_Subset() {

	// Given a Pageset:
	p1, _ := page.LoadVirtualString("my/pages/here/a.md", "# First!")
	p2, _ := page.LoadVirtualString("my/pages/here/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("my/pages/there/d.md", "# Third!")
	ps, err := pageset.New([]*page.Page{p1, p2, p3})
	if err != nil {
		panic(err)
	}

	// ...we commonly extract a set below a specific path:
	subset := ps.PathSubset("/here", "my/pages")
	fmt.Println(subset.Len())

	// ...upon which we can perform our normal sorts and selects:
	for i, p := range subset.ByPath() {
		fmt.Println(i, p.Title())
	}
	// Output:
	// 2
	// 0 First!
	// 1 Second!

}
