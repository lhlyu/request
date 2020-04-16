package test

import (
	"github.com/lhlyu/request/v2"
	"testing"
)

const (
	avatarUrl = "https://avatars1.githubusercontent.com/u/41979509?v=4"
	apiUrl    = "https://api.github.com/users/lhlyu"
	picUrl    = "https://cdn.jsdelivr.net/gh/lhlyu/pb@master/b/43.jpg"
)

func TestNew(t *testing.T) {
	request.New().
		Debug().
		Get(apiUrl).
		Then(func(r *request.SuperReceiver) {
			t.Log(r.GetBodyString())
		})
}
