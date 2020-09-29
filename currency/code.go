package currency

import (
	"regexp"
)

// CurrencyCode [ISO 4217 specification](https://en.wikipedia.org/wiki/ISO_4217)所规定的3位大写字母货币代码，如CNY。
type CurrencyCode string

var codeMatch = regexp.MustCompile(`^[A-Z]{3}$`).MatchString

func (c CurrencyCode) Invalid() bool {
	return !codeMatch(string(c))
}

func (c CurrencyCode) String() string {
	if c.Invalid() {
		return "<Invalid CurrencyCode>" + string(c)
	}
	return string(c)
}
