package snowflakeid

import (
	"context"
	"runtime"
	"sync"
	"testing"
)

func TestGenerateID(t *testing.T) {
	topologyIn := int64(5)
	workerIn := int64(6)
	tsIn := int64(1716316122622)
	c := Builder{
		Layout: Layout{
			TimeBits:     41,
			TopologyBits: 5,
			WorkerBits:   5,
			SequenceBits: 12,
			Epoch:        10,
		},
		UseRandSeed: false,
		WorkerGen: func(_, _ uint64) (topolog uint64, worker uint64, err error) {
			worker = uint64(workerIn)
			topolog = uint64(topologyIn)
			return
		},
		TimeGen: func(_ context.Context) int64 {
			return tsIn
		},
	}

	gen, err := c.Build()
	if err != nil {
		t.Fatal(err)
	}
	id := gen(context.TODO())

	epoch, topology, worker, seq, err := Parse(c.Layout, id)
	if err != nil {
		t.Fatal(err)
	}

	if epoch != tsIn {
		t.Fatal(epoch)
	}
	if topology != topologyIn {
		t.Fatal(topology)
	}
	if worker != workerIn {
		t.Fatal(worker)
	}
	if seq != 0 {
		t.Fatal()
	}

	// ts := time.UnixMilli(epoch)
	// fmt.Printf("ts: %v\n", ts)
}

func TestDuplicate(t *testing.T) {

	count := 10000
	t.Parallel()
	var wg sync.WaitGroup
	cpu := runtime.NumCPU()
	ch := make(chan int64, cpu*count*2)
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// defer func() {
			// 	fmt.Printf("i exit: %v\n", i)
			// }()
			for j := 0; j < count; j++ {
				// fmt.Printf("i: %v-%v\n", i, j)
				ch <- GenerateID(context.TODO())
				// fmt.Printf("j: %v-%v\n", i, j)
			}
		}()
	}
	wg.Wait()

	// fmt.Println("all exit")

	m := map[int64]struct{}{}
	var last int64
	c := 0
	for {
		select {
		case i := <-ch:
			if _, ok := m[i]; ok {
				t.Fatal("重复:", i, len(m))
			}
			m[i] = struct{}{}
			if i < last {
				c++
			}

			last = i
		default:
			if len(m) != cpu*count {
				t.Fatal("个数不及预期")
			}
			t.Log("顺序混乱", c) // todo 优化?
			return
		}
	}

}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < 10000; i++ {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				id := GenerateID(context.TODO())
				if id < 1 {
					b.Fatal()
				}
			}
		})
	}
}
