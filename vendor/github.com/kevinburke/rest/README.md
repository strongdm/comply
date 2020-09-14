# rest

This library contains a HTTP client, and a number of useful middlewares for
writing a HTTP client and server in Go. For more information and package
documentation, please [see the godoc documentation][gddo].

### Client

The `Client` struct makes it easy to interact with a JSON API.

```go
client := rest.NewClient("username", "password", "http://ipinfo.io")
req, _ := client.NewRequest("GET", "/json", nil)
type resp struct {
    City string `json:"city"`
    Ip   string `json:"ip"`
}
var r resp
client.Do(req, &r)
fmt.Println(r.Ip)
```

### Transport

Use the `rest.Transport` as the `http.Transport` to easily inspect the raw HTTP
request and response. Set `DEBUG_HTTP_TRAFFIC=true` in your environment to dump
HTTP requests and responses to stderr.

### Defining Custom Error Responses

`rest` exposes a number of HTTP error handlers - for example,
`rest.ServerError(w, r, err)` will write a 500 server error to w. By default,
these error handlers will write a generic JSON response over the wire, using
fields specified by the [HTTP problem spec][spec].

You can define a custom error handler if you like (say if you want to return
a HTML server error, or 404 error or similar) by calling RegisterHandler:

```go
rest.RegisterHandler(500, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    err := rest.CtxErr(r)
    fmt.Println("Server error:", err)
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(500)
    w.Write([]byte("<html><body>Server Error</body></html>"))
}))
```

[spec]: https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-03

### Debugging

Set the `DEBUG_HTTP_TRAFFIC` environment variable to print out all
request/response traffic being made by the client.

`rest` also includes a `Transport` that is a drop in for a `http.Transport`,
but includes support for debugging HTTP requests. Add it like so:

```go
client := http.Client{
    Transport: &rest.Transport{
        Debug: true,
        Output: os.Stderr,
        Transport: http.DefaultTransport,
    },
}
```

## Donating

Donations free up time to make improvements to the library, and respond to
bug reports. You can send donations via Paypal's "Send Money" feature to
kev@inburke.com. Donations are not tax deductible in the USA.

[gddo]: https://godoc.org/github.com/kevinburke/rest
