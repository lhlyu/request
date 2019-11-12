package benchmark

import (
	"github.com/lhlyu/request"
	"testing"
)

func T() {
	m := map[interface{}]interface{}{
		2:     1,
		true:  "xxx",
		false: true,
		3.14:  "x",
		"hel": 1,
	}
	request.InterToMap(m)
}

// BenchmarkT-4   	 1000000	      2014 ns/op
// BenchmarkT-4   	 1000000	      1999 ns/op
func BenchmarkT(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		T()
	}
}
