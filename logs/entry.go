package logs

import (
	"fmt"
	"regexp"

	"github.com/junhwong/goost/logs/timestamp"
)

type Fields map[string]*Field

// Entry 单条日志记录
type Entry struct {
	Handler      Handler             // 处理器
	Time         timestamp.Timestamp // 时间戳
	Level        Level               // 日志级别
	Message      string              // 消息
	Prefix       string              // 前缀，如：app.service.module or db.mysql.name
	Tags         Fields              // 标签数据 strings([0-9a-zA-z_\-.]+), numbers(int,float,byte), bool(true,false)
	Data         Fields              // 附加数据
	Fields       []*Field
	IsPrintEntry bool
}

func (entry *Entry) handle() {
	if entry.Handler == nil {
		return
	}
	entry.Handler.Handle(entry)
	entry.Handler = nil
	// if err := e.Handler.Handle(e); err != nil {
	// 	Crash(fmt.Errorf("log: handle error %v", err))
	// }
}

func (e *Entry) log(lvl Level, args []interface{}) {

	if len(args) == 0 {
		return
	}

	e.Level = lvl
	if e.Time == 0 {
		e.Time = timestamp.Now()
	}

	defer e.handle()
	format, ok := args[0].(string)
	if ok {
		args = args[1:]
	}
	a := []interface{}{}
	for _, f := range args {
		if f, ok := f.(*Field); ok {
			e.Fields = append(e.Fields, f)
		} else {
			a = append(a, f)
		}
	}
	if len(a) == 0 {
		e.Message = format
	} else {
		e.Message = fmt.Sprintf(format, a...)
	}
}

var tagNameRule = regexp.MustCompile(`^\w[\w\.\-]+\w$`)

func checkTagName(name string) {
	if !tagNameRule.MatchString(name) {
		Crash(fmt.Errorf("log: Invalid tag name: %s", name))
	}
}

func (e *Entry) Debug(args ...interface{}) { e.log(DEBUG, args) }
func (e *Entry) Info(args ...interface{})  { e.log(INFO, args) }
func (e *Entry) Warn(args ...interface{})  { e.log(WARN, args) }
func (e *Entry) Error(args ...interface{}) { e.log(ERROR, args) }
func (e *Entry) Fatal(args ...interface{}) { e.log(FATAL, args) }
func (e *Entry) Trace(args ...interface{}) { e.log(TRACE, args) }
