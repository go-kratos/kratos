package model

const (
	// ExpMulti exp multi
	ExpMulti = 100
	// level floor conf.
	level1   = 1
	level2   = 200
	level3   = 1500
	level4   = 4500
	level5   = 10800
	level6   = 28800
	levelMax = -1
)

// BuildLevel build level by LevelInfo
func (lv *LevelInfo) BuildLevel(exp int64, sexp bool) {
	exp = exp / ExpMulti
	switch {
	case exp < level1:
		lv.Cur = 0
		lv.Min = 0
		lv.NextExp = level1
	case exp < level2:
		lv.Cur = 1
		lv.Min = level1
		lv.NextExp = level2
	case exp < level3:
		lv.Cur = 2
		lv.Min = level2
		lv.NextExp = level3
	case exp < level4:
		lv.Cur = 3
		lv.Min = level3
		lv.NextExp = level4
	case exp < level5:
		lv.Cur = 4
		lv.Min = level4
		lv.NextExp = level5
	case exp < level6:
		lv.Cur = 5
		lv.Min = level5
		lv.NextExp = level6
	default:
		lv.Cur = 6
		lv.Min = level6
		lv.NextExp = levelMax
	}
	if sexp {
		lv.NowExp = int32(exp)
	}
}
