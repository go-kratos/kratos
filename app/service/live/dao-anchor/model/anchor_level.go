package model

import (
	"errors"
	"math"
)

const MaxAnchorLevel = 40

var anchorLevelScoreTable = []int64{
	0,
	50,
	200,
	470,
	920,
	2100,
	4060,
	7160,
	11760,
	18060,
	27160,
	39610,
	56410,
	78810,
	109810,
	154810,
	226810,
	319810,
	442810,
	602810,
	816810,
	1138810,
	1594810,
	2214810,
	3004810,
	3984810,
	5229810,
	6909810,
	9013810,
	11883810,
	15613810,
	20613810,
	27313810,
	36413810,
	47813810,
	62013810,
	79513810,
	99513810,
	122013810,
	147013810,
	int64(math.MaxInt64),
}

// Returns the index of first element that **greater** than the given value.
func upperBound(list []int64, value int64) int {
	count := len(list)
	first:= 0
	for count > 0 {
		i := first
		step := count / 2
		i += step
		if value >= list[i] {
			first = i + 1
			count -= step + 1
		} else {
			count = step
		}
	}

	return first
}

// GetAnchorLevel returns anchor level, i.e. Lv1 ~ Lv40, corresponding to the given score.
func GetAnchorLevel(score int64) (int64, error) {
	if score < 0 {
		return 0, errors.New("invalid anchor score")
	}

	return int64(upperBound(anchorLevelScoreTable, score)), nil
}

// GetLevelScoreInfo returns left & right score of a given level.
func GetLevelScoreInfo(lv int64) (left, right int64, err error) {
	if lv < 1 || lv >= int64(len(anchorLevelScoreTable)) {
		return 0, 0, errors.New("invalid request level")
	}

	left = anchorLevelScoreTable[lv-1]
	right = anchorLevelScoreTable[lv] - 1

	return
}

// GetAnchorLevelColor returns level color.
func GetAnchorLevelColor(lv int64) (int64, error) {
	if lv < 1 || lv > MaxAnchorLevel {
		return 0, errors.New("invalid request level")
	}

	q := lv / 10
	r := lv % 10
	if r == 0 {
		q -= 1
	}

	return q, nil
}