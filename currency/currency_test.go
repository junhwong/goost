package currency

import (
	"encoding/json"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestParse(t *testing.T) {
	tests := []struct {
		in        string
		expect    Money
		expectErr bool
		cut, max  int
	}{
		{
			in:        "CNY12.1201",
			expect:    121201,
			expectErr: false,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12.1200",
			expect:    121200,
			expectErr: false,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12..1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "XYZ12.1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        ".1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        " CNY12.1200 ",
			expectErr: false,
			expect:    121200,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12.1200A",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12. 1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNYY12.1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12Y.1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "CNY12 .1200",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "XYZ",
			expectErr: true,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "1",
			expect:    10000,
			expectErr: false,
			cut:       -1,
			max:       -1,
		},
		{
			in:        "0.1",
			expect:    1000,
			expectErr: false,
			cut:       -1,
			max:       1,
		},
		{
			in:        "TBD0.12",
			expect:    1000,
			expectErr: false,
			cut:       1,
			max:       -1,
		},
		{
			in:        "-0.1",
			expect:    -1000,
			expectErr: false,
			cut:       -1,
			max:       -1,
		},
	}
	// TODO: 需要更多测试用例
	_, _ = Define("TBD", 2)
	for _, it := range tests {
		v, _, err := Parse(it.in, it.cut, it.max)
		if err != nil {
			if it.expectErr {
				continue
			}
			t.Errorf("parse %q, unexpect error: %v", it.in, err)
		}
		if it.expectErr {
			t.Errorf("parse %q, expect error, got: nil", it.in)
			continue
		}
		if v != it.expect {
			t.Errorf("parse %q, expect %v, got: %v", it.in, int(it.expect), v.RawInt())
		}
	}

}

func runDefineSync() {
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := rand.Intn(1000)
			time.Sleep(time.Nanosecond * time.Duration(r))
			Define("TBD", 4)
			r = rand.Intn(1000)
			time.Sleep(time.Nanosecond * time.Duration(r))
			Define("TBD", 4)
			Currency("ENY")

		}()
	}
	wg.Wait()
}

func TestDefineSync(t *testing.T) {
	runDefineSync()
}

func BenchmarkDefineSync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runDefineSync()
	}
}

func TestExchange(t *testing.T) {
	// 大VS小
	q, _ := Define("JPY", FromFloat(15.4527))
	m := Money(20000) //CNY

	e := m.ExchangeTo(q) // 转换为JPY
	if q.Format(e) != "JPY30.9054" {
		t.Fatal()
	}
	e = e.ExchangeFrom(q) // 转换为CNY
	if e != m {
		t.Fatal()
	}
	if _, err := m.ExchangeToE(nil); err == nil {
		t.Fatal()
	}
	if _, err := m.ExchangeFromE(nil); err == nil {
		t.Fatal()
	}
	// 小VS大
	q, _ = Define("USD", FromFloat(0.1465))
	m = Money(20000) //CNY

	e = m.ExchangeTo(q) // 转换为USD
	if q.Format(e) != "USD0.2930" {
		t.Fatal()
	}
	e = e.ExchangeFrom(q) // 转换为CNY
	if e != m {
		t.Fatal()
	}
}

func TestPrecent(t *testing.T) {
	//CNY0.9888 1*0.9
	m := Money(10000)
	if baseCurrency.Format(m.Precent(-120)) != "CNY-1.2000" {
		t.Fatal()
	}
	if baseCurrency.Format(m.Precent(90)) != "CNY0.9000" {
		t.Fatal()
	}
}

func TestFormat(t *testing.T) {
	testCases := []struct {
		in     Money
		prec   uint8
		expect string
	}{
		{
			in:     FromFloat(2),
			prec:   4,
			expect: "2.0000",
		},
		{
			in:     FromFloat(2),
			prec:   2,
			expect: "2.00",
		},
		{
			in:     FromFloat(2),
			prec:   0,
			expect: "2",
		},
		{
			in:     FromFloat(-2),
			prec:   0,
			expect: "-2",
		},
	}
	for _, tC := range testCases {
		s := tC.in.Format(BaseCurrencyCode, tC.prec, true)
		if s == tC.expect {
			continue
		}
		t.Errorf("formating %s: expect %q, got: %q", tC.in, tC.expect, s)
	}
}

func TestMarshalJSON(t *testing.T) {
	obj := struct {
		M Money
	}{}
	s := `{"M":"CNY1.0000"}`
	err := json.Unmarshal([]byte(s), &obj)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(&obj)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != s {
		t.Fatal(string(b))
	}

	// nolint
	Define("USD", FromFloat(0.1465))
	s = `{"M":"USD1.0000"}`
	err = json.Unmarshal([]byte(s), &obj)
	if err != nil {
		t.Fatal(err)
	}

	if obj.M.String() != "CNY6.8259" {
		t.Fatal(obj.M.String())
	}
}
func TestInt(t *testing.T) {
	// 922_337_203_685 = 9.22千亿, 央行万亿以上
	// 验证
	t.Logf("%d", 9223372036854775800/10000000)
}
