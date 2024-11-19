package buffering

import "sync"

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(Buffer)
		},
	}
}
func GetBuffer() *Buffer {
	buf := bufferPool.Get().(*Buffer)
	buf.Reset()
	return buf
}
func PutBuffer(buf *Buffer) { // todo 限制buf大小
	if buf == nil {
		return
	}
	bufferPool.Put(buf)
}
func UseBuffer(fn func(buf *Buffer) error) error {
	buf := GetBuffer()
	defer PutBuffer(buf)
	return fn(buf)
}
