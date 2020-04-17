package test

import (
	"github.com/lhlyu/request/v2"
	"testing"
	"time"
)

const (
	apiUrl = "https://api.github.com/users/lhlyu"
	picUrl = "https://cdn.jsdelivr.net/gh/lhlyu/pb@master/b/43.jpg"
)

func TestNew(t *testing.T) {
	body := request.New().
		Debug().
		Get(apiUrl).GetBodyString()
	t.Log(body)
}

type User struct {
	Id int `json:"id"`
}

func TestSend(t *testing.T) {
	body := request.New().
		Timeout(time.Second*10).
		Send(`{"id":123}`).
		Send(&User{1}).
		Send("id=1").
		Post("http://localhost:8080/article/%d", 1).GetBodyString()
	t.Log(body)
}

func dosomething() {
	time.Sleep(time.Second)
}

// 异步请求
func TestAsynch(t *testing.T) {
	r := request.New().
		Asynch(). // 开启异步
		Get(apiUrl)

	// 可以做别的事
	dosomething()

	// 等待
	body := r.Await().GetBodyString()

	t.Log(body)
}

// 保存文件
func TestSave(t *testing.T) {
	err := request.New().
		Get(picUrl).Save("./test.jpg")
	t.Log(err)
}

// 链式调用
func TestThen(t *testing.T) {
	request.New().
		Get(apiUrl).Then(func(r request.Receiver) {
		t.Log(r.GetStatus())
	}).Then(func(r request.Receiver) {
		t.Log(r.GetBodyString())
	}).Then(func(r request.Receiver) {
		u := &User{}
		r.BodyUnmarshal(u)
		t.Log(u.Id)
		if u.Id > 100000 {
			r.Stop() // 会导致后面的操作停止
		}
	}).Then(func(r request.Receiver) {
		// 正则匹配
		vals := r.BodyCompile("\\d+")
		for _, v := range vals {
			t.Log(v)
		}
	})
}
