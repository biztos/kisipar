// filesystemprovider_test.go -- tests for filesystem provider
// --------------------------

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

func Test_NewFileSystemProvider(t *testing.T) {

	assert := assert.New(t)

	fsp := kisipar.NewFileSystemProvider(kisipar.FileSystemProviderConfig{
		ContentDir: "/any/where",
	})

	// Note that the the update time is initialized.
	assert.Regexp(
		"^<FileSystemProvider with 0 items at /any/where, updated .*>$",
		fsp.String(), "stringifies as expected")

}

func Test_FileSystemProvider_LoadTemplates_NoTemplateDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	assert.Nil(err, "no error")
	assert.Nil(fsp.Template(), "no master template set")

}

func Test_FileSystemProvider_LoadTemplates_NoSuchDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{TemplateDir: "nosuchdir"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp(".*no such file or directory$", err.Error())
	}
}

func Test_FileSystemProvider_LoadTemplates_DirNotDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{TemplateDir: "README.md"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Equal("Not a directory: README.md", err.Error())
	}
}

func Test_FileSystemProvider_LoadTemplates_TemplateError(t *testing.T) {

	assert := assert.New(t)

	dir := filepath.Join("test", "fsp-bad-templates")
	config := kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .* Template", err.Error())
	}
}

func Test_FileSystemProvider_LoadContent_NoSuchDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{ContentDir: "nosuchdir"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Regexp(".*no such file or directory$", err.Error())
	}
}

func Test_FileSystemProvider_LoadContent_DirNotDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{ContentDir: "README.md"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Equal("Not a directory: README.md", err.Error())
	}
}

func Test_FileSystemProvider_LoadContent_Success(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{
		ContentDir: filepath.Join("test", "fsp-content"),
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	assert.Nil(err, "no error")

	exp := []string{
		"/dupe",
		"/foo",
		"/index",
		"/other",
		"/other.txt",
		"/foo/bar",
		"/foo/s.js",
		"/foo/bar/baz",
		"/foo/bother/data.json",
		"/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}
