package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type UploadType string

// JSON specifies you'd like to upload JSON data.
var JSON UploadType = "application/json"

// FormURLEncoded specifies you'd like to upload form-urlencoded data.
var FormURLEncoded UploadType = "application/x-www-form-urlencoded"

const Version = "2.2"

var ua string

func init() {
	gv := strings.Replace(runtime.Version(), "go", "", 1)
	ua = fmt.Sprintf("rest-client/%s (https://github.com/kevinburke/rest) go/%s (%s/%s)",
		Version, gv, runtime.GOOS, runtime.GOARCH)
}

// Client is a generic Rest client for making HTTP requests.
type Client struct {
	// Username for use in HTTP Basic Auth
	ID string
	// Password for use in HTTP Basic Auth
	Token string
	// HTTP Client to use for making requests
	Client *http.Client
	// The base URL for all requests to this API, for example,
	// "https://fax.twilio.com/v1"
	Base string
	// Set UploadType to JSON or FormURLEncoded to control how data is sent to
	// the server. Defaults to FormURLEncoded.
	UploadType UploadType
	// ErrorParser is invoked when the client gets a 400-or-higher status code
	// from the server. Defaults to rest.DefaultErrorParser.
	ErrorParser func(*http.Response) error

	useBearerAuth bool
}

// NewClient returns a new Client with the given user and password. Base is the
// scheme+domain to hit for all requests.
func NewClient(user, pass, base string) *Client {
	return &Client{
		ID:          user,
		Token:       pass,
		Client:      defaultHttpClient,
		Base:        base,
		UploadType:  JSON,
		ErrorParser: DefaultErrorParser,
	}
}

// NewBearerClient returns a new Client configured to use Bearer authentication.
func NewBearerClient(token, base string) *Client {
	return &Client{
		ID:            "",
		Token:         token,
		Client:        defaultHttpClient,
		Base:          base,
		UploadType:    JSON,
		ErrorParser:   DefaultErrorParser,
		useBearerAuth: true,
	}
}

var defaultDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}

// DialSocket configures c to use the provided socket and http.Transport to
// dial a Unix socket instead of a TCP port.
//
// If transport is nil, the settings from DefaultTransport are used.
func (c *Client) DialSocket(socket string, transport *http.Transport) {
	dialSock := func(ctx context.Context, proto, addr string) (conn net.Conn, err error) {
		return defaultDialer.DialContext(ctx, "unix", socket)
	}
	if transport == nil {
		ht := http.DefaultTransport.(*http.Transport)
		transport = &http.Transport{
			Proxy:                 ht.Proxy,
			MaxIdleConns:          ht.MaxIdleConns,
			IdleConnTimeout:       ht.IdleConnTimeout,
			TLSHandshakeTimeout:   ht.TLSHandshakeTimeout,
			ExpectContinueTimeout: ht.ExpectContinueTimeout,
			DialContext:           dialSock,
		}
	}
	if c.Client == nil {
		// need to copy this so we don't modify the default client
		c.Client = &http.Client{
			Timeout: defaultHttpClient.Timeout,
		}
	}
	switch tp := c.Client.Transport.(type) {
	// TODO both of these cases clobbber the existing transport which isn't
	// ideal.
	case nil, *Transport:
		c.Client.Transport = &Transport{
			RoundTripper: transport,
			Debug:        DefaultTransport.Debug,
			Output:       DefaultTransport.Output,
		}
	case *http.Transport:
		c.Client.Transport = transport
	default:
		panic(fmt.Sprintf("could not set DialSocket on unknown transport: %#v", tp))
	}
}

// NewRequest creates a new Request and sets basic auth based on the client's
// authentication information.
func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	// see for example https://github.com/meterup/github-release/issues/1 - if
	// the path contains the full URL including the base, strip it out
	path = strings.TrimPrefix(path, c.Base)
	req, err := http.NewRequest(method, c.Base+path, body)
	if err != nil {
		return nil, err
	}
	switch {
	case c.useBearerAuth && c.Token != "":
		req.Header.Add("Authorization", "Bearer "+c.Token)
	case !c.useBearerAuth && (c.ID != "" || c.Token != ""):
		req.SetBasicAuth(c.ID, c.Token)
	}
	req.Header.Add("User-Agent", ua)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")
	if method == "POST" || method == "PUT" {
		uploadType := c.UploadType
		if uploadType == "" {
			uploadType = JSON
		}
		req.Header.Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", uploadType))
	}
	return req, nil
}

// Do performs the HTTP request. If the HTTP response is in the 2xx range,
// Unmarshal the response body into v. If the response status code is 400 or
// above, attempt to Unmarshal the response into an Error. Otherwise return
// a generic http error.
func (c *Client) Do(r *http.Request, v interface{}) error {
	var res *http.Response
	var err error
	if c.Client == nil {
		res, err = defaultHttpClient.Do(r)
	} else {
		res, err = c.Client.Do(r)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		if c.ErrorParser != nil {
			return c.ErrorParser(res)
		}
		return DefaultErrorParser(res)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if v == nil || res.StatusCode == http.StatusNoContent {
		return nil
	} else {
		return json.Unmarshal(resBody, v)
	}
}

// DefaultErrorParser attempts to parse the response body as a rest.Error. If
// it cannot do so, return an error containing the entire response body.
func DefaultErrorParser(resp *http.Response) error {
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rerr := new(Error)
	err = json.Unmarshal(resBody, rerr)
	if err != nil {
		return fmt.Errorf("invalid response body: %s", string(resBody))
	}
	if rerr.Title == "" {
		return fmt.Errorf("invalid response body: %s", string(resBody))
	} else {
		rerr.Status = resp.StatusCode
		return rerr
	}
}
