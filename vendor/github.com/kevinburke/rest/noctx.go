// +build !go1.7

package rest

import "net/http"

// CtxErr returns an error that's been stored in the Request context. This is
// a no-op in older versions of Go.
func CtxErr(r *http.Request) error {
	return nil
}

// ctxSetErr sets err in the request context, and returns the new request.
func ctxSetErr(r *http.Request, err error) *http.Request {
	Logger.Info("cannot set error in context", "err", err)
	return r
}

// CtxDomain returns a domain that's been set on the request. This is a no-op
// in older versions of Go.
func CtxDomain(r *http.Request) string {
	return ""
}

// ctxSetDomain sets a domain in the request context, and returns the new
// request.
func ctxSetDomain(r *http.Request, domain string) *http.Request {
	Logger.Info("cannot set domain in context", "domain", domain)
	return r
}
