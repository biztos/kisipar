// handler_examples_test.go - examples for http handlers
// ------------------------

package site_test

import (
	"github.com/biztos/kisipar/site"

	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleSite_NewServeMux() {

	// Create a Site with a custom override.  We start with a standard site,
	// in this case an empty one:
	s, err := site.New("")
	if err != nil {
		log.Fatal(err)
	}

	// This is the override we want: /world -> HELLO CRUEL WORLD
	h := &site.PatternHandler{
		Pattern: "/world",
		Function: func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "HELLO CRUEL WORLD\n")
		}}

	// We already have a ServeMux set up as the Server's Handler in our Site,
	// but we don't know what's in it and we don't want to put in something
	// that conflicts with the existing multiplexer. So we give it a new
	// one built with our override.
	mux := s.NewServeMux(h)

	// Here we attach it back to the Site's Server, though we could of course
	// go do something else with it, such as serve a very Kisipar-like site
	// from some other server framework.
	s.SetHandler(mux)

	// Let's prove our Mux is doing what we want:
	req, err := http.NewRequest("GET", "http://example.com/world", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// HELLO CRUEL WORLD
}
