package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (rq *Request) setHeader(){
	return
}

func (rq *Request) setCookie(){
	return
}

func (rq *Request) setPostData(){
	return
}

func (rq *Request) setTransport(){
	return
}

func (rq *Request) setFiles(){

}

func (rq *Request) setBody(){
	rq.r.Body = ioutil.NopCloser(strings.NewReader(""))
	return
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