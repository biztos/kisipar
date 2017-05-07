// handler_test.go

package kisipar_test

import (
	// Standard:
	"bytes"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	// Helpful:
	"github.com/biztos/testig"
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_NewHandler_ErrorNoSite(t *testing.T) {

	panicky := func() { kisipar.NewHandler(nil) }
	testig.AssertPanicsWith(t, panicky, "Site must not be nil.",
		"Nil Site caught")

}

func Test_NewHandler_ErrorNoProvider(t *testing.T) {

	panicky := func() { kisipar.NewHandler(&kisipar.Site{}) }
	testig.AssertPanicsWith(t, panicky, "Site.Provider must not be nil.",
		"Nil Provider caught")

}

func Test_NewHandler_Success(t *testing.T) {

	assert := assert.New(t)

	s := &kisipar.Site{Provider: &kisipar.StandardProvider{}}
	_, err := kisipar.NewHandler(s)
	if !assert.Nil(err, "no error") {
		assert.FailNow(err.Error())
	}

}

func Test_Handler_Error_Fallback(t *testing.T) {

	assert := assert.New(t)

	// We need a Site with a Provider that has no error templates.
	s := &kisipar.Site{
		Config:   &kisipar.Config{Port: 1000},
		Provider: kisipar.NewStandardProvider(),
	}
	h, err := kisipar.NewHandler(s)
	if !assert.Nil(err, "no error") {
		assert.FailNow(err.Error())
	}

	// Rig up test:
	// TODO: logger!
	r := httptest.NewRequest("GET", "http://biztos.com/foo", nil)
	w := httptest.NewRecorder()
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	// Test it with a fake code as well, just because:
	h.Error(w, r, 599, "oops", "badness")

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(599, resp.StatusCode, "code passed through")
	assert.Equal("text/plain; charset=utf-8", resp.Header.Get("Content-Type"),
		"text/plain served")
	exp := "oops\n"
	assert.Equal(exp, string(body), "body as expected")
	assert.Regexp("^.* GET http://biztos.com/foo 599 oops: badness", buf.String(), "error logged as expected")
}

func Test_Handler_Error_WithTemplate(t *testing.T) {

	assert := assert.New(t)

	// We need a Site with a Provider that has an error template.
	yaml := `# yaml
templates:
    /errors/599.html: |
        Path: {{ .Page.Path }}
        Error: {{ .Page.Title }}
        Detail: {{ .Page.HTML }}`
	sp, err := kisipar.StandardProviderFromYAML(yaml)
	if err != nil {
		assert.FailNow(err.Error())
	}

	s := &kisipar.Site{
		Config:   &kisipar.Config{Port: 1000},
		Provider: sp,
	}
	h, err := kisipar.NewHandler(s)
	if !assert.Nil(err, "no error") {
		assert.FailNow(err.Error())
	}

	// Rig up test:
	// TODO: logger!
	r := httptest.NewRequest("GET", "http://biztos.com/foo", nil)
	w := httptest.NewRecorder()
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	// Test it with a fake code as well, just because:
	h.Error(w, r, 599, "oops", "badness")

	// NOTE: We do not check the headers because (apparently) the recorder
	// doesn't automatically set them.
	exp := "Path: /errors/599\nError: oops\nDetail: badness"
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(599, resp.StatusCode, "code passed through")
	assert.Equal(exp, string(body), "body as expected")
	assert.Regexp("^.* GET http://biztos.com/foo 599 oops: badness", buf.String(), "error logged as expected")
}
