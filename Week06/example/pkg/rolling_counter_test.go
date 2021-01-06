package pkg

import (
	"testing"
	"time"
)

func TestRollingCounter(t *testing.T) {
	counter := NewRollingCounter(10)
	for i, v := range []int64{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		t.Logf("输出: %v, %v", i, v)
		counter.IncrN(v)
		time.Sleep(1 * time.Second)
	}
	t.Logf("sun is %v", counter.Sum())
}
