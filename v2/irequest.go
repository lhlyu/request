package v2

import "net/http"

// 接口定义
type IRequest interface {
	// 设置基础url
	SetBaseUrl(s string) IRequest
	SetBaseUrlf(format string, v ...interface{}) IRequest

	// 设置单个路由
	SetUrl(s string) IRequest
	SetUrlf(format string, v ...interface{}) IRequest
	// 移除路由
	RemoveUrls(ss ...string) IRequest
	// 添加一个路由
	AddUrl(s string) IRequest
	AddUrlf(format string, v ...interface{}) IRequest
	// 设置多个路由
	SetUrls(ss ...string) IRequest

	// 设置是否异步
	SetAsynch(enable bool) IRequest

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

	// 设置Cookie,用法:
	SetCookie(v interface{}) IRequest
	// 添加cookie
	AddCookie(key string, value interface{}) IRequest

	// 设置请求参数,用法:
	SetQueryParam(v interface{}) IRequest
	// 添加请求参数
	AddQueryParam(key string, value interface{}) IRequest

	// 设置body参数,用法:
	SetBodyData(v interface{}) IRequest

	// 设置请求拦截器
	SetInterceptors(f func(r IRequest)) IRequest

	// 设置错误处理
	SetErrorHandler(f func(err error)) IRequest

	// -------------- 自定义 -----------------
	SetClient(c *http.Client) IRequest
	SetRequest(r *http.Request) IRequest
	SetTransport(t *http.Transport) IRequest
}

func NewRequest() IRequest {
	c := &http.Request{}
	return nil
}
