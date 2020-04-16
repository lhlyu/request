package request

import (
	"net/http"
	"net/url"
	"time"
)

type superAgentRetryable struct {
	RetryableStatus []int
	RetryerTime     time.Duration
	RetryerCount    int
	Attempt         int
	Enable          bool
}

type SuperAgent struct {
	Url        string
	Method     string
	Header     http.Header
	TargetType string
	ForceType  string
	Data       map[string]interface{}
	// 优先级 大于 Data
	SliceData []interface{}
	// url参数
	QueryData url.Values
	// 文件
	FileData []File
	// 是否将send的内容反射到 queryUrl
	BounceToRawString bool
	RawString         string
	Client            *http.Client
	Transport         *http.Transport
	Cookies           []*http.Cookie
	Errors            []error
	BasicAuth         struct{ Username, Password string }
	isDebug           bool
	CurlCommand       bool
	logger            Logger
	// 重试设置
	Retryable            superAgentRetryable
	DoNotClearSuperAgent bool
	// 标识，是否是克隆
	isClone bool
	// 是否异步
	isAsynch bool
	// 错误处理
	errHandler func(e error)
}

type File struct {
	Filename  string
	Fieldname string
	Data      []byte
}
