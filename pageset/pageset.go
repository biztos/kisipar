// pageset/pageset.go - core Kisipar Pageset code.
// ------------------

// Package pageset defines a set of pages in a Kisipar web site.
//
//     TODO:
//     * Rethink array / slice sorting stuff for byPath et al.  Benchmark mem!
//     * Tags() []string -> fetch all tags
//     * ByMetaString(key string) -> ordered by this meta string in page
//     * SubsetForTags([]string) -> all having all of these tags
//     * SubsetForMetaString(key,val string) -> all having meta string of this
//     * ByTime (and thus a decent Data func in Page)
//     * Desc and Lower versions of Path sort (ByPathDesc, ByPathLower[Desc])
//     * Desc and Lower versions of ModTime and Date sorts.
package pageset

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/biztos/kisipar/page"
)

type cache struct {

	// Listed pages, which are the basis of sorted sets:
	listed []*page.Page

	// Sorted sets:
	byPath    []*page.Page
	byCreated []*page.Page
	byModTime []*page.Page
	byTime    []*page.Page

	// Subsets as maps, since we don't know the usage in advance:
	pathSubsets map[string]*Pageset
	tagSubsets  map[string]*Pageset

	// Tags are a bit expensive to extract:
	tags []string
}

func (c *cache) clearAll() {
	c.listed = nil
	c.byPath = nil
	c.byCreated = nil
	c.byModTime = nil

	c.pathSubsets = map[string]*Pageset{}
	c.tagSubsets = map[string]*Pageset{}

	c.tags = nil

}

// A Pageset defines a set of Pages that may sorted and manipulated as needed.
// Pagesets can be as big as an entire Kisipar site, or as small as zero
// pages.  The Pageset expects that its Pages are not manipulated outside of
// the Pageset context.
//
// Each Page must be unique by extension-stripped Path, meaning a Pageset
// should not have Pages at both "foo/bar.md" and "foo/bar.txt".
//
// Pages are normally accessed, e.g. in a template, using one of the sorting
// "By*" methods, which remove unlisted pages.
type Pageset struct {

	// We keep a map of extension-stripped Paths to Pages, but don't allow it
	// to be manipulated directly:
	pageMap map[string]*page.Page

	// We cache aggressively, as it's only lists of pointers:
	cache *cache
}

// New creates a Pageset with the provided slice of Pages.  Each Page must
// be unique by extension-stripped path key.
func New(pages []*page.Page) (*Pageset, error) {

	pm := map[string]*page.Page{}
	for _, p := range pages {

		key := strings.TrimSuffix(p.Path, filepath.Ext(p.Path))
		if pm[key] != nil {
			return nil, fmt.Errorf("Duplicate path for %s: %s", key, p.Path)
		}
		pm[key] = p
	}

	return &Pageset{
		pageMap: pm,
		cache: &cache{
			pathSubsets: map[string]*Pageset{},
			tagSubsets:  map[string]*Pageset{},
		},
	}, nil
}

// Len returns the number of Pages in the Pageset.  Note that this includes
// Unlisted Pages.
func (ps *Pageset) Len() int {
	return len(ps.pageMap)
}

// Page returns the Page at the given path key, or nil if none is defined.
func (ps *Pageset) Page(key string) *page.Page {
	return ps.pageMap[key]
}

// AddPage adds a Page to the Pageset, and clears the sorting and subset
// caches.  If a Page for the Page's Path exists in the Pageset it is
// replaced.
func (ps *Pageset) AddPage(p *page.Page) {

	// TODO: selective unCache if we are replacing by path.
	// unCacheNonPath?  Or what? uncache(all bool)
	key := strings.TrimSuffix(p.Path, filepath.Ext(p.Path))
	ps.pageMap[key] = p
	ps.cache.clearAll()

}

// RemovePage removes the Page at the given normalized path key from the
// Pageset, and clears the sorting and subset caches.
func (ps *Pageset) RemovePage(key string) {

	delete(ps.pageMap, key)
	ps.cache.clearAll()

}

// RefreshPage refreshes a Page at the given path key, by reloading it,
// loading it into the Pageset, or removing it if it is no longer on disk.
// If the page is not found or any filesystem or parse error occurs, the
// error is returned; nil is returned on success.
func (ps *Pageset) RefreshPage(key string) error {

	// Refresh the page if it exists, removing it from the Pageset on error.
	if p := ps.Page(key); p != nil {
		oldModTime := p.ModTime
		if err := p.Refresh(); err != nil {
			delete(ps.pageMap, key)
			ps.cache.clearAll()

			if !isReallyNotExist(err) {
				// A parse error, or similar, occurred.
				return err
			}
		} else {
			if p.ModTime != oldModTime {
				ps.cache.clearAll()
			}
			return err
		}

	}

	// If we don't (any longer) have it, let's try to get it.
	// (It's a supported, if obscure, use-case to delete foo.md and have
	// foo.txt loaded in its place.)
	p, err := page.LoadAny(key)
	if err == nil {
		ps.pageMap[key] = p
		ps.cache.clearAll()
		return nil
	} else if isReallyNotExist(err) {
		return os.ErrNotExist
	} else {
		return err
	}

}

// ...for when IsNotExist isn't enough:
func isReallyNotExist(err error) bool {

	if os.IsNotExist(err) {
		return true
	}
	if strings.HasSuffix(err.Error(), "not a directory") {
		return true
	}
	return false
}

func (ps *Pageset) listedPages() []*page.Page {

	if len(ps.cache.listed) == 0 {
		pages := make([]*page.Page, len(ps.pageMap))
		i := 0
		for _, p := range ps.pageMap {
			if p.Unlisted == false {
				pages[i] = p
				i++
			}
		}
		ps.cache.listed = pages[0:i]
	}

	return ps.cache.listed

}

// Tags returns the unique lowercase tags found in the Pageset's Pages,
// alpha-sorted.
// Unlisted pages are ignored. The result is cached for future use.
func (ps *Pageset) Tags() []string {

	if len(ps.cache.tags) == 0 {
		have := map[string]bool{}
		tags := []string{}
		pages := ps.listedPages()
		for _, p := range pages {
			for _, t := range p.Tags() {
				t = strings.ToLower(t)
				if have[t] == false {
					have[t] = true
					tags = append(tags, t)
				}
			}
		}
		sort.Strings(tags)
		ps.cache.tags = tags

	}
	return ps.cache.tags

}

// ListedSubset returns a new Pageset containing only those pages that are
// not marked Unlisted.
func (ps *Pageset) ListedSubset() *Pageset {

	subset, err := New(ps.listedPages())
	if err != nil {
		// This can only be programmer error, since you should not be able
		// to get contradictory pages into the same pageset using the
		// normal methods.
		panic("ListedSubset failed for Pageset: " + err.Error())
	}

	return subset

}

// TagSubset returns a new Pageset containing only those pages that contain
// the given tag, case-insensitively.  The result is cached for future use.
func (ps *Pageset) TagSubset(tag string) *Pageset {

	tag = strings.ToLower(tag)
	if subset := ps.cache.tagSubsets[tag]; subset != nil {
		return subset
	}

	pages := []*page.Page{}
	for _, p := range ps.pageMap {
		for _, t := range p.Tags() {
			if strings.ToLower(t) == tag {
				pages = append(pages, p)
				break
			}
		}
	}
	subset, err := New(pages)
	if err != nil {
		// This can only be programmer error, since you should not be able
		// to get contradictory pages into the same pageset using the
		// normal methods.
		panic("TagSubset failed for Pageset: " + err.Error())
	}
	ps.cache.tagSubsets[tag] = subset

	return subset

}

// PathSubset returns a new Pageset containing those Pages whose Path has the
// provided prefix.  If a trim string is defined then it is trimmed from the
// beginning of each path before comparison.
// The subset is cached for future use.
// NOTE: any change (AddPage, RemovePage, RefreshPage) to the resulting
// Pageset will *NOT* affect the origin Pageset.  Thus in a normal website
// context, these operations should only be performed on the master Pageset.
// TODO: this is stupid, let the Pageset know about its trim string in general
func (ps *Pageset) PathSubset(prefix, trim string) *Pageset {

	// We need an equatable thing for our key, apparently; not []string.
	// So maybe just this...
	key := trim + "\n" + prefix
	if subset := ps.cache.pathSubsets[key]; subset != nil {
		return subset
	}

	pages := []*page.Page{}
	for _, p := range ps.pageMap {
		path := strings.TrimPrefix(p.Path, trim)
		if strings.HasPrefix(path, prefix) {
			pages = append(pages, p)
		}
	}
	subset, err := New(pages)
	if err != nil {
		// This can only be programmer error, since you should not be able
		// to get contradictory pages into the same pageset using the
		// normal methods.
		panic("PathSubset failed for Pageset: " + err.Error())
	}
	ps.cache.pathSubsets[key] = subset

	return subset

}

// PageSlice returns a slice of the input pages with the given start and end
// positions.  If start is out of bounds for the slice then an empty slice
// is returned.  If end is out of bounds then it is set to the length of the
// slice, i.e. to the end. A negative end is omitted, i.e. [start:].
//
// This function is intended for use in templates, e.g. in order to show the
// first N items of a Pageset:
//
//    {{ with .Pageset }}{{ $pp := .Pageslice .ByTime 0 20 }}...{{ end }}
func (ps *Pageset) PageSlice(pages []*page.Page, start, end int) []*page.Page {

	if start > end && end > 0 {
		panic("start may not be greater than end")
	}
	if start < 0 {
		panic("start may not be negative")
	}
	if start > len(pages) {
		return pages[0:0]
	}
	if end < 0 || end > len(pages) {
		return pages[start:]
	}

	return pages[start:end]
}

// --------------------------------------------------------------------------
// SORTING!
// --------------------------------------------------------------------------
// Sorting in Go is *fun*! (Not fun.  But arguably powerful.)

// SORT BY PATH, RESPECTING DEPTH AND INDEXES:
type byPath []*page.Page

func (s byPath) Len() int {
	return len(s)
}
func (s byPath) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPath) Less(i, j int) bool {

	// First, compare the path depths.
	pp1 := strings.Split(s[i].Path, string(filepath.Separator))
	pp2 := strings.Split(s[j].Path, string(filepath.Separator))
	if len(pp1) != len(pp2) {
		return len(pp1) < len(pp2)
	}

	// Second, give precedence to indexes.
	i1 := s[i].IsIndex
	i2 := s[j].IsIndex
	if i1 && !i2 {
		return true
	} else if i2 && !i1 {
		return false
	}

	// Fall back to string comparison.
	return s[i].Path < s[j].Path
}

// ByPath returns an array of the Pageset's Pages sorted by Path: the first
// sort key is the length of the Path, the second is whether the page is
// an index (e.g. "foo/index.md") and the third is the path itself.
// Unlisted Pages are excluded.  The result is cached for future use.
func (ps *Pageset) ByPath() []*page.Page {

	if ps.cache.byPath != nil {
		return ps.cache.byPath
	}

	listed := ps.listedPages()
	pages := make([]*page.Page, len(listed))
	copy(pages, listed)
	sort.Sort(byPath(pages))
	ps.cache.byPath = pages

	return pages
}

// SORT BY TIME: CREATED (FALLBACK: FILE MOD TIME; TIEBREAKER: PATH)
type byCreated []*page.Page

func (s byCreated) Len() int {
	return len(s)
}
func (s byCreated) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCreated) Less(i, j int) bool {

	iTime := s[i].Created()
	if iTime == nil {
		iTime = &s[i].ModTime
	}
	jTime := s[j].Created()
	if jTime == nil {
		jTime = &s[j].ModTime
	}

	if iTime.Equal(*jTime) {
		// Tiebreaker: path comparison using our path sorting logic.
		bp := byPath{s[i], s[j]}
		return bp.Less(0, 1)
	}

	return iTime.After(*jTime)

}

// ByCreated returns an array of the Pageset's Pages sorted by Created time,
// newest first, as set in the page meta. If no Created time is set, the
// Page's ModTime is used. Ties (to the nanosecond) are broken by a smart
// Path comparison.  Unlisted Pages are excluded.  The result is cached for
// future use.
func (ps *Pageset) ByCreated() []*page.Page {

	if ps.cache.byCreated != nil {
		return ps.cache.byCreated
	}

	listed := ps.listedPages()
	pages := make([]*page.Page, len(listed))
	copy(pages, listed)
	sort.Sort(byCreated(pages))
	ps.cache.byCreated = pages

	return pages
}

// SORT BY TIME: MOD TIME (TIEBREAKER: PATH)
type byModTime []*page.Page

func (s byModTime) Len() int {
	return len(s)
}
func (s byModTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byModTime) Less(i, j int) bool {

	iTime := s[i].ModTime
	jTime := s[j].ModTime

	if iTime.Equal(jTime) {
		// Tiebreaker: path comparison using our path sorting logic.
		bp := byPath{s[i], s[j]}
		return bp.Less(0, 1)
	}

	return iTime.After(jTime)

}

// ByModTime returns an array of the Pageset's Pages sorted by ModTime, newest
// first.  Unlisted Pages are excluded.  The result is cached for future use.
// For pages with exactly the same modification time (to the nanosecond) the
// secondary sort key is the Path.
func (ps *Pageset) ByModTime() []*page.Page {

	if ps.cache.byModTime != nil {
		return ps.cache.byModTime
	}

	listed := ps.listedPages()
	pages := make([]*page.Page, len(listed))
	copy(pages, listed)
	sort.Sort(byModTime(pages))
	ps.cache.byModTime = pages

	return pages

}

// SORT BY TIME: DATE FUNCTION (TIEBREAKER: PATH)
type byTime []*page.Page

func (s byTime) Len() int {
	return len(s)
}
func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTime) Less(i, j int) bool {

	iTime := s[i].Time()
	jTime := s[j].Time()

	if iTime.Equal(*jTime) {
		// Tiebreaker: path comparison using our path sorting logic.
		bp := byPath{s[i], s[j]}
		return bp.Less(0, 1)
	}

	return iTime.After(*jTime)

}

// ByTime returns an array of the Pageset's Pages sorted by the Time function,
// newest first.  Unlisted Pages are excluded.  The result is cached for
// future use.
// For pages with exactly the same Time() (to the nanosecond) the secondary
// sort key is the Path.
func (ps *Pageset) ByTime() []*page.Page {

	if ps.cache.byTime != nil {
		return ps.cache.byTime
	}

	listed := ps.listedPages()
	pages := make([]*page.Page, len(listed))
	copy(pages, listed)
	sort.Sort(byTime(pages))
	ps.cache.byTime = pages

	return pages

}

// Reverse reverses a slice of Pages, returning a new slice.
func Reverse(pages []*page.Page) []*page.Page {

	last := len(pages) - 1
	reversed := make([]*page.Page, len(pages))
	for i, p := range pages {
		reversed[last-i] = p
	}
	return reversed
}
