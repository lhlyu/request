package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

type SuperReceiver struct {
	Resp *http.Response
	Body []byte
	Errs []error

	once sync.Once
	resp chan *http.Response
}

func newSuperReceiver(isAsynch bool) *SuperReceiver {
	resp := &SuperReceiver{}
	if isAsynch {
		resp.resp = make(chan *http.Response, 0)
	}
	return resp
}

func (s *SuperReceiver) GetResponse() *http.Response {
	s.Resp.Body = ioutil.NopCloser(bytes.NewBuffer(s.GetBody()))
	return s.Resp
}

func (s *SuperReceiver) GetBody() []byte {
	s.once.Do(func() {
		if s.Resp == nil {
			s.Body = nil
			return
		}
		byts, e := ioutil.ReadAll(s.Resp.Body)
		if e != nil {
			s.Errs = append(s.Errs, e)
			return
		}
		defer s.Resp.Body.Close()
		s.Body = byts
		s.Resp.Body = ioutil.NopCloser(bytes.NewBuffer(byts))
	})
	return s.Body
}

func (s *SuperReceiver) GetBodyString() string {
	return string(s.GetBody())
}

func (s *SuperReceiver) BodyUnmarshal(v interface{}) error {
	return json.Unmarshal(s.GetBody(), v)
}

func (s *SuperReceiver) BodyCompile(pattern string) []string {
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(s.GetBodyString(), -1)
}

func (r *SuperReceiver) GetStatus() int {
	if r.Resp == nil {
		return -1
	}
	return r.Resp.StatusCode
}

func (s *SuperReceiver) Await() Receiver {
	if v, ok := <-s.resp; ok {
		s.Resp = v
	}
	return s
}

// 保存文件
func (s *SuperReceiver) Save(path string) error {
	byts := s.GetBody()
	if len(byts) == 0 {
		return nil
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(byts))
	return err
}

func (s *SuperReceiver) Then(fn func(r Receiver)) Receiver {
	fn(s)
	return s
}

func (s *SuperReceiver) Error() error {
	if len(s.Errs) > 0 {
		return s.Errs[0]
	}
	return nil
}

func (s *SuperReceiver) Errors() []error {
	return s.Errs
}
