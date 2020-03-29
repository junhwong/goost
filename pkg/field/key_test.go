package field

import "testing"

func TestIsValidKeyName(t *testing.T) {
	testCases := []struct {
		name string
		pass bool
	}{
		{
			name: "foo",
			pass: true,
		},
		{
			name: "foo.bar0",
			pass: true,
		},
		{
			name: "foo..bar0",
			pass: false,
		},
		{
			name: "fo_o.b-ar",
			pass: true,
		},
		{
			name: "fo_o.b-ar.oth",
			pass: true,
		},
		{
			name: "foo.123",
			pass: false,
		},
		{
			name: "123.bar",
			pass: false,
		},
		{
			name: ".",
			pass: false,
		},
		{
			name: "-._",
			pass: false,
		},
		{
			name: "*dfd5#",
			pass: false,
		},
		{
			name: "中文",
			pass: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if IsValidKeyName(tC.name) != tC.pass {
				t.Fail()
			}
		})
	}
}

func TestInt(t *testing.T) {
	// 922_337_203_685 = 9.22千亿, 央行万亿以上
	// 验证
	t.Logf("%d", 9223372036854775800/10000000)
}
