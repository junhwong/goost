package snowflakeid

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/junhwong/goost/runtime"
)

func WorkerWithHash(data []byte, max uint64) (uint64, error) {
	if max == 0 {
		return 0, fmt.Errorf("max cannot be 0")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("data cannot be nil or empty")
	}
	id := uint64(runtime.HashCode(data)) % max
	// id := crc32.ChecksumIEEE(data) % uint32(max)
	return uint64(id), nil
}

// WorkerWithPodIDOrHostname 获取工作节点ID.
// 算法:
//
//	hash(env.podid || hostname) % max
func WorkerWithPodIDOrHostname(d, max uint64) (uint64, uint64, error) {
	name, ok := os.LookupEnv("POD_ID")
	if !ok || name == "" {
		var err error
		name, err = os.Hostname()
		if err != nil {
			return 0, 0, err
		}
	}
	w, err := WorkerWithHash([]byte(name), max)
	if err != nil {
		return 0, 0, err
	}
	return 0, w, nil
}

// UTCMillisecond 返回 UTC 时间的毫秒时间戳
func UTCMillisecond(ctx context.Context) int64 {
	return UnixMilli(time.Now())
}

func UnixMilli(nt time.Time) int64 {
	if _, offs := nt.Zone(); offs != 0 {
		nt = nt.Add(time.Second * time.Duration(-offs))
	}
	return nt.UnixMilli()
}

var EPOCH2020 uint64 = 1577836800000 // since 2020-01-01T00:00:00Z

// func init() {
// 	time.s
// 	nt, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
// 	if err != nil {
// 		panic(err)
// 	}
// 	Epoch2020 = uint64(UnixMilli(nt))
// }
