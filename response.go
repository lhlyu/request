package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type IResponse interface {
	GetRequest() *http.Request          // 获取请求体
	GetResponse() *http.Response        // 获取响应体
	GetBody() string                    // 获取body内容
	GetStatus() int                     // 获取响应状态码
	BodyUnmarshal(v interface{}) error  // 数据解析
}

type response struct {
	resp       *http.Response
}

func (r *response) GetRequest() *http.Request{
	return r.resp.Request
}

func (r *response) GetResponse() *http.Response{
	return r.resp
}

func (r *response) GetBody() string{
	bts,_ := ioutil.ReadAll(r.resp.Body)
	return string(bts)
}

func (r *response) GetStatus() int{
	return r.resp.StatusCode
}

func (r *response) BodyUnmarshal(v interface{}) error{
	bts,_ := ioutil.ReadAll(r.resp.Body)
	return json.Unmarshal(bts,v)
}