package currency

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reg = regexp.MustCompile(`^(?P<c>[A-Z]{3})?(?P<i>\-?\d+)(?P<d>\.\d+)?$`)

// Parse 解析货币字符串。
// cutDecimalLen 参数小于0将忽略，否则直接裁剪实际小数长度到指定值之内(忽略后缀0)。
// maxDecimalLen 参数小于0将忽略，否则判断实际小数长度(裁剪之后，如果cutDecimalLen生效)是否<=指定值，否则返回错误(忽略后缀0)。
func Parse(s string, cutDecimalLen, maxDecimalLen int) (Money, *CurrencyPair, error) {
	values := reg.FindStringSubmatch(strings.TrimSpace(s))
	if len(values) < 2 {
		return 0, nil, fmt.Errorf("parsing %q: invalid currency syntax", s)
	}
	names := reg.SubexpNames()

	def := baseCurrency
	vs := ""
	var err error
	for i, name := range names {
		if values[i] == "" {
			continue
		}
		switch name {
		case "c":
			def = Currency(values[i])
			if def == nil || def.Invalid() {
				return 0, nil, fmt.Errorf("parsing %q: invalid currency code: %s", s, values[i])
			}
		case "i":
			vs += values[i]
		case "d":
			prec := strings.TrimRight(values[i], "0")
			if cutDecimalLen > -1 && len(prec)-1 > cutDecimalLen {
				prec = prec[0 : cutDecimalLen+1]
			}
			prec = strings.TrimRight(prec, "0")
			if maxDecimalLen > -1 && len(prec)-1 > maxDecimalLen {
				return 0, nil, fmt.Errorf("parsing %q: precision must be ≤%d, actual:%d", s, maxDecimalLen, len(prec)-1)
			}
			if prec != "." {
				vs += prec
			}
		}
	}

	f, err := strconv.ParseFloat(vs, 64)
	if err != nil {
		return 0, nil, err
	}
	m := Money(mut(f, Precision))
	return m, def, err
}
