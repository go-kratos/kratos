package upcrmmodel

import (
	"github.com/siddontang/go-mysql/mysql"
	"go-common/library/time"
)

const (
	//TaskStateStart 0
	TaskStateStart = 0
	//TaskStateFinish 1
	TaskStateFinish = 1
	//TaskStateError 2
	TaskStateError = 2
)

const (
	//TaskTypeCreditDaily 1
	TaskTypeCreditDaily = 1
	//TaskTypeScoreSectionDaily 2
	TaskTypeScoreSectionDaily = 2
	//TaskTypeSignTaskCalculate 3
	TaskTypeSignTaskCalculate = 3
	//TaskTypeSignCheckDue 4
	TaskTypeSignCheckDue = 4
)

var (
	//TimeFmtMysql mysql time format
	TimeFmtMysql = mysql.TimeFormat
	//TimeFmtDate with only date
	TimeFmtDate = "2006-01-02"
)

//TaskInfo  struct
type TaskInfo struct {
	ID           uint32 `gorm:"column:id"`
	GenerateDate string
	TaskType     int8
	StartTime    time.Time
	EndTime      time.Time
	TaskState    int16
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
}
