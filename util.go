package request

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

func getStrType(s string) int {
	if json.Valid([]byte(s)) {
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
			index := strings.Index(v, "//")
			if index > -1 {
				v = v[:index]
			}
			vs := strings.TrimSpace(v)
			if vs == "" {
				continue
			}

			vArr := strings.SplitN(vs, ":", 2)
			if len(vArr) < 2 {
				continue
			}
			m[strings.TrimSpace(vArr[0])] = strings.TrimSpace(vArr[1])
		}
	case _QS:
		values, _ := url.ParseQuery(s)
		for k, v := range values {
			m[k] = v[0]
		}
	case _JSON:
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			log.Println(err)
		}

	}
	return m
}
