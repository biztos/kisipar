// standardprovider_test.go

package kisipar_test

import (
	// Standard:
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"testing"
	"time"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

// This is presumably the minimal Pather:
type TestPather string

func (p TestPather) Path() string { return string(p) }

// This is presumably as minimal a working Pather-Stubber as you can make:
type TestPatherStubber string

func (p TestPatherStubber) Path() string       { return string(p) }
func (p TestPatherStubber) IsPageStub() bool   { return false }
func (p TestPatherStubber) TypeString() string { return "TestPatherStubber" }
func (p TestPatherStubber) Stub() kisipar.Stub { return p }

func Test_InterfaceConformity(t *testing.T) {

	// This will crash if anything doesn't match.
	var f1 = func(ds kisipar.Provider) {
		t.Log(ds)
	}
	f1(&kisipar.StandardProvider{})

	// ...and so on...
	var f2 = func(x kisipar.PageStub) {
		t.Log(x)
	}
	f2(&kisipar.StandardPageStub{})

}

func Test_BasicStub(t *testing.T) {

	assert := assert.New(t)

	bs := kisipar.NewBasicStub("/foo")
	assert.Equal("/foo", bs.Path(), "Path works")
	assert.Equal("BasicStub", bs.TypeString(), "TypeString works")
	assert.False(bs.IsPageStub(), "IsPageStub works")

}

func Test_StandardPageStub(t *testing.T) {

	assert := assert.New(t)

	ps := &kisipar.StandardPageStub{}
	assert.Equal("StandardPageStub", ps.TypeString(), "TypeString works")
	assert.True(ps.IsPageStub(), "IsPageStub works")
	assert.True(ps.IsPageStub(), "IsPage works")

}

func Test_NewStandardContent(t *testing.T) {

	assert := assert.New(t)

	sc := kisipar.NewStandardContent("/foo", "text/foo", "helo", time.Unix(1, 2))

	assert.Equal("/foo", sc.Path(), "Path works")
	assert.Equal("text/foo", sc.ContentType(), "ContentType works")
	assert.Equal(time.Unix(1, 2), sc.ModTime(), "ModTime works")
	rs := sc.ReadSeeker()

	p := make([]byte, 4)
	n, err := rs.Read(p)
	if !assert.Nil(err) {
		t.Log(err)
	}
	assert.Equal(n, 4, "content bytes read from ReadSeeker")
	assert.Equal("helo", string(p), "expected bytes read from ReadSeeker")

}

func Test_NewStandardContent_ModTimeDefault(t *testing.T) {

	assert := assert.New(t)

	sc := kisipar.NewStandardContent("/foo", "text/foo", "helo", time.Time{})

	assert.WithinDuration(time.Now(), sc.ModTime(), time.Second,
		"ModTime now for zero time")

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
	assert.Equal(template.HTML("<h1>foo</h1>"), p.HTML(), "HTML")
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
		"html":    "<h1>hey</h1>",
	}

	tStr := map[string]string{
		"path":    "string",
		"title":   "string",
		"tags":    "string slice",
		"created": "Time",
		"updated": "Time",
		"meta":    "string-interface map",
		"html":    "string",
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

func Test_StandardPageFromData_EmptyPathError(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{"path": ""}
	_, err := kisipar.StandardPageFromData(input)
	if assert.Error(err) {
		assert.Equal("path may not be an empty string", err.Error(),
			"error as expected")
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

func Test_StandardPage_Stub(t *testing.T) {

	assert := assert.New(t)

	src := kisipar.NewStandardPage(
		"/foo",                                  // path
		"The Foo",                               // title
		[]string{"boo", "hoo"},                  // tags
		time.Unix(0, 0),                         // created
		time.Unix(10000, 0),                     // updated
		map[string]interface{}{"helo": "WORLD"}, // meta
		"<h1>foo</h1>",                          // html
	)

	// This is a bit wonky but you'd normally not have to deal with it
	// directly: Go thinks this is a plain Stub but it's actually a PageStub;
	// we have to cast it as such before we can access its other methods.
	s := src.Stub()
	p, _ := s.(kisipar.PageStub)

	// The stub is identical to the page but also has stubby stuff.
	if !assert.True(s.IsPageStub(), "it's a page stub") {
		t.Fatal(fmt.Sprintf("%T", s))
	}

	assert.Equal("The Foo", p.Title(), "Title")
	assert.Equal([]string{"boo", "hoo"}, p.Tags(), "Tags")
	assert.Equal(time.Unix(0, 0), p.Created(), "Created")
	assert.Equal(time.Unix(10000, 0), p.Updated(), "Updated")
	assert.Equal(map[string]interface{}{"helo": "WORLD"}, p.Meta(), "Meta")

}

func Test_StandardPage_FlexMetaString(t *testing.T) {

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

	assert.Equal("WORLD", p.FlexMetaString("HELO"), "FlexMetaString")
}

func Test_NewStandardProvider(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()

	assert.Regexp("<StandardProvider with 0 items, updated .*>",
		sp.String(), "Stringifies as expected")
}

func Test_StandardProviderFromYaml_YamlError(t *testing.T) {

	assert := assert.New(t)

	yaml := `# I am (not really proper) YAML!
foo: {-- ,]`

	_, err := kisipar.StandardProviderFromYaml(yaml)
	if assert.Error(err) {
		assert.Regexp("yaml", err, "Error is useful")
	}
}

func Test_StandardProviderFromYaml_TemplateError(t *testing.T) {

	assert := assert.New(t)

	yaml := `# I am YAML with a bad template!
pages:
    /foo/bar:
        title: I am the Foo Bar!
templates:
    any/random/tmpl.html: |
        {{ foreach .Nope }}
`

	_, err := kisipar.StandardProviderFromYaml(yaml)
	if assert.Error(err) {
		assert.Regexp("template", err, "Error is useful")
	}
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

func Test_StandardProvider_GetSince_ErrNotExist(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	_, err := sp.GetSince("nada", time.Now())
	if assert.Error(err) {
		assert.Equal(kisipar.ErrNotExist, err, "ErrNotExist returned")
		assert.True(kisipar.IsNotExist(err), "IsNotExist satisfied")
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

func Test_StandardProvider_GetStub_ErrNotExist(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	_, err := sp.GetStub("nada")
	if assert.Error(err) {
		assert.Equal(kisipar.ErrNotExist, err, "ErrNotExist returned")
		assert.True(kisipar.IsNotExist(err), "IsNotExist satisfied")
	}

}

func Test_StandardProvider_GetStub_ErrNotStubber(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	p := TestPather("stubless")
	sp.Add(p)
	_, err := sp.GetStub("stubless")
	if assert.Error(err) {
		assert.Equal(kisipar.ErrNotStubber, err, "ErrNotStubber returned")
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

func Test_StandardProvider_GetPageStubs(t *testing.T) {

	assert := assert.New(t)

	p1 := kisipar.NewStandardPage(
		"/foo/p1",                               // path
		"The Foo 1",                             // title
		[]string{"boo", "hoo"},                  // tags
		time.Unix(0, 0),                         // created
		time.Unix(10000, 0),                     // updated
		map[string]interface{}{"helo": "WORLD"}, // meta
		"<h1>foo</h1>",                          // html
	)
	p2 := kisipar.NewStandardPage(
		"/foo/p2",                               // path
		"The Foo 2",                             // title
		[]string{"boo", "hoo"},                  // tags
		time.Unix(0, 0),                         // created
		time.Unix(10000, 0),                     // updated
		map[string]interface{}{"helo": "WORLD"}, // meta
		"<h1>foo</h1>",                          // html
	)

	sp := kisipar.NewStandardProvider()
	sp.Add(TestPatherStubber("/foo"))
	sp.Add(TestPatherStubber("/foodie"))
	sp.Add(TestPather("/foo/bar/baz")) // not stubber...
	sp.Add(p1)
	sp.Add(p2)

	got := sp.GetPageStubs("/foo")
	assert.Equal(2, len(got), "got expected number of items")
	exp := []string{"/foo/p1", "/foo/p2"}
	paths := []string{}
	for _, v := range got {
		paths = append(paths, v.Path())
	}
	assert.Equal(exp, paths, "got stubs at expected paths")

}

func Test_StandardProvider_GetStubs(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	sp.Add(TestPatherStubber("/foo"))
	sp.Add(TestPatherStubber("/foodie"))
	sp.Add(TestPather("/foo/bar/baz")) // not stubber...
	sp.Add(TestPatherStubber("/foment"))

	got := sp.GetStubs("/foo")
	assert.Equal(2, len(got), "got expected number of items")
	exp := []string{"/foo", "/foodie"}
	paths := []string{}
	for _, v := range got {
		paths = append(paths, v.Path())
	}
	assert.Equal(exp, paths, "got stubs at expected paths")

}

func Test_StandardProvider_GetAll(t *testing.T) {

	assert := assert.New(t)

	sp := kisipar.NewStandardProvider()
	sp.Add(TestPatherStubber("/foo"))
	sp.Add(TestPatherStubber("/foodie"))
	sp.Add(TestPather("/foo/bar/baz")) // not stubber...
	sp.Add(TestPatherStubber("/foment"))

	got := sp.GetAll("/foo")
	assert.Equal(3, len(got), "got expected number of items")
	exp := []string{"/foo", "/foodie", "/foo/bar/baz"}
	paths := []string{}
	for _, v := range got {
		paths = append(paths, v.Path())
	}
	assert.Equal(exp, paths, "got stubs at expected paths")

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
