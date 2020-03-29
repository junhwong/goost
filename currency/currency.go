package currency

// Currency 货币。系统代表本位币，始终是1，如人民币 1元=10000毫。4位小数 0.0001
//
// bug() 比特币的最小单位是一个中本聪，等于0.00000001个比特币。
type Currency int64

const (
	currencyBase  Currency = 10000
	currencyBasef float64  = 10000

	// CurrencyMaxValue 货币的最大值， 90.22兆
	CurrencyMaxValue Currency = 9223372036854775800
	// CurrencyMinValue 货币的最小值
	CurrencyMinValue Currency = -9223372036854775800
)

// BaseCurrencyCode 本位币代码
var BaseCurrencyCode string = "CNY"

// CurrencyPair 货币对。 https://www.investopedia.com/terms/c/currencypair.asp
type CurrencyPair struct {
	// 本位币 ISO Currency Code
	BaseCode string
	// 报价货币 ISO Currency Code
	QuoteCode string
	// 汇率, *10000
	ExchangeRate int32
}

// 转换为本位币
func (cp CurrencyPair) Exchange(quote Currency) Currency {
	return 0
}

func ParseE(fullValue interface{}) (Currency, error) {
	return 0, nil
	//   f,err:= cast.float64E(fullValue)
	//   if err!=nil{
	//     return 0, err
	//   }
	//   return Currency(f * currencyBasef)
}

func Parse(fullValue interface{}) Currency {
	v, _ := ParseE(fullValue)
	return v
}

// 转换为交易金额字符串，2位小数。
func (c Currency) TradeString() string {
	return ""
}

// 4位小数金额。
func (c Currency) String() string {
	return ""
}
