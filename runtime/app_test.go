package runtime

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	builder := lifecycle{}
	builder.Append(func(ctx context.Context) {
		time.Sleep(time.Microsecond * 5)
		fmt.Println("Append 1")
	})
	builder.AppendServing(func(ctx context.Context, onStarted func()) {
		onStarted()
		time.Sleep(time.Microsecond * 3)
		fmt.Println("AppendServing 1")
	})
	builder.Append(func(ctx context.Context) {
		time.Sleep(time.Microsecond * 2)
		fmt.Println("Append 2")
	})

	// contexts := []*hookCtx{}
	// for range builder {
	// 	next, cancel := context.WithCancel(start)
	// 	contexts = append(contexts, &hookCtx{
	// 		ctx:    next,
	// 		cancel: cancel,
	// 	})
	// }
	// n := len(contexts) - 1
	// fmt.Printf("n: %v\n", n)

	// for i, h := range builder {
	// 	h.hookCtx = contexts[n-i]
	// 	if h.serving {
	// 		wg.Add(1)
	// 	}
	// }

	// var servingCancelOnce sync.Once
	// // for _, h := range builder {
	// // 	if h.serving {
	// // 		wg.Add(1)
	// // 	}
	// // 	go h.doRun(&wg, &servingCancelOnce)
	// // }

	// var next func()
	// var i = 0
	// var mu sync.Mutex
	// next = func() {
	// 	mu.Lock()
	// 	// fmt.Printf("i: %v\n", i)
	// 	if i > n {
	// 		mu.Unlock()
	// 		return
	// 	}

	// 	h := builder[i]
	// 	i++
	// 	mu.Unlock()

	// 	h.doRun(&wg, &servingCancelOnce, next)

	// }
	//
	start := context.TODO()
	var wg sync.WaitGroup
	var m sync.Map
	next := builder.build(start, &wg, func(s string) {}, &m)
	next()

	wg.Wait()
	fmt.Println("wg all done")
	time.Sleep(time.Second)

}

func TestRef(t *testing.T) {
	var f any

	// tv := reflect.TypeOf(f)
	f = struct{}{}
	f = func() {}
	rv := reflect.ValueOf(f)
	if !rv.IsValid() || rv.Kind() != reflect.Func {
		return
	}
	fn := runtime.FuncForPC(rv.Pointer()).Name()
	fmt.Printf("fn: %v\n", fn)
}
