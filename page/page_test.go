// page/page_test.go -- general tests for a Kisipar Page.
// -----------------

package page_test

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/kisipar/page"
)

func Test_New_RequiresPath(t *testing.T) {

	assert := assert.New(t)

	_, err := page.New("")
	if assert.Error(err, "error returned") {
		assert.Equal("page.New requires a source path.",
			err.Error(), "error as expected")
	}

}

func Test_New_PathRequiresExtension(t *testing.T) {

	assert := assert.New(t)

	_, err := page.New("/foo/bar")
	if assert.Error(err, "error returned") {
		assert.Equal("No file extension in source path: /foo/bar",
			err.Error(), "error as expected")
	}
}

func Test_NewPageLoad_PathNotFound(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("nonesuch", "nowhere.md")
	p, err := page.New(path)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Load()
	if assert.Error(err, "error returned") {
		// NOTE: we don't care if this is a little user-unfriendly, because
		// normally you would know your file is there before setting the Page.
		assert.Equal("stat nonesuch/nowhere.md: no such file or directory",
			err.Error(), "error as expected")
	}
}

func Test_New_Success(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("anything-you-want", "a-page.md")
	p, err := page.New(path)
	if assert.Nil(err, "no error returned") {
		assert.False(p.Virtual, "page is not virtual")
		assert.False(p.Unlisted, "page is not unlisted")
		assert.False(p.IsIndex, "page is not an index")
		assert.Nil(p.Source, "source not read yet")

	}

}

func Test_NewPageLoad_ReadError(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "not-a-page.md")
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.New(path)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Load()
	if assert.Error(err, "error returned") {
		assert.Regexp("^.* is a directory",
			err.Error(), "error as expected")
	}

}

func Test_NewPageLoad_Success(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte("boo"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	p, err := page.New(path)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Load()
	if assert.Nil(err, "no error returned") {
		assert.Equal("boo", string(p.Source), "page source read from file")
		assert.Equal(info.ModTime().UTC(), p.ModTime, "mod time set correctly")
		assert.False(p.Virtual, "page is not virtual")
		assert.False(p.Unlisted, "page is not unlisted")
	}

}

func Test_IsIndex(t *testing.T) {

	assert := assert.New(t)

	index_paths := []string{
		"index.md",
		"INDEX.Md",
		"index.txt",
		"index.markdown",
		filepath.Join("foo", "index.md"),
	}
	for _, path := range index_paths {
		p, err := page.New(path)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(p.IsIndex, "IsIndex true for "+path)
	}

	nonindex_paths := []string{
		"foo.txt",
		filepath.Join("index", "foo.md"),
		"not-separator*index.md",
		filepath.Join("index.md", "foo.md"),
	}
	for _, path := range nonindex_paths {
		p, err := page.New(path)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(p.IsIndex, "IsIndex false for "+path)
	}

}

func Test_Author(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Author": "David Bowie",
		},
	}

	assert.Equal("David Bowie", p.Author(), "Author comes from Meta")

}

func Test_Time_FallBackToModTime(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	p := &page.Page{
		ModTime: mt,
		Meta:    map[string]interface{}{},
	}

	assert.Equal(&mt, p.Time(), "time fallback is ModTime")
}

func Test_Time_CreatedOnly(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	ct := time.Unix(200, 300)
	p := &page.Page{
		ModTime: mt,
		Meta: map[string]interface{}{
			"Created": ct.String(),
		},
	}

	exp := ct.String()
	got := p.Time().String()
	assert.Equal(exp, got, "time is Created if no Updated")
}

func Test_Time_UpdatedOnly(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	ut := time.Unix(200, 300)
	p := &page.Page{
		ModTime: mt,
		Meta: map[string]interface{}{
			"Updated": ut.String(),
		},
	}

	exp := ut.String()
	got := p.Time().String()
	assert.Equal(exp, got, "time is Updated if no Created")
}

func Test_Time_CreatedNewer(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	ct := time.Unix(300, 400)
	ut := time.Unix(200, 300)
	p := &page.Page{
		ModTime: mt,
		Meta: map[string]interface{}{
			"Created": ct.String(),
			"Updated": ut.String(),
		},
	}

	exp := ct.String()
	got := p.Time().String()
	assert.Equal(exp, got, "time is Created if newer")
}

func Test_Time_UpdatedNewer(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	ct := time.Unix(300, 400)
	ut := time.Unix(400, 500)
	p := &page.Page{
		ModTime: mt,
		Meta: map[string]interface{}{
			"Created": ct.String(),
			"Updated": ut.String(),
		},
	}

	exp := ut.String()
	got := p.Time().String()
	assert.Equal(exp, got, "time is Updated if newer")
}

func Test_Time_UpdatedNewer_Nanosecond(t *testing.T) {

	assert := assert.New(t)

	mt := time.Unix(100, 200)
	ct := time.Unix(300, 999999990)
	ut := time.Unix(300, 999999991)
	p := &page.Page{
		ModTime: mt,
		Meta: map[string]interface{}{
			"Created": ct.String(),
			"Updated": ut.String(),
		},
	}

	exp := ut.String()
	got := p.Time().String()
	assert.Equal(exp, got, "time is Updated if newer by one nanosecond")
}

func Test_Created(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Created": "2016.05.08",
		},
	}

	exp := time.Date(2016, 5, 8, 0, 0, 0, 0, time.UTC)
	ts := p.Created()

	assert.Equal(exp, *ts, "Created meta parsed as expected")

}

func Test_Updated(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Updated": "2016.05.08",
		},
	}

	exp := time.Date(2016, 5, 8, 0, 0, 0, 0, time.UTC)
	ts := p.Updated()

	assert.Equal(exp, *ts, "Updated meta parsed as expected")

}

func Test_Description(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Description": "David Bowie rocks!",
		},
	}

	assert.Equal("David Bowie rocks!", p.Description(),
		"Description comes from Meta")

}

func Test_Title(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Path: filepath.Join("foo", "bar-baz.md"),
		Meta: map[string]interface{}{
			"Title": "215 Celsius",
		},
	}

	assert.Equal("215 Celsius", p.Title(), "Title comes from Meta")
	p.Meta["Title"] = nil
	p.Meta["title"] = "fallback"
	assert.Equal("fallback", p.Title(), "Title falls back in Meta")
	p.Meta["title"] = nil
	assert.Equal("bar-baz", p.Title(), "Title falls back to Path filename")
	p.Path = ""
	assert.Equal("", p.Title(),
		"Title returns empty string when nothing hits")

}

func Test_Summary(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Summary": "A cactus story.",
		},
	}

	assert.Equal("A cactus story.", p.Summary(), "Summary comes from Meta")

}

func Test_Keywords(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Keywords": "cactus, city",
		},
	}

	assert.Equal("cactus, city", p.Keywords(), "Keywords comes from Meta")

}

func Test_Tags(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Tags": []string{"cactus", "city"},
		},
	}

	assert.Equal([]string{"cactus", "city"}, p.Tags(),
		"Tags come from Meta as pure array")

}

func Test_Tags_FromInterface(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"Tags": []interface{}{"cactus", "city"},
		},
	}

	assert.Equal([]string{"cactus", "city"}, p.Tags(),
		"Tags come from Meta as pure array")

}

func Test_Tags_Fallback(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{
		Meta: map[string]interface{}{
			"tags": "cactus, city",
		},
	}

	assert.Equal([]string{"cactus", "city"}, p.Tags(),
		"Tags come from Meta with fallback and split")

}

func Test_Load_RequiresPath(t *testing.T) {

	assert := assert.New(t)

	_, err := page.Load("")
	if assert.Error(err, "error returned") {
		assert.Equal("page.Load requires a source path.",
			err.Error(), "error as expected")
	}
}

func Test_Load_PathRequiresExtension(t *testing.T) {

	assert := assert.New(t)

	_, err := page.Load("/foo/bar")
	if assert.Error(err, "error returned") {
		assert.Equal("No file extension in source path: /foo/bar",
			err.Error(), "error as expected")
	}
}

func Test_Load_PathNotFound(t *testing.T) {

	assert := assert.New(t)

	path := filepath.Join("nonesuch", "nowhere.md")
	_, err := page.Load(path)
	if assert.Error(err, "error returned") {
		// NOTE: we don't care if this is a little user-unfriendly, because
		// normally you would know your file is there before setting the Page.
		assert.Equal("stat nonesuch/nowhere.md: no such file or directory",
			err.Error(), "error as expected")
	}
}

func Test_Load_ReadError(t *testing.T) {

	assert := assert.New(t)

	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "not-a-page.md")
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	_, err := page.Load(path)
	if assert.Error(err, "error returned") {
		assert.Regexp("^.* is a directory",
			err.Error(), "error as expected")
	}

}

func Test_Load_ParseError(t *testing.T) {

	assert := assert.New(t)

	input := "# Test\n\n```json\n{ bad: [1,2 }\n```\n\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	_, err := page.Load(path)
	if assert.Error(err, "error returned") {
		assert.Regexp("invalid character",
			err.Error(), "error as expected")
	}
}

func Test_Load_Success(t *testing.T) {

	assert := assert.New(t)

	input := `# Test Page

    { "id": 1234 }
 
Here we are.
`
	expMeta := map[string]interface{}{
		"Title": "Test Page",
		"id":    1234,
	}
	expContent := `<h1>Test Page</h1>

<p>Here we are.</p>
`
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	p, err := page.Load(path)
	if assert.Nil(err, "no error returned") {
		assert.Equal(input, string(p.Source), "page source read from file")
		assert.Equal(info.ModTime().UTC(), p.ModTime, "mod time set correctly")
		assert.Equal(expMeta, p.Meta, "meta set correctly")
		assert.Equal(expContent, string(p.Content), "content set correctly")
		assert.False(p.Virtual, "page is not virtual")
		assert.False(p.Unlisted, "page is not unlisted")
	}

}

func Test_LoadAny_RequiresPath(t *testing.T) {

	assert := assert.New(t)

	_, err := page.LoadAny("")
	if assert.Error(err, "error returned") {
		assert.Equal("page.LoadAny requires a source path.",
			err.Error(), "error as expected")
	}
}

func Test_LoadAny_NotFound(t *testing.T) {

	assert := assert.New(t)

	// An empty dir will of course not have the path in it.
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "nonesuch")

	p, err := page.LoadAny(path)

	assert.Nil(p, "no Page returned")
	if assert.Error(err, "error returned") {
		assert.True(os.IsNotExist(err), "error is an IsNotExist error")
	}

}

func Test_LoadAny_ParseError(t *testing.T) {

	assert := assert.New(t)
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "bad")
	fdata := []byte("# here\n\n    Bad: [a,b,c{\n\n")
	fpath := path + ".md"
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	p, err := page.LoadAny(path)
	assert.Nil(p, "no Page returned")
	if assert.Error(err, "error returned") {
		assert.False(os.IsNotExist(err), "error is NOT an IsNotExist error")
		assert.Regexp("yaml", err.Error(), "error as expected")
	}

}

func Test_LoadAny_Success_LowerCaseExtension(t *testing.T) {

	assert := assert.New(t)

	// Let's have a file this time.
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "have")
	fdata := []byte("# here\n\n    Title: Here!\n\ni am.")
	fpath := path + ".md"
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	p, err := page.LoadAny(path)
	assert.Nil(err, "no error returned")
	if assert.NotNil(p, "Page returned") {
		assert.Equal("Here!", p.Title(), "page title set")
		assert.Equal("<h1>here</h1>\n\n<p>i am.</p>\n", string(p.Content),
			"page content set")
	}

}

func Test_LoadAny_Success_UpperCaseExtension(t *testing.T) {

	assert := assert.New(t)

	// The default ExtParsers list has .MD as well as .md, for instance.
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "have")
	fpath := path + ".MD"
	fdata := []byte("# here\n\n    Title: Here!\n\ni am.")
	if err := ioutil.WriteFile(fpath, fdata, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	p, err := page.LoadAny(path)
	assert.Nil(err, "no error returned")
	if assert.NotNil(p, "Page returned") {
		assert.Equal("Here!", p.Title(), "page title set")
		assert.Equal("<h1>here</h1>\n\n<p>i am.</p>\n", string(p.Content),
			"page content set")
	}

}

func Test_LoadVirtual_RequiresPath(t *testing.T) {

	assert := assert.New(t)

	_, err := page.LoadVirtual("", []byte{})
	if assert.Error(err, "error returned") {
		assert.Equal("page.LoadVirtual requires a source path.",
			err.Error(), "error as expected")
	}
}

func Test_LoadVirtual_PathRequiresExtension(t *testing.T) {

	assert := assert.New(t)

	_, err := page.LoadVirtual("/foo/bar", []byte{})
	if assert.Error(err, "error returned") {
		assert.Equal("No file extension in source path: /foo/bar",
			err.Error(), "error as expected")
	}
}

func Test_LoadVirtual_ParseError(t *testing.T) {

	assert := assert.New(t)

	input := "# Test\n\n```json\n{ bad: [1,2 }\n```\n\n"
	_, err := page.LoadVirtual("/any/path.md", []byte(input))
	if assert.Error(err, "error returned") {
		assert.Regexp("invalid character",
			err.Error(), "error as expected")
	}

}

func Test_LoadVirtual_Success(t *testing.T) {

	assert := assert.New(t)

	input := `# Test Page

    { "id": 1234 }
 
Here we are.
`
	expMeta := map[string]interface{}{
		"Title": "Test Page",
		"id":    1234,
	}
	expContent := `<h1>Test Page</h1>

<p>Here we are.</p>
`
	p, err := page.LoadVirtual("/any/path.md", []byte(input))
	if assert.Nil(err, "no error returned") {
		assert.Equal(input, string(p.Source), "page source read from file")
		assert.Equal(expMeta, p.Meta, "meta set correctly")
		assert.Equal(expContent, string(p.Content), "content set correctly")
		assert.True(p.Virtual, "page is virtual")
		assert.False(p.Unlisted, "page is not unlisted")

		// Mod time is a bit tricky, since it is set to time.Now() but we are
		// of course no longer living in the now. Testify to the rescue!
		assert.WithinDuration(time.Now(), p.ModTime, time.Second,
			"mod time is now (give or take one second)")

	}

}

func Test_LoadVirtualString_SuccessWithJSON(t *testing.T) {

	assert := assert.New(t)

	input := `# Test Page

    { "id": 1234, "Tags": ["foo","bar"] }
 
Here we are.
`
	expMeta := map[string]interface{}{
		"Title": "Test Page",
		"id":    1234,
		"Tags":  []interface{}{"foo", "bar"},
	}
	expContent := `<h1>Test Page</h1>

<p>Here we are.</p>
`
	p, err := page.LoadVirtualString("/any/path.md", input)
	if assert.Nil(err, "no error returned") {
		assert.Equal(input, string(p.Source), "page source read from file")
		assert.Equal(expMeta, p.Meta, "meta set correctly")
		assert.Equal(expContent, string(p.Content), "content set correctly")
		assert.True(p.Virtual, "page is virtual")
		assert.False(p.Unlisted, "page is not unlisted")

		// Mod time is a bit tricky, since it is set to time.Now() but we are
		// of course no longer living in the now. Testify to the rescue!
		assert.WithinDuration(time.Now(), p.ModTime, time.Second,
			"mod time is now (give or take one second)")

	}

}

func Test_LoadVirtualString_SuccessWithYAML(t *testing.T) {

	assert := assert.New(t)

	input := `# Test Page

    # Some YAML eh?
    Title: Test Page
    id: 1234
    Tags: [foo, bar]

Here we are.
`
	expMeta := map[string]interface{}{
		"Title": "Test Page",
		"id":    1234,
		"Tags":  []interface{}{"foo", "bar"},
	}
	expContent := `<h1>Test Page</h1>

<p>Here we are.</p>
`
	p, err := page.LoadVirtualString("/any/path.md", input)
	if assert.Nil(err, "no error returned") {
		assert.Equal(input, string(p.Source), "page source read from file")
		assert.Equal(expMeta, p.Meta, "meta set correctly")
		assert.Equal(expContent, string(p.Content), "content set correctly")
		assert.True(p.Virtual, "page is virtual")
		assert.False(p.Unlisted, "page is not unlisted")

		// Mod time is a bit tricky, since it is set to time.Now() but we are
		// of course no longer living in the now. Testify to the rescue!
		assert.WithinDuration(time.Now(), p.ModTime, time.Second,
			"mod time is now (give or take one second)")

	}

}

func Test_Refresh_Virtual(t *testing.T) {

	assert := assert.New(t)

	// W/o path it will fail if we try to do a load.
	p := &page.Page{Virtual: true}
	err := p.Refresh()
	assert.Nil(err, "no error on virtual refresh that would otherwise fail")

}

func Test_Refresh_SameTimes(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	// Fiddle with the values so we can tell it's not reloaded.
	p.Source = []byte("s")
	p.Meta = map[string]interface{}{"m": "X"}
	p.Content = template.HTML("c")

	// Refresh should be a noop.
	err = p.Refresh()
	assert.Nil(err, "no error on Refresh")

	assert.Equal("s", string(p.Source), "page source not reloaded")
	assert.Equal("c", string(p.Content), "page content not reset")
	assert.Equal(map[string]interface{}{"m": "X"},
		p.Meta, "page meta not reset")

}

func Test_Refresh_NotFound(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	assert.Nil(err, "no error on Load")

	// Remove the temp file so it can not be reloaded:
	os.RemoveAll(dir)

	// Fiddle with the values so we can tell it's not reloaded.
	p.Source = []byte("s")
	p.Meta = map[string]interface{}{"m": "X"}
	p.Content = template.HTML("c")

	// We should get a checkable error here:
	err = p.Refresh()
	assert.Error(err, "have error on Refresh")
	assert.True(os.IsNotExist(err), "error passes os.IsNotExist")

	assert.Equal("s", string(p.Source), "page source not reloaded")
	assert.Equal("c", string(p.Content), "page content not reset")
	assert.Equal(map[string]interface{}{"m": "X"},
		p.Meta, "page meta not reset")

}

func Test_Refresh_ParseFailure(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
	if err := ioutil.WriteFile(path, []byte(input), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	p, err := page.Load(path)
	assert.Nil(err, "no error on Load")

	// Rewrite the file.
	fresh := "# Fresher\n\n```json\n{ id: [bad,\n```\n\nHere\n"
	if err := ioutil.WriteFile(path, []byte(fresh), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Tweak the mod time just in case we're running very fast.
	p.ModTime = time.Unix(0, 0)

	// Fiddle with the values so we can tell it's not reloaded.
	p.Source = []byte("s")
	p.Meta = map[string]interface{}{"m": "X"}
	p.Content = template.HTML("c")

	// Refresh should return a parse error.
	err = p.Refresh()
	if assert.Error(err, "error returned from Refresh") {
		assert.Regexp("invalid character", err.Error(), "error is useful")
	}

	assert.Equal("s", string(p.Source), "page source not reloaded")
	assert.Equal("c", string(p.Content), "page content not reset")
	assert.Equal(map[string]interface{}{"m": "X"},
		p.Meta, "page meta not reset")

}

func Test_Refresh_Success(t *testing.T) {

	assert := assert.New(t)

	input := "# Test page\n\nHere\n"
	dir, derr := ioutil.TempDir("", "kisipar-page-test-")
	if derr != nil {
		t.Fatal(derr)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "a-page.MD")
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

	// Refresh should work.
	err = p.Refresh()
	assert.Nil(err, "no error on Refresh")

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

func Test_LoadVirtual_SetUnlisted(t *testing.T) {

	assert := assert.New(t)

	source := `
# unlist-me

    Created: 2016-01-02
    Unlisted: true
`
	p, err := page.LoadVirtualString("/foo.md", source)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(p.Unlisted, "Unlisted property set from page Meta")

}

func Test_LoadVirtual_SetUnlisted_UppercaseFallback(t *testing.T) {

	assert := assert.New(t)

	source := `
# unlist-me

    Created: 2016-01-02
    UNLISTED: true
`
	p, err := page.LoadVirtualString("/foo.md", source)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(p.Unlisted, "Unlisted property set from page Meta (upper)")

}

func Test_LoadVirtual_SetUnlisted_LowerFallback(t *testing.T) {

	assert := assert.New(t)

	source := `
# unlist-me

    Created: 2016-01-02
    unlisted: true
`
	p, err := page.LoadVirtualString("/foo.md", source)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(p.Unlisted, "Unlisted property set from page Meta (lower)")

}

func Test_String(t *testing.T) {

	assert := assert.New(t)

	allYaml := `
# have all!

    Created: 2016-01-02
    Updated: 2016-02-03
`
	hasAll, err := page.LoadVirtualString("/foo.md", allYaml)
	if err != nil {
		t.Fatal(err)
	}
	hasAll.ModTime = time.Unix(1000, 2000).UTC() // for predictability

	assert.Equal("/foo.md: have all! "+
		"(Created: 2016-01-02 00:00:00 +0000 UTC; "+
		"Updated: 2016-02-03 00:00:00 +0000 UTC; "+
		"ModTime: 1970-01-01 00:16:40.000002 +0000 UTC)", hasAll.String(),
		"stringifies correctly with all times")

	hasNone, err := page.LoadVirtualString("/foo.md", "# have none!")
	if err != nil {
		t.Fatal(err)
	}
	hasNone.ModTime = time.Unix(1000, 2000).UTC() // for predictability

	assert.Equal("/foo.md: have none! "+
		"(ModTime: 1970-01-01 00:16:40.000002 +0000 UTC)", hasNone.String(),
		"stringifies correctly with only ModTime")

}
