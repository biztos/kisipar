// filesystemprovider.go -- File System Based Provider for Kisipar.
// ---------------------

package kisipar

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// Own stuff:
	"github.com/biztos/frostedmd"
)

// FileSystemProviderConfig defines the configuration options for a
// FileSystemProvider.
type FileSystemProviderConfig struct {
	ContentDir      string         // All content, parseable and static.
	TemplateDir     string         // Templates, if any.
	TemplateTheme   string         // Theme for default templates, if any.
	Exclude         *regexp.Regexp // Exclude paths matching this regexp.
	Include         *regexp.Regexp // Include paths matching this regexp.
	AllowMetaErrors bool           // Don't fail on Page metadata errors.
	AutoRefresh     bool           // Watch filesystem and refresh
}

// NewFileSystemProviderConfig returns a FileSystemProviderConfig based on
// the provided data, accepting strings for Regexp values.  Wrong-typed
// items, as well as unknown keys, are treated as errors.
func NewFileSystemProviderConfig(d map[string]interface{}) (*FileSystemProviderConfig, error) {
	expected := map[string]bool{
		"ContentDir":      true,
		"TemplateDir":     true,
		"TemplateTheme":   true,
		"Exclude":         true,
		"Include":         true,
		"AllowMetaErrors": true,
		"AutoRefresh":     true,
	}
	for k, _ := range d {
		if !expected[k] {
			return nil,
				fmt.Errorf("Unexpected FileSystemProviderConfig key: %s", k)
		}
	}

	cdir, err := mapString(d, "ContentDir")
	if err != nil {
		return nil, err
	}
	tdir, err := mapString(d, "TemplateDir")
	if err != nil {
		return nil, err
	}
	theme, err := mapString(d, "TemplateTheme")
	if err != nil {
		return nil, err
	}
	exclude, err := mapRegexp(d, "Exclude")
	if err != nil {
		return nil, err
	}
	include, err := mapRegexp(d, "Include")
	if err != nil {
		return nil, err
	}
	metaerrs, err := mapBool(d, "AllowMetaErrors")
	if err != nil {
		return nil, err
	}
	refresh, err := mapBool(d, "AutoRefresh")
	if err != nil {
		return nil, err
	}

	cfg := &FileSystemProviderConfig{
		ContentDir:      cdir,
		TemplateDir:     tdir,
		TemplateTheme:   theme,
		Exclude:         exclude,
		Include:         include,
		AllowMetaErrors: metaerrs,
		AutoRefresh:     refresh,
	}
	return cfg, nil

}

// This is arguably useful elsewhere... or I'm over-engineering the shit...
// What I really want is to be able to marshal a map[string]interface
// to a generic struct, while magically handling things like regexp.
func mapString(d map[string]interface{}, k string) (string, error) {
	v := d[k]
	if v == nil {
		return "", nil
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("%s must be a string, not %T.", k, v)
}
func mapBool(d map[string]interface{}, k string) (bool, error) {
	v := d[k]
	if v == nil {
		return false, nil
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("%s must be a bool, not %T.", k, v)
}
func mapRegexp(d map[string]interface{}, k string) (*regexp.Regexp, error) {
	v := d[k]
	if v == nil {
		return nil, nil
	}
	if r, ok := v.(*regexp.Regexp); ok {
		return r, nil
	}
	if s, ok := v.(string); ok {
		r, err := regexp.Compile(s)
		if err != nil {
			return nil,
				fmt.Errorf("%s is not a valid regexp string: %s", k, err.Error())
		}
		return r, nil
	}

	return nil,
		fmt.Errorf("%s is neither a *regexp.Regexp nor a string, but a %T.", k, v)
}

// FileSystemProvider is a Provider loaded from the local file system. Its
// items are stored in memory as Pages or Files, and may optionally be auto-
// refreshed when the source directory changes.
//
// It is the recommended Provider to use for developing templates, and is
// also useful for simpler, mostly static web sites.
type FileSystemProvider struct {
	*StandardProvider

	// Note to self: do this more! By not having a simple struct be a pointer
	// we avoid all the crap around catching nil values, while losing...
	// nothing.
	// OR not?
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

	tmpl := template.Must(template.New("").Funcs(FuncMap()).Parse(""))
	walker := func(path string, info os.FileInfo, err error) error {

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

// LoadInternalTemplates loads internal templates, which are available in
// variations or themes.  The theme may be specified in the config's
// TemplateTheme property.
func (fsp *FileSystemProvider) LoadInternalTemplates() error {

	theme := fsp.config.TemplateTheme
	if theme == "" {
		theme = "default"
	}

	// We must have internal assets for this or we return an error.
	// The assets must match be under the theme path and end in .html,
	// as we aren't doing anything complicated (at least not intentionally)
	// with the built-in templates.
	prefix := "templates/" + theme + "/"
	paths := []string{}
	for _, name := range AssetNames() {
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".html") {
			paths = append(paths, name)
		}
	}
	if len(paths) == 0 {
		return fmt.Errorf("No templates available for theme %s.", theme)
	}

	// Parse the template assets.  Path should already be normalized to slash.
	// We panic instead of surfacing errors because the internal templates
	// should all be error-free (and this is exercised in the unit tests).
	tmpl := template.Must(template.New("").Funcs(FuncMap()).Parse(""))
	for _, path := range paths {

		b := MustAsset(path)
		path = strings.TrimPrefix(path, prefix)
		template.Must(tmpl.New(path).Parse(string(b)))

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
// By default all listed items are included except those starting with a dot.
// Additional exclusions can be defined by setting the Exclude and Include
// Regexp items in the config; these are applied to the relative path of each
// file on disk.  If a path matches both expressions it will be excluded.
//
// Markdown and YAML files are added as StandardPages, while all other
// files are added as StandardFiles.  The file extensions recognized are
// .yaml, .yml, and .md.  (Note that HTML can be used as-is in Markdown files,
// allowing verbatim content to be included in a StandardPage.)
//
// Page files have their extensions stripped from their request paths, so
// "/foo/bar/baz.md" has a Path of "/foo/bar/baz".
//
// Markdown conversion uses Frosted Markdown, allowing for a data block at
// the top of the file:
//
// https://github.com/biztos/frostedmd
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

	mdparser := frostedmd.New()
	walker := func(path string, info os.FileInfo, err error) error {

		// Check exclusions/inclusions first (we assume this is faster than
		// doing file system checks).
		relpath, err := filepath.Rel(config.ContentDir, path)
		if config.Exclude != nil && config.Exclude.MatchString(relpath) {
			return nil
		}
		if config.Include != nil && !config.Include.MatchString(relpath) {
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
			rpath = strings.TrimSuffix(rpath, ext)
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			res, err := mdparser.Parse(b)
			if err != nil && !config.AllowMetaErrors {
				return fmt.Errorf("%s: %s", path, err)
			}
			// Get some things from the meta.
			// TODO: additional meta stuff for tags etc?
			title := FlexMappedString(res.Meta, "title")
			p, _ := StandardPageFromData(map[string]interface{}{
				"path":  rpath,
				"title": title,
				// "tags":  []string{"foo", "bar"},
				//  "created": time.Now(),
				//  "updated": time.Now(),
				"meta": res.Meta,
				"html": string(res.Content),
			})
			fsp.Add(p)

		} else if ext == ".yml" || ext == ".yaml" {
			// YAML page.
			// TODO: parse it as a page, obviously.
			rpath = strings.TrimSuffix(rpath, ext)
			fsp.Add(NewStandardFile(rpath, path))
		} else {
			// File to be served as-is.
			fsp.Add(NewStandardFile(rpath, path))

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
func NewFileSystemProvider(cfg *FileSystemProviderConfig) *FileSystemProvider {

	return &FileSystemProvider{NewStandardProvider(), *cfg}

}

// LoadFileSystemProvider returns a FileSystemProvider intialized from the
// given config with templates and items loaded, via LoadTemplates and
// LoadContent; or the first error that occurs.
//
// Templates are read from the configured TemplateDir, and all available items
// are read from the ContentDir.  TemplateDir is optional, allowing for the
// use of built-in templates.
//
// If the configured AutoRefresh is true, then both TemplateDir and ContentDir
// are watched for changes and the content updated when needed.
func LoadFileSystemProvider(cfg *FileSystemProviderConfig) (*FileSystemProvider, error) {

	fsp := NewFileSystemProvider(cfg)

	// We are more likely to hit template errors than content errors so we
	// start with Templates.
	if cfg.TemplateDir == "" {
		if err := fsp.LoadInternalTemplates(); err != nil {
			return nil, err
		}
	} else {
		if err := fsp.LoadTemplates(); err != nil {
			return nil, err
		}
	}

	// Then load content.
	if err := fsp.LoadContent(); err != nil {
		return nil, err
	}

	// TODO: set up autorefresh, etc.
	return fsp, nil

}

// String returns a log-friendly description of the Provider.
func (fsp *FileSystemProvider) String() string {

	return fmt.Sprintf("<FileSystemProvider with %d items at %s, updated %s>",
		len(fsp.items), fsp.config.ContentDir, fsp.updated)
}
