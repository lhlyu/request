package main

import "github.com/lhlyu/request"

func main() {
	rq := request.NewRequest()
	rq.SetDebug(true)
	rq.DoHttp()

}

