package utils

import "math"

func ToPtr[T any](t T) *T {
	return &t
}

func RoundTo(n float64, decimals uint32) float64 {
	const pow = 10

	decimalPlaces := math.Pow(pow, float64(decimals))

	return math.Round(n*decimalPlaces) / decimalPlaces
}
