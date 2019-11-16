package request

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	body       []byte

}

func newResponse(resp *http.Response,err error) IResponse{
	r := &response{
		resp:resp,
		body:make([]byte,0),
	}
	if err != nil{
		log.Println("request error :",err)
		return r
	}
	defer resp.Body.Close()
	bts,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("read body error :",err)
		return r
	}
	r.body = bts
	return r
}

func (r *response) GetRequest() *http.Request{
	return r.resp.Request
}

func (r *response) GetResponse() *http.Response{
	return r.resp
}

func (r *response) GetBody() string{
	return string(r.body)
}

func (r *response) GetStatus() int{
	return r.resp.StatusCode
}

func (r *response) BodyUnmarshal(v interface{}) error{
	return json.Unmarshal(r.body,v)
}