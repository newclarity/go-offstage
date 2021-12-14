package testclient

import (
	"encoding/json"
	"flag"
	"fmt"
	"testing"
)

const testURL = "https://api.github.com/users/google/repos"

type repo struct {
	RepoId      int    `json:"id"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
}

var reposFunc = func(body []byte) (interface{}, error) {
	repos := make([]repo, 0)
	err := json.Unmarshal(body, &repos)
	return repos, err
}

var refresh = flag.Bool("refresh", false, "Refresh recorded results")

const ResponseAsRecordedFilename = "response-as-recorded.json"

func getExpectedResponses(t *testing.T) []*ExpectedResponse {
	return []*ExpectedResponse{
		{
			Fixture:     NewFixture(t, ResponseAsRecordedFilename),
			URL:         testURL,
			StatusCode:  200,
			ContentType: "application/json; charset=utf-8",
		},
		{
			URL:        fmt.Sprintf("%s/not-there", testURL),
			StatusCode: 404,
		},
		{
			URL:        "http://localhost/not-there",
			StatusCode: 0,
		},
	}
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestMockTransportGet(t *testing.T) {
	if !*refresh {
		t.Log("-refresh=false: Skipping refresh of recorded results")
	} else {
		tc := NewTestContext(TestContextArgs{
			TestManager:      t,
			TestURL:          testURL,
			ValidateBodyFunc: reposFunc,
			Filename:         "response-as-recorded.json",
		})
		tc.RefreshJSONFixture()
	}
	var testNum = 1
	for _, er := range getExpectedResponses(t) {
		tc := NewTestContext(TestContextArgs{
			TestManager:      t,
			ExpectedResponse: er,
			TestNum:          testNum,
			ValidateBodyFunc: reposFunc,
		})
		_, err := tc.TestJSONGet(NewClientWithTransport(&MockTransport{
			ExpectedResponse: er,
		}))
		if err != nil {
			tc.Error("Failed to HTTP GET", err)
		}
		testNum++
	}
}
