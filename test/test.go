package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	for i := 0; i < 20; i++ {
		fmt.Println(GetRandomString(8))
	}

}

func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		time.Sleep(time.Nanosecond)
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
