package request

const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

const (
	TypeJSON       = "json"
	TypeXML        = "xml"
	TypeUrlencoded = "urlencoded"
	TypeForm       = "form"
	TypeFormData   = "form-data"
	TypeHTML       = "html"
	TypeText       = "text"
	TypeMultipart  = "multipart"
	TypeStream     = "stream"
)

var Types = map[string]string{
	TypeJSON:       "application/json",
	TypeXML:        "application/xml",
	TypeForm:       "application/x-www-form-urlencoded",
	TypeFormData:   "application/x-www-form-urlencoded",
	TypeUrlencoded: "application/x-www-form-urlencoded",
	TypeHTML:       "text/html",
	TypeText:       "text/plain",
	TypeMultipart:  "multipart/form-data",
	TypeStream:     "application/octet-stream", // 用于文件下载
}
