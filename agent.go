// The MIT License (MIT)
//
// Copyright (c) 2014 Theeraphol Wattanavekin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
	"moul.io/http2curl"
)

var DisableTransportSwap = false

func New() *SuperAgent {
	cookiejarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookiejarOptions)

	debug := os.Getenv("DEBUG") == "1"

	s := &SuperAgent{
		TargetType:        TypeJSON,
		Data:              make(map[string]interface{}),
		Header:            http.Header{},
		RawString:         "",
		SliceData:         []interface{}{},
		QueryData:         url.Values{},
		FileData:          make([]File, 0),
		BounceToRawString: false,
		Client:            &http.Client{Jar: jar},
		Transport:         &http.Transport{},
		Cookies:           make([]*http.Cookie, 0),
		Errors:            nil,
		BasicAuth:         struct{ Username, Password string }{},
		isDebug:           debug,
		CurlCommand:       false,
		logger:            log.New(os.Stderr, "[gorequest]", log.LstdFlags),
		isClone:           false,
		isAsynch:          false,
		errHandler: func(e error) {
			log.Println(e.Error())
		},
	}
	s.Transport.DisableKeepAlives = true
	return s
}

func (s *SuperAgent) Clone() *SuperAgent {
	clone := &SuperAgent{
		Url:                  s.Url,
		Method:               s.Method,
		Header:               http.Header(cloneMapArray(s.Header)),
		TargetType:           s.TargetType,
		ForceType:            s.ForceType,
		Data:                 shallowCopyData(s.Data),
		SliceData:            shallowCopyDataSlice(s.SliceData),
		QueryData:            url.Values(cloneMapArray(s.QueryData)),
		FileData:             shallowCopyFileArray(s.FileData),
		BounceToRawString:    s.BounceToRawString, // 是否需要反射到原字符串
		RawString:            s.RawString,
		Client:               s.Client,
		Transport:            s.Transport,
		Cookies:              shallowCopyCookies(s.Cookies),
		Errors:               shallowCopyErrors(s.Errors),
		BasicAuth:            s.BasicAuth,
		isDebug:              s.isDebug,
		CurlCommand:          s.CurlCommand,
		logger:               s.logger,
		Retryable:            copyRetryable(s.Retryable),
		DoNotClearSuperAgent: true,
		isClone:              true,
		isAsynch:             s.isAsynch,
		errHandler:           s.errHandler,
	}
	return clone
}

func (s *SuperAgent) Timeout(timeout time.Duration) *SuperAgent {
	s.safeModifyHttpClient()
	s.Client.Timeout = timeout
	return s
}

func (s *SuperAgent) SetErrHandler(fn func(e error)) *SuperAgent {
	s.errHandler = fn
	return s
}

func (s *SuperAgent) Debug() *SuperAgent {
	s.isDebug = true
	return s
}

func (s *SuperAgent) SetCurlCommand(enable bool) *SuperAgent {
	s.CurlCommand = enable
	return s
}

func (s *SuperAgent) Asynch() *SuperAgent {
	s.isAsynch = true
	return s
}

// 每次请求前都清除上一次的数据
func (s *SuperAgent) SetDoNotClearSuperAgent(enable bool) *SuperAgent {
	s.DoNotClearSuperAgent = enable
	return s
}

func (s *SuperAgent) SetLogger(logger Logger) *SuperAgent {
	s.logger = logger
	return s
}

func (s *SuperAgent) SetHeader(param string, value string) *SuperAgent {
	s.Header.Set(param, value)
	return s
}

func (s *SuperAgent) AppendHeader(param string, value string) *SuperAgent {
	s.Header.Add(param, value)
	return s
}

// 设置重试 重试次数，重试等待时间 期望状态码
func (s *SuperAgent) Retry(retryerCount int, retryerTime time.Duration, statusCode ...int) *SuperAgent {
	for _, code := range statusCode {
		statusText := http.StatusText(code)
		if len(statusText) == 0 {
			s.Errors = append(s.Errors, errors.New("StatusCode '"+strconv.Itoa(code)+"' doesn't exist in http package"))
		}
	}

	s.Retryable = struct {
		RetryableStatus []int
		RetryerTime     time.Duration
		RetryerCount    int
		Attempt         int
		Enable          bool
	}{
		statusCode,
		retryerTime,
		retryerCount,
		0,
		true,
	}
	return s
}

func (s *SuperAgent) SetBasicAuth(username string, password string) *SuperAgent {
	s.BasicAuth = struct{ Username, Password string }{username, password}
	return s
}

func (s *SuperAgent) AddCookie(c *http.Cookie) *SuperAgent {
	s.Cookies = append(s.Cookies, c)
	return s
}

func (s *SuperAgent) AddSimpleCookie(name, val string) *SuperAgent {
	s.Cookies = append(s.Cookies, &http.Cookie{
		Name:  name,
		Value: val,
	})
	return s
}

func (s *SuperAgent) AddCookies(cookies []*http.Cookie) *SuperAgent {
	s.Cookies = append(s.Cookies, cookies...)
	return s
}

func (s *SuperAgent) Type(typeStr string) *SuperAgent {
	if _, ok := Types[typeStr]; ok {
		s.ForceType = typeStr
	} else {
		s.Errors = append(s.Errors, errors.New("Type func: incorrect type \""+typeStr+"\""))
	}
	return s
}

func (s *SuperAgent) Query(content interface{}) *SuperAgent {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.queryString(v.String())
	case reflect.Struct:
		s.queryStruct(v.Interface())
	case reflect.Map:
		s.queryMap(v.Interface())
	default:
	}
	return s
}

func (s *SuperAgent) Param(key string, value string) *SuperAgent {
	s.QueryData.Add(key, value)
	return s
}

// 传输层安全协议配置 http2
func (s *SuperAgent) TLSClientConfig(config *tls.Config) *SuperAgent {
	s.safeModifyTransport()
	s.Transport.TLSClientConfig = config
	return s
}

func (s *SuperAgent) Proxy(proxyUrl string) *SuperAgent {
	parsedProxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		s.Errors = append(s.Errors, err)
	} else if proxyUrl == "" {
		s.safeModifyTransport()
		s.Transport.Proxy = nil
	} else {
		s.safeModifyTransport()
		s.Transport.Proxy = http.ProxyURL(parsedProxyUrl)
	}
	return s
}

// 重定向策略设置
func (s *SuperAgent) RedirectPolicy(policy func(req *http.Request, via []*http.Request) error) *SuperAgent {
	s.safeModifyHttpClient()
	s.Client.CheckRedirect = func(r *http.Request, v []*http.Request) error {
		vv := make([]*http.Request, len(v))
		for i, r := range v {
			vv[i] = r
		}
		return policy(r, vv)
	}
	return s
}

func (s *SuperAgent) Send(content interface{}) *SuperAgent {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.SendString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.SendString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.SendString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float64:
		s.SendString(strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.Float32:
		s.SendString(strconv.FormatFloat(v.Float(), 'f', -1, 32))
	case reflect.Bool:
		s.SendString(strconv.FormatBool(v.Bool()))
	case reflect.Struct:
		s.SendStruct(v.Interface())
	case reflect.Slice:
		s.SendSlice(makeSliceOfReflectValue(v))
	case reflect.Array:
		s.SendSlice(makeSliceOfReflectValue(v))
	case reflect.Ptr:
		s.Send(v.Elem().Interface())
	case reflect.Map:
		s.SendMap(v.Interface())
	default:
		return s
	}
	return s
}

func (s *SuperAgent) SendSlice(content []interface{}) *SuperAgent {
	s.SliceData = append(s.SliceData, content...)
	return s
}

func (s *SuperAgent) SendMap(content interface{}) *SuperAgent {
	return s.SendStruct(content)
}

func (s *SuperAgent) SendStruct(content interface{}) *SuperAgent {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				s.Data[k] = v
			}
		}
	}
	return s
}

func (s *SuperAgent) SendString(content string) *SuperAgent {
	if !s.BounceToRawString {
		var val interface{}
		d := json.NewDecoder(strings.NewReader(content))
		d.UseNumber()
		if err := d.Decode(&val); err == nil {
			switch v := reflect.ValueOf(val); v.Kind() {
			case reflect.Map:
				for k, v := range val.(map[string]interface{}) {
					s.Data[k] = v
				}
			case reflect.Slice:
				s.SendSlice(val.([]interface{}))
			default:
				s.BounceToRawString = true
			}
		} else if formData, err := url.ParseQuery(content); err == nil {
			for k, formValues := range formData {
				for _, formValue := range formValues {
					if val, ok := s.Data[k]; ok {
						var strArray []string
						strArray = append(strArray, string(formValue))
						switch oldValue := val.(type) {
						case []string:
							strArray = append(strArray, oldValue...)
						case string:
							strArray = append(strArray, oldValue)
						}
						s.Data[k] = strArray
					} else {
						s.Data[k] = formValue
					}
				}
			}
			s.TargetType = TypeForm
		} else {
			s.BounceToRawString = true
		}
	}
	s.RawString += content
	return s
}

func (s *SuperAgent) SendFile(file interface{}, args ...string) *SuperAgent {

	filename := ""
	fieldname := "file"
	argsLen := len(args)

	if argsLen >= 1 && len(args[0]) > 0 {
		filename = strings.TrimSpace(args[0])
	}
	if argsLen >= 2 && len(args[1]) > 0 {
		fieldname = strings.TrimSpace(args[1])
	}
	if fieldname == "file" || fieldname == "" {
		fieldname = "file" + strconv.Itoa(len(s.FileData)+1)
	}

	switch v := reflect.ValueOf(file); v.Kind() {
	case reflect.String:
		pathToFile, err := filepath.Abs(v.String())
		if err != nil {
			s.Errors = append(s.Errors, err)
			return s
		}
		if filename == "" {
			filename = filepath.Base(pathToFile)
		}
		data, err := ioutil.ReadFile(v.String())
		if err != nil {
			s.Errors = append(s.Errors, err)
			return s
		}
		s.FileData = append(s.FileData, File{
			Filename:  filename,
			Fieldname: fieldname,
			Data:      data,
		})
	case reflect.Slice:
		slice := makeSliceOfReflectValue(v)
		if filename == "" {
			filename = "filename"
		}
		f := File{
			Filename:  filename,
			Fieldname: fieldname,
			Data:      make([]byte, len(slice)),
		}
		for i := range slice {
			f.Data[i] = slice[i].(byte)
		}
		s.FileData = append(s.FileData, f)
	case reflect.Ptr:
		if argsLen == 1 {
			return s.SendFile(v.Elem().Interface(), args[0])
		}
		if argsLen >= 2 {
			return s.SendFile(v.Elem().Interface(), args[0], args[1])
		}
		return s.SendFile(v.Elem().Interface())
	default:
		if v.Type() == reflect.TypeOf(os.File{}) {
			osfile := v.Interface().(os.File)
			if filename == "" {
				filename = filepath.Base(osfile.Name())
			}
			data, err := ioutil.ReadFile(osfile.Name())
			if err != nil {
				s.Errors = append(s.Errors, err)
				return s
			}
			s.FileData = append(s.FileData, File{
				Filename:  filename,
				Fieldname: fieldname,
				Data:      data,
			})
			return s
		}

		s.Errors = append(s.Errors, errors.New("SendFile currently only supports either a string (path/to/file), a slice of bytes (file content itself), or a os.File!"))
	}

	return s
}

// ------------------------------- 请求 -------------------------------------------

func (s *SuperAgent) Http(method, targetUrl string, a ...interface{}) *SuperReceiver {
	targetUrl = fmt.Sprintf(targetUrl, a...)
	s.clearSuperAgent()
	s.Method = method
	s.Url = targetUrl
	s.Errors = nil
	return s.do()
}

func (s *SuperAgent) Get(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(GET, targetUrl, a...)
}

func (s *SuperAgent) Post(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(POST, targetUrl, a...)
}

func (s *SuperAgent) Head(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(HEAD, targetUrl, a...)
}

func (s *SuperAgent) Put(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(PUT, targetUrl, a...)
}

func (s *SuperAgent) Delete(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(DELETE, targetUrl, a...)
}

func (s *SuperAgent) Patch(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(PATCH, targetUrl, a...)
}

func (s *SuperAgent) Options(targetUrl string, a ...interface{}) *SuperReceiver {
	return s.Http(OPTIONS, targetUrl, a...)
}

// 统一请求
func (s *SuperAgent) do() *SuperReceiver {
	res := newSuperReceiver(s.isAsynch)
	var (
		errs []error
		resp *http.Response
	)
	if s.isAsynch {
		go func() {
			for {
				resp, errs = s.getResponseBytes()
				if errs != nil {
					break
				}
				// 是否可以重试请求，设置重试的次数
				if s.isRetryableRequest(resp) {
					resp.Header.Set("Retry-Count", strconv.Itoa(s.Retryable.Attempt))
					break
				}
			}
			res.resp <- resp
			res.Errs = append(res.Errs, errs...)
		}()
		return res
	}

	for {
		resp, errs = s.getResponseBytes()
		if errs != nil {
			break
		}
		// 是否可以重试请求，设置重试的次数
		if s.isRetryableRequest(resp) {
			resp.Header.Set("Retry-Count", strconv.Itoa(s.Retryable.Attempt))
			break
		}
	}
	res.Resp = resp
	res.Errs = append(res.Errs, errs...)
	return res
}

func (s *SuperAgent) AsCurlCommand() (string, error) {
	req, err := s.makeRequest()
	if err != nil {
		return "", err
	}
	cmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}
	return cmd.String(), nil
}
