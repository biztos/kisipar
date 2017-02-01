// standardprovider.go -- the StandardProvider type and frienp.
// -------------------
// TODO: MetaAsset or something like that, so you can have a set of pics
// with titles and metas.  I want that for slideshows anyway, and they
// will be way easier to test from a virtual provider.
// ** StandardAsset? **
// ALSO: some concept of sort order for them, maybe?
package kisipar

import (
	// Standard Library:
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	// Third-party packages:
	"gopkg.in/yaml.v2"
)

// BasicStub is a minimal immutable non-Page Stub.
type BasicStub struct {
	path string
}

// Path returns the request path of the stub.
func (s *BasicStub) Path() string { return s.path }

// TypeString returns the stringified type of the stub.
func (s *BasicStub) TypeString() string { return fmt.Sprintf("%T", s) }

// IsPageStub returns false for the BasicStub; use the StandardPageStub for
// a simple page-based stub.
func (s *BasicStub) IsPageStub() bool { return false }

// NewBasicStub returns a pointer to a BasicStub with the given path.
func NewBasicStub(rpath string) *BasicStub { return &BasicStub{rpath} }

// StandardPageStub is a stub based on StandardPage.  This means that every
// StandardPageStub is both a Stub and a Page.
type StandardPageStub StandardPage

// IsPage always returns true for a StandardPageStub.
func (sps *StandardPageStub) IsPage() bool { return true }

// StandardPage is an immutable Page that exists entirely in memory.
type StandardPage struct {
	path    string
	title   string
	tags    []string
	created time.Time
	updated time.Time
	meta    map[string]interface{}
	html    string
}

// Stub returns a StandardPageStub from the StandardPage.
func (p *StandardPage) Stub() *StandardPageStub {
	return &StandardPageStub{
		path:    p.Path(),
		title:   p.Title(),
		tags:    p.Tags(),
		created: p.Created(),
		updated: p.Updated(),
		meta:    p.Meta(),
		html:    "",
	}
}

// Path returns the canonical request path of the page.
func (p *StandardPage) Path() string {
	return p.path
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

// Meta retuns the full Meta Map.  It may be nil.
func (p *StandardPage) Meta() map[string]interface{} {
	return p.meta
}

// MetaString returns a string value from the page's Meta Map for the given
// key.  Lookup is case-sensitive.  The value is stringified per %v in
// fmt.Sprintf. If the mapped value or the map itself is nil then the empty
// string is returned.
func (p *StandardPage) MetaString(key string) string {
	if p.meta == nil {
		return ""
	}
	val := p.meta[key]
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

// MetaStrings returns a slice of string values from the page's Meta Map. If
// the value is already a []string, that is returned; if it is a slice then
// each value is stringified as in MetaString and that slice of strings
// returned; if it is a single value, that value is stringified via
// MetaString and returned in a slice of one; and if the value is nil (or the
// map itself is nil) an empty slice is returned.
func (p *StandardPage) MetaStrings(key string) []string {
	if p.meta == nil {
		return []string{}
	}
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
func NewStandardPage(title string, tags []string, created time.Time,
	updated time.Time, meta map[string]interface{}) *StandardPage {

	return &StandardPage{
		title:   title,
		tags:    tags,
		created: created,
		updated: updated,
		meta:    meta,
	}
}

// StandardPageFromData returns a pointer to a StandardPage with its internal
// properties set according to the key-value pairs in the provided data map d.
// All properties are optional, and unknown properties are ignored, allowing
// the data map to be used for other things.  An error is returned if the
// value has the wrong type.
//
// Keys should be lowercase:
//
//  map[string]interface{}{
//      "title": "Hello World",
//      "tags": []string{"foo","bar"},
//      "created": time.Now(),
//      "updated": time.Now(),
//      "meta": map[string]interface{}{"foo":"bar"},
//  }
func StandardPageFromData(d map[string]interface{}) (*StandardPage, error) {

	p := &StandardPage{}
	for k, v := range d {
		switch k {
		case "title":
			if val, ok := v.(string); ok {
				p.title = val
			} else {
				return nil, wrongTypeError(k, v, "string")
			}
		case "tags":
			if val, ok := v.([]string); ok {
				p.tags = val
			} else {
				return nil, wrongTypeError(k, v, "string slice")
			}
		case "created":
			if val, ok := v.(time.Time); ok {
				p.created = val
			} else {
				return nil, wrongTypeError(k, v, "Time")
			}
		case "updated":
			if val, ok := v.(time.Time); ok {
				p.updated = val
			} else {
				return nil, wrongTypeError(k, v, "Time")
			}
		case "meta":
			if val, ok := v.(map[string]interface{}); ok {
				p.meta = val
			} else {
				return nil, wrongTypeError(k, v, "string-interface map")
			}
		default:
			// Ignore unknown keys, they're harmless.
		}

	}

	return p, nil
}

func wrongTypeError(k string, v interface{}, want string) error {
	return fmt.Errorf("Wrong type for %s: %T not %s", k, v, want)
}

// StandardProvider is an opaque Provider that exists entirely in memory.
// It can be used as the base for any other type of Provider that does not
// need special features.  Items are added using the interface-specific
// Add* methods, and removed with the Remove method; these are safe for
// use by concurrent goroutines.
type StandardProvider struct {
	items    map[string]Pather
	modtimes map[string]time.Time
	template *template.Template
	created  time.Time
	updated  time.Time
	mutex    sync.RWMutex
}

// NewStandardProvider returns an empty StandardProvider to be populated
// via AddPage et al.  Directly creating a StandardProvider may lead to
// runtime errors in its methop; use this function instead.
func NewStandardProvider() *StandardProvider {
	return &StandardProvider{
		items:    map[string]Pather{},
		modtimes: map[string]time.Time{},
		created:  time.Now(),
		updated:  time.Now(),
		mutex:    sync.RWMutex{},
	}
}

// Add adds a Pather of any kind to the StandardProvider at its Path.  If any
// other item exists at that path it will be overridden.  Add is safe for
// concurrent use.
func (sp *StandardProvider) Add(p Pather) {

	modtime := time.Now()
	if f, ok := p.(File); ok {
		if info, err := os.Stat(f.FilePath()); err == nil {
			modtime = info.ModTime()
		}
	}

	sp.mutex.Lock()

	sp.items[p.Path()] = p
	sp.modtimes[p.Path()] = modtime
	sp.updated = time.Now()

	sp.mutex.Unlock()
}

// StandardProviderFromYaml returns a StandardProvider with its pages and
// data read from the supplied yaml string.  The structure should be:
//    pages:
//      /path/to/foo:
// TODO!
//
// This is useful for testing and for placeholder and/or generated sites
// with text-only content.
func StandardProviderFromYaml(in string) (*StandardProvider, error) {

	meta := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(in), &meta)
	if err != nil {
		return nil, err
	}
	p := NewStandardProvider()

	// TODO: pages, etc... all items...

	return p, nil
}

// String returns a log-friendly description of the Provider.
func (sp *StandardProvider) String() string {
	return fmt.Sprintf("<StandardProvider with %d items, updated %s>",
		len(sp.items), sp.updated)
}

// Get returns an item from the Provider if available; an error if not.
func (sp *StandardProvider) Get(rpath string) (Pather, error) {

	sp.mutex.RLock()
	item := sp.items[rpath]
	sp.mutex.RUnlock()

	if item == nil {
		return nil, ErrNotExist
	}
	return item, nil

}

// GetSince returns an item from the Provider if available and more recent
// than the since time; or an error if not. The caller should in particular
// watch for the ErrNotModified error, which usually indicates a 304 Not
// Modified response should be sent to the client.
func (sp *StandardProvider) GetSince(rpath string, since time.Time) (Pather, error) {

	sp.mutex.RLock()
	item := sp.items[rpath]
	modtime := sp.modtimes[rpath]
	sp.mutex.RUnlock()

	if item == nil {
		return nil, ErrNotExist
	}
	if modtime.After(since) {
		return item, nil
	}
	return nil, ErrNotModified
}

// GetStub returns a Stub item if available and implements the Stubber
// interface; an error if not.
func (sp *StandardProvider) GetStub(rpath string) (Stub, error) {
	sp.mutex.RLock()
	item, err := sp.Get(rpath)
	sp.mutex.RUnlock()

	if err != nil {
		return nil, err
	}
	if s, ok := item.(Stubber); ok {
		return s.Stub(), nil
	}
	return nil, ErrNotStubber

}

// GetUnder returns a slice of Stub items for everything "under" the given
// prefix, i.e. everything with a path having the given prefix.  It is
// usually wise to terminate the prefix with a slash, but this is up to the
// template author.
//
// TODO: new tye for []Stub, an interface, something that can deal with
// iterators and so on.  TBD really.
func (sp *StandardProvider) GetUnder(prefix string) ([]Stub, error) {

	sp.mutex.RLock()
	stubs := []Stub{}
	for path, item := range sp.items {
		if strings.HasPrefix(path, prefix) {
			if s, ok := item.(Stubber); ok {
				stubs = append(stubs, s.Stub())
			}
		}
	}
	sp.mutex.RUnlock()

	return stubs, nil

}

// Template returns the top-level template collection for the
// StandardProvider. It may be nil.
func (sp *StandardProvider) Template() *template.Template {
	return sp.template
}

// TemplateFor returns the template that should be used for the given
// Page.  It may be nil.  Template-selection logic is up to the Provider,
// and in this case it is very simple:
//   * The "template" MetaString, or
//   * The "Template" MetaString, or
//   * The template with a matching path, or
//   * The template with the longest matching basedir to the Page's basedir.
//
// Only internal templates are considered.  It is up to the handler to choose
// among external templates, though the function is exposed for reuse: cf.
// PageTemplate.
func (sp *StandardProvider) TemplateFor(p Page) *template.Template {
	return PageTemplate(sp.template, p)
}

// TODO: MOVE TO TEMPLATES OR MAYBE TO PROVIDER:
func PageTemplate(tmpl *template.Template, p Page) *template.Template {
	if tmpl == nil {
		return nil
	}
	// var tpath string
	// ...
	return nil

}

// TODO AND TO MOVE:
type FileProviderConfig struct {
	RootDir           string
	TemplateDir       string
	ExcludeExtensions []string
	AutoRefresh       bool
}

// FileProvider is a Provider using only the file system.  It is the
// recommended Provider to use for developing templates, and is also
// useful for set-and-forget sites such as placeholders or smaller archives.
// It is NOT recommended for large sites or any case where performance is key.
type FileProvider struct {
	root string
}

// NewFileProvider returns a FileProvider intialized from the given
// root path.
func NewFileProvider(root string) *FileProvider {
	return &FileProvider{}
}

// String returns a log-friendly description of the Provider.
func (p *FileProvider) String() string {
	return fmt.Sprintf("<FileProvider in %s>", p.root)
}
