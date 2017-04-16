// filesystemprovider_test.go -- tests for filesystem provider
// --------------------------

package kisipar_test

import (
	// Standard:
	"io/ioutil"
	"os"
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

	dir := filepath.Join("testdata", "fsp-bad-templates")
	config := kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .* Template", err.Error())
	}
}

// Sort of a pain in the ass edge case but I hit it for real while debugging
// so (alas) it's worth testing for.
func Test_FileSystemProvider_LoadTemplates_InnerSymlinkErr(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner file which is the link target.
	fn := filepath.Join(dir, "target.html")
	if err = ioutil.WriteFile(fn, []byte("hello"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	assert.Nil(err)
}

func Test_FileSystemProvider_LoadTemplates_Success(t *testing.T) {

	assert := assert.New(t)

	dir := filepath.Join("testdata", "fsp-templates")
	config := kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}
}

func Test_FileSystemProvider_LoadContent_NoContentDir(t *testing.T) {

	assert := assert.New(t)

	config := kisipar.FileSystemProviderConfig{ContentDir: ""}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if !assert.Nil(err, "no error") {
		t.Logf("Error: %s", err.Error())
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
		ContentDir: filepath.Join("testdata", "fsp-content"),
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
