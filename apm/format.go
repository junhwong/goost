package apm

import (
	"bytes"

	"github.com/junhwong/goost/apm/field"
)

// Formatter 表示一个格式化器。
type Formatter interface {
	// Format 格式化一条日志。
	//
	// 注意：不要缓存 `entry`, `dest` 对象，因为它们可能是池化对象。
	Format(entry *field.Field, dest *bytes.Buffer) (err error)
}
