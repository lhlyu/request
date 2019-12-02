package main

import (
	"fmt"
	"github.com/lhlyu/request"
	"io/ioutil"
	"net/http"
)

func main() {
	//f1()
	//f2()
	//f3()
	//f4()
	//f5()
	//f6()
	//f7()
	//f8()
	//f9()
	//f10()
	//f11()
	//f12()
	//f13()
	//f14()
	f15()
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
	request.NewRequest().
		SetTimeOut(5). // 5 second
		SetBaseUrl("http://localhost:8080").
		SetUrl("/api/users"). // real url = baseUrl + Url
		DoGet().AssertStatus(200)
}

func f4() {
	request.NewRequest().
		SetProxy("http://localhost:8081").
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

func f9() {
	request.NewRequest().
		SetUrl("http://localhost:8080/api/articles").
		SetDebug(true).
		AddHeader("userId", 1).
		SetParam(`{"pageNum":1,"pageSize":10}`). // 这里的参数会拼接到 url ? 后面
		DoGet().Then(func(resp *http.Response) { // 自定义处理响应体
		if resp.StatusCode == 200 {
			fmt.Println("自定义处理响应体")
			bts, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(bts))
		}
	})
}

// 上传文件

func f10() {
	// eg1:
	request.NewRequest().
		SetUrl("http://localhost:8080/api/upload").
		SetFile(`file: E:\x.png`).
		DoPostFile().AssertStatus(200)

	// eg2:
	request.NewRequest().
		SetUrl("http://localhost:8080/api/upload").
		AddFile("file1", "E:\\x1.png"). // 文件1
		AddFile("file2", "E:\\x2.png"). // 文件2
		AddData("userId", 1).           // 其他参数
		DoPostFile().AssertStatus(200)
}

type Msg struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func f11() {

	result := &Msg{}

	m := map[string]interface{}{
		"userId": 1,
		"name":   "tt",
	}

	request.NewRequest().
		SetUrl("http://localhost:8080").
		SetData(m).
		DoPost().
		BodyUnmarshal(result)
}

func f12() {
	resp := request.NewRequest().
		SetUrl("http://localhost:8080").
		DoGet()
	err := resp.Error()
	fmt.Println(err)
	fmt.Println(resp.GetStatus())
}

func f13() {
	resp := request.NewRequest().SetHeader(`
Accept-Encoding:gzip, deflate, br  // 这种头会自动转成  identity
X-Token:frerf3f3f3f.dfrf3f4f.4f34f3r34r3435
Connection:keep-alive
Pragma:no-cache
Accept:*/*
`).SetUrl("http://localhost:8080").SetDebug(true).DoGet()
	fmt.Println(resp.GetStatus())
	resp.Then(func(resp *http.Response) {
		bts, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bts))
	}).Then(func(resp *http.Response) {
		fmt.Println(resp.Cookies())
	}).AssertStatus(0)
}

func f14() {
	request.NewRequest().
		Get("http://localhost:8080/sadadasd", nil).
		OnSuccess(func(resp request.IResponse) { // if 200 <= status < 300 then
			fmt.Println("success:", resp.GetStatus())
		}).OnError(func(resp request.IResponse) { // else then
		fmt.Println("error:", resp.GetStatus())
	})
}

func f15() {
	request.NewRequest().
		Get("https://www.163.com", nil).
		OnSuccess(func(resp request.IResponse) {
			imgs := resp.BodyCompile("(https?.*\\.jpg)")
			for k, v := range imgs {
				fmt.Println(k, v)
			}
		})
}
