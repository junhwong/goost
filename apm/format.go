package apm

import (
	"bytes"
)

// Formatter 表示一个格式化器。
type Formatter interface {
	// Format 格式化一条日志。
	//
	// 注意：不要缓存 `entry`, `dest` 对象，因为它们是池化对象。
	Format(entry Entry, dest *bytes.Buffer) (err error)
}
