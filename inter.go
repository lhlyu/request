package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (*Request) getDuration(times int) time.Duration {
	return DURATION * time.Duration(times)
}

func (rq *Request) setUrl() {
	apiUrl := rq.u
	if strings.Index(rq.u, "//") < 1 {
		if rq.b != "" {
			apiUrl = rq.b + rq.u
		}
	}
	uUrl, err := url.Parse(apiUrl)
	if rq.errHandler(err) {
		return
	}
	uValues := uUrl.Query()
	for k, v := range rq.p.Param {
		uValues.Set(k, v)
	}
	uUrl.RawQuery = uValues.Encode()
	rq.r.URL = uUrl
	return
}

func (rq *Request) setHeader() {
	for k, v := range rq.p.Header {
		rq.r.Header.Set(k, v)
	}
	return
}

func (rq *Request) getHeader(key string) string {
	return rq.r.Header.Get(key)
}

func (rq *Request) setCookie() {
	for _, v := range rq.p.Cookie {
		rq.r.AddCookie(v)
	}
	return
}

func (rq *Request) setFormData() {
	if len(rq.p.Data) == 0 {
		return
	}
	uValue := url.Values{}
	for k, v := range rq.p.Data {
		uValue.Set(k, anyToString(v))
	}
	rq.r.Body = ioutil.NopCloser(strings.NewReader(uValue.Encode()))
	return
}

func (rq *Request) setTransport() {
	rq.c.Timeout = rq.getDuration(rq.p.TimeOut)
	if rq.p.Proxy != "" {
		rq.t.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(rq.p.Proxy)
		}
	}
	rq.c.Transport = rq.t
	return
}

func (rq *Request) setFiles() {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for k, v := range rq.p.FileData {
		f, err := os.Open(v)
		defer f.Close()
		if err != nil {
			log.Printf("open %s file is err : %v\n", v, err)
			return
		}
		formData, err := writer.CreateFormFile(k, f.Name())
		if err != nil {
			log.Printf("create form field %v\n", k)
			return
		}

		_, err = io.Copy(formData, f)
		if err != nil {
			log.Printf("copy %s file is err : %v\n", v, err)
			return
		}
	}

	for k, v := range rq.p.Data {
		_ = writer.WriteField(k, anyToString(v))
	}
	err := writer.Close()
	if err != nil {
		log.Printf("close writer is err : %v\n", err)
		return
	}
	rq.r.Header.Set("Content-Type", writer.FormDataContentType())
	rq.r.Body = ioutil.NopCloser(strings.NewReader(body.String()))
}

func (rq *Request) setBody() {
	if len(rq.p.Data) == 0 {
		return
	}
	rq.r.Body = ioutil.NopCloser(strings.NewReader(rq.ObjToJson(rq.p.Data)))
	return
}

func (rq *Request) ObjToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bts, err := json.Marshal(v)
	if rq.errHandler(err) {
		return ""
	}
	return string(bts)
}

func (rq *Request) handler(f func()) {
	if rq.d {
		f()
	}

}

func (rq *Request) errHandler(err error) bool {
	if rq.d {
		if err != nil {
			log.Println("ERROR | ", err)
			return true
		}
	}
	return false
}

// send request
func (rq *Request) send() IResponse {

	rq.r.Method = rq.m
	rq.c.Timeout = rq.getDuration(rq.p.TimeOut)

	rq.setUrl()
	rq.setTransport()
	rq.setHeader()
	rq.setCookie()

	if rq.getHeader(CONTENT_TYPE) == X_WWW_FORM_URLENCODED {
		rq.setFormData()
	} else if rq.getHeader(CONTENT_TYPE) == FORM_DATA {
		rq.setFiles()
	} else {
		rq.setBody()
	}

	resp, err := rq.c.Do(rq.r)
	return newResponse(resp, err)
}
