package testclient

import (
	"io/ioutil"
	"testing"
)

func TestMockTransportLoad(t *testing.T) {
	f := NewFixture(t, "loadable.json")
	mt := &MockTransport{
		ExpectedResponse: &ExpectedResponse{Fixture: f},
	}
	r, err := mt.LoadBody()
	if err != nil {
		t.Errorf("Failed to load %s into MockTransport: %s",
			f.Filepath(),
			err)
	}
	var b []byte
	b, err = ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Failed to read %s: %s",
			f.Filepath(),
			err)
	}
	expected := "\"Hello World\""
	if string(b) != expected {
		t.Errorf("Expected '%s', got '%s'",
			expected,
			string(b))
	}
}
