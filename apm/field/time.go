package field

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var CST *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	CST = loc
}

func ParseTimeZone(s string) (*time.Location, error) {
	if len(s) == 0 {
		return time.Local, nil
	}
	switch strings.ToLower(s) {
	case "utc":
		return time.UTC, nil
	case "local":
		return time.Local, nil
	case "cst", "asia/shanghai":
		return CST, nil
	default:
		loc, err := time.LoadLocation(s)
		if err != nil {
			return nil, err
		}
		return loc, nil
	}
}

var timeLayoutMap = map[string][]string{
	"rfc3339":  {time.RFC3339Nano, time.RFC3339},
	"datetime": {time.DateTime, "01/02/2006 15:04:05"},
	"date":     {time.DateOnly, "01/02/2006"}, //MM/dd/yyyy
}
var (
	timeReg = regexp.MustCompile("[a-zA-Z]+")
	fReg    = regexp.MustCompile("^[fF]+$")
	zReg    = regexp.MustCompile("^z+$")
	zZReg   = regexp.MustCompile("^Z+$")
)

func replaceToGoTimeTempl(l string) string {

	l = timeReg.ReplaceAllStringFunc(l, func(s string) string {
		// https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/custom-date-and-time-format-strings
		switch s {
		case "yyyy":
			return "2006"
		case "yy":
			return "06"
		case "M":
			return "1"
		case "MM":
			return "01"
		case "dd":
			return "02"
		case "h":
			return "3"
		case "hh":
			return "03"
		case "H": // 不支持
			return "15"
		case "HH":
			return "15"
		case "m":
			return "4"
		case "mm":
			return "04"
		case "s":
			return "5"
		case "ss":
			return "05"
		default:
			s = fReg.ReplaceAllStringFunc(s, func(s string) string {
				switch len(s) {
				case 3, 6, 9:
					s = strings.ReplaceAll(s, "f", "9")
					s = strings.ReplaceAll(s, "F", "0")
				}
				return s
			})
			s = zReg.ReplaceAllStringFunc(s, func(s string) string {
				switch len(s) {
				case 1:
					return "-07"
				case 2:
					return "-07:00"
				case 3:
					return "-07:00:00"
				case 4:
					return "-0700"
				case 5:
					return "-070000"
				}
				return s
			})
			s = zZReg.ReplaceAllStringFunc(s, func(s string) string {
				switch len(s[1:]) {
				case 1:
					return "Z07"
				case 2:
					return "Z07:00"
				case 3:
					return "Z07:00:00"
				case 4:
					return "Z0700"
				case 5:
					return "Z070000"
				}
				return s
			})
		}
		return s
	})
	return l
}

func ParseTimeLayouts(a []string) []string {
	if len(a) == 0 {
		return nil
	}
	var layouts []string
	for _, l := range a {
		if p := timeLayoutMap[strings.ToLower(l)]; len(p) > 0 {
			layouts = append(layouts, p...)
		} else {
			// 2006-01-02T15:04:05.999999999Z07:00
			l = strings.ReplaceAll(l, "%yyyy", "2006")
			l = strings.ReplaceAll(l, "%yyyy", "2006")
			r := regexp.MustCompile("%[a-zA-Z]+")
			r.ReplaceAllStringFunc(l, func(s string) string {
				switch s[1:] {
				case "yyyy":
					return "2006"
				case "yy":
					return "06"
				case "MM":
					return "01"
				case "dd":
					return "02"
				case "hh":
					return "15"
				case "mm":
					return "04"
				case "ss":
					return "05"
				case "f":
					return "999999999"
				}
				return s
			})
			layouts = append(layouts, l)
			// todo 转义
			// https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/standard-date-and-time-format-strings
			// https://docs.python.org/zh-cn/3/library/time.html
			// https://www.elastic.co/guide/en/beats/filebeat/current/processor-timestamp.html
		}
	}
	return layouts
}

func ParseTime(s string, layouts []string, loc *time.Location) (time.Time, error) {
	if len(layouts) == 0 {
		for _, v := range timeLayoutMap {
			for _, l := range v {
				v, err := time.ParseInLocation(l, s, loc)
				if err != nil {
					continue
				}
				return v, nil
			}
		}
	}
	for _, l := range layouts {
		v, err := time.ParseInLocation(l, s, loc)
		if err != nil {
			continue
		}
		return v, nil
	}
	// 转换失败
	return time.Time{}, fmt.Errorf("unable convert to time: %q", s)
}
