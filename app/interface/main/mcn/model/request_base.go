package model

import (
	"math"
)

const (
// dateFmt     = "20060102"
// dateTimeFmt = "20060102_150405"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

//ToFixed fix float precision
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output // since go 1.9 doesn't have a math.Round function...
}

// floatFormat format float to string
// func floatFormat(f float64) string {
// 	return strconv.FormatFloat(f, 'f', 2, 64)
// }

// intFormat format int to string
// func intFormat(i int64) string {
// 	return strconv.Itoa(int(i))
// }

//PageArg page arg
type PageArg struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

//PageResult page result
type PageResult struct {
	Page       int `json:"page"`
	TotalCount int `json:"total_count"`
}

//CheckPageValidation check the page validte, return limit offset
func (arg *PageArg) CheckPageValidation() (limit, offset int) {
	if arg.Page < 1 {
		arg.Page = 1
	}
	if arg.Size > 1000 || arg.Size <= 0 {
		arg.Size = 10
	}
	limit = arg.Size
	offset = (arg.Page - 1) * limit
	return
}

//ToPageResult cast to page result
func (arg *PageArg) ToPageResult(total int) (res PageResult) {
	res.TotalCount = total
	res.Page = arg.Page
	return
}

//ExportArg export arg
type ExportArg struct {
	Export string `form:"export"`
}

//ExportFormat export format
func (e *ExportArg) ExportFormat() string {
	return e.Export
}
