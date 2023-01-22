// feed_test.go -- tests for the kisipar Atom feed generator.
// ------------

package site_test

import (
	// Standard:
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Kisipar:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/site"
)

func Test_Feed_SiteDefaults(t *testing.T) {

	assert := assert.New(t)

	s, err := site.LoadVirtual(
		nil, // config
		nil, // pages
		nil, // templates
	)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	assert.Equal("Anonymous Kisipar Site", f.Title)
	assert.Equal("http://localhost:8020/feed.xml", f.ID)
	if assert.NotNil(f.Author, "author set") {
		assert.Equal("Anonymous Kisipar Fan", f.Author.Name)
	}

}

func Test_Feed_EntryDefaults(t *testing.T) {

	assert := assert.New(t)

	p, err := page.LoadVirtualString("/foo.md", "There is not title.")
	if err != nil {
		t.Fatal(err)
	}
	s, err := site.LoadVirtual(
		nil,             // config
		[]*page.Page{p}, // pages
		nil,             // templates
	)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	if assert.Equal(1, len(f.Entry), "one entry") {
		e := f.Entry[0]
		assert.Equal("foo", e.Title)
		assert.Equal("http://localhost:8020/foo", e.ID)
		if assert.NotNil(e.Author, "author set") {
			assert.Equal("Anonymous Kisipar Fan", e.Author.Name)
		}
	}

}

func Test_Feed_Minimal(t *testing.T) {

	assert := assert.New(t)

	yaml := `# FEED TEST
Name: Feed Test Site from YAML
Pages:
    alpha.md: "# Not an index."
    beta.md: "# Also not."
    cronie/bad.md: "# Deeper non-index."
    cronie/index.md: "# An index!"
    dooder.md: "# Another non-index."
`

	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	if assert.Equal(5, len(f.Entry), "five entries (index included)") {
		exp_titles := []string{
			"Not an index.",
			"Also not.",
			"Another non-index.",
			// For ties by time, deeper is later and index is first
			"An index!",
			"Deeper non-index.",
		}
		titles := []string{
			f.Entry[0].Title,
			f.Entry[1].Title,
			f.Entry[2].Title,
			f.Entry[3].Title,
			f.Entry[4].Title,
		}
		assert.Equal(exp_titles, titles, "order as expected")
	}

}

func Test_Feed_TruncateToFeedItems(t *testing.T) {

	assert := assert.New(t)

	yaml := `# FEED TEST
Name: Feed Test Site from YAML
FeedItems: 3
Pages:
    a.md: "# One"
    b.md: "# Two"
    d.md: "# Three"
    e.md: "# Four"`

	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	assert.Equal(3, len(f.Entry), "FeedItems sets limit")

}

func Test_Feed_EntryIsPublished(t *testing.T) {

	assert := assert.New(t)

	yaml := `# FEED TEST
Name: Feed Test Site from YAML
Pages:
    a.md: "# One"
    b.md: |
        # Two
        
            Created: 1970-10-01
`

	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	if assert.Equal(2, len(f.Entry), "two items") {
		assert.Zero(f.Entry[0].Published, "first has blank Published attr")
		assert.Equal("1970-10-01T00:00:00+00:00", string(f.Entry[1].Published),
			"second has correct Published attr")

	}

}

func Test_Feed_EntryHasSummary(t *testing.T) {

	assert := assert.New(t)

	yaml := `# FEED TEST
Name: Feed Test Site from YAML
Pages:
    a.md: "# One"
    b.md: |
        # Two
        
            Summary: This here is my summary.
`

	s, err := site.LoadVirtualYaml(yaml)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Feed(nil)
	if assert.Equal(2, len(f.Entry), "two items") {
		assert.Zero(f.Entry[0].Published, "first has nil Summary attr")
		sum2 := f.Entry[1].Summary
		if assert.NotZero(sum2, "seocond has a Summary") {
			assert.Equal("text", sum2.Type, "Type is text")
			assert.Equal("This here is my summary.", sum2.Body,
				"Body set from page meta")
		}

	}

}
