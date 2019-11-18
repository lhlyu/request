package request

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"
)

type IRequest interface {
	SetDebug(debug bool) *Request                        // 打印日志
	SetBaseUrl(baseUrl string) *Request                  // 这个会拼接在url前面
	SetUrl(apiUrl string) *Request                       // 设置Url
	SetUrlf(urlFormat string, v ...interface{}) *Request // 格式化设置Url
	SetMethod(method string) *Request                    // 设置请求方法

	SetTimeOut(second int) *Request   // 设置超时时间
	SetProxy(proxy string) *Request   // 设置代理
	SetParam(v interface{}) *Request  // 设置url参数,? 后面的参数
	SetHeader(v interface{}) *Request // 设置头
	SetData(v interface{}) *Request   // 设置body参数
	SetFile(v interface{}) *Request   // 设置文件

	AddCookies(cookies ...*Cookie) *Request           // 添加Cookie
	AddParam(k string, v interface{}) *Request        // 添加url参数,? 后面的参数
	AddHeader(k string, v interface{}) *Request       // 添加头部参数
	AddCookieSimple(k string, v interface{}) *Request // 添加简单的cookie
	AddData(k string, v interface{}) *Request         // 添加body 参数
	AddFile(field, filePath string) *Request          // 添加文件

	SetContentType(tp string) *Request // 设置 content-type

	SetTransport(transport *http.Transport) *Request // 自定义transport
	SetRequest(r *http.Request) *Request             // 自定义请求体
	SetClient(c *http.Client) *Request               // 自定义HTTP client
	SetOption(option *Option) *Request               // 自定义参数

	DoHttp() IResponse                                        // 发起请求
	DoGet(apiUrl string, param interface{}) IResponse         // GET 请求
	DoPost(apiUrl string, data interface{}) IResponse         // POST 请求
	DoPostForm(apiUrl string, formData interface{}) IResponse // POST FORM 请求
}

type Request struct {
	c *http.Client
	r *http.Request
	t *http.Transport
	p *Option
	b string
	u string
	m string
	d bool
}

type Option struct {
	Proxy    string                 // 代理地址
	TimeOut  int                    // 超时时间 秒
	Files    []string               // 文件路径 (暂未实现)
	Cookie   []*Cookie              // cookies
	Param    map[string]string      // url参数， ? 后面的参数
	Header   map[string]string      // 请求头
	Data     map[string]interface{} // body 参数 json形式
	FileData map[string]*os.File    // 文件
}

type (
	MSI    = map[string]interface{}
	MSS    = map[string]string
	Cookie = http.Cookie
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
			DisableCompression:    true,
			MaxConnsPerHost:       50,
		},
	}
	r := &http.Request{
		Header: make(http.Header),
	}
	t := &http.Transport{}
	p := &Option{
		Cookie: make([]*Cookie, 0),
		Files:  make([]string, 0),
		Param:  make(map[string]string),
		Data:   make(map[string]interface{}),
		Header: map[string]string{
			CONTENT_TYPE: APPLICATION_JSON,
		},
		FileData: make(map[string]*os.File),
	}
	return &Request{
		c: c,
		r: r,
		t: t,
		p: p,
		m: GET,
	}
}
