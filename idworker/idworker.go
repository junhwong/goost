package idworker

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	defaultGen IDGenerator
)

func init() {
	rand.Seed(time.Now().UnixNano())
	gen, err := NewLongSnowflakeIdBuilder().Build()
	if err != nil {
		panic(err)
	}
	defaultGen = gen
}

// ID 表示一个分布式唯一标识
type ID string // 数据库主键，配合模糊查询

type IDGenerator func() ID

// NextId 返回一个新的ID
func NextId() ID {
	return defaultGen()
}

func SetDefault(gen IDGenerator) error {
	if gen == nil {
		return fmt.Errorf("gen cannot be nil")
	}
	defaultGen = gen
	return nil
}