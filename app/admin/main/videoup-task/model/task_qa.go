package model

import "time"

//qa task state & type & log business/type.
const (
	QAStateWait   = int16(1)
	QAStateFinish = int16(2)

	QATypeVideo = int8(1)

	LogQATask      = 111
	LogQATaskVideo = 1
)

//QAStates states.
var (
	QAStates = map[int16]string{
		QAStateWait:   "待质检",
		QAStateFinish: "已质检",
	}
)

//QATask qatask
type QATask struct {
	ID       int64     `json:"id"`
	State    int16     `json:"state"`
	Type     int8      `json:"type"`
	DetailID int64     `json:"detail_id"`
	UID      int64     `json:"uid"`
	FTime    time.Time `json:"ftime"`
	CTime    time.Time `json:"ctime"`
}
