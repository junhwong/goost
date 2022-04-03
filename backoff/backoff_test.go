package backoff

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))

	b := Backoff{
		Min:    time.Second,
		Max:    time.Second * 60,
		Jitter: time.Second,
	}

	for i := 0; i < 10; i++ {
		fmt.Println(b.Duration())
	}
}
func TestB(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))

	b := Backoff{
		Min:    time.Second,
		Max:    time.Second * 60,
		Jitter: time.Second,
	}
	b.Reset(10)
	for b.Backoff() {
		fmt.Println(time.Now())
	}
}
