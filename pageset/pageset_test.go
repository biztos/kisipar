// pageset/pageset_test.go - main tests for Kisipar Pagesets
// -----------------------
// TODO: test the mutex logic somehow, it's pretty mysterious right now.

package pageset_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
)

// The function below is copied from site/site_test.go.
//
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

func Test_New_DupeError(t *testing.T) {

	assert := assert.New(t)

	pages := []*page.Page{
		&page.Page{Path: "/foo/bar.md"},
		&page.Page{Path: "/foo/baz.md"},
		&page.Page{Path: "/foo/bar.md"},
		&page.Page{Path: "/foo/baz.md"},
	}

	_, err := pageset.New(pages)
	if assert.Error(err, "error on dupes") {
		assert.Equal("Duplicate path for /foo/bar: /foo/bar.md", err.Error(),
			"error as expected")
	}
}

func Test_New_Success_Empty(t *testing.T) {

	assert := assert.New(t)

	// Empty pagesets are fine.
	ps, err := pageset.New([]*page.Page{})
	assert.Nil(err, "no error returned")
	assert.Equal(0, ps.Len(), "zero pages in set")

}

func Test_New_Success(t *testing.T) {

	assert := assert.New(t)

	// Virtual pages are easy to handle here.
	p1, err := page.LoadVirtualString("/a.md", "# First!")
	p2, err := page.LoadVirtualString("/b.md", "# Second!")
	p3, err := page.LoadVirtualString("/c.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}

	ps, err := pageset.New(pages)
	assert.Nil(err, "no error returned")
	assert.Equal(3, ps.Len(), "three pages in set")

}

func Test_Page_NotFound(t *testing.T) {

	assert := assert.New(t)

	// ...much be worth abstracting this:
	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}

	ps, err := pageset.New(pages)
	assert.Nil(err, "no error returned")
	assert.Nil(ps.Page("nonesuch"), "not-found Page returns nil")

}

func Test_Page_Found(t *testing.T) {

	assert := assert.New(t)

	// ...much be worth abstracting this:
	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}

	ps, err := pageset.New(pages)
	assert.Nil(err, "no error returned")
	assert.Equal(p2, ps.Page("/b"), "found Page returns exected Page")

}

func Test_AddPage(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	pages := []*page.Page{p1, p2}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(2, ps.Len(), "size is two")

	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	ps.AddPage(p3)
	assert.Equal(3, ps.Len(), "size is three after adding one")

	p2x, _ := page.LoadVirtualString("/b.md", "# Second Redux!")
	ps.AddPage(p2x)
	assert.Equal(3, ps.Len(), "size still three after adding existing one")
	assert.Equal(p2x, ps.Page("/b"), "page replaced with added page")

	// TODO: test that it's included in various sorts etc.

}

func Test_RemovePage(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	pages := []*page.Page{p1, p2}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(2, ps.Len(), "size is two")
	assert.Equal(p2, ps.Page("/b"), "second page as expected")

	ps.RemovePage("/b")
	assert.Equal(1, ps.Len(), "size is one after removal")
	assert.Nil(ps.Page("/b"), "second page now nil")

	// TODO: test that it's removed from various sorts etc.

}

func Test_RefreshPage_NoPageForPathKey(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	pages := []*page.Page{p1, p2}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(2, ps.Len(), "size is two")

	err = ps.RefreshPage("/c")
	if assert.Error(err, "error returned when no page for path key") {
		assert.True(os.IsNotExist(err), "IsNotExist true for error")
	}

}

func Test_RefreshPage_PathNotFound(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/nonesuch.md", "# First!")
	p1.Virtual = false
	pages := []*page.Page{p1}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, ps.Len(), "size is one")

	err = ps.RefreshPage("/nonesuch")
	if assert.Error(err, "error returned refreshing a not-found page") {
		assert.True(os.IsNotExist(err), "it's a IsNotExist error")
	}

}

func Test_RefreshPage_ParseError(t *testing.T) {

	assert := assert.New(t)

	input := "# Test Page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Rewrite the file.
	fresh := "# Fresher\n\n```json\n{id: [1,2}\n```\n\nHere."
	if err := ioutil.WriteFile(path, []byte(fresh), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Tweak the mod time just in case we're running very fast.
	p.ModTime = time.Unix(0, 0)

	// Get it into our pageset now.
	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, ps.Len(), "size is one")

	// Refresh should NOT work.
	err = ps.RefreshPage(key)
	if assert.Error(err, "error thrown for bad meta block") {
		assert.Regexp("invalid character", err.Error(),
			"error as expected")
	}
	expMeta := map[string]interface{}{
		"Title": "Test Page",
	}
	expContent := `<h1>Test Page</h1>

<p>Here</p>
`

	assert.Equal([]byte(input), p.Source, "page source not reloaded")
	assert.Equal(expMeta, p.Meta, "page meta reset")
	assert.Equal(expContent, string(p.Content), "page content reset")
	assert.Equal(time.Unix(0, 0), p.ModTime, "mod time not reset")

}

func Test_RefreshPage_NoChange(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	p, err := page.Load(path)
	assert.Nil(err, "no error on Load")

	ps, err := pageset.New([]*page.Page{p})
	assert.Nil(err, "no error returned")
	assert.Equal(1, ps.Len(), "size is one")

	err = ps.RefreshPage(key)
	assert.Nil(err, "no error refreshing an already-fresh page")

}

func Test_RefreshPage_Changed(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	assert.Nil(err, "no error on Load")

	// Rewrite the file.
	fresh := `# Fresher

    { "Author": "John Woo" }

Here!
`
	if err := ioutil.WriteFile(path, []byte(fresh), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Tweak the mod time just in case we're running very fast.
	p.ModTime = time.Unix(0, 0)

	// Get it into our pageset now.
	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}

	// Refresh should work.
	err = ps.RefreshPage(key)
	assert.Nil(err, "no error refreshing a stale")

	expMeta := map[string]interface{}{
		"Title":  "Fresher",
		"Author": "John Woo",
	}
	expContent := `<h1>Fresher</h1>

<p>Here!</p>
`

	assert.Equal([]byte(fresh), p.Source, "page source reloaded")
	assert.Equal(expMeta, p.Meta, "page meta reset")
	assert.Equal(expContent, string(p.Content), "page content reset")
	assert.WithinDuration(time.Now(), p.ModTime, time.Second,
		"mod time is updated to now (give or take one second)")

}

func Test_RefreshPage_Disappeared(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Get it into our pageset now.
	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(ps.Page(key), "have page for key before refresh")

	// Zap the file.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	// Refresh should give us a not-found response.
	err = ps.RefreshPage(key)
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "error is an IsNotExist")
	}
	assert.Nil(ps.Page(key), "no page for key found after refresh")

}

func Test_RefreshPage_LoadNew(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Start with a blank set.
	ps, err := pageset.New([]*page.Page{})
	if err != nil {
		t.Fatal(err)
	}

	// Refresh should yield a new page in the pageset.
	err = ps.RefreshPage(key)
	assert.Nil(err, "no error returned")
	p := ps.Page(key)
	if assert.NotNil(p, "page returned for key after refresh") {
		assert.Equal("Test page", p.Title(), "page as expected")
	}

}

func Test_RefreshPage_LoadReplacement(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	fdata := []byte("# Test page\n\nHere\n")
	if err := ioutil.WriteFile(path, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Get it into our pageset now.
	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(ps.Page(key), "have page for key before refresh")

	// Remove that and put in a new one as a .txt file.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	path = key + ".txt"
	fdata = []byte("# Replaced page\n\nHere\n")
	if err := ioutil.WriteFile(path, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Refresh should yield the new page in the pageset.
	err = ps.RefreshPage(key)
	assert.Nil(err, "no error returned")
	p = ps.Page(key)
	if assert.NotNil(p, "page returned for key after refresh") {
		assert.Equal("Replaced page", p.Title(), "page as expected")
	}

}

func Test_RefreshPage_RefreshError(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".MD"
	fdata := []byte("# Test page\n\nHere\n")
	if err := ioutil.WriteFile(path, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Tweak the mod time just in case we're running very fast.
	p.ModTime = time.Unix(0, 0)

	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}

	// Write the new file with bad data.
	bad := []byte("# Fresher\n\n    Bad: {x[{\n\nxx\n")
	if err := ioutil.WriteFile(path, bad, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Refresh should error out.
	err = ps.RefreshPage(key)
	if assert.Error(err, "error on refresh") {
		assert.False(os.IsNotExist(err), "not an IsNotExist error")
		assert.Regexp("yaml", err.Error(), "useful error")
	}

}

func Test_RefreshPage_LoadError(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".md"
	fdata := []byte("# Fresher\n\n    Bad: {x[{\n\nxx\n")
	if err := ioutil.WriteFile(path, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Start with a plank set.
	ps, err := pageset.New([]*page.Page{})
	if err != nil {
		t.Fatal(err)
	}

	// Refresh should error out.
	err = ps.RefreshPage(key)
	if assert.Error(err, "error on refresh") {
		assert.False(os.IsNotExist(err), "not an IsNotExist error")
		assert.Regexp("yaml", err.Error(), "useful error")
	}

}

func Test_RefreshPage_NonDirInPath(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)

	key := filepath.Join(dir, "a-page")
	path := key + ".md"
	fdata := []byte("# Test page\n\nHere\n")
	if err := ioutil.WriteFile(path, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Get it into our pageset now.
	ps, err := pageset.New([]*page.Page{p})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(ps.Page(key), "have page for key before refresh")

	// Refresh on a key under this path should return a not-found error.
	badkey := filepath.Join(path, "mattersnot")
	err = ps.RefreshPage(badkey)
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "error is IsNotExist-y")
	}

}

func Test_ByPath(t *testing.T) {

	assert := assert.New(t)

	input := []string{
		"foo/aaa/bbb.md", // comes after foo/bar because longer
		"foo/bar.md",
		"abacus.md",  // comes *after* index.md
		"foo/BBQ.md", // we are case-sensitive
		"bar/index.md",
		"bar/foo.md",
		"bar/INDEX.MD", // weird edge case but possible
		"index.md",
		"other.md",
		"foo_bar.md", // comes before foo/bar because not path separator
	}
	exp_paths := []string{
		"index.md",
		"abacus.md",
		"foo_bar.md",
		"other.md",
		"bar/INDEX.MD",
		"bar/index.md",
		"bar/foo.md",
		"foo/BBQ.md",
		"foo/bar.md",
		"foo/aaa/bbb.md",
	}
	pages := make([]*page.Page, len(input))
	for i, path := range input {
		p, _ := page.LoadVirtualString(path, "# "+path)
		pages[i] = p
	}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{}
	for _, p := range ps.ByPath() {
		paths = append(paths, p.Path)
	}
	assert.Equal(exp_paths, paths, "results sorted correctly")

	// And again, for the cached version:
	cached := []string{}
	for _, p := range ps.ByPath() {
		cached = append(cached, p.Path)
	}
	assert.Equal(exp_paths, cached, "cached results sorted correctly")

}

func Test_ByPath_NoPages(t *testing.T) {

	assert := assert.New(t)
	pages := []*page.Page{}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal([]*page.Page{}, ps.ByPath(), "results empty list")
}

func Test_ByCreated(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "    Created: 2016-01-02\n\n# 1")
	p2, _ := page.LoadVirtualString("/b.md", "    Created: 2016-01-03\n\n# 2")
	p3, _ := page.LoadVirtualString("/c.md", "    Created: 2016-01-04\n\n# 3")
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p3, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly")

	// Now it's been cached, we assume:
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly (cached)")
}

func Test_ByCreated_ModTimeFallbackSome(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "    Created: 2016-01-02\n\n# 1")
	p2, _ := page.LoadVirtualString("/b.md", "    Created: 2016-01-03\n\n# 2")
	p3, _ := page.LoadVirtualString("/c.md", "# 3") // Now is new!
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p3, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly")

	// Now it's been cached, we assume:
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly (cached)")
}

func Test_ByCreated_ModTimeFallbackAll(t *testing.T) {

	assert := assert.New(t)

	// NOTE: we expect at least one ns difference on the automatic mod times.
	p1, _ := page.LoadVirtualString("/a.md", "# 1")
	p2, _ := page.LoadVirtualString("/b.md", "# 2")
	p3, _ := page.LoadVirtualString("/c.md", "# 3") // Now is new!
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p3, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly")

	// Now it's been cached, we assume:
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly (cached)")
}

func Test_ByCreated_PathOnTie(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/b/c.md", "    Created: 2016-01-02\n\n# 1")
	p2, _ := page.LoadVirtualString("/b/a.md", "    Created: 2016-01-02\n\n# 2")
	p3, _ := page.LoadVirtualString("/c/d.md", "    Created: 2016-01-03\n\n# 3")
	p4, _ := page.LoadVirtualString("/c/e.md", "    Created: 2016-01-03\n\n# 4")
	pages := []*page.Page{p4, p3, p1, p2}
	sorted := []*page.Page{p3, p4, p2, p1}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByCreated(), "results sorted correctly")

}

func Test_ByCreated_NoPages(t *testing.T) {

	assert := assert.New(t)
	pages := []*page.Page{}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal([]*page.Page{}, ps.ByCreated(), "results empty list")
}

func Test_ByModTime(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p1.ModTime = time.Unix(100, 100)
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p2.ModTime = time.Unix(200, 200)
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	p3.ModTime = time.Unix(300, 300)
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p3, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByModTime(), "results sorted correctly")

	// Now it's been cached, we assume:
	assert.Equal(sorted, ps.ByModTime(), "results sorted correctly (cached)")
}

func Test_ByModTime_Nanoseconds(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p1.ModTime = time.Unix(100, 101)
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p2.ModTime = time.Unix(100, 102)
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	p3.ModTime = time.Unix(100, 103)
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p3, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByModTime(), "results sorted correctly")
}

func Test_ByModTime_PathOnTie(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p1.ModTime = time.Unix(100, 100)
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p2.ModTime = time.Unix(200, 200)
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	p3.ModTime = time.Unix(300, 300)
	p4, _ := page.LoadVirtualString("/d.md", "# Fourth!")
	p4.ModTime = time.Unix(300, 300)
	pages := []*page.Page{p4, p3, p1, p2}
	sorted := []*page.Page{p3, p4, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByModTime(), "results sorted correctly")

}

func Test_ByModTime_NoPages(t *testing.T) {

	assert := assert.New(t)
	pages := []*page.Page{}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal([]*page.Page{}, ps.ByModTime(), "results empty list")
}

func Test_ByTime(t *testing.T) {

	assert := assert.New(t)

	s1 := `# First!
    Created: 2013-01-01
    Updated: 2016-01-01
`
	s2 := `# Second!
    Updated: 2015-01-01
`
	s3 := `# Third!
    Created: 2014-01-01
`
	p1, _ := page.LoadVirtualString("/z.md", s1)
	p1.ModTime = time.Unix(100, 100)
	p2, _ := page.LoadVirtualString("/y.md", s2)
	p2.ModTime = time.Unix(200, 200)
	p3, _ := page.LoadVirtualString("/x.md", s3)
	p3.ModTime = time.Unix(300, 300)
	pages := []*page.Page{p3, p1, p2}
	sorted := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByTime(), "results sorted correctly")

	// Now it's been cached, we assume:
	assert.Equal(sorted, ps.ByTime(), "results sorted correctly (cached)")
}

func Test_ByTime_PathOnTie(t *testing.T) {

	assert := assert.New(t)

	// Just fall back to ModTime for now.

	p1, _ := page.LoadVirtualString("/a.md", "# First!")
	p1.ModTime = time.Unix(100, 100)
	p2, _ := page.LoadVirtualString("/b.md", "# Second!")
	p2.ModTime = time.Unix(200, 200)
	p3, _ := page.LoadVirtualString("/c.md", "# Third!")
	p3.ModTime = time.Unix(300, 300)
	p4, _ := page.LoadVirtualString("/d.md", "# Fourth!")
	p4.ModTime = time.Unix(300, 300)
	pages := []*page.Page{p4, p3, p1, p2}
	sorted := []*page.Page{p3, p4, p2, p1}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(sorted, ps.ByTime(), "results sorted correctly")

}

func Test_ByTime_NoPages(t *testing.T) {

	assert := assert.New(t)
	pages := []*page.Page{}

	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal([]*page.Page{}, ps.ByTime(), "results empty list")
}

func Test_PathSubset(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	subpages := []*page.Page{p1, p2}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	subset := ps.PathSubset("/a", "")
	assert.Equal(subpages, subset.ByPath(), "subset as expected")

	// Second call is cached, though we don't know that here.
	subset = ps.PathSubset("/a", "")
	assert.Equal(subpages, subset.ByPath(), "subset as expected second time")

}

func Test_PathSubset_Trim(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/foo/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/foo/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/foo/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	subpages := []*page.Page{p1, p2}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	subset := ps.PathSubset("/a", "/foo")
	assert.Equal(subpages, subset.ByPath(), "subset as expected")

	// However, if we look for the trimmed thing itself we will not find it.
	empty, err := pageset.New([]*page.Page{})
	if err != nil {
		t.Fatal(err)
	}
	subset = ps.PathSubset("/foo", "/foo")
	assert.Equal(empty, subset, "subset empty as expected")

}

func Test_PathSubset_PanicsOnConflict(t *testing.T) {

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	// You would be insane to do this in real life, but it's possible.
	// (AFAICT it's acceptable, idiomatic Go to have a constructor function
	// that sets mutable properties, the mutation of which could put your
	// object in a bad state.  Presumably because of the time saved accessing
	// the properties vs. having an actual accessor method.)
	ps.Page("/a/a").Path = "/a/b.md"
	AssertPanicsWith(t, func() { ps.PathSubset("/a", "") },
		"PathSubset failed for Pageset: Duplicate path for /a/b: /a/b.md",
		"expected panic for dupe page key")

}

func Test_Reverse(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	reversed := []*page.Page{p3, p2, p1}
	res := pageset.Reverse(pages)

	assert.Equal(reversed, res, "reversed as expected")
	assert.NotEqual(pages, res, "not the original slice")

}

func Test_PageSlice_NegativeStart(t *testing.T) {

	pages := []*page.Page{}
	ps := &pageset.Pageset{}

	AssertPanicsWith(t, func() { ps.PageSlice(pages, -1, 1) },
		"start may not be negative",
		"expected panic")
}

func Test_PageSlice_StartAboveEnd(t *testing.T) {

	pages := []*page.Page{}
	ps := &pageset.Pageset{}

	AssertPanicsWith(t, func() { ps.PageSlice(pages, 2, 1) },
		"start may not be greater than end",
		"expected panic")
}

func Test_PageSlice_HighEnd(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	pp := ps.PageSlice(pages, 1, 99)
	assert.Equal([]*page.Page{p2, p3}, pp, "expected pages returned")
}

func Test_PageSlice_HighStart(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	pp := ps.PageSlice(pages, 4, 4)
	assert.Equal([]*page.Page{}, pp, "expected pages returned")
}

func Test_PageSlice_NegativeEnd(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	pp := ps.PageSlice(pages, 1, -1)
	assert.Equal([]*page.Page{p2, p3}, pp, "expected pages returned")
}

func Test_PageSlice_NormalRange(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/c/d.md", "# Third!")
	p4, _ := page.LoadVirtualString("/c/e.md", "# Fourth!")
	p5, _ := page.LoadVirtualString("/c/f.md", "# Fifth!")
	pages := []*page.Page{p1, p2, p3, p4, p5}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	pp := ps.PageSlice(pages, 1, 4)
	assert.Equal([]*page.Page{p2, p3, p4}, pp, "expected pages returned")
}

func Test_ListedSubset_PanicEdgeCase(t *testing.T) {

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	pages := []*page.Page{p1, p2}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	// Now create some programmer error (or sitemaster error, maybe):
	p2.Path = p1.Path

	AssertPanicsWith(t, func() { ps.ListedSubset() },
		"ListedSubset failed for Pageset: Duplicate path for /a/a: /a/a.md",
		"expected panic")

}

func Test_ListedSubset(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a/a.md", "# First!")
	p2, _ := page.LoadVirtualString("/a/b.md", "# Second!")
	p3, _ := page.LoadVirtualString("/a/c.md", "# Second!")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	p2.Unlisted = true

	exp_pages := []*page.Page{p1, p3}
	exp_ps, err := pageset.New(exp_pages)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(exp_ps, ps.ListedSubset(), "listed subset as expected")

}

func Test_TagSubset_PanicEdgeCase(t *testing.T) {

	pdata := `# A Page!

    Title: Pagey Page
    Tags: [foo, bar]

Etc.`
	p1, _ := page.LoadVirtualString("/a/a.md", pdata)
	p2, _ := page.LoadVirtualString("/a/b.md", pdata)
	pages := []*page.Page{p1, p2}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}

	// Now create some programmer error (or sitemaster error, maybe):
	p2.Path = p1.Path
	AssertPanicsWith(t, func() { ps.TagSubset("foo") },
		"TagSubset failed for Pageset: Duplicate path for /a/a: /a/a.md",
		"expected panic")

}

func Test_TagSubSet(t *testing.T) {

	assert := assert.New(t)

	pdata := `# A Page!

    Title: Pagey Page
    Tags: [foo, bar]

Etc.`

	p1, _ := page.LoadVirtualString("/a.md", pdata)
	p2, _ := page.LoadVirtualString("/b.md", "no such tag")
	p3, _ := page.LoadVirtualString("/c.md", pdata)
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	p2.Unlisted = true

	exp_pages := []*page.Page{p1, p3}
	exp_ps, err := pageset.New(exp_pages)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(exp_ps, ps.TagSubset("foo"), "tag subset as expected")

	// Second is cached; exercise it.
	assert.Equal(exp_ps, ps.TagSubset("foo"), "second round (cached) ok")

}

func Test_Tags(t *testing.T) {

	assert := assert.New(t)

	p1, _ := page.LoadVirtualString("/a.md", "# 1\n\n    Tags: [foo,bar]")
	p2, _ := page.LoadVirtualString("/b.md", "# 2\n\n    Tags: [foo,boo]")
	p3, _ := page.LoadVirtualString("/c.md", "# 3\n\n    Tags: [ZOO,bar]")
	pages := []*page.Page{p1, p2, p3}
	ps, err := pageset.New(pages)
	if err != nil {
		t.Fatal(err)
	}
	p2.Unlisted = true

	exp := []string{"bar", "foo", "zoo"}
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(exp, ps.Tags(), "tags as expected")

	// Second is cached; exercise it.
	assert.Equal(exp, ps.Tags(), "cached tags as expected")

}
