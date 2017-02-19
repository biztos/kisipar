// standardprovider_test.go

package kisipar_test

import (
	// Standard:
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

// This is presumably as minimal a working Pather-Stubber as you can make:
type TestPatherStubber string

func (p TestPatherStubber) Path() string       { return string(p) }
func (p TestPatherStubber) IsPageStub() bool   { return false }
func (p TestPatherStubber) TypeString() string { return "TestPatherStubber" }
func (p TestPatherStubber) Stub() kisipar.Stub { return p }

func Test_InterfaceConformity(t *testing.T) {

	// This will crash if anything doesn't match.
	var f = func(ds kisipar.Provider) {
		t.Log(ds)
	}
	f(&kisipar.StandardProvider{})

}

func Test_IsNotExist(t *testing.T) {

	assert := assert.New(t)

	assert.True(kisipar.IsNotExist(kisipar.ErrNotExist), "local ErrNotExist")
	assert.True(kisipar.IsNotExist(os.ErrNotExist), "os.ErrNotExist")
	assert.False(kisipar.IsNotExist(errors.New("other")), "other error")

}

func Test_NewStandardPage(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"/foo",                                  // path
		"The Foo",                               // title
		[]string{"boo", "hoo"},                  // tags
		time.Unix(0, 0),                         // created
		time.Unix(10000, 0),                     // updated
		map[string]interface{}{"helo": "WORLD"}, // meta
		"<h1>foo</h1>",                          // html
	)

	assert.Equal("The Foo", p.Title(), "Title")
	assert.Equal([]string{"boo", "hoo"}, p.Tags(), "Tags")
	assert.Equal(time.Unix(0, 0), p.Created(), "Created")
	assert.Equal(time.Unix(10000, 0), p.Updated(), "Updated")
	assert.Equal(map[string]interface{}{"helo": "WORLD"}, p.Meta(), "Meta")

}

func Test_StandardPageFromData(t *testing.T) {

	assert := assert.New(t)

	cr := time.Unix(0, 0)
	up := time.Now()
	input := map[string]interface{}{
		"path":           "/foo/bar",
		"id":             "possibly-unique",
		"title":          "Hello World",
		"tags":           []string{"foo", "bar"},
		"created":        cr,
		"updated":        up,
		"meta":           map[string]interface{}{"foo": "bar"},
		"something else": &kisipar.StandardPage{},
	}

	p, err := kisipar.StandardPageFromData(input)
	if assert.Nil(err, "no error") {
		assert.Equal("Hello World", p.Title(), "Title")
		assert.Equal([]string{"foo", "bar"}, p.Tags(), "Tags")
		assert.Equal(cr, p.Created(), "Created")
		assert.Equal(up, p.Updated(), "Updated")
		assert.Equal(map[string]interface{}{"foo": "bar"}, p.Meta(), "Meta")
	}

}

func Test_StandardPageFromData_TypeErrors(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{
		"path":    "/foo/bar",
		"title":   "Hello World",
		"tags":    []string{"foo", "bar"},
		"created": time.Time{},
		"updated": time.Time{},
		"meta":    map[string]interface{}{"foo": "bar"},
	}

	tStr := map[string]string{
		"path":    "string",
		"title":   "string",
		"tags":    "string slice",
		"created": "Time",
		"updated": "Time",
		"meta":    "string-interface map",
	}
	type NotMyType struct{}
	for k, _ := range input {
		// nice dumb copy
		in2 := map[string]interface{}{}
		for k2, v2 := range input {
			in2[k2] = v2
		}

		// now override the one thing
		in2[k] = NotMyType{}

		_, err := kisipar.StandardPageFromData(in2)
		if assert.Error(err) {

			exp := fmt.Sprintf("Wrong type for %s: kisipar_test.NotMyType not %s",
				k, tStr[k])
			assert.Equal(exp, err.Error(), "error as expected")
		} else {
			t.Fatalf("No error but we reset %s!", k)
		}
	}

}

func Test_StandardPageFromData_MinimalData(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{"path": "/foo/bar"}

	p, err := kisipar.StandardPageFromData(input)
	if assert.Nil(err, "no error") {
		assert.Equal("/foo/bar", p.Path(), "Path")
		assert.Zero("", p.Title(), "Title")
		assert.Zero(p.Tags(), "Tags")
		assert.Zero(p.Created(), "Created")
		assert.Zero(p.Updated(), "Updated")
		assert.Zero(p.Meta(), "Meta")
	}

}

func Test_StandardPageFromData_EmptyData(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{}

	_, err := kisipar.StandardPageFromData(input)
	if assert.Error(err, "error returned") {
		assert.Equal("path not set", err.Error(), "error as expected")
	}

}

func Test_StandardPage_MetaString(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"/foo",                 // path
		"The Foo",              // title
		[]string{"boo", "hoo"}, // tags
		time.Unix(0, 0),        // created
		time.Unix(10000, 0),    // updated

		// meta:
		map[string]interface{}{
			"helo": "WORLD",
			"pi":   3.1415,
			"nada": nil,
			"ts":   time.Unix(0, 0),
		},

		// html
		"<h1>helo</h1>",
	)

	expTs := time.Unix(0, 0).String() // includes local TS
	assert.Equal("WORLD", p.MetaString("helo"), "string -> string")
	assert.Equal("3.1415", p.MetaString("pi"), "float -> string")
	assert.Equal("", p.MetaString("nada"), "real nil -> string")
	assert.Equal("", p.MetaString("nonesuch"), "missing -> string")
	assert.Equal(expTs, p.MetaString("ts"), "time -> string")

	// nil source
	p = &kisipar.StandardPage{}
	assert.Equal("", p.MetaString("any"), "nil meta -> empty string")
}

func Test_StandardPage_MetaStrings(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"/foo",                 // path
		"The Foo",              // title
		[]string{"boo", "hoo"}, // tags
		time.Unix(0, 0),        // created
		time.Unix(10000, 0),    // updated

		// meta:
		map[string]interface{}{
			"strings": []string{"fee", "fi", "fo"},
			"single":  "HOLA",
			"nada":    nil,
			"ts":      time.Unix(0, 0),
			"mixed":   []interface{}{time.Unix(0, 0), 3.1415, "flubber"},
		},

		// html
		"<h1>helo</h1>",
	)

	expTs := time.Unix(0, 0).String() // includes local TS
	assert.Equal([]string{"fee", "fi", "fo"}, p.MetaStrings("strings"), "strings")
	assert.Equal([]string{"HOLA"}, p.MetaStrings("single"), "one string")

	assert.Equal([]string{}, p.MetaStrings("nada"), "real nil")
	assert.Equal([]string{}, p.MetaStrings("nonesuch"), "missing")

	assert.Equal([]string{expTs}, p.MetaStrings("ts"), "single time")

	assert.Equal([]string{expTs, "3.1415", "flubber"},
		p.MetaStrings("mixed"), "mixed slice")

	// nil source
	p = &kisipar.StandardPage{}
	assert.Equal([]string{}, p.MetaStrings("any"),
		"nil meta -> empty string slice")
}

func Test_NewStandardProvider(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()

	assert.Regexp("<StandardProvider with 0 items, updated .*>",
		sp.String(), "Stringifies as expected")
}

func Test_StandardProviderFromYaml(t *testing.T) {

	assert := assert.New(t)

	yaml := `# I am YAML!
pages:
    /foo/bar:
        title: I am the Foo Bar!
        tags: [foo,bar]
        created: 2016-01-02T15:04:05Z
        updated: 2017-02-02T15:04:05Z
        content: |
            This is the foo, the bar, the baz and
            the bat if you like.  For sanity's sake
            let's not let it be Markdown.
    /baz/bat:
        title: The BazzerBat
        tags: [foo,bazzers,badgers]
content:
    /js/goober.js:
        type: application/javascript
        content: |
            window.alert('hello world');
templates:
    any/random/tmpl.html: |
        <!doctype html>
        <script src="/js/goober.js"></script>
        <h1>Hello {{ .Title }}</h1>
`

	sp, err := kisipar.StandardProviderFromYaml(yaml)
	if assert.Nil(err, "no error") {

		assert.Regexp("<StandardProvider with 3 items, updated .*>",
			sp.String(), "Stringifies as expected")
		return
		foo, err := sp.Get("/foo/bar")
		if assert.Nil(err, "no err getting first item") {
			assert.Implements((*kisipar.Page)(nil), foo, "it's a Page")
		}
		bat, err := sp.Get("/baz/bat")
		if assert.Nil(err, "no err getting second item") {
			assert.Implements((*kisipar.Page)(nil), bat, "it's a Page")
		}
		goob, err := sp.Get("/js/goober.js")
		if assert.Nil(err, "no err getting third item") {
			assert.Implements((*kisipar.Content)(nil), goob, "it's a Content")
		}

		tmpl := sp.Template()
		assert.NotNil(tmpl, "Template not nil")

	}
}

func Test_StandardProvider_Add(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()

	assert.Equal(0, sp.Count(), "zero items")
	u := sp.Updated()
	sp.Add(TestPatherStubber("dummy"))
	assert.Equal(1, sp.Count(), "one item")
	assert.True(u.Before(sp.Updated()), "Updated moves forward")

}

func Test_StandardProvider_Get(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	p := TestPatherStubber("dummy")
	sp.Add(p)

	got, err := sp.Get("dummy")
	if assert.Nil(err, "no error getting item") {
		assert.Equal(p, got, "got expected item")
	}

	got, err = sp.Get("not dummy")
	if assert.NotNil(err, "error for missing item") {
		assert.Equal(kisipar.ErrNotExist, err, "standard error")
		// And obviously:
		assert.True(kisipar.IsNotExist(err), "IsNotExist true for error")
	}

}

func Test_StandardProvider_GetSince(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	o := TestPatherStubber("older")
	sp.Add(o)
	ts := sp.Updated()
	n := TestPatherStubber("newer")
	sp.Add(n)

	f := kisipar.NewStandardFile("file", "standardprovider_test.go")
	sp.Add(f)

	got, err := sp.GetSince("newer", ts)
	if assert.Nil(err, "no error getting newer item") {
		assert.Equal(n, got, "got expected item")
	}

	got, err = sp.GetSince("older", ts)
	if assert.NotNil(err, "error for older item") {
		assert.Equal(kisipar.ErrNotModified, err, "standard error")
	}

	// The file, however, is special: it's the newest thing added, but its
	// mod time is (by definition) older than the running of this test.
	got, err = sp.GetSince("file", ts)
	if assert.NotNil(err, "error for file") {
		assert.Equal(kisipar.ErrNotModified, err, "standard error")
	}

}

func Test_StandardProvider_GetStub(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	p := TestPatherStubber("dummy")
	sp.Add(p)

	got, err := sp.GetStub("dummy")
	if assert.Nil(err, "no error getting item") {
		assert.Equal(p, got, "got expected item")
	}

}

func Test_StandardProvider_GetStubs(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	sp.Add(TestPatherStubber("foo"))
	sp.Add(TestPatherStubber("foodie"))
	sp.Add(TestPatherStubber("foo/bar/baz"))
	sp.Add(TestPatherStubber("foment"))

	got := sp.GetStubs("foo")
	assert.Equal(3, len(got), "got expected number of items")

}

// Mostly tested under templates_test.go, cf. TemplatesFromData.
func Test_StandardProvider_TemplateFor(t *testing.T) {

	assert := assert.New(t)

	yaml := `# I am YAML!
pages:
    /foo/bar:
        title: Anything
templates:
    foo/bar.html: |
        HERE
`

	sp, err := kisipar.StandardProviderFromYaml(yaml)
	if err != nil {
		panic(err)
	}
	p, err := sp.Get("/foo/bar")
	if err != nil {
		panic(err)
	}
	page, ok := p.(*kisipar.StandardPage)
	if !ok {
		panic(p)
	}
	tmpl := sp.TemplateFor(page)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE\n", b.String(), "content as expected")
		}
	}
}
