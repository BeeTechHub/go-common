package utils

import "math"

func RoundFloat(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}
