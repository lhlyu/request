package request

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Request struct {
	c *http.Client
	r *http.Request
}

type Option struct {
	Method string
	Url    string
	Param  interface{} // map[string]interface{}
	Data   interface{} // map[string]interface{}

}

const (
	DURATION                  = time.Second
	DEFAULT_RESPONSE_TIME_OUT = DURATION * 5 // 默认超时时间 5s
)

func NewRequestDefualt() *Request {
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, DEFAULT_RESPONSE_TIME_OUT)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(DEFAULT_RESPONSE_TIME_OUT)) // 设置与连接关联的读写期限,0表示不会超时
				return conn, nil
			},
			ResponseHeaderTimeout: DEFAULT_RESPONSE_TIME_OUT, // 响应时间,0表示不会超时
		},
	}
	r := &http.Request{}
	return &Request{
		c: c,
		r: r,
	}
}

func NewRequest(client *http.Client) *Request {
	return &Request{
		c: client,
	}
}

func (*Request) getDuration(times int) time.Duration {
	return DURATION * time.Duration(times)
}

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
func (rq *Request) Ask(method, url string, header interface{}) *Request {
	return rq
}

// get request
func (rq *Request) Get(apiUrl string, params ...interface{}) *Request {
	return rq
}

// post request
func (rq *Request) Post(apiUrl string, datas ...interface{}) *Request {
	return rq
}
