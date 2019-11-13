package request

import "time"

const (
	DURATION                  = time.Second
	DEFAULT_RESPONSE_TIME_OUT = DURATION * 5 // 默认超时时间 5s
)

// error tip
const (
	url_empty = "url is empty"
)
