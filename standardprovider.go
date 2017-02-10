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
	"errors"
	"fmt"
	"html/template"
	"io"
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

// StandardContent is an immutable piece of arbitrary content stored as a
// string internally, with no metadata other than the ContentType.  It
// implements the Content interface.
type StandardContent struct {
	rpath   string
	ctype   string
	mtime   time.Time
	content string
}

// Path returns the request path of the content.
func (c *StandardContent) Path() string {
	return c.rpath
}

// ContentType returns the content-type of the content.  If it is the empty
// string, the server will guess, which is likely to be inefficient.
func (c *StandardContent) ContentType() string {
	return c.ctype
}

// ModTime returns the modificatio time of the content. If it is the zero
// time then the current time is returned.
func (c *StandardContent) ModTime() time.Time {
	if c.mtime.IsZero() {
		return time.Now()
	}
	return c.mtime
}

// ReadSeeker returns an io.ReadSeeker created from the content string of the
// content item.
func (c *StandardContent) ReadSeeker() io.ReadSeeker {
	return strings.NewReader(c.content)
}

// NewStandardContent returns a pointer to a StandardContent item with the
// given request path, content type, content string and mod time.
func NewStandardContent(rpath, ctype, content string, mtime time.Time) *StandardContent {
	return &StandardContent{
		rpath:   rpath,
		ctype:   ctype,
		content: content,
		mtime:   mtime,
	}
}

// StandardPage is an immutable Page that exists entirely in memory.
type StandardPage struct {
	rpath   string
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
		rpath:   p.Path(),
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
	return p.rpath
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

// HTML retuns the rendered HTML of the StandardPage. It may be empty.
// Note that StandardPage does not have any rendering logic of its own.
func (p *StandardPage) HTML() string {
	return p.html
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
func NewStandardPage(rpath, title string, tags []string,
	created, updated time.Time, meta map[string]interface{},
	html string) *StandardPage {

	return &StandardPage{
		rpath:   rpath,
		title:   title,
		tags:    tags,
		created: created,
		updated: updated,
		meta:    meta,
		html:    html,
	}
}

// StandardPageFromData returns a pointer to a StandardPage with its
// internal properties set according to the key-value pairs in the provided
// data map m. All properties are optional except for path, and unknown
// properties are ignored, allowing the data map to be used for other
// things. An error is returned if any value has the wrong type, or if path
// is either undefined or the empty string.
//
// Keys should be lowercase:
//
//  map[string]interface{}{
//      "path": "/foo/bar",
//      "title": "Hello World",
//      "tags": []string{"foo","bar"},
//      "created": time.Now(),
//      "updated": time.Now(),
//      "meta": map[string]interface{}{"foo":"bar"},
//  }
func StandardPageFromData(m map[string]interface{}) (*StandardPage, error) {

	p := &StandardPage{}

	path := m["path"]
	if path == nil {
		return nil, errors.New("path not set")
	}
	if val, ok := path.(string); ok {
		p.rpath = val
	} else {
		return nil, wrongTypeError("path", path, "string")
	}
	if p.rpath == "" {
		return nil, errors.New("path may not be an empty string")
	}

	// The rest are optional.
	for k, v := range m {
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

// StandardFile is a simple implementation of the File interface.
type StandardFile struct {
	rpath string
	fpath string
}

// Path returns the request path corresponding to the file.
func (sf *StandardFile) Path() string { return sf.rpath }

// FilePath returns the path to the file on the file system.
func (sf *StandardFile) FilePath() string { return sf.fpath }

// NewStandardFile returns a StandardFile for the given request path and
// filesystem path.  The file itself is not checked.
func NewStandardFile(rpath, fpath string) *StandardFile {
	return &StandardFile{
		rpath: rpath,
		fpath: fpath,
	}
}

// StandardProvider is a mostly-opaque Provider that exists entirely in
// memory.
//
// It can be used as the base for any other type of Provider that does not
// need special features within the storage and retrieval of items.  Items
// are added using the Add method, and removed with the Remove method; these
// are safe for use by concurrent goroutines.
//
// The one exposed property, Meta, allows for other types based on this type
// to store arbitrary data (typically a configuration).  It is not used by
// any StandardProvider methods.
type StandardProvider struct {
	Meta     map[string]interface{}
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

// StandardProviderFromYaml returns a StandardProvider with Pages, Content and
// templates read from the supplied YAML input.  The YAML structure is:
//    pages:
//      /path/to/foo:
//         title: I am Page Foo
//         tags: [foo,yaml,demo]
//         meta:
//             foo: true
//         html: |
//              I am content!
//              Multi-line!
//    content:
//      /path/to/bar.js:
//         type: application/javascript
//         content: |
//              var x = 'Hello!';
//              window.alert(x);
//    templates:
//      /some/template/path: |
//          Hello {{ .Title }}
//
// The date format is not flexible.  Pages and Content with the same path
// will be overwritten, most likely in a random order.
//
// In case of parse errors, the first error encountered is returned as-is.
//
// This is useful for testing and for placeholder and/or generated sites
// with simple content.  It is NOT recommended complex scenarios.
func StandardProviderFromYaml(src string) (*StandardProvider, error) {

	type pageFromYaml struct {
		Title   string
		Tags    []string
		Created time.Time
		Updated time.Time
		Meta    map[string]interface{}
		Html    string
	}
	type contentFromYaml struct {
		Type    string
		Content string
	}
	type fromYaml struct {
		Pages     map[string]*pageFromYaml
		Content   map[string]*contentFromYaml
		Templates map[string]string
	}
	target := fromYaml{}

	if err := yaml.Unmarshal([]byte(src), &target); err != nil {
		return nil, err
	}

	sp := NewStandardProvider()
	for path, item := range target.Pages {
		p := NewStandardPage(
			path,
			item.Title,
			item.Tags,
			item.Created,
			item.Updated,
			item.Meta,
			item.Html,
		)

		sp.Add(p)
	}
	for path, item := range target.Content {
		c := NewStandardContent(path, item.Type, item.Content, time.Now())
		sp.Add(c)
	}

	// Set up any templates.
	if target.Templates != nil {
		tmpl, err := TemplatesFromData(target.Templates)
		if err != nil {
			return nil, err
		}
		sp.SetTemplate(tmpl)
	}

	return sp, nil

}

// String returns a log-friendly description of the Provider.
func (sp *StandardProvider) String() string {

	var s string
	if len(sp.items) != 1 {
		s = "s"
	}
	return fmt.Sprintf("<StandardProvider with %d item%s, updated %s>",
		len(sp.items), s, sp.updated)
}

// Count returns the total number of items in the StandardProvider.
func (sp *StandardProvider) Count() int {
	return len(sp.items)
}

// Updated returns the last time at which the StandardProvider was updated.
func (sp *StandardProvider) Updated() time.Time {
	return sp.updated
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
// Page.  It is a convenience wrapper for the PageTemplate function.
func (sp *StandardProvider) TemplateFor(p Page) *template.Template {
	return PageTemplate(sp.template, p)
}

// SetTemplate sets the internal template returned by Template.
func (sp *StandardProvider) SetTemplate(tmpl *template.Template) {
	sp.mutex.Lock()
	sp.template = tmpl
	sp.updated = time.Now()
	sp.mutex.Unlock()
}
