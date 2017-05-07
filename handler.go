// handlers.go -- kisipar http handler logic
// -----------

package kisipar

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// Handler is an http.Handler that responds to requests with content from a
// Provider.  It may be used for the catch-all / root handler ("/").
type Handler struct {
	Site *Site
}

// NewHandler returns a Handler using the given Site, or an error if the
// Site is nil or its Provider is nil.
func NewHandler(s *Site) (*Handler, error) {
	if s == nil {
		panic("Site must not be nil.")
	}
	if s.Provider == nil {
		panic("Site.Provider must not be nil.")
	}
	return &Handler{Site: s}, nil
}

// ServeHTTP implements the http.Handler interface, and responds to a request
// using the standard Kisipar logic:
//
//   1. For a request with the root path ("/"), a template named "/index.html"
//      is looked up via the Provider. If there is no such template then a
//      text/plain response is served with the Site name as sole content,
//      with the name defaulting to "Kisipar."
//
//   2. If the Site has a configured StaticDir then the file is checked
//      using the Site's StaticPath method; see there for configuration
//      details.  If StaticPath returns a path and no error, the file is
//      delivered; errors are ignored.
//
//   3. If there is a template matching the request path then we execute that
//      template.  The template is sought via the Provider's TemplateFor
//      method.
//
//   4. If the Provider returns a page via Get() or GetSince() then it is
//      served according to its type (Redirect, File, Content, or Page).
//      GetSince is used if there is an If-Modified-Since header in the
//      request and its timestamp can be parsed; in that case, a response of
//      304 Not Modified may be served.
//
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s := h.Site
	path := r.URL.Path

	// Root path is special:
	if path == "/" {
		if tmpl := s.Provider.TemplateForPath("/index.html"); tmpl != nil {

			dot := &Dot{
				Request:  r,
				Config:   s.Config,
				Provider: s.Provider,
			}
			h.ServeTemplate(w, dot, tmpl, http.StatusOK)
			return
		}

		name := h.Site.Config.Name
		if name == "" {
			name = "Kisipar."
		}
		fmt.Fprintln(w, name+"\n")
		return
	}

	// Static files (and potentially dirs) are http's problem:
	if file, err := s.StaticPath(path); err == nil {
		http.ServeFile(w, r, file)
		return
	}

	// We should have quick access to the template if it's there.
	if tmpl := s.Provider.TemplateForPath(path); tmpl != nil {
		dot := &Dot{
			Request:  r,
			Config:   s.Config,
			Provider: s.Provider,
		}
		h.ServeTemplate(w, dot, tmpl, http.StatusOK)
		return
	}

	// Freshness is up to the Provider.  We ignore range requests.
	if r.Header.Get("Range") == "" {
		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if ts, err := http.ParseTime(ims); err == nil {
				// We have a time, use it.
				p, err := h.Site.Provider.GetSince(r.URL.Path, ts)
				if err != nil && !IsNotExist(err) {
					// Uh-oh, data-provider error!
					h.Error(w, r, http.StatusInternalServerError,
						"Data provider error.", err.Error())
					return
				}
				h.ServeItem(w, r, p)
				return
			}
		}
	}

	// Final check: get it from the provider.
	p, err := h.Site.Provider.Get(r.URL.Path)
	if err != nil {
		if IsNotExist(err) {
			h.Error(w, r, http.StatusNotFound, "Not Found", "")
		} else {
			// Data-provider error!
			h.Error(w, r, http.StatusInternalServerError,
				"Data provider error.", err.Error())
		}
		return
	}

	h.ServeItem(w, r, p)
}

// ServeItem serves a Pather which may be a Page, Content, File, Path,
// PathHandler or Redirect.
func (h *Handler) ServeItem(w http.ResponseWriter, r *http.Request, p Pather) {

	if h.Site == nil {
		panic("Site must not be nil.")
	}
	if h.Site.Provider == nil {
		panic("Site.Provider must not be nil.")
	}

	switch item := p.(type) {
	case Page:

		if tmpl := h.Site.Provider.TemplateFor(item); tmpl != nil {
			dot := &Dot{
				Request:  r,
				Config:   h.Site.Config,
				Provider: h.Site.Provider,
				Page:     item,
			}
			h.ServeTemplate(w, dot, tmpl, http.StatusOK)
		}

		// In normal operation you'd always have a template, because the
		// provider will do something reasonable for the default.  But if
		// not then not: serve the plain HTML without any template and hope
		// for the best.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(item.HTML()))
	case Content:
		if ct := item.ContentType(); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		http.ServeContent(w, r, item.Path(), item.ModTime(), item.ReadSeeker())
	case File:
		http.ServeFile(w, r, item.FilePath())
	case PathHandler:
		item.Func()(w, r)
	case Redirect:
		code := http.StatusFound
		if item.Permanent() {
			code = http.StatusMovedPermanently
		}
		http.Redirect(w, r, item.Location(), code)
	default:
		detail := fmt.Sprintf(
			"Pather at %s has type %T; handler can not Serve it.",
			p.Path(), p)
		h.Error(w, r, http.StatusInternalServerError,
			"Data provider error.", detail)
	}

}

// Error replies to the request with the given HTTP error code, error
// message, and optional detail. A StandardPage is created to hold the
// error, with a Path of "/errors/<code>" (e.g. "error/404") and msg as its
// Title. The detail is stored in the Page HTML. A template is sought for
// the error based on the Path as noted above. If none is found, then
// http.Error is used to serve the error; however, note that the default
// template, if implemented, is also valid for serving errors.
//
// The error is logged to the standard logger.
func (h *Handler) Error(w http.ResponseWriter, r *http.Request, code int, msg, detail string) {

	log.Printf("ERROR: %s %s %s %d %s: %s",
		r.RemoteAddr, r.Method, r.URL, code, msg, detail)

	p, _ := StandardPageFromData(
		map[string]interface{}{
			"path":    fmt.Sprintf("/errors/%d", code),
			"title":   msg,
			"created": time.Now(),
			"updated": time.Now(),
			"html":    detail,
		},
	)
	if tmpl := h.Site.Provider.TemplateFor(p); tmpl != nil {
		dot := &Dot{
			Request:  r,
			Config:   h.Site.Config,
			Provider: h.Site.Provider,
			Page:     p,
		}
		h.ServeTemplate(w, dot, tmpl, code)
		return
	}

	// Fall back to simple error.
	http.Error(w, msg, code)
	return

}

// ServeTemplate executes template t with Dot d, serving the result or a
// generic error message (the detailed error message is logged).  This
// method allows template errors to be caught without partial output written
// to the client; however, it is inefficient.  If the Site Config
// FastTemplates is true then the template will be rendered directly and,
// in case of errors, partially.
//
// The Content-Type is set to the standard MIME type for the template name's
// extension, or text/html by default.
func (h *Handler) ServeTemplate(w http.ResponseWriter, dot *Dot, tmpl *template.Template, code int) {

	// Set the content type based on the template name.
	ct := "text/html; charset=utf-8"
	name := tmpl.Name()
	ext := filepath.Ext(name)
	if ext != "" {
		if found := mime.TypeByExtension(name); found != "" {
			ct = found
		}
	}

	// As a special case, we make sure text/plain gets a charset.
	if ct == "text/plain" {
		ct = "text/plain; charset=utf-8"
	}

	w.Header().Set("Content-Type", ct)
	w.WriteHeader(code)

	// The unsafe route for the speed-conscious:
	if h.Site.Config.FastTemplates {
		err := tmpl.Execute(w, dot)
		if err != nil {
			log.Printf("ERROR: Template error in %s: %s", name, err.Error())
		}
		return
	}

	// The safer, saner, slower, memory-wasting route:
	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, dot); err != nil {
		log.Printf("ERROR: Template error in %s: %s", name, err.Error())
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		return
	}
	w.Write(b.Bytes())

}
