// config.go - configuration for kisipar sites
// ---------
// (Largely TODO...)

package kisipar

import (
	"errors"
	"io/ioutil"

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
}

// LoadConfig loads a single Config from a YAML file.  The first error
// encountered is returned.  Note that the keys in the YAML file should be
// in lowercase.
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
	return cfg, nil
}
