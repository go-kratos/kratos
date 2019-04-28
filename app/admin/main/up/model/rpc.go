package model

import "time"

//ArgSpecial arg
type ArgSpecial struct {
	GroupID int64
}

//ArgInfo arg
type ArgInfo struct {
	Mid  int64
	From int
}

//ArgMidWithDate arg
type ArgMidWithDate struct {
	Mid  int64
	Date time.Time // 需要查询的时间，如果不填，会使用当前最新的数据
}
