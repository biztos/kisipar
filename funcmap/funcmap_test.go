// funcmap_test.go - test the template-accessible functions.
// ---------------

package funcmap_test

import (
	// Standard:
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/funcmap"
)

func Test_New(t *testing.T) {

	assert := assert.New(t)

	fm := funcmap.New()

	exp_funcs := []string{
		// products of our own tortured logic:
		"truncints",
		"truncintsto",
		"sortints",
		"sortintsasc",
		"sortintsdesc",
		"reverseints",
		"intrange",
		"pathdepth",
		"indent",

		// gtf freebies:
		"replace",
		"default",
		"length",
		"lower",
		"upper",
		"truncatechars",
		"urlencode",
		"wordcount",
		"divisibleby",
		"lengthis",
		"trim",
		"capfirst",
		"pluralize",
		"yesno",
		"rjust",
		"ljust",
		"center",
		"filesizeformat",
		"apnumber",
		"intcomma",
		"ordinal",
		"first",
		"last",
		"join",
		"slice",
		"random",
		"striptags",
	}

	for _, s := range exp_funcs {

		assert.NotNil(fm[s], "func defined for "+s)
	}
}

func Test_TruncInts(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{}, funcmap.TruncInts([]int{}),
		"empty truncates to empty slice")

	assert.Equal([]int{}, funcmap.TruncInts([]int{4}),
		"single truncates to empty slice")

	assert.Equal([]int{3}, funcmap.TruncInts([]int{3, 2}),
		"two elements truncate to one")

	// This *should* be more efficient than dealing with array copies.
	i := []int{22, 33, 44}
	r := funcmap.TruncInts(i)
	i[1] = 999
	assert.Equal(999, r[1], "slice is of same underlying array")

}

func Test_TruncIntsTo(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{}, funcmap.TruncIntsTo(1, []int{}),
		"empty truncates to empty slice")

	assert.Equal([]int{4}, funcmap.TruncIntsTo(1, []int{4}),
		"single truncates to single at 1")

	assert.Equal([]int{}, funcmap.TruncIntsTo(0, []int{4}),
		"single truncates to nothing at 0")

	assert.Equal([]int{3}, funcmap.TruncIntsTo(1, []int{3, 2}),
		"two elements truncate to one at 1")

	assert.Equal([]int{3, 2}, funcmap.TruncIntsTo(2, []int{3, 2, 1}),
		"three elements truncate to two at 2")

	assert.Equal([]int{3, 2}, funcmap.TruncIntsTo(-1, []int{3, 2, 1}),
		"three elements truncate to two at -1")

	assert.Equal([]int{3}, funcmap.TruncIntsTo(-2, []int{3, 2, 1}),
		"three elements truncate to one at -2")

	assert.Equal([]int{3}, funcmap.TruncIntsTo(-2, []int{3, 2, 1}),
		"three elements truncate to zero at -3")

	assert.Equal([]int{3}, funcmap.TruncIntsTo(-2, []int{3, 2, 1}),
		"three elements truncate to zero at -4")

	// This *should* be more efficient than dealing with array copies.
	// i := []int{22, 33, 44}
	// r := funcmap.TruncIntsTo(i, 2)
	// i[1] = 999
	// assert.Equal(999, r[1], "slice is of same underlying array")

}

func Test_ReverseInts(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{}, funcmap.ReverseInts([]int{}),
		"empty stays empty")

	assert.Equal([]int{1, 2, 3, 4}, funcmap.ReverseInts([]int{4, 3, 2, 1}),
		"ordered reverses")

	assert.Equal([]int{9, 1, 3, 1, 9}, funcmap.ReverseInts([]int{9, 1, 3, 1, 9}),
		"result is not sorted")

	// Prove we copy.
	i := []int{22, 33, 44}
	r := funcmap.ReverseInts(i)
	i[1] = 999
	assert.Equal(33, r[1], "return slice is a copy")

}

func Test_SortInts(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{}, funcmap.SortInts([]int{}),
		"empty stays empty")

	assert.Equal([]int{1, 2, 3, 4}, funcmap.SortInts([]int{4, 3, 2, 1}),
		"ordered reverses")

	assert.Equal([]int{1, 1, 3, 9, 9}, funcmap.SortInts([]int{9, 1, 3, 1, 9}),
		"result is sorted")

	// Prove we copy.
	i := []int{22, 33, 44}
	r := funcmap.SortInts(i)
	i[1] = 999
	assert.Equal(33, r[1], "return slice is a copy")

}

func Test_SortIntsDesc(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{}, funcmap.SortIntsDesc([]int{}),
		"empty stays empty")

	assert.Equal([]int{4, 3, 2, 1}, funcmap.SortIntsDesc([]int{1, 2, 3, 4}),
		"ordered reverses")

	assert.Equal([]int{9, 9, 3, 1, 1}, funcmap.SortIntsDesc([]int{9, 1, 3, 1, 9}),
		"result is sorted")

	// Prove we copy.
	i := []int{22, 33, 44}
	r := funcmap.SortIntsDesc(i)
	i[1] = 999
	assert.Equal(33, r[1], "return slice is a copy")

}

func Test_IntRange(t *testing.T) {

	assert := assert.New(t)

	assert.Equal([]int{4}, funcmap.IntRange(4, 4),
		"IntRange(4,4)-> 4")

	assert.Equal([]int{1, 2, 3, 4}, funcmap.IntRange(1, 4),
		"IntRange(1,4)-> 1 2 3 4")

	assert.Equal([]int{4, 3, 2, 1}, funcmap.IntRange(4, 1),
		"IntRange(4,1) -> 4 3 2 1")

	assert.Equal([]int{-2, -1, 0, 1}, funcmap.IntRange(-2, 1),
		"IntRange(-2, 1) -> -2 -1 0 1")

	assert.Equal([]int{0, -1, -2, -3}, funcmap.IntRange(0, -3),
		"IntRange(0, -3) -> 0 -1 -2 -3")

}

func Test_PathDepth(t *testing.T) {

	assert := assert.New(t)

	assert.Equal(0, funcmap.PathDepth("no slashes here"),
		"no slashes -> zero")

	assert.Equal(1, funcmap.PathDepth("/standard"),
		"leading slashe -> one")

	assert.Equal(4, funcmap.PathDepth("/fee/fi/foe/fum"),
		"multiple")
}

func Test_Indent(t *testing.T) {

	assert := assert.New(t)

	assert.Equal("", funcmap.Indent(0),
		"zero -> empty")

	assert.Equal("    ", funcmap.Indent(1),
		"one -> 4 sp.")

	assert.Equal("            ", funcmap.Indent(3),
		"three -> 4 sp.")
}
