// serving.go - web-serving functions of the Kisipar site.
// ----------

package site

// Serve serves the Site, either in TLS mode or insecure.  For most purposes
// insecure will be preferable, as you should have another layer between
// Kisipar and the open internet.
//
// An error is always returned, as per http.ListenAndServe.
//
// Serve does NOT block to wait for the server to finish, as one may serve
// multiple Kisipar sites from a single application.  In order to block in
// the normal fashion, wrap the final Serve call in log.Fatal or similar.
func (s *Site) Serve() error {

	if s.Server == nil {
		panic("Serve called but Server is nil.") // ...at least helpful.
	}

	if s.ServeTLS {
		return s.Server.ListenAndServeTLS(s.CertFile, s.KeyFile)

	} else {
		return s.Server.ListenAndServe()
	}

}
