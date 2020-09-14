// Package rest implements responses and a HTTP client for API consumption.
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	log "github.com/inconshreveable/log15"
)

const jsonContentType = "application/json; charset=utf-8"

// Logger logs information about incoming requests.
var Logger log.Logger = log.New()

// Error implements the HTTP Problem spec laid out here:
// https://tools.ietf.org/html/rfc7807
type Error struct {
	// The main error message. Should be short enough to fit in a phone's
	// alert box. Do not end this message with a period.
	Title string `json:"title"`

	// Id of this error message ("forbidden", "invalid_parameter", etc)
	ID string `json:"id"`

	// More information about what went wrong.
	Detail string `json:"detail,omitempty"`

	// Path to the object that's in error.
	Instance string `json:"instance,omitempty"`

	// Link to more information about the error (Zendesk, API docs, etc).
	Type string `json:"type,omitempty"`

	// HTTP status code of the error.
	Status int `json:"status,omitempty"`
}

func (e *Error) Error() string {
	return e.Title
}

func (e *Error) String() string {
	if e.Detail != "" {
		return fmt.Sprintf("rest: %s. %s", e.Title, e.Detail)
	} else {
		return fmt.Sprintf("rest: %s", e.Title)
	}
}

var handlerMap = make(map[int]http.Handler)
var handlerMu sync.RWMutex

// RegisterHandler registers the given HandlerFunc to serve HTTP requests for
// the given status code. Use CtxErr and CtxDomain to retrieve extra values set
// on the request in f (if any).
//
// Despite registering the handler for the code, f is responsible for calling
// WriteHeader(code) since it may want to set response headers first.
//
// To delete a Handler, call RegisterHandler with nil for the second argument.
func RegisterHandler(code int, f http.Handler) {
	handlerMu.Lock()
	defer handlerMu.Unlock()
	switch f {
	case nil:
		delete(handlerMap, code)
	default:
		handlerMap[code] = f
	}
}

// ServerError logs the error to the Logger, and then responds to the request
// with a generic 500 server error message. ServerError panics if err is nil.
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusInternalServerError]
	handlerMu.RUnlock()
	if ok {
		r = ctxSetErr(r, err)
		f.ServeHTTP(w, r)
	} else {
		defaultServerError(w, r, err)
	}
}

var serverError = Error{
	Status: http.StatusInternalServerError,
	ID:     "server_error",
	Title:  "Unexpected server error. Please try again",
}

func defaultServerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		panic("rest: no error to log")
	}
	Logger.Error("Server error", "code", 500, "method", r.Method, "path", r.URL.Path, "err", err)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(serverError); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 500, "err", err)
	}
}

var notFound = Error{
	Title:  "Resource not found",
	ID:     "not_found",
	Status: http.StatusNotFound,
}

// NotFound returns a 404 Not Found error to the client.
func NotFound(w http.ResponseWriter, r *http.Request) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusNotFound]
	handlerMu.RUnlock()
	if ok {
		f.ServeHTTP(w, r)
	} else {
		defaultNotFound(w, r)
	}
}

func defaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusNotFound)
	nf := notFound
	nf.Instance = r.URL.Path
	if err := json.NewEncoder(w).Encode(nf); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 404, "err", err)
	}
}

// BadRequest logs a 400 error and then returns a 400 response to the client.
func BadRequest(w http.ResponseWriter, r *http.Request, err *Error) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusBadRequest]
	handlerMu.RUnlock()
	if ok {
		r = ctxSetErr(r, err)
		f.ServeHTTP(w, r)
	} else {
		defaultBadRequest(w, r, err)
	}
}

var gone = Error{
	Title:  "Resource is gone",
	ID:     "gone",
	Status: http.StatusGone,
}

// Gone responds to the request with a 410 Gone error message
func Gone(w http.ResponseWriter, r *http.Request) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusGone]
	handlerMu.RUnlock()
	if ok {
		f.ServeHTTP(w, r)
	} else {
		defaultGone(w, r)
	}
}

func defaultGone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusGone)
	g := gone
	g.Instance = r.URL.Path
	if err := json.NewEncoder(w).Encode(g); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 404, "err", err)
	}
}

func defaultBadRequest(w http.ResponseWriter, r *http.Request, err *Error) {
	if err == nil {
		panic("rest: no error to write")
	}
	if err.Status == 0 {
		err.Status = http.StatusBadRequest
	}
	Logger.Info("Bad request", "code", 400, "method", r.Method, "path", r.URL.Path, "err", err)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 400, "err", err)
	}
}

var notAllowed = Error{
	Title:  "Method not allowed",
	ID:     "method_not_allowed",
	Status: http.StatusMethodNotAllowed,
}

var authenticate = Error{
	Title:  "Unauthorized. Please include your API credentials",
	ID:     "unauthorized",
	Status: http.StatusUnauthorized,
}

// NotAllowed returns a generic HTTP 405 Not Allowed status and response body
// to the client.
func NotAllowed(w http.ResponseWriter, r *http.Request) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusMethodNotAllowed]
	handlerMu.RUnlock()
	if ok {
		f.ServeHTTP(w, r)
	} else {
		defaultNotAllowed(w, r)
	}
}

func defaultNotAllowed(w http.ResponseWriter, r *http.Request) {
	e := notAllowed
	e.Instance = r.URL.Path
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 405, "err", err)
	}
}

// Forbidden returns a 403 Forbidden status code to the client, with the given
// Error object in the response body.
func Forbidden(w http.ResponseWriter, r *http.Request, err *Error) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusForbidden]
	handlerMu.RUnlock()
	if ok {
		r = ctxSetErr(r, err)
		f.ServeHTTP(w, r)
	} else {
		defaultForbidden(w, r, err)
	}
}

func defaultForbidden(w http.ResponseWriter, r *http.Request, err *Error) {
	if err.ID == "" {
		err.ID = "forbidden"
	}
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 403, "err", err)
	}
}

// NoContent returns a 204 No Content message.
func NoContent(w http.ResponseWriter) {
	// No custom handler since there's no custom behavior.
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

// Unauthorized sets the Domain in the request context
func Unauthorized(w http.ResponseWriter, r *http.Request, domain string) {
	handlerMu.RLock()
	f, ok := handlerMap[http.StatusUnauthorized]
	handlerMu.RUnlock()
	if ok {
		r = ctxSetDomain(r, domain)
		f.ServeHTTP(w, r)
	} else {
		defaultUnauthorized(w, r, domain)
	}
}

func defaultUnauthorized(w http.ResponseWriter, r *http.Request, domain string) {
	err := authenticate
	err.Instance = r.URL.Path
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, domain))
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		Logger.Info("Couldn't write error", "path", r.URL.Path, "code", 401, "err", err)
	}
}
