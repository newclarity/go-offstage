package testclient

type ExpectedResponse struct {
	URL         string      `json:"url"`
	StatusCode  int         `json:"status_code"`
	ContentType string      `json:"content_type"`
	Body        interface{} `json:"body"`
	Fixture     *Fixture    `json:"-"`
}

func NewExpectedResponse(u string, sc int) *ExpectedResponse {
	er := ExpectedResponse{
		URL:        u,
		StatusCode: sc,
	}
	return &er
}
func (tr *ExpectedResponse) Filepath() string {
	if tr.Fixture == nil {
		return ""
	}
	return tr.Fixture.Filepath()
}
