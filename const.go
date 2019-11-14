package request

import "time"

const _version = "v1.0.0"

const (
	DURATION                  = time.Second
	DEFAULT_RESPONSE_TIME_OUT = DURATION * 5 // 默认超时时间 5s
)

// error tip
const (
	url_empty = "url is empty"
)

// http method
const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	HEAD   = "HEAD"
	TRACE  = "TRACE"
	OPTION = "OPTION"
)

// string type
const (
	_KV_LINE = iota
	_QS
	_JSON
)

// header key
const (
	CONTENT_TYPE = "Content-Type"
)

// header value
const (
	APPLICATION_JSON = "application/json"  //  json参数放body
	X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"   // 一般运用于表单提交 参数放body， & 拼接
	FORM_DATA = "multipart/form-data"   // 文件上传  参数放body， 参数用分隔符 --------
)