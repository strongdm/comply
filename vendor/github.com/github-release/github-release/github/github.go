// Package github is a mini-library for querying the GitHub v3 API that
// takes care of authentication (with tokens only) and pagination.
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/kevinburke/rest"
	"github.com/tomnomnom/linkheader"
)

const DefaultBaseURL = "https://api.github.com"

// Set to values > 0 to control verbosity, for debugging.
var VERBOSITY = 0

// DoAuthRequest ...
//
// TODO: This function is amazingly ugly (separate headers, token, no API
// URL constructions, et cetera).
func DoAuthRequest(method, url, mime, token string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := newAuthRequest(method, url, mime, token, headers, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Client collects a few options that can be set when contacting the GitHub
// API, such as authorization tokens. Methods called on Client will supply
// these options when calling the API.
type Client struct {
	client *rest.Client
}

// NewClient creates a new Client for use with the Github API.
func NewClient(username, token string, client *rest.Client) Client {
	c := Client{}
	if client == nil {
		c.client = rest.NewClient(username, token, DefaultBaseURL)
	} else {
		c.client = client
	}
	return c
}

// SetBaseURL updates the client's base URL, if baseurl is a non-empty value.
func (c Client) SetBaseURL(baseurl string) {
	// This is lazy, because the caller always tries to override the base URL
	// with EnvApiEndpoint and we used to ignore that at the top of Get, but we
	// don't do that now. So instead just filter out/ignore the empty string.
	if baseurl == "" {
		return
	}
	c.client.Base = baseurl
}

// Get fetches uri (relative URL) from the GitHub API and unmarshals the
// response into v. It takes care of pagination transparantly.
func (c Client) Get(uri string, v interface{}) error {
	rc, err := c.getPaginated(uri)
	if err != nil {
		return err
	}
	defer rc.Close()
	var r io.Reader = rc
	if VERBOSITY > 0 {
		vprintln("BODY:")
		r = io.TeeReader(rc, os.Stderr)
	}

	// Github may return paginated responses. If so, githubGetPaginated will
	// return a reader which yields the concatenation of all pages. These
	// reponses are _separate_ JSON arrays. Standard json.Unmarshal() or
	// json.Decoder.Decode() will not have the expected result when
	// unmarshalling into v. For example, a 2-page response:
	//
	//   1. [{...}, {...}, {...}]
	//   2. [{...}]
	//
	// If v is a slice type, we'd like to decode the four objects from the
	// two pages into a single slice. However, if we just use
	// json.Decoder.Decode(), that won't work. v will be overridden each
	// time.
	//
	// For this reason, we use two very ugly things.
	//
	//   1. We analyze v with reflect to see if it's a slice.
	//   2. If so, we use the json.Decoder token API and reflection to
	//      dynamically add new elements into the slice, ignoring the
	//      boundaries between JSON arrays.
	//
	// This is a lot of work, and feels very stupid. An alternative would be
	// removing the outermost ][ in the intermediate responses, which would
	// be even more finnicky. Another alternative would be to explicitly
	// expose a pagination API, forcing clients of this code to deal with
	// it. That's how the go-github library does it. But why solve a problem
	// sensibly if one can power through it with reflection (half-joking)?

	sl := reflect.Indirect(reflect.ValueOf(v)) // Get the reflect.Value of the slice so we can append to it.
	t := sl.Type()
	if t.Kind() != reflect.Slice {
		// Not a slice, not going to handle special pagination JSON stream
		// semantics since it likely wouldn't work properly anyway. If this
		// is a non-paginated stream, it should work.
		return json.NewDecoder(r).Decode(v)
	}
	t = t.Elem() // Extract the type of the slice's elements.

	// Use streaming Token API to append all elements of the JSON stream
	// arrays (pagination) to the slice.
	for dec := json.NewDecoder(r); ; {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil // Natural end of the JSON stream.
			}
			return err
		}
		vprintf("TOKEN %T: %v\n", tok, tok)
		// Check for tokens until we get an opening array brace. If we're
		// not in an array, we can't decode an array element later, which
		// would result in an error.
		if tok != json.Delim('[') {
			continue
		}

		// Read the array, appending all elements to the slice.
		for dec.More() {
			it := reflect.New(t) // Interface to a valid pointer to an object of the same type as the slice elements.
			if err := dec.Decode(it.Interface()); err != nil {
				return err
			}
			vprintf("OBJECT %T: %v\n", it.Interface(), it)
			sl.Set(reflect.Append(sl, it.Elem()))
		}
	}
}

var defaultHttpClient *http.Client

func init() {
	defaultHttpClient = &http.Client{
		Transport: rest.DefaultTransport,
	}
}

// Caller is responsible for reading and closing the response body.
func (c Client) do(r *http.Request) (*http.Response, error) {
	// Pulled this out of client.go:Do because we need to read the response
	// headers.
	var res *http.Response
	var err error
	if c.client.Client == nil {
		res, err = defaultHttpClient.Do(r)
	} else {
		res, err = c.client.Client.Do(r)
	}
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		// both of these consume res.Body
		if c.client.ErrorParser != nil {
			return nil, c.client.ErrorParser(res)
		}
		return nil, rest.DefaultErrorParser(res)
	}
	return res, nil
}

const uaPart = "github-release/" + VERSION

func (c Client) NewRequest(method, uri string, body io.Reader) (*http.Request, error) {
	req, err := c.client.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	ua := req.Header.Get("User-Agent")
	if ua == "" {
		req.Header.Set("User-Agent", uaPart)
	} else {
		req.Header.Set("User-Agent", uaPart+" "+ua)
	}
	return req, nil
}

// getPaginated returns a reader that yields the concatenation of the
// paginated responses to a query (URI).
//
// TODO: Rework the API so we can cleanly append per_page=100 as a URL
// parameter.
func (c Client) getPaginated(uri string) (io.ReadCloser, error) {
	// Parse the passed-in URI to make sure we don't lose any values when
	// setting our own params.
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	v := u.Query()
	v.Set("per_page", "100") // The default is 30, this makes it less likely for Github to rate-limit us.
	u.RawQuery = v.Encode()
	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	vprintln("GET (top-level)", resp.Request.URL, "->", resp)

	// If the HTTP response is paginated, it will contain a Link header.
	links := linkheader.Parse(resp.Header.Get("Link"))
	if len(links) == 0 {
		return resp.Body, nil // No pagination.
	}

	// In this case, fetch all pages and concatenate them.
	r, w := io.Pipe()
	done := make(chan struct{})               // Backpressure from the pipe writer.
	responses := make(chan *http.Response, 5) // Allow 5 concurrent HTTP requests.
	responses <- resp

	// URL fetcher goroutine. Fetches paginated responses until no more
	// pages can be found. Closes the write end of the pipe if fetching a
	// page fails.
	go func() {
		defer close(responses) // Signal that no more requests are coming.
		for len(links) > 0 {
			nextLinkURL := nextLink(links)
			if nextLinkURL == "" {
				return // We're done.
			}

			req, err := c.NewRequest("GET", nextLinkURL, nil)
			if err != nil {
				w.CloseWithError(err)
				return
			}
			resp, err := c.do(req)
			if err != nil {
				w.CloseWithError(err)
				return
			}
			links = linkheader.Parse(resp.Header.Get("Link"))
			if err != nil {
				w.CloseWithError(err)
				return
			}
			select {
			case <-done:
				return // The body concatenator goroutine signals it has stopped.
			case responses <- resp: // Schedule the request body to be written to the pipe.
			}
		}
	}()

	// Body concatenator goroutine. Writes each response into the pipe
	// sequentially. Closes the write end of the pipe if the HTTP status is
	// not 200 or the body can't be read.
	go func() {
		defer func() {
			// Drain channel and close bodies, stop leaks.
			for resp := range responses {
				resp.Body.Close()
			}
		}()
		defer close(done) // Signal that we're done writing all requests, or an error occurred.
		for resp := range responses {
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				w.CloseWithError(fmt.Errorf("expected '200 OK' but received '%v' (url: %s)", resp.Status, resp.Request.URL))
				return
			}
			_, err := io.Copy(w, resp.Body)
			resp.Body.Close()
			if err != nil {
				vprintln("error: io.Copy: ", err)
				w.CloseWithError(err)
				return
			}
		}
		w.Close()
	}()

	return r, nil
}

// Create a new request that sends the auth token.
func newAuthRequest(method, url, mime, token string, headers map[string]string, body io.Reader) (*http.Request, error) {
	vprintln("creating request:", method, url, mime, token)

	var n int64 // content length
	var err error
	if f, ok := body.(*os.File); ok {
		// Retrieve the content-length and buffer up if necessary.
		body, n, err = materializeFile(f)
		if err != nil {
			return nil, err
		}
	}

	// TODO find all of the usages and replace with the Client.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", token)
	req.Header.Set("User-Agent", uaPart)

	// net/http automatically does this if req.Body is of type
	// (bytes.Reader|bytes.Buffer|strings.Reader). Sadly, we also need to
	// handle *os.File.
	if n != 0 {
		vprintln("setting content-length to", n)
		req.ContentLength = n
	}

	if mime != "" {
		req.Header.Set("Content-Type", mime)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// nextLink returns the HTTP header Link annotated with 'next', "" otherwise.
func nextLink(links linkheader.Links) string {
	for _, link := range links {
		if link.Rel == "next" && link.URL != "" {
			return link.URL
		}
	}
	return ""
}
