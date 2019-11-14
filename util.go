package request

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// 下划线转连接符
func underlineToConnector(s string) string {
	return strings.Replace(s, "_", "-", -1)
}

// 小驼峰
func toSmallTitle(s string) string {
	if s == "" {
		return ""
	}
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	re, _ := regexp.Compile("\\s+")
	s = re.ReplaceAllString(s, "")
	if len(s) == 0 {
		return ""
	}
	s = strings.ToLower(s[0:1]) + s[1:]
	return s
}

func UploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return nil, err
}



func getStrType(s string) int {
	if json.Valid([]byte(s)){
		return _JSON
	}
	if strings.Contains(s, ":") {
		return _KV_LINE
	} else if strings.Contains(s, "=") {
		return _QS
	}
	return -1
}

func strToMSI(s string, strType int) MSI {
	m := make(MSI)
	switch strType {
	case _KV_LINE:
		sArr := strings.Split(s, "\n")
		for _, v := range sArr {
			vArr := strings.SplitN(v, ":", 2)
			m[strings.TrimSpace(vArr[0])] = strings.TrimSpace(vArr[1])
		}
	case _QS:
		values, _ := url.ParseQuery(s)
		for k, v := range values {
			m[k] = v[0]
		}
	case _JSON:
		if err := json.Unmarshal([]byte(s),&m);err != nil{
			log.Println(err)
		}

	}
	return m
}


