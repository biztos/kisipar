// main_test.go - attempts at testing the kisipar command
// ------------
// Or, an attempt at getting proper test coverage by overriding os.Exit.
//
// NOTE: not using assert here as this may turn into a blog post; stdlib only!

package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/biztos/kisipar"
)

type testRecorder struct {
	ExitCode  int
	LogBuffer *bytes.Buffer
	T         *testing.T
}

func (r *testRecorder) AssertExitedWith(code int) {
	if r.ExitCode == -1 {
		r.T.Errorf("Did not (apparently) exit; expected %d.", code)
	} else if r.ExitCode != code {
		r.T.Errorf("Exited with wrong code: expected %d, got %d",
			code, r.ExitCode)
	} else {
		r.T.Logf("Exited with expected code: %d", code)
	}
}

func (r *testRecorder) AssertLoggedRegexp(rs string) {
	re := regexp.MustCompile(rs)
	got := r.LogBuffer.String()
	if re.MatchString(got) {
		r.T.Logf("Log buffer correct:\n%s\nMatches: %s", got, rs)
	} else {
		r.T.Errorf("Log buffer incorrect.\n%s\nDoes not match: %s", got, rs)
	}
}

func (r *testRecorder) AssertLoggedString(s string) {
	got := r.LogBuffer.String()
	if got != s {
		r.T.Errorf("Log buffer incorrect.\nExp: '%s'\nGot: '%s'", s, got)
	} else {
		r.T.Logf("Log buffer correct: %s", s)
	}
}

func (r *testRecorder) Exit(code int) {
	r.ExitCode = code
}

func prepTestRecorder(t *testing.T) *testRecorder {
	r := &testRecorder{
		ExitCode:  -1,
		LogBuffer: new(bytes.Buffer),
		T:         t,
	}
	kisipar.LAUNCH_SERVERS = false // pending a better idea...
	log.SetOutput(r.LogBuffer)
	log.SetFlags(0)
	EXIT_FUNCTION = r.Exit

	return r
}

func Test_BadPath(t *testing.T) {

	r := prepTestRecorder(t)

	os.Args = []string{
		"kisipar",                     // the binary (ignored here)
		"no-such-path-here-we-assume", // our test site path (nonexistent)
	}
	main()

	r.AssertExitedWith(1)
	r.AssertLoggedRegexp("^Site error.*no such file or directory")

}

func Test_Success(t *testing.T) {

	r := prepTestRecorder(t)

	// We have a minimal test site handy:
	path := filepath.Join("test_data", "site_1")

	// The call is thus:
	os.Args = []string{
		"kisipar", // the binary (ignored here)
		path,      // our test site path
	}
	main()

	r.AssertExitedWith(0)
	r.AssertLoggedString("Test Site One: listening on port 8081.\n")
}
