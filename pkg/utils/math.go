package utils

import "math"

func RoundTo(f float64, n int) float64 {
	return math.Round(f*math.Pow(10, float64(n))) / math.Pow(10, float64(n))
}
