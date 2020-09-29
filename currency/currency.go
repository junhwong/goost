package currency

import (
	"errors"
	"strings"
	"sync/atomic"
)

const (
	// MoneyMaxValue 货币的最大值， 90.22兆
	MoneyMaxValue Money = 9223372036854775800
	// MoneyMinValue 货币的最小值
	MoneyMinValue Money = -9223372036854775800
	// Precision 用于存储和转换货币相对于系统的精度，也用于format。
	// 最大值是8。
	// 如人民币: 1元=10000毫。4位小数为 0.0001元，即 1Currency=0.0001元
	Precision uint8 = 4
	// BaseCurrencyCode 表示本位币代码。
	BaseCurrencyCode CurrencyCode = "CNY"
)

func init() {
	// 如果编译改变该值将会检测到它的准确性
	if BaseCurrencyCode.Invalid() {
		panic("currency: Invalid BaseCurrencyCode: " + string(BaseCurrencyCode))
	}
	defines.Store(map[CurrencyCode]*CurrencyPair{})
}

// // CurrencyPair 表示一个相对于本位币的货币对。
// type CurrencyPair interface {
// 	Code() CurrencyCode
// 	Rate() Money
// 	Format(money Money, withoutCode ...bool) string
// 	Invalid() bool
// }

// CurrencyPair 表示一个相对于本位币的货币对。
type CurrencyPair struct {
	code CurrencyCode
	rate Money
}

func (p *CurrencyPair) Code() CurrencyCode {
	return p.code
}
func (p *CurrencyPair) Rate() Money {
	return p.rate
}
func (p *CurrencyPair) Invalid() bool {
	if p == nil || p.rate < 1 {
		return true
	}
	return p.code.Invalid()
}
func (p *CurrencyPair) Format(m Money, withoutCode ...bool) string {
	code := CurrencyCode("")
	if p != nil {
		code = p.code
	}
	return m.Format(code, Precision, withoutCode...)
}

var (
	baseCurrency = &CurrencyPair{BaseCurrencyCode, 1}
	defines      = atomic.Value{}
)

func Define(code string, exchangeRate Money) (*CurrencyPair, error) {
	c := CurrencyCode(strings.ToUpper(strings.TrimSpace(code)))
	if c.Invalid() {
		return nil, errors.New("Invalid currency code: " + code)
	}
	if c == BaseCurrencyCode {
		return nil, errors.New("Cannot set base currency: " + code)
	}
	dirty := map[CurrencyCode]*CurrencyPair{}
	if p := defines.Load().(map[CurrencyCode]*CurrencyPair); p != nil {
		for k, v := range p {
			dirty[k] = v
		}
	}
	def := &CurrencyPair{
		code: c,
		rate: exchangeRate,
	}
	dirty[c] = def
	defines.Store(dirty)
	return def, nil
}

func Currency(code string) *CurrencyPair {
	c := CurrencyCode(code)
	if c == BaseCurrencyCode {
		return baseCurrency
	}
	if p := defines.Load().(map[CurrencyCode]*CurrencyPair); p != nil {
		def := p[c]
		if def != nil {
			return def
		}
	}
	return nil
}
