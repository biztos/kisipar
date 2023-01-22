// serving_test.go -- tests for the kisipar site serving functions.
// ---------------

package site_test

import (
	// Standard:
	"errors"
	"fmt"
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/site"
)

// A fake Server! Because we aren't going to write unit tests for the http
// package here. :-)
type FakeServer struct {
	ServeTLS bool
}

func (s *FakeServer) ListenAndServe() error {
	return errors.New("INSECURE")
}
func (s *FakeServer) ListenAndServeTLS(c, k string) error {
	return fmt.Errorf("SECURE: %s & %s", c, k)
}

func Test_ServeWithoutServerPanics(t *testing.T) {

	s, err := site.New("") // all defaults, no TLS
	if err != nil {
		t.Fatal(err)
	}
	s.Server = nil

	AssertPanicsWith(t, func() { s.Serve() },
		"Serve called but Server is nil.",
		"Serve without Server panics as expected")

}

func Test_ServeWithoutTLS(t *testing.T) {

	assert := assert.New(t)

	s, err := site.New("") // all defaults, no TLS
	if err != nil {
		t.Fatal(err)
	}
	s.Server = &FakeServer{} // lest we wait forever...

	err = s.Serve()
	if assert.Error(err, "error returned") {
		assert.Equal("INSECURE", err.Error(), "error as expected")
	}

}

func Test_ServeWithTLS(t *testing.T) {

	assert := assert.New(t)

	s, err := site.New("") // all defaults, no TLS
	if err != nil {
		t.Fatal(err)
	}
	s.Server = &FakeServer{}
	s.CertFile = "ima_cert"
	s.KeyFile = "ima_key"

	s.ServeTLS = true

	err = s.Serve()
	if assert.Error(err, "error returned") {
		assert.Equal("SECURE: ima_cert & ima_key", err.Error(),
			"error as expected")
	}

}
