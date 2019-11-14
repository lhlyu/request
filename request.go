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
	t *http.Transport
	p *Option
	b string // baseUrl
	u string // url
	m string // method
	d bool   // debug
}

type Option struct {
	Proxy    string      // http://xxx:8081   代理请求地址
	TimeOut  int         // second
	Files    []string    // file path
	Param    map[string]string
	Header   map[string]string
	Cookie   map[string]string
	Data     map[string]interface{}
	FormData map[string]string
}

type (
	MSI = map[string]interface{}
	MSS = map[string]string

)

func NewRequest() IRequest {
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
	//c.Transport = &http.Transport{
	//	Proxy: func(request *http.Request) (url *url.URL, e error) {
	//		return url.Parse(  "http://202.105.233.40:80")
	//	},
	//}
	r := &http.Request{}
	t := &http.Transport{}
	p := &Option{
		Param: make(map[string]string),
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Data: make(map[string]interface{}),
		FormData: make(map[string]string),
		Files:make([]string,0),
	}
	return &Request{
		c: c,
		r: r,
		t: t,
		p: p,
		m: GET,
	}
}

// https://www.cnblogs.com/mafeng/p/7068837.html
