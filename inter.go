package request

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (*Request) getDuration(times int) time.Duration {
	return DURATION * time.Duration(times)
}

func (rq *Request) setUrl(){
	apiUrl := rq.u
	if strings.Index(rq.u,"//") < 1{
		if rq.b != ""{
			apiUrl = rq.b + rq.u
		}
	}
	uUrl,err := url.Parse(apiUrl)
	if rq.errHandler(err){
		return
	}
	uValues := uUrl.Query()
	for k,v := range rq.p.Param{
		uValues.Set(k,v)
	}
	uUrl.RawQuery = uValues.Encode()
	rq.r.URL = uUrl
	return
}

func (rq *Request) setHeader(){
	for k, v := range rq.p.Header {
		rq.r.Header.Set(k,v)
	}
	return
}

func (rq *Request) setCookie(){
	for _, v := range rq.p.Cookie {
		rq.r.AddCookie(v)
	}
	return
}

func (rq *Request) setFormData(){
	if len(rq.p.FormData) == 0{
		return
	}
	uValue := url.Values{}
	for k, v := range rq.p.FormData {
		uValue.Set(k, v)
	}
	rq.r.Body = ioutil.NopCloser(strings.NewReader(uValue.Encode()))
	return
}

func (rq *Request) setTransport(){
	rq.c.Timeout = rq.getDuration(rq.p.TimeOut)
	if rq.p.Proxy != "" {
		rq.t.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(rq.p.Proxy)
		}
	}
	rq.c.Transport = rq.t
	return
}

// todo
func (rq *Request) setFiles(){

}

func (rq *Request) setBody(){
	if len(rq.p.Data) == 0{
		return
	}
	rq.r.Body = ioutil.NopCloser(strings.NewReader(rq.ObjToJson(rq.p.Data)))
	return
}

func (rq *Request) ObjToJson(v interface{}) string{
	if v == nil{
		return ""
	}
	bts,err := json.Marshal(v)
	if rq.errHandler(err){
		return ""
	}
	return string(bts)
}

func (rq *Request) handler(f func()){
	if rq.d{
		f()
	}

}

func (rq *Request) errHandler(err error) bool{
	if rq.d{
		if err != nil{
			log.Println("ERROR | ",err)
			return true
		}
	}
	return false
}

// send request
func (rq *Request) send() IResponse{

	rq.c.Timeout = rq.getDuration(rq.p.TimeOut)

	rq.setTransport()
	rq.setCookie()
	rq.setFiles()

	rq.r.Method = rq.m

	rq.setUrl()
	rq.setHeader()
	rq.setFormData()
	rq.setBody()

	resp,err := rq.c.Do(rq.r)
	return newResponse(resp,err)
}