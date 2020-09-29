package currency

import (
	"errors"
	"strconv"
	"strings"
)

// Money 表示一个整数化的货币值。系统代表本位币，始终是1。
type Money int64

func (m Money) RawInt() int64 {
	return int64(m)
}

func (m Money) RawFloat() float64 {
	return float64(m)
}

func (m Money) Float() float64 {
	return div(float64(m), Precision)
}

func (m Money) String() string {
	return m.Format(BaseCurrencyCode, Precision)
}

// Exchange 使用报价汇率转换
func (m Money) ExchangeToE(quote *CurrencyPair) (Money, error) {
	if quote.Invalid() {
		return 0, errors.New("Invalid quote currency pair")
	}
	return Money(div(m.RawFloat()*quote.Rate().RawFloat(), Precision)), nil
}

// Exchange 汇率转换
func (m Money) ExchangeTo(quote *CurrencyPair) (quoteMoney Money) {
	quoteMoney, _ = m.ExchangeToE(quote)
	return
}

// Exchange 汇率转换
func (m Money) ExchangeFromE(quote *CurrencyPair) (Money, error) {
	if quote.Invalid() {
		return 0, errors.New("Invalid quote currency pair")
	}
	// m/ q.Rate() * Money(10000)
	return Money(mut(m.RawFloat()/quote.Rate().RawFloat(), Precision)), nil
}

// Exchange 汇率转换
func (m Money) ExchangeFrom(quote *CurrencyPair) (quoteMoney Money) {
	quoteMoney, _ = m.ExchangeFromE(quote)
	return
}

// Precent 计算百分比
func (m Money) Precent(p int) (quoteMoney Money) {
	prec := float64(p) / float64(100)
	return Money(float2Float(m.RawFloat() * prec))
}

func (m Money) Format(c CurrencyCode, prec uint8, withoutCode ...bool) string {
	if prec > Precision {
		prec = Precision
	}
	exinc := false
	for _, it := range withoutCode {
		exinc = it
	}
	s := strconv.FormatFloat(m.Float(), 'f', int(prec)+1, 64)
	arr := strings.SplitN(s, ".", 2)
	rs := c.String()
	if exinc {
		rs = ""
	}
	rs += arr[0]
	if prec > 0 {
		rs += "." + arr[1][:prec]
	}
	return rs
}

//MarshalJSON 实现它的json序列化方法
func (r Money) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

//UnmarshalJSON 实现它的json序列化方法
func (r *Money) UnmarshalJSON(b []byte) error {
	v, p, err := Parse(strings.Trim(string(b), `"`), -1, -1)
	if err != nil {
		return err
	}
	if p.Code() != BaseCurrencyCode {
		v, err = v.ExchangeFromE(p)
		if err != nil {
			return err
		}
	}
	*r = v
	return nil
}

func FromFloat(raw float64) Money {
	return Money(mut(raw, Precision))
}

func FromInt(raw int64) Money {
	return Money(mut(float64(raw), Precision))
}
