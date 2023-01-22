// site/site_test.go
// -----------------
//
// TODO: make all the config tests etc. NOT use the filesystem.
// (arguably: make New not check the FS in the first place, but put
// the config loader into Load)
package site_test

import (
	// Standard:
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/site"
)

// TODO: build this into testify.assert, send a pull request?
// TODO: ...or at least put it in our own library, say "testy" or something.
// TODO: ...consider a "testy" library anyway, since I'm less and less
//       thrilled with assert?
// (Because hey, it really sucks to not test the *nature* of the panic.)
// (And because hey, I'm kinda testy...)
func AssertPanicsWith(t *testing.T, f func(), exp, msg string) {

	panicked := false
	got := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				got = fmt.Sprintf("%s", r)
			}
		}()
		f()
	}()

	if !panicked {
		assert.Fail(t, "Function did not panic.", msg)
		t.FailNow()
	} else if got != exp {
		errMsg := fmt.Sprintf(
			"Panic not as expected:\n  expected: %s\n    actual: %s",
			exp, got)
		assert.Fail(t, errMsg, msg)
	}
}

func Test_New_WithoutPath(t *testing.T) {

	assert := assert.New(t)

	s, err := site.New("")
	assert.Nil(err, "no error returned")
	assert.NotNil(s, "Site returned")
	assert.NotNil(s.Config, "Config set")
}

func Test_New_PathNotFound(t *testing.T) {

	assert := assert.New(t)

	_, err := site.New("noSuchPath")
	if assert.Error(err, "error returned") {
		assert.Equal("stat noSuchPath: no such file or directory",
			err.Error(), "error as expected")
	}
}

func Test_New_PathNotDir(t *testing.T) {

	assert := assert.New(t)

	_, err := site.New("site.go")
	if assert.Error(err, "error returned") {
		assert.Equal("site.go: not a directory",
			err.Error(), "error as expected")
	}
}

func Test_New_NoConfigFile(t *testing.T) {

	assert := assert.New(t)

	// We allow this, mostly in order to allow the app to create a default
	// server for a collection of pages that have nothing else.
	// (This may be a spurious use-case but it seems one might have a set of
	// .md files and want to just run something like "kisipar -p mydir"...)
	s, err := site.New("./")
	assert.Nil(err, "no error returned")
	assert.NotNil(s, "site returned")
	assert.NotNil(s.Config, "empty config set")

}

func Test_New_YamlNotFile(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	if err := os.Mkdir(cpath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	_, err := site.New(outer)
	if assert.Error(err, "error returned") {
		assert.Regexp("Config error.*directory",
			err.Error(), "error as expected")
	}
}

func Test_New_YamlConfig(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	cdata := []byte(`# test config
Name: test
Port: 1234
Host: fancyhost.com
BaseURL: any-base-works
`)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(outer)
	assert.Nil(err, "no error returned")
	if assert.NotNil(s, "site returned") {
		assert.Equal("test", s.Name, "name sticks")
		assert.Equal(1234, s.Port, "port sticks")
		assert.Equal("fancyhost.com", s.Host, "host sticks")
		assert.Equal("any-base-works", s.BaseURL, "baseurl sticks")
	}
}

func Test_New_YamlConfig_AllDefaults(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	cdata := []byte("# test config (all defaults)")
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(outer)
	assert.Nil(err, "no error returned")
	if assert.NotNil(s, "site returned") {
		assert.Equal("Anonymous Kisipar Site", s.Name, "name defaults")
		assert.Equal(8020, s.Port, "port defaults")
		assert.Equal("localhost", s.Host, "host defaults")
		assert.Equal("http://localhost:8020", s.BaseURL, "baseurl derived")
	}
}

func Test_New_YamlConfig_BaseURLFromHost(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	cdata := []byte(`# test config
Name: test
Port: 1234
Host: fancyhost.com
`)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(outer)
	assert.Nil(err, "no error returned")
	if assert.NotNil(s, "site returned") {
		assert.Equal("test", s.Name, "name sticks")
		assert.Equal(1234, s.Port, "port sticks")
		assert.Equal("fancyhost.com", s.Host, "host sticks")
		assert.Equal("http://fancyhost.com", s.BaseURL,
			"baseurl derived correctly")
	}
}

func Test_New_YamlConfig_BaseURLFromPort(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	cdata := []byte(`# test config
Name: test
Port: 1234
`)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(outer)
	assert.Nil(err, "no error returned")
	if assert.NotNil(s, "site returned") {
		assert.Equal("test", s.Name, "name sticks")
		assert.Equal(1234, s.Port, "port sticks")
		assert.Equal("localhost", s.Host, "host defaults")
		assert.Equal("http://localhost:1234", s.BaseURL,
			"baseurl derived correctly")
	}
}

func Test_New_YamlConfig_TLS(t *testing.T) {

	assert := assert.New(t)

	outer, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(outer)
	cpath := filepath.Join(outer, "config.yaml")
	cdata := []byte(`# test config
Name: test
Host: foo.com
KeyFile: /key/not/checked/yet
CertFile: /cert/not/checked/yet
`)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(outer)
	assert.Nil(err, "no error returned")
	if assert.NotNil(s, "site returned") {
		assert.Equal("/key/not/checked/yet", s.KeyFile, "keyfile sticks")
		assert.Equal("/cert/not/checked/yet", s.CertFile, "certfile sticks")
		assert.True(s.ServeTLS, "expect to serve TLS")
		assert.Equal("https://foo.com", s.BaseURL,
			"baseurl derived with https")
	}
}

func Test_New_BadPageExtensions(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	cpath := filepath.Join(dir, "config.yaml")
	cdata := []byte(`# good yaml, bad list of strings:
Name: Bad Page Extensions
PageExtensions: [foo,1,3.4,x]
`)
	if err := ioutil.WriteFile(cpath, cdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	_, err := site.New(dir)
	if assert.Error(err, "error returned") {
		assert.Equal("Config PageExtensions must be a list of strings; "+
			"item 1 is int.", err.Error(), "error useful")
	}
}

func Test_URL(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{}

	// Nothing set:
	assert.Equal("/foo/bar", s.URL("/foo/bar"),
		"empty Site defaults with no BaseURL")

	// Normal case:
	s.BaseURL = "all-your-base"
	assert.Equal("all-your-base/foo/bar", s.URL("/foo/bar"),
		"standard case works with leading slash")
	assert.Equal("all-your-base/foo/bar", s.URL("foo/bar"),
		"standard case works without leading slash")
	assert.Equal("all-your-base/foo/bar", s.URL("foo/bar/"),
		"standard case works with trailing slash")

}

func Test_Href_NilPage(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{}

	assert.Equal("", s.Href(nil), "Href for nil is empty string")

}

func Test_Href(t *testing.T) {

	assert := assert.New(t)

	type testCase struct {
		s   string
		p   string
		exp string
	}
	testCases := []testCase{
		{"", "/foo/bar.md", "/foo/bar"},
		{"/foo", "/foo/bar.md", "/bar"},
		{"/foo/bar", "/foo/bar/index.md", "/"},
		{"", "/foo/index.md", "/foo"},
		{"/foo", "/foo/index.md", "/"},
	}

	for _, tc := range testCases {
		s := &site.Site{PagePath: tc.s}
		p, err := page.LoadVirtualString(tc.p, "# ANY CONTENT")
		if err != nil {
			t.Fatal(err)
		}
		res := s.Href(p)
		assert.Equal(tc.exp, res,
			"%s -> %s for s.PagePath = '%s'", tc.p, tc.exp, tc.s)

	}

}

func Test_PageURL(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{BaseURL: "https://urbase", PagePath: "/fakepages"}

	assert.Equal("", s.PageURL(nil), "empty string for nil page")

	p1, err := page.LoadVirtualString("/fakepages/foo.md", "# ANYTHING")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("https://urbase/foo", s.PageURL(p1),
		"expected url for page under PagePath")

	p2, err := page.LoadVirtualString("/any/foo.md", "# ANYTHING")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("https://urbase/any/foo", s.PageURL(p2),
		"expected url for page not under PagePath")

	p3, err := page.LoadVirtualString("/any/foo/index.md", "# ANYTHING")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("https://urbase/any/foo", s.PageURL(p3),
		"expected url for page not under PagePath, but is index")

	// Edge case: what if our path has no leading slash?
	p4, err := page.LoadVirtualString("any/foo/bar.md", "# ANYTHING")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("https://urbase/any/foo/bar", s.PageURL(p4),
		"expected url for page not under PagePath, sans leading slash")

	// Edge case: what if your BaseURL has a trailing slash?
	p5, err := page.LoadVirtualString("/any/foo/bar.md", "# ANYTHING")
	if err != nil {
		t.Fatal(err)
	}
	s.BaseURL = "foo://bar/boo/"
	assert.Equal("foo://bar/boo/any/foo/bar", s.PageURL(p5),
		"expected url for page not under PagePath, BaseURL ends with slash")
}
