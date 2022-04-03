package apm

import (
	"strconv"
	"strings"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
	"github.com/spf13/cast"
)

type Entry = field.Fields

func GetLevel(entry Entry) (lvl level.Level) {
	if entry == nil {
		return
	}
	if f := entry.Get(LevelKey); f != nil {
		lvl = level.FromInt(cast.ToInt(f.Value))
	}
	return
}

func genCodefile(method string, file string, line int) string {
	if method == "main.main" {
		method = "main.go"
		file = ""
	}
	i := strings.LastIndex(file, "/")
	if i > 0 {
		file = file[i:]
		i = strings.LastIndex(method, "/")
		pkg := method
		if i > 0 {
			pkg = method[i:]
			method = method[:i]
		} else {
			method = ""
		}
		pkg = strings.SplitN(pkg, ".", 2)[0]
		method = method + pkg + file
	}

	if method != "" && line != 0 {
		method += ":" + strconv.Itoa(line)
	}
	return method
}
