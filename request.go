package request

import (
	"net/http"
)

type Request struct {
	c *http.Client
	r *http.Request
	t *http.Transport
	p *Option
	b string // baseUrl
}

type Option struct {
	Proxy   string      // http://xxx:8081   代理请求地址
	TimeOut int         // second
	Param   interface{} // map[interface{}]interface{} or struct or queryString                 //
	Data    interface{} // map[interface{}]interface{} or struct or queryString or jsonString
	Header  interface{} // map[interface{}]interface{} or struct
	Cookie  interface{} // map[interface{}]interface{} or struct
}

func NewRequest() *Request {
	c := &http.Client{
		//Transport: &http.Transport{
		//	DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
		//		conn, err := net.DialTimeout(network, addr, DEFAULT_RESPONSE_TIME_OUT)
		//		if err != nil {
		//			return nil, err
		//		}
		//		conn.SetDeadline(time.Now().Add(DEFAULT_RESPONSE_TIME_OUT)) // 设置与连接关联的读写期限,0表示不会超时
		//		return conn, nil
		//	},
		//	ResponseHeaderTimeout: DEFAULT_RESPONSE_TIME_OUT, // 响应时间,0表示不会超时
		//},
	}
	//c.Transport = &http.Transport{
	//	Proxy: func(request *http.Request) (url *url.URL, e error) {
	//		return url.Parse(  "http://202.105.233.40:80")
	//	},
	//}
	r := &http.Request{}
	t := &http.Transport{}
	p := &Option{}
	return &Request{
		c: c,
		r: r,
		t: t,
		p: p,
	}
}

// https://www.cnblogs.com/mafeng/p/7068837.html
