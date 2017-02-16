// filesystemprovider.go -- File System Based Provider for Kisipar.
// ---------------------

package kisipar

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	// Third-party packages:
	//	"gopkg.in/yaml.v2"

	// Own stuff:
	"github.com/biztos/frostedmd"
)

// FileSystemProviderConfig defines the configuration options for a
// FileSystemProvider.
type FileSystemProviderConfig struct {
	ContentDir      string   // All content, parseable and static.
	TemplateDir     string   // Templates, if any.
	ExcludePrefixes []string // Ignore anything with these prefixes.
	ExcludeSuffixes []string // Ignore anything with these suffixes.
	Strict          bool     // Fail on Page metadata parsing errors.
	AutoRefresh     bool     // Watch filesystem and refresh
}

// FileSystemProvider is a Provider loaded from the local file system. Its
// items are stored in memory as Pages or Files, and may optionally be auto-
// refreshed when the source directory changes.
//
// It is the recommended Provider to use for developing templates, and is
// also useful for set-and-forget sites such as placeholders or smaller
// archives.
type FileSystemProvider struct {
	*StandardProvider

	// Note to self: do this more! By not having a simple struct be a pointer
	// we avoid all the crap around catching nil values, while losing...
	// nothing.
	config FileSystemProviderConfig
}

// LoadTemplates loads the templates from the configred TemplateDir and sets
// them in the FileSystemProvider.  The first error encountered is returned.
// If the config is not set, or the TemplateDir is the empty string, then
// no action is taken.
func (fsp *FileSystemProvider) LoadTemplates() error {

	config := fsp.config
	if config.TemplateDir == "" {
		return nil // NOOP
	}

	dirInfo, err := os.Stat(config.TemplateDir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("Not a directory: %s", config.TemplateDir)
	}

	tmpl, _ := template.New("").Funcs(FuncMap()).Parse("")
	walker := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// It might be a link to a dir, or something missing...
		realInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if realInfo.IsDir() {
			return nil
		}

		// Now the fun begins!
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error reading %s: %v", path, err)
		}

		path = filepath.ToSlash(strings.TrimPrefix(path, config.TemplateDir))
		path = strings.TrimPrefix(path, string(filepath.Separator))

		if _, err := tmpl.New(path).Parse(string(b)); err != nil {
			return fmt.Errorf("Template %s failed: %v", path, err)
		}

		return nil
	}

	err = filepath.Walk(config.TemplateDir, walker)
	if err != nil {
		return fmt.Errorf("Error walking %s: %v", config.TemplateDir, err)
	}

	fsp.SetTemplate(tmpl)

	return nil

}

// LoadContent loads all items in the FileSystemProvider's configured
// ContentDir, adding StandardPage or StandardFile items as appropriate.
// The first error encountered is returned.
//
// If the ContentDir is the empty string then no action is taken.
//
// By default all listed items are included except those starting with a dot;
// this can be changed by setting ExcludePrefixes and/or ExcludeSuffixes in
// the config.  Note that these exclusions apply to both files and
// directories.
//
// Markdown and YAML files are added as StandardPages, while all other
// files are added as StandardFiles.  The file extensions recognized are
// .yaml, .yml, and .md.
//
// Page files have their extensions stripped from their request paths, so
// "/foo/bar/baz.md" has a Path of "/foo/bar/baz".
//
// Markdown conversion uses Frosted Markdown, allowing for a data block at
// the top of the file:
//
// https://github.com/biztos/frostedmd
//
// Anything not loaded here is removed, in a blocking operation, after the
// fact.
//
// TODO: mutex locking strategy here!  Just duplicate Add or what?
// (don't want to deadlock)
func (fsp *FileSystemProvider) LoadContent() error {

	config := fsp.config
	if config.ContentDir == "" {
		return nil // NOOP
	}

	dirInfo, err := os.Stat(config.ContentDir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("Not a directory: %s", config.ContentDir)
	}

	// We can't use the standard Add function because we need to lock
	// everything long enough to replace all the items (i.e. remove things).
	// Instead we have to build our items in place, then swap them.  Easy.
	mdparser := frostedmd.New()
	items := map[string]Pather{}
	walker := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// Watch for symlinks... though there must be a better way, no?
		realInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if realInfo.IsDir() {
			return nil
		}

		// We will keep it now, whatever it is.  For which we need its
		// request-style path.
		rpath := filepath.ToSlash(strings.TrimPrefix(path, config.ContentDir))
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" {
			// Markdown page.
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			res, err := mdparser.Parse(b)
			if err != nil && config.Strict {
				return err
			}

			// Get some things from the meta.
			title := FlexMappedString(res.Meta, "title")
			p, _ := StandardPageFromData(map[string]interface{}{
				"path":  rpath,
				"title": title,
				// "tags":  []string{"foo", "bar"},
				//  "created": time.Now(),
				//  "updated": time.Now(),
				"meta": res.Meta,
			})
			items[rpath] = p

		} else if ext == ".yml" || ext == ".yaml" {
			// YAML page.
		} else {
			// File to be served as-is.
			items[rpath] = NewStandardFile(rpath, path)

		}

		return nil
	}

	err = filepath.Walk(config.ContentDir, walker)
	if err != nil {
		return fmt.Errorf("Error walking %s: %v", config.ContentDir, err)
	}

	return nil

}

// NewFileSystemProvider returns a FileSystemProvider with the provided
// config.  No action is taken; use LoadFileSystemProvider in most cases.
func NewFileSystemProvider(config FileSystemProviderConfig) *FileSystemProvider {

	return &FileSystemProvider{NewStandardProvider(), config}

}

// LoadFileSystemProvider returns a FileSystemProvider intialized from the
// given config with templates and items loaded, via LoadTemplates and
// LoadItems; or the first error that occurs.
//
// Templates are read from the configured TemplateDir, and all available items
// are read from the ContentDir.  TemplateDir is optional, allowing for the
// use of built-in default templates.
//
// If the configured AutoRefresh is true, then both TemplateDir and ContentDir
// are watched for changes and the content updated when needed.
func LoadFileSystemProvider(config FileSystemProviderConfig) (*FileSystemProvider, error) {

	fsp := NewFileSystemProvider(config)

	// Load templates first, if defined.
	if err := fsp.LoadTemplates(); err != nil {
		return nil, err
	}

	// TODO: load shit up, set up autorefresh, etc.
	return fsp, nil

}

// String returns a log-friendly description of the Provider.
func (fsp *FileSystemProvider) String() string {

	return fmt.Sprintf("<FileSystemProvider with %d items at %s, updated %s>",
		len(fsp.items), fsp.config.ContentDir, fsp.updated)
}
