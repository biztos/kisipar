// handlers.go -- kisipar http handler logic
// -----------

package kisipar

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Handler is an http.Handler that responds to requests with content from a
// Provider.
type Handler struct {
	site *Site
}

// NewHandler returns a Handler using the given Site, or an error if the
// Site is nil.
func NewHandler(s *Site) (*Handler, error) {
	if s == nil {
		return nil, errors.New("Site must not be nil")
	}
	return &Handler{s}, nil
}

// Error replies to the request with the given HTTP error code, error
// message, and optional detail. A StandardPage is created to hold the
// error, with a Path of "/error/<code>" (e.g. "error/404") and msg as its
// Title. The detail is stored in the Page HTML. A template is sought for
// the error based on the Path as noted above. If none is found, then
// http.Error is used to serve the error; however, note that the default
// template, if implemented, is also valid for serving errors.
//
// The error is logged to the standard logger.
func (h *Handler) Error(w http.ResponseWriter, r *http.Request, code int, msg, detail string) {

	log.Printf("%s %s %s %d %s: %s",
		r.RemoteAddr, r.Method, r.URL, code, msg, detail)

	p, _ := StandardPageFromData(
		map[string]interface{}{
			"path":    fmt.Sprintf("/error/%d", code),
			"title":   msg,
			"created": time.Now(),
			"updated": time.Now(),
			"html":    detail,
		},
	)
	tmpl := h.site.Provider.TemplateFor(p)
	if tmpl == nil {
		// Fall back to simple error.
		http.Error(w, msg, code)
		return
	}

	// Let's not forget our HTTP status!
	w.WriteHeader(code)
	tmpl.Execute(w, p)

}
