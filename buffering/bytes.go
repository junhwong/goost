package buffering

import "sync"

var bytesPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 1024)
		return &buf
	},
}

func PutBytes(buf []byte) {
	if buf == nil {
		return
	}
	// 将切片重置到最大容量
	buf = buf[:cap(buf)]
	bytesPool.Put(&buf) // 存储指针（非切片值）
}

// 获取切片工具函数
func GetBytes() []byte {
	ptr := bytesPool.Get().(*[]byte) // 从池中取指针
	return *ptr
}

func UseBytes(fn func(buf []byte) error) error {
	buf := GetBytes()
	defer PutBytes(buf)
	return fn(buf)
}
