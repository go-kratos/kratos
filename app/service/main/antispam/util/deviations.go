package util

import "math"

// Max .
func Max(vars []int64) (maxVar int64) {
	for _, i := range vars {
		if i > maxVar {
			maxVar = i
		}
	}
	return
}

// Expectation .
func Expectation(randomVars []float64) float64 {
	if len(randomVars) == 0 {
		return 0
	}

	var sum float64
	for _, rv := range randomVars {
		sum += rv
	}

	return sum / float64(len(randomVars))
}

// StdDeviation .
func StdDeviation(randomVars []float64) float64 {
	if len(randomVars) == 0 {
		return 0
	}
	return math.Sqrt(Deviation(randomVars))
}

// Deviation .
func Deviation(randomVars []float64) float64 {
	if len(randomVars) == 0 {
		return 0
	}
	var total float64
	expec := Expectation(randomVars)
	for _, rv := range randomVars {
		total += math.Pow(rv-expec, 2.0)
	}
	return total / float64(len(randomVars))
}

// Normallization .
func Normallization(randomVars []int64) []float64 {
	if len(randomVars) == 0 {
		return nil
	}
	maxVal := Max(randomVars)
	if maxVal == 0 || maxVal == 1 {
		return nil
	}
	res := make([]float64, 0, len(randomVars))
	for _, rv := range randomVars {
		res = append(res, math.Log10(float64(rv))/math.Log10(float64(maxVal)))
	}
	return res
}
