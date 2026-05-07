package util

import "github.com/shopspring/decimal"

func SafeDivision[T int | int64 | int32](molecule, denominator T) T {
	if denominator == 0 {
		return 0
	}
	return molecule / denominator
}

func SafeFloatDivision[T int | int64 | int32 | float64](molecule, denominator T) float64 {
	if denominator == 0 {
		return 0
	}
	v := decimal.NewFromFloat(float64(molecule) / float64(denominator))
	vv, _ := v.Round(2).Float64()
	return vv
}
