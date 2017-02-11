// crutch.go - test coverage crutch for poorly designed shit like bindata.
// ---------
// OK, if you hate it so much why not write your own?
// (Maybe I should.)
//
// FUCKING BINDATA!
// There are enough unreachable test cases that hey, fuckit.  Write our own.

package kisipar

import (
	"errors"
)

// CrutchUnderBindata calls untestable code in bindata.
func CrutchUnderBindata() {

	// First up: we need to trigger an error here:
	// func bindataRead(data []byte, name string) ([]byte, error) {
	//             gz, err := gzip.NewReader(bytes.NewBuffer(data))
	//             if err != nil {
	//                     return nil, fmt.Errorf("Read %q: %v", name, err)
	//             }
	bindataReadBadData := []byte("probably not gzipped")
	_, err := bindataRead(bindataReadBadData, "foo")
	if err == nil {
		panic("no error for bindataReadBadData")
	}

	// Next, from the same func:
	// var buf bytes.Buffer
	//   _, err = io.Copy(&buf, gz)
	//   clErr := gz.Close()
	//
	//   if err != nil {
	//           return nil, fmt.Errorf("Read %q: %v", name, err)
	//   }
	//   if clErr != nil {
	//           return nil, err
	//   }
	// Crap, is this even possible?
	// io.Copy wants to copy gz into buf until an error occurs or
	// EOF (of gz) is reached; gz is already valid at this point.
	// Maybe if the gzip footer is wrong?  Indeed.  Fucking hell guys.
	bindataReadBadChecksum := []byte{
		0x1f, 0x8b, 0x8, 0x8, 0x1a, 0xb3, 0x90, 0x58, 0x0, 0x3, 0x74, 0x69,
		0x6e, 0x79, 0x2e, 0x6a, 0x73, 0x6f, 0x6e, 0x0, 0xab, 0x56, 0xca,
		0x48, 0xcd, 0xc9, 0xc9, 0x57, 0xb2, 0x52, 0xa, 0xf7, 0xf, 0xf2, 0x71,
		0x51, 0xaa, 0xe5, 0x2, 0x0,
		// Footer, 8 bytes:
		0x56, 0x33, 0x4f, 0x22, 0x33, 0x0, 0x0, 0x0,
	}
	_, err = bindataRead(bindataReadBadChecksum, "foo")
	if err == nil {
		panic("no error for bindataReadBadChecksum")
	}
	// Lovely.  Now how do we make it fail to Close?
	// It appears we can not:
	// * we already have a valid gzip.Reader
	// * io.Copy will read its data until EOF
	// * it already did that without gzip complaining
	// * Close will only fail if the footer is screwed up
	//
	// TODO: if this really can't be beat then autogen it away.

	// Now back to something easier to cheat:
	// func templatesDefaultHtml() (*asset, error) {
	//         bytes, err := templatesDefaultHtmlBytes()
	//         if err != nil {
	//                 return nil, err
	//         }
	// The data is stored in a regular variable so...
	old_templatesDefaultHtml := _templatesDefaultHtml
	defer func() { _templatesDefaultHtml = old_templatesDefaultHtml }()

	// ARGH WE WILL HAVE ONE OF THESE FUCKING THINGS FOR EVERY FUCKING
	// ASSET!!! Generator needed obviously.  The crutch func will go there.
	_templatesDefaultHtml = []byte("not quite gzip")
	_, err = templatesDefaultHtmlBytes()
	if err == nil {
		panic("no error for templatesDefaultHtmlBytes")
	}
	_, err = templatesDefaultHtml()
	if err == nil {
		panic("no error for templatesDefaultHtml")
	}

	// Ditto here, generation needed: function errors.
	retErr := func() (*asset, error) { return nil, errors.New("ERR HERE") }
	old_bindata := _bindata
	defer func() { _bindata = old_bindata }()
	_bindata = map[string]func() (*asset, error){
		"templates/default.html": retErr,
	}
	_, err = Asset("templates/default.html")
	if err == nil {
		panic("no error for Asset")
	}
	_, err = AssetInfo("templates/default.html")
	if err == nil {
		panic("no error for AssetInfo")
	}
	err = RestoreAsset("anydir", "templates/default.html")
	if err == nil {
		panic("no error for RestoreAsset for bad asset func")
	}

}
