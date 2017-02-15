// binsanity_test.go - auto-generated; edit at thine own peril!
//
// More info: https://github.com/biztos/binsanity

package kisipar_test

import (
	"fmt"
	"testing"

	"github.com/biztos/kisipar"
)

func TestAssetNames(t *testing.T) {
	names := kisipar.AssetNames()
	t.Log(names)
}

func TestAsset(t *testing.T) {

	// Not found:
	missing := "---* no such asset we certainly hope *---"
	_, err := kisipar.Asset(missing)
	if err == nil {
		t.Fatal("No error for missing asset.")
	}
	if err.Error() != "Asset "+missing+" not found" {
		t.Fatal("Wrong error for missing asset: ", err.Error())
	}

	// Found (each one):
	for _, name := range kisipar.AssetNames() {
		// NOTE: it would be nice to test the non-zero sizes but it's possible
		// to have empty files, so...
		_, err := kisipar.Asset(name)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestMustAsset(t *testing.T) {

	// Not found:
	missing := "---* no such asset we certainly hope *---"
	exp := "Asset ---* no such asset we certainly hope *--- not found"
	panicky := func() { kisipar.MustAsset(missing) }
	AssertPanicsWith(t, panicky, exp, "panic for not found")

	// Found (each one):
	for _, name := range kisipar.AssetNames() {
		// NOTE: it would be nice to test the non-zero sizes but it's possible
		// to have empty files, so...
		_ = kisipar.MustAsset(name)
	}
}

// For a more useful version of this see: https://github.com/biztos/testig
func AssertPanicsWith(t *testing.T, f func(), exp string, msg string) {

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
		t.Fatalf("Function did not panic: %s", msg)
	} else if got != exp {

		t.Fatalf("Panic not as expected: %s\n  expected: %s\n    actual: %s",
			msg, exp, got)
	}

	// (In go testing, success is silent.)

}
