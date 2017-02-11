// bindata_test.go -- tests for generated bindata
//
// TODO: generate this!  Use Perl if needed.
package kisipar_test

import (
	// Standard:
	"io/ioutil"
	"os"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

// Testability is part of the craft.  The bindata folks did us a fine service
// but fell short here.
func Test_CrutchUnderBindata(t *testing.T) {

	assert := assert.New(t)

	assert.NotPanics(kisipar.CrutchUnderBindata, "no panic")
}

func Test_Asset(t *testing.T) {

	assert := assert.New(t)

	_, err := kisipar.Asset("foo")
	if assert.Error(err, "error for missing asset") {
		assert.Equal("Asset foo not found", err.Error(), "error as expected")
	}

	b, err := kisipar.Asset("templates/default.html")
	if assert.Nil(err, "no error for known asset") {
		assert.NotZero(b, "data returned")

	}
}

func Test_MustAsset(t *testing.T) {

	assert := assert.New(t)
	panicky := func() { kisipar.MustAsset("bla") }
	nopanic := func() { kisipar.MustAsset("templates/default.html") }

	assert.Panics(panicky, "panic for missing asset")
	assert.NotPanics(nopanic, "no panic for known asset")
}

func Test_AssetInfo(t *testing.T) {

	assert := assert.New(t)

	_, err := kisipar.AssetInfo("foo")
	if assert.Error(err, "error for missing asset") {
		assert.Equal("AssetInfo foo not found", err.Error(),
			"error as expected")
	}

	info, err := kisipar.AssetInfo("templates/default.html")
	if assert.Nil(err, "no error for known asset") {
		assert.NotNil(info, "info returned")
		assert.NotZero(info.Name(), "Name not zero")
		assert.NotZero(info.Size(), "Size not zero")
		assert.NotZero(info.Mode(), "Mode not zero")
		assert.NotZero(info.ModTime(), "ModTime not zero")
		assert.False(info.IsDir(), "IsDir is false")

		// Apparently no underlying data source.
		// (When would that not be nil?  Curious.)
		assert.Nil(info.Sys(), "Sys is nil")
	}
}

func Test_AssetNames(t *testing.T) {

	assert := assert.New(t)

	names := kisipar.AssetNames()
	assert.NotZero(names, "have names")
}

func Test_AssetDir(t *testing.T) {

	assert := assert.New(t)

	_, err := kisipar.AssetDir("foo")
	if assert.Error(err, "error for missing asset") {
		assert.Equal("Asset foo not found", err.Error(),
			"error as expected")
	}

	assets, err := kisipar.AssetDir("templates")
	if assert.Nil(err, "no error for known asset dir") {
		assert.NotZero(assets, "have assets")
	}

}

func Test_RestoreAsset(t *testing.T) {

	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "kisipar_test_bindata_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = kisipar.RestoreAsset(dir, "foo")
	if assert.Error(err, "error for missing asset") {
		assert.Equal("Asset foo not found", err.Error(),
			"error as expected")
	}

	err = kisipar.RestoreAsset(dir, "templates/default.html")
	if assert.Nil(err, "no error") {
		// TODO: Check that we have what we want in the dir.

	}
}

func Test_RestoreAssets(t *testing.T) {

	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "kisipar_test_bindata_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = kisipar.RestoreAssets(dir, "templates")
	if assert.Nil(err, "no error") {
		// TODO: Check that we have what we want in the dir.

	}
}
