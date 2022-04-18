package testcase

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"github.com/newclarity/go-offstage"
	"github.com/newclarity/go-offstage/testclient"
	"net/http"
)

type TestCase struct {
	Server  *offstage.Server
	Client  *testclient.Client
	Headers http.Header
}

// NewTestCase instantiates a new test case
func NewTestCase() (tc *TestCase) {
	return &TestCase{
		Client:  testclient.NewClient(),
		Headers: make(http.Header, 0),
	}
}

func (tc *TestCase) SetServer(s *offstage.Server) {
	tc.Server = s
}

func (tc *TestCase) SetServerPort(sp int) {
	tc.Server.Port = sp
}

func (tc *TestCase) SetHeader(name string, format string, args ...interface{}) {
	h := fmt.Sprintf(format, args...)
	hs, ok := tc.Headers[name]
	if !ok {
		tc.Headers[name] = []string{h}
	} else {
		tc.Headers[name] = append(hs, h)
	}
}

func (tc *TestCase) GET(url string, headers http.Header) (r *Response, err error) {
	var _r *http.Response
	for range only.Once {
		_r, err = tc.Client.GET(url, headers)
		if err != nil {
			// TODO Wrap error
			break
		}
		r = NewResponse(_r)
	}
	return r, err
}

func (tc *TestCase) POST(url string, body string, headers http.Header) (r *Response, err error) {
	var _r *http.Response
	for range only.Once {
		_r, err = tc.Client.POST(url, body, headers)
		if err != nil {
			// TODO Wrap error
			break
		}
		r = NewResponse(_r)
	}
	return r, err
}

func (tc *TestCase) GetPort() int {
	return tc.Server.Port
}

// GetURL the URL for this service
func (tc *TestCase) GetURL(path string) string {
	return fmt.Sprintf("http://localhost:%d%s",
		tc.GetPort(),
		path) // tc.GetPath()
}
