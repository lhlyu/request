package request

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	error_url_is_empty     = errors.New("url or baseUrl is empty")
	error_method_is_empty  = errors.New("request method is empty")
	error_client_is_nil    = errors.New("http client is nil")
	error_request_is_nil   = errors.New("http request is nil")
	error_transport_is_nil = errors.New("http transport is nil")
	error_not_json         = errors.New("param is not json string")
)

// 接口实现
type Request struct {
	C            *http.Client
	R            *http.Request
	T            *http.Transport
	Debug        bool                                        // 调试
	BaseUrl      string                                      // 基础url
	Url          string                                      // 路由 或 完整请求地址，支持多个
	Timeout      int                                         // 超时时间 秒
	Proxy        string                                      // 代理地址
	Method       string                                      // 请求方法
	Header       http.Header                                 // 请求头
	Cookies      []*http.Cookie                              // cookies
	QueryParam   url.Values                                  // url 参数
	Body         string                                      // body 参数
	Interceptors func(r *http.Request, c *http.Client) error // 请求拦截器
	ErrorHandler func(err error)                             // 错误处理
	Err          error
}

// ------------- 私有方法 ---------------------

// 错误处理
func (r *Request) errHandler(err error) bool {
	r.Err = err
	if r.ErrorHandler != nil {
		r.ErrorHandler(err)
		return true
	}
	return false
}

// 检查
func (r *Request) check() error {
	if r.Url == "" && r.BaseUrl == "" {
		return error_url_is_empty
	}
	if r.Method == "" {
		return error_method_is_empty
	}
	if r.C == nil {
		return error_client_is_nil
	}
	if r.R == nil {
		return error_request_is_nil
	}
	return nil
}

func (r *Request) debug() {
	if r.Debug {
		buf := bytes.Buffer{}
		buf.WriteString("[request]\n")
		buf.WriteString(fmt.Sprintf("url: %s\n", r.R.URL.String()))
		buf.WriteString(fmt.Sprintf("method: %s\n\n", r.R.Method))
		buf.WriteString("header:\n")
		for k, h := range r.Header {
			buf.WriteString("\t" + k + ": ")
			for _, v := range h {
				buf.WriteString(v + ";")
			}
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
		buf.WriteString("body:\n")
		buf.WriteString(fmt.Sprint("\t", r.Body))
		log.Println(buf.String())
	}
}

// 获取请求地址
// http://localhost:8080/api?a=1
// {"Scheme":"http","Opaque":"","User":null,"Host":"localhost:8080","Path":"/api","RawPath":"","ForceQuery":false,"RawQuery":"a=1","Fragment":""}
func (r *Request) getUrl() (*url.URL, error) {
	apiUrl := r.BaseUrl + r.Url
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	for k, v := range r.QueryParam {
		for _, w := range v {
			u.Query().Add(k, w)
		}
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u, nil
}

// ----------------------- 实现 ----------------------

func (r *Request) SetDebug(enable bool) IRequest {
	r.Debug = enable
	return r
}

func (r *Request) SetBaseUrl(s string) IRequest {
	s = strings.TrimSpace(s)
	r.BaseUrl = s
	return r
}

func (r *Request) SetBaseUrlf(format string, v ...interface{}) IRequest {
	r.BaseUrl = fmt.Sprintf(format, v...)
	return r
}

func (r *Request) SetUrl(s string) IRequest {
	s = strings.TrimSpace(s)
	r.Url = s
	return r
}

func (r *Request) SetUrlf(format string, v ...interface{}) IRequest {
	r.Url = fmt.Sprintf(format, v...)
	return r
}

func (r *Request) SetTimeout(s int) IRequest {
	r.Timeout = s
	return r
}

func (r *Request) SetProxy(s string) IRequest {
	r.Proxy = s
	return r
}

func (r *Request) SetMethod(method string) IRequest {
	r.Method = method
	return r
}

func (r *Request) SetHeader(v interface{}) IRequest {
	m, err := toMSI(v)
	if err != nil {
		r.errHandler(err)
		return r
	}
	for k, v := range m {
		r.Header.Set(k, anyToString(v))
	}
	return r
}

func (r *Request) AddHeader(key string, value interface{}) IRequest {
	r.Header.Add(key, anyToString(value))
	return r
}

func (r *Request) SetContentType(s string) IRequest {
	return r.AddHeader("content-type", s)
}

func (r *Request) AddCookie(v *http.Cookie) IRequest {
	r.Cookies = append(r.Cookies, v)
	return r
}

func (r *Request) AddSimpleCookie(key string, value interface{}) IRequest {
	r.Cookies = append(r.Cookies, &http.Cookie{
		Name:  key,
		Value: anyToString(value),
	})
	return r
}

func (r *Request) SetQueryParam(v interface{}) IRequest {
	m, err := toMSI(v)
	if err != nil {
		r.errHandler(err)
		return r
	}
	for k, v := range m {
		r.QueryParam.Set(k, anyToString(v))
	}
	return r
}

func (r *Request) AddQueryParam(key string, value interface{}) IRequest {
	r.QueryParam.Add(key, anyToString(value))
	return r
}

func (r *Request) SetBody(s string) IRequest {
	r.Body = s
	return r
}

//// json 字符串
//func (r *Request) SetJsonBody(s string) IRequest {
//	if !json.Valid([]byte(s)){
//		r.errHandler(error_not_json)
//		return r
//	}
//	r.Body = s
//	return r
//}
//
//// query url
//func (r *Request) SetQueryBody(query string) IRequest {
//	r.Body = url.PathEscape(query)
//	return r
//}
//
//// line param  to json
//func (r *Request) SetLineJsonBody(line string) IRequest {
//	rows := strings.Split(line,"\n")
//	reg := regexp.MustCompile("[:=]")
//	m := make(map[string]string)
//	for _, row := range rows {
//		index := strings.Index(row,"//")
//		row = row[:index]
//		row = strings.TrimSpace(row)
//		if row == ""{
//			continue
//		}
//		cols := reg.Split(row,2)
//		if len(cols) != 2 {
//			cols = append(cols,"")
//		}
//		m[cols[0]] = cols[1]
//	}
//	byts,err := json.Marshal(m)
//	if err != nil{
//		r.errHandler(err)
//		return r
//	}
//	r.Body = string(byts)
//	return r
//}
//
//// line param  to query
//func (r *Request) SetLineQueryBody(line string) IRequest {
//	rows := strings.Split(line,"\n")
//	reg := regexp.MustCompile("[:=]")
//	u := url.Values{}
//	for _, row := range rows {
//		index := strings.Index(row,"//")
//		row = row[:index]
//		row = strings.TrimSpace(row)
//		if row == ""{
//			continue
//		}
//		cols := reg.Split(row,2)
//		if len(cols) != 2 {
//			cols = append(cols,"")
//		}
//		u.Add(cols[0],cols[1])
//	}
//	r.Body = u.Encode()
//	return r
//}

func (r *Request) SetInterceptors(f func(r *http.Request, c *http.Client) error) IRequest {
	r.Interceptors = f
	return r
}

func (r *Request) SetErrorHandler(f func(err error)) IRequest {
	r.ErrorHandler = f
	return r
}

func (r *Request) SetClient(c *http.Client) IRequest {
	r.C = c
	return r
}

func (r *Request) SetRequest(re *http.Request) IRequest {
	r.R = re
	return r
}

func (r *Request) SetTransport(t *http.Transport) IRequest {
	r.T = t
	return r
}

// -------------- 请求 ------------------

func (r *Request) DoHttp() IResponse {
	res := &Response{}
	if err := r.check(); err != nil {
		r.errHandler(err)
		return nil
	}
	// 设置超时
	r.C.Timeout = time.Duration(r.Timeout) * time.Second
	// 设置代理
	if r.Proxy != "" {
		if r.T == nil {
			r.errHandler(error_transport_is_nil)
			return nil
		}
		r.T.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(r.Proxy)
		}
	}
	if r.T != nil {
		r.C.Transport = r.T
	}

	// 设置 方法
	r.R.Method = r.Method
	// 设置 url
	apiUrl, err := r.getUrl()
	if err != nil {
		r.errHandler(err)
		return nil
	}
	r.R.URL = apiUrl

	// 设置 header
	r.R.Header = r.Header
	// 设置 cookie
	for _, v := range r.Cookies {
		r.R.AddCookie(v)
	}
	// 设置 body
	r.R.Body = ioutil.NopCloser(strings.NewReader(r.Body))
	if r.Interceptors != nil {
		err = r.Interceptors(r.R, r.C)
		if err != nil {
			r.errHandler(err)
			return nil
		}
	}
	r.debug()
	time.Sleep(time.Second * 2)
	resp, err := r.C.Do(r.R)
	if err != nil {
		r.errHandler(err)
		return nil
	}

	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.errHandler(err)
		return nil
	}
	res.respBody = string(byts)
	return res
}

func (r *Request) DoGet() IResponse {
	r.Method = GET
	return r.DoHttp()
}

func (r *Request) DoPost() IResponse {
	r.Method = POST
	return r.DoHttp()
}

func (r *Request) Get(format string, v ...interface{}) IResponse {
	r.SetUrlf(format, v...)
	return r.DoGet()
}

func (r *Request) Post(apiUrl string, body string) IResponse {
	r.SetUrl(apiUrl)
	r.SetBody(body)
	return r.DoPost()
}

// ----------------------- 异步请求
func (r *Request) Async() IResponse {
	res := &Response{
		asyncBody: make(chan string, 0),
	}
	go func() {
		if err := r.check(); err != nil {
			r.errHandler(err)
			return
		}
		// 设置超时
		r.C.Timeout = time.Duration(r.Timeout) * time.Second
		// 设置代理
		if r.Proxy != "" {
			if r.T == nil {
				r.errHandler(error_transport_is_nil)
				return
			}
			r.T.Proxy = func(request *http.Request) (i *url.URL, e error) {
				return url.Parse(r.Proxy)
			}
		}
		if r.T != nil {
			r.C.Transport = r.T
		}

		// 设置 方法
		r.R.Method = r.Method
		// 设置 url
		apiUrl, err := r.getUrl()
		if err != nil {
			r.errHandler(err)
			return
		}
		r.R.URL = apiUrl

		time.Sleep(time.Second * 2)
		// 设置 header
		r.R.Header = r.Header
		// 设置 cookie
		for _, v := range r.Cookies {
			r.R.AddCookie(v)
		}
		// 设置 body
		r.R.Body = ioutil.NopCloser(strings.NewReader(r.Body))
		if r.Interceptors != nil {
			err = r.Interceptors(r.R, r.C)
			if err != nil {
				r.errHandler(err)
				return
			}
		}
		r.debug()
		resp, err := r.C.Do(r.R)
		if err != nil {
			r.errHandler(err)
			return
		}

		byts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			r.errHandler(err)
			return
		}
		res.asyncBody <- string(byts)
	}()
	return res
}
