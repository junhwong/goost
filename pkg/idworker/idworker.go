package idworker

// IdWorker 用于生成分布式ID.
type IdWorker interface {
	// NextId 生成并返回下一个ID.
	NextId() int64
}
