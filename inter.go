package request

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

func (*Request) getDuration(times int) time.Duration {
	return DURATION * time.Duration(times)
}

func (rq *Request) Qs(apiUrl string, param interface{}) (*url.URL, error) {
	apiUrl = strings.TrimSpace(apiUrl)
	if apiUrl == "" {
		return nil, errors.New(url_empty)
	}
	u, e := url.Parse(rq.b + apiUrl)
	if e != nil {
		return nil, e
	}
	if param == nil {
		return u, nil
	}
	m := InterToMapStr(param)
	if len(m) == 0 {
		return u, nil
	}
	values := url.Values{}
	for k, v := range u.Query() {
		values.Set(k, v[0])
	}
	for k, v := range m {
		values.Set(k, v)
	}
	u.RawQuery = values.Encode()
	return u, nil
}
