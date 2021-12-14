package testclient

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mikeschinkel/go-only"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

// NewTestContext creates a new *TextContext from a TextContextArgs{}
func NewTestContext(args TestContextArgs) *TestContext {
	tc := TestContext{
		T:                args.TestManager,
		ExpectedResponse: args.ExpectedResponse,
		TestNum:          args.TestNum,
		ValidateBodyFunc: args.ValidateBodyFunc,
	}
	if args.TestURL == "" {
		args.TestURL = "http://localhost"
	}
	if args.StatusCode == 0 {
		args.StatusCode = 200
	}
	if tc.ExpectedResponse == nil {
		tc.ExpectedResponse = NewExpectedResponse(args.TestURL, args.StatusCode)
	}
	if args.Fixture != nil {
		tc.ExpectedResponse.Fixture = args.Fixture
	} else if args.Filename != "" {
		tc.ExpectedResponse.Fixture = NewFixture(tc.T, args.Filename)
	}
	return &tc
}

type ValidateBodyFunc func([]byte) (interface{}, error)

// TestContext embedded *testing.T and adds a ExpectedResponse
type TestContext struct {
	*testing.T
	ExpectedResponse *ExpectedResponse
	TestNum          int
	ValidateBodyFunc ValidateBodyFunc
}

// TestContextArgs is used to pass arguments to NewTestContext()
type TestContextArgs struct {
	TestManager      *testing.T
	ExpectedResponse *ExpectedResponse
	TestNum          int
	ValidateBodyFunc ValidateBodyFunc
	TestURL          string
	StatusCode       int
	Fixture          *Fixture
	Filename         string
}

// Error calls t.Errorf() and includes the URL in TextContext.ExpectedResponse
func (tc *TestContext) Error(message string, err error) {
	t := tc.T
	t.Helper()
	t.Errorf("%s: %s",
		fmt.Sprintf("%s from %s", message, tc.ExpectedResponse.URL),
		err)
}

// RefreshJSONFixture
func (tc *TestContext) RefreshJSONFixture() {
	for range only.Once {

		rt := NewRefreshTransport(tc.ExpectedResponse)
		client := NewClientWithTransport(rt)
		response, err := tc.TestJSONGet(client)
		if err != nil {
			tc.Error("Failed to HTTP GET", err)
			break
		}

		body := response.Body.(string)

		var marker = uuid.New().String()
		response.Body = marker

		data, err := json.MarshalIndent(response, "", "\t")
		if err != nil {
			tc.Error("Failed to marshal response JSON", err)
			break
		}

		var repos interface{}
		repos, err = tc.ValidateBodyFunc([]byte(body))
		if err != nil {
			tc.Error("Failed to unmarshal JSON from body", err)
		}

		bodydata, err := json.MarshalIndent(repos, "   ", "\t")
		if err != nil {
			tc.Error("Failed to marshal body data JSON", err)
			break
		}

		data = []byte(strings.Replace(
			string(data),
			fmt.Sprintf("\"%s\"", marker),
			string(bodydata),
			1))

		err = ioutil.WriteFile(tc.ExpectedResponse.Filepath(), data, os.ModePerm)
		if err != nil {
			tc.Error(fmt.Sprintf("Failed to write json to %s", tc.ExpectedResponse.Fixture.Filename), err)
			break
		}

	}
}

// TestJSONGet
func (tc *TestContext) TestJSONGet(client *Client) (resp *ExpectedResponse, err error) {
	for range only.Once {
		var hr *http.Response
		testURL := tc.ExpectedResponse.URL
		tc.Run(fmt.Sprintf("Get_URL_via_HTTP[%s]", testURL), func(t *testing.T) {
			for range only.Once {
				hr, err = client.GET(testURL, nil)
				if ignoreConnectionRefused(tc, err) {
					err = nil
					break
				}
				if err == nil {
					break
				}
			}
		})
		if hr == nil {
			err = nil
			break
		}

		tc.Run(fmt.Sprintf("Check_StatusCode[%s]", testURL), func(t *testing.T) {
			for range only.Once {

				if hr.StatusCode != tc.ExpectedResponse.StatusCode {
					tc.Error(
						fmt.Sprintf("Got status code %d, expected %d",
							hr.StatusCode,
							tc.ExpectedResponse.StatusCode),
						err)
				}
			}
		})

		if hr.StatusCode != 200 {
			break
		}

		resp = tc.ExpectedResponse

		var body []byte
		tc.Run(fmt.Sprintf("Read_Body[%s]", testURL), func(t *testing.T) {
			body, err = ioutil.ReadAll(hr.Body)
			if err != nil {
				tc.Error("Failed to read body", err)
			}
			if len(body) == 0 {
				tc.Error("Failed due to empty body returned", err)
			}
			resp.Body = string(body)
		})

		if cts, ok := hr.Header[ContentTypeHeader]; ok {
			switch len(cts) {
			case 1:
				resp.ContentType = cts[0]
			default:
				tc.Error(
					fmt.Sprintf("no '%s' header", ContentTypeHeader),
					fmt.Errorf("!"))
			}
		}

		tc.Run(fmt.Sprintf("Body_Is_Valid_JSON[%s]", testURL), func(t *testing.T) {
			_, err = tc.ValidateBodyFunc(body)
			if err != nil {
				tc.Error("Failed to unmarshal JSON from body", err)
			}
		})

	}
	return resp, err
}

// ignoreConnectionRefused
func ignoreConnectionRefused(tc *TestContext, err error) (ignore bool) {
	return tc.ExpectedResponse.StatusCode == 0 && strings.Contains(err.Error(), "connection refused")
}
