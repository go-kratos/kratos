package model

import (
	"sort"
)

// MeanOperationNameValue .
type MeanOperationNameValue struct {
	Tag   map[string]string
	Value float64
}

// SortRank Desc
func SortRank(values []MeanOperationNameValue) {
	sort.Slice(values, func(i, j int) bool {
		return values[i].Value > values[j].Value
	})
}
