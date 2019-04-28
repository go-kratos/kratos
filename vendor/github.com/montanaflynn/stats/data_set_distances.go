package stats

import (
	"math"
)

// Validate data for distance calculation
func validateData(dataPointX, dataPointY []float64) error {
	if len(dataPointX) == 0 || len(dataPointY) == 0 {
		return EmptyInput
	}

	if len(dataPointX) != len(dataPointY) {
		return SizeErr
	}
	return nil
}

// Computes Chebyshev distance between two data sets
func ChebyshevDistance(dataPointX, dataPointY []float64) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	var tempDistance float64
	for i := 0; i < len(dataPointY); i++ {
		tempDistance = math.Abs(dataPointX[i] - dataPointY[i])
		if distance < tempDistance {
			distance = tempDistance
		}
	}
	return distance, nil
}

//
// Computes Euclidean distance between two data sets
//
func EuclideanDistance(dataPointX, dataPointY []float64) (distance float64, err error) {

	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	distance = 0
	for i := 0; i < len(dataPointX); i++ {
		distance = distance + ((dataPointX[i] - dataPointY[i]) * (dataPointX[i] - dataPointY[i]))
	}
	return math.Sqrt(distance), nil
}

//
// Computes Manhattan distance between two data sets
//
func ManhattanDistance(dataPointX, dataPointY []float64) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	distance = 0
	for i := 0; i < len(dataPointX); i++ {
		distance = distance + math.Abs(dataPointX[i]-dataPointY[i])
	}
	return distance, nil
}

//
// Computes minkowski distance between two data sets.
//
// Input:
//    dataPointX: First set of data points
//    dataPointY: Second set of data points. Length of both data
//                sets must be equal.
//    lambda:     aka p or city blocks; With lambda = 1
//                returned distance is manhattan distance and
//                lambda = 2; it is euclidean distance. Lambda
//                reaching to infinite - distance would be chebysev
//                distance.
// Output:
//     Distance or error
//
func MinkowskiDistance(dataPointX, dataPointY []float64, lambda float64) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	for i := 0; i < len(dataPointY); i++ {
		distance = distance + math.Pow(math.Abs(dataPointX[i]-dataPointY[i]), lambda)
	}
	distance = math.Pow(distance, float64(1/lambda))
	if math.IsInf(distance, 1) == true {
		return math.NaN(), InfValue
	}
	return distance, nil
}
