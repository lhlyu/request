package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"
)

func T() {

}

// BenchmarkT-4   	 1000000	      2014 ns/op
// BenchmarkT-4   	 1000000	      1999 ns/op

// BenchmarkT-4   	     500	   3041264 ns/op
// BenchmarkT-4   	     500	   3034985 ns/op
// BenchmarkT-4   	  500000	      2792 ns/op

// BenchmarkT-4   	     500	   3561649 ns/op
// BenchmarkT-4   	     500	   3556508 ns/op
// BenchmarkT-4   	  300000	      3515 ns/op
func BenchmarkT(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		JsonToStr()
	}
}

type A struct {
	B int
}

func JsonToStr() {
	for i := 0; i < 10; i++ {
		a := &A{
			B: i,
		}
		json.Marshal(a)
	}
}

func fmtToStr() {
	for i := 0; i < 10; i++ {
		a := &A{
			B: i,
		}
		fmt.Sprint(a)
	}
}
