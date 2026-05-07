package util

import "github.com/shopspring/decimal"

// 保留1位小数
func Decimal(value float64) float64 {
	v := decimal.NewFromFloat(value)
	vv, _ := v.Round(1).Float64()
	return vv
}
