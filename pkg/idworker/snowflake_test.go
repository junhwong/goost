package idworker

import (
	"sync"
	"testing"
)

func TestGenSnowflakeID(t *testing.T) {
	// id: 198695121796444187,worder id: 123,sn: 27,epoch: 47372608613
	id := genSnowflakeId(123, 47372608613, 27)
	if id != 198695121796444187 {
		t.Fatalf("id:%v=>%v", 198695121796444187, id)
	}
	epoch := id >> 22
	if epoch != 47372608613 {
		t.Fatalf("epoch:%v=>%v", 47372608613, epoch)
	}
	worderID := id / 4096 % 1024
	if worderID != 123 {
		t.Fatalf("worderID:%v=>%v", 123, worderID)
	}
	sn := id % 4096
	if sn != 27 {
		t.Fatalf("sn:%v=>%v", 27, sn)
	}
}

func benchmarkSnowflake(b *testing.B) int64 {
	ch := make(chan int64, 10000*200)
	var wg sync.WaitGroup
	for i := 1; i < 100; i++ {
		wg.Add(1)
		go func(wid int64) {
			defer wg.Done()
			w, err := NewSnowflakeIdWorker(wid)
			if err != nil {
				panic(err)
			}
			for x := 0; x < 10000; x++ {
				id := w.NextId()
				ch <- id
			}
		}(int64(i))
	}
	wg.Wait()
	close(ch)
	var id int64 = -1
	var count int64
LOOP:
	for {
		select {
		case it := <-ch:
			if it == 0 {
				break LOOP
			}
			count++
			if id == -1 {
				id = it
				continue LOOP
			}
			if id == it {
				b.Fatalf("Duplicate ID:%v", id)
			}
			id = it
		default:
			break
		}
	}
	return count

}

func BenchmarkSnowflake(b *testing.B) {
	var count int64
	for i := 0; i < b.N; i++ {
		count += benchmarkSnowflake(b)
	}
	b.Logf("count:%v", count)
}
