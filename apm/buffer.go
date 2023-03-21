package apm

import (
	"bytes"
	"sync"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}
func GetBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}
func PutBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	bufferPool.Put(buf)
}
func UseBuffer(fn func(buf *bytes.Buffer) error) error {
	buf := GetBuffer()
	defer PutBuffer(buf)
	return fn(buf)
}
