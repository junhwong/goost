package runtime

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	builder := lifecycle{}
	builder.Append(func(ctx context.Context) {
		time.Sleep(time.Microsecond * 5)
		fmt.Println("Append 1")
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
	// start := context.TODO()
	// var wg sync.WaitGroup
	// var m sync.Map
	builder.Wait() //start, &wg, func(s string) {}, &m
	// next()

	// wg.Wait()
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

func TestLifecycle(t *testing.T) {

	l, done := NewLifecycle(context.TODO())
	defer done()
	b := bytes.NewBuffer(make([]byte, 0))
	l.Append(func(ctx context.Context) {
		fmt.Fprint(b, "1")
		<-ctx.Done()
		fmt.Fprint(b, "1")
	})
	l.Append(func(ctx context.Context) {
		fmt.Fprint(b, "2")
		time.Sleep(time.Millisecond * 100)
		fmt.Fprint(b, "2")
	})
	l.Append(func(ctx context.Context) {
		fmt.Fprint(b, "3")
		<-ctx.Done()
		fmt.Fprint(b, "3")
	})

	l.Wait()

	if b.String() != "123231" {
		t.Fatal(b.String())
	}
}
