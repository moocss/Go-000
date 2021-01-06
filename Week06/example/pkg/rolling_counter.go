package pkg

import (
	"sync"
	"time"
)

// RollingCounter 滑动窗口计数器
type RollingCounter struct {
	buckets  map[int64]*Bucket // 秒级为单位的桶, 字典的key为时间戳
	interval int64             // 时间周期

	mux sync.RWMutex // 慎重使用 *sync.RWMutex 不可复制类型的变量: https://www.haohongfan.com/post/2020-12-22-struct-sync/
}

// Bucket 为计数桶
type Bucket struct {
	Count int64
}

// NewRollingCounter 实例化滑动窗口计数器
func NewRollingCounter(interval int64) *RollingCounter {
	return &RollingCounter{
		buckets:  make(map[int64]*Bucket),
		interval: interval,
	}
}

// GetCurrentBucket 获取当前桶
func (rc *RollingCounter) getCurrentBucket() *Bucket {
	now := time.Now().Unix()

	// 判断当前时间是否有桶存在, 有就直接返回
	if b, ok := rc.buckets[now]; ok {
		return b
	}

	// 否则, 创建新桶
	b := new(Bucket)
	rc.buckets[now] = b

	return b
}

func (rc *RollingCounter) removeOldBuckets() {
	t := time.Now().Unix() - rc.interval
	for timestamp := range rc.buckets {
		if timestamp <= t {
			delete(rc.buckets, timestamp)
		}
	}
}

// Incr 累加计数器
func (rc *RollingCounter) Incr() {
	rc.IncrN(1)
}

// IncrN 累加计数器
func (rc *RollingCounter) IncrN(i int64) {
	rc.mux.Lock()
	defer rc.mux.Unlock()

	bucket := rc.getCurrentBucket()
	bucket.Count += i
	rc.removeOldBuckets()
}

// Sum 累计, 将过去N秒内各窗口的数值相加
func (rc *RollingCounter) Sum() int64 {
	var sum int64
	t := time.Now().Unix() - rc.interval

	rc.mux.RLock()
	defer rc.mux.RUnlock()

	for timestamp, bucket := range rc.buckets {
		if timestamp >= t {
			sum += bucket.Count
		}
	}

	return sum
}

// Max 获取过去N秒内的最大值
func (rc *RollingCounter) Max() int64 {
	var max int64
	t := time.Now().Unix() - rc.interval

	rc.mux.RLock()
	defer rc.mux.RUnlock()

	for timestamp, bucket := range rc.buckets {
		if timestamp >= t {
			if bucket.Count > max {
				max = bucket.Count
			}
		}
	}

	return max
}

// Min
func (rc *RollingCounter) Min() int64 {
	var min int64
	t := time.Now().Unix() - rc.interval

	rc.mux.RLock()
	defer rc.mux.RUnlock()

	for timestamp, bucket := range rc.buckets {
		if timestamp >= t {
			if min == 0 {
				min = bucket.Count
				continue
			}
			if bucket.Count < min {
				min = bucket.Count
			}
		}
	}

	return min
}

// Avg 获取过去N秒内所有窗口的平均值。
func (rc *RollingCounter) Avg() float64 {
	return float64(rc.Sum()) / float64(rc.interval)
}

// Stats
func (rc *RollingCounter) Stats() map[int64]*Bucket {
	return rc.buckets
}
