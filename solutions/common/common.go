package common

import "math"

// GetDigits returns the digits of n between lowerTenPow and upperTenPow.
// GetDigits(5287, 1, 3) == 28
// GetDigits(5287, 2, 10) == 52
func GetDigits(n, lowerTenPow, upperTenPow int) int {
	if lowerTenPow >= upperTenPow {
		panic("Ya dingus")
	}
	return n / int(math.Pow(10, float64(lowerTenPow))) % int(math.Pow(10, float64(upperTenPow-lowerTenPow)))
}

// GetDigit is a convenience function to get a single digit.
// getDigit(5287, 0) == 7
// getDigit(5287, 2) == 2
// getDigit(5287, 4) == 0
func GetDigit(n, tenPow int) int {
	return GetDigits(n, tenPow, tenPow+1)
}
