// binsanity.go - auto-generated; edit at thine own peril!
//
// More info: https://github.com/biztos/binsanity

package kisipar

import "fmt"

// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {
	if b := data[name]; b != nil {
		return b, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset returns the byte content of the asset for the given name, or
// panics if no such asset is available.
func MustAsset(name string) []byte {
	b, err := Asset(name)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// AssetNames returns the alpha-sorted names of the assets.
func AssetNames() []string {
	return names
}

// The names, sorted:
var names = []string{
	"templates/README.md",
	"templates/debug/default.html",
	"templates/default/default.html",
	"templates/naked/default.html",
	"templates/wonky/default.html",
}

// The data itself (long lines ahead):
var data = map[string][]byte{
	"templates/README.md": []byte{0x23, 0x20, 0x52, 0x45, 0x41, 0x44, 0x4d, 0x45, 0x3a, 0x20, 0x6b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x2d, 0x69, 0x6e, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x2e, 0xa, 0xa, 0x54, 0x68, 0x65, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x2d, 0x69, 0x6e, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x20, 0x61, 0x72, 0x65, 0x20, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x65, 0x64, 0x20, 0x69, 0x6e, 0x20, 0x22, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x22, 0x20, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x20, 0x20, 0x4f, 0x6e, 0x6c, 0x79, 0x20, 0x74, 0x68, 0x6f, 0x73, 0x65, 0xa, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x20, 0x73, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x20, 0x62, 0x65, 0x20, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x20, 0x69, 0x6e, 0x74, 0x6f, 0x20, 0x61, 0x20, 0x67, 0x69, 0x76, 0x65, 0x6e, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x73, 0x65, 0x74, 0x2e, 0xa, 0xa, 0x42, 0x75, 0x69, 0x6c, 0x74, 0x2d, 0x69, 0x6e, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x20, 0x6d, 0x6f, 0x73, 0x74, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x61, 0x73, 0x73, 0x75, 0x6d, 0x65, 0x20, 0x61, 0x6e, 0x79, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x20, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x20, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x20, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x3b, 0x20, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x73, 0xa, 0x61, 0x6e, 0x64, 0x20, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x73, 0x20, 0x6d, 0x75, 0x73, 0x74, 0x20, 0x62, 0x65, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x20, 0x69, 0x6e, 0x2e, 0xa, 0xa, 0x54, 0x68, 0x65, 0x20, 0x67, 0x6f, 0x61, 0x6c, 0x20, 0x69, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x61, 0x62, 0x6c, 0x79, 0x20, 0x75, 0x73, 0x65, 0x66, 0x75, 0x6c, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x74, 0x68, 0x65, 0x73, 0x65, 0x20, 0x75, 0x73, 0x65, 0x20, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3a, 0xa, 0xa, 0x31, 0x2e, 0x20, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x20, 0x61, 0x20, 0x62, 0x6c, 0x6f, 0x67, 0x20, 0x28, 0x6f, 0x72, 0x20, 0x73, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x29, 0x20, 0x77, 0x69, 0x74, 0x68, 0x6f, 0x75, 0x74, 0x20, 0x61, 0x6e, 0x79, 0x20, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x2e, 0xa, 0x32, 0x2e, 0x20, 0x44, 0x65, 0x62, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0xa, 0x33, 0x2e, 0x20, 0x44, 0x65, 0x62, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x20, 0x6b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x20, 0x69, 0x74, 0x73, 0x65, 0x6c, 0x66, 0x2e, 0xa},
	"templates/debug/default.html": []byte{0x7b, 0x7b, 0x2f, 0x2a, 0xa, 0x20, 0x20, 0x20, 0x20, 0xa, 0x20, 0x20, 0x20, 0x20, 0x64, 0x65, 0x62, 0x75, 0x67, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x2d, 0x20, 0x22, 0x64, 0x65, 0x62, 0x75, 0x67, 0x22, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x4b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x68, 0x65, 0x20, 0x64, 0x65, 0x62, 0x75, 0x67, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x69, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x2e, 0x2e, 0x2e, 0x20, 0x64, 0x65, 0x62, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x21, 0x20, 0x44, 0x4f, 0x20, 0x4e, 0x4f, 0x54, 0x20, 0x55, 0x53, 0x45, 0x20, 0x69, 0x6e, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2d, 0x66, 0x61, 0x63, 0x69, 0x6e, 0x67, 0x20, 0x73, 0x69, 0x74, 0x65, 0x73, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x28, 0x4f, 0x62, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x6c, 0x79, 0x2e, 0x29, 0xa, 0xa, 0x2a, 0x2f, 0x7d, 0x7d, 0x3c, 0x21, 0x64, 0x6f, 0x63, 0x74, 0x79, 0x70, 0x65, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x3c, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x20, 0x74, 0x79, 0x70, 0x65, 0x3d, 0x22, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x63, 0x73, 0x73, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x62, 0x6f, 0x64, 0x79, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x66, 0x6f, 0x6e, 0x74, 0x2d, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x3a, 0x20, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x3b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x57, 0x4f, 0x4e, 0x4b, 0x59, 0x3a, 0x20, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x64, 0x69, 0x76, 0x20, 0x69, 0x64, 0x3d, 0x22, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0x7b, 0x20, 0x2e, 0x48, 0x74, 0x6d, 0x6c, 0x20, 0x7d, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x64, 0x69, 0x76, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x3c, 0x2f, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa},
	"templates/default/default.html": []byte{0x7b, 0x7b, 0x2f, 0x2a, 0xa, 0x20, 0x20, 0x20, 0x20, 0xa, 0x20, 0x20, 0x20, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x2d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x4f, 0x6e, 0x65, 0x20, 0x54, 0x72, 0x75, 0x65, 0x20, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x4b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x68, 0x65, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x67, 0x68, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x2c, 0x20, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2d, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x73, 0x65, 0x74, 0x20, 0x75, 0x73, 0x69, 0x6e, 0x67, 0xa, 0x20, 0x20, 0x20, 0x20, 0x61, 0x20, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x20, 0x6f, 0x66, 0x20, 0x4a, 0x61, 0x76, 0x61, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x2a, 0x6e, 0x6f, 0x2a, 0x20, 0x63, 0x6f, 0x6f, 0x6b, 0x69, 0x65, 0x73, 0x20, 0x28, 0x68, 0x65, 0x6e, 0x63, 0x65, 0x20, 0x6e, 0x6f, 0x20, 0x63, 0x6f, 0x6f, 0x6b, 0x69, 0x65, 0x20, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x29, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x69, 0x74, 0x73, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2c, 0x20, 0x77, 0x68, 0x69, 0x63, 0x68, 0x20, 0x69, 0x73, 0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x61, 0x6e, 0x79, 0x20, 0x50, 0x61, 0x67, 0x65, 0x20, 0x74, 0x68, 0x61, 0x74, 0x20, 0x68, 0x61, 0x73, 0x20, 0x6e, 0x6f, 0xa, 0x20, 0x20, 0x20, 0x20, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x20, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0xa, 0xa, 0x2a, 0x2f, 0x7d, 0x7d, 0x3c, 0x21, 0x64, 0x6f, 0x63, 0x74, 0x79, 0x70, 0x65, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x3c, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x48, 0x74, 0x6d, 0x6c, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x3c, 0x2f, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa},
	"templates/naked/default.html": []byte{0x7b, 0x7b, 0x2f, 0x2a, 0xa, 0x20, 0x20, 0x20, 0x20, 0xa, 0x20, 0x20, 0x20, 0x20, 0x6e, 0x61, 0x6b, 0x65, 0x64, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x2d, 0x20, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x72, 0x64, 0x20, 0x22, 0x6e, 0x61, 0x6b, 0x65, 0x64, 0x22, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x4b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x68, 0x65, 0x20, 0x6e, 0x61, 0x6b, 0x65, 0x64, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x6d, 0x61, 0x6b, 0x65, 0x73, 0x20, 0x6e, 0x6f, 0x20, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x20, 0x74, 0x6f, 0x20, 0x77, 0x72, 0x61, 0x70, 0x20, 0x6f, 0x72, 0x20, 0x62, 0x65, 0x61, 0x75, 0x74, 0x69, 0x66, 0x79, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0xa, 0x20, 0x20, 0x20, 0x20, 0x48, 0x54, 0x4d, 0x4c, 0x2e, 0xa, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x4f, 0x44, 0x4f, 0x3a, 0x20, 0x61, 0x20, 0x6e, 0x61, 0x6b, 0x65, 0x64, 0x20, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x20, 0x70, 0x61, 0x67, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x73, 0x6f, 0x72, 0x74, 0x3a, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x65, 0x76, 0x65, 0x72, 0x79, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x20, 0x62, 0x65, 0x6c, 0x6f, 0x77, 0xa, 0x20, 0x20, 0x20, 0x20, 0x28, 0x6d, 0x65, 0x61, 0x6e, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x20, 0x68, 0x61, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x67, 0x6f, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x20, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x29, 0x2e, 0xa, 0xa, 0x2a, 0x2f, 0x7d, 0x7d, 0x3c, 0x21, 0x64, 0x6f, 0x63, 0x74, 0x79, 0x70, 0x65, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x3c, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x48, 0x74, 0x6d, 0x6c, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x3c, 0x2f, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa},
	"templates/wonky/default.html": []byte{0x7b, 0x7b, 0x2f, 0x2a, 0xa, 0x20, 0x20, 0x20, 0x20, 0xa, 0x20, 0x20, 0x20, 0x20, 0x77, 0x6f, 0x6e, 0x6b, 0x79, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x2d, 0x20, 0x22, 0x77, 0x6f, 0x6e, 0x6b, 0x79, 0x22, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x4b, 0x69, 0x73, 0x69, 0x70, 0x61, 0x72, 0x2e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x54, 0x68, 0x65, 0x20, 0x77, 0x6f, 0x6e, 0x6b, 0x79, 0x20, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x69, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x20, 0x20, 0x59, 0x4d, 0x4d, 0x56, 0x3b, 0x20, 0x45, 0x26, 0x4f, 0x45, 0x3b, 0x20, 0x6e, 0x6f, 0x20, 0x67, 0x75, 0x61, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65, 0x73, 0x2e, 0xa, 0xa, 0x2a, 0x2f, 0x7d, 0x7d, 0x3c, 0x21, 0x64, 0x6f, 0x63, 0x74, 0x79, 0x70, 0x65, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x3c, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x20, 0x74, 0x79, 0x70, 0x65, 0x3d, 0x22, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x63, 0x73, 0x73, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x62, 0x6f, 0x64, 0x79, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x62, 0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x3a, 0x20, 0x79, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x3b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x57, 0x4f, 0x4e, 0x4b, 0x59, 0x3a, 0x20, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x64, 0x69, 0x76, 0x20, 0x69, 0x64, 0x3d, 0x22, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0x7b, 0x20, 0x2e, 0x48, 0x74, 0x6d, 0x6c, 0x20, 0x7d, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x64, 0x69, 0x76, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x2f, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x3c, 0x2f, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa},
}
