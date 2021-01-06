package pkg

import (
	"testing"
	"time"
)

func TestRollingCounter(t *testing.T) {
	// 窗口周期为10秒
	counter := NewRollingCounter(10)

	for _, v := range []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		counter.IncrN(v)
		time.Sleep(1 * time.Second)
	}

	for i, n := range counter.Stats() {
		t.Logf("[%s]: %v\n", time.Unix(i, 0).Format("2006-01-02 15:04:05"), n.Count)
	}

	t.Logf("sum: %v, max: %v, min: %v, avg: %v", counter.Sum(), counter.Max(), counter.Min(), counter.Avg())
}
