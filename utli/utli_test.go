// utli_test.go - misc utility tests.
// ------------

package utli_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/kisipar/utli"
)

func Test_ParseTimeString(t *testing.T) {

	assert := assert.New(t)

	// We only care about the mechanism.
	orig := utli.TIME_PARSING_FORMAT_STRINGS
	defer func() { utli.TIME_PARSING_FORMAT_STRINGS = orig }()

	utli.TIME_PARSING_FORMAT_STRINGS = []string{
		"2006......02.......04",
		"2006.02.04",
		"2006!!!02!!!04",
	}

	// Success
	parsed := utli.ParseTimeString("2015.07.14")
	if assert.NotNil(parsed, "time returned on good parse") {
		assert.Equal(2015, parsed.Year(), "year correct")
	}

	// Failure (no format matches)
	notparsed := utli.ParseTimeString("2015/07/14")
	assert.Nil(notparsed, "nil returned on failed parse")

	// Bonus: we can parse our own default stringification by default!
	utli.TIME_PARSING_FORMAT_STRINGS = orig
	someday := time.Unix(123, 456)
	defparsed := utli.ParseTimeString(someday.String())
	exp := someday.String()
	got := defparsed.String()
	assert.Equal(exp, got, "default string parsed correctly")
}
