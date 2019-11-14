package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
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
	rq.handler(func() {
		log.Println("requestUrl | ",uUrl.String())
	})
	rq.r.URL = uUrl
	return
}


func (rq *Request) setHeader(){
	headers := http.Header{}
	for k, v := range rq.p.Header {
		headers.Set(k,v)
	}
	rq.r.Header = headers
	return
}

func (rq *Request) setCookie(){
	for _, v := range rq.p.Cookie {
		rq.r.AddCookie(v)
	}
	return
}

func (rq *Request) setPostData(){
	return
}

func (rq *Request) setTransport(){
	rq.c.Transport = rq.t
	return
}

func (rq *Request) setFiles(){

}

func (rq *Request) setBody(){
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
func (rq *Request) send(){

	rq.c.Timeout = rq.getDuration(rq.p.TimeOut)

	rq.setTransport()
	rq.setCookie()
	rq.setFiles()

	rq.r.Method = rq.m

	rq.setUrl()
	rq.setHeader()
	rq.setPostData()
	rq.setBody()

	resp,err := rq.c.Do(rq.r)
	if rq.errHandler(err){
		return
	}
	//httputil.DumpRequest()
	bts,_ := httputil.DumpResponse(resp,true)
	buf := &bytes.Buffer{}
	json.Indent(buf,bts,"  ","   ")
	fmt.Println("debug = ",string(bts))
}