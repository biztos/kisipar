// datasource_test.go

package kisipar_test

import (
	// Standard:
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func TestInterfaceConformity(t *testing.T) {

	// This will crash if anything doesn't match.
	var f = func(ds kisipar.DataSource) {
		t.Log(ds)
	}
	f(&kisipar.FileDataSource())

}

func TestFoo(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, 1, "placeholder")
}
