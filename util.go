package request

import (
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

const (
	_KV_LINE = iota
	_QS
)

func getStrType(s string) int {
	if strings.Contains(s, ":") {
		return _KV_LINE
	} else if strings.Contains(s, "=") {
		return _QS
	}
	return -1
}

func strToMap(s string, strType int) map[string]string {
	m := make(map[string]string)
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
	}
	return m
}
