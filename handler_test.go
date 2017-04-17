// handler_test.go

package kisipar_test

import (
	// Standard:
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_NewHandler(t *testing.T) {

	assert := assert.New(t)

	h := kisipar.NewHandler(nil)
	assert.NotNil(h)

}
