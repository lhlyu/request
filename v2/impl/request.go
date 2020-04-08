package impl

import (
	"fmt"
	request "github.com/lhlyu/request/v2/v2"
	"net/http"
	"net/url"
	"time"
)

// 接口实现
type Request struct {
	C            *http.Client
	R            *http.Request
	T            *http.Transport
	BaseUrl      string                   // 基础url
	Url          string                   // 路由 或 完整请求地址，支持多个
	Asynch       bool                     // 是否异步
	Timeout      int                      // 超时时间 秒
	Proxy        string                   // 代理地址
	Method       string                   // 请求方法
	Header       http.Header              // 请求头
	Cookies      []*http.Cookie           // cookies
	QueryParam   url.Values               // url 参数
	BodyData     interface{}              // body 参数
	Interceptors func(r request.IRequest) // 请求拦截器
	ErrorHandler func(err error)          // 错误处理
}

func (r *Request) SetBaseUrl(s string) request.IRequest {
	r.BaseUrl = s
	return r
}

func (r *Request) SetBaseUrlf(format string, v ...interface{}) request.IRequest {
	r.BaseUrl = fmt.Sprintf(format, v...)
	return r
}

func (r *Request) SetUrl(s string) request.IRequest {
	r.Url = s
	return r
}

func (r *Request) SetUrlf(format string, v ...interface{}) request.IRequest {
	r.Url = fmt.Sprintf(format, v...)
	return r
}

func (r *Request) SetAsynch(enable bool) request.IRequest {
	r.Asynch = enable
	return r
}

func (r *Request) SetTimeout(s int) request.IRequest {
	r.Timeout = s
	return r
}

func (r *Request) SetProxy(s string) request.IRequest {
	r.Proxy = s
	return r
}

func (r *Request) SetMethod(method string) request.IRequest {
	r.Method = method
	return r
}

func (r *Request) SetHeader(v interface{}) request.IRequest {
	m, err := toMSI(v)
	if err != nil {
		r.ErrorHandler(err)
		return r
	}
	for k, v := range m {
		r.Header.Set(k, anyToString(v))
	}
	return r
}

func (r *Request) AddHeader(key string, value interface{}) request.IRequest {
	r.Header.Add(key, anyToString(value))
	return r
}

func (r *Request) SetContentType(s string) request.IRequest {
	return r.AddHeader("content-type", s)
}

func (r *Request) AddCookie(v *http.Cookie) request.IRequest {
	r.Cookies = append(r.Cookies, v)
	return r
}

func (r *Request) AddSimpleCookie(key string, value interface{}) request.IRequest {
	r.Cookies = append(r.Cookies, &http.Cookie{
		Name:  key,
		Value: anyToString(value),
	})
	return r
}

func (r *Request) SetQueryParam(v interface{}) request.IRequest {
	m, err := toMSI(v)
	if err != nil {
		r.ErrorHandler(err)
		return r
	}
	for k, v := range m {
		r.QueryParam.Set(k, anyToString(v))
	}
	return r
}

func (r *Request) AddQueryParam(key string, value interface{}) request.IRequest {
	r.QueryParam.Add(key, anyToString(value))
	return r
}

func (r *Request) SetBodyData(v interface{}) request.IRequest {
	r.BodyData = v
	return r
}

func (r *Request) SetInterceptors(f func(r request.IRequest)) request.IRequest {
	r.Interceptors = f
	return r
}

func (r *Request) SetErrorHandler(f func(err error)) request.IRequest {
	r.ErrorHandler = f
	return r
}

func (r *Request) SetClient(c *http.Client) request.IRequest {
	r.C = c
	return r
}

func (r *Request) SetRequest(re *http.Request) request.IRequest {
	r.R = re
	return r
}

func (r *Request) SetTransport(t *http.Transport) request.IRequest {
	r.T = t
	return r
}

// -------------- 请求 ------------------
func (r *Request) DoHttp() request.IResponse {
	// 设置超时
	r.C.Timeout = time.Duration(r.Timeout) * time.Second
	// 设置代理
	if r.Proxy != "" {
		r.T.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(r.Proxy)
		}
	}
	// 设置 方法
	r.R.Method = r.Method
	// 设置 url
	baseUrl, err := url.Parse(r.BaseUrl)
	if err != nil {
		r.ErrorHandler(err)
		return nil
	}
	u, err := url.Parse(r.Url)
	if err != nil {
		r.ErrorHandler(err)
		return nil
	}

	// 设置 header
	// 设置 cookie
	// 设置 body

	return nil
}

func (r *Request) DoGet() request.IResponse {
	return nil
}

func (r *Request) DoPost() request.IResponse {
	return nil
}

func (r *Request) Get(format string, v ...interface{}) request.IResponse {
	return nil
}

func (r *Request) Post(apiUrl string, body interface{}) request.IResponse {
	return nil
}
