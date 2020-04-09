package request

import (
	"log"
	"net/http"
	"net/url"
)

// 请求接口定义
type IRequest interface {
	// 是否开启调试
	SetDebug(enable bool) IRequest
	// 设置基础url
	SetBaseUrl(s string) IRequest
	SetBaseUrlf(format string, v ...interface{}) IRequest

	// 设置单个路由
	SetUrl(s string) IRequest
	SetUrlf(format string, v ...interface{}) IRequest

	// 设置超时(秒)
	SetTimeout(s int) IRequest

	// 设置代理地址
	SetProxy(s string) IRequest

	// 设置请求方法
	SetMethod(method string) IRequest

	// 设置请求头,用法:
	SetHeader(v interface{}) IRequest
	// 添加请求头
	AddHeader(key string, value interface{}) IRequest
	// 设置内容类型
	SetContentType(s string) IRequest

	// 添加cookie
	AddCookie(v *http.Cookie) IRequest
	// 添加简单cookie
	AddSimpleCookie(key string, value interface{}) IRequest

	// 设置请求参数,用法:
	SetQueryParam(v interface{}) IRequest
	// 添加请求参数
	AddQueryParam(key string, value interface{}) IRequest

	// 设置body参数,用法:
	SetBody(s string) IRequest

	// 设置请求拦截器
	SetInterceptors(f func(r *http.Request, c *http.Client) error) IRequest

	// 设置错误处理
	SetErrorHandler(f func(err error)) IRequest

	// -------------- 自定义 -----------------
	SetClient(c *http.Client) IRequest
	SetRequest(r *http.Request) IRequest
	SetTransport(t *http.Transport) IRequest

	// -------------- 请求 ------------------
	DoHttp() IResponse
	DoGet() IResponse
	DoPost() IResponse
	Get(format string, v ...interface{}) IResponse
	Post(apiUrl string, body string) IResponse
	// --------------- 异步 -----------------
	Async() IResponse
}

func NewRequest() IRequest {
	re := &Request{
		C:          &http.Client{},
		R:          &http.Request{},
		T:          &http.Transport{},
		Method:     GET,
		Header:     make(http.Header),
		QueryParam: make(url.Values),
		ErrorHandler: func(err error) {
			log.Println(err)
		},
	}
	return re
}
