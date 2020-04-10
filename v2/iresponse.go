package request

import "net/http"

type IResponse interface {
	GetRequest() *http.Request
	GetResponse() *http.Response

	// 获取返回响应内容
	GetBody() string
	// 获取响应状态码
	GetStatus() int
	// 校验状态码,默认判断200
	IsStatus(status ...int) bool
	// body内容解析
	BodyUnmarshal(v interface{}) error
	// 正则匹配body
	BodyCompile(pattern string) []string
	// 接收异步
	Await() IResponse
	// 错误
	Error() error
}
