package kisipar_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/biztos/kisipar"
	"github.com/stretchr/testify/assert"
)

func tmpSite(config string) string {
	if config == "" {
		config = `# tmpSite config.yaml
Name: tmpSite
Port: 8086
Host: tmpsite.org
`
	}

	dir, err := ioutil.TempDir("", "kisipar-test-")
	if err != nil {
		panic(err)
	}
	cpath := filepath.Join(dir, "config.yaml")
	cdata := []byte(config)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		panic(err)
	}

	// subdirs:
	if err := os.Mkdir(filepath.Join(dir, "templates"), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "pages"), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "static"), os.ModePerm); err != nil {
		panic(err)
	}

	return dir

}

type FakeServer struct {
	ServeTLS bool
	Name     string
}

func (s *FakeServer) ListenAndServe() error {
	return fmt.Errorf("INSECURE %s", s.Name)
}
func (s *FakeServer) ListenAndServeTLS(c, k string) error {
	return fmt.Errorf("SECURE: %s & %s", c, k)
}

func Test_Load_RequiresPaths(t *testing.T) {
	assert := assert.New(t)
	_, err := kisipar.Load()
	if assert.Error(err, "error returned") {
		assert.Equal("kisipar.Load requires at least one site path.",
			err.Error(), "error as expected")
	}
}

func Test_Load_MissingTemplates(t *testing.T) {
	assert := assert.New(t)

	path := tmpSite("TemplatePath: NonesuchDir")
	defer os.RemoveAll(path)

	_, err := kisipar.Load(path)
	if assert.Error(err, "error returned") {
		assert.Regexp("TemplatePath not found",
			err.Error(), "error as expected")
	}
}

func Test_Load_SingleSite(t *testing.T) {

	assert := assert.New(t)

	path := tmpSite("")
	defer os.RemoveAll(path)

	k, err := kisipar.Load(path)
	if assert.Nil(err, "no error returned") {
		if assert.Equal(len(k.Sites), 1, "one site in Sites") {
			assert.Equal("tmpSite", k.Sites[0].Name, "site config loaded")
		}
	}
}

func Test_Load_MultiSite(t *testing.T) {

	assert := assert.New(t)

	path1 := tmpSite("Name: tmpSite1\nPort: 1000")
	defer os.RemoveAll(path1)
	path2 := tmpSite("Name: tmpSite2\nPort: 2000")
	defer os.RemoveAll(path2)

	k, err := kisipar.Load(path1, path2)
	if assert.Nil(err, "no error returned") {
		if assert.Equal(len(k.Sites), 2, "two sites in Sites") {
			assert.Equal("tmpSite1", k.Sites[0].Name, "site config 1 loaded")
			assert.Equal("tmpSite2", k.Sites[1].Name, "site config 2 loaded")
		}
	}
}

func Test_Load_ErrorDuplicatePort(t *testing.T) {

	assert := assert.New(t)

	path1 := tmpSite("Name: tmpSite1\nPort: 1000")
	defer os.RemoveAll(path1)
	path2 := tmpSite("Name: tmpSite2\nPort: 2000")
	defer os.RemoveAll(path2)
	path3 := tmpSite("Name: tmpSite3\nPort: 1000")
	defer os.RemoveAll(path3)

	_, err := kisipar.Load(path1, path2, path3)
	if assert.Error(err, "error returned") {
		assert.Equal("Duplicate Port 1000: "+path1+" vs. "+path3,
			err.Error(), "error as expected")
	}
}

// Yeah, something is nondeterministic here, should see if maybe that
// last buffer is not being flushed, or a goroutine is dropped maybe?
func xxx_Test_Serve_MultiSite(t *testing.T) {

	assert := assert.New(t)

	path1 := tmpSite("Name: tmpSite1\nPort: 8088")
	defer os.RemoveAll(path1)
	path2 := tmpSite("Name: tmpSite2\nPort: 8099")
	defer os.RemoveAll(path2)

	k, err := kisipar.Load(path1, path2)
	if err != nil {
		t.Fatal(err)
	}

	// Now fake out the Server so we don't actually launch one.
	// (The things we do for unit testing...)
	k.Sites[0].Server = &FakeServer{Name: "FIRST"}
	k.Sites[1].Server = &FakeServer{Name: "LAST"}

	// ...now serve, catching output.
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetFlags(0)
	k.Serve()
	exp := `tmpSite1: listening on port 8088.
tmpSite2: listening on port 8099.
tmpSite2 (port 8099): INSECURE LAST
`

	assert.Equal(exp, buf.String(), "servers logged as expected")

}
