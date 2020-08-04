package idworker

import (
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
//		hash(env.podid || hostname) % max
func WorkerWithPodIDOrHostname(max uint64) (uint64, error) {
	name, ok := os.LookupEnv("POD_ID")
	if !ok || name == "" {
		var err error
		name, err = os.Hostname()
		if err != nil {
			return 0, err
		}
	}
	return WorkerWithHash([]byte(name), max)
}

// UTCMillisecond 返回 UTC 时间的毫秒时间戳
func UTCMillisecond() uint64 {
	nt := time.Now()
	if _, offs := nt.Zone(); offs != 0 {
		nt = nt.Add(time.Second * time.Duration(-offs))
	}
	return uint64(nt.UnixNano()) / 1e6
}
