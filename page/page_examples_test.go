// page/page_examples_test.go -- examples for Page documentation.
// --------------------------

package page_test

import (
	"fmt"
	"github.com/biztos/kisipar/page"
)

func Example() {

	// Most common use case: load a Markdown file using the default parser.
	p, err := page.Load("example.md")
	if err != nil {
		panic(err)
	}

	fmt.Println("Title:", p.Title())
	fmt.Println("Author:", p.Author())
	fmt.Println("Hamand:", p.MetaString("Hamand"))
	fmt.Println(p.Content)

	// Output:
	// Title: Example Kisipar Page
	// Author: Thelonius Mönch
	// Hamand: EGGS!
	// <h1>Example Kisipar Page</h1>
	//
	// <h2>Welcome to the Example Page!</h2>
	//
	// <p>A paragraph is all you need.</p>
}
