package field

import (
	"strconv"
	"strings"
)

type RoundMode string

const (
	RoundHalfUp   RoundMode = "ROUND_HALF_UP"   // 四舍五入
	RoundHalfDown RoundMode = "ROUND_HALF_DOWN" // 五舍六入
	RoundHalfEven RoundMode = "ROUND_HALF_EVEN" // default 四舍六入五取偶法 银行家舍入法，其实质是一种四舍六入五取偶（又称四舍六入五留双）法
	RoundTrunc    RoundMode = "ROUND_TRUNC"     // 截取指定位数小数
	RoundDown     RoundMode = "ROUND_DOWN"      // 去掉小数部分取整
	RoundUp       RoundMode = "ROUND_UP"        // 取右边最近的整数
)

// 舍入.
//
// FIXME: 目前是使用字符串处理,效率不高
func Round(v float64, prec int, m RoundMode) float64 {
	switch m {
	case RoundHalfEven:
		s := strconv.FormatFloat(v, 'f', prec, 64)
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return v
	case RoundUp, RoundDown:
		prec = 0
	}

	parts := strings.SplitN(strconv.FormatFloat(v, 'f', -1, 64), ".", 2)
	if len(parts) == 1 {
		return v
	}

	if b := len(parts[1]) <= prec; b || m == RoundTrunc {
		if b {
			return v
		}
		v, err := strconv.ParseFloat(parts[0]+"."+parts[1][:prec], 64)
		if err != nil {
			panic(err)
		}
		return v
	}

	afr := parts[1][prec:]
	parts[1] = parts[1][:prec]
	var carry func()
	out := false
	carry = func() {
		l := len(parts[1]) - 1
		if l < 0 {
			f, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				panic(err)
			}
			if v < 0 {
				v = f - 1
			} else {
				v = f + 1
			}

			out = true
			return
		}
		c := parts[1][l]
		if c == '9' {
			parts[1] = parts[1][:l]
			carry()
			return
		}
		i, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		parts[1] = parts[1][:l] + strconv.Itoa(i+1)
	}

	switch m {
	case RoundHalfUp:
		if canCarray(afr, '5', '4') {
			carry()
		}
	case RoundHalfDown:
		if canCarray(afr, '6', '5') {
			carry()
		}
	case RoundDown:
		parts[1] = "0"
	case RoundUp:
		parts[1] = ""
		carry()
	default:
		panic("unknown rounding mode: " + m)
	}

	if out {
		return v
	}
	v, err := strconv.ParseFloat(parts[0]+"."+parts[1], 64)
	if err != nil {
		panic(err)
	}
	return v
}

var canCarray func(s string, a, b byte) bool

func init() {
	canCarray = func(s string, a, b byte) bool {
		if len(s) == 0 {
			return false
		}
		if s[0] >= a {
			return true
		}
		if len(s) < 2 || s[1] < b {
			return false
		}
		return canCarray(s[1:], a, b)
	}
}
