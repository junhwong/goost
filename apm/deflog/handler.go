package deflog

import (
	"sort"

	"github.com/junhwong/goost/apm"
)

type handlerSlice []apm.Handler

func (x handlerSlice) Len() int           { return len(x) }
func (x handlerSlice) Less(i, j int) bool { return x[i].Priority() > x[j].Priority() }
func (x handlerSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x handlerSlice) Sort()              { sort.Sort(x) }
func (x handlerSlice) handle(entry apm.Entry) {
	size := x.Len()
	crt := 0
	var next func()
	next = func() {
		if crt >= size {
			return
		}
		h := x[crt]
		crt++
		h.Handle(entry, next)
	}
	next()
}
