package backoff

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// 回避算法
type Backoff struct {
	Min, Max time.Duration // 回避等待的时间, 默认 1s-1m
	Jitter   time.Duration // 是否启用抖动, 0则禁用

	factor  int64
	retries int64
}

// 获取本轮需要等待的时长
func (b *Backoff) Duration() time.Duration {
	//b.factor
	old := b.factor
	f := old
	if f < 1 {
		f = 1
	}
	if b.Min < 1 {
		b.Min = time.Second
	}
	if b.Max < 1 {
		b.Max = time.Minute
	}
	d := b.Min * time.Duration(f)
	f <<= 1
	if b.Jitter > 0 {
		j := time.Duration(rand.Int63n(10))
		if j <= 0 {
			j = 1
		}
		d += time.Duration(rand.Int63n(f)) * b.Jitter / j
	}

	for {
		if b := atomic.CompareAndSwapInt64(&b.factor, old, f); b {
			break
		}
	}

	if d < b.Min {
		return b.Min
	}
	if d > b.Max {
		return b.Max
	}

	return time.Duration(d)
}

// 等待。返回true表示未到达最大尝试次数
func (b *Backoff) Backoff() bool {
	if b.retries <= 0 {
		return false
	}
	time.Sleep(b.Duration())
	r := b.retries
	rr := r - 1
	for {
		if b := atomic.CompareAndSwapInt64(&b.retries, r, rr); b {
			break
		}
	}
	return true
}

func (b *Backoff) Reset(retries int64) {
	old := b.factor
	for {
		if b := atomic.CompareAndSwapInt64(&b.factor, old, 1); b {
			break
		}
	}
	r := b.retries
	for {
		if b := atomic.CompareAndSwapInt64(&b.retries, r, retries); b {
			break
		}
	}
}
