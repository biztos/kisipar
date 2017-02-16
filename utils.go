// utils.go -- general-purpose utilities useful in Kisipar.
// --------
// NOTE: these may wander off into the "utli" package at some point.

package kisipar

import (
	"fmt"
	"reflect"
	"strings"
)

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

func stringVal(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}

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

// FlexMapedStrings returns a slice of strings as in MappedStrings, with
// case variations from FlexMappedValue.
func FlexMapedStrings(m map[string]interface{}, k string) []string {

	return stringVals(FlexMappedValue(m, k))

}
