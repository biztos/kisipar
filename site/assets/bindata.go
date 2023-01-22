// Code generated by go-bindata.
// sources:
// data/demosite/templates/default.html
// DO NOT EDIT!

package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _demositeTemplatesDefaultHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x54\xcd\x6e\xdb\x30\x0c\xbe\xef\x29\x58\x2f\xd8\xc9\x3f\xdd\x06\xec\xe0\x2a\x19\xb6\x1e\xd6\xcb\x82\x02\xdb\x65\x47\x4d\x66\x22\x61\xb2\x64\x48\x4a\x53\x23\xf0\xbb\x8f\xb2\xdb\x24\x4e\xdc\x44\x17\x9b\xe2\xdf\x47\xf2\x13\xd9\x4d\x65\x45\x68\x1b\x04\x19\x6a\xbd\x78\xc7\x86\x0f\xd0\x61\x12\x79\x35\xfc\xf6\x62\x50\x41\xe3\x62\xb7\x83\xad\x0a\x12\xf2\x47\xbe\x46\xe8\x3a\x92\xf3\xdf\x51\x43\xff\x90\x01\x89\x68\xaa\x97\xfb\x5f\x2a\x60\xbe\xe4\x75\xd4\xb1\x62\xf0\x3f\xc4\xf3\xa1\x25\xaf\x98\x7b\x9e\x04\x7c\x0e\x85\xf0\x3e\x39\xe8\xe3\xf9\x6b\xab\x16\x76\xa3\xab\x78\x56\xd6\x84\x6c\xc5\x6b\xa5\xdb\x12\x7e\xa0\x75\x6b\xc5\x53\xf0\xe8\xd4\xea\xee\xcc\x38\x86\xce\xb8\x56\x6b\x53\x82\x40\x13\xd0\x8d\x6d\xba\x91\xd4\x38\x4c\x41\xd8\x0a\xaf\xa5\xfd\x89\x46\xdb\x14\x6a\x6b\xac\x6f\xb8\xc0\xf3\xc4\xc2\x6a\xeb\x4a\x58\x3b\x44\x73\x29\xa5\xfc\x98\xca\x4f\xa9\xfc\x7c\x2d\x63\xf2\x80\xfa\x09\x83\x12\x1c\x96\xb8\xc1\x24\x85\xfd\x05\x15\xcf\x8d\xcf\x26\x3a\x30\x4e\xf5\x3e\x0e\x2d\x1d\x3e\x1e\xc3\x44\xc6\xad\xaa\x82\x2c\xe1\xcb\xed\x6d\xf3\x7c\xb9\x97\x1a\x57\xe1\xdc\xa2\xe6\x34\x0c\xd2\x92\x3b\xf0\x4d\xb0\x17\xd1\xdc\xdb\xa6\x75\x6a\x2d\xa7\x80\xf4\xa5\xf7\x14\x29\x41\x05\xca\x29\xde\x0a\xc5\x8a\xde\xec\x85\xb4\xc5\x81\xb5\x2c\xd2\xe7\x88\x70\x37\x59\x06\x82\x1b\xd8\x22\xfc\x43\x6c\x68\x40\x75\x4d\x84\xf0\x5f\x21\xcb\x0e\x66\xa7\x04\x3f\xf8\x57\xea\x09\x54\x35\x4f\xa2\xe2\x84\xa8\x91\xed\xf7\x84\x98\xc2\x8d\x5c\x0a\xf2\x19\x85\x1e\x1e\xc7\xf1\xcd\xcc\xd3\x2b\x81\x72\x3e\x3c\x97\x13\x65\x8f\x64\xd6\xf8\x5e\xff\x3a\xb5\xb7\x30\x91\xee\x08\x16\xdb\xe8\x51\x66\xc7\x0d\xd5\x33\x6b\x62\x28\x8a\x98\x7f\x6f\x1f\x39\x05\xef\xc6\x33\x61\x5a\x2d\x18\x07\xe9\x70\x35\x4f\x5e\xd1\xe5\x0f\x24\x46\xd7\xae\x4b\xe2\x02\x98\x35\xfb\x27\xcf\x0a\xbe\x60\x05\x39\x5d\x28\x92\x15\xc7\x50\xae\xf6\x64\x5f\xd3\x9e\x1e\x27\xcd\xfe\x20\x48\x71\xd7\xf7\x7c\x69\xb7\xf9\x1f\xe4\x2e\x2e\x9f\xd3\x81\x9c\xad\xa9\x6f\x9b\x20\xad\x1b\x04\xd4\x1e\x27\xb6\xd4\x14\xf8\x3d\x5c\x56\x0c\x84\x22\x8e\xf5\x4b\xf2\x7f\x00\x00\x00\xff\xff\x17\xa1\xe8\x9b\x3c\x05\x00\x00")

func demositeTemplatesDefaultHtmlBytes() ([]byte, error) {
	return bindataRead(
		_demositeTemplatesDefaultHtml,
		"demosite/templates/default.html",
	)
}

func demositeTemplatesDefaultHtml() (*asset, error) {
	bytes, err := demositeTemplatesDefaultHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "demosite/templates/default.html", size: 1340, mode: os.FileMode(420), modTime: time.Unix(1468413952, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"demosite/templates/default.html": demositeTemplatesDefaultHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"demosite": &bintree{nil, map[string]*bintree{
		"templates": &bintree{nil, map[string]*bintree{
			"default.html": &bintree{demositeTemplatesDefaultHtml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
