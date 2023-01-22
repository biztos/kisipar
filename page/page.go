// page/page.go - Kisipar page core.
// ------------

// Package page defines a single page in a Kisipar web site.
package page

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/biztos/kisipar/utli"
)

// The StringArrayer interface enables custom types, usually named array
// types, to generate string arrays via MetaStringArray.
type StringArrayer interface {
	StringArray() []string
}

// A Page represents a single page, to be rendered in a template.  Pages
// should be created with New, Load or LoadVirtual.
type Page struct {

	// The data source:
	Path     string    // path as loaded; arbitrary unique for Virtual pages.
	Virtual  bool      // does the Page exist on disk?
	Unlisted bool      // should the Page be included in site listings?
	IsIndex  bool      // is the page an Index?
	ModTime  time.Time // file mod time; instantiation time for Virtual.
	Source   []byte    // raw, unparsed source data.

	// The standard rendered HTML content:
	Content template.HTML

	// The data extracted from the meta block:
	Meta map[string]interface{}

	// Some operations need to be atomic since we are normally operating
	// inside a web server.  Ergo:
	// TODO: really?
	mutex *sync.Mutex
}

// Load loads a page and parses it using the Parser associated with its
// extension in ExtParsers.
func Load(path string) (*Page, error) {

	if path == "" {
		return nil, errors.New("page.Load requires a source path.")
	}

	page, err := New(path)
	if err != nil {
		return nil, err
	}
	if err := page.Load(); err != nil {
		return nil, err
	}
	if err := page.Parse(); err != nil {
		return nil, err
	}

	return page, nil

}

// LoadAny loads the first page it finds for the source path with an extension
// listed in ExtParsers.  Extensions are matched exactly, meaning that
// ExtParsers must have entries of every supported extension case.
func LoadAny(spath string) (*Page, error) {
	if spath == "" {
		return nil, errors.New("page.LoadAny requires a source path.")
	}
	for _, ep := range ExtParsers {
		page, err := Load(spath + ep.Ext)
		if err == nil {
			return page, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, os.ErrNotExist
}

// LoadVirtual returns a page with a virtual path, i.e. not necessarily
// corresponding to any file on disk.  The provided path should indicate
// the source type, e.g. "virtual/document.md" for Markdown.
func LoadVirtual(path string, input []byte) (*Page, error) {

	if path == "" {
		return nil, errors.New("page.LoadVirtual requires a source path.")
	}

	page, err := New(path)
	if err != nil {
		return nil, err
	}

	page.Virtual = true
	page.Source = input
	page.ModTime = time.Now().UTC()

	if err := page.Parse(); err != nil {
		return nil, err
	}

	return page, nil
}

// LoadVirtualString is shorthand for LoadVirtual(path,[]byte(string).
func LoadVirtualString(path, input string) (*Page, error) {
	return LoadVirtual(path, []byte(input))
}

// New returns a Page with source loaded from the provided path, but not yet
// parsed.  The path must have an extension.
func New(path string) (*Page, error) {

	if path == "" {
		return nil, errors.New("page.New requires a source path.")
	}

	ext := filepath.Ext(path)
	if ext == "" {
		return nil, fmt.Errorf("No file extension in source path: %s", path)
	}

	// Is it an index?
	idxBase := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ext))
	isIndex := idxBase == "index"

	// Don't forget to Load if non-virtual!
	return &Page{
		Path:    path,
		IsIndex: isIndex,
		mutex:   &sync.Mutex{},
	}, nil

}

// Load loads the source data from the Page's path, but does not parse it.
func (p *Page) Load() error {

	path := p.Path

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	p.Source = b
	p.ModTime = info.ModTime().UTC()

	return nil
}

// Parse parses the Page's Source data into Meta and Content using the first
// Parser in ExtParsers to match the Page's Path extension, or the
// DefaultParser if none matches. Extension matching is case-insensitive.
//
// If the Meta defines a boolean "Unlisted" (or "unlisted" or "UNLISTED")
// with a value of true, the Page's Unlisted property is set to true.
func (p *Page) Parse() error {

	// First let the parser do the heavy lifting:
	res, err := p.getParser().Parse(p.Source)
	if err != nil {
		return err
	}
	p.Meta = res.Meta()
	p.Content = template.HTML(res.Content())

	if p.MetaBool("Unlisted") {
		p.Unlisted = true
	}

	return nil

}

func (p *Page) getParser() Parser {
	ext := strings.ToLower(filepath.Ext(p.Path))
	for _, ep := range ExtParsers {
		if ext == strings.ToLower(ep.Ext) {
			return ep.Parser
		}
	}
	return DefaultParser

}

// String returns a hopefully-useful stringification of the page, for logging
// and debugging purposes.
func (p *Page) String() string {

	times := []string{}
	if t := p.Created(); t != nil {
		times = append(times, "Created: "+t.String())
	}
	if t := p.Updated(); t != nil {
		times = append(times, "Updated: "+t.String())
	}
	times = append(times, "ModTime: "+p.ModTime.String())
	return fmt.Sprintf("%s: %s (%s)",
		p.Path, p.Title(), strings.Join(times, "; "))

}

// Refresh reloads the page if it is not Virtual and the modtime of the source
// file is different than the current ModTime.  In case of file errors,
// including not-found, the error will be returned without the page
// being modified, and the caller must handle the error.
func (p *Page) Refresh() error {
	if p.Virtual {
		return nil
	}

	// We don't want multiple threads trying to do the load simultaneously,
	// because that's asking for trouble even if the property resets are the
	// only *actual* point of contention.
	p.mutex.Lock()
	info, err := os.Stat(p.Path)
	if err != nil {
		p.mutex.Unlock()
		return err
	}
	if info.ModTime().UTC() == p.ModTime {
		p.mutex.Unlock()
		return nil
	}

	fresh, err := Load(p.Path)
	if err != nil {
		p.mutex.Unlock()
		return err
	}

	// We have new page content.
	p.Source = fresh.Source
	p.ModTime = fresh.ModTime
	p.Meta = fresh.Meta
	p.Content = fresh.Content

	p.mutex.Unlock()
	return nil

}

// Time returns the newer of the page's Created and Updated meta times;
// if neither is set or neither is parseable, returns the ModTime.
func (p *Page) Time() *time.Time {
	created := p.Created()
	updated := p.Updated()
	if created == nil && updated == nil {
		return &p.ModTime
	}
	if created == nil || (updated != nil && updated.After(*created)) {
		return updated
	}

	return created
}

// Created returns the page's creation time as defined in the meta.
// Created is simply shorthand for MetaTime("Created").
func (p *Page) Created() *time.Time {
	return p.MetaTime("Created")
}

// Updated returns the page's update time as defined in the meta.
// Updated is simply shorthand for MetaTime("Updated").
func (p *Page) Updated() *time.Time {
	return p.MetaTime("Updated")
}

// Title returns the Title string from the Page's Meta block, or the file
// name (without extension) of the Path if no Title is available in the Meta.
func (p *Page) Title() string {
	title := p.MetaString("Title")
	if title == "" && p.Path != "" {
		ext := filepath.Ext(p.Path)
		base := filepath.Base(p.Path)
		title = strings.TrimSuffix(base, ext)
	}
	return title
}

// Author returns the Author string from the Page's Meta block.
// Author is simply shorthand for MetaString("Author").
func (p *Page) Author() string {
	return p.MetaString("Author")
}

// Description returns the Description string from the Page's Meta block.
// Description is simply shorthand for MetaString("Description").
func (p *Page) Description() string {
	return p.MetaString("Description")
}

// Summary returns the Summary string from the Page's Meta block.
// Summary is simply shorthand for MetaString("Summary").
func (p *Page) Summary() string {
	return p.MetaString("Summary")
}

// Keywords returns the Keywords string from the Page's Meta block.
// (Keywords are generally used in the HTML meta elements, which is why this
// is a string and not a []string. For an array, use Tags.)
// Keywords is simply shorthand for MetaString("Keywords").
func (p *Page) Keywords() string {
	return p.MetaString("Keywords")
}

// Tags returns a list of tags from the Page's Meta block, as an array of
// strings.  This is normally defined in the meta as either an array of
// strings directly, or a comma-delimited list in string form such as
// "foo, bar, baz" -- this function supports both forms within the Meta
// map but the implmentation is up to the Parser.
// Tags is simply shorthand for MetaStringArray("Tags")
func (p *Page) Tags() []string {
	return p.MetaStringArray("Tags")
}

// MetaBool returns a boolean value from the Page's Meta block, with
// undefined and non-boolean values treated as false. If there is no match
// for the exact key, lookup is attempted on the lowercase and uppercase
// versions, in that order.
func (p *Page) MetaBool(key string) bool {

	val := p.metaValForKey(key)
	if val == nil {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	default:
		return false
	}

}

// MetaString returns a string value from the Meta, with non-string values
// stringified.  If there is no match for the exact key, lookup is attempted
// on the lowercase and uppercase versions, in that order.  If there is no
// value for the key in any form, the empty string is returned.  Thus
// MetaString should not be used to determine the presence of a key, only
// in cases where an empty string and nil/not-found have identical meaning.
func (p *Page) MetaString(key string) string {

	val := p.metaValForKey(key)
	if val == nil {
		return ""
	}

	return stringify(val)

}

// MetaTime returns a time value from the Meta, if the requested key holds
// a string value that can be parsed into a time.  Nil is returned in all
// other cases.
func (p *Page) MetaTime(key string) *time.Time {

	return utli.ParseTimeString(p.MetaString(key))

}

// MetaStringArray returns an array of string values from the Meta, or an
// empty array if there is value for the key. If the value is already a
// []string, it is returned as-is. If the value is a string, it is split
// with a comma delimiter and the resulting items are trimmed of leading and
// trailing spaces. If the value is an array of integers or floats then its
// contents are stringified. For any type that implements the StringArrayer
// interface, the result of its StringArray function is returned.  All other
// types return an empty array. Key lookup follows the logic of MetaString.
func (p *Page) MetaStringArray(key string) []string {

	val := p.metaValForKey(key)
	if val == nil {
		return []string{}
	}

	// NOTE: some of these cases are quite redundant in order to avoid the
	// overhead of reflection, since we (weirdly) lose the ability to deal
	// with a slice when we handle more than one slice type in the case.
	// TODO: benchmark this vs reflection, see which is faster, blog it!
	switch v := val.(type) {
	case []string:
		return v
	case string:
		res := strings.Split(v, ",")
		for i, s := range res {
			res[i] = strings.TrimSpace(s)
		}
		return res
	case []int:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []int32:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []int64:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []uint:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []uint32:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []uint64:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []float32:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []float64:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	case []interface{}:
		res := make([]string, len(v))
		for i, e := range v {
			res[i] = stringify(e)
		}
		return res
	default:
		if s, ok := val.(StringArrayer); ok {
			return s.StringArray()
		}
		return []string{}
	}
}

// NOTE: this might not be much better than using fmt.Sprintf...
func stringify(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// TODO: useful formatting of times
		if i, ok := v.(fmt.Stringer); ok {
			return i.String()
		}

		// Fall back to whatever fmt can do for us; this seems like the
		// most useful thing for random incoming values, getting us strings
		// like "&{Whatever:thing}" for a pointer to a struct with a Whatever
		// string set to "thing" (but maybe we're overthinking this).
		return fmt.Sprintf("%+v", v)
	}
}

func (p *Page) metaValForKey(key string) interface{} {

	v := p.Meta[key]
	if v == nil {
		v = p.Meta[strings.ToLower(key)]
		if v == nil {
			v = p.Meta[strings.ToUpper(key)]
		}
	}
	return v
}
