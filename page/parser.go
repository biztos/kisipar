// page/parser.go - Kisipar page parsing.
// --------------

package page

import (
	"github.com/biztos/kisipar/frostedmd"
)

// The Parser interface defines a struct capable of parsing source data into
// a meta map and HTML content.
type Parser interface {
	Parse([]byte) (ParseResult, error)
}

// The ParseResult interface defines a struct representing a Parser's Parse
// output.
type ParseResult interface {
	Meta() map[string]interface{}
	Content() []byte
}

// ExtParser pairs an extension to a Parser for use in the ordered ExtParsers.
type ExtParser struct {
	Ext    string
	Parser Parser
}

// ExtParsers defines the parsers that will be used for various extensions in
// Load and LoadVirtual.  Note that extensions are matched case-insensitively,
// but LoadAny is case-sensitive; thus if you specifically expecte to look up
// a particular cased extension (".Txt" for instance) then you need to set
// that up here in ExtParsers.
//
// The standard values should be sufficient for any normal Kisipar site.
var ExtParsers = []*ExtParser{
	{".md", &MdParser{}},       // .md -> Markdown
	{".MD", &MdParser{}},       // .MD -> Markdown (hello MSDOS!)
	{".markdown", &MdParser{}}, // .markdown - rare but not unheard of!
	{".MARKDOWN", &MdParser{}}, //
	{".txt", &MdParser{}},      // Text can also be Markdow, and vice-
	{".TXT", &MdParser{}},      // versa.
}

// LimitExtParsers removes any ExtParser entries from ExtParsers not matching
// the provided set of extensions, putting the resulting ExtParsers in the
// order of the extlist received.  This is useful for keeping the behavior
// of LoadAny in sync with any site-wide preload operations.
//
// Note that this is destructive: once limited, the original set of ExtParsers
// is gone.
func LimitExtParsers(extlist []string) {
	have := map[string]*ExtParser{}
	for _, ep := range ExtParsers {
		have[ep.Ext] = ep
	}

	keepers := []*ExtParser{}
	for _, ext := range extlist {
		if ep := have[ext]; ep != nil {
			keepers = append(keepers, ep)
		}
	}

	ExtParsers = keepers

}

// MdParser is the standard Markdown parser, using the frostedmd parsing
// logic based on blackfriday.
type MdParser struct{}

// Parse implements the Parser interface for the MdParser type.
func (p *MdParser) Parse(b []byte) (ParseResult, error) {

	return frostedmd.MarkdownCommon(b)

}

// VerbatimParser is a parser that simply returns its input with an empty
// meta map.
type VerbatimParser struct{}

// VerbatimParseResult is the result of a verbatim parse.
type VerbatimParseResult struct {
	meta    map[string]interface{}
	content []byte
}

// Meta implements the ParseResult interface for the VerbatimParseResult type.
func (vpr *VerbatimParseResult) Meta() map[string]interface{} {
	return vpr.meta
}

// Content implements the ParseResult interface for the VerbatimParseResult
// type.
func (vpr *VerbatimParseResult) Content() []byte {
	return vpr.content
}

// Parse implements the Parser interface for the VerbatimParser type.
func (v *VerbatimParser) Parse(b []byte) (ParseResult, error) {

	return &VerbatimParseResult{
		meta:    map[string]interface{}{},
		content: b,
	}, nil
}

// DefaultParser defines the Parser to use if no more specific Parser is
// available.
var DefaultParser Parser = &VerbatimParser{}
