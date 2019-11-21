package request

import "time"

const _version = "v1.4.0"

const (
	DURATION                  = time.Second
	DEFAULT_RESPONSE_TIME_OUT = DURATION * 5 // 默认超时时间 5s
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
	_KV_LINE = iota // 行
	_QS             // queryUrl
	_JSON           // json
)

// header key
const (
	CONTENT_TYPE = "Content-Type"
)

// header value
const (
	APPLICATION_JSON      = "application/json"                  //  json参数放body
	X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded" // 一般运用于表单提交 参数放body， & 拼接
	FORM_DATA             = "multipart/form-data"               // 文件上传  参数放body， 参数用分隔符 --------
)

//var Types = map[string]string{
//	"html":       "text/html",
//	"json":       "application/json",
//	"xml":        "application/xml",
//	"text":       "text/plain",
//	"urlencoded": "application/x-www-form-urlencoded",
//	"form":       "application/x-www-form-urlencoded",
//	"form-data":  "application/x-www-form-urlencoded",
//	"multipart":  "multipart/form-data",
//}
