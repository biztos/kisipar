// kisipar.go -- kisipar top-level description.
// ----------

// Package kisipar implements a general-purpose web server for
// smallish websites.  It directly supports file-based web sites written
// in Markdown as well as several more obscure use-cases; and it is designed
// to be a useful content-delivery layer for more complex sites.
//
// This being the modern Interweb, opinions vary on what "smallish" means,
// and benchmarks are TODO, but a normal-sized blog should easily fit.
//
// A user's guide of sorts is available at:
//
// https://kisipar.biztos.com/
//
// For more technical information please see the project page on GitHub:
//
// https://github.com/biztos/kisipar
package kisipar

import (
	"github.com/biztos/kisipar/site"
)

// NewSite returns a new site.Site initialized with the given config file,
// or the first error encountered.
func NewSite(file string) (*site.Site, error) {

	cfg, err := site.LoadConfig(file)
	if err != nil {
		return nil, err
	}
	s, err := site.NewSite(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}
