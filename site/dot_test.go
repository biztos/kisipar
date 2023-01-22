// site/dot_test.go - main tests for the Kisipar Dot.
// ----------------

package site_test

import (
	"net/http"
	"testing"

	// Third-party packages:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/site"
)

func Test_Dot_SetRegister(t *testing.T) {

	assert := assert.New(t)

	d := &site.Dot{Register: 123}

	o := d.SetRegister(321)
	assert.Equal(123, o, "old value returned")
	assert.Equal(321, d.Register, "new value set")

}

func Test_Dot_URL(t *testing.T) {

	assert := assert.New(t)

	emptyDot := &site.Dot{}
	assert.Equal("", emptyDot.URL(), "URL empty for nil Request")

	r, _ := http.NewRequest("GET", "http://example.com/foo/bar?baz+bat", nil)
	rDot := &site.Dot{Request: r}
	assert.Equal("http://example.com/foo/bar?baz+bat", rDot.URL(),
		"URL matches for normal Request")
}

func Test_Dot_Template_InsufficientDataShouldPanic(t *testing.T) {

	panickyDot := &site.Dot{}
	panickyFunc := func() { panickyDot.Template() }

	// Need a Site:
	AssertPanicsWith(t, panickyFunc, "Site is nil",
		"Template call panics for empty Dot")

	// Need a Template too:
	panickyDot.Site = &site.Site{}
	AssertPanicsWith(t, panickyFunc, "Site.Template is nil",
		"Template call panics for Dot with an empty Site (no Template)")

}

func Test_Dot_Template_PageMetaShouldOverride(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: Not an index.
    bar.md: |
        # Bar!
        
            Template: foo
        
        OK then.
Templates:
    index: INDEX TEMPLATE
    foo: FOO TEMPLATE
    bar: BAR TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "http://example.com/bar", nil)
	dot := &site.Dot{
		Site:    s,
		Request: r,
		Page:    s.Pageset.Page("bar"),
		Pageset: s.Pageset,
	}

	tmpl := dot.Template()
	if assert.NotNil(tmpl, "template returned") {
		assert.Equal("foo", tmpl.Name(),
			"foo is bar's template per the meta")

	}

}

func Test_Dot_Template_StandardMatch(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: Not an index.
    bar.md: We are bar...
Templates:
    index: INDEX TEMPLATE
    foo: FOO TEMPLATE
    bar: BAR TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "http://example.com/bar", nil)
	dot := &site.Dot{
		Site:    s,
		Request: r,
		Page:    s.Pageset.Page("bar"),
		Pageset: s.Pageset,
	}

	tmpl := dot.Template()
	if assert.NotNil(tmpl, "template returned") {
		assert.Equal("bar", tmpl.Name(), "bar is bar's template")

	}

}

// TODO: break this down into smaller chunks...
func Test_Dot_Template(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: Not an index.
    bar.md: We are bar...
    foo/bar.md: Also bar, under foo.
Templates:
    index: INDEX TEMPLATE
    single: SINGLE TEMPLATE
    foo: FOO TEMPLATE
    foo/bar: FOO BAR TEMPLATE
    bar/single: BAR SINGLE TEMPLATE
    bar/index: BAR INDEX TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "http://example.com/bar", nil)
	dot := &site.Dot{
		Site:    s,
		Request: r,
		Page:    s.Pageset.Page("bar"), // any page will do.
		Pageset: s.Pageset,
	}

	// Exact match for the URL (case insensitive):
	dot.Request, _ = http.NewRequest("GET", "http://example.com/FOO/bar", nil)
	assert.Equal("foo/bar", dot.Template().Name(),
		"right template for exact match")

	// Match at one level up:
	dot.Request, _ = http.NewRequest("GET", "http://xxx/foo/bar/baz", nil)
	assert.Equal("foo/bar", dot.Template().Name(),
		"right template for exact match one level up")

	// Match for alt single at level:
	dot.Pageset = nil
	dot.Request, _ = http.NewRequest("GET", "http://example.com/bar", nil)
	assert.Equal("bar/single", dot.Template().Name(),
		"right template for single match at level")

	// Match for alt single one level up:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/bar/xxx", nil)
	assert.Equal("bar/single", dot.Template().Name(),
		"right template for single match one up")

	// Fallback to main single:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/x/x", nil)
	assert.Equal("single", dot.Template().Name(),
		"right template for single match at root")

	// For index matches we need a pageset:
	dot.Pageset = s.Pageset

	// Match for index at level:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/bar", nil)
	assert.Equal("bar/index", dot.Template().Name(),
		"right template for index match one up")

	// Match for index one level up:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/bar/x", nil)
	assert.Equal("bar/index", dot.Template().Name(),
		"right template for index match one up")

	// Fallback to main index:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/x/x", nil)
	assert.Equal("index", dot.Template().Name(),
		"right template for index match at root")

	// Main index for top request:
	dot.Request, _ = http.NewRequest("GET", "http://example.com/", nil)
	assert.Equal("index", dot.Template().Name(),
		"right template for index match at root (GET /)")

}

func Test_Dot_Template_IndexForNoIndexPage(t *testing.T) {

	assert := assert.New(t)

	yaml := `# TEST
Name: Test Virtual Site
Pages:
    foo.md: Not an index.
Templates:
    index: INDEX TEMPLATE`
	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	dot := &site.Dot{
		Site:    s,
		Request: r,
		Page:    nil,
		Pageset: s.Pageset,
	}

	tmpl := dot.Template()
	if assert.NotNil(tmpl, "template returned") {
		assert.Equal("index", tmpl.Name(),
			"'tis the index template")

	}

}
