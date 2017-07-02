// dot.go -- The Dot, passed to Kisipar templates.

package site

import (
	"net/http"

	// Own stuff:
	"github.com/biztos/kisipar/provider"
)

// Dot is the structure passed to all Kisipar templates for execution (it is
// usually addressed as "." in the templates, hence the name).
type Dot struct {
	Request  *http.Request     // The Request being replied to.
	Config   *Config           // The site configuration.
	Page     provider.Page     // The Page, if there is one, or:
	Stubs    []provider.Stub   // The slice of Stubs if there isn't.
	Provider provider.Provider // The Provider, for queries etc.
}
