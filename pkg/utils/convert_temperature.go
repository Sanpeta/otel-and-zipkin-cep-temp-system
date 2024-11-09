package utils

import "math"

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func RoundToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}
