## doc

#### request

```
NewRequest() IRequest  // 新建一个请求
```

#### option定义

```go
type Option struct {
	Proxy    string                 // 代理地址
	TimeOut  int                    // 超时时间 秒
	Cookie   []*Cookie              // cookies
	Param    map[string]string      // url参数， ? 后面的参数
	Header   map[string]string      // 请求头
	Data     map[string]interface{} // body 参数 json形式
	FileData map[string]string      // 文件
}
```

#### request 方法

```go 
SetDebug(debug bool) IRequest                           // 打印日志
SetBaseUrl(baseUrl string) IRequest                     // 这个会拼接在url前面
SetUrl(apiUrl string) IRequest                          // 设置Url
SetUrlf(urlFormat string, v ...interface{}) IRequest    // 格式化设置Url
SetMethod(method string) IRequest                       // 设置请求方法

SetTimeOut(second int) IRequest                         // 设置超时时间
SetProxy(proxy string) IRequest                         // 设置代理
SetParam(v interface{}) IRequest                        // 设置url参数,? 后面的参数
SetHeader(v interface{}) IRequest                       // 设置头
SetData(v interface{}) IRequest                         // 设置body参数
SetFile(v interface{}) IRequest                         // 设置文件

AddCookies(cookies ...*Cookie) IRequest                 // 添加Cookie
AddParam(k string, v interface{}) IRequest              // 添加url参数,? 后面的参数
AddHeader(k string, v interface{}) IRequest             // 添加头部参数
AddCookieSimple(k string, v interface{}) IRequest       // 添加简单的cookie
AddData(k string, v interface{}) IRequest               // 添加body 参数
AddFile(field, filePath string) IRequest                // 添加文件

SetContentType(tp string) IRequest                      // 设置 content-type

SetTransport(transport *http.Transport) IRequest        // 自定义transport
SetRequest(r *http.Request) IRequest                    // 自定义请求体
SetClient(c *http.Client) IRequest                      // 自定义HTTP client
SetOption(option *Option) IRequest                      // 自定义参数

Get(apiUrl string, param interface{}) IResponse         // GET 请求
Post(apiUrl string, data interface{}) IResponse         // POST 请求
PostForm(apiUrl string, formData interface{}) IResponse // POST FORM 请求
DoHttp() IResponse                                      // 发起请求
DoPost() IResponse                                      // 发起post请求
DoGet() IResponse                                       // 发起get请求
DoPostForm() IResponse                                  // 发起post form请求
DoPostFile() IResponse                                  // 发起post multipart/form-data 请求
```

#### response 方法

```go
GetRequest() *http.Request         // 获取请求体
GetResponse() *http.Response       // 获取响应体
GetBody() string                   // 获取body内容
GetStatus() int                    // 获取响应状态码
IsStatusOk() bool                  // 响应状态码 == 200
AssertStatus(status int)           // 断言状态码
BodyUnmarshal(v interface{}) error // 数据解析
Then(func(resp *http.Response))    // 自定义处理   response
Error() error                      // 返回请求错误
```