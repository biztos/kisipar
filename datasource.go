// datasource.go
package kisipar

import (
	// Standard Library:
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	// Third-party packages:
	"github.com/olebedev/config"
)

// Page defines a single page renderable as an HTML fragment.  It is
// typically returned from a DataSource and made available in a template.
type Page interface {
	Id() string     // Unique ID of the Page.
	Title() string  // Title of the Page.
	Tags() []string // List of Tags applicable to the Page.
	HTML() string   // Rendered HTML fragment of the Page.
}

// Data defines an arbitrary set of bytes with a content-type as returned
// by the DataSource.
type Data interface {
	ContentType() string
	Bytes() []byte
}

// File defines a file on disk which can be further processed, e.g. by the
// http.ServeFile handler.
type File interface {
	Path() string
}

// DataSource is a provider of Pages, Files and Data, usually as an interface
// to external resoures such as a database, cloud environment, cache,
// filesystem or all of the above.
//
// When processing a file with a given
type DataSource interface {
	String() string                    // For logging / debugging.
	Has(rpath string) bool             // Check for existence of an item.
	Page(rpath string) (Page, error)   // Fetch a Page for a given path.
	Data(rpath string) (Data, error)   // Fetch a Page for a given path.
	File(rpath string) (string, error) // Return the path to the File o
	Handler() http.Handler             // For accepting any commands.
}

// StandardPage is a type of Page that can hold its data in memory and/or
// in a (Frosted) Markdown, YAML or HTML file.
type StandardPage struct {
	RequestPath string    // Path to the resource in a request: foo/bar
	FilePath    string    // Path to the file on disk: foo/bar.md
	ModTime     time.Time // Last read/init/etc. time.

	MetaData map[string]interface{}
	Content  string // HTML content, rendered if applicable.
	Virtual  bool   // true if NOT backed by the file system.
}

// Load loads the page from the file at FilePath.  Supported file types are
// identified by extension and parsed: .md (or .txt) for (Frosted) Markdown,
// .yaml (or .yml) for YAML.
//
// Any other file type is simply loaded as-is and treated as rendered
// HTML-safe content, which of course might not be the case. The DataSource
// is responsible for limiting which file types can become Pages.
// This allows for, say, SVG illustrations to be treated as Pages while
// retaining the .svg extension.
func (p *StandardPage) Load() error {

	// Get the data first, then decide what to do with it.
	b, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		return err
	}

	// We are somewhat liberal with the file extensions.
	ext := strings.ToLower(p.FilePath.Ext())
	if ext == ".md" || ext == ".txt" {
		return p.ParseMarkdown(b)
	} else if ext == ".yaml" || ext == ".yml" {
		return p.ParseYaml(b)
	} else {
		p.Content = string(b)
	}

}

// StandardPageFromFile returns a StandardPage for the request path rpath,
// with loaded from the given file path.  The file must be of a type supported
// by the Load function.
func StandardPageFromFile(rpath, file string) (*StandardPage, error) {

	p := &StandardPage{
		RequestPath: rpath,
		FilePath:    file,
		ModTime:     time.Now(),
	}
	if err := p.Load(); err != nil {
		return nil, err
	}
	return p, nil

}

// VirtualDataSource is a DataSource that exists entirely in memory.  It
// is primarily useful for testing, but might have other uses as well.
// (What might those be? Dynamically creating a site based on read-once data?
// Making a placeholder that exists entirely as a config file? Remotely
// updating a small site via an API?)
type VirtualDataSource struct {
	pages map[string]StandardPage
	data  map[string]Data
}

// NewVirtualDataSource returns a VirtualDataSource with its pages initialized
// from the provided yaml string as per the Initialize function.
func NewVirtualDataSource(yaml string) {
	return &FileDataSource{}
}

// FileDataSource is a DataSource using only the file system.  It is the
// recommended DataSource to use for developing templates, and is also
// useful for set-and-forget sites such as placeholders or smaller archives.
// It is NOT recommended for large sites or any case where performance is key.
type FileDataSource struct {
	root   string
	config *config.Config
}

// String returns a useful description of the DataSource,
func (ds *FileDataSource) String() string {
	return fmt.Sprintf("<FileDataSource in %s>", ds.root)
}

// Handler returns a not-found handler as there are no applicable commands.
func (ds *FileDataSource) Handler() http.Handler {
	return http.NotFoundHandler()
}

// NewFileDataSource returns a FileDataSource intialized from the given
// root path.
func NewFileDataSource(root string) *FileDataSource {
	return &FileDataSource{}
}
