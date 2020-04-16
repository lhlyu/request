package test

import (
	"github.com/lhlyu/request/v2"
	"testing"
)

const (
	avatarUrl = "https://avatars1.githubusercontent.com/u/41979509?v=4"
	apiUrl    = "https://api.github.com/users/lhlyu"
)

func TestNew(t *testing.T) {
	request.New().
		Debug().
		Get("https://cdn.jsdelivr.net/gh/lhlyu/pb@master/b/43.jpg").
		Then(func(r *request.SuperReceiver) {
			e := r.Save("./43.jpg")
			if e != nil {
				t.Log(e)
			}
		})
}
