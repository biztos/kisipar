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
// error encountered.  The configured Provider will be initialized, and
// a Mux and Server set up to handle requests.  If any of these is already
// set up, use the corresponding Init methods on a directly constructed
// Site for granular control:
//
//    s := &Site{Config: myConfig, Provider: myProvider}
//    if err := s.InitMux(); err != nil {
//        panic(err)
//    }
//    // ...and so on.
func NewSite(cfg *Config) (*Site, error) {

	if cfg == nil {
		return nil, errors.New("Config must not be nil")
	}

	// Can't create a Site without minimum config to listen on:
	if cfg.Port == 0 {
		return nil, errors.New("Config.Port must not be zero")
	}
	s := &Site{
		Config: cfg,
		Server: &http.Server{},
		Mux:    http.NewServeMux(),
	}

	// Run the initializers in order:
	if err := s.InitProvider(); err != nil {
		return nil, err
	}
	if err := s.InitMux(); err != nil {
		return nil, err
	}
	if err := s.InitServer(); err != nil {
		return nil, err
	}

	return s, nil
}

// InitProvider intializes the Provider in a Site based on the Config's
// Provider and ProviderConfig properties. The first error encountered is
// returned.
func (s *Site) InitProvider() error {

	// Set up the Provider.
	switch s.Config.Provider {
	case "":
		return errors.New("Provider missing from Config.")
	case "filesystem":
		cfg, err := NewFileSystemProviderConfig(s.Config.ProviderConfig)
		if err != nil {
			return fmt.Errorf("ProviderConfig error: %s", err)
		}
		fsp, err := LoadFileSystemProvider(cfg)
		if err != nil {
			return err
		}
		s.Provider = fsp
		return nil

	default:
		return errors.New("Unsupported Provider: " + s.Config.Provider)
	}

}

// InitMux intializes the Mux in a Site based on the Config and the Provider.
// The first error encountered is returned.
func (s *Site) InitMux() error {

	return nil

}

// InitServer intializes the Server in a Site based on the Config and the
// Mux.  The first error encountered is returned.
func (s *Site) InitServer() error {

	return nil

}
