# request
go http request

### 使用

> go get -v github.com/lhlyu/request

### doc

[链接](./DOC.md)

### 示例

- GET 请求

```go
func main(){
    rq := request.NewRequest()
    resp := rq.Get("http://localhost:8080","a=1&b=2&c=3")
    fmt.Println(resp.GetStatus())
    fmt.Println(resp.GetBody())
    

    // 第二种方式 支持结构体
    
    type P struct{
        A  int  `json:"a"`
        B  int  `json:"b"`
        C  int  `json:"c"`
    }
    p := &P{1,2,3}
    rq.Get("http://localhost:8080",p)
    
    // 第三种方式
    
    rq.Get("http://localhost:8080",`
    a :  1    //这是一个注释会自动去除
    b:2
             c:3
    d // 这个参数无用会自动去除
    `)
    
    // 第四种 json字符串
    rq.Get("http://localhost:8080",`{"a":1,"b":"2"}`)

}
```

- POST 请求类型，参数会放到body里

```go
func main(){
    rq := request.NewRequest()
    // 同样支持四种形式的传参，参数都会放到 body
    rq.Post("http://localhost:8080","a=1&b=2&c=3")
    rq.Post("http://localhost:8080",`{"a":1,"b":"2"}`)
    rq.Post("http://localhost:8080",p)
    rq.Post("http://localhost:8080",`
    a :  1    //这是一个注释会自动去除
    b:2
           c:3
    d // 这个参数无用会自动去除
    `)
}
```

- POST FORM 同上

- 其他示例

```go
func main(){
    rq := request.NewRequest()
    rq.SetDebug(true)
    rq.SetTimeOut(15)
    rq.SetBaseUrl("http://localhost:8080")
    rq.SetUrl("/api/articles")
    rq.SetParam(`
        pageSize:10
        pageNum:1
    `)
    rq.SetMethod(request.GET)  // 默认是GET 
    rq.AddHeader("token","xxx")
    resp := rq.DoHttp()
    fmt.Println(resp.GetStatus())
    fmt.Println(resp.GetBody())
}
```

```go
func main() {
	rq := request.NewRequest()
	rq.SetUrl("http://xxxx.com/upload")
	rq.SetFile(`file:C:\xxx.jpg`)
	resp := rq.DoPostFile()
	fmt.Println(resp.GetBody())
	fmt.Println(resp.GetStatus())
}
```