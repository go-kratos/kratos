package model

import (
	"math"
	"strconv"
)

// CalNum2SufStr int64转换带后缀字符串（K,W,E）
func CalNum2SufStr(n int64) string {
	var f float64
	var s string
	if n > 1000 {
		f = float64(n) / 1000
		s = strconv.FormatFloat(math.Ceil(f), 'f', 0, 64) + "k"
	} else {
		f = float64(n)
		s = strconv.FormatFloat(f, 'f', 0, 64)
	}
	return s
}
