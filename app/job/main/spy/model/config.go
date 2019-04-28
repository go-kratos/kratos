package model

import "time"

// config properties.
const (
	LimitBlockCount = "limitBlockCount"
	LessBlockScore  = "lessBlockScore"
	AutoBlock       = "autoBlock"
	AutoBlockOpen   = 1
)

const (
	//PunishmentTypeBlock 封禁
	PunishmentTypeBlock = 1
	//BlockReasonSize block reason size
	BlockReasonSize = 3
	//BlockLockKey cycle block
	BlockLockKey = "cycleblock"
	//ReportJobKey report job
	ReportJobKey = "reportjob"
	//DefLockTime def.
	DefLockTime int64 = 60
)

// Config def.
type Config struct {
	ID       int64
	Property string
	Name     string
	Val      string
	Ctime    time.Time
}
