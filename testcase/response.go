package testcase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	Indent bool
}

func NewResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
		Indent:   true,
	}
}

func (r *Response) ReadBody() []byte {
	var body []byte
	body, _ = ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	return body
}

func (r *Response) UnmarshalTo(v interface{}) {
	_ = json.Unmarshal(r.ReadBody(), &v)
}

func (r *Response) ReadBodyToJSON() string {
	var v interface{}
	r.UnmarshalTo(&v)
	var body []byte
	if r.Indent {
		body, _ = json.MarshalIndent(v, "", "   ")
	} else {
		body, _ = json.Marshal(v)
	}
	return string(body)
}
