// standardprovider.go -- the StandardProvider type and friends.
// -------------------
// TODO: MetaAsset or something like that, so you can have a set of pics
// with titles and metas.  I want that for slideshows anyway, and they
// will be way easier to test from a virtual provider.
// ** StandardAsset? **
// ALSO: some concept of sort order for them, maybe?

package provider

import (
	// Standard Library:
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	// Third-party packages:
	"gopkg.in/yaml.v2"
)

// YAML-parsing helper types.
// TODO: build a proper marshaler instead of this hack.
type stdPageFromYAML struct {
	Title   string
	Tags    []string
	Created time.Time
	Updated time.Time
	Meta    map[string]interface{}
	Html    string
}
type stdContentFromYAML struct {
	Type    string
	Content string
}
type stdFromYAML struct {
	Pages     map[string]*stdPageFromYAML
	Content   map[string]*stdContentFromYAML
	Templates map[string]string
}

// StandardPather is a minimal Pather.
type StandardPather struct {
	path string
}

// Path returns the path of the StandardPather.
func (p *StandardPather) Path() string { return p.path }

// NewStandardPather returns a StandardPather with path p.
func NewStandardPather(p string) *StandardPather {
	return &StandardPather{p}
}

// StandardStub is a minimal immutable non-Page Stub.
type StandardStub struct {
	path string
}

// Path returns the request path of the stub.
func (s *StandardStub) Path() string { return s.path }

// TypeString returns the stringified type of the stub: "StandardStub"
func (s *StandardStub) TypeString() string { return "StandardStub" }

// IsPageStub returns false for the StandardStub; use the StandardPageStub for
// a simple page-based stub.
func (s *StandardStub) IsPageStub() bool { return false }

// NewStandardStub returns a pointer to a StandardStub with the given path.
func NewStandardStub(rpath string) *StandardStub { return &StandardStub{rpath} }

// StandardPageStub is a stub based on a StandardPage.
type StandardPageStub struct {
	StandardPage
}

// TypeString returns the stringified type of the stub: "StandardPageStub"
func (s *StandardPageStub) TypeString() string { return "StandardPageStub" }

// IsPageStub returns true for the StandardPageStub.
func (s *StandardPageStub) IsPageStub() bool { return true }

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

// ModTime returns the modification time of the content. If it is the zero
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
func (p *StandardPage) Stub() Stub {
	return &StandardPageStub{*p}
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
func (p *StandardPage) HTML() template.HTML {
	return template.HTML(p.html)
}

// MetaString returns a string value from the page's Meta Map for the given
// key.  Lookup is case-sensitive.  The MappedString function is used.
func (p *StandardPage) MetaString(key string) string {
	return MappedString(p.meta, key)
}

// FlexMetaString returns a string value from the page's Meta Map for the
// given key, checking case variations. The FlexMappedString function is
// used.
func (p *StandardPage) FlexMetaString(key string) string {
	return FlexMappedString(p.meta, key)
}

// MetaStrings returns a slice of string values from the page's Meta Map,
// using the MappedStrings function.
func (p *StandardPage) MetaStrings(key string) []string {
	return MappedStrings(p.meta, key)
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

// StandardPageFromYAML returns a page from YAML data in src, or the first
// error encountered, with the assigned path.
func StandardPageFromYAML(path, src string) (*StandardPage, error) {

	target := stdPageFromYAML{}

	if err := yaml.Unmarshal([]byte(src), &target); err != nil {
		return nil, err
	}

	p := NewStandardPage(
		path,
		target.Title,
		target.Tags,
		target.Created,
		target.Updated,
		target.Meta,
		target.Html,
	)
	return p, nil

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
//      "html": "<h1>boo!</h1>",
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
		case "html":
			if val, ok := v.(string); ok {
				p.html = val
			} else {
				return nil, wrongTypeError(k, v, "string")
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

// StandardPathError is a minimal PathError implementation and may be used
// (usually via NewStandardPathError) for error responses from fetch methods
// such as Get.
type StandardPathError struct {
	path          string
	code          int
	message       string
	publicDetail  string
	privateDetail string
}

// Path returns the path of the StandardPathError in the Provider.
func (p *StandardPathError) Path() string { return p.path }

// Code returns the HTTP error code of the StandardPathError.
func (p *StandardPathError) Code() int { return p.code }

// Message returns the HTTP status message of the StandardPathError.
func (p *StandardPathError) Message() string { return p.message }

// PublicDetail returns the public-facing detail message of the
// StandardPathError.
func (p *StandardPathError) PublicDetail() string { return p.publicDetail }

// PrivateDetail returns the private, loggable detail message of the
// StandardPathError.
func (p *StandardPathError) PrivateDetail() string { return p.privateDetail }

// Error returns a stringified version of the StandardPathError suitable for
// non-http error handling.  Note that this does NOT include the PrivateDetail
// for the error.
func (p *StandardPathError) Error() string {
	return fmt.Sprintf("%d %s: %s (%s)", p.code, p.message, p.publicDetail, p.path)
}

// NewStandardPathError returns a PathError with the given path and code,
// and the positional messages:
//   Message (defaults to the standard message for the code)
//   PublicDetail
//   PrivateDetail
// It is normal to call NewStandardPathError with no msgStrings.  If more than
// one string is sent as a PrivateDetail then they are assumed to be sprintf-
// style (fmt,args) arguments.
func NewStandardPathError(path string, code int, msgStrings ...string) *StandardPathError {

	msg := ""
	public := ""
	private := ""
	if len(msgStrings) > 0 {
		msg = msgStrings[0]
	}
	if len(msgStrings) > 1 {
		public = msgStrings[1]
	}
	if len(msgStrings) > 2 {
		private = msgStrings[2]
	}
	if len(msgStrings) > 3 {
		f := msgStrings[2]
		sargs := msgStrings[3:]
		fargs := make([]interface{}, len(sargs))
		for idx, s := range sargs {
			fargs[idx] = interface{}(s)
		}
		private = fmt.Sprintf(f, fargs...)
	}

	if msg == "" {
		msg = http.StatusText(code)
	}
	return &StandardPathError{
		path:          path,
		code:          code,
		message:       msg,
		publicDetail:  public,
		privateDetail: private,
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
//
// The items are always returned in lexical (alpha) order by Path.
type StandardProvider struct {
	Meta     map[string]interface{}
	items    map[string]Pather
	paths    PathStrings
	modtimes map[string]time.Time
	template *template.Template
	created  time.Time
	updated  time.Time
	mutex    sync.RWMutex
}

// NewStandardProvider returns an empty StandardProvider to be populated
// via Add.  Directly creating a StandardProvider may lead to runtime errors
// in its methods; use this function instead.
func NewStandardProvider() *StandardProvider {
	return &StandardProvider{
		items:    map[string]Pather{},
		paths:    PathStrings{},
		modtimes: map[string]time.Time{},
		created:  time.Now(),
		updated:  time.Now(),
		mutex:    sync.RWMutex{},
	}
}

// StandardProviderFromYAML returns a StandardProvider with Pages, Content and
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
func StandardProviderFromYAML(src string) (*StandardProvider, error) {

	target := stdFromYAML{}

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
	rpath := p.Path()

	sp.mutex.Lock()
	sp.items[rpath] = p
	sp.modtimes[rpath] = modtime
	sp.updated = time.Now()
	sp.paths = sp.paths.Add(rpath)
	sp.mutex.Unlock()
}

// Get returns an item from the Provider if available; an error if not.
func (sp *StandardProvider) Get(rpath string) (Pather, PathError) {

	sp.mutex.RLock()
	item := sp.items[rpath]
	sp.mutex.RUnlock()

	if item == nil {
		return nil, NewStandardPathError(rpath, http.StatusNotFound)
	}
	return item, nil

}

// GetSince returns an item from the Provider if available and more recent
// than the since time; or an error if not. The caller should in particular
// watch for the ErrNotModified error, which usually indicates a 304 Not
// Modified response should be sent to the client.
func (sp *StandardProvider) GetSince(rpath string, since time.Time) (Pather, PathError) {

	sp.mutex.RLock()
	item := sp.items[rpath]
	modtime := sp.modtimes[rpath]
	sp.mutex.RUnlock()

	if item == nil {
		return nil, NewStandardPathError(rpath, http.StatusNotFound)
	}
	if modtime.After(since) {
		return item, nil
	}
	return nil, NewStandardPathError(rpath, http.StatusNotModified)
}

// GetStub returns a Stub item if available and implements the Stubber
// interface; an error if not.  An item that is not a Stubber is treated
// as not-found.
//
// The expected use-case for this is providing a preview of a large object,
// e.g. an image thumbnail with metadata, though a specialized Stub
// implementation.
func (sp *StandardProvider) GetStub(rpath string) (Stub, PathError) {
	sp.mutex.RLock()
	item, err := sp.Get(rpath)
	sp.mutex.RUnlock()

	if err != nil {
		return nil, err
	}
	if s, ok := item.(Stubber); ok {
		return s.Stub(), nil
	}
	return nil, NewStandardPathError(rpath, http.StatusNotFound)

}

// PathsUnder returns a slice of all item paths "under" the given prefix.
func (sp *StandardProvider) PathsUnder(prefix string) []string {

	sp.mutex.RLock()

	paths := []string{}

	for _, path := range sp.paths {
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
		}
	}
	sp.mutex.RUnlock()

	return paths

}

// GetPages returns a slice of all Pages with paths "under" the given prefix,
// following the same logic as GetAll.
func (sp *StandardProvider) GetPages(prefix string) []Page {

	paths := sp.PathsUnder(prefix)
	pages := []Page{}
	sp.mutex.RLock()
	for _, path := range paths {
		if item, _ := sp.Get(path); item != nil {
			if p, ok := item.(Page); ok {
				pages = append(pages, p)
			}
		}
	}
	sp.mutex.RUnlock()
	return pages
}

// GetPageStubs returns a slice of all PageStubs "under" the given prefix,
// following the logic described in GetStubs.
//
// This is the standard way of retrieving a list of Pages below the current
// directory.
func (sp *StandardProvider) GetPageStubs(prefix string) []PageStub {
	paths := sp.PathsUnder(prefix)
	pagestubs := []PageStub{}
	sp.mutex.RLock()
	for _, path := range paths {
		if s, _ := sp.GetStub(path); s != nil {
			if p, ok := s.(PageStub); ok {
				pagestubs = append(pagestubs, p)
			}
		}
	}
	sp.mutex.RUnlock()
	return pagestubs

}

// GetStubs returns a slice of Stub items for everything "under" the given
// prefix, i.e. everything with a path beginning with that prefix, that
// implement the Stubber interface.
//
// It is usually wise to terminate the prefix with a slash, but this is up to
// the template author.  If the prefix is the empty string then all available
// Stubs will be returned.
//
// If nothing exists under the prefix, and empty slice is returned.  Thus
// any "not found" must be implemented by the caller.
//
// Items are returned in string-sorted order by path.  Any item that does not
// implement the Stubber interface is ignored.
func (sp *StandardProvider) GetStubs(prefix string) []Stub {

	paths := sp.PathsUnder(prefix)
	stubs := []Stub{}
	sp.mutex.RLock()
	for _, path := range paths {
		if s, _ := sp.GetStub(path); s != nil {
			stubs = append(stubs, s)
		}
	}
	sp.mutex.RUnlock()

	return stubs

}

// GetAll returns a slice of Pathers representing all items "under" the
// given prefix, i.e. everything with a path having the given prefix. It is
// usually wise to terminate the prefix with a slash, but this is up to the
// template author. If the prefix is the empty string then all available
// Stubs will be returned.
//
// Items that implement the Stubber interface are returned as stubs; all
// others are returned as-is.
//
// If nothing exists under the prefix, and empty slice is returned. Thus
// any "not found" must be implemented by the caller.
//
// This means the caller must care what kind of data is being returned, as
// the only generally valid assuption is that every item is a Pather.
//
// TODO: iterator!
func (sp *StandardProvider) GetAll(prefix string) []Pather {

	paths := sp.PathsUnder(prefix)
	pathers := make([]Pather, len(paths))
	sp.mutex.RLock()
	for idx, path := range paths {
		if s, _ := sp.GetStub(path); s != nil {
			pathers[idx] = s
		} else {
			p, _ := sp.Get(path)
			pathers[idx] = p
		}
	}
	sp.mutex.RUnlock()
	return pathers

}

// Template returns the top-level template collection for the
// StandardProvider. It may be nil.
func (sp *StandardProvider) Template() *template.Template {
	return sp.template
}

// TemplateFor returns the template that should be used for the given
// Pather.  If the Pather is a Page then PageTemplate is used; otherwise
// PathTemplate.
func (sp *StandardProvider) TemplateFor(p Pather) *template.Template {
	if page, ok := p.(Page); ok {
		return PageTemplate(sp.template, page)
	} else {
		return PathTemplate(sp.template, p.Path())
	}
}

// TemplateForPath returns the template for the given path, via PathTemplate.
func (sp *StandardProvider) TemplateForPath(p string) *template.Template {
	return PathTemplate(sp.template, p)
}

// SetTemplate sets the internal template returned by Template.
func (sp *StandardProvider) SetTemplate(tmpl *template.Template) {
	sp.mutex.Lock()
	sp.template = tmpl
	sp.updated = time.Now()
	sp.mutex.Unlock()
}

// Paths returns an ordered list of the paths of all items known to the
// StandardProvider.  It is often the case that all paths without extensions
// are Pages.
func (sp *StandardProvider) Paths() []string {
	return []string(sp.paths)
}
