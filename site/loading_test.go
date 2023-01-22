// loading_test.go - test the loading functions for the Kisipar site.
// ---------------

package site_test

import (
	// Standard:
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	// Third-party:
	"github.com/olebedev/config"
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/site"
)

func Test_LoadPages(t *testing.T) {

	assert := assert.New(t)

	exp_titles := []string{
		"Index Page",
		"Other Page",
		"Bar Index!",
		"Bar Abacus!",
		"Bar Boomerang!",
		"Foo Bar Page",
		"Foo Bat Page",
		"Foo Baz Page",
		"Stuff and Bother!",
	}
	titles := []string{}

	path := filepath.Join("test_data", "full_site")
	s, err := site.New(path)
	assert.Nil(err, "no error initializing site")

	err = s.LoadPages()
	assert.Nil(err, "no error on LoadPages")
	for _, p := range s.Pageset.ByPath() {
		titles = append(titles, p.Title())
	}

	assert.Equal(exp_titles, titles,
		"pages loaded as expected (by title check)")

}

func Test_LoadPages_ErrorNoSuchDir(t *testing.T) {

	assert := assert.New(t)

	s := &site.Site{PagePath: "no/such/thing/here"}
	err := s.LoadPages()
	if assert.Error(err, "error loading page with missing PagePath") {
		assert.True(os.IsNotExist(err), "error true for os.IsNotExist")
	}

}

func Test_LoadPages_ErrorBadPage(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	pdir := filepath.Join(dir, "pages")
	perr := os.Mkdir(pdir, os.ModePerm)
	if perr != nil {
		t.Fatal(perr)
	}
	fpath := filepath.Join(pdir, "foo.md")
	fdata := []byte("# xx\n\n```json\n{ foo: [}\n```\n\n")
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	s, err := site.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = s.LoadPages()
	if assert.Error(err, "error loading page with bad meta") {
		assert.Regexp("invalid character", err.Error(), "error makes sense")
	}

}

func Test_LoadPages_NoPagePath(t *testing.T) {

	assert := assert.New(t)

	// Basic thing: sets the Pageset
	s := &site.Site{}
	err := s.LoadPages()
	assert.Nil(err, "no error on LoadPages")
	if assert.NotNil(s.Pageset, "Pageset has content") {
		assert.Equal(0, s.Pageset.Len(), "Pageset has zero size")
	}
}

func Test_LoadPages_PagesetReplaced(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("test_data", "full_site")
	s, err := site.New(path)
	assert.Nil(err, "no error initializing site")

	s.PagePath = ""

	// ...we should consider testing more variations here, as this test
	// hits the coverage only because we know what's inside the box.
	err = s.LoadPages()
	assert.Nil(err, "no error on LoadPages")
	if assert.NotNil(s.Pageset, "Pageset has content") {
		assert.Equal(0, s.Pageset.Len(), "Pageset has zero size")
	}
}

func Test_LoadTemplates(t *testing.T) {

	assert := assert.New(t)

	exp_names := []string{
		"", // the nameless default
		"bar/index",
		"bar/single",
		"error",
		"foo/bar",
		"index",
		"shared/foot",
		"shared/head",
		"single",
	}
	names := []string{}

	path := filepath.Join("test_data", "full_site")
	s, err := site.New(path)
	assert.Nil(err, "no error initializing site")

	err = s.LoadTemplates()
	assert.Nil(err, "no error on LoadTemplates")
	for _, t := range s.Template.Templates() {
		names = append(names, t.Name())
	}
	sort.Strings(names)

	assert.Equal(exp_names, names,
		"templates loaded as expected (by name check)")
}

func Test_LoadTemplates_ErrorBadDefaultTemplate(t *testing.T) {

	assert := assert.New(t)

	old_default := site.DEFAULT_TEMPLATE
	defer func() { site.DEFAULT_TEMPLATE = old_default }()

	site.DEFAULT_TEMPLATE = `{{ no-good }}`

	path := filepath.Join("test_data", "full_site")
	s, err := site.New(path)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadTemplates()
	if assert.Error(err, "error returned for bad default template") {
		assert.Regexp("^DEFAULT_TEMPLATE.*bad character",
			err.Error(), "error is useful")
	}

}

func Test_LoadTemplates_ErrorDupeName(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	tdir := filepath.Join(dir, "templates")
	terr := os.Mkdir(tdir, os.ModePerm)
	if terr != nil {
		t.Fatal(terr)
	}
	tpath := filepath.Join(tdir, "foo.html")
	tdata := []byte("{{ .Site.Name }}")
	if err := ioutil.WriteFile(tpath, tdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tpath2 := filepath.Join(tdir, "foo.tmpl") // upperce isn't enough on OSX
	if err := ioutil.WriteFile(tpath2, tdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	s, err := site.New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadTemplates()
	if assert.Error(err, "error returned for dupe template") {
		assert.Regexp("^Duplicate template for foo",
			err.Error(), "error is useful")
	}

}

func Test_LoadTemplates_ErrorFileRead(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	tdir := filepath.Join(dir, "templates")
	terr := os.Mkdir(tdir, os.ModePerm)
	if terr != nil {
		t.Fatal(terr)
	}
	tpath := filepath.Join(tdir, "foo.html")
	tdata := []byte("{{ .Site.Name }}")
	var modeNope os.FileMode = 0000
	if err := ioutil.WriteFile(tpath, tdata, modeNope); err != nil {
		t.Fatal(err)
	}

	s, err := site.New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadTemplates()
	if assert.Error(err, "error returned for unreadable file") {
		// TODO: make this useful cross-platform
		assert.Regexp("permission",
			err.Error(), "error is useful")
	}

}

func Test_LoadTemplates_ErrorDirNotFound(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	s, err := site.New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadTemplates()
	if assert.Error(err, "error returned for no templates dir") {
		// TODO: make this useful cross-platform
		assert.Regexp("^TemplatePath not found",
			err.Error(), "error is useful")
	}

}

func Test_LoadTemplates_ErrorDirNotReadable(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	// Oh my, this is cumbersome.  We need to catch unexpected stat errors
	// and we want to beautify the IsNotExist case, but that leaves us having
	// to create another error for os.Stat.  Apparently an unreadable path
	// may do the trick (and unreadable file will not).
	var modeNope os.FileMode = 0000
	sdir := filepath.Join(dir, "unreadable")
	if serr := os.Mkdir(sdir, modeNope); serr != nil {
		t.Fatal(serr)
	}
	s := &site.Site{TemplatePath: filepath.Join(sdir, "foo")}

	err := s.LoadTemplates()
	if assert.Error(err, "error returned for unreadable templates") {
		// TODO: make this useful cross-platform
		assert.Regexp("permission",
			err.Error(), "error is useful")
	}

}

func Test_LoadTemplates_ErrorDirNotDir(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "templates")
	if err := ioutil.WriteFile(fpath, []byte("x"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	s, err := site.New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadTemplates()
	if assert.Error(err, "error returned for non-dir templates") {
		// TODO: make this useful cross-platform
		assert.Regexp("^Not a directory",
			err.Error(), "error is useful")
	}

}

func Test_LoadVirtual_ErrorBadPageExtensionsList(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()

	cfg, err := config.ParseYaml(`# good yaml, bad list of strings:
Name: Bad Page Extensions
PageExtensions: [foo,1,3.4,x]
`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = site.LoadVirtual(cfg, []*page.Page{}, nil)
	if assert.Error(err, "error returned") {
		assert.Equal("Config PageExtensions must be a list of strings; "+
			"item 1 is int.", err.Error(), "error useful")
	}

}

func Test_LoadVirtual_ErrorBadPageExtensionsType(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()

	cfg, err := config.ParseYaml(`# good yaml, bad list of strings:
Name: Bad Page Extensions
PageExtensions: "ima_notlist"
`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = site.LoadVirtual(cfg, []*page.Page{}, nil)
	if assert.Error(err, "error returned") {
		assert.Equal("Config PageExtensions is not a list.", err.Error(),
			"error useful")
	}

}

func Test_LoadVirtual_ErrorDupePage(t *testing.T) {

	assert := assert.New(t)

	cfg, _ := config.ParseYaml("")
	p1, _ := page.LoadVirtualString("foo.md", "# here")
	p2, _ := page.LoadVirtualString("foo.md", "# here")

	_, err := site.LoadVirtual(cfg, []*page.Page{p1, p2}, nil)
	if assert.Error(err, "error returned") {
		assert.Equal("Duplicate path for foo: foo.md", err.Error(),
			"error useful")
	}

}

func Test_LoadVirtual_ErrorBadDefaultTemplate(t *testing.T) {

	assert := assert.New(t)

	old_default := site.DEFAULT_TEMPLATE
	defer func() { site.DEFAULT_TEMPLATE = old_default }()

	site.DEFAULT_TEMPLATE = `{{ no-good }}`
	cfg, _ := config.ParseYaml("")
	p1, _ := page.LoadVirtualString("foo.md", "# here")
	p2, _ := page.LoadVirtualString("foo.md", "# here")

	_, err := site.LoadVirtual(cfg, []*page.Page{p1, p2}, nil)
	if assert.Error(err, "error returned") {
		assert.Regexp("^DEFAULT_TEMPLATE", err.Error(), "error useful")
	}

}

func Test_LoadVirtual_Success(t *testing.T) {

	assert := assert.New(t)

	// TODO: with template!
	cfg, _ := config.ParseYaml("")
	p1, _ := page.LoadVirtualString("foo.md", "# here")
	p2, _ := page.LoadVirtualString("bar.md", "# here")

	s, err := site.LoadVirtual(cfg, []*page.Page{p1, p2}, nil)
	assert.Nil(err, "no error returned")
	assert.NotNil(s, "site returned")

}

func Test_LoadVirtual_SuccessWithNils(t *testing.T) {

	assert := assert.New(t)

	s, err := site.LoadVirtual(nil, nil, nil)
	assert.Nil(err, "no error returned")
	assert.NotNil(s, "site returned")
	assert.NotNil(s.Config, "Config not nil in site")
	assert.NotNil(s.Pageset, "Pageset not nil in site")
	assert.NotNil(s.Template, "Template not nil in site")

}

func Test_LoadVirtual_SuccessWithPageExtensions(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()

	// TODO: with template!
	cfg, _ := config.ParseYaml(`PageExtensions: [.foo, .bar]`)
	p1, _ := page.LoadVirtualString("foo.md", "# here")
	p2, _ := page.LoadVirtualString("bar.md", "# here")

	s, err := site.LoadVirtual(cfg, []*page.Page{p1, p2}, nil)
	assert.Nil(err, "no error returned")
	assert.NotNil(s, "site returned")
	assert.Equal([]string{".foo", ".bar"}, s.PageExtensions,
		"PageExtensions stick")
}

func Test_Load(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("test_data", "full_site")
	s, err := site.Load(path)
	assert.Nil(err, "no error on load of "+path)
	assert.NotNil(s, "site returned")

}

func Test_Load_ErrorNoPath(t *testing.T) {

	assert := assert.New(t)

	_, err := site.Load("")
	if assert.Error(err, "error on Load with empty path") {
		assert.Equal("path must not be empty", err.Error(),
			"error makes sense")
	}
}

func Test_Load_ErrorBadConfig(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	fpath := filepath.Join(dir, "config.yaml")
	fdata := []byte("bad: [a-b,}\n")
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	_, err := site.Load(dir)
	if assert.Error(err, "error on Load with bad config") {
		assert.Regexp("^Config error.*yaml", err.Error(),
			"error makes sense")
	}
}

func Test_Load_ErrorBadTemplate(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	tdir := filepath.Join(dir, "templates")
	terr := os.Mkdir(tdir, os.ModePerm)
	if terr != nil {
		t.Fatal(terr)
	}
	fpath := filepath.Join(tdir, "foo.html")
	fdata := []byte("{{ no-such-stuff }}")
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	_, err := site.Load(dir)
	if assert.Error(err, "error on Load with bad template") {
		assert.Regexp("template.*foo", err.Error(),
			"error makes sense")
	}
}

func Test_Load_ErrorBadPage(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-site-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	tdir := filepath.Join(dir, "templates")
	terr := os.Mkdir(tdir, os.ModePerm)
	if terr != nil {
		t.Fatal(terr)
	}
	tpath := filepath.Join(tdir, "foo.html")
	tdata := []byte("{{ .Site.Name }}")
	if err := ioutil.WriteFile(tpath, tdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Now, the flawed page:
	pdir := filepath.Join(dir, "pages")
	perr := os.Mkdir(pdir, os.ModePerm)
	if perr != nil {
		t.Fatal(perr)
	}
	ppath := filepath.Join(pdir, "foo.md")
	pdata := []byte("# foo\n\n```json\n{ foo: [}\n```\n\n")
	if err := ioutil.WriteFile(ppath, pdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	_, err := site.Load(dir)
	if assert.Error(err, "error on Load with bad page meta") {
		assert.Regexp("invalid character", err.Error(),
			"error makes sense")
	}
}

func Test_LoadVirtualYaml_ErrorBadYaml(t *testing.T) {

	assert := assert.New(t)

	yaml := `foo: { bar [x:`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Regexp("yaml", err.Error(), "error is useful")
	}

}

func Test_LoadVirtualYaml_ErrorBadPageType(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: 1234`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Equal("Non-string value in Pages: foo.md",
			err.Error(), "error as expected")
	}

}

func Test_LoadVirtualYaml_ErrorBadTemplateType(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Templates:
    foo: 1234`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Equal("Non-string value in Templates: foo",
			err.Error(), "error as expected")
	}

}

func Test_LoadVirtualYaml_ErrorBadPageData(t *testing.T) {

	assert := assert.New(t)

	// Well this is a little annoying... now equivalent of <<HERE?
	pageData := "\"# Fresher\\n\\n```json\\n{ id: [bad,\\n```\\n\\nHere\\n\""
	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: ` + pageData
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Regexp("^Page foo.md: invalid character",
			err.Error(), "error is useful")
	}

}

func Test_LoadVirtualYaml_ErrorBadTemplateData(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Templates:
    foo: "{{ not-valid"
`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Regexp("^Template foo:",
			err.Error(), "error is useful")
	}

}

func Test_LoadVirtualYaml_ErrorBadDefaultTemplate(t *testing.T) {

	assert := assert.New(t)

	orig := site.DEFAULT_TEMPLATE
	defer func() { site.DEFAULT_TEMPLATE = orig }()

	site.DEFAULT_TEMPLATE = "{{ not-valid"
	yaml := `Name: Test Virtual Site`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Regexp("^DEFAULT_TEMPLATE",
			err.Error(), "error is useful")
	}

}

func Test_LoadVirtualYaml_ErrorBadUnlistedPaths(t *testing.T) {

	assert := assert.New(t)

	yaml := `Name: Test Virtual Site
UnlistedPaths: not-a-list`
	_, err := site.LoadVirtualYaml(yaml)
	if assert.Error(err, "error on load") {
		assert.Equal("Config UnlistedPaths is not a list.",
			err.Error(), "error is useful")
	}

}

func Test_LoadVirtualYaml_Success(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: Not an index.
    bar.md: Also not.
Templates:
    index: INDEX TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if assert.Nil(err, "no error on load") {
		assert.NotNil(s.Template.Lookup("index"), "have index template")
		pp := s.Pageset.ByPath()
		if assert.Equal(2, len(pp), "two pages parsed") {
			assert.Equal(pp[0].ModTime, pp[1].ModTime, "ModTime same for both")
		}
	}

}

func Test_LoadVirtualYaml_SuccessWithPageMeta(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: |
        # I am Foo
        
            Title: The Foo
            Created: 1970-10-01
        
        Back to the content.
Templates:
    index: INDEX TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if assert.Nil(err, "no error on load") {
		pp := s.Pageset.ByPath()
		if assert.Equal(1, len(pp), "one page parsed") {
			assert.Equal("The Foo", pp[0].Title(), "Title set from meta")
			ct := time.Unix(23587200, 0).UTC()
			assert.Equal(&ct, pp[0].Created(), "Created set from meta")
			assert.Equal("<h1>I am Foo</h1>\n\n<p>Back to the content.</p>\n",
				string(pp[0].Content), "content set after meta")
		}
	}

}

func Test_LoadVirtualYaml_SuccessWithoutPagesOrTemplates(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site`
	s, err := site.LoadVirtualYaml(yaml)
	if assert.Nil(err, "no error on load") {
		assert.Equal(0, len(s.Pageset.ByPath()), "zero pages parsed")
		assert.Nil(s.Template.Lookup("index"), "have index template")
	}

}

func Test_LoadVirtualYaml_SuccessWithUnlistedPaths(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
UnlistedPaths: [/nope, /bar]
Pages:
    /foo.md: The Foo
    /bar.md: The Bar
    /baz.md: The Baz
Templates:
    index: INDEX TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if assert.Nil(err, "no error on load") {
		listed := s.Pageset.ListedSubset()
		assert.Nil(listed.Page("/bar"), "unlisted is set as such")
		assert.NotNil(listed.Page("/foo"), "first listed still there")
		assert.NotNil(listed.Page("/baz"), "second listed still there")
	}

}
