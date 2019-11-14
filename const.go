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