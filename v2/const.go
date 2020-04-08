package request

import "time"

const _version = "v2.0.0"

const (
	DURATION                  = time.Second
	DEFAULT_RESPONSE_TIME_OUT = DURATION * 5 // 默认超时时间 5s
)

// http method
const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	HEAD    = "HEAD"
	TRACE   = "TRACE"
	OPTIONS = "OPTIONS"
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

// 异步
// 复用
// 日志
// body数据格式
// 多个请求
// 请求拦截器 interceptors
// 响应拦截器 interceptors

/**
//跳过证书验证
tr := &http.Transport{
TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

//http cookie接口
cookieJar, _ := cookiejar.New(nil)

c := &http.Client{
Jar:       cookieJar,
Transport: tr,
}

c.Get("https://192.168.0.199:8443/config/app: ")

*/

/**
func main() {
	url := "http://localhost:8080"
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil) // http client get 请求
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background()) // 获取一个上下文
	reqest = reqest.WithContext(ctx)                        // 设置当前请求的上下文
	go func() {                                             // 一定条件 比如超时等等
		cancel() // 撤销当前请求
	}()
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println("Fatal error ", err.Error()) // 取消反馈
		return
	}
	defer response.Body.Close()
}
*/
