// handler.go - kisipar site multiplexing http handler (mux)
// ----------

package site

import (
	// Standard:
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	// Kisipar:
	"github.com/biztos/kisipar/page"
	"github.com/biztos/kisipar/pageset"
)

// PatternHandler defines an HTTP handler or handler function together with
// the pattern used by its multiplexer (Mux).  Use this to wrap arguments to
// NewServeMux, via NewPatternHandler and NewPatternHandlerFunc.  Function
// and Handler should not both be defined, and the Pattern may not be an
// empty string.
type PatternHandler struct {
	Pattern  string
	Function func(http.ResponseWriter, *http.Request)
	Handler  http.Handler
}

// NewServeMux creates and returns a multiplexing HTTP request handler -- an
// http.ServeMux -- based on the Site's properties. This is set in the
// default Server created by the New function, but it may also be called in
// order to use the standard request-handling logic elsewhere.
//
// If no handler is provided for the pattern "/" then one will be obtained
// from the MainHandler function.
//
// Panics if duplicate patterns are passed in the handlers, or if any passed
// PatternHandler is malformed. (Such a case would unambiguously indicate
// programmer error.)
func (s *Site) NewServeMux(handlers ...*PatternHandler) *http.ServeMux {

	// Register any override handlers:
	havePattern := map[string]bool{}
	for i, h := range handlers {
		p := h.Pattern
		if p == "" {
			msg := fmt.Sprintf("Empty handler pattern at position %d.", i)
			panic(msg)
		}
		if havePattern[h.Pattern] {
			msg := fmt.Sprintf("Duplicate handler pattern: %s", p)
			panic(msg)
		}
		if h.Function != nil && h.Handler != nil {
			msg := fmt.Sprintf("Redundant handler and function for pattern: %s", p)
			panic(msg)
		}
		if h.Function == nil && h.Handler == nil {
			msg := fmt.Sprintf("Neither handler nor function for pattern: %s", p)
			panic(msg)
		}
		havePattern[h.Pattern] = true
	}
	mux := http.NewServeMux()
	for _, h := range handlers {
		if h.Handler != nil {
			mux.Handle(h.Pattern, h.Handler)
		} else {
			mux.HandleFunc(h.Pattern, h.Function)
		}
	}

	// Set up the standard handlers based on the site's setup:
	// TODO ...

	// Register the main fallback handler if it was not already set.
	// TODO: real stuff of course.
	if !havePattern["/"] {
		mux.HandleFunc("/", s.MainHandler())
	}
	return mux
}

// MainHandler returns an HTTP handler function applying the Kisipar logic
// to the current Site.  That logic, in a nutshell, is:
//
// 1. Static files take priority, followed by Pages, then page assets.
//
// 2. Static paths containing no extension look for an "index.html" file under
//    the implied directory, e.g.: "/foo" will match "/foo/index.html".
//
// 3. Requests containing a dot (".") anywhere in the cleaned path are not
//    considered potential Pages; those not containing any extension are not
//    considered potential page-asset files.
//
// 4. Pages are sought at the path key first, then at its index, e.g.:
//    "foo/bar" -> "foo/bar.md" OR "foo/bar/index.md" (where the ".md"
//    extension is the first match from the Site's PageExtensions).
//
// 5. Directories are treated as not-found.
//
// 6. Page source files are available as page assets *only* if the Site's
//    ServePageSources property is set to true; otherwise no page asset
//    with an extension matching the Site's PageExtensions will be found.
//
// 7. If the top of the site is not otherwise handled, a simple default page
//    is served.
func (s *Site) MainHandler() func(w http.ResponseWriter, req *http.Request) {

	// TODO: figure out what to do about news feeds, contact forms, any
	// redirects or even executables, etc; and set that stuff up in a way
	// that's cheap to check later.  Probably boolean closures.

	return func(w http.ResponseWriter, req *http.Request) {

		// Cleaning the path is a bit expensive, so we do it once.
		rpath := path.Clean(req.URL.Path)

		// TODO: special cases as needed

		// Static beats everything else, because you may need to drop in
		// a static file in an emergency.  This includes special cases such
		// as the feed.
		if s.handleStatic(w, req, rpath) {
			return
		}

		// Feed is handled if set, and can not be overridden by a Page or
		// index.
		if s.handleFeed(w, req, rpath) {
			return
		}

		// Check for a proper Page, or index.
		if s.handlePage(w, req, rpath) {
			return
		}

		// Check for page-level assets, which are a special kind of static
		// file.
		if s.handleAsset(w, req, rpath) {
			return
		}

		// Nothing left, so it's four oh four.
		//
		// TODO: cache the 404's up to some point in time so we can just
		// skip all the above checks when we know we're getting 404's.
		// (However, not infinitely, or it's just a different DoS vector.)
		// (Consider a pre-sized hash for that, with a pre-sized key array,
		// and a mutex for updates. Hash is of ages, so you can time it out.
		// If you get a "new" 404 then you put it back into a random slot
		// in the hash.  This might (or might not) be a good defense against
		// randomized attacks.  The hash+array lookup is much faster than
		// a stat call to disk, presumably even with the mutex and lookup
		// costs, though that might be worth benchmarking.)
		s.sendNotFound(w, req)

	}

}

// PageForPath returns a Page from the Site's Pageset for the given cleaned
// request Path.
//
// Exact matches on virtual pages take precedence, but they *must* be on
// virtual pages: a non-virtual page with a key matching a request path
// would be a dangerous coincidence, as the Site's PagePath is part of the
// key.
//
// Otherwise a match is sought in the Pageset after first refreshing the path
// key.  First the path itself is sought, then the path's index: "foo/bar"
// then "foo/bar/index".
//
// Pages not found will result in os.IsNotExist style errors; parse or
// filesystem errors are returned as-is.
//
// TODO: subject this to some kind of PagesetTTL so we don't hit the disk
// all the time for duplicate or, worse, random requests.
// ...but don't call it that!  Because it implies -- and we might want this
// -- that we would reload the whole Pageset every so often.  (Odd use-case
// vs. bouncing the server.)
// PageRefreshInterval?  AND ALSO TODO: whether to check disk at all.
// PagesRefresh, PageRefreshInterval
//
// TODO: have a set of prefixes that are always not-found, or maybe always
// something else, in case there are persistent annoying attacks that get
// 404 errors (which you'd otherwise want to monitor, to watch for broken
// incoming links)
func (s *Site) PageForPath(rpath string) (*page.Page, error) {

	// It's an odd corner case but if you don't have any Pageset, everything
	// is by definition not-found.
	if s.Pageset == nil {
		return nil, os.ErrNotExist
	}

	idxpath := path.Join(rpath, "index")

	// Virtual pages must match on the sub-path: /foo, not /site/pages/foo.
	if p := s.Pageset.Page(rpath); p != nil && p.Virtual == true {
		return p, nil
	}
	if p := s.Pageset.Page(idxpath); p != nil && p.Virtual == true {
		return p, nil
	}

	// Everything else is checked against pages on disk, if we have any.
	if s.PagePath == "" {
		return nil, os.ErrNotExist
	}

	key := filepath.Join(s.PagePath, filepath.FromSlash(rpath))
	idxkey := filepath.Join(s.PagePath, filepath.FromSlash(idxpath))

	// Exact match takes precedence.
	if err := s.Pageset.RefreshPage(key); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if p := s.Pageset.Page(key); p != nil {
		return p, nil
	}

	// Index is our fallback.  Not-found errors are caught here too.
	if err := s.Pageset.RefreshPage(idxkey); err != nil {
		return nil, err
	}

	// Obvioiusly this creates a tiny race condition between RefreshPage
	// and Page, but the alternative (so far) appears to be untested code.
	// TODO: figure out how to avoid this, ideally without a bunch of stupid
	// test rigging.  Maybe put the lock here instead of in RefreshPage?
	// (That might be a good idea.  Where does the lock belong anyway?)
	return s.Pageset.Page(idxkey), nil
}

// TODO: gussy these up! Use templates and Dots if available, etc.
func (s *Site) sendInternalServerError(w http.ResponseWriter, req *http.Request, err error) {
	// TODO: log the error, don't put it here.
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ERROR: %s\n", err.Error())

}
func (s *Site) sendNotFound(w http.ResponseWriter, req *http.Request) {
	// TODO: log the error IF that is set.  So, something like:
	// if s.LogNotFoundErrors { ... }
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "NOT FOUND")
}

// Send a static page, or error out, and let the caller know whether the
// response has been served.  Not-found is *not* an error in this case, and
// directories are treated as not-found.  Note that the "stat" step is
// required, AFAICT, one way or the other, since we need the ModTime.
func (s *Site) handleStatic(w http.ResponseWriter, req *http.Request, rpath string) bool {

	if s.StaticPath == "" {
		return false
	}

	// Directories are assumed to mean their HTML index, allowing one to
	// override dynamic content in a pinch.  Yes, this is a real use-case!
	if filepath.Ext(rpath) == "" {
		rpath = path.Join(rpath, "index.html")
	}

	fpath := filepath.Join(s.StaticPath, filepath.FromSlash(rpath))
	info, err := os.Stat(fpath)
	if err != nil {
		// Nonexistence and other path errors are treated as not-found.
		return false
	}

	// Double check for "foo.dir" et al:
	if info.IsDir() {
		return false
	}

	// We have a file, but can we read it?  If we can't, bail out through our
	// custom 500 response instead of leaking config info (and ugliness) on
	// the generic 403.  Thus we don't use http.ServeFile() anymore.
	file, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("Open error on static %s: %s", rpath, err.Error())
		s.sendInternalServerError(w, req, err)
		return true
	}
	defer file.Close()
	http.ServeContent(w, req, fpath, info.ModTime(), file)
	return true

}

// Send a News Feed if the path matches the one set in the site config.
// TODO: sub-feeds, tag-feeds, etc.  Handle them here; config TBD.
func (s *Site) handleFeed(w http.ResponseWriter, req *http.Request, rpath string) bool {

	if s.NoFeed || rpath != s.FeedPath {
		return false
	}

	f := s.Feed(nil)
	data, err := xml.MarshalIndent(&f, "", "    ")
	if err != nil {
		s.sendInternalServerError(w, req, err)
		return true
	}

	// Good enough for now! Send it.
	w.Header().Set("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.sendInternalServerError(w, req, err)
		return true
	}

	return true

}

// Send a Page, or error out, returning true if the response has been served.
func (s *Site) handlePage(w http.ResponseWriter, req *http.Request, rpath string) bool {

	// We don't bother with anything that has a dot in its path.
	// Want dot dirs? Sorry!
	if strings.Contains(rpath, ".") {
		return false
	}

	// We take the given page if we have it, but if we don't we might still
	// have an index.
	p, err := s.PageForPath(rpath)
	if err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("Page error for %s: %s", rpath, err.Error())
		s.sendInternalServerError(w, req, err)
		return true
	}

	// Shall we have a Pageset?  And if so, which one?
	// TODO: make sure pageset logic is working (it might well not be)
	var ps *pageset.Pageset
	if rpath == "/" {
		ps = s.Pageset
	} else if p != nil && p.IsIndex {
		prefix := filepath.Dir(p.Path) + string(os.PathSeparator)
		ps = s.Pageset.PathSubset(prefix, s.PagePath)
	} else if p == nil {

		fpath := filepath.Join(s.PagePath, filepath.FromSlash(rpath))
		prefix := fpath + string(os.PathSeparator)
		ps = s.Pageset.PathSubset(prefix, s.PagePath)
		if ps.Len() == 0 {
			// No such subset, ergo no index page to handle.
			return false
		}

	}

	// Prep & Serve, we should be good here.
	dot := &Dot{
		Request: req,
		Page:    p,
		Pageset: ps,
		Site:    s,
		Now:     time.Now(),
	}

	// PROBLEM: we want control of the header, which we lose if write to the
	// thing first, but if we don't then we buffer the stupid rendered page.
	// (Anyway I suppose we want a caching option, so this is probably fine.)
	// (Maybe this doesn't matter at all with Nginx fronting us?)
	tmpl := dot.Template()
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, dot); err != nil {
		s.sendInternalServerError(w, req, err)
		return true
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// TODO: figure out whether we need to (want to) set the content-length,
	// especially in a reverse-proxy context a la Nginx.
	if _, err := w.Write(buf.Bytes()); err != nil {

		// It's not clear how this might be triggered in real life, but
		// we should catch it just the same.  However sending an internal
		// server error is not really an option anymore.
		msg := fmt.Sprintf("Write failed for %s: %s", rpath, err.Error())
		panic(msg)
	}

	// TODO: log something obvious here, similar to Apache common but in
	// nice JSON format.
	return true

}

// Send a page asset, or error out.  Behaves the same as sendAsset but also
// limits the type of page that can be served based on site settings.
func (s *Site) handleAsset(w http.ResponseWriter, req *http.Request, rpath string) bool {

	if s.PagePath == "" {
		return false
	}

	ext := filepath.Ext(rpath)
	if ext == "" {
		return false
	}

	// Specifically disallow anything in our extensions list unless we are
	// serving sources.  NOTE: this is just about as fast as using a map,
	// according to the benchmark, and easier to maintain.
	if !s.ServePageSources {
		for _, pext := range s.PageExtensions {
			if ext == pext {
				return false
			}
		}
	}

	fpath := filepath.Join(s.PagePath, filepath.FromSlash(rpath))
	info, err := os.Stat(fpath)
	if err != nil {
		// Nonexistence and other path errors are treated as not-found.
		return false
	}

	// Directories are not served from static, thus allowing you to have a
	// dynamic page at "foo/bar" with a static resource at "foo/bar/x.js"
	// if you so desire.
	// TODO: consider (and document it!) whether to allow serving of static
	// directories, so you could store collections of stuff there if you want.
	// Things like build files say.
	if info.IsDir() {
		return false
	}

	// We have a file, but can we read it?  If we can't, bail out through our
	// custom 500 response instead of leaking config info (and ugliness) on
	// the generic 403.  Thus we don't use http.ServeFile() anymore.
	file, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("Open error on static %s: %s", rpath, err.Error())
		s.sendInternalServerError(w, req, err)
		return true
	}
	defer file.Close()
	http.ServeContent(w, req, fpath, info.ModTime(), file)
	return true

}

// SetHandler sets the http.Handler in the Site's Server.  Panics if the
// Server is not an http.Server.  This function should be used for overriding
// the Handler in an otherwise standard Server.
func (s *Site) SetHandler(h http.Handler) {
	svr, ok := s.Server.(*http.Server)
	if !ok {
		panic("Server is not an http.Server.")
	}
	svr.Handler = h
}

// ServeHTTP calls the ServeHTTP method on the Site's Server's Handler. Panics
// if the Server is not an http.Server.  This function exists mostly to
// support testing, but may also be useful in exotic use-cases.
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svr, ok := s.Server.(*http.Server)
	if !ok {
		panic("Server is not an http.Server.")
	}
	svr.Handler.ServeHTTP(w, r)
}
