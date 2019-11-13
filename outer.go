package request

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// 设置超时时间,超时包括连接时间，任何重定向和读取响应正文,0表示不会超时
func (rq *Request) SetTimeOut(second int) *Request {
	rq.c.Timeout = rq.getDuration(second)
	return rq
}

// 设置 Transport
func (rq *Request) SetTransport(transport *http.Transport) *Request {
	if transport == nil {
		return rq
	}
	rq.c.Transport = transport
	return rq
}

// 设置request
func (rq *Request) SetHttpRequest(request *http.Request) *Request {
	if request == nil {
		return rq
	}
	rq.r = request
	return rq
}

// 设置header
func (rq *Request) SetHeader(header interface{}) *Request {

	return rq
}

// 设置cookie
func (rq *Request) SetCookie(header interface{}) *Request {

	return rq
}

// 设置cookiejar
func (rq *Request) SetCookieJar(header interface{}) *Request {
	return rq
}

// 通用请求
func (rq *Request) Ask(method, apiUrl string, option *Option) {
	if option == nil {
		option = rq.p
	}
	u, err := rq.Qs(apiUrl, option.Param)
	if err != nil {
		log.Println(err)
		return
	}
	rq.r.Method = method
	rq.r.URL = u
	rq.r.Header = InterToHeader(option.Header)

	rq.c.Timeout = rq.getDuration(option.TimeOut)

	if option.Proxy != "" {
		rq.t.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(option.Proxy)
		}
	}

	resp, err := rq.c.Do(rq.r)
	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bts), err)
	return
}
