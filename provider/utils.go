// utils.go -- kisipar Provider utility functions and types.
// --------
// NOTE: these may wander off into the "utli" package at some point.

package provider

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// PathStrings is a sortable set of paths, e.g. the paths of the Pather items
// in a Provider.  Its first sort key is the depth of the path, i.e. the
// number of separators it contains; the second is the string itself.
type PathStrings []string

// Len returns the length, as per sort.Interface.
func (ps PathStrings) Len() int {
	return len(ps)
}

// Less reports whether i comes before j, as per sort.Interface.
func (ps PathStrings) Less(i, j int) bool {

	return PathLess(ps[i], ps[j])

}

// PathLess reports whether path a should be sorted before path b.
func PathLess(a, b string) bool {

	if a == b {
		return false
	}

	// First, compare the path depths.
	ci := strings.Count(a, "/")
	cj := strings.Count(b, "/")
	if ci != cj {
		return ci < cj
	}

	// Fall back to string comparison.
	return a < b

}

// Swap swaps two items, as per sort.Interface.
func (ps PathStrings) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// Add adds an item to the PathStrings if it does not already exist, returning
// the new PathStrings.  The PathStrings must already be sorted; populating
// the PathStrings via Add is one way to guarantee this.
func (ps PathStrings) Add(s string) PathStrings {

	pos := sort.Search(len(ps), func(i int) bool {
		return !PathLess(ps[i], s)
	})
	if pos == len(ps) || ps[pos] != s {

		// In the best case we are just manipulating slices here; in the
		// worst we are extending the array.  We trust append to handle this
		// efficiently.

		// First we make sure there's space for the new one:
		ps = append(ps, "")

		// Then we copy (dst,src) the elements above the insert
		// point one position higher.
		copy(ps[pos+1:], ps[pos:])

		// Finally we put our new element in at the position.
		ps[pos] = s
	}

	return ps
}

// Remove removes an item from the PathStrings if it exists, returning the new
// PathStrings.
func (ps PathStrings) Remove(s string) PathStrings {
	pos := sort.Search(len(ps), func(i int) bool {
		return !PathLess(ps[i], s)
	})
	if pos < len(ps) && ps[pos] == s {
		ps = append(ps[:pos], ps[pos+1:]...)
	}

	return ps

}

// FlexMappedValue returns the raw value in map m for key k, trying variations
// on character case when nil values are found. In order, the variations tried
// are: Title Case, UPPERCASE, and lowercase.
//
// This is useful in handling user input for things like meta blocks in YAML
// files.
func FlexMappedValue(m map[string]interface{}, k string) interface{} {

	if v := m[k]; v != nil {
		return v
	}
	if v := m[strings.Title(k)]; v != nil {
		return v
	}
	if v := m[strings.ToUpper(k)]; v != nil {
		return v
	}
	if v := m[strings.ToLower(k)]; v != nil {
		return v
	}

	return nil

}

// MappedString returns the string representation of the map value m for key
// string k.
//
// If the value is nil, an empty string is returned.
//
// If the value is a string, that is returned.  Otherwise the %v
// representation of the value is returned.
func MappedString(m map[string]interface{}, k string) string {

	if val := m[k]; val != nil {
		return stringVal(val)
	}
	return ""
}

// val --> string
func stringVal(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}

// val --> []string
func stringVals(val interface{}) []string {
	if val == nil {
		return []string{}
	}
	if s, ok := val.([]string); ok {
		return s
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:

		slice := reflect.ValueOf(val)
		s := make([]string, slice.Len())

		for i := 0; i < slice.Len(); i++ {
			s[i] = stringVal(slice.Index(i))
		}
		return s
	default:
		return []string{stringVal(val)}
	}
}

// FlexMappedString returns the string representation of the map value m for
// string k as in MappedString with case variations from FlexMappedValue.
func FlexMappedString(m map[string]interface{}, k string) string {

	return stringVal(FlexMappedValue(m, k))

}

// MappedStrings returns a slice of strings corresponding to value of k in m.
// If the value is already a []string, that is returned; if it is a slice
// then each value is stringified as in MappedString and that slice of
// strings returned; if it is a single value, that value is stringified by
// MappedString and returned in a slice of one; and if the value is nil (or
// the map itself is nil) an empty slice is returned.
func MappedStrings(m map[string]interface{}, k string) []string {

	return stringVals(m[k])
}

// FlexMappedStrings returns a slice of strings as in MappedStrings, with
// case variations from FlexMappedValue.
func FlexMappedStrings(m map[string]interface{}, k string) []string {

	return stringVals(FlexMappedValue(m, k))

}
