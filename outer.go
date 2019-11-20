package request

import (
	"fmt"
	"net/http"
)

func (rq *Request) SetDebug(debug bool) IRequest {
	rq.d = debug
	return rq
}

func (rq *Request) SetBaseUrl(baseUrl string) IRequest {
	rq.b = baseUrl
	return rq
}

func (rq *Request) SetUrl(apiUrl string) IRequest {
	rq.u = apiUrl
	return rq
}

func (rq *Request) SetUrlf(urlFormat string, v ...interface{}) IRequest {
	rq.u = fmt.Sprintf(urlFormat, v...)
	return rq
}

func (rq *Request) SetMethod(method string) IRequest {
	rq.m = method
	return rq
}

// 设置超时时间,超时包括连接时间，任何重定向和读取响应正文,0表示不会超时
func (rq *Request) SetTimeOut(second int) IRequest {
	rq.p.TimeOut = second
	return rq
}

func (rq *Request) SetProxy(proxy string) IRequest {
	rq.p.Proxy = proxy
	return rq
}

// map | string(query | lines | json)
func (rq *Request) SetParam(v interface{}) IRequest {
	rq.p.Param = toMSS(v)
	return rq
}

func (rq *Request) AddParam(k string, v interface{}) IRequest {
	rq.p.Param[k] = anyToString(v)
	return rq
}

// map | string(query | lines | json)
func (rq *Request) SetHeader(v interface{}) IRequest {
	rq.p.Header = toMSS(v)
	return rq
}

func (rq *Request) AddHeader(k string, v interface{}) IRequest {
	rq.p.Header[k] = anyToString(v)
	return rq
}

func (rq *Request) AddCookies(cookies ...*Cookie) IRequest {
	if cookies != nil {
		rq.p.Cookie = append(rq.p.Cookie, cookies...)
	}
	return rq
}

func (rq *Request) AddCookieSimple(k string, v interface{}) IRequest {
	rq.p.Cookie = append(rq.p.Cookie, &Cookie{
		Name:  k,
		Value: anyToString(v),
	})
	return rq
}

func (rq *Request) SetData(v interface{}) IRequest {
	rq.p.Data = toMSI(v)
	return rq
}

func (rq *Request) AddData(k string, v interface{}) IRequest {
	rq.p.Data[k] = v
	return rq
}

func (rq *Request) SetFile(v interface{}) IRequest {
	m := toMSS(v)
	for k, v := range m {
		rq.p.FileData[k] = v
	}
	return rq
}

func (rq *Request) AddFile(field, filePath string) IRequest {
	rq.p.FileData[field] = filePath
	return rq
}

// 设置 content-type
func (rq *Request) SetContentType(tp string) IRequest {
	rq.AddHeader(CONTENT_TYPE, tp)
	return rq
}

// 设置 Transport
func (rq *Request) SetTransport(transport *http.Transport) IRequest {
	rq.t = transport
	return rq
}

func (rq *Request) SetRequest(r *http.Request) IRequest {
	rq.r = r
	return rq
}

func (rq *Request) SetClient(c *http.Client) IRequest {
	rq.c = c
	return rq
}

func (rq *Request) SetOption(option *Option) IRequest {
	rq.p = option
	return rq
}

func (rq *Request) DoHttp() IResponse {
	return rq.send()
}

func (rq *Request) Get(apiUrl string, param interface{}) IResponse {
	rq.SetMethod(GET)
	rq.AddHeader(CONTENT_TYPE, APPLICATION_JSON)
	rq.SetUrl(apiUrl)
	rq.SetParam(param)
	return rq.send()
}

func (rq *Request) Post(apiUrl string, data interface{}) IResponse {
	rq.SetMethod(POST)
	rq.AddHeader(CONTENT_TYPE, APPLICATION_JSON)
	rq.SetUrl(apiUrl)
	rq.SetData(data)
	return rq.send()
}

func (rq *Request) PostForm(apiUrl string, formData interface{}) IResponse {
	rq.SetMethod(POST)
	rq.AddHeader(CONTENT_TYPE, X_WWW_FORM_URLENCODED)
	rq.SetUrl(apiUrl)
	rq.SetData(formData)
	return rq.send()
}

func (rq *Request) DoGet() IResponse {
	rq.SetMethod(GET)
	rq.AddHeader(CONTENT_TYPE, APPLICATION_JSON)
	return rq.send()
}

func (rq *Request) DoPost() IResponse {
	rq.SetMethod(POST)
	rq.AddHeader(CONTENT_TYPE, APPLICATION_JSON)
	return rq.send()
}

func (rq *Request) DoPostForm() IResponse {
	rq.SetMethod(POST)
	rq.AddHeader(CONTENT_TYPE, X_WWW_FORM_URLENCODED)
	return rq.send()
}

func (rq *Request) DoPostFile() IResponse {
	rq.SetMethod(POST)
	rq.AddHeader(CONTENT_TYPE, FORM_DATA)
	return rq.send()
}
