// Package assets provides embedded asset data for the Kisipar site package.
//   TODO: get to 100% test coverage of the stupid package.
//   TODO: consider building that into go-bindata with an option -tests
package assets

// MustAssetString is a convenience wrapper for MustAsset, returning a string
// instead of a byte slice.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}
