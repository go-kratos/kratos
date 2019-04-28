package model

import "time"

//ArgSpecial special arg
type ArgSpecial struct {
	GroupID int64
}

//ArgInfo info arg
type ArgInfo struct {
	Mid  int64
	From int
}

//ArgMidWithDate arg mid with date
type ArgMidWithDate struct {
	Mid  int64
	Date time.Time // 需要查询的时间，如果不填，会使用当前最新的数据
}

//ArgUpSwitch info arg
type ArgUpSwitch struct {
	Mid   int64
	State int
	From  int
	IP    string
}
