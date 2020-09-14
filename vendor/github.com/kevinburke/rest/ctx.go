// +build go1.7

package rest

import (
	"context"
	"net/http"
)

type ctxVar int

var errCtx ctxVar = 0
var domainCtx ctxVar = 1

// CtxErr returns an error that's been stored in the Request context.
func CtxErr(r *http.Request) error {
	val := r.Context().Value(errCtx)
	if val == nil {
		return nil
	}
	return val.(error)
}

// ctxSetErr sets err in the request context, and returns the new request.
func ctxSetErr(r *http.Request, err error) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), errCtx, err))
}

// CtxDomain returns a domain that's been set on the request. Use it to get the
// domain set on a 401 error handler.
func CtxDomain(r *http.Request) string {
	val := r.Context().Value(domainCtx)
	if val == nil {
		return ""
	}
	return val.(string)
}

// ctxSetDomain sets a domain in the request context, and returns the new
// request.
func ctxSetDomain(r *http.Request, domain string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), domainCtx, domain))
}
