package apm

import (
	"strings"
	"sync"
)

var (
	trimFieldPrefix = []string{"apm."}
	settingsMu      sync.Mutex
)

func AddTrimFieldPrefix(prefix ...string) {
	if len(prefix) == 0 {
		return
	}
	settingsMu.Lock()
	tmp := append(trimFieldPrefix, prefix...)
	settingsMu.Unlock()
	trimFieldPrefix = tmp
}

func TrimFieldNamePrefix(s string) string {
	prefixs := trimFieldPrefix // copy ptr
	for _, prefix := range prefixs {
		s = strings.TrimSpace(strings.TrimPrefix(s, prefix))
	}
	return s
}
