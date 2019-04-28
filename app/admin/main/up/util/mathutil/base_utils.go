package mathutil

//EPSILON very small number
var EPSILON float32 = 0.00000001

//FloatEquals float equals
func FloatEquals(a, b float32) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

//EPSILON64 very small number
var EPSILON64 = 0.00000001

//Float64Equals float equal
func Float64Equals(a, b float64) bool {
	if (a-b) < EPSILON64 && (b-a) < EPSILON64 {
		return true
	}
	return false
}

//Min min
func Min(a, b int) int {
	if a < b {
		return a
	}
	if b < a {
		return b
	}
	return a
}

//Max max
func Max(a, b int) int {
	if a > b {
		return a
	}
	if b > a {
		return b
	}
	return a
}

// IntersectInt64 intersect map
func IntersectInt64(as, bs []int64) []int64 {
	i := make([]int64, 0, Max(len(as), len(bs)))
	for _, a := range as {
		for _, b := range bs {
			if a == b {
				i = append(i, a)
			}
		}
	}
	return i
}
