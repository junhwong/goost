package currency

import (
	"fmt"
	"math"
	"strconv"
)

func float2Float(num float64) float64 {
	float_num, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", num), 64)
	return float_num
}
func mut(f float64, prec uint8) float64 {
	return float2Float(f * math.Pow10(int(prec)))
}
func div(f float64, prec uint8) float64 {
	return float2Float(f / math.Pow10(int(prec)))
}
