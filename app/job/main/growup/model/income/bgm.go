package income

//import (
//	"go-common/library/time"
//)

// BGM background music
type BGM struct {
	ID     int64  `json:"id"`
	AID    int64  `json:"aid"`
	CID    int64  `json:"cid"`
	SID    int64  `json:"sid"`
	MID    int64  `json:"uid"`
	JoinAt string `json:"join_time"`
	Title  string `json:"title"`
}
