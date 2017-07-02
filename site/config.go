// config.go - configuration for kisipar sites
// ---------

package site

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	// Third-party packages:
	"gopkg.in/yaml.v2"
)

// Config defines the configuration of a single Kisipar site.
type Config struct {
	Port           int                    // Port on which to listen.
	Name           string                 // Public-facing name of the site.
	Owner          string                 // Public-facing name of the owner
	Provider       string                 // Name of Provider to use.
	ProviderConfig map[string]interface{} // Configuration for the Provider.
	Dir            string                 // Root directory e.g. for:
	StaticDir      string                 // Directory for static files.
	ListStatic     bool                   // Serve directory listings?
	FastTemplates  bool                   // Faster (less safe) templates.
}

// LoadConfig loads a single Config from a YAML file.  The first error
// encountered is returned.  Note that the keys in the YAML file should be
// in lowercase.  If no Dir is specified then it is set to the absolute path
// of the directory of the file.
func LoadConfig(file string) (*Config, error) {

	if file == "" {
		return nil, errors.New("No config file specified")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	if cfg.Dir == "" {
		// NOTE: under what condition would Abs return an error?
		// If we ever figure that out we can test for it.
		cfg.Dir, _ = filepath.Abs(filepath.Dir(file))
	}

	if cfg.Provider == "" {
		cfg.Provider = "filesystem"
	}

	return cfg, nil
}
