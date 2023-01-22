// loading.go - loading functions for the Kisipar site.
// ----------

package site

import (
	// Standard library:
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	// Third-party packages:
	"github.com/olebedev/config"

	// Kisipar packages:
	"github.com/biztos/kisipar/funcmap"
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
)

// At various points we need to unlist pages according to the site config.
// TODO: test for cross-platform, might need to use filepath here.
func (s *Site) unlistByPath(pages ...*page.Page) {
	if len(s.UnlistedPaths) == 0 {
		return
	}
	for _, p := range pages {
		path := strings.TrimPrefix(p.Path, s.PagePath)
		for _, prefix := range s.UnlistedPaths {
			if strings.HasPrefix(path, prefix) {
				p.Unlisted = true
				break
			}
		}
	}
}

// Load initializes a Site at the given directory path and loads its pages and
// templates.
func Load(path string) (*Site, error) {

	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	site, err := New(path)
	if err != nil {
		return nil, err
	}
	if err := site.LoadTemplates(); err != nil {
		return nil, err
	}
	if err := site.LoadPages(); err != nil {
		return nil, err
	}
	// TODO: static checks?
	return site, nil
}

// LoadPages loads the pages at the Site's PagePath into the Pageset.  If the
// Site already has a Pageset, it will be replaced.  Any file under the
// PagePath whose extension matches one of the PageExtensions will be
// loaded.
func (s *Site) LoadPages() error {

	wantExt := map[string]bool{}
	for _, e := range s.PageExtensions {
		wantExt[e] = true
	}
	s.Pageset, _ = pageset.New([]*page.Page{})
	if s.PagePath != "" {
		visit := func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() && wantExt[filepath.Ext(path)] {
				p, err := page.Load(path)
				if err != nil {
					return err
				}
				s.unlistByPath(p)
				s.Pageset.AddPage(p)
			}
			return nil
		}
		err := filepath.Walk(s.PagePath, visit)
		if err != nil {
			return err
		}
	}

	return nil

}

// LoadTemplates loads the templates under the Site's TemplatePath, putting
// them all into the Site's Template property.  The template names are the
// filepaths, lowercased and stripped of both the TemplatePath prefix and
// file extension.  Files may have any extension, but the cleaned path must
// be unique or an error will be returned.
func (s *Site) LoadTemplates() error {

	// The top (root) template has no name, and holds the default template.
	fm := funcmap.New()
	tmpl, err := template.New("").Funcs(fm).Parse(DEFAULT_TEMPLATE)
	if err != nil {
		return errors.New("DEFAULT_TEMPLATE: " + err.Error())
	}

	// Every file in the directory is a template, regardless of extension.
	if s.TemplatePath != "" {

		// filepath.Walk panics if the directory doesn't exist, so we check
		// that first.
		info, err := os.Stat(s.TemplatePath)
		if os.IsNotExist(err) {
			return fmt.Errorf("TemplatePath not found: %s", s.TemplatePath)
		}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("Not a directory: %s", s.TemplatePath)
		}

		// Now is the time on Sprockets when we walk:
		have := map[string]bool{}
		visit := func(path string, f os.FileInfo, err error) error {

			if !f.IsDir() {
				// Add the template under its cleaned path name.
				// TODO: make this platform-neutral, with slashes in the name
				// in all cases.
				name := strings.TrimPrefix(
					strings.ToLower(
						strings.TrimPrefix(
							strings.TrimSuffix(path, filepath.Ext(path)),
							s.TemplatePath)), "/")
				if have[name] {
					return fmt.Errorf("Duplicate template for %s: %s",
						name, path)
				}
				have[name] = true
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if _, err := tmpl.New(name).Parse(string(b)); err != nil {
					return err
				}
			}
			return nil
		}
		err = filepath.Walk(s.TemplatePath, visit)
		if err != nil {
			return err
		}

	}

	s.Template = tmpl
	return nil
}

// Load initializes a virtual site containing the provided pages, with cfg
// as its Config.  A nil Config is acceptable, as is an empty array of pages
// and a nil template.
func LoadVirtual(cfg *config.Config, pages []*page.Page,
	tmpl *template.Template) (*Site, error) {

	// Set an empty config if we have none.
	if cfg == nil {
		cfg, _ = config.ParseYaml("")
	}
	// En empty Path ensures we will not go looking for files on disk later.
	site := &Site{Path: "", Config: cfg}
	err := site.setup()
	if err != nil {
		return nil, err
	}

	// Set the default template if we have nothing yet.
	if tmpl == nil {
		fm := funcmap.New()
		tmpl, err = template.New("").Funcs(fm).Parse(DEFAULT_TEMPLATE)
		if err != nil {
			return nil, fmt.Errorf("DEFAULT_TEMPLATE: %s" + err.Error())
		}
	}
	site.Template = tmpl

	// Ingest the pages, if any:
	site.unlistByPath(pages...)
	ps, err := pageset.New(pages)
	if err != nil {
		return nil, err
	}
	site.Pageset = ps

	return site, nil

}

// LoadVirtualYaml initializes a virtual site from a YAML input file that
// contains the configuration and, optionally, maps of the Pages and
// Templates. All Pages will have the same ModTime.
func LoadVirtualYaml(yaml string) (*Site, error) {
	cfg, err := config.ParseYaml(yaml)
	if err != nil {
		return nil, err
	}

	// Using Now as the ModTime of each Page causes trouble, because there are
	// small differences in load time and the pages are a hash in the YAML,
	// so no guaranteed order.
	mtime := time.Now().UTC()

	// TODO: fix up path handling so we can feed in "foo/bar.md" etc, and
	// have it Do the Right Thing on e.g. Windows.  Just in principle.
	pages := []*page.Page{}
	if pp, _ := cfg.Map("Pages"); pp != nil {
		for k, v := range pp {
			if s, ok := v.(string); ok {
				p, err := page.LoadVirtualString(k, s)
				if err != nil {
					return nil, fmt.Errorf("Page %s: %s", k, err.Error())
				}
				p.ModTime = mtime
				pages = append(pages, p)
			} else {
				return nil,
					fmt.Errorf("Non-string value in Pages: %s", k)
			}
		}
	}

	fm := funcmap.New()
	tmpl, err := template.New("").Funcs(fm).Parse(DEFAULT_TEMPLATE)
	if err != nil {
		return nil, fmt.Errorf("DEFAULT_TEMPLATE: %s" + err.Error())
	}
	if tt, _ := cfg.Map("Templates"); tt != nil {
		for k, v := range tt {
			if s, ok := v.(string); ok {
				if _, err := tmpl.New(k).Parse(s); err != nil {
					return nil, fmt.Errorf("Template %s: %s", k, err.Error())
				}
			} else {
				return nil,
					fmt.Errorf("Non-string value in Templates: %s", k)
			}
		}
	}

	return LoadVirtual(cfg, pages, tmpl)

}
