package util

import (
	"math"
	"math/big"
)

// Div div
func Div(x, y float64) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Quo(a, b)
	d, _ := c.Float64()
	return d
}

// Mul mul
func Mul(x, y float64) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Mul(a, b)
	d, _ := c.Float64()
	return d
}

// DivWithRound div with round
func DivWithRound(x, y float64, places int) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Quo(a, b)
	d, _ := c.Float64()
	return Round(d, places)
}

// MulWithRound mul with round
func MulWithRound(x, y float64, places int) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Mul(a, b)
	d, _ := c.Float64()
	return Round(d, places)
}

// Round round
func Round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
