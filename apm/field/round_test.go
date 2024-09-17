package field

import "testing"

// TestRound 测试 Round 函数的不同情况
// TODO 添加更多测试用例， 目前是AI辅助生成
func TestRound(t *testing.T) {
	testCases := []struct {
		value     float64
		precision int
		mode      RoundMode
		expected  float64
	}{
		// 四舍五入 (RoundHalfUp)
		{1.44, 1, RoundHalfUp, 1.4},
		{1.45, 1, RoundHalfUp, 1.5},
		{1.45, 0, RoundHalfUp, 2.0},
		{1.55, 1, RoundHalfUp, 1.6},
		{1.55, 0, RoundHalfUp, 2.0},
		{1.5, 0, RoundHalfUp, 2.0},
		{2.5, 0, RoundHalfUp, 3.0},
		{1.555, 2, RoundHalfUp, 1.56},
		{1.555, 1, RoundHalfUp, 1.6},
		{1.555, 0, RoundHalfUp, 2.0},
		{1.455, 2, RoundHalfUp, 1.46},
		{1.455, 1, RoundHalfUp, 1.5},
		{1.455, 0, RoundHalfUp, 2.0},

		// 五舍六入 (RoundHalfDown)
		{1.45, 1, RoundHalfDown, 1.4},
		{1.55, 1, RoundHalfDown, 1.5},
		{1.55, 0, RoundHalfDown, 1.0},
		{1.5, 0, RoundHalfDown, 1.0},
		{2.5, 0, RoundHalfDown, 2.0},
		{4.5, 0, RoundHalfDown, 4.0},
		{1.555, 2, RoundHalfDown, 1.55},
		{1.555, 1, RoundHalfDown, 1.5},
		{1.555, 0, RoundHalfDown, 1.0},
		{1.455, 2, RoundHalfDown, 1.45},
		{1.455, 1, RoundHalfDown, 1.4},
		{1.455, 0, RoundHalfDown, 1.0},

		// 四舍六入五取偶法 (RoundHalfEven)
		{1.45, 0, RoundHalfEven, 1.0},
		{2.5, 0, RoundHalfEven, 2.0},
		{3.5, 0, RoundHalfEven, 4.0},
		{4.5, 0, RoundHalfEven, 4.0},
		{1.555, 2, RoundHalfEven, 1.55}, // todo 是否正确？1.56 1.55
		{1.555, 1, RoundHalfEven, 1.6},
		{1.555, 0, RoundHalfEven, 2.0},
		{2.555, 2, RoundHalfEven, 2.56},
		{2.555, 1, RoundHalfEven, 2.6},
		{2.555, 0, RoundHalfEven, 3.0},
		{1.455, 2, RoundHalfEven, 1.46},
		{1.455, 1, RoundHalfEven, 1.5},
		{1.455, 0, RoundHalfEven, 1.0},

		// 截取指定位数小数 (RoundTrunc)
		{1.45678, 2, RoundTrunc, 1.45},
		{1.45678, 3, RoundTrunc, 1.456},
		{1.45678, 4, RoundTrunc, 1.4567},
		{1.45678, 5, RoundTrunc, 1.45678},
		{-1.45678, 2, RoundTrunc, -1.45},
		{-1.45678, 3, RoundTrunc, -1.456},
		{-1.45678, 4, RoundTrunc, -1.4567},
		{-1.45678, 5, RoundTrunc, -1.45678},

		// 去掉小数部分取整 (RoundDown)
		{1.45678, 0, RoundDown, 1.0},
		{-1.45678, 0, RoundDown, -1.0},
		{1.45678, 2, RoundDown, 1.0},
		{-1.45678, 2, RoundDown, -1.0},
		{1.45678, 3, RoundDown, 1.0},
		{-1.45678, 3, RoundDown, -1.0},

		// 取右边最近的整数 (RoundUp)
		{1.45678, 0, RoundUp, 2.0},
		{-1.45678, 0, RoundUp, -2.0},
		{1.45678, 2, RoundUp, 2},
		{-1.45678, 2, RoundUp, -2},
		{1.45678, 3, RoundUp, 2},
		{-1.45678, 3, RoundUp, -2},

		// 边界情况
		{0.0, 0, RoundHalfUp, 0.0},
		{0.0, 2, RoundHalfUp, 0.0},
		{0.0, 0, RoundHalfDown, 0.0},
		{0.0, 2, RoundHalfDown, 0.0},
		{0.0, 0, RoundHalfEven, 0.0},
		{0.0, 2, RoundHalfEven, 0.0},
		{0.0, 0, RoundTrunc, 0.0},
		{0.0, 2, RoundTrunc, 0.0},
		{0.0, 0, RoundDown, 0.0},
		{0.0, 2, RoundDown, 0.0},
		{0.0, 0, RoundUp, 0.0},
		{0.0, 2, RoundUp, 0.0},

		// 负数情况
		{-1.5, 0, RoundHalfUp, -2.0},
		{-1.5, 0, RoundHalfDown, -1.0},
		{-1.5, 0, RoundHalfEven, -2.0},
		{-1.5, 0, RoundTrunc, -1.0},
		{-1.5, 0, RoundDown, -1.0},
		{-1.5, 0, RoundUp, -2.0},
	}

	for _, tc := range testCases {
		result := Round(tc.value, tc.precision, tc.mode)
		if result != tc.expected {
			t.Errorf("Round(%v, %d, %v) = %v; want %v", tc.value, tc.precision, tc.mode, result, tc.expected)
		}
	}
}
