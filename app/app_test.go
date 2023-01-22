// app_test.go - tests for the Kisipar app package.
// -----------
// NOTE: since we do main()-like tricks here, some things are tested in
// app_examples_test.go; a Better Way (mock main?) would be nice!

package app_test

import (
	// Standard Library:
	"fmt"
	"os"
	"testing"

	// Third Party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/app"
)

// TODO: consider a "testy" library! Or at least a sub-package here!
// TODO: ...and remove the "assert" dep if so.
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

func Test_Usage(t *testing.T) {

	assert := assert.New(t)
	exp := `Foo Bar.

Usage:
  foobar [options] <SITEPATH>...
  foobar -h | --help
  foobar -v | --version

Options:
  -h --help     Show this screen.
  -v --version  Show version.

Version:
  This is Foo Bar version 3.2.1.
`
	assert.Equal(exp, app.Usage("Foo Bar", "3.2.1", "foobar"),
		"DocOpty-usage as expected")
}

func Test_GetOpts_EmptyStringsShouldPanic(t *testing.T) {

	AssertPanicsWith(t, func() { app.GetOpts("", "") },
		"Heading and Usage strings must not be empty.",
		"panics with both strings empty")

	AssertPanicsWith(t, func() { app.GetOpts("", "xxx") },
		"Heading and Usage strings must not be empty.",
		"panics with heading empty")

	AssertPanicsWith(t, func() { app.GetOpts("xxx", "") },
		"Heading and Usage strings must not be empty.",
		"panics with usage empty")
}

func Test_GetOpts_SingleSitePath(t *testing.T) {

	// Just in case you have an app that wants to limit usage to a single
	// site, but still use the standard interfaces.
	assert := assert.New(t)

	usage := `XXX

Usage:
  xxx [options] <SITEPATH>
  xxx -h | --help
  xxx -v | --version

Options:
  -h --help     Show this screen.
  -v --version  Show version.
`

	os.Args = []string{"xxx", "dummypath"}
	opts := app.GetOpts("xxx", usage)
	assert.Equal([]string{"dummypath"}, opts.SitePaths, "single path parsed")
}

func Test_GetOpts_MultiSitePath(t *testing.T) {

	assert := assert.New(t)

	usage := `XXX

Usage:
  xxx [options] <SITEPATH>...
  xxx -h | --help
  xxx -v | --version

Options:
  -h --help     Show this screen.
  -v --version  Show version.
`

	os.Args = []string{"xxx", "path1", "path2", "path3"}
	opts := app.GetOpts("xxx", usage)
	assert.Equal([]string{"path1", "path2", "path3"}, opts.SitePaths,
		"multi path parsed")
}

func Test_Run_PathError(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"xxx", "no-such-path-here-we-hope"}
	err := app.Run("XXX", "1.2.3", "xxx")
	if assert.Error(err, "error returned") {
		assert.Regexp("no such file or directory", err.Error(),
			"error useful")
	}

}
