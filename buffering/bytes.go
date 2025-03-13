package buffering

import "sync"

var bytesPool = &sync.Pool{
	New: func() any {
		return make([]byte, 0, 1024*1024)
	},
}
