// site/site.go
// ------------

// Package site defines a Kisipar web site, and contains most of its actual
// server logic.
package site

import (
	// Standard library:
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	// Third-party packages:
	"github.com/olebedev/config"

	// Kisipar packages:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
	"github.com/biztos/kisipar/site/assets"
)

var DEFAULT_NAME = "Anonymous Kisipar Site"
var DEFAULT_OWNER = "Anonymous Kisipar Fan"
var DEFAULT_PORT = 8020
var DEFAULT_PAGE_EXTENSIONS = []string{
	".md",
	".MD", // Trying to be a good MSDOS citizen here...
	".markdown",
	".MARKDOWN",
	".txt",
	".TXT",
}

// TODO: instead of just having one default template, have a whole set of them
// that can be available to other templates, under the "kisipar" namespace
// in the loaded templates.  (Still have the minimalist default though.)
var DEFAULT_TEMPLATE = assets.MustAssetString("demosite/templates/default.html")

var DEFAULT_READ_TIMEOUT = 10 * time.Second
var DEFAULT_WRITE_TIMEOUT = 10 * time.Second
var DEFAULT_MAX_HEADER_BYTES = 1 << 20

var DEFAULT_FEED_PATH = "/feed.xml"
var DEFAULT_FEED_ITEMS = 20

// A SiteServer can be a custom implementation as long as it provides the
// standard APIs for serving with and without Transport Layer Security.
type SiteServer interface {
	ListenAndServe() error
	ListenAndServeTLS(string, string) error
}

// A Site represents a Kisipar web site ready for serving.
type Site struct {

	// The Path is the base path of the site.  Kisipar sites are contained in
	// single directories.
	Path string

	// PagePath is the path under which Page content is located, i.e. the
	// Markdown files and (optionally) their resources.
	PagePath string

	// PageExtensions are the recognized page content file extensions.  Files
	// with these extensions are parsed as Pages.
	PageExtensions []string

	// UnlistedPaths are the path prefixes under which to automatically set
	// Pages to Unlisted.
	UnlistedPaths []string

	// ServePageSources determines whether the source files (e.g. "foo.md")
	// can be served as assets when requested directly.  If false, page assets
	// with extensions listed in PageExtensions will be treated as Not Found.
	ServePageSources bool

	// TemplatePath is the path under which Templates are located.
	TemplatePath string

	// StaticPath is the path under which Static content is located.
	StaticPath string

	// The Name can be anything you want, and is the primary identifier of the
	// site in the logs.  It is also used in templates and news feeds.
	Name string

	// The Owner is the display name of the site owner, e.g. "John Q. Doe".
	Owner string

	// The Email is where contact notices, if any, are sent.
	Email string

	// The Host should be the *external* host name, such as "example.com" --
	// this is used to construct a standard BaseURL if none is specified.
	Host string

	// The BaseURL is used for generating links to the site's pages.
	BaseURL string

	// FeedPath is the URL path to the site's Atom feed, if available.
	// FeedTitle is the Title to use in the Feed, if not the site's Name.
	// FeedItems specifies the maximum number of items to include in the
	// feed.  To turn off this feature, set NoFeed to a true value.
	// NOTE: the feed excludes unlisted pages; cf. UnlistedPaths.
	// TODO: (maybe) support multiple feeds, e.g. "/foo/*" vs "/bar/*"
	FeedPath  string
	FeedTitle string
	FeedItems int
	NoFeed    bool

	// The Port determines where the server will listen, and ServeTLS dictates
	// whether we listen on HTTP or HTTPS.
	Port     int
	ServeTLS bool

	// If ServeTLS is true, then CertFile and KeyFile must contain paths to
	// a valid cert and key, respectively.
	CertFile string
	KeyFile  string

	// Timeout and header-reading limits for the Server include:
	ReadTimeout    time.Duration // Config is in seconds.
	WriteTimeout   time.Duration // Config is in seconds.
	MaxHeaderBytes int

	// The Config is used to set the properties above.
	Config *config.Config

	// The Server is used to serve the Site; by default it will be set to
	// a standard http.Server using the configurable properties above.
	Server SiteServer

	// The Pageset contains all the known Pages for dynamic generation.
	Pageset *pageset.Pageset

	// The Template containing all the shared templates as well as all the
	// specific page templates.
	Template *template.Template
}

// New initializes a Site at the given directory path.  A config file in YAML
// format, named "config.yaml", is sought under the path.  If path is the
// empty string or no config file is present then a blank Config is used.
//
// Expected config values are:
//
//   Name           # Name of the site; default: Anonymous Kisipar Site
//   Owner          # Owner of the site; default: Anonymous Kisipar Fan
//   Host           # Host from which we're serving; default: localhost
//   Port           # Port on which to serve; default: 8020
//   BaseURL        # BaseURL for all site URLs; default: derived.
//   CertFile       # TLS only: path to the cert file
//   KeyFile        # TLS only: path to the key file
//   PagePath       # relative path for pages; default: pages
//   UnlistedPaths  # path (prefixes) for unlisted pages
//   TemplatePath   # relative path for templates; default: templates
//   StaticPath     # relative path for static content; default: static
//   FeedPath       # URL path for Atom feed; standard default: /feed.xml
//   FeedTitle      # Title for the Atom feed, if not the site Name
//   FeedItems      # Number of items in the Atom feed; standard default: 20
//   NoFeed         # boolean switch to disable the Atom feed
//
// Sensible, but not necessarily perfect, defaults are calculated as needed;
// most can be overridden via the package variables.
//
// All configs are available via the site's Config property.
func New(path string) (*Site, error) {

	var s *Site
	if path == "" {
		s = &Site{Config: config.Must(config.ParseYaml(""))}
	} else {

		dirInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if !dirInfo.IsDir() {
			return nil, fmt.Errorf("%s: not a directory", path)
		}

		s = &Site{Path: path}
		if err := s.setConfig(); err != nil {
			return nil, fmt.Errorf("Config error: %s", err.Error())
		}
	}

	// There appears to be a bug in the config package, where it can get
	// a nil root which then causes lookups to fail "wrongly."
	// TODO: pull request to Lebedev for this, it should be easy to fix.
	if s.Config.Root == nil {
		s.Config.Root = map[string]interface{}{}
	}

	if err := s.setup(); err != nil {
		return nil, err
	}

	return s, nil

}

func (s *Site) setConfig() error {

	// The config file must be in YAML.  Deal with it.
	yamlPath := filepath.Join(s.Path, "config.yaml")
	if _, err := os.Stat(yamlPath); !os.IsNotExist(err) {
		// We have a YAML config.
		cfg, err := config.ParseYamlFile(yamlPath)
		if err != nil {
			return err
		}
		s.Config = cfg
		return nil
	}

	// We are sadly lacking a config, but we are able to run on defaults
	// if necessary.
	s.Config = config.Must(config.ParseYaml(""))
	return nil

}

func (s *Site) setup() error {

	// Do we have directories?
	if s.Path == "" {
		s.PagePath = ""
		s.TemplatePath = ""
		s.StaticPath = ""
	} else {
		// TODO: consider only defaulting these if the directories are
		// present, otherwise having (e.g.) no templates or no static.
		pdir := s.Config.UString("PagePath", "pages")
		tdir := s.Config.UString("TemplatePath", "templates")
		sdir := s.Config.UString("StaticPath", "static")

		// NOTE: there is no support for full paths here; everything is by
		// definition relative to the site Path.
		s.PagePath = filepath.Join(s.Path, pdir)
		s.TemplatePath = filepath.Join(s.Path, tdir)
		s.StaticPath = filepath.Join(s.Path, sdir)

	}

	// What files are considered Pages?
	// (Be a little paranoid here in case of user error defining the list.)
	pExt, err := s.configStringList("PageExtensions")
	if err != nil {
		return err
	} else if len(pExt) == 0 {
		pExt = DEFAULT_PAGE_EXTENSIONS
	}
	s.PageExtensions = pExt
	page.LimitExtParsers(pExt)

	// Is anything Unlisted based on its path?
	s.UnlistedPaths, err = s.configStringList("UnlistedPaths")
	if err != nil {
		return err
	}

	// Shall we serve sources?
	s.ServePageSources = s.Config.UBool("ServePageSources", false)

	// Standard Server setup:
	s.Port = s.Config.UInt("Port", DEFAULT_PORT)
	s.ReadTimeout = s.cfgDuration("ReadTimeout", DEFAULT_READ_TIMEOUT)
	s.WriteTimeout = s.cfgDuration("WriteTimeout", DEFAULT_WRITE_TIMEOUT)
	s.MaxHeaderBytes = s.Config.UInt("MaxHeaderBytes", DEFAULT_MAX_HEADER_BYTES)

	// Are we secure?
	s.ServeTLS = false
	s.CertFile = s.Config.UString("CertFile", "")
	s.KeyFile = s.Config.UString("KeyFile", "")
	if s.CertFile != "" && s.KeyFile != "" {
		s.ServeTLS = true
	}

	// Properties used in templates and other magic:
	s.Name = s.Config.UString("Name", DEFAULT_NAME)
	s.Host = s.Config.UString("Host", "localhost")
	s.Owner = s.Config.UString("Owner", DEFAULT_OWNER)
	s.Email = s.Config.UString("Email", "")
	baseurl := s.Config.UString("BaseURL", "")
	if baseurl == "" {
		scheme := "http"
		if s.ServeTLS {
			scheme += "s"
		}
		if s.Host == "localhost" {
			baseurl = fmt.Sprintf("%s://localhost:%d", scheme, s.Port)
		} else {
			baseurl = fmt.Sprintf("%s://%s", scheme, s.Host)
		}
	}
	s.BaseURL = strings.TrimSuffix(baseurl, "/")

	// Shall we serve a news feed?
	s.NoFeed = s.Config.UBool("NoFeed", false)
	if !s.NoFeed {
		s.FeedPath = s.Config.UString("FeedPath", DEFAULT_FEED_PATH)
		s.FeedTitle = s.Config.UString("FeedTitle", s.Name)
		s.FeedItems = s.Config.UInt("FeedItems", DEFAULT_FEED_ITEMS)
	}

	// The Server needs a sane default of course; in very custom situations,
	// of which Testing is the most obvious, it may be overridden.
	s.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.Port),
		Handler:        s.NewServeMux(),
		ReadTimeout:    s.ReadTimeout,
		WriteTimeout:   s.WriteTimeout,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	return nil

}

func (s *Site) cfgDuration(key string, def time.Duration) time.Duration {

	secs := s.Config.UInt(key)
	if secs == 0 {
		return def
	}

	return time.Duration(secs) * time.Second

}

// It's a bit annoying that the config package doesn't do this.  Maybe worth
// a pull request.
func isConfigTypeError(e error) bool {

	return strings.HasPrefix(e.Error(), "Type mismatch")
}

func (s *Site) configStringList(key string) ([]string, error) {

	if s.Config.Root == nil {
		return []string{}, nil
	}

	v, err := s.Config.List(key)
	if err != nil {
		if isConfigTypeError(err) {
			return nil, fmt.Errorf("Config %s is not a list.", key)
		} else {
			return []string{}, nil
		}
	}

	list := make([]string, len(v))
	for i, val := range v {
		if str, ok := val.(string); ok {
			list[i] = str
		} else {
			return nil, fmt.Errorf(
				"Config %s must be a list of strings; item %d is %T.",
				key, i, val)
		}
	}

	return list, nil
}

// URL returns a full URL for the given path, based on the Site's BaseURL.
func (s *Site) URL(path string) string {

	return strings.TrimSuffix(s.BaseURL, "/") + "/" +
		strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")

}

// Href returns the URL path (realtive URL) for a Page.
func (s *Site) Href(p *page.Page) string {

	if p == nil {
		return ""
	}

	rpath := strings.TrimPrefix(p.Path, s.PagePath)
	if p.IsIndex {
		rpath = filepath.Dir(rpath)
	} else {
		rpath = strings.TrimSuffix(rpath, filepath.Ext(rpath))
	}

	return rpath

}

// PageURL returns the full URL for a Page, based on the Site's BaseURL.
func (s *Site) PageURL(p *page.Page) string {
	if p == nil {
		return ""
	}

	rpath := strings.TrimPrefix(p.Path, s.PagePath)
	if p.IsIndex {
		rpath = filepath.Dir(rpath)
	} else {
		rpath = strings.TrimSuffix(rpath, filepath.Ext(rpath))
	}
	return strings.TrimSuffix(s.BaseURL, "/") +
		"/" +
		strings.TrimPrefix(rpath, "/")

}
