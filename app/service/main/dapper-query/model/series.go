package model

// SeriesItem .
type SeriesItem struct {
	Field string
	Rows  []*float64
}

// Series data
type Series struct {
	Timestamps []int64
	Items      []*SeriesItem
}
