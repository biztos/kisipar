// templates_test.go -- tests for general Provider template logic
//
package provider_test

import (
	// Standard:
	"bytes"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar/provider"
)

func Test_TemplateThemes(t *testing.T) {

	assert := assert.New(t)

	exp := []string{"debug", "default", "naked", "wonky"}
	assert.Equal(exp, provider.TemplateThemes(), "themes as expected")
}

func Test_TemplatesFromData_NilMap(t *testing.T) {

	assert := assert.New(t)

	_, err := provider.TemplatesFromData(nil)
	if assert.Error(err, "error returned") {
		assert.Equal("TemplatesFromData input may not be nil.",
			err.Error(), "error as expected")
	}

}

func Test_TemplatesFromData(t *testing.T) {

	assert := assert.New(t)

	input := map[string]string{
		"foo":          "the foo",
		"bar/html":     "the bar html no ext",
		"baz/bat.html": "the baz-bat, realistic-ishly",
	}

	tmpl, err := provider.TemplatesFromData(input)
	if assert.Nil(err, "no error returned") {
		for k, v := range input {
			if got := tmpl.Lookup(k); assert.NotNil(got, "got "+k) {
				var b bytes.Buffer
				if assert.Nil(got.Execute(&b, nil), "executes without error") {
					assert.Equal(v, b.String(), "content as expected")
				}
			}
		}
	}

}

func Test_TemplatesFromData_BadTemplate(t *testing.T) {

	assert := assert.New(t)

	input := map[string]string{"foo": "broken: {{ nosuchfunctiondefined }}"}

	_, err := provider.TemplatesFromData(input)
	if assert.Error(err, "got error") {
		exp := "Template foo failed: template: foo:1: function " +
			"\"nosuchfunctiondefined\" not defined"
		assert.Equal(exp, err.Error(), "error as expected")
	}

}

func Test_PageTemplate_NilTemplate(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{"path": "/"})
	if err != nil {
		t.Fatal(err)
	}

	tmpl := provider.PageTemplate(nil, p)
	assert.Nil(tmpl, "nil template in, nil template out")
}

func Test_PageTemplate_TemplateInMeta_TitleCase(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/",
		"meta": map[string]interface{}{
			"Template": "foo/bar.html",
		}},
	)
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar.html": "HERE",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
	}
}

func Test_PageTemplate_TemplateInMeta_LowerCase(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/",
		"meta": map[string]interface{}{
			"template": "foo/bar.html",
		}},
	)
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar.html": "HERE",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
	}
}

func Test_PageTemplate_TemplateInMeta_UpperCase(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/",
		"meta": map[string]interface{}{
			"TEMPLATE": "foo/bar.html",
		}},
	)
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar.html": "HERE",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
	}
}

func Test_PageTemplate_PathMatch_Exact(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/foo/bar": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/foo/bar", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_PathMatch_NoSlash(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_PathMatch_WithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/foo/bar.html": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/foo/bar.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_PathMatch_NoSlashWithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar.html": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_BestGuess_Exact(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/foo/bar": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/foo/bar", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_BestGuess_NoSlash(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_BestGuess_WithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/foo/bar.html": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/foo/bar.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_BestGuess_NoSlashWithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"foo/bar.html": "HERE AT FOO BAR",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("foo/bar.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT FOO BAR", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_BestGuess_NoTopLevelSlash(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/": "HERE AT TOP",
	})
	tmpl := provider.PageTemplate(master, p)
	assert.Nil(tmpl, "got nothing")
}

func Test_PageTemplate_BestGuess_NoTopLevelDot(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		".": "HERE AT TOP",
	})
	tmpl := provider.PageTemplate(master, p)
	assert.Nil(tmpl, "got nothing")
}

func Test_PageTemplate_BestGuess_NoTopLevelEmpty(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"": "HERE AT TOP",
	})
	tmpl := provider.PageTemplate(master, p)
	assert.Nil(tmpl, "got nothing")
}

func Test_PageTemplate_Default_Exact(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/default": "HERE AT DEFAULT",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/default", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT DEFAULT", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_Default_NoSlash(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"default": "HERE AT DEFAULT",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("default", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT DEFAULT", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_Default_WithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"/default.html": "HERE AT DEFAULT",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("/default.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT DEFAULT", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_Default_NoSlashWithExtension(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/foo/bar/baz/bat",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"default.html": "HERE AT DEFAULT",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("default.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT DEFAULT", b.String(), "content as expected")
		}
	}
}

func Test_PageTemplate_Default_TopOfSite(t *testing.T) {

	assert := assert.New(t)

	p, err := provider.StandardPageFromData(map[string]interface{}{
		"path": "/",
	})
	if err != nil {
		t.Fatal(err)
	}

	master, _ := provider.TemplatesFromData(map[string]string{
		"default.html": "HERE AT DEFAULT",
	})
	tmpl := provider.PageTemplate(master, p)
	if assert.NotNil(tmpl, "got template") {
		assert.Equal("default.html", tmpl.Name(), "right template returned")
		var b bytes.Buffer
		if assert.Nil(tmpl.Execute(&b, nil), "executes without error") {
			assert.Equal("HERE AT DEFAULT", b.String(), "content as expected")
		}
	}
}
