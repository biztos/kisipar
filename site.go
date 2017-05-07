// site.go -- A single Kisipar site.

package kisipar

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	s := &Site{Config: cfg}

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
// returned.  If the Provider is already set then no action is taken.
func (s *Site) InitProvider() error {

	if s.Provider != nil {
		return nil
	}

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

// InitMux intializes the Mux (request multiplexer) in a Site based on the
// Config and the Provider. The first error encountered is returned.
// If the Mux is already set then no action is taken.
func (s *Site) InitMux() error {

	// Do nothing if we already are mux-ified.
	if s.Mux != nil {
		return nil
	}

	// Default handler for now.
	// TODO: stuff based on configs (contact, news...)
	handler, err := NewHandler(s)
	if err != nil {
		return err
	}
	s.Mux = http.NewServeMux()
	s.Mux.Handle("/", handler)

	return nil

}

// InitServer intializes the Server in a Site based on the Config and the
// Mux.  The first error encountered is returned.  If the Server is already
// set then no action is taken.
func (s *Site) InitServer() error {

	if s.Server != nil {
		return nil
	}

	if s.Mux == nil {
		return errors.New("Mux must not be nil.")
	}

	s.Server = &http.Server{
		Handler: s.Mux,
		Addr:    fmt.Sprintf(":%d", s.Config.Port),
	}

	return nil

}

// StaticPath looks for the given file in the StaticDir of the Site's
// Config. If StaticDir is relative then it is assumed to be under the Dir
// of the Site's Config. Directories are allowed only if ListStatic is
// true in the Config; if it is false then any dir-like request is given an
// extension of ".html". The path is returned, or any error returned from
// os.Stat; the error is testable with os.IsNotExist.
func (s *Site) StaticPath(file string) (string, error) {

	// You should have a config, normally, but we can live without it here.
	if s.Config == nil {
		return "", os.ErrNotExist
	}

	if s.Config.StaticDir == "" {
		return "", os.ErrNotExist
	}

	dir := s.Config.StaticDir
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(s.Config.Dir, s.Config.StaticDir)
	}

	if s.Config.ListStatic == false {
		file = strings.TrimSuffix(file, "/") + ".html"
	}

	file = filepath.Join(dir, file)

	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	// Disallow directories unless we are being super liberal.
	if s.Config.ListStatic == false && info.IsDir() {
		return "", os.ErrNotExist
	}

	return file, nil

}

// Serve serves the site via ListenAndServe(TLS) according to its
// configuration.  A brief message is sent to the log before listening starts.
func (s *Site) Serve() error {

	if s.Server == nil {
		return errors.New("Server must not be nil.")
	}
	if s.Config == nil {
		return errors.New("Config must not be nil.")
	}

	name := s.Config.Name
	if name == "" {
		name = "anonymous site"
	}
	log.Printf("Kisipar listening on %s for %s.\n", s.Server.Addr, name)

	// TODO: switch for TLS etc if we care.
	return s.Server.ListenAndServe()
}
