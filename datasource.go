// datasource.go
// -------------
// TODO:
// * some logic for 304 not modified esp. since the cache is in the DS.
//   (give it a TS?  or what? ideally make it optional...)
package kisipar

import (
	// Standard Library:
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"time"

	// Third-party packages:
	"gopkg.in/yaml.v2"
)

type ItemType string

// TODO: consider some better enum approach...
const (
	NoItem      ItemType = ""
	PageItem    ItemType = "Page"
	DataItem    ItemType = "Data"
	FileItem    ItemType = "File"
	RedirItem   ItemType = "Redir"
	HandlerItem ItemType = "Handler"
)

var (

	// Not-realy-error errors:
	ErrNotExist    = errors.New("item does not exist")
	ErrNotModified = errors.New("not modified")
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

// Page defines a single page renderable as an HTML fragment.  It is
// typically returned from a DataSource and made available in a template.
// A Page may (and often will) provide many more methods for use in templates;
// the interface defines the minimum required set used by default templates
// and in logging.
type Page interface {
	Id() string                   // Unique ID of the Page; empty if N/A.
	Title() string                // Title of the Page.
	Tags() []string               // List of Tags applicable to the Page.
	Created() time.Time           // Creation time (IsZero if N/A).
	Updated() time.Time           // Update time (IsZero if N/A).
	Meta() map[string]interface{} // Available metadata for the Page.
	MetaString(string) string     // Return a Meta value as a string.
	MetaStrings(string) []string  // Return a Meta value as a string slice.
	HTML() string                 // Rendered HTML fragment of the Page.
}

// Data defines an arbitrary set of bytes with a content-type as returned
// by the DataSource.
// TODO: io.Reader instead of bytes?  Might be useful for blobs...
type Data interface {
	ContentType() string
	Bytes() []byte
}

// File defines a file on disk which can be further processed, e.g. by the
// http.ServeFile handler.
type File interface {
	Path() string // path on disk
}

// Redir defines a redirection to send to the client.  The supported types are
// 301 Moved Permanently and 302 Found, with Permanent indicating which to
// use.
type Redir interface {
	Permanent() bool
	Location() string
}

// Handler defines a custom handler for a response, and provides a system
// for extending response handling at the DataSource level.
type Handler http.Handler

// DataSource is a provider of the items available at request paths.
// For the default handlers, the items must be limited to those listed in
// the ItemType constants.  If nothing is found for a given request, the
// error returned by the method should be ErrNotExist or another that
// satisfies IsNotExist; other errors are treated as Internal Server Errors.
//
// A DataSource may be entirely virtual (cf. VitualDataSource) but in most
// cases it is an interface to a persistent data store of some type:
// database, cloud environment, cache, filesystem or all of the above.
//
// Optimization of access is the responsibility of the DataSource.  In
// particular, a useful DataSource should respond very quickly to Has
// calls, e.g. by caching in memory all known resources.  In normal operation,
// Has is called first, and then a type-specific method such as Page will be
// called if applicable.  Thus for every request, including DDoS and probing,
// the DataSource's Has method will be called once.
//
// Template service is optional, and enabled or disabled in the site
// configuration.
type DataSource interface {
	// Meta:
	String() string // For logging / debugging.

	// Single request responses:
	Has(rpath string) ItemType
	GetPage(rpath string) (Page, error)
	GetData(rpath string) (Data, error)
	GetFile(rpath string) (File, error)
	GetRedir(rpath string) (Redir, error)
	GetHandler(rpath string) (Handler, error)

	// Batches:
	GetPages(prefix string) ([]Page, error) // Pages under a prefix

	// Templates:
	Template(time.Time) (*template.Template, error)
	TemplateFor(rpath string) (*template.Template, error)
}

// StandardPage is an immutable Page that exists entirely in memory.
// It is the preferred Page implementation for all standard uses.
type StandardPage struct {
	id      string
	title   string
	tags    []string
	created time.Time
	updated time.Time
	meta    map[string]interface{}
	html    string
}

// Id returns the (possibly unique) identifier of the page.
func (p *StandardPage) Id() string {
	return p.id
}

// Title returns the title of the page.
func (p *StandardPage) Title() string {
	return p.title
}

// Tags returns the tags of the page.  Order is not guaranteed.
func (p *StandardPage) Tags() []string {
	return p.tags
}

// Created returns the creation timestamp of the page.  The zero time is
// the equivalent of nil.
func (p *StandardPage) Created() time.Time {
	return p.created
}

// Updated returns the update timestamp of the page.  The zero time is
// the equivalent of nil.
func (p *StandardPage) Updated() time.Time {
	return p.updated
}

// MetaString returns a string value from the page's Meta Map for the given
// key.  Lookup is case-sensitive.  The value is stringified per %v in
// fmt.Sprintf. If the mapped value is nil then the empty string is returned.
func (p *StandardPage) MetaString(key string) string {
	val := p.meta[key]
	return fmt.Sprintf("%v", val)
}

// MetaStrings returns a slice of string values from the page's Meta Map. If
// the value is already a []string, that is returned; if it is a slice then
// each value is stringified as in MetaString and that slice of strings
// returned; if it is a single value, that value is stringified via
// MetaString and returned in a slice of one; and if the value is nil an
// empty slice is returned.
func (p *StandardPage) MetaStrings(key string) []string {
	val := p.meta[key]
	if val == nil {
		return []string{}
	}
	if s, ok := val.([]string); ok {
		return s
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:

		slice := reflect.ValueOf(val)
		s := make([]string, slice.Len())

		for i := 0; i < slice.Len(); i++ {
			s[i] = fmt.Sprintf("%v", slice.Index(i))
		}
		return s
	default:
		return []string{p.MetaString(key)}
	}

}

// NewStandardPage returns a pointer to a StandardPage with its internal
// properties set according to the arguments.
func NewStandardPage(id, title string, tags []string, created time.Time,
	updated time.Time, meta map[string]interface{}) *StandardPage {

	return &StandardPage{
		id:      id,
		title:   title,
		tags:    tags,
		created: created,
		updated: updated,
		meta:    meta,
	}
}

// StandardPageFromYaml returns a pointer to an initialized StandardPage
// with its Id, Title, and Tags set from meta using the MetaString and
// MetaStrings methods; and the Created and Updated values set if they are
// present in meta and are of type time.Time. These should be keyed in
// lowercase:
//
//  map[string]interface{}{
//      "id": "possibly-unique",
//      "title": "Hello World",
//      "tags": []string{"foo","bar"},
//      "created": time.Now(),
//  }
func StandardPageFromYaml(html string, meta map[string]interface{}) *StandardPage {

	p := &StandardPage{
		meta: meta,
		html: html,
	}

	p.id = p.MetaString("id")
	p.title = p.MetaString("title")
	p.tags = p.MetaStrings("tags")

	if val := meta["created"]; val != nil {
		if t, ok := val.(time.Time); ok {
			p.created = t
		}
	}
	if val := meta["udated"]; val != nil {
		if t, ok := val.(time.Time); ok {
			p.updated = t
		}
	}

	return p

}

// VirtualDataSource is a DataSource that exists entirely in memory.  It
// is primarily useful for testing, but might have other uses as well.
// (What might those be? Dynamically creating a site based on read-once data?
// Making a placeholder that exists entirely as a config file? Remotely
// updating a small site via an API?)
type VirtualDataSource struct {
	// or just... items map[string]interface{}? We add/delete via api...
	pages    map[string]*StandardPage
	data     map[string]*StandardPage
	files    map[string]*StandardPage
	redirs   map[string]*StandardPage
	handlers map[string]*StandardPage
	// TODO: templates!
	created time.Time
}

// NewVirtualDataSource returns an empty VirtualDataSource to be populated
// with the AddData and AddPage methods.
func NewVirtualDataSource() *VirtualDataSource {
	return &VirtualDataSource{
		pages:    map[string]*StandardPage{},
		data:     map[string]*StandardPage{},
		files:    map[string]*StandardPage{},
		redirs:   map[string]*StandardPage{},
		handlers: map[string]*StandardPage{},
		created:  time.Now(),
	}
}

// VirtualDataSourceFromYaml returns a VirtualDataSource with its pages and
// data read from the supplied yaml string.  The structure should be:
//    pages:
//      /path/to/foo:
// TODO!
//
// This is useful for testing and for placeholder and/or generated sites
// with text-only content.
func VirtualDataSourceFromYaml(in string) (*VirtualDataSource, error) {

	meta := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(in), &meta)
	if err != nil {
		return nil, err
	}
	ds := NewVirtualDataSource()

	// TODO: pages, etc... all items...

	return ds, nil
}

// String returns a log-friendly description of the DataSource.
func (ds *VirtualDataSource) String() string {
	return fmt.Sprintf("<VirtualDataSource: %d pages, %d data>",
		len(ds.pages), len(ds.data))
}

// Has returns the ItemType of the path within the DataSource.
func (ds *VirtualDataSource) Has(rpath string) ItemType {
	//
	return NoItem
}

// GetPage returns a Page if available from the DataSource; an error if not.
func (ds *VirtualDataSource) GetPage(rpath string) (Page, error) {
	return nil, ErrNotExist
}

// GetData returns a Data if available from the DataSource; an error if not.
func (ds *VirtualDataSource) GetData(rpath string) (Data, error) {
	return nil, ErrNotExist
}

// GetFile returns a File if available from the DataSource; an error if not.
func (ds *VirtualDataSource) GetFile(rpath string) (File, error) {
	return nil, ErrNotExist
}

// GetRedir returns a Redir if available from the DataSource; an error if not.
func (ds *VirtualDataSource) GetRedir(rpath string) (Redir, error) {
	return nil, ErrNotExist
}

// GetHandler returns a Handler if available from the DataSource; an error if
// not.
func (ds *VirtualDataSource) GetHandler(rpath string) (Handler, error) {
	return nil, ErrNotExist
}

// GetPages returns a slice of Page items whose paths have the given prefix,
// or ErrNotExist if none have.
func (ds *VirtualDataSource) GetPages(rpath string) ([]Page, error) {
	return nil, ErrNotExist
}

// Template compiles and returns the template collection if it has changed
// since last.  If it has not, ErrNotModified is returned.
func (ds *VirtualDataSource) Template(last time.Time) (*template.Template, error) {
	return nil, ErrNotModified
}

// TemplateFor returns the template to be used for the given path, or an
// error (possibly ErrNotExist, which means use the default).
func (ds *VirtualDataSource) TemplateFor(rpath string) (*template.Template, error) {
	return nil, ErrNotExist
}

// FileDataSource is a DataSource using only the file system.  It is the
// recommended DataSource to use for developing templates, and is also
// useful for set-and-forget sites such as placeholders or smaller archives.
// It is NOT recommended for large sites or any case where performance is key.
type FileDataSource struct {
	root string
}

// NewFileDataSource returns a FileDataSource intialized from the given
// root path.
func NewFileDataSource(root string) *FileDataSource {
	return &FileDataSource{}
}

// String returns a log-friendly description of the DataSource.
func (ds *FileDataSource) String() string {
	return fmt.Sprintf("<FileDataSource in %s>", ds.root)
}
