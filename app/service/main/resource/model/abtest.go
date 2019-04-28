package model

import "hash/crc32"

// AbTest struct
type AbTest struct {
	ID          int64  `json:"groupId"`
	Name        string `json:"groupName"`
	Threshold   int64  `json:"flowPercent"`
	ParamValues string `json:"-"`
	UTime       int64  `json:"-"`
}

// AbTestIn check build in test
func (ab *AbTest) AbTestIn(buvid string) (ok bool) {
	ration := crc32.ChecksumIEEE([]byte(buvid)) % 100
	if ration < uint32(ab.Threshold) {
		ok = true
	}
	return
}
