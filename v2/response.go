package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Response struct {
	req       *http.Request
	resp      *http.Response
	reqBody   string
	respBody  string
	asyncBody chan string
	err       error
}

func (r *Response) GetRequest() *http.Request {
	if r.resp.Body == nil {
		r.resp.Body = ioutil.NopCloser(strings.NewReader(r.reqBody))
	}
	return r.req
}

func (r *Response) GetResponse() *http.Response {
	if r.resp.Body == nil {
		r.resp.Body = ioutil.NopCloser(strings.NewReader(r.respBody))
	}
	return r.resp
}

func (r *Response) GetBody() string {
	return r.respBody
}

func (r *Response) GetStatus() int {
	return r.resp.StatusCode
}

func (r *Response) IsStatus(status ...int) bool {
	if len(status) == 0 {
		return r.resp.StatusCode == 200
	}
	for _, v := range status {
		if v == r.resp.StatusCode {
			return true
		}
	}
	return false
}

func (r *Response) BodyUnmarshal(v interface{}) error {
	return json.Unmarshal([]byte(r.respBody), v)
}

func (r *Response) BodyCompile(pattern string) []string {
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(r.GetBody(), -1)
}

func (r *Response) Await() <-chan string {
	return r.asyncBody
}

func (r *Response) Error() error {
	return r.err
}
