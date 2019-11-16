## doc

#### request

```
NewRequest() IRequest  // 新建一个请求
```

#### option定义

```go
type Option struct {
	Proxy    string                     // 代理地址
	TimeOut  int                        // 超时时间 秒 
	Files    []string                   // 文件路径 (暂未实现)
	Cookie   []*Cookie                  // cookies  
	Param    map[string]string          // url参数， ? 后面的参数
	Header   map[string]string          // 请求头
	Data     map[string]interface{}     // body 参数 json形式
	FormData map[string]string          // body 参数 &拼接形式
}
```

#### request 方法

```go 
SetDebug(debug bool) *Request                             // 打印日志
SetBaseUrl(baseUrl string) *Request                       // 这个会拼接在url前面
SetUrl(apiUrl string) *Request                            // 设置Url
SetUrlf(urlFormat string, v ...interface{}) *Request      // 格式化设置Url
SetMethod(method string) *Request                         // 设置请求方法

SetTimeOut(second int) *Request                           // 设置超时时间
SetProxy(proxy string) *Request                           // 设置代理
SetParam(v interface{}) *Request                          // 设置url参数,? 后面的参数
SetHeader(v interface{}) *Request                         // 设置头
SetData(v interface{}) *Request                           // 设置json形式 的 body参数
SetFormData(v interface{}) *Request                       // 设置 url参数形式的 body参数
SetFiles(filePaths ...string) *Request                    // 设置文件 (未实现)

AddCookies(cookies ...*Cookie) *Request                   // 添加Cookie
AddParam(k string, v interface{}) *Request                // 添加url参数,? 后面的参数
AddHeader(k string, v interface{}) *Request               // 添加头部参数
AddCookieSimple(k string, v interface{}) *Request         // 添加简单的cookie
AddData(k string, v interface{}) *Request                 // 添加json形式的body 参数
AddFormData(k string, v interface{}) *Request             // 添加 url参数形式的 body参数
AddFile(filePath string) *Request                         // 添加文件 (未实现)

SetTransport(transport *http.Transport) *Request          // 自定义transport
SetRequest(r *http.Request) *Request                      // 自定义请求体
SetClient(c *http.Client) *Request                        // 自定义HTTP client
SetOption(option *Option) *Request                        // 自定义参数

DoHttp() IResponse                                        // 发起请求
DoGet(apiUrl string, param interface{}) IResponse         // GET 请求
DoPost(apiUrl string, data interface{}) IResponse         // POST 请求
DoPostForm(apiUrl string, formData interface{}) IResponse // POST FORM 请求
```

#### response 方法

```go
GetRequest() *http.Request          // 获取请求体
GetResponse() *http.Response        // 获取响应体
GetBody() string                    // 获取body内容
GetStatus() int                     // 获取响应状态码
BodyUnmarshal(v interface{}) error  // 数据解析
```