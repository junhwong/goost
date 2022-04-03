package apm

import "sort"

// 日志项处理器
type Handler interface {

	// 优先级. 值越大越优先
	Priority() int

	// 处理日志
	Handle(entry Entry, next func())
}

type handlerSlice []Handler

func (x handlerSlice) Len() int           { return len(x) }
func (x handlerSlice) Less(i, j int) bool { return x[i].Priority() > x[j].Priority() }
func (x handlerSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x handlerSlice) Sort()              { sort.Sort(x) }
