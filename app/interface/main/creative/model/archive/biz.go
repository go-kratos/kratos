package archive

import "go-common/library/time"

// VoteType 投票业务
const VoteType int8 = 2

// BIZ business
type BIZ struct {
	ID     int64     `json:"id"`
	Type   int32     `json:"type"`
	Aid    int64     `json:"aid"`
	Uid    int64     `json:"uid"`
	State  int8      `json:"state"`
	Remark string    `json:"remark"`
	Data   string    `json:"data"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}
