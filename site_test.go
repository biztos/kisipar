// site_test.go

package kisipar_test

import (
	// Standard:
	"path/filepath"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_NewSite_ErrorNoConfig(t *testing.T) {

	assert := assert.New(t)

	_, err := kisipar.NewSite(nil)
	if assert.Error(err) {
		assert.Equal("Config must not be nil", err.Error(), "error useful")
	}

}

func Test_NewSite_ErrorNoPort(t *testing.T) {

	assert := assert.New(t)

	cfg := &kisipar.Config{}
	_, err := kisipar.NewSite(cfg)
	if assert.Error(err) {
		assert.Equal("Config.Port must not be zero", err.Error(),
			"error useful")
	}

}

func Test_NewSite_ErrorNoProvider(t *testing.T) {

	assert := assert.New(t)

	cfg := &kisipar.Config{Port: 1234}
	_, err := kisipar.NewSite(cfg)
	if assert.Error(err) {
		assert.Equal("Provider missing from Config.", err.Error(),
			"error useful")
	}

}

func Test_NewSite_ErrorUnsupportedProvider(t *testing.T) {

	assert := assert.New(t)

	cfg := &kisipar.Config{Provider: "nonesuch", Port: 1111}
	_, err := kisipar.NewSite(cfg)
	if assert.Error(err) {
		assert.Equal("Unsupported Provider: nonesuch", err.Error(),
			"error useful")
	}

}

func Test_NewSite_ErrorBadProviderConfig(t *testing.T) {

	assert := assert.New(t)

	cfg := &kisipar.Config{
		Port:           1234,
		Provider:       "filesystem",
		ProviderConfig: map[string]interface{}{"Fremd": "Körper"},
	}
	_, err := kisipar.NewSite(cfg)
	if assert.Error(err) {
		assert.Equal("ProviderConfig error: Unexpected FileSystemProviderConfig key: Fremd",
			err.Error(), "error useful")
	}

}

func Test_NewSite_ErrorBadTemplates(t *testing.T) {

	assert := assert.New(t)

	cfg := &kisipar.Config{
		Port:     1234,
		Provider: "filesystem",
		ProviderConfig: map[string]interface{}{
			"TemplateDir": filepath.Join("testdata", "fsp-bad-templates"),
		},
	}
	_, err := kisipar.NewSite(cfg)
	if assert.Error(err) {
		assert.Regexp("Template", err.Error(), "error useful")
	}

}

func Test_NewSite_Success(t *testing.T) {

	assert := assert.New(t)

	// A FileSystemProvider can have an empty config, it will just end up
	// with the standard internal template.  Additional InitProvider
	// scenarios are tested separately.
	cfg := &kisipar.Config{
		Port:     1234,
		Provider: "filesystem",
	}
	s, err := kisipar.NewSite(cfg)
	if !assert.Nil(err, "no error") {
		assert.FailNow(err.Error())
	}
	if assert.NotNil(s, "got Site") {
		assert.Equal(cfg, s.Config, "Config kept")
		p, ok := s.Provider.(*kisipar.FileSystemProvider)
		if assert.True(ok, "Provider set") {
			t.Log(p)
		}
		// TODO: mux, server!
	}

}
