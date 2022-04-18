package testcase

import (
	"github.com/newclarity/go-offstage"
)

var currentTestCase *TestCase

func InitiateNewTestCase() {
	offstage.InTesting = true
	currentTestCase = NewTestCase()
}

//goland:noinspection GoUnusedExportedFunction
func GetCurrentTestCase() *TestCase {
	return currentTestCase
}

func SetTestCaseHeader(name, format string, args ...interface{}) {
	currentTestCase.SetHeader(name, format, args...)
}

func SetTestCaseServer(s *offstage.Server) {
	currentTestCase.SetServer(s)
}

func DoTestCaseGET(path string) (r *Response, err error) {
	tc := currentTestCase
	return tc.GET(tc.GetURL(path), tc.Headers)
}

func DoTestCasePOST(path, body string) (r *Response, err error) {
	tc := currentTestCase
	return tc.POST(tc.GetURL(path), body, tc.Headers)
}

type Args struct {
	Organization  string
	Repository    string
	CommitterName string
}

func ConfigureServerFunc(args Args) func(s *offstage.Server) {
	return func(s *offstage.Server) {
		SetTestCaseServer(s)
		c := s.Context
		c.SetOrganizationName(args.Organization)
		c.SetRepositoryName(args.Repository)
		c.SetCommitterName(args.CommitterName)
	}
}

func StartTestCaseServer(start offstage.StartFunc, configure offstage.ConfigureFunc) {
	currentTestCase.Server = start(configure)
}

func StopTestCaseServer() {
	currentTestCase.Server.StopServer()
}
