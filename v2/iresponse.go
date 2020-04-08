package v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type IResponse interface {
	GetRequest() *http.Request                // 获取请求体
	GetResponse() *http.Response              // 获取响应体
	GetBody() string                          // 获取body内容
	GetStatus() int                           // 获取响应状态码
	IsStatusOk() bool                         // 响应状态码 == 200
	AssertStatus(status int)                  // 断言状态码
	BodyUnmarshal(v interface{}) error        // 数据解析
	Then(func(resp *http.Response)) IResponse // 自定义处理   response
	Error() error                             // 返回请求错误
	OnSuccess(func(resp IResponse)) IResponse // 成功回调
	OnError(func(resp IResponse)) IResponse   // 失败回调
	BodyCompile(pattern string) []string      // 正则匹配body
	Print()                                   // 打印请求和响应
}

type response struct {
	resp *http.Response
	body []byte
	url  string
	err  error
}

func newResponse(resp *http.Response, url string, err error) IResponse {
	if resp == nil {
		resp = &http.Response{}
	}
	r := &response{
		resp: resp,
		body: make([]byte, 0),
		url:  url,
		err:  err,
	}
	if err != nil {
		log.Println("request error :", err)
		return r
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body error :", err)
		return r
	}
	r.body = bts
	return r
}

func (r *response) GetRequest() *http.Request {
	return r.resp.Request
}

func (r *response) GetResponse() *http.Response {
	r.resp.Body = ioutil.NopCloser(strings.NewReader(r.GetBody()))
	return r.resp
}

func (r *response) GetBody() string {
	return string(r.body)
}

func (r *response) GetStatus() int {
	return r.resp.StatusCode
}

func (r *response) IsStatusOk() bool {
	return r.resp.StatusCode == 200
}

func (r *response) AssertStatus(status int) {
	fmt.Println(r.resp.StatusCode == status)
}

func (r *response) BodyUnmarshal(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *response) Then(f func(resp *http.Response)) IResponse {
	rsp := r.GetResponse()
	rsp.Body = ioutil.NopCloser(strings.NewReader(r.GetBody()))
	f(rsp)
	return r
}

func (r *response) Error() error {
	return r.err
}

func (r *response) OnSuccess(f func(resp IResponse)) IResponse {
	status := r.GetStatus()
	if status >= 200 && status < 300 {
		f(r)
	}
	return r
}
func (r *response) OnError(f func(resp IResponse)) IResponse {
	status := r.GetStatus()
	if status < 200 || status >= 300 {
		f(r)
	}
	return r
}

func (r *response) BodyCompile(pattern string) []string {
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(r.GetBody(), -1)
}

func (r *response) Print() {
	fmt.Println("Request Header")
	fmt.Println(r.GetRequest().Method, r.url)
	for k, v := range r.GetRequest().Header {
		fmt.Printf("%s : %s\n", k, strings.Join(v, ";"))
	}
	fmt.Println()

	fmt.Println("Response Header")
	for k, v := range r.GetResponse().Header {
		fmt.Printf("%s : %s\n", k, strings.Join(v, ";"))
	}
	fmt.Println()
	fmt.Println("Response Body")
	fmt.Println(r.GetBody())
}

/**

1开头：信息状态码

2开头：成功状态码

3开头：重定向状态码

4开头：客户端错误状态码

5开头：服务端错误状态码

*/
