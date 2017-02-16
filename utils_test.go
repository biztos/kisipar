// utils_test.go -- tests for kisipar general utilities
// -------------

package kisipar_test

import (
	// Standard:
	"bytes"
	"fmt"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_MappedString_NilMap(t *testing.T) {

	assert := assert.New(t)

	assert.Zero(kisipar.MappedString(nil, "x"), "empty string returned")
}

func Test_MappedString_NotInMap(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": "foo"}
	assert.Zero(kisipar.MappedString(m, "y"), "empty string returned")
}

func Test_MappedString_String(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": "foo"}
	assert.Equal("foo", kisipar.MappedString(m, "x"), "string as expected")
}

func Test_MappedString_Int(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": 12345}
	assert.Equal("12345", kisipar.MappedString(m, "x"), "string as expected")
}

func Test_MappedString_Stringer(t *testing.T) {

	assert := assert.New(t)

	// We happen to know this is a Stringer:
	var b bytes.Buffer
	fmt.Fprint(&b, "here")
	m := map[string]interface{}{"x": &b}
	assert.Equal("here", kisipar.MappedString(m, "x"), "string as expected")
}

func Test_FlexMappedString_NoMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": ""}
	assert.Zero(kisipar.FlexMappedString(m, "y"), "empty string returned")

}

func Test_FlexMappedString_ExactMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "x"), "string as expected")

}

func Test_FlexMappedString_TitleMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"Foo": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "foo"), "string as expected")

}

func Test_FlexMappedString_UpperMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"FOO": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "foo"), "string as expected")

}

func Test_FlexMappedString_LowerMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"foo": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "FOO"), "string as expected")

}
