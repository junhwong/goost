package apm

import (
	"context"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/cast"
)

// 对标准库 `runtime.Caller` 的封装
func Caller(depth int) (info CallerInfo) {
	doCaller(depth, &info)
	return
}

func doCaller(depth int, info *CallerInfo) {
	info.depth = depth + 1
	info.pc, info.File, info.Line, info.Ok = runtime.Caller(info.depth)

	if info.Ok {
		info.Method = runtime.FuncForPC(info.pc).Name()
		// info.Method, info.Package = split(runtime.FuncForPC(info.pc).Name())
	}
	// info.File, info.Path = split(info.File)
	return
}

func split(s string) (string, string) {
	i := strings.LastIndex(s, "/")
	if i > 0 {
		return s[i+1:], s[:i]
	}
	return s, ""
}

// 函数调用的名称等简单信息
type CallerInfo struct {
	// Path    string
	File string
	// Package string
	Method string
	Line   int

	depth int
	pc    uintptr
	Ok    bool
}

func (info CallerInfo) Caller() string {
	p := info.File
	f := ""
	if i := strings.LastIndex(p, "/"); i > 0 {
		f = p[i+1:]
		p = p[:i]
	}
	if i := strings.LastIndex(p, "/"); i > 0 {
		p = p[i+1:]
	}
	p = path.Join(p, f)
	if info.Line > 0 {
		p += ":" + cast.ToString(info.Line)
	}
	return p
}

func WithCaller(ctx context.Context, depth ...int) context.Context {
	d := 2
	if len(depth) > 0 {
		d = depth[len(depth)-1]
	}
	info := Caller(d)
	if setter, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		setter.Set(callerContextKey, &info)
	} else {
		ctx = context.WithValue(ctx, callerContextKey, &info)
	}
	return ctx
}
func CallerFrom(ctx context.Context) *CallerInfo {
	obj, _ := ctx.Value(callerContextKey).(*CallerInfo)
	return obj
}
