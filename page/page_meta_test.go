// page/page_meta_test.go - tests for the Meta map of a Page.
// ----------------------

package page_test

import (
	"fmt"
	"github.com/biztos/kisipar/page"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Anything struct {
	No int
}

func (a Anything) String() string {
	if a.No == 0 {
		return "I could be anything."
	}
	return fmt.Sprintf("Anything %d.", a.No)
}

type Anythings []Anything

func (a Anythings) StringArray() []string {
	res := make([]string, len(a))
	for i, e := range a {
		res[i] = e.String()
	}
	return res
}

type Stringless struct {
	Whatever string
}

func Test_MetaString_KeyNotFound(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{}} // empty meta

	assert.Equal("", p.MetaString("foo"), "empty string for no-such-key")

	// Same should happen if it's empty.
	p.Meta["foo"] = ""
	assert.Equal("", p.MetaString("foo"), "empty string for empty-string val")

	// Same should happen if it's not a standard case.
	p.Meta["Bar"] = ""
	assert.Equal("", p.MetaString("bar"), "empty string for weird key case")

}

func Test_MetaString_CaseVariants(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{
		"foo": "lower foo",
		"Foo": "title foo",
		"FOO": "upper foo",
		"BAR": "upper bar",
	}}

	assert.Equal("title foo", p.MetaString("Foo"), "exact match: title")
	assert.Equal("lower foo", p.MetaString("foo"), "exact match: lower")
	assert.Equal("upper foo", p.MetaString("FOO"), "exact match: upper")

	assert.Equal("lower foo", p.MetaString("foO"), "fallback: lower")
	assert.Equal("upper bar", p.MetaString("bar"), "fallback: upper")

}

func Test_MetaString_NonStringTypes(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{
		"int":        int(1),
		"uint":       uint(2),
		"int32":      int32(3),
		"int64":      int64(4),
		"uint32":     uint32(5),
		"uint64":     uint64(6),
		"float32":    float32(7.123),
		"float64":    float64(8.1234567),
		"true":       true,
		"false":      false,
		"time":       time.Unix(0, 0).UTC(),
		"anything":   &Anything{},
		"stringless": &Stringless{"thing"},
	}}

	assert.Equal("1", p.MetaString("int"), "int converts")
	assert.Equal("2", p.MetaString("uint"), "uint converts")
	assert.Equal("3", p.MetaString("int32"), "int32 converts")
	assert.Equal("4", p.MetaString("int64"), "int64 converts")
	assert.Equal("5", p.MetaString("uint32"), "uint32 converts")
	assert.Equal("6", p.MetaString("uint64"), "uint64 converts")
	assert.Equal("7.123", p.MetaString("float32"), "uint32 converts")
	assert.Equal("8.1234567", p.MetaString("float64"), "uint64 converts")
	assert.Equal("true", p.MetaString("true"), "true converts")
	assert.Equal("false", p.MetaString("false"), "false converts")

	// Times will likely be used a lot; we will probably want some conversion
	// options available in the templates.
	assert.Equal("1970-01-01 00:00:00 +0000 UTC", p.MetaString("time"),
		"time converts")

	// Anything that implements fmt.Stringer:
	assert.Equal("I could be anything.", p.MetaString("anything"),
		"random struct is stringified")

	// Something that does not, and thus passes through to fmt.Sprintf:
	assert.Equal("&{Whatever:thing}", p.MetaString("stringless"),
		"random struct is stringified")
}

func Test_MetaStringArray(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{
		"stringArray":    []string{"eenie", "meenie", "  foe  "},
		"plainString":    "eenie, meenie,   foe  ",
		"intArray":       []int{1, 33, 22},
		"int32Array":     []int32{1, 33, 22},
		"int64Array":     []int64{1, 33, 22},
		"uintArray":      []uint{1, 33, 22},
		"uint32Array":    []uint32{1, 33, 22},
		"uint64Array":    []uint64{1, 33, 22},
		"float32Array":   []float32{1.23, 3.45, 2.34},
		"float64Array":   []float64{1.2345678, 33.0000000001, 22.334455},
		"interfaceArray": []interface{}{"foo", "bar"},
		"StringArrayer": Anythings{
			Anything{No: 3},
			Anything{No: 1},
			Anything{No: 2},
		},
		"NotStringArrayer": Anything{No: 1},
	}}

	// not-found
	assert.Equal([]string{},
		p.MetaStringArray("nonesuch"),
		"not-found returns empty string array")

	// strings as-is
	assert.Equal([]string{"eenie", "meenie", "foe"},
		p.MetaStringArray("plainString"),
		"source plain string split and trimmed")

	// strings split
	assert.Equal([]string{"eenie", "meenie", "  foe  "},
		p.MetaStringArray("stringArray"),
		"source string array returned as is, unaltered")

	// ints
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("intArray"),
		"source int array stringified")
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("int32Array"),
		"source int32 array stringified")
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("int64Array"),
		"source int64 array stringified")

	// uints
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("uintArray"),
		"source uint array stringified")
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("uint32Array"),
		"source uint32 array stringified")
	assert.Equal([]string{"1", "33", "22"},
		p.MetaStringArray("uint64Array"),
		"source uint64 array stringified")

	// floats
	assert.Equal([]string{"1.23", "3.45", "2.34"},
		p.MetaStringArray("float32Array"),
		"source float32 array stringified")
	assert.Equal([]string{"1.2345678", "33.0000000001", "22.334455"},
		p.MetaStringArray("float64Array"),
		"source float64 array stringified")

	// stringifiable interface (as for tags, etc -- pretty common)
	assert.Equal([]string{"foo", "bar"},
		p.MetaStringArray("interfaceArray"),
		"source interface array stringified")

	// random struct array that implements StringArrayer
	assert.Equal([]string{"Anything 3.", "Anything 1.", "Anything 2."},
		p.MetaStringArray("StringArrayer"),
		"source array of random structs stringified via StringArray")

	// random struct that does not implement StringArrayer
	assert.Equal([]string{},
		p.MetaStringArray("NotStringArrayer"),
		"source random struct not stringified; no StringArray")
}

func Test_MetaTime(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{
		"ts_good":  "2016-05-08",
		"ts_bad":   "not-a-date-isit",
		"ts_empty": "",
	}}

	// not-found
	assert.Nil(p.MetaTime("nonesuch"), "nil for not-found")

	// empty
	assert.Nil(p.MetaTime("ts_empty"), "nil for empty string")

	// bad
	assert.Nil(p.MetaTime("ts_bad"), "nil for bad time string (not parsed)")

	// good
	ts := p.MetaTime("ts_good")
	if assert.NotNil(ts, "good time parsed") {
		assert.Equal(2016, ts.Year(), "year as expected")
		assert.Equal(5, int(ts.Month()), "month as expected")
		assert.Equal(8, ts.Day(), "day as expected")
	}

}

func Test_MetaBool(t *testing.T) {

	assert := assert.New(t)

	p := &page.Page{Meta: map[string]interface{}{
		"b_true":  true,
		"b_false": false,
		"b_bad":   "true",
		"b_nil":   nil,
		"asItIs":  true,  // case as passed
		"Asitis":  false, // not a default
		"asitis":  true,  // first default
		"ASITIS":  false, // second default
		"THIS":    true,  // second default (for use)
	}}

	// not-found
	assert.False(p.MetaBool("nonesuch"), "false for not-found")

	// nil
	assert.False(p.MetaBool("b_nil"), "false for nil")

	// bad
	assert.False(p.MetaBool("b_bad"), "false for non-bool")

	// true
	assert.True(p.MetaBool("b_true"), "true for bool true")

	// false
	assert.False(p.MetaBool("b_false"), "false for bool false")

	// prefer as given
	assert.True(p.MetaBool("asItIs"), "true for key-case-as-passed")

	// fallback to lc
	assert.True(p.MetaBool("AsItIs"), "true for lowercase fallback")

	// fallback to uc
	assert.True(p.MetaBool("this"), "true for uppercase fallback")

}
