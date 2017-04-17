// dot.go -- The Dot, passed to Kisipar templates.

package kisipar

import (
	"net/http"
)

// Dot is the structure passed to all Kisipar templates for execution (it is
// usually addressed as "." in the templates, hence the name).
type Dot struct {
	Request  *http.Request // The Request being replied to.
	Config   *Config       // The site configuration.
	Page     Page          // The Page, if there is one, or:
	Stubs    []Stub        // The slice of Stubs if there isn't.
	Provider Provider      // The Provider, for queries etc.
}
