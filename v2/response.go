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
	async     bool
	asyncBody chan string
	err       error
}

func (r *Response) check() {
	if r.async {
		r.Await()
	}
}

func (r *Response) GetRequest() *http.Request {
	r.check()
	if r.resp.Body == nil {
		r.resp.Body = ioutil.NopCloser(strings.NewReader(r.reqBody))
	}
	return r.req
}

func (r *Response) GetResponse() *http.Response {
	r.check()
	if r.resp.Body == nil {
		r.resp.Body = ioutil.NopCloser(strings.NewReader(r.respBody))
	}
	return r.resp
}

func (r *Response) GetBody() string {
	r.check()
	return r.respBody
}

func (r *Response) GetStatus() int {
	r.check()
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

func (r *Response) Await() IResponse {
	if r.async {
		r.respBody = <-r.asyncBody
		r.async = false
	}
	return r
}

func (r *Response) Error() error {
	return r.err
}
