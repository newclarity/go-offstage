package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"net/http"
	"net/url"
)

// Client provides an instance of the RequestDoer interface
type Client struct {
	*http.Client
}

func NewClient() *Client {
	return NewClientWithTransport(http.DefaultTransport)
}

func NewClientWithTransport(rt http.RoundTripper) *Client {
	return &Client{
		Client: &http.Client{
			Transport: rt,
		},
	}
}

// GET sends an HTTP(S) GET request to the URL with provided headers.
func (c *Client) GET(url string, headers http.Header) (resp *http.Response, err error) {
	return c.do(http.MethodGet, url, headers, nil)
}

// PUT sends an an HTTP(S) PUT to the URL with the body with provided headers
func (c *Client) PUT(url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	return c.requestWithBody(http.MethodPut, url, body, headers)
}

// POST sends an HTTP(S) POST request to the URL with the body and provided headers
func (c *Client) POST(url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	return c.requestWithBody(http.MethodPost, url, body, headers)
}

// DELETE sends an HTTP(S) DELETE request to the URL with provided headers.
func (c *Client) DELETE(url string, headers http.Header) (resp *http.Response, err error) {
	return c.do(http.MethodDelete, url, headers, nil)
}

// requestWithBody sends either an HTTP(S) POST or PUT request to the URL with the body and provided headers
func (c *Client) requestWithBody(method, url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	var jsonBytes []byte
	for range only.Once {

		var scheme string
		scheme, err = getURLScheme(url)
		if err != nil {
			break
		}

		jsonBytes, err = json.Marshal(body)
		if err != nil {
			err = fmt.Errorf("unable to marshal JSON body for %s POST request", scheme)
			break
		}

		resp, err = c.do(method, url, headers, bytes.NewReader(jsonBytes))

	}
	return resp, err
}

// do calls on http client to do an HTTP(S) request
//goland:noinspection GoUnusedParameter
func (c *Client) do(method, url string, headers http.Header, body interface{}) (resp *http.Response, err error) {
	for range only.Once {
		var scheme string
		scheme, err = getURLScheme(url)
		if err != nil {
			break
		}

		var req *http.Request
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			err = fmt.Errorf("unable to instantiate %s %s request",
				scheme,
				method,
			)
			break
		}

		req.Header = headers
		resp, err = c.Do(req)
		if err != nil {
			err = fmt.Errorf("unable to perform %s %s request: %s",
				scheme,
				method,
				err,
			)
		}

	}
	return resp, err
}

// getURLScheme returns the scheme from a URL
func getURLScheme(u string) (s string, err error) {
	for range only.Once {
		uo, err := url.Parse(u)
		if err != nil {
			err = fmt.Errorf("unable to parse URL '%s': %s",
				u,
				err,
			)
		}
		s = uo.Scheme
	}
	return s, err
}
