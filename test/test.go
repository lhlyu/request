package main

import (
	"encoding/json"
	"fmt"
	"github.com/lhlyu/request"
)

func main() {
	m := map[interface{}]interface{}{
		2:     1,
		true:  "xxx",
		false: true,
		3.14:  "x",
		"hel": 1,
	}
	c := request.InterToMap(m)
	bts, _ := json.Marshal(c)
	fmt.Println(string(bts))
}
