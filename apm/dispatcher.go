package apm

import (
	"sync/atomic"
	"time"

	"github.com/junhwong/goost/apm/field"
)

func loop() {
	for e := range queue {
		if e == nil {
			return
		}
		handleEntry(e)
		queuewg.Done()
	}
}

func handleEntry(entry *field.Field) {
	if entry == nil {
		return
	}

	// var once sync.Once
	var release = func() {
		if entry == nil {
			return
		}
		field.Release(entry)
	}

	handlers := GetHandlers()
	if len(handlers) == 0 {
		release()
		return
	}

	size := int32(handlers.Len())
	var crt atomic.Int32

	// defer release()?

	var next func()
	next = func() {
		i := crt.Add(1) - 1
		if i >= size {
			release()
			return
		}
		h := handlers[i]
		h.Handle(entry, next, release)
	}
	next()
}

func Dispatch(e *field.Field) {
	if e == nil {
		return
	}
	queuewg.Add(1)
	queue <- e
}

func Flush() {
	time.Sleep(time.Nanosecond) // sync: WaitGroup is reused before previous Wait has returned
	queuewg.Wait()
}
