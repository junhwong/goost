package app

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type commandFn = func()

var (
	defaultCommand commandFn
	mu             sync.Mutex
	commands       = list.New()
)

func parse(args []string) map[string]string {
	arguments := make(map[string]string)
	for _, a := range args {
		arr := strings.Split(a, "=")
		if len(arr) == 1 {
			arguments[arr[0]] = "="
			continue
		}
		arguments[arr[0]] = arr[1]
	}
	return arguments
}

func PrintUsage(commandName string) {
	fmt.Printf("Usage: %s <command>\n\n", strings.Trim(envApplication.String(), `"`))
	elem := commands.Front()
	for elem != nil {
		cmd := elem.Value.(*CommandDirective)
		if len(cmd.keys) == 0 {
			elem = elem.Next()
			continue
		}
		help := cmd.usage(cmd)
		if help != "" {
			help = "\t: " + help
		}
		// if help != "" && cmd.defValue != nil && cmd.defValue != "" && cmd.defValue != false && cmd.defValue != 0 {
		// 	help += fmt.Sprintf("(default: %v)", cmd.defValue)
		// }
		fmt.Printf("%#v%s\n", strings.Join(cmd.keys, ","), help)
		elem = elem.Next()
	}

}

var commandNameRule = regexp.MustCompile(`^[a-zA-Z][\w\-\._]*`)
var commandOptionRule = regexp.MustCompile(`^\-+[a-zA-Z][\w\-\._]*`)

// Command 声明一个 CLI 指令，其执行顺序为声明时的顺序。
//
// `command` 用于定义选项和子命令，格式为：
// 	- 带子命令 `SUBCMD -SHORT,--LONG`
// 	- 默认选项 `-SHORT,--LONG`(或之一)。
//
// `defValue` 给定的初始值同时也是定义其选项的数据类型。
//
// `handle` 要执行的方法，参数为接收到的输入。
//
// `usage` 用于打印帮助时的信息，可以是 `string` 或 `UsageFormatter`。
func Command(command string, defValue interface{}, handle func(interface{}), usage ...interface{}) {
	if handle == nil {
		fatalf("handle is required: %s", command)
	}
	command = strings.TrimSpace(command)
	var part []string
	if len(command) > 0 {
		part = strings.Split(strings.TrimSpace(command), " ")
	}

	commandName := "default"
	switch {
	case len(part) == 2:
		commandName = part[0]
		part = strings.Split(part[1], ",")
	case len(part) == 1:
		if strings.HasPrefix(part[0], "-") {
			part = strings.Split(part[0], ",")
		} else {
			commandName = part[0]
			part = []string{}
		}
	case len(part) == 0:
	default:
		fatalf("invalid command: %s", command)
	}

	if !commandNameRule.MatchString(commandName) {
		fatalf("invalid command name: `%s` with `%s`", commandName, command)
	}

	options := make([]string, 0)
	switch {
	case len(part) == 2:
		if !commandOptionRule.MatchString(part[0]) || strings.HasPrefix(part[0], "--") {
			fatalf("invalid option name: `%s` part of `%s`", part[0], command)
		}
		if !commandOptionRule.MatchString(part[1]) || !strings.HasPrefix(part[1], "--") {
			fatalf("invalid option name: `%s` with `%s`", part[1], command)
		}
		if part[0] == part[1] {
			fatalf("invalid options: `%s` with `%s`", part[1], command) // TODO:
		}
		options = append(options, part[0])
		options = append(options, part[1])
	case len(part) == 1:
		if !commandOptionRule.MatchString(part[0]) {
			fatalf("invalid option name: `%v` of `%s`", part, command)
		}
		options = append(options, part[0])
	case len(part) != 0:
		fatalf("invalid options: %s", command)
	}
	var help UsageFormatter = func(*CommandDirective) string {
		return ""
	}
	if len(usage) != 0 && usage[0] != nil {
		switch v := usage[0].(type) {
		case string:
			help = func(*CommandDirective) string {
				return v
			}
		case UsageFormatter:
			help = v
		case func(*CommandDirective) string:
			help = v
		default:
			fmt.Printf("Unsupport UsageFormatter: %+v\n", v)
		}
	}
	defineCmd(commandName, options, defValue, handle, help)
}

type UsageFormatter func(*CommandDirective) string
type CommandDirective struct {
	CommandName  string
	keys         []string
	defValue     interface{}
	handle       func(interface{})
	usage        UsageFormatter
	parser       func(string) (interface{}, error)
	isOption     bool
	isDefaultcmd bool
	dataType     string
	Options      []*CommandDirective
	Children     []*CommandDirective
}

var commandMap = make(map[string]*CommandDirective)
var rootCmd = &CommandDirective{
	CommandName: "",
	Options:     make([]*CommandDirective, 0),
	Children:    make([]*CommandDirective, 0),
	isOption:    false,
}

func defineCmd(commandName string, options []string, value interface{}, handle func(interface{}), help UsageFormatter) {
	mu.Lock()
	defer mu.Unlock()

	cmd, ok := commandMap[commandName]
	if !ok {
		cmd = &CommandDirective{
			CommandName:  commandName,
			Options:      make([]*CommandDirective, 0),
			Children:     make([]*CommandDirective, 0),
			isOption:     false,
			isDefaultcmd: commandName == "default",
		}
		commandMap[commandName] = cmd
	} else if len(cmd.Options) == 0 && len(options) == 0 {
		fmt.Println(cmd)
		fatalf("command duplicated: %s", commandName)
	}

	if len(options) == 0 {
		cmd.handle = handle
		cmd.defValue = value
		cmd.usage = help
	} else {
		sub := &CommandDirective{
			CommandName: commandName,
			Options:     make([]*CommandDirective, 0),
			Children:    make([]*CommandDirective, 0),
			keys:        options,
			isOption:    true,
			defValue:    value,
			handle:      handle,
			usage:       help,
		}
		cmd.Children = append(cmd.Children, sub)
		for _, key := range options {
			key = commandName + " " + key
			if _, ok := commandMap[key]; ok {
				fatalf("command duplicated: %s", key)
			}
			commandMap[key] = sub
			cmd.Options = append(cmd.Options, sub)
		}
		cmd = sub
	}

	switch value := value.(type) {
	case string:
		cmd.dataType = "string"
		cmd.parser = func(s string) (interface{}, error) {
			if s != "" {
				return s, nil
			}
			return value, nil
		}
	case bool:
		cmd.dataType = "bool"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			return strconv.ParseBool(s)
		}
	case uint16, uint, uint32, uint64:
		cmd.dataType = "uinteger"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			return strconv.ParseUint(s, 10, 64)
		}
	case int16, int, int32, int64:
		cmd.dataType = "integer"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			return strconv.ParseInt(s, 10, 64)
		}
	case float32, float64:
		cmd.dataType = "float"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			return strconv.ParseFloat(s, 64)
		}
	case time.Duration:
		cmd.dataType = "duration"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				return time.Duration(v), nil
			}
			s = strings.ToLower(s)
			var b time.Duration
			switch {
			case strings.HasSuffix(s, "s"):
				b = time.Second
			case strings.HasSuffix(s, "m"):
				b = time.Minute
			case strings.HasSuffix(s, "h"):
				b = time.Hour
			case strings.HasSuffix(s, "d"):
				b = time.Hour * 24
			default:
				return 0, fmt.Errorf("invalid time suffix, must be one of [s,m,h,d]: %s", s)
			}
			if v, err := strconv.ParseFloat(s[:len(s)-1], 64); err == nil {
				return time.Duration(v * float64(b)), nil
			} else {
				return 0, fmt.Errorf("failed to parse duration number: %s", err)
			}
		}
	case time.Time:
		cmd.dataType = "time"
		cmd.parser = func(s string) (interface{}, error) {
			if s == "" {
				return value, nil
			}
			if v, err := time.Parse(time.RFC3339, s); err == nil {
				return v, nil
			}
			if v, err := time.Parse(time.ANSIC, s); err == nil {
				return v, nil
			}
			if v, err := time.Parse(time.RFC850, s); err == nil {
				return v, nil
			}
			if v, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
				return v, nil
			}
			if v, err := time.Parse("2006-01-02", s); err == nil {
				return v, nil
			}
			return time.Time{}, fmt.Errorf("failed to parse time: %s", s)
		}
	default:
		fatalf("the command value type not supported: %v of %s", value, cmd.CommandName)
	}
	commands.PushBack(cmd)
}

func fatalf(format string, v ...interface{}) {
	fmt.Printf("app: "+format+"\n", v...)
	if debugMode {
		debug.PrintStack()
	}
	os.Exit(-1)
}
