// config_test.go -- tests for Kisipar Config
// --------------

package site_test

import (
	// Standard:
	"os"
	"path/filepath"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar/site"
)

func Test_LoadConfig_ErrorNoFile(t *testing.T) {

	assert := assert.New(t)

	_, err := site.LoadConfig("")
	if assert.Error(err) {
		assert.Equal("No config file specified", err.Error(),
			"error is useful")
	}

}

func Test_LoadConfig_ErrorFileNotExist(t *testing.T) {

	assert := assert.New(t)

	_, err := site.LoadConfig("nonesuch.yaml")
	if assert.Error(err) {
		assert.True(os.IsNotExist(err), "error isa IsNotExist")
		assert.Regexp("nonesuch.yaml", err.Error(),
			"error is useful")
	}

}

func Test_LoadConfig_ErrorBadYAML(t *testing.T) {

	assert := assert.New(t)

	file := filepath.Join("testdata", "broken-config.yaml")
	_, err := site.LoadConfig(file)
	if assert.Error(err) {
		assert.False(os.IsNotExist(err), "error nota IsNotExist")
		assert.Regexp("yaml", err.Error(), "error is useful")
	}

}

func Test_LoadConfig_Success(t *testing.T) {

	assert := assert.New(t)

	file := filepath.Join("testdata", "fsp-config.yaml")
	cfg, err := site.LoadConfig(file)
	if !assert.Nil(err, "no error") {
		t.Fatal(err)
	}
	assert.Equal(8080, cfg.Port, "Port set")
	assert.Equal("FSP Test Site", cfg.Name, "Name set")
	assert.Equal("Kisipar", cfg.Owner, "Owner set")
	assert.Equal("filesystem", cfg.Provider, "Provider set")

	pconfig := map[string]interface{}{
		"templates": "fsp-templates",
		"content":   "fsp-content",
	}
	assert.Equal(pconfig, cfg.ProviderConfig, "ProviderConfig set")

}
