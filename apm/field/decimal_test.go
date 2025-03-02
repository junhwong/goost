package field

import (
	"math"
	"testing"
)

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		input  float64
		expect bool
	}{
		// 安全范围内的测试
		{123.456, true},
		{123.0, false},

		// 2^53 边界测试
		{maxSafeInteger, false},
		{maxSafeInteger - 0.5, false},
		{minSafeInteger, false},
		{minSafeInteger + 0.5, false},

		// 超过安全范围的精确测试
		{math.Pow(2, 53) + 1.5, true},  // 实际存储为 9007199254740994
		{math.Pow(2, 100) + 0.5, true}, // 非常大的数但仍可检测小数
		{math.Pow(2, 100) + 1.0, true}, // 偶数间隔的整数值
	}

	for _, tt := range tests {
		got := hasDecimal(tt.input)
		if got != tt.expect {
			t.Errorf("input: %v, expect: %v, got: %v", tt.input, tt.expect, got)
		}
	}
}
