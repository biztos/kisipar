// utils_test.go -- tests for kisipar general utilities
// -------------

package kisipar_test

import (
	// Standard:
	"bytes"
	"fmt"
	"sort"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_PathStrings_Sort(t *testing.T) {

	assert := assert.New(t)

	ps := kisipar.PathStrings{"/a/x", "/zzzzzz", "/a/b/c"}
	sort.Sort(ps)

	exp := kisipar.PathStrings{"/zzzzzz", "/a/x", "/a/b/c"}
	assert.Equal(exp, ps, "sorted as expected")

}

func Test_PathStrings_Add(t *testing.T) {

	assert := assert.New(t)

	ps := kisipar.PathStrings{}
	ps = ps.Add("/a/x")
	ps = ps.Add("/zzzzzz")
	ps = ps.Add("/a/b/c")
	ps = ps.Add("/aaa")     // at start
	ps = ps.Add("/a/b/c/d") // at end
	ps = ps.Add("/zzzzzz")  // exists

	exp := kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	assert.Equal(exp, ps, "added as expected")

}

func Test_PathStrings_Remove(t *testing.T) {

	assert := assert.New(t)

	ps := kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	exp := ps

	assert.Equal(exp, ps.Remove("/nonesuch"),
		"removal of missing item does nothing")

	ps = kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	exp = kisipar.PathStrings{"/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	assert.Equal(exp, ps.Remove("/aaa"), "removal of first item")

	ps = kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	exp = kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c"}
	assert.Equal(exp, ps.Remove("/a/b/c/d"), "removal of last item")

	ps = kisipar.PathStrings{"/aaa", "/zzzzzz", "/a/x", "/a/b/c", "/a/b/c/d"}
	exp = kisipar.PathStrings{"/aaa", "/a/x", "/a/b/c", "/a/b/c/d"}
	assert.Equal(exp, ps.Remove("/zzzzzz"), "removal of middle item")

	ps = kisipar.PathStrings{"/aaa"}
	exp = kisipar.PathStrings{}
	assert.Equal(exp, ps.Remove("/aaa"), "removal of only item")

	ps = kisipar.PathStrings{}
	exp = kisipar.PathStrings{}
	assert.Equal(exp, ps.Remove("/aaa"), "removal from empty")
}

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
	assert.Equal("here", kisipar.FlexMappedString(m, "x"),
		"string as expected")

}

func Test_FlexMappedString_TitleMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"Foo": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "foo"),
		"string as expected")

}

func Test_FlexMappedString_UpperMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"FOO": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "foo"),
		"string as expected")

}

func Test_FlexMappedString_LowerMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"foo": "here"}
	assert.Equal("here", kisipar.FlexMappedString(m, "FOO"),
		"string as expected")

}

func Test_MappedStrings_NilMap(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]string{}, kisipar.MappedStrings(nil, "x"),
		"empty slice returned")
}

func Test_MappedStrings_NotInMap(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": "foo"}
	assert.Equal([]string{}, kisipar.MappedStrings(m, "y"),
		"empty slice returned")
}

func Test_MappedStrings_StringSlice(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": []string{"a", "b"}}
	assert.Equal([]string{"a", "b"}, kisipar.MappedStrings(m, "x"),
		"slice as expected")
}

func Test_MappedStrings_IntSlice(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": []int{12, 34, 5}}
	assert.Equal([]string{"12", "34", "5"}, kisipar.MappedStrings(m, "x"),
		"slice as expected")
}

func Test_MappedStrings_StringerSlice(t *testing.T) {

	assert := assert.New(t)

	// We happen to know this is a Stringer:
	var b1 bytes.Buffer
	fmt.Fprint(&b1, "here")
	var b2 bytes.Buffer
	fmt.Fprint(&b2, "there")
	m := map[string]interface{}{"x": []interface{}{&b1, &b2}}
	assert.Equal([]string{"here", "there"}, kisipar.MappedStrings(m, "x"),
		"slice as expected")
}

func Test_MappedStrings_MixedSlice(t *testing.T) {

	assert := assert.New(t)

	var b1 bytes.Buffer
	fmt.Fprint(&b1, "here")
	m := map[string]interface{}{"x": []interface{}{&b1, "there", 3.14}}
	assert.Equal([]string{"here", "there", "3.14"},
		kisipar.MappedStrings(m, "x"),
		"slice as expected")
}

func Test_FlexMappedStrings_NoMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": ""}
	assert.Equal([]string{}, kisipar.FlexMappedStrings(m, "y"),
		"empty slice returned")

}

func Test_FlexMappedStrings_ExactMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"x": []string{"a", "b"}}
	assert.Equal([]string{"a", "b"}, kisipar.FlexMappedStrings(m, "x"),
		"slice as expected")

}

func Test_FlexMappedStrings_TitleMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"Foo": []string{"a", "b"}}
	assert.Equal([]string{"a", "b"}, kisipar.FlexMappedStrings(m, "foo"),
		"slice as expected")

}

func Test_FlexMappedStrings_UpperMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"FOO": []string{"a", "b"}}
	assert.Equal([]string{"a", "b"}, kisipar.FlexMappedStrings(m, "foo"),
		"slice as expected")

}

func Test_FlexMappedStrings_LowerMatch(t *testing.T) {

	assert := assert.New(t)

	m := map[string]interface{}{"foo": []string{"a", "b"}}
	assert.Equal([]string{"a", "b"}, kisipar.FlexMappedStrings(m, "FOO"),
		"slice as expected")

}
