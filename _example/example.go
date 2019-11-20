package main

import "github.com/lhlyu/request"

func main() {
	f1()
}

func f1() {
	request.NewRequest().Get("http://www.baidu.com", nil).AssertStatus(200) // 打印状态码是否与200相等
}

func f2() {
	request.NewRequest().
		SetHeader(`token:fdhdufiwuf`).
		AddHeader("client", "android").
		AddCookieSimple("userId", "1").
		SetUrl("http://www.baidu.com").
		SetMethod(request.GET).DoHttp().AssertStatus(200)
}

func f3() {
	request.NewRequest().SetTimeOut(5). // 5 second
						SetBaseUrl("http://localhost:8080").
						SetUrl("/api/users"). // real url = baseUrl + Url
						DoGet().AssertStatus(200)
}

func f4() {
	request.NewRequest().SetProxy("http://localhost:8081").
		SetUrl("http://localhost:8080").
		DoGet().AssertStatus(200)
}

func f5() {
	request.NewRequest().
		SetUrl("http://localhost:8080").
		SetContentType(request.X_WWW_FORM_URLENCODED).
		SetData(`{"userId":1,"name":"tt"}`).
		DoHttp().AssertStatus(200)
}

func f6() {
	request.NewRequest().
		SetUrl("http://localhost:8080").
		SetData(`userId:1            // k-v 以行为一组
                    name:tt`).
		DoHttp().AssertStatus(200)
}

type A struct {
	UserId int    `json:"userId"`
	Name   string `json:"name"`
}

func f7() {
	request.NewRequest().
		SetUrl("http://localhost:8080").
		SetData(&A{1, "tt"}).
		DoHttp().AssertStatus(200)
}

func f8() {
	request.NewRequest().
		SetUrl("http://localhost:8080").
		SetData("userId=1&name=tt").
		DoPost().AssertStatus(200)
}
