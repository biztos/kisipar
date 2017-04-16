// provider.go -- the Provider interface et al.
// -----------

package kisipar

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

var (

	// Not-realy-error errors:
	ErrNotExist    = errors.New("item does not exist")
	ErrNotModified = errors.New("not modified")
	ErrNotStubber  = errors.New("item is not a stubber")
)

// IsNotExist returns true if err matches ErrNotExist or certain other known
// errors that may be passed through from the DataSource.
func IsNotExist(err error) bool {
	if err == ErrNotExist {
		return true
	}
	if os.IsNotExist(err) {
		return true
	}
	return false
}

// Provider provides items of any type that correspond to request paths.
// It may also provide templates.
//
// TODO: the following can't be gauranteed so maybe a SortedStubs method
// to make a copy of the array and sort it?
//   Slices returned from batch requests
//   should be backed by new arrays and thus safe to sort.
type Provider interface {

	// Meta:
	String() string // For logging / debugging.

	// Single requests:
	Get(string) (Pather, error)                 // Fetch item at path.
	GetSince(string, time.Time) (Pather, error) // If-Modified-Since.
	GetStub(string) (Stub, error)               // If one Stub needed.

	// Batch requests:
	GetStubs(prefix string) []Stub // Fetch Stubs under a prefix.
	GetAll(prefix string) []Pather // Fetch Stubs or raw items.

	// TODO:
	// GetTagged
	// Find  (with any interface)
	// ...?

	// Templates:
	Template() *template.Template
	TemplateFor(Page) *template.Template
}

// Pather is the minimum interface for a Provider item.  Every item must
// know its request path.  An item that satisfies no more specific interface
// will normally be served as a 204 No Content response.
type Pather interface {
	Path() string
}

// Stubber is the interface for generating a Stub from an item within a
// Provider.
type Stubber interface {
	Stub() Stub
}

// Redirect is a Provider item indicating that the client should be
// redirected to another location. The supported HTTP redirect responses are
// 301 Moved Permanently and 302 Found, with Permanent indicating which to
// use.
type Redirect interface {
	Path() string
	Permanent() bool
	Location() string
}

// File is a Provider item that points to a file on disk, to be sent to the
// client through http.ServeFile. It is useful for local caching of larger
// and/or immutable assets such as images.
type File interface {
	Path() string
	FilePath() string
}

// Content is a Provider item used to serve arbitrary content through
// http.ServeContent.
type Content interface {
	Path() string
	ContentType() string
	ModTime() time.Time
	ReadSeeker() io.ReadSeeker
}

// Handler defines a custom handler for a response, providing for full
// customization of responses within the Provider itself.  The standard
// use-case for this is writing to the Provider through POST requests.
type Handler interface {
	Path() string
	Func() func(http.ResponseWriter, *http.Request)
}

// Page is a single content page renderable as an HTML fragment in a template.
// A Page may (and often will) provide many more methods for use in templates;
// the interface defines the minimum required set used by default templates
// and in logging.
type Page interface {
	Path() string                 // Request path for Pather interface.
	Title() string                // Title of the Page.
	Tags() []string               // List of Tags applicable to the Page.
	Created() time.Time           // Creation time (IsZero if N/A).
	Updated() time.Time           // Update time (IsZero if N/A).
	Meta() map[string]interface{} // Available metadata for the Page.
	MetaString(string) string     // Return a Meta value as a string.
	MetaStrings(string) []string  // Return a Meta value as a string slice.
	HTML() template.HTML          // Rendered HTML fragment of the Page.
}

// Stub is included exclusively in lists (including lists of one).  The
// template should have an expection of what sort of items will be stubbed
// for a given list, as the Stub will normally provide additional methods
// useful in a listing items.
type Stub interface {
	Path() string       // Request path for Pather interface.
	IsPageStub() bool   // True if it's really the stub of a Page.
	TypeString() string // Type as string, useful to templates.
}

// PageStub is a more detailed Stub, to be used for Pages. In practice this
// is usually implemented as an alias to the page type. A PageStub may be
// used for non-Page types that provide page-like metadata and content, e.g.
// images or map entries. The IsPageStub method should return true if the
// PageStub is backing a traditional Page, and should be trusted over
// interface checks.
//
// TODO: consider ditching IsPageStub in favor of something else, like maybe
// HasPage or IsTruePage in order to avoid confusion.
type PageStub interface {
	Path() string                 // Request path for Pather interface.
	Title() string                // Title of the Page.
	Tags() []string               // List of Tags applicable to the Page.
	Created() time.Time           // Creation time (IsZero if N/A).
	Updated() time.Time           // Update time (IsZero if N/A).
	Meta() map[string]interface{} // Available metadata for the Page.
	MetaString(string) string     // Return a Meta value as a string.
	MetaStrings(string) []string  // Return a Meta value as a string slice.
	IsPageStub() bool             // True if it's really the stub of a Page,
	TypeString() string           // Type as string, useful to templates.
}
