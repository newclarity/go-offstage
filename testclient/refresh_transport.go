package testclient

import (
	"net/http"
)

// MockTransport provides a mock for RoundTrip() for testing
type RefreshTransport struct {
	http.RoundTripper
	ExpectedResponse *ExpectedResponse
}

func NewRefreshTransport(er *ExpectedResponse) *RefreshTransport {
	return &RefreshTransport{
		RoundTripper:     http.DefaultTransport,
		ExpectedResponse: er,
	}
}
