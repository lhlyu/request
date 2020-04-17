package request

import (
	"crypto/tls"
	"net/http"
	"time"
)

// 请求体接口定义
type Agent interface {
	// 开启debug
	Debug() Agent
	// 克隆
	Clone() Agent
	// 设置超时
	Timeout(timeout time.Duration) Agent
	// 是否打印curl命令
	Curl() Agent
	// 开启异步
	Asynch() Agent
	// 每次请求前都清空Agent
	SetDoNotClearSuperAgent(enable bool) Agent
	// 设置请求头
	SetHeader(param string, value string) Agent
	// 添加请求头
	AddHeader(param string, value string) Agent
	// 设置重试： 重试次数,间隔时间,期望返回的状态码列表
	Retry(retryerCount int, retryerTime time.Duration, statusCode ...int) Agent
	// 设置auth
	SetBasicAuth(username string, password string) Agent
	// 添加cookie
	AddCookie(c *http.Cookie) Agent
	// 添加简单cookie
	AddSimpleCookie(name, val string) Agent
	// 添加cookies
	AddCookies(cookies []*http.Cookie) Agent
	// 设置请求数据类型
	Type(typeStr string) Agent
	// 设置url参数
	Query(content interface{}) Agent
	// 添加url参数
	AddParam(key string, value string) Agent
	// 传输层安全协议配置 http2
	TLSClientConfig(config *tls.Config) Agent
	// 设置代理地址
	Proxy(proxyUrl string) Agent
	// 设置重定向策略
	RedirectPolicy(policy func(req *http.Request, via []*http.Request) error) Agent
	// 发送内容，一般放在body,自动判断类型
	Send(content interface{}) Agent
	// 指定发送切片
	SendSlice(content []interface{}) Agent
	// 指定发送map
	SendMap(content interface{}) Agent
	// 指定发送struct
	SendStruct(content interface{}) Agent
	// 指定发送string
	SendString(content string) Agent
	// 发送文件,允许输入文件路径，文件流, [文件名,Field]
	SendFile(file interface{}, args ...string) Agent

	// 通用请求方法
	Http(method, targetUrl string, a ...interface{}) Receiver
	// GET请求
	Get(targetUrl string, a ...interface{}) Receiver
	// POST请求
	Post(targetUrl string, a ...interface{}) Receiver
	// HEAD请求
	Head(targetUrl string, a ...interface{}) Receiver
	// PUT请求
	Put(targetUrl string, a ...interface{}) Receiver
	// DELETE请求
	Delete(targetUrl string, a ...interface{}) Receiver
	// PATCH请求
	Patch(targetUrl string, a ...interface{}) Receiver
	// OPTIONS请求
	Options(targetUrl string, a ...interface{}) Receiver

	// 返回curl
	AsCurlCommand() (string, error)
}
