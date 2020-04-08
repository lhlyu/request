package v2

import (
	"fmt"
	"net/http"
	"net/url"
)

// 接口实现
type Request struct {
	c            *http.Client
	r            *http.Request
	t            *http.Transport
	BaseUrl      string           // 基础url
	Urls         map[string]bool  // 路由 或 完整请求地址，支持多个
	Asynch       bool             // 是否异步
	Timeout      int              // 超时时间 秒
	Proxy        string           // 代理地址
	Method       string           // 请求方法
	Header       *http.Header     // 请求头
	Cookies      []*http.Cookie   // cookies
	QueryParam   *url.Values      // url 参数
	BodyData     interface{}      // body 参数
	Interceptors func(r IRequest) // 请求拦截器
	ErrorHandler func(err error)  // 错误处理
}

func (r *Request) SetBaseUrl(s string) IRequest {
	r.BaseUrl = s
	return nil
}

func (r *Request) SetBaseUrlf(format string, v ...interface{}) IRequest {
	r.BaseUrl = fmt.Sprintf(format, v...)
	return nil
}

func (r *Request) SetUrl(s string) IRequest {
	r.Urls[s] = true
	return nil
}

func (r *Request) SetUrlf(format string, v ...interface{}) IRequest {
	r.Urls[fmt.Sprintf(format, v...)] = true
	return nil
}

func (r *Request) RemoveUrls(ss ...string) IRequest {
	for _, v := range ss {
		if _, ok := r.Urls[v]; ok {
			delete(r.Urls, v)
		}
	}
	return nil
}

func (r *Request) AddUrl(s string) IRequest {
	r.Urls[s] = true
	return nil
}

func (r *Request) AddUrlf(format string, v ...interface{}) IRequest {
	r.Urls[fmt.Sprintf(format, v...)] = true
	return nil
}

func (r *Request) SetUrls(ss ...string) IRequest {
	for _, v := range ss {
		r.Urls[v] = true
	}
	return nil
}

func (r *Request) SetAsynch(enable bool) IRequest {
	r.Asynch = enable
	return nil
}

func (r *Request) SetTimeout(s int) IRequest {
	r.Timeout = s
	return nil
}

func (r *Request) SetProxy(s string) IRequest {
	r.Proxy = s
	return nil
}

func (r *Request) SetMethod(method string) IRequest {
	r.Method = method
	return nil
}

func (r *Request) SetHeader(v interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) AddHeader(key string, value interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) SetContentType(s string) IRequest {
	return r.AddHeader("content-type", s)
}

func (r *Request) SetCookie(v interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) AddCookie(key string, value interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) SetQueryParam(v interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) AddQueryParam(key string, value interface{}) IRequest {
	// todo
	return nil
}

func (r *Request) SetBodyData(v interface{}) IRequest {
	r.BodyData = v
	return nil
}

func (r *Request) SetInterceptors(f func(r IRequest)) IRequest {
	r.Interceptors = f
	return nil
}

func (r *Request) SetErrorHandler(f func(err error)) IRequest {
	r.ErrorHandler = f
	return nil
}
