// site.go -- A single Kisipar site.

package kisipar

import (
	"errors"
	"fmt"
	"net/http"
)

// Site represents a single Kisipar web site.
type Site struct {
	Config   *Config
	Provider Provider
	Server   *http.Server
	Mux      *http.ServeMux
	// TODO: logger?
}

// NewSite returns an initialized Site for the given Config, or the first
// error encountered.
func NewSite(cfg *Config) (*Site, error) {
	if cfg == nil {
		return nil, errors.New("Config required")
	}
	s := &Site{
		Config: cfg,
		Server: &http.Server{},
		Mux:    http.NewServeMux(),
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s, nil

}

// Init initializes the Site based on the Config, returning the first error
// encountered.
func (s *Site) Init() error {

	// Set up the Provider.
	switch s.Config.Provider {
	case "":
		return errors.New("Provider missing from Config.")
	case "filesystem":
		cfg, err := FileSystemProviderConfigFromData(s.Config.ProviderConfig)
		if err != nil {
			return fmt.Errorf("ProviderConfig error: %s", err)
		}
		fsp, err := LoadFileSystemProvider(cfg)
		if err != nil {
			return err
		}
		s.Provider = fsp

	default:
		return errors.New("Unsupported Provider: " + s.Config.Provider)
	}

	// Set up the Mux (multiplexer).
	// TODO
	return nil
}
