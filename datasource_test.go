// datasource_test.go

package kisipar_test

import (
	// Standard:
	"errors"
	"os"
	"testing"
	"time"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_InterfaceConformity(t *testing.T) {

	// This will crash if anything doesn't match.
	var f = func(ds kisipar.DataSource) {
		t.Log(ds)
	}
	f(&kisipar.VirtualDataSource{})

}

func Test_IsNotExist(t *testing.T) {

	assert := assert.New(t)

	assert.True(kisipar.IsNotExist(kisipar.ErrNotExist), "local ErrNotExist")
	assert.True(kisipar.IsNotExist(os.ErrNotExist), "os.ErrNotExist")
	assert.False(kisipar.IsNotExist(errors.New("other")), "other error")

}

func Test_NewStandardPage(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"foo-id",               // id
		"The Foo",              // title
		[]string{"boo", "hoo"}, // tags
		time.Unix(0, 0),        // created
		time.Unix(10000, 0),    // updated
		map[string]interface{}{"helo": "WORLD"},
	)

	assert.Equal("foo-id", p.Id(), "Id")
	assert.Equal("The Foo", p.Title(), "Title")
	assert.Equal([]string{"boo", "hoo"}, p.Tags(), "Tags")
	assert.Equal(time.Unix(0, 0), p.Created(), "Created")
	assert.Equal(time.Unix(10000, 0), p.Updated(), "Updated")
	//	assert.Equal(map[string]interface{}{"helo": "WORLD"}, p.Meta(), "Meta")

}
