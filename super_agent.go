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
	"fmt"
	"io/ioutil"
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

// ------------------ 定义SuperAgent 实现 Agent 接口 ---------------------
type SuperAgent struct {
	// 请求url
	Url string
	// 请求方法
	Method string
	// 请求头
	Header     http.Header
	TargetType string
	ForceType  string
	Data       map[string]interface{}
	// 优先级 大于 Data
	SliceData []interface{}
	// url参数
	QueryData url.Values
	// 文件
	FileData []SuperFile
	// 是否将send的内容反射到 queryUrl
	BounceToRawString bool

	RawString string
	Client    *http.Client
	Transport *http.Transport
	Cookies   []*http.Cookie
	Errors    []error
	BasicAuth struct{ Username, Password string }
	// 是否打印curl命令
	isCurl bool
	// 重试设置
	Retryable superAgentRetryable

	DoNotClearSuperAgent bool
	// 是否开启debug
	isDebug bool
	// 标识，是否是克隆
	isClone bool
	// 是否异步
	isAsynch bool
}

// 请求重试设置
type superAgentRetryable struct {
	RetryableStatus []int
	RetryerTime     time.Duration
	RetryerCount    int
	Attempt         int
	Enable          bool
}

// 上传文件
type SuperFile struct {
	Filename  string
	Fieldname string
	Data      []byte
}

// 新建
func New() Agent {
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
		FileData:          make([]SuperFile, 0),
		BounceToRawString: false,
		Client:            &http.Client{Jar: jar},
		Transport:         &http.Transport{},
		Cookies:           make([]*http.Cookie, 0),
		Errors:            nil,
		BasicAuth:         struct{ Username, Password string }{},
		isDebug:           debug,
		isCurl:            false,
		isClone:           false,
		isAsynch:          false,
	}
	s.Transport.DisableKeepAlives = true
	return s
}

// --------------- 实现Agent接口 ------------------

func (s *SuperAgent) Debug() Agent {
	s.isDebug = true
	return s
}

func (s *SuperAgent) Clone() Agent {
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
		isCurl:               s.isCurl,
		Retryable:            copyRetryable(s.Retryable),
		DoNotClearSuperAgent: true,
		isClone:              true,
		isAsynch:             s.isAsynch,
	}
	return clone
}

func (s *SuperAgent) Timeout(timeout time.Duration) Agent {
	s.safeModifyHttpClient()
	s.Client.Timeout = timeout
	return s
}

func (s *SuperAgent) Curl() Agent {
	s.isCurl = true
	return s
}

func (s *SuperAgent) Asynch() Agent {
	s.isAsynch = true
	return s
}

// 每次请求前都清除上一次的数据
func (s *SuperAgent) SetDoNotClearSuperAgent(enable bool) Agent {
	s.DoNotClearSuperAgent = enable
	return s
}

func (s *SuperAgent) SetHeader(param string, value string) Agent {
	s.Header.Set(param, value)
	return s
}

func (s *SuperAgent) AddHeader(param string, value string) Agent {
	s.Header.Add(param, value)
	return s
}

// 设置重试 重试次数，重试等待时间 期望状态码
func (s *SuperAgent) Retry(retryerCount int, retryerTime time.Duration, statusCode ...int) Agent {
	for _, code := range statusCode {
		statusText := http.StatusText(code)
		if len(statusText) == 0 {
			s.Errors = append(s.Errors, e("Retry func: ", error_status_not_exist, code))
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

func (s *SuperAgent) SetBasicAuth(username string, password string) Agent {
	s.BasicAuth = struct{ Username, Password string }{username, password}
	return s
}

func (s *SuperAgent) AddCookie(c *http.Cookie) Agent {
	s.Cookies = append(s.Cookies, c)
	return s
}

func (s *SuperAgent) AddSimpleCookie(name, val string) Agent {
	s.Cookies = append(s.Cookies, &http.Cookie{
		Name:  name,
		Value: val,
	})
	return s
}

func (s *SuperAgent) AddCookies(cookies []*http.Cookie) Agent {
	s.Cookies = append(s.Cookies, cookies...)
	return s
}

func (s *SuperAgent) Type(typeStr string) Agent {
	if _, ok := Types[typeStr]; ok {
		s.ForceType = typeStr
	} else {
		s.Errors = append(s.Errors, e("Type func: ", error_type_not_support))
	}
	return s
}

func (s *SuperAgent) Query(content interface{}) Agent {
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

func (s *SuperAgent) AddParam(key string, value string) Agent {
	s.QueryData.Add(key, value)
	return s
}

// 传输层安全协议配置 http2
func (s *SuperAgent) TLSClientConfig(config *tls.Config) Agent {
	s.safeModifyTransport()
	s.Transport.TLSClientConfig = config
	return s
}

func (s *SuperAgent) Proxy(proxyUrl string) Agent {
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
func (s *SuperAgent) RedirectPolicy(policy func(req *http.Request, via []*http.Request) error) Agent {
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

func (s *SuperAgent) Send(content interface{}) Agent {
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

func (s *SuperAgent) SendSlice(content []interface{}) Agent {
	s.SliceData = append(s.SliceData, content...)
	return s
}

func (s *SuperAgent) SendMap(content interface{}) Agent {
	return s.SendStruct(content)
}

func (s *SuperAgent) SendStruct(content interface{}) Agent {
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

func (s *SuperAgent) SendString(content string) Agent {
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

func (s *SuperAgent) SendFile(file interface{}, args ...string) Agent {

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
		s.FileData = append(s.FileData, SuperFile{
			Filename:  filename,
			Fieldname: fieldname,
			Data:      data,
		})
	case reflect.Slice:
		slice := makeSliceOfReflectValue(v)
		if filename == "" {
			filename = "filename"
		}
		f := SuperFile{
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
			s.FileData = append(s.FileData, SuperFile{
				Filename:  filename,
				Fieldname: fieldname,
				Data:      data,
			})
			return s
		}

		s.Errors = append(s.Errors, e("SendFile func: ", error_file_type_not_support))
	}

	return s
}

// ------------------------------- 请求 -------------------------------------------

func (s *SuperAgent) Http(method, targetUrl string, a ...interface{}) Receiver {
	targetUrl = fmt.Sprintf(targetUrl, a...)
	s.clearSuperAgent()
	s.Method = method
	s.Url = targetUrl
	s.Errors = nil
	return s.do()
}

func (s *SuperAgent) Get(targetUrl string, a ...interface{}) Receiver {
	return s.Http(GET, targetUrl, a...)
}

func (s *SuperAgent) Post(targetUrl string, a ...interface{}) Receiver {
	return s.Http(POST, targetUrl, a...)
}

func (s *SuperAgent) Head(targetUrl string, a ...interface{}) Receiver {
	return s.Http(HEAD, targetUrl, a...)
}

func (s *SuperAgent) Put(targetUrl string, a ...interface{}) Receiver {
	return s.Http(PUT, targetUrl, a...)
}

func (s *SuperAgent) Delete(targetUrl string, a ...interface{}) Receiver {
	return s.Http(DELETE, targetUrl, a...)
}

func (s *SuperAgent) Patch(targetUrl string, a ...interface{}) Receiver {
	return s.Http(PATCH, targetUrl, a...)
}

func (s *SuperAgent) Options(targetUrl string, a ...interface{}) Receiver {
	return s.Http(OPTIONS, targetUrl, a...)
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
