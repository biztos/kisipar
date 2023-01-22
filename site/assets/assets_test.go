package assets_test

import (
	"github.com/biztos/kisipar/site/assets"
	"github.com/stretchr/testify/assert"
	"testing"
)

const KNOWN_GOOD_ASSET_NAME = "demosite/templates/default.html"
const KNOWN_BAD_ASSET_NAME = "no-such-thing-here"

func Test_MustAssetString_Success(t *testing.T) {

	assert := assert.New(t)

	var s string
	assert.NotPanics(func() {
		s = assets.MustAssetString(KNOWN_GOOD_ASSET_NAME)
	}, "no panic from MustAssetString")
	assert.NotZero(s, "data returned")
}

func Test_MustAssetString_Panic(t *testing.T) {

	assert := assert.New(t)

	var s string
	assert.Panics(func() {
		s = assets.MustAssetString(KNOWN_BAD_ASSET_NAME)
	}, "panic from MustAssetString")
	assert.Zero(s, "no data returned")
}
