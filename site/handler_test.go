// handler_test.go -- tests for the kisipar multiplexing http handler (mux).
// ---------------

package site_test

import (
	// Standard:
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/site"
)

type FailResponseWriter struct{}

func (w *FailResponseWriter) Write(b []byte) (int, error) { return 0, errors.New("FAIL") }
func (w *FailResponseWriter) WriteHeader(code int)        { return }
func (w *FailResponseWriter) Header() http.Header         { return http.Header{} }

func ReqAndRec(t *testing.T, url string) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	return req, rec

}

type FakeHandler struct {
	served int
}

func (h *FakeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.served++
	fmt.Fprintf(w, "FAKE\n")
}

func Test_NewServeMux(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{}
	m := s.NewServeMux()
	assert.NotNil(m, "mux returned")

}

func Test_NewServeMux_Overrides(t *testing.T) {

	assert := assert.New(t)

	h1 := &site.PatternHandler{
		Pattern: "/here",
		Function: func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "HERE\n")
		},
	}
	h2 := &site.PatternHandler{
		Pattern: "/none",
		Handler: http.NotFoundHandler(),
	}
	h3 := &site.PatternHandler{
		Pattern: "/",
		Function: func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "FALLBACK\n")
		},
	}

	s, err := site.New("")
	if err != nil {
		t.Fatal(err)
	}
	m := s.NewServeMux(h1, h2, h3)
	assert.NotNil(m, "mux returned")
	s.SetHandler(m)

	// And let's prove it works, not that the coverage metric cares.
	type res struct {
		url  string
		code int
		body string
	}
	tests := []res{
		res{url: "http://example.com/here", code: 200, body: "HERE\n"},
		res{url: "http://example.com/none", code: 404,
			body: "404 page not found\n"},
		res{url: "http://example.com/other", code: 200, body: "FALLBACK\n"},
	}
	for _, test := range tests {
		url := test.url
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		assert.Equal(test.code, w.Code, "code as expected for "+url)
		assert.Equal(test.body, w.Body.String(), "body as expected for "+url)
	}

}

func Test_NewServeMux_PanicDuplicatePatterns(t *testing.T) {

	h1 := &site.PatternHandler{
		Pattern:  "/this",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	h2 := &site.PatternHandler{
		Pattern:  "/that",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	h3 := &site.PatternHandler{
		Pattern:  "/this",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	s := &site.Site{}

	AssertPanicsWith(t, func() { s.NewServeMux(h1, h2, h3) },
		"Duplicate handler pattern: /this", "NewServeMux panics as expected")

}

func Test_NewServeMux_PanicEmptyPattern(t *testing.T) {

	h1 := &site.PatternHandler{
		Pattern:  "/this",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	h2 := &site.PatternHandler{
		Pattern:  "/that",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	h3 := &site.PatternHandler{
		Pattern:  "",
		Function: func(w http.ResponseWriter, req *http.Request) {},
	}
	s := &site.Site{}
	AssertPanicsWith(t, func() { s.NewServeMux(h1, h2, h3) },
		"Empty handler pattern at position 2.",
		"NewServeMux panics as expected")
}

func Test_NewServeMux_PanicNils(t *testing.T) {

	h := &site.PatternHandler{Pattern: "/nope"}
	s := &site.Site{}
	AssertPanicsWith(t, func() { s.NewServeMux(h) },
		"Neither handler nor function for pattern: /nope",
		"NewServeMux panics as expected")
}

func Test_NewServeMux_PanicRedundant(t *testing.T) {

	h := &site.PatternHandler{
		Pattern:  "/toomuch",
		Function: func(w http.ResponseWriter, req *http.Request) {},
		Handler:  http.NotFoundHandler(),
	}

	s := &site.Site{}
	AssertPanicsWith(t, func() { s.NewServeMux(h) },
		"Redundant handler and function for pattern: /toomuch",
		"NewServeMux panics as expected")

}

func Test_NewServeMux_SimpleSuccess(t *testing.T) {

	assert := assert.New(t)

	s, err := site.LoadVirtualYaml("") // all defaults, no TLS
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(s.Server, "Server set")

	svr, ok := s.Server.(*http.Server)
	assert.True(ok, "Server is an http.Server")
	assert.NotNil(svr.Handler, "Handler set in Server")
	assert.IsType(&http.ServeMux{}, svr.Handler, "It's a ServeMux!")

	// Prove it works, minimally.
	type res struct {
		url  string
		code int
	}
	tests := []res{
		res{url: "http://example.com/", code: 200},
		res{url: "http://example.com/none", code: 404},
	}
	for _, test := range tests {
		url := test.url
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		if !assert.Equal(test.code, w.Code, "code as expected for "+url) {
			t.Logf(w.Body.String())
		}

	}

}

func Test_SetHandler(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{
		Server: &FakeServer{},
	}

	AssertPanicsWith(t, func() { s.SetHandler(&FakeHandler{}) },
		"Server is not an http.Server.",
		"SetHandler with unexpected Server type panics")

	s.Server = &http.Server{}
	assert.NotPanics(func() { s.SetHandler(&FakeHandler{}) },
		"no panic with an http.Server")

}

func Test_ServeHTTP(t *testing.T) {

	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/any", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	s := &site.Site{
		Server: &FakeServer{},
	}

	AssertPanicsWith(t, func() { s.ServeHTTP(w, req) },
		"Server is not an http.Server.",
		"ServeHTTP with unexpected Server type panics")

	s.Server = &http.Server{Handler: &FakeHandler{}}
	assert.NotPanics(func() { s.ServeHTTP(w, req) },
		"no panic with an http.Server")

	assert.Equal(200, w.Code, "FakeHandler sends 200 OK")
	assert.Equal("FAKE\n", w.Body.String(), "FakeHandler sends body")

}

func Test_MainHandler(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{
		Server: &FakeServer{},
	}
	h := s.MainHandler()
	assert.NotNil(h, "generic MainHandler not nil")
}

func Test_MainHandler_StaticFile(t *testing.T) {

	assert := assert.New(t)

	// First we test that it *does* load the file, then we test that it
	// does not if there is no StaticPath when the function is created.
	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Standard file found:
	req, w := ReqAndRec(t, "http://example.com/js/kisipar.js")
	handler(w, req)
	assert.Equal(200, w.Code, "200 response recorded")
	assert.Equal("var kisipar = {\n    hello: \"szavassz világ!\"\n}",
		w.Body.String(), "static file written")

	// Standard file not found:
	req, w = ReqAndRec(t, "http://example.com/js/nonesuch.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded for missing file")

	// Now make another with no static files.
	s.StaticPath = ""
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/js/kisipar.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")

	// No such dir:
	s.StaticPath = filepath.Join(path, "no_such_thing")
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/js/kisipar.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")

	// Not a dir:
	s.StaticPath = filepath.Join(path, "config.yaml")
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/js/kisipar.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")
	t.Logf("Body: %s", w.Body.String())

}

func Test_MainHandler_StaticDir(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	req, w := ReqAndRec(t, "http://example.com/js")
	handler(w, req)
	assert.Equal(404, w.Code,
		"404 response recorded for plain static dir w/o extension")

	// again with an extension
	req, w = ReqAndRec(t, "http://example.com/js/special.dir")
	handler(w, req)
	assert.Equal(404, w.Code,
		"404 response recorded for static dir with extension")

	// TODO: some basic content check on the default (pretty) error page.
	// TODO: prove it was NOT logged -- and maybe then prove it WAS, dep.
	// on settings.
}

// The PathErrors from os.Stat are no longer testable since they are all 404.
// TODO: come up with another way to test a file-open error.
func TODO_Test_MainHandler_StaticFileError(t *testing.T) {

	// This one is rather convoluted because the only obvious way to trigger
	// a file error is for its *parent* folder to be unreadable.  Or, maybe
	// there's a way better approach I am just missing.
	assert := assert.New(t)

	// Make a basic site and initialize it:
	path, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	spath := filepath.Join(path, "static")
	if err := os.Mkdir(spath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	fpath := filepath.Join(spath, "noperm.txt")
	b := []byte("OOPS READABLE")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Remove our ability to see the file, and then try to serve it.
	var modeNone os.FileMode = 0000
	if err := os.Chmod(spath, modeNone); err != nil {
		t.Fatal(err)
	}

	req, w := ReqAndRec(t, "http://example.com/noperm.txt")
	handler(w, req)

	assert.Equal(500, w.Code, "500 response recorded for bad permissions")

	// TODO: some basic content check on the default (pretty) error page.
	// TODO: prove the error was logged

}

func Test_MainHandler_StaticFileNotReadable(t *testing.T) {

	// Problem: if we get as far as http.ServeFile() with an unreadable file,
	// we will get a generic 403.  That leaks site info (i.e. which file is
	// static) and also (I *think*) forces us to use the generic error
	// response content.
	assert := assert.New(t)

	// Make a basic site and initialize it:
	path, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	spath := filepath.Join(path, "static")
	if err := os.Mkdir(spath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	fpath := filepath.Join(spath, "noperm.txt")
	b := []byte("OOPS READABLE")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Remove our ability to read the file, and then try to serve it.
	var modeNone os.FileMode = 0000
	if err := os.Chmod(fpath, modeNone); err != nil {
		t.Fatal(err)
	}

	req, w := ReqAndRec(t, "http://example.com/noperm.txt")
	handler(w, req)

	assert.Equal(500, w.Code, "500 response recorded for bad permissions")
	t.Logf("Body: %s", w.Body.String())
	// TODO: some basic content check on the default (pretty) error page.
	// TODO: prove the error was logged

}

func Test_MainHandler_AssetFile(t *testing.T) {

	assert := assert.New(t)

	// First we test that it *does* load the file, then we test that it
	// does not if there is no PagePath when the function is created.
	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Standard file found:
	req, w := ReqAndRec(t, "http://example.com/foo/boogie.js")
	handler(w, req)
	assert.Equal(200, w.Code, "200 response recorded for standard asset")
	assert.Equal("window.alert('BOOGIE!');\n",
		w.Body.String(), "page asset written")

	// Standard file not found:
	req, w = ReqAndRec(t, "http://example.com/foo/nonesuch.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded for missing file")

	// Source file not found:
	s.ServePageSources = false // as it should already be...
	req, w = ReqAndRec(t, "http://example.com/foo/bar.md")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded for source file")

	// Source file found (allowed):
	exp := "# Foo Bar Page\n\n    Author: Jimmini Criquette\n\nI am the Bar of the Foo!"
	s.ServePageSources = true
	req, w = ReqAndRec(t, "http://example.com/foo/bar.md")
	handler(w, req)
	assert.Equal(200, w.Code, "200 response recorded for source file")
	assert.Equal(exp, w.Body.String(), "page asset written")

	// Now make another with no PagePath so nothing will be looked up.
	s.PagePath = ""
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/foo/boogie.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")

	// No dir:
	s.PagePath = filepath.Join(path, "no_such_thing")
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/foo/boogie.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")

	// Not a dir:
	s.PagePath = filepath.Join(path, "config.yaml")
	handler = s.MainHandler()
	req, w = ReqAndRec(t, "http://example.com/foo/boogie.js")
	handler(w, req)
	assert.Equal(404, w.Code, "404 response recorded")
	t.Logf("Body: %s", w.Body.String())

}

func Test_MainHandler_AssetDir(t *testing.T) {

	assert := assert.New(t)

	// TODO: is there a way to do this without the friggin' file system?
	// (Probably there is not, since we need to look there for assets...)
	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	req, w := ReqAndRec(t, "http://example.com/foo/assets")
	handler(w, req)
	assert.Equal(404, w.Code,
		"asset-only dir (no pages) is NOT FOUND")

	// And another with an extension.
	req, w = ReqAndRec(t, "http://example.com/foo/stuff.dir")
	handler(w, req)
	assert.Equal(404, w.Code,
		"404 response recorded for unindexed pages dir with extension")

	// TODO: some basic content check on the default (pretty) error page.
	// TODO: prove it was NOT logged -- and maybe then prove it WAS, dep.
	// on settings.
	// TODO: at least consider allowing dotted dirs.  Is it a big deal?
}

func Test_MainHandler_AssetFileNotReadable(t *testing.T) {

	// Problem: if we get as far as http.ServeFile() with an unreadable file,
	// we will get a generic 403.  That leaks site info (i.e. which file is
	// pages) and also (I *think*) forces us to use the generic error
	// response content.
	assert := assert.New(t)

	// Make a basic site and initialize it:
	path, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	spath := filepath.Join(path, "pages")
	if err := os.Mkdir(spath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	fpath := filepath.Join(spath, "noperm.js")
	b := []byte("OOPS READABLE")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Remove our ability to read the file, and then try to serve it.
	var modeNone os.FileMode = 0000
	if err := os.Chmod(fpath, modeNone); err != nil {
		t.Fatal(err)
	}

	req, w := ReqAndRec(t, "http://example.com/noperm.js")
	handler(w, req)

	assert.Equal(500, w.Code, "500 response recorded for bad permissions")
	t.Logf("Body: %s", w.Body.String())
	// TODO: some basic content check on the default (pretty) error page.
	// TODO: prove the error was logged

}

func Test_MainHandler_WriterWriteFailure(t *testing.T) {

	// Sometimes the most twisted edge cases are the easiest to test...

	// Our standard test site on disk:
	// TODO: at least consider having a "create virtual test site" func here
	// so we don't rely on the filesystem for stuff that could just as well
	// be done by virtual sites.
	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	req, _ := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	w := &FailResponseWriter{}

	AssertPanicsWith(t, func() { handler(w, req) },
		"Write failed for /foo/bar: FAIL",
		"ResponseWriter Write failure panics")
}

func Test_MainHandler_StandardPage(t *testing.T) {

	assert := assert.New(t)

	// Our standard test site on disk:
	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	req, w := ReqAndRec(t, "http://example.com/foo/bar")
	handler(w, req)

	exp := `<!doctype html>
<html>
    <head>
        <title>Foo Bar Page - Testing Kisipar</title>
    </head>
    <body>
        <h1>Foo Bar Page - Testing Kisipar</h1>
<h2>I am the foo bar!</h2>
<div id="Content">
<h1>Foo Bar Page</h1>

<p>I am the Bar of the Foo!</p>

</div>

        <div id="Footer">
            Copyright &copy; 2016
            Jimmini Criquette
            
        </div>
    </body>
</html>
`
	assert.Equal(200, w.Code, "code 200 sent")
	assert.Equal(exp, w.Body.String(), "body as expected")

}

func Test_MainHandler_StandardPageParseFailure(t *testing.T) {

	assert := assert.New(t)

	// NOTE: it would be nice to not deal with any of the disk stuff, but
	// then it becomes really hard to test, say, a page-refresh failure since
	// that can't happen to virtual pages.
	// TODO: think about this more, there must be a smarter way.

	// Create a temp site with one good page, which will turn bad.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(pdir, "page.md")
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	// Rewrite garbage to the page, and set the Page in memory to be "old."
	bad := []byte("# Bad!\n\n    foo: { x [ y,\n\\nBar!")
	if err := ioutil.WriteFile(fpath, bad, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s.Pageset.ByPath()[0].ModTime = time.Unix(1, 1)

	req, w := ReqAndRec(t, "http://example.com/page")
	handler(w, req)

	// TODO: real(-ish) error pages!
	exp := "ERROR: Page error for /page: yaml: did not find expected ',' or '}'\n"

	assert.Equal(500, w.Code, "code 500 sent")
	assert.Equal(exp, w.Body.String(), "body as expected")

}

func Test_MainHandler_TemplateExecuteFailure(t *testing.T) {

	assert := assert.New(t)

	// NOTE: although it would make this particular case harder to test,
	// it would probably be wise to execute all templates on load, so you
	// don't end up with a bogus template in some obscure part of the site.

	// Create a temp site with one good page, and a bad template
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(pdir, "page.md")
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tpath := filepath.Join(tdir, "page.html")
	if err := ioutil.WriteFile(tpath, []byte("{{ .Bad }}"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.MainHandler()

	req, w := ReqAndRec(t, "http://example.com/page")
	handler(w, req)

	// TODO: real(-ish) error pages! And logging!
	assert.Equal(500, w.Code, "code 500 sent")
	assert.Regexp("^ERROR: template: page", w.Body.String(), "body as expected")

}

func Test_MainHandler_ServesMainIndexWithPagesetForIndexPage(t *testing.T) {

	assert := assert.New(t)

	// REMEMBER: leading slashes for virtual pages!
	yaml := `Name: Test
Pages:
    /foo/bar.md: "# Foo Bar"
    /foo/index.md: "# Foo Index"
Templates:
    index: |
        INDEX
        {{with .Pageset}}HAVE PAGESET{{end}}
        {{with .Page}}HAVE PAGE{{end}}
    single: |
        SINGLE
        {{with .Pageset}}HAVE PAGESET{{end}}
        {{with .Page}}HAVE PAGE{{end}}
`
	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}

	handler := s.MainHandler()

	req, w := ReqAndRec(t, "http://example.com/foo")
	handler(w, req)

	assert.Equal(200, w.Code, "code 200 sent")
	assert.Equal("INDEX\nHAVE PAGESET\nHAVE PAGE\n", w.Body.String(),
		"body as expected")

	for _, p := range s.Pageset.ByPath() {
		t.Log(p.Path)
	}
}

func Test_PageForPath_NotFoundWithoutPageset(t *testing.T) {

	assert := assert.New(t)

	s, err := site.LoadVirtualYaml("Name: Test")
	if err != nil {
		t.Fatal(err)
	}

	s.Pageset = nil

	p, err := s.PageForPath("any")
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "...of the IsNotExist persuasion")
	}
	assert.Nil(p, "no page returned")

}

func Test_PageForPath_ExactMatch_Fresh(t *testing.T) {

	// This is the most common case of a "found" page: an exact match,
	// and no change on disk since we loaded the site.

	assert := assert.New(t)

	// Create a minimal site on disk.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(pdir, "page.md")
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	p, err := s.PageForPath("/page")
	assert.Nil(err, "no error")
	if assert.NotNil(p, "page returned") {
		assert.Equal("Here!", p.Title(), "page parsed correctly")
	}
}

func Test_PageForPath_ExactMatch_Virtual(t *testing.T) {

	assert := assert.New(t)

	p, err := page.LoadVirtual("/foo/bar/baz.md", []byte("# Here!\n\ni am."))
	if err != nil {
		t.Fatal(err)
	}
	s, err := site.LoadVirtual(nil, []*page.Page{p}, nil)
	if err != nil {
		t.Fatal(err)
	}

	pfp, err := s.PageForPath("/foo/bar/baz")
	assert.Nil(err, "no error")
	if assert.NotNil(pfp, "page returned") {
		assert.Equal("Here!", pfp.Title(), "page parsed correctly")
	}

}

func Test_PageForPath_ExactMatch_Disappeared(t *testing.T) {

	assert := assert.New(t)

	// Create a minimal site on disk.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(pdir, "page.md")
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	// ...and kill it:
	if err := os.Remove(fpath); err != nil {
		t.Fatal(err)
	}

	p, err := s.PageForPath("/page")
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "error is IsNotExist-y")
	}
	assert.Nil(p, "no page returned")
}

func Test_PageForPath_ExactMatch_ParseError(t *testing.T) {

	assert := assert.New(t)

	// Create a minimal site on disk.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	key := filepath.Join(pdir, "page")
	fpath := key + ".md"
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	// ...make it anew, but badly:
	b = []byte("# Here!\n\n    bad: { [ (xxx...\n\n")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// ...and deal with our super-speedy test execution:
	s.Pageset.Page(key).ModTime = time.Unix(0, 0)

	p, err := s.PageForPath("/page")
	if assert.Error(err, "error returned") {
		assert.False(os.IsNotExist(err), "error is NOT IsNotExist-y")
		assert.Regexp("yaml", err.Error(), "error is useful")
	}
	assert.Nil(p, "no page returned")
}

func Test_PageForPath_IndexMatch_Fresh(t *testing.T) {

	assert := assert.New(t)

	// Create a minimal site on disk.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	idir := filepath.Join(pdir, "foo")
	if err := os.Mkdir(idir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(idir, "index.md")
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	p, err := s.PageForPath("/foo")
	assert.Nil(err, "no error")
	if assert.NotNil(p, "page returned") {
		assert.Equal("Here!", p.Title(), "page parsed correctly")
	}
}

func Test_PageForPath_IndexMatch_Virtual(t *testing.T) {

	assert := assert.New(t)

	p, err := page.LoadVirtual("/foo/index.md", []byte("# Here!\n\ni am."))
	if err != nil {
		t.Fatal(err)
	}
	s, err := site.LoadVirtual(nil, []*page.Page{p}, nil)
	if err != nil {
		t.Fatal(err)
	}

	pfp, err := s.PageForPath("/foo")
	assert.Nil(err, "no error")
	if assert.NotNil(pfp, "page returned") {
		assert.Equal("Here!", pfp.Title(), "page parsed correctly")
	}
}

func Test_PageForPath_IndexMatch_ParseError(t *testing.T) {

	assert := assert.New(t)

	// Create a minimal site on disk.
	dir, err := ioutil.TempDir("", "kisipar-site-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	if err := os.Mkdir(pdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tdir := filepath.Join(dir, "templates")
	if err := os.Mkdir(tdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	idir := filepath.Join(pdir, "foo")
	if err := os.Mkdir(idir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	key := filepath.Join(idir, "index")
	fpath := key + ".md"
	b := []byte("# Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	// ...make it anew, but badly:
	b = []byte("# Here!\n\n    bad: { [ (xxx...\n\n")
	if err := ioutil.WriteFile(fpath, b, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// ...and deal with our super-speedy test execution:
	s.Pageset.Page(key).ModTime = time.Unix(0, 0)

	p, err := s.PageForPath("/foo")
	if assert.Error(err, "error returned") {
		assert.False(os.IsNotExist(err), "error is NOT IsNotExist-y")
		assert.Regexp("yaml", err.Error(), "error is useful")
	}
	assert.Nil(p, "no page returned")
}

func Test_PageForPath_NoMatch_Virtual(t *testing.T) {

	assert := assert.New(t)

	p, err := page.LoadVirtual("/foo/bar/baz.md", []byte("# Here!\n\ni am."))
	if err != nil {
		t.Fatal(err)
	}
	s, err := site.LoadVirtual(nil, []*page.Page{p}, nil)
	if err != nil {
		t.Fatal(err)
	}

	pfp, err := s.PageForPath("/page")
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "error is IsNotExist-y")
	}
	assert.Nil(pfp, "no page returned")

}
