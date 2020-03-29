package idworker

import (
	"fmt"
	"testing"
	"time"
)

// func TestGenSnowflakeID(t *testing.T) {
// 	// id: 198695121796444187,worder id: 123,sn: 27,epoch: 47372608613
// 	id := genSnowflakeId(123, 47372608613, 27)
// 	if id != 198695121796444187 {
// 		t.Fatalf("id:%v=>%v", 198695121796444187, id)
// 	}
// 	epoch := id >> 22
// 	if epoch != 47372608613 {
// 		t.Fatalf("epoch:%v=>%v", 47372608613, epoch)
// 	}
// 	worderID := id / 4096 % 1024
// 	if worderID != 123 {
// 		t.Fatalf("worderID:%v=>%v", 123, worderID)
// 	}
// 	sn := id % 4096
// 	if sn != 27 {
// 		t.Fatalf("sn:%v=>%v", 27, sn)
// 	}
// }

// func benchmarkSnowflake(b *testing.B) int64 {
// 	ch := make(chan int64, 10000*200)
// 	var wg sync.WaitGroup
// 	for i := 1; i < 100; i++ {
// 		wg.Add(1)
// 		go func(wid int64) {
// 			defer wg.Done()
// 			w, err := NewSnowflakeIdWorker(wid)
// 			if err != nil {
// 				panic(err)
// 			}
// 			for x := 0; x < 10000; x++ {
// 				id := w.NextId()
// 				ch <- id
// 			}
// 		}(int64(i))
// 	}
// 	wg.Wait()
// 	close(ch)
// 	var id int64 = -1
// 	var count int64
// LOOP:
// 	for {
// 		select {
// 		case it := <-ch:
// 			if it == 0 {
// 				break LOOP
// 			}
// 			count++
// 			if id == -1 {
// 				id = it
// 				continue LOOP
// 			}
// 			if id == it {
// 				b.Fatalf("Duplicate ID:%v", id)
// 			}
// 			id = it
// 		default:
// 			break
// 		}
// 	}
// 	return count

// }

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Log(NextId())
	}
}

func TestX(t *testing.T) {
	ep := uint64(time.Unix(4294967296-1, 0).UnixNano()) / 1e9
	idMax := (uint64(1) << 53) - 1
	id := ep & idMax
	id <<= 21
	id |= 31 << 16
	id |= 65535
	t.Log(id, "id")
	t.Log(idMax, "idMax", (1<<33)-2)
	t.Log(ep, "ep")
	t.Log(UTCMillisecond() / 1e3)
}

// 2242-03-16 20:56:32
//9007199254740990
//9007199254740991
//9007199252643839

// 3324563795345409
// 6649127649880571904
// 31322794543022085
func TestShort(t *testing.T) {
	b := SnowflakeIdBuilder{
		TimeBits:     32,
		WorkerBits:   5,
		SequenceBits: 16,
		WorkerIDGen:  WorkerWithPodIDOrHostname,
		TimeGen: func() uint64 {
			return UTCMillisecond() / 1e3 // 秒，最大2106-02-07 14:28:15
		},
	}
	gen, err := b.Build()
	if err != nil {
		panic(err)
	}
	t.Log(gen())
	t.Log(gen())
	t.Log(gen())
}

func TestLong(t *testing.T) {
	b := SnowflakeIdBuilder{
		TimeBits:     41,
		WorkerBits:   10,
		SequenceBits: 12,
		WorkerIDGen:  WorkerWithPodIDOrHostname,
		TimeGen: func() uint64 {
			epoch := UTCMillisecond()
			twepoch := uint64(1577808000000) // 2020-01-01 00:00:00
			return epoch - twepoch           // 2089-09-07 23:47:35 最大
		},
	}
	gen, err := b.Build()
	if err != nil {
		panic(err)
	}
	t.Log(gen())
	t.Log(gen())
	t.Log(gen())
}

func TestA(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			for j := 0; j < 10000; j++ {
				NextId()
			}
		})

	}
}
