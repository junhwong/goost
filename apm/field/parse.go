package field

import (
	"fmt"
	"regexp"
)

var (
	labelValueSepPatt = regexp.MustCompile(`\s|\,|;|\:`)
	labelValuePatt    = regexp.MustCompile(`\%?[\w\-_\.]+`)
)

// 解析标签值, 事情符合独立无特殊字符。方便作为索引使用。
// `ignoreInvalidSegment`为`true`时, 不符合的只忽略而不返回错误.
func ParseLabelValue(s string, ignoreInvalidSegment ...bool) ([]string, error) {
	var arr []string
	ignoredErr := false
	if len(ignoreInvalidSegment) > 0 {
		ignoredErr = ignoreInvalidSegment[len(ignoreInvalidSegment)-1]
	}
	for _, v := range labelValueSepPatt.Split(s, -1) {
		if len(v) == 0 {
			continue
		}
		if !labelValuePatt.MatchString(v) {
			if ignoredErr {
				continue
			}
			return nil, fmt.Errorf("invalid label value part %q of %q", v, s)
		}

		arr = append(arr, v)
	}
	return arr, nil
}
