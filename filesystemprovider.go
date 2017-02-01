// filesystemprovider.go -- File System Based Provider for Kisipar.
// ---------------------

package kisipar

import (
	"fmt"
)

// FileSystemProviderConfig defines the configuration options for a
// FileSystemProvider.
type FileSystemProviderConfig struct {
	ContentDir      string   // All content, parseable and static.
	TemplateDir     string   // Templates, if any.
	ExcludePrefixes []string // Ignore anything with these prefixes.
	ExcludeSuffixes []string // Ignore anything with these suffixes.
	AutoRefresh     bool
}

// FileSystemProvider is a Provider loaded from the local file system.  It
// lives in memory but may optionally be auto-refreshed when the source
// directory changes. It is the recommended Provider to use for developing
// templates, and is also useful for set-and-forget sites such as
// placeholders or smaller archives.
type FileSystemProvider StandardProvider

// NewFileProvider returns a FileProvider intialized from the given config,
// or the first error that occurs in loading and parsing content and
// templates.
//
// By default all listed items are included except those starting with a dot;
// this can be changed by setting ExcludedPrefixes and/or ExcludedSuffixes in
// the config.
//
// Page files have their extensions stripped from their request paths, so
// "/foo/bar/baz.md" has a Path of "/foo/bar/baz".
//
// Index files ("index.html", "index.md") are mapped to their parent directory
// path if and only if no   Directories are not loaded as separate items; anything with
// the
//
// HTML, YAML and Markdown files are added as StandardPages, while all other
// files are added as StandardFiles.
func NewFileSystemProvider(config *FileSystemProviderConfig) (*FileSystemProvider, error) {
	fsp := &FileSystemProvider{
		Meta: map[string]interface{}{"config": config},
	}
	// TODO: load shit up, set up autorefresh, etc.
	return fsp, nil

}

// String returns a log-friendly description of the Provider.
func (fsp *FileSystemProvider) String() string {
	return fmt.Sprintf("<FileSystemProvider with %d items, updated %s>",
		len(fsp.items), fsp.updated)
}
