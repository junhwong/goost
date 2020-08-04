package idworker

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNextId(t *testing.T) {
	for i := 0; i < 50000; i++ {
		NextId()
	}
	t.Log("done")
}

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Log(NextId())
	}
}

func TestA(t *testing.T) {
	// 生产10000个ID大约2秒(含文件IO)
	t.Parallel()
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			f, err := os.OpenFile("./tmp.id", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			for j := 0; j < 100; j++ {
				id := NextId()

				_, err = f.WriteString(strconv.FormatUint(uint64(id), 10) + "\n")
				if err != nil {
					t.Fatal(err)
				}
				f.Sync()
			}
		})
	}
}

func TestParseSnowflakeId(t *testing.T) {
	id := NextId()
	t.Log(id)
	epoch, workerId, seq, err := ParseSnowflakeId(DefaultLayout, id)
	t.Log("parse:", epoch, workerId, seq, err)
	t.Log(time.Unix(int64(epoch/1000), int64(epoch-epoch/1000)).UTC())
}
