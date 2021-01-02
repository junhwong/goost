package apm

// var (
// 	itemReg = regexp.MustCompile(`\$\{.+?\}`)
// 	reg     = regexp.MustCompile(`^\$\{\s*([\w]+)(\.[\w_\-]+)?(\?)?\s*(\|\s*.+?)?\s*\}$`)
// )

// type Handler interface {
// 	Handle(*Entry)
// }

// // FormatHandler 用于格式化并输出日志
// type FormatHandler struct {
// 	Out     io.Writer
// 	Format  string
// 	Getters []func(*Entry) interface{}
// 	Level   Level
// }

// func (hdr *FormatHandler) Handle(entry *Entry) {
// 	if entry.Level < hdr.Level {
// 		return
// 	}
// 	out := hdr.Out
// 	if out == nil {
// 		switch entry.Level {
// 		case DEBUG, INFO, TRACE:
// 			out = os.Stdout
// 		default:
// 			out = os.Stderr
// 		}
// 	}
// 	values := []interface{}{}
// 	for _, getter := range hdr.Getters {
// 		values = append(values, getter(entry))
// 	}
// 	_, err := fmt.Fprintf(out, hdr.Format, values...) // TODO: error handle?
// 	if err != nil {
// 		log.Printf("FormatHandler.Handle: %v", err)
// 	}
// }

// func NewFormatHandler(lvl Level, templete string, out io.Writer, parse ...func(templete string) (format string, getters []func(*Entry) interface{})) *FormatHandler {
// 	hdr := &FormatHandler{Level: lvl, Out: out}
// 	p := buildTemplete
// 	if len(parse) > 0 && parse[0] == nil {
// 		p = parse[0]
// 	}
// 	hdr.Format, hdr.Getters = p(templete)
// 	return hdr
// }

// func buildTemplete(templete string) (format string, getters []func(*Entry) interface{}) {
// 	getters = []func(*Entry) interface{}{}
// 	format = itemReg.ReplaceAllStringFunc(templete, func(s string) string {
// 		m := reg.FindAllStringSubmatch(s, -1)
// 		if len(m) != 1 || len(m[0]) <= 1 {
// 			return s
// 		}
// 		key := m[0][1]
// 		var subkey, ft string
// 		var opt bool
// 		if len(m[0]) > 2 && m[0][2] != "" {
// 			subkey = m[0][2][1:]
// 		}
// 		if len(m[0]) > 3 && m[0][3] != "" {
// 			if m[0][3] == "?" {
// 				opt = true
// 			} else {
// 				ft = m[0][3][1:]
// 			}
// 		}
// 		if len(m[0]) > 4 && m[0][4] != "" {
// 			ft = m[0][4][1:]
// 		}
// 		placeholder, getter, ok := parseTempleteItem(key, subkey, opt, strings.TrimSpace(ft))
// 		if !ok {
// 			return s
// 		}
// 		getters = append(getters, getter)
// 		return placeholder
// 	})
// 	return
// }

// func parseTempleteItem(key, subkey string, opt bool, format string) (placeholder string, getter func(*Entry) interface{}, ok bool) {
// 	ok = true
// 	switch key {
// 	case "time":
// 		if format == "" {
// 			format = time.RFC3339Nano
// 		} else {
// 			// 参考：
// 			// go time format
// 			// http://momentjs.cn/docs/#/displaying/
// 			switch format {
// 			case "ANSIC":
// 				format = time.ANSIC
// 			case "RFC3339":
// 				format = time.RFC3339
// 			case "Kitchen":
// 				format = time.Kitchen
// 			case "Stamp":
// 				format = time.Stamp
// 			case "StampMilli":
// 				format = time.StampMilli
// 			case "StampMicro":
// 				format = time.StampMicro
// 			default:
// 				plist := [][2]string{
// 					// 年
// 					[2]string{"yyyy", "2006"},
// 					[2]string{"yy", "06"},
// 					// 月
// 					[2]string{"MMMM", "January"},
// 					[2]string{"MMM", "Jan"},
// 					[2]string{"MM", "01"},
// 					[2]string{"M", "1"},
// 					// 星期
// 					[2]string{"W", "Monday"},
// 					[2]string{"w", "Mon"},
// 					// 日
// 					[2]string{"dd", "02"},
// 					[2]string{"d", "2"},
// 					// 时
// 					[2]string{"HH", "15"},
// 					[2]string{"hh", "03"},
// 					[2]string{"h", "3"},
// 					// 分
// 					[2]string{"mm", "04"},
// 					[2]string{"m", "4"},
// 					// 秒
// 					[2]string{"ss", "05"},
// 					[2]string{"s", "5"},
// 					// 毫秒/微妙/纳秒 TODO: 不显示前面的 .
// 					[2]string{"SSS", ".000000000"},
// 					[2]string{"SS", ".000000"},
// 					[2]string{"S", ".000"},
// 					// 时区
// 					[2]string{"ZZZ", "MST"},
// 					[2]string{"ZZ", "Z07:00"},
// 					[2]string{"Z", "Z0700"},
// 					// [2]string{"z", "Z"}, // TODO: 直接输出 Z
// 					// 上午/下午
// 					[2]string{"PM", "PM"},
// 					[2]string{"pm", "pm"},
// 				}
// 				for _, it := range plist {
// 					src := it[0]
// 					dst := it[1]
// 					format = strings.ReplaceAll(format, src, dst)
// 				}
// 			}
// 		}
// 		getter = func(e *Entry) interface{} {
// 			return e.Time.Format(format)
// 		}
// 		placeholder = "s"
// 	case "level":
// 		if format != "" {
// 			switch format {
// 			case "so":
// 				format = time.RFC3339Nano
// 			}
// 		}
// 		switch subkey {
// 		case "code":
// 			getter = func(e *Entry) interface{} {
// 				return uint16(e.Level)
// 			}
// 			if format == "" {
// 				format = "d"
// 			}
// 		case "short":
// 			getter = func(e *Entry) interface{} {
// 				return e.Level.Short()
// 			}
// 			if format == "" {
// 				format = "s"
// 			}
// 		default:
// 			getter = func(e *Entry) interface{} {
// 				return e.Level
// 			}
// 		}
// 		if format == "" {
// 			format = "-5s"
// 		}
// 		placeholder = format
// 	case "message":
// 		getter = func(e *Entry) interface{} {
// 			if opt && e.Message == "" {
// 				return ""
// 			}
// 			return e.Message
// 		}
// 		if format == "" {
// 			format = "s"
// 		}
// 		placeholder = format
// 	case "tags":
// 		if subkey != "" {
// 			getter = func(e *Entry) interface{} {
// 				f := e.Tags[subkey]
// 				if opt && (f == nil || f.Value == nil) {
// 					return ""
// 				}
// 				return f.Value
// 			}
// 		} else {
// 			getter = func(e *Entry) interface{} {
// 				if opt && len(e.Tags) == 0 {
// 					return ""
// 				}
// 				return e.Tags // TODO: 展开
// 			}
// 		}
// 		if format == "" {
// 			format = "v"
// 		}
// 		placeholder = format
// 	case "data":
// 		if subkey != "" {
// 			getter = func(e *Entry) interface{} {
// 				f := e.Data[subkey]
// 				if opt && (f == nil || f.Value == nil) {
// 					return ""
// 				}
// 				return f.Value
// 			}
// 		} else {
// 			getter = func(e *Entry) interface{} {
// 				if opt && len(e.Data) == 0 {
// 					return ""
// 				}
// 				return e.Data // TODO: 展开
// 			}
// 		}
// 		if format == "" {
// 			format = "v"
// 		}
// 		placeholder = format
// 	default:
// 		ok = false
// 	}
// 	placeholder = "%" + placeholder
// 	return
// }
