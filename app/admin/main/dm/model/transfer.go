package model

import "go-common/library/time"

// all variable used in dm transfer
const (
	TransferJobStateAll      = int8(-1)
	TransferJobStatInit      = int8(0)
	TransferJobStatFinished  = int8(1)
	TransferJobStatFailed    = int8(2)
	TransferJobStatTransfing = int8(3)
)

// TransList transfer list info
type TransList struct {
	ID    int64     `json:"id"`    //弹幕转移ID
	From  int64     `json:"from"`  //来源Cid
	To    int64     `json:"to"`    //目标Cid
	State int64     `json:"state"` //弹幕转移状态
	Title string    `json:"title"` //来源稿件标题
	Ctime time.Time `json:"ctime"` //转移开始时间
}

// TransListRes return transfer list and page info
type TransListRes struct {
	Result []*TransList `json:"result"`
	Page   *PageInfo    `json:"page"`
}

// TransferJobInfo dm transfer info
type TransferJobInfo struct {
	ID      int64
	FromCID int64
	ToCID   int64
	MID     int64
	Offset  float64
	State   int8
	Ctime   time.Time
	Mtime   time.Time
}
