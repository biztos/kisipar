// feed.go - news feed generation for the Kisipar site.
// -------
// TODO: deal with remote content etc, obviously this is going to require
// some kind of re-parse no?
//
// USEFUL: https://validator.w3.org/feed/docs/atom.html

package site

import (

	// Standard:
	"encoding/xml"
	"time"

	// Third-Party:
	"golang.org/x/tools/blog/atom"

	// Kisipar:
	"github.com/biztos/kisipar/pageset"
)

// SaneEntry extends a Google Atom Entry to support an xml:base attribute,
// thus making it minimally sane.  WTF, GOOG?
type SaneEntry struct {
	atom.Entry
	XMLBase string `xml:"xml:base,attr"`
}

// SaneFeed allows a Google-style Feed to have an xml:base attribute per
// entry, and thus not be COMPLETELY EFFING INSANE.
// (TODO: write own, this shit is not so fancy.)
type SaneFeed struct {
	XMLName xml.Name     `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string       `xml:"title"`
	ID      string       `xml:"id"`
	Link    []atom.Link  `xml:"link"`
	Updated atom.TimeStr `xml:"updated"`
	Author  *atom.Person `xml:"author"`
	Entry   []*SaneEntry `xml:"entry"`
}

// Feed returns an Atom Feed (as a SaneFeed) for the Site based on its
// configuration.
// If a Pageset is provided, that is used for the Entries; otherwise
// the main Site Pageset is used.
// Pages are used in descending Time order, i.e. Updated > Created; ModTime
// is the fallback. The Feed's Updated field is the time of the newest
// entry.
// TODO: An optional Template for the Feed.
// TODO: Support more parts of the Person for author.
func (s *Site) Feed(ps *pageset.Pageset) *SaneFeed {

	if ps == nil {
		ps = s.Pageset
	}

	f := &SaneFeed{
		Title: s.FeedTitle,
		ID:    s.BaseURL + s.FeedPath,
		Author: &atom.Person{
			Name: s.Owner,
		},
		Updated: atom.Time(time.Now()), // default if no entries
		Link: []atom.Link{
			atom.Link{
				Rel:  "self",
				Href: s.BaseURL + s.FeedPath,
			},
			atom.Link{
				Rel:  "alternate",
				Href: s.BaseURL,
			},
		},
	}
	entries := make([]*SaneEntry, s.FeedItems)

	pages := ps.ByTime()
	total := len(pages)
	if total > s.FeedItems {
		total = s.FeedItems
	}
	for i, p := range pages {
		ts := p.Time()
		if i == 0 {
			f.Updated = atom.Time(*ts)
		}
		if i == total {
			break
		}

		// If we have no page Author, try the site Owner.
		// TODO: collect all authors for the feed-level Author tags.
		// (For which probably have to rewrite the atom library right?)
		author := p.Author()
		if author == "" {
			author = s.Owner
		}

		// We always have content.
		// TODO: fix it up, or if not needed then put it down there.
		// (links are fine but need to be absolute)
		content := p.Content

		// Use the self-URL as the xml:base until proven wrong.
		selfurl := s.PageURL(p)
		e := atom.Entry{
			Title:   p.Title(),
			ID:      s.PageURL(p),
			Updated: atom.Time(*ts),
			Author: &atom.Person{
				Name: author,
			},
			Content: &atom.Text{
				Type: "html",
				Body: string(content),
			},
			Link: []atom.Link{
				atom.Link{
					Rel:  "self",
					Href: selfurl,
				},
			},
		}

		// Published?  (This is pretty useful but not required.)
		if pub := p.Created(); pub != nil {
			e.Published = atom.Time(*pub)
		}

		// Summarized?
		if summary := p.Summary(); summary != "" {
			e.Summary = &atom.Text{
				Type: "text",
				Body: string(summary),
			}
		}

		entries[i] = &SaneEntry{
			e,
			selfurl,
		}

	}
	f.Entry = entries[:total]
	return f

}
