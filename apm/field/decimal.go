package field

import (
	"math"
	"strconv"
	"strings"
)

const (
	maxSafeInteger = 9007199254740992.0  // +2^53
	minSafeInteger = -9007199254740992.0 // -2^53
)

func hasDecimal(f float64) bool {

	// 处理特殊值
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}

	// 安全范围内使用快速判断
	if f > minSafeInteger && f < maxSafeInteger {
		return f != math.Trunc(f)
	}

	if f == minSafeInteger || f == maxSafeInteger {
		return false
	}

	s := strconv.FormatFloat(f, 'g', -1, 64)
	return strings.Contains(s, ".")
}
