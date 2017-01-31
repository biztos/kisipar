// datasource_test.go

package kisipar_test

import (
	// Standard:
	"errors"
	"fmt"
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
	f(&kisipar.StandardDataSource{})

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
		"foo-id",                                // id
		"The Foo",                               // title
		[]string{"boo", "hoo"},                  // tags
		time.Unix(0, 0),                         // created
		time.Unix(10000, 0),                     // updated
		map[string]interface{}{"helo": "WORLD"}, // meta
	)

	assert.Equal("foo-id", p.Id(), "Id")
	assert.Equal("The Foo", p.Title(), "Title")
	assert.Equal([]string{"boo", "hoo"}, p.Tags(), "Tags")
	assert.Equal(time.Unix(0, 0), p.Created(), "Created")
	assert.Equal(time.Unix(10000, 0), p.Updated(), "Updated")
	assert.Equal(map[string]interface{}{"helo": "WORLD"}, p.Meta(), "Meta")

}

func Test_StandardPageFromData(t *testing.T) {

	assert := assert.New(t)

	cr := time.Unix(0, 0)
	up := time.Now()
	input := map[string]interface{}{
		"id":             "possibly-unique",
		"title":          "Hello World",
		"tags":           []string{"foo", "bar"},
		"created":        cr,
		"updated":        up,
		"meta":           map[string]interface{}{"foo": "bar"},
		"something else": &kisipar.StandardPage{},
	}

	p, err := kisipar.StandardPageFromData(input)
	if assert.Nil(err, "no error") {
		assert.Equal("possibly-unique", p.Id(), "Id")
		assert.Equal("Hello World", p.Title(), "Title")
		assert.Equal([]string{"foo", "bar"}, p.Tags(), "Tags")
		assert.Equal(cr, p.Created(), "Created")
		assert.Equal(up, p.Updated(), "Updated")
		assert.Equal(map[string]interface{}{"foo": "bar"}, p.Meta(), "Meta")
	}

}

func Test_StandardPageFromData_TypeErrors(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{
		"id":      "possibly-unique",
		"title":   "Hello World",
		"tags":    []string{"foo", "bar"},
		"created": time.Time{},
		"updated": time.Time{},
		"meta":    map[string]interface{}{"foo": "bar"},
	}

	tStr := map[string]string{
		"id":      "string",
		"title":   "string",
		"tags":    "string slice",
		"created": "Time",
		"updated": "Time",
		"meta":    "string-interface map",
	}
	type NotMyType struct{}
	for k, _ := range input {
		// nice dumb copy
		in2 := map[string]interface{}{}
		for k2, v2 := range input {
			in2[k2] = v2
		}

		// now override the one thing
		in2[k] = NotMyType{}

		_, err := kisipar.StandardPageFromData(in2)
		if assert.Error(err) {

			exp := fmt.Sprintf("Wrong type for %s: kisipar_test.NotMyType not %s",
				k, tStr[k])
			assert.Equal(exp, err.Error(), "error as expected")
		} else {
			t.Fatalf("No error but we reset %s!", k)
		}
	}

}

func Test_StandardPageFromData_EmptyData(t *testing.T) {

	assert := assert.New(t)

	input := map[string]interface{}{}

	p, err := kisipar.StandardPageFromData(input)
	if assert.Nil(err, "no error") {
		assert.Zero("", p.Id(), "Id")
		assert.Zero("", p.Title(), "Title")
		assert.Zero(p.Tags(), "Tags")
		assert.Zero(p.Created(), "Created")
		assert.Zero(p.Updated(), "Updated")
		assert.Zero(p.Meta(), "Meta")
	}

}

func Test_StandardPage_MetaString(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"foo-id",               // id
		"The Foo",              // title
		[]string{"boo", "hoo"}, // tags
		time.Unix(0, 0),        // created
		time.Unix(10000, 0),    // updated

		// meta:
		map[string]interface{}{
			"helo": "WORLD",
			"pi":   3.1415,
			"nada": nil,
			"ts":   time.Unix(0, 0),
		},
	)

	expTs := time.Unix(0, 0).String() // includes local TS
	assert.Equal("WORLD", p.MetaString("helo"), "string -> string")
	assert.Equal("3.1415", p.MetaString("pi"), "float -> string")
	assert.Equal("", p.MetaString("nada"), "real nil -> string")
	assert.Equal("", p.MetaString("nonesuch"), "missing -> string")
	assert.Equal(expTs, p.MetaString("ts"), "time -> string")

	// nil source
	p = &kisipar.StandardPage{}
	assert.Equal("", p.MetaString("any"), "nil meta -> empty string")
}

func Test_StandardPage_MetaStrings(t *testing.T) {

	assert := assert.New(t)

	p := kisipar.NewStandardPage(
		"foo-id",               // id
		"The Foo",              // title
		[]string{"boo", "hoo"}, // tags
		time.Unix(0, 0),        // created
		time.Unix(10000, 0),    // updated

		// meta:
		map[string]interface{}{
			"strings": []string{"fee", "fi", "fo"},
			"single":  "HOLA",
			"nada":    nil,
			"ts":      time.Unix(0, 0),
			"mixed":   []interface{}{time.Unix(0, 0), 3.1415, "flubber"},
		},
	)

	expTs := time.Unix(0, 0).String() // includes local TS
	assert.Equal([]string{"fee", "fi", "fo"}, p.MetaStrings("strings"), "strings")
	assert.Equal([]string{"HOLA"}, p.MetaStrings("single"), "one string")

	assert.Equal([]string{}, p.MetaStrings("nada"), "real nil")
	assert.Equal([]string{}, p.MetaStrings("nonesuch"), "missing")

	assert.Equal([]string{expTs}, p.MetaStrings("ts"), "single time")

	assert.Equal([]string{expTs, "3.1415", "flubber"},
		p.MetaStrings("mixed"), "mixed slice")

	// nil source
	p = &kisipar.StandardPage{}
	assert.Equal([]string{}, p.MetaStrings("any"),
		"nil meta -> empty string slice")
}
