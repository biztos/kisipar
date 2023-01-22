// bindata_test.go - because 100% dammit.
// ---------------
// TODO: fork bindata, make it do this automatically.

package assets_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/kisipar/site/assets"
)

func Test_AssetNames(t *testing.T) {

	assert := assert.New(t)

	names := assets.AssetNames()

	assert.Contains(names, KNOWN_GOOD_ASSET_NAME, "known good name present")
	assert.NotContains(names, KNOWN_BAD_ASSET_NAME, "known bad name absent")

}

func Test_AssetInfo_Success(t *testing.T) {

	assert := assert.New(t)

	info, err := assets.AssetInfo(KNOWN_GOOD_ASSET_NAME)
	assert.Nil(err, "no error for good asset")
	assert.NotNil(info, "info returned")

	// Exercise the info struct a bit, which is stupid but necessry for our
	// coverage metrics.
	assert.Equal(KNOWN_GOOD_ASSET_NAME, info.Name(), "Name works")
	assert.NotZero(info.Size(), "Size works")
	assert.NotZero(info.Mode(), "Mode works")
	assert.NotZero(info.ModTime(), "ModTime works")
	assert.False(info.IsDir(), "IsDir works")
	assert.Nil(info.Sys(), "Sys works")
}

func Test_AssetInfo_ErrorNotFound(t *testing.T) {

	assert := assert.New(t)

	info, err := assets.AssetInfo(KNOWN_BAD_ASSET_NAME)
	if assert.Error(err, "error returned") {
		assert.Regexp("^AssetInfo .* not found$", err.Error(), "error useful")
	}
	assert.Nil(info, "no info returned")

}

func Test_Asset_Success(t *testing.T) {

	assert := assert.New(t)

	b, err := assets.Asset(KNOWN_GOOD_ASSET_NAME)
	assert.Nil(err, "no error")
	assert.NotZero(b, "asset data returned")
}

func Test_Asset_ErrorNotFound(t *testing.T) {

	assert := assert.New(t)

	b, err := assets.Asset(KNOWN_BAD_ASSET_NAME)
	if assert.Error(err, "error returned") {
		assert.Regexp("^Asset .* not found$", err.Error(), "error useful")
	}
	assert.Zero(b, "asset data empty")
}

func Test_MustAsset_Success(t *testing.T) {

	assert := assert.New(t)

	b := assets.MustAsset(KNOWN_GOOD_ASSET_NAME)
	assert.NotZero(b, "asset data returned (no panic)")
}

func Test_MustAsset_Panic(t *testing.T) {

	assert := assert.New(t)

	var b []byte
	assert.Panics(func() {
		b = assets.MustAsset(KNOWN_BAD_ASSET_NAME)
	}, "panic for missing asset")
	assert.Zero(b, "no data returned")

}

func Test_AssetDir_ErrorNotFound(t *testing.T) {

	assert := assert.New(t)

	names, err := assets.AssetDir("/no/such/thing/here")
	if assert.Error(err, "error returned") {
		assert.Regexp("^Asset .* not found$", err.Error(), "error useful")
	}
	assert.Zero(names, "names empty")
}

func Test_AssetDir_Success(t *testing.T) {

	assert := assert.New(t)

	names, err := assets.AssetDir(filepath.Dir(KNOWN_GOOD_ASSET_NAME))
	assert.Nil(err, "no error returned")
	assert.Contains(names, filepath.Base(KNOWN_GOOD_ASSET_NAME),
		"good asset found")
}
