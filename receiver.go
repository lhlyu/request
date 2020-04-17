package request

import "net/http"

// 响应接口定义
type Receiver interface {
	// 获取原生响应体
	GetResponse() *http.Response
	// 获取body内容
	GetBody() []byte
	// 获取body内容
	GetBodyString() string
	// body内容转对象
	BodyUnmarshal(v interface{}) error
	// 正则匹配body内容
	BodyCompile(pattern string) []string
	// 获取响应状态码，没有则返回 -1
	GetStatus() int
	// 异步等待响应结果返回
	Await() Receiver
	// 保存到文件
	Save(fl interface{}) error
	// 后续操作
	Then(fn func(r Receiver)) Receiver
	// 停止后续操作
	Stop()
	// 获取最后一个错误
	Error() error
	// 获取所有错误
	Errors() []error
}
