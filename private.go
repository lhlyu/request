package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"moul.io/http2curl"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func cloneMapArray(old map[string][]string) map[string][]string {
	newMap := make(map[string][]string, len(old))
	for k, vals := range old {
		newMap[k] = make([]string, len(vals))
		for i := range vals {
			newMap[k][i] = vals[i]
		}
	}
	return newMap
}
func shallowCopyData(old map[string]interface{}) map[string]interface{} {
	if old == nil {
		return nil
	}
	newData := make(map[string]interface{})
	for k, val := range old {
		newData[k] = val
	}
	return newData
}
func shallowCopyDataSlice(old []interface{}) []interface{} {
	if old == nil {
		return nil
	}
	newData := make([]interface{}, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyFileArray(old []File) []File {
	if old == nil {
		return nil
	}
	newData := make([]File, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyCookies(old []*http.Cookie) []*http.Cookie {
	if old == nil {
		return nil
	}
	newData := make([]*http.Cookie, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyErrors(old []error) []error {
	if old == nil {
		return nil
	}
	newData := make([]error, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}

func copyRetryable(old superAgentRetryable) superAgentRetryable {
	newRetryable := old
	newRetryable.RetryableStatus = make([]int, len(old.RetryableStatus))
	for i := range old.RetryableStatus {
		newRetryable.RetryableStatus[i] = old.RetryableStatus[i]
	}
	return newRetryable
}

func (s *SuperAgent) clearSuperAgent() {
	if s.DoNotClearSuperAgent {
		return
	}
	s.Url = ""
	s.Method = ""
	s.Header = http.Header{}
	s.Data = make(map[string]interface{})
	s.SliceData = []interface{}{}
	s.QueryData = url.Values{}
	s.FileData = make([]File, 0)
	s.BounceToRawString = false
	s.RawString = ""
	s.ForceType = ""
	s.TargetType = TypeJSON
	s.Cookies = make([]*http.Cookie, 0)
	s.Errors = nil
}

func (s *SuperAgent) queryStruct(content interface{}) *SuperAgent {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				var queryVal string
				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}
					queryVal = string(j)
				}
				s.QueryData.Add(k, queryVal)
			}
		}
	}
	return s
}

func (s *SuperAgent) queryString(content string) *SuperAgent {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		if queryData, err := url.ParseQuery(content); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					s.QueryData.Add(k, string(queryValue))
				}
			}
		} else {
			s.Errors = append(s.Errors, err)
		}
	}
	return s
}

func (s *SuperAgent) queryMap(content interface{}) *SuperAgent {
	return s.queryStruct(content)
}

func (s *SuperAgent) isRetryableRequest(resp *http.Response) bool {
	if s.Retryable.Enable && s.Retryable.Attempt < s.Retryable.RetryerCount && contains(resp.StatusCode, s.Retryable.RetryableStatus) {
		time.Sleep(s.Retryable.RetryerTime)
		s.Retryable.Attempt++
		return false
	}
	return true
}

func (s *SuperAgent) makeRequest() (*http.Request, error) {
	var (
		req           *http.Request
		contentType   string
		contentReader io.Reader // body 内容
		err           error
	)

	if s.Method == "" {
		return nil, errors.New("No method specified")
	}

	switch s.TargetType {
	case TypeJSON:
		var contentJson []byte
		if s.BounceToRawString {
			contentJson = []byte(s.RawString)
		} else if len(s.Data) != 0 {
			contentJson, _ = json.Marshal(s.Data)
		} else if len(s.SliceData) != 0 {
			contentJson, _ = json.Marshal(s.SliceData)
		}
		if contentJson != nil {
			contentReader = bytes.NewReader(contentJson)
			contentType = "application/json"
		}
	case TypeForm, TypeFormData, TypeUrlencoded:
		var contentForm []byte
		if s.BounceToRawString || len(s.SliceData) != 0 {
			contentForm = []byte(s.RawString)
		} else {
			formData := changeMapToURLValues(s.Data)
			contentForm = []byte(formData.Encode())
		}
		if len(contentForm) != 0 {
			contentReader = bytes.NewReader(contentForm)
			contentType = "application/x-www-form-urlencoded"
		}
	case TypeText:
		if len(s.RawString) != 0 {
			contentReader = strings.NewReader(s.RawString)
			contentType = "text/plain"
		}
	case TypeXML:
		if len(s.RawString) != 0 {
			contentReader = strings.NewReader(s.RawString)
			contentType = "application/xml"
		}
	case TypeMultipart:
		var (
			buf = &bytes.Buffer{}
			mw  = multipart.NewWriter(buf)
		)

		if s.BounceToRawString {
			fieldName := s.Header.Get("data_fieldname")
			if fieldName == "" {
				fieldName = "data"
			}
			fw, _ := mw.CreateFormField(fieldName)
			fw.Write([]byte(s.RawString))
			contentReader = buf
		}

		if len(s.Data) != 0 {
			formData := changeMapToURLValues(s.Data)
			for key, values := range formData {
				for _, value := range values {
					fw, _ := mw.CreateFormField(key)
					fw.Write([]byte(value))
				}
			}
			contentReader = buf
		}

		if len(s.SliceData) != 0 {
			fieldName := s.Header.Get("json_fieldname")
			if fieldName == "" {
				fieldName = "data"
			}
			h := make(textproto.MIMEHeader)
			fieldName = strings.Replace(strings.Replace(fieldName, "\\", "\\\\", -1), `"`, "\\\"", -1)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, fieldName))
			h.Set("Content-Type", "application/json")
			fw, _ := mw.CreatePart(h)
			contentJson, err := json.Marshal(s.SliceData)
			if err != nil {
				return nil, err
			}
			fw.Write(contentJson)
			contentReader = buf
		}

		if len(s.FileData) != 0 {
			for _, file := range s.FileData {
				fw, _ := mw.CreateFormFile(file.Fieldname, file.Filename)
				fw.Write(file.Data)
			}
			contentReader = buf
		}

		mw.Close()

		if contentReader != nil {
			contentType = mw.FormDataContentType()
		}
	default:
		return nil, errors.New("TargetType '" + s.TargetType + "' could not be determined")
	}

	if req, err = http.NewRequest(s.Method, s.Url, contentReader); err != nil {
		return nil, err
	}
	for k, vals := range s.Header {
		for _, v := range vals {
			req.Header.Add(k, v)
		}

		if strings.EqualFold(k, "Host") {
			req.Host = vals[0]
		}
	}

	if len(contentType) != 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置url参数
	q := req.URL.Query()
	for k, v := range s.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()

	if s.BasicAuth != struct{ Username, Password string }{} {
		req.SetBasicAuth(s.BasicAuth.Username, s.BasicAuth.Password)
	}

	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func (s *SuperAgent) getResponseBytes() (*http.Response, []error) {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	if len(s.Errors) != 0 {
		return nil, s.Errors
	}
	switch s.ForceType {
	case TypeJSON, TypeForm, TypeXML, TypeText, TypeMultipart:
		s.TargetType = s.ForceType
	default:
		contentType := s.Header.Get("Content-Type")
		for k, v := range Types {
			if contentType == v {
				s.TargetType = k
			}
		}
	}

	if len(s.Data) != 0 && len(s.SliceData) != 0 {
		s.BounceToRawString = true
	}

	req, err = s.makeRequest()
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, s.Errors
	}

	if !DisableTransportSwap {
		s.Client.Transport = s.Transport
	}

	if s.isDebug {
		dump, err := httputil.DumpRequest(req, true)
		s.logger.SetPrefix("[http] ")
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Request: %s", string(dump))
		}
	}

	if s.CurlCommand {
		curl, err := http2curl.GetCurlCommand(req)
		s.logger.SetPrefix("[curl] ")
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("CURL command line: %s", curl)
		}
	}
	resp, err = s.Client.Do(req)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, s.Errors
	}

	if s.isDebug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}

	if err != nil {
		return nil, []error{err}
	}
	return resp, nil
}
