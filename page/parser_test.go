// page/page_parsing_test.go - tests for Page parsing features.
// -------------------------

package page_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/kisipar/page"
)

type TestParser struct{}
type TestParseResult struct {
	meta    map[string]interface{}
	content []byte
}

func (tpr *TestParseResult) Meta() map[string]interface{} {
	return tpr.meta
}
func (tpr *TestParseResult) Content() []byte {
	return tpr.content
}
func (t *TestParser) Parse([]byte) (page.ParseResult, error) {
	res := &TestParseResult{
		meta:    map[string]interface{}{"Title": "tested"},
		content: []byte("test parsed"),
	}
	return res, nil
}

func Test_VerbatimParser(t *testing.T) {

	assert := assert.New(t)

	input := `# Could be Markdown!

    { "id" : 1234 }

Might be, but we do not care.`

	v := &page.VerbatimParser{}
	res, err := v.Parse([]byte(input))
	assert.Nil(err, "no error on verbatim Parse")
	assert.Equal(map[string]interface{}{}, res.Meta(),
		"empty meta map returned")
	assert.Equal(input, string(res.Content()), "verbatim content returned")

}

func Test_Parse_ExtParserFound(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()
	page.ExtParsers = []*page.ExtParser{
		{".foo", &page.VerbatimParser{}},
		{".here", &TestParser{}},
	}
	p := &page.Page{Path: "/some/thing.here", Source: []byte("RAW")}

	err := p.Parse()
	assert.Nil(err, "no error parsing with the test parser")
	assert.Equal(p.Title(), "tested", "meta set by parser")
	assert.Equal("test parsed", string(p.Content), "content set by parser")

}

func Test_Parse_ExtParserCaseInsensitive(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()
	page.ExtParsers = []*page.ExtParser{
		{".HERE", &TestParser{}},
		{".here", &page.VerbatimParser{}},
		{".Here", &page.VerbatimParser{}},
	}
	p := &page.Page{Path: "/some/thing.Here", Source: []byte("RAW")}

	err := p.Parse()
	assert.Nil(err, "no error parsing with the test parser")
	assert.Equal(p.Title(), "tested", "meta set by parser")
	assert.Equal("test parsed", string(p.Content), "content set by parser")

}

func Test_Parse_ExtParserDefault(t *testing.T) {

	assert := assert.New(t)

	origDefaultParser := page.DefaultParser
	defer func() { page.DefaultParser = origDefaultParser }()
	page.DefaultParser = &TestParser{}
	p := &page.Page{Path: "/some/thing.here", Source: []byte("RAW")}

	err := p.Parse()
	assert.Nil(err, "no error parsing with the test parser as default")
	assert.Equal(p.Title(), "tested", "meta set by parser")
	assert.Equal("test parsed", string(p.Content), "content set by parser")

}

func Test_Parse_LimitExtParsers(t *testing.T) {

	assert := assert.New(t)

	origExtParsers := page.ExtParsers
	defer func() { page.ExtParsers = origExtParsers }()

	// Empty! (Would you ever do this?)
	exp := []*page.ExtParser{}
	page.LimitExtParsers([]string{})
	assert.Equal(exp, page.ExtParsers, "limit to empty -> empty")

	// Just .md
	page.ExtParsers = origExtParsers
	exp = []*page.ExtParser{origExtParsers[0]}
	page.LimitExtParsers([]string{".md"})
	assert.Equal(exp, page.ExtParsers, "limit to .md -> one")

	// A couple.
	page.ExtParsers = origExtParsers
	exp = []*page.ExtParser{origExtParsers[0], origExtParsers[4]}
	page.LimitExtParsers([]string{".md", ".txt"})
	assert.Equal(exp, page.ExtParsers, "limit to .md+.txt -> two")

}
