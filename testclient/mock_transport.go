package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"io"
	"io/ioutil"
	"net/http"
)

const ContentTypeHeader = "Content-Type"

// MockTransport provides a mock for RoundTrip() for testing
type MockTransport struct {
	ExpectedResponse *ExpectedResponse
}

// RoundTrip mocks the HTTP request-response process.
//goland:noinspection GoUnusedParameter
func (mt *MockTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	for range only.Once {
		header := make(http.Header, 0)
		header.Add(ContentTypeHeader, mt.ExpectedResponse.ContentType)
		var body io.ReadCloser
		body, err = mt.LoadBody()
		if err != nil {
			break
		}
		res = &http.Response{
			StatusCode: mt.ExpectedResponse.StatusCode,
			Header:     header,
			// See https://gist.github.com/crgimenes/92d851b944ca2e459da7daa5c44801bf
			Body: body,
		}

	}
	return res, err
}

func (mt *MockTransport) LoadBody() (rc io.ReadCloser, err error) {
	var action string
	for range only.Once {
		fp := mt.ExpectedResponse.Filepath()
		var b []byte
		b, err = ioutil.ReadFile(fp)
		if err != nil {
			action = "read"
			break
		}
		r := struct {
			Body interface{} `json:"body"`
		}{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			action = "load JSON"
			break
		}
		var body []byte
		body, err = json.Marshal(r.Body)
		if err != nil {
			action = "marshal body to JSON"
			break
		}
		rc = ioutil.NopCloser(bytes.NewReader(body))
	}
	if err != nil {
		err = fmt.Errorf("unable to %s from %s: %s",
			action,
			mt.ExpectedResponse.Filepath(),
			err)

	}
	return rc, err
}
