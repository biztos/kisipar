// app.go - rigging for the kisipar binary app

// Package app defines the Kisipar web server application.
package app

import (

	// Standard Library:
	"fmt"

	// Third-party:
	"github.com/docopt/docopt-go"

	// Kisipar:
	"github.com/biztos/kisipar"
)

// Options represents the standard set of command-line options, which are
// parsed according to the usage specification provided to GetOpts.
type Options struct {
	SitePaths []string
}

// Run runs the application with the standard options and the provided name,
// version, and binary name; and the default usage spec from Usage.  The first
// error encountered is returned.
//
//  func main() {
//      if err := kisipar.Run("Foobar Thingy","1.2.3","foobar"); err != nil
//         log.Fatal(err)
//      }
//      log.Println("DONE.")
//  }
func Run(name, version, binary string) error {

	opts := GetOpts(name+" "+version, Usage(name, version, binary))
	k, err := kisipar.Load(opts.SitePaths...)
	if err != nil {
		return err
	}
	k.Serve()
	return nil
}

// GetOpts processes arguments, exiting if the result is not ready to serve.
// The usage argument must be a DocOpt-style specification (and help text),
// and the heading argument must be the application name including version,
// e.g. "Kisipar 1.0.0" or "Naval Fate 2.0".
func GetOpts(heading, usage string) *Options {
	if heading == "" || usage == "" {
		panic("Heading and Usage strings must not be empty.")
	}
	args, _ := docopt.Parse(
		usage,
		nil,     // use default os args
		true,    // enable help option
		heading, // the version string
		false,   // do NOT require options first
		true,    // let DocOpt exit for version, help, user error
	)

	opts := &Options{}

	// We *should* have a SITEPATH arg set one way or the other.
	if sp := args["<SITEPATH>"]; sp != nil {
		// It *might* be a single string.
		if s, ok := sp.(string); ok {
			opts.SitePaths = []string{s}
		} else if ss, ok := sp.([]string); ok {
			// More commonly, it's a (small) set of them.
			opts.SitePaths = ss
		}
	}

	// Further opts are TODO... MVP first!

	return opts
}

// Usage returns the standard DocOpt-style usage specification with the given
// name as the proper app name with version.
func Usage(name, version, binary string) string {

	f := `%s.

Usage:
  %s [options] <SITEPATH>...
  %s -h | --help
  %s -v | --version

Options:
  -h --help     Show this screen.
  -v --version  Show version.

Version:
  This is %s version %s.
`

	return fmt.Sprintf(f,
		name,    // heading
		binary,  // usage
		binary,  // usage: help
		binary,  // usage: version
		name,    // version
		version, // version
	)

}
