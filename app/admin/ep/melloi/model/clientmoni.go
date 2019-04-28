package model

import "time"

//ClientMoni model for performance test container cpu monitor
type ClientMoni struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ScriptID   int       `json:"script_id" form:"script_id"`
	ReportSuID int       `json:"report_su_id" form:"report_su_id"`
	JobName    string    `json:"job_name" form:"job_name"`
	JobNameAll string    `json:"job_name_all" form:"job_name_all"`
	CPUUsed    string    `json:"cpu_used" form:"cpu_used"`
	ElapsdTime int       `json:"elapsd_time"`
	Ctime      time.Time `json:"ctime" form:"ctime"`
	Mtime      time.Time `json:"mtime" form:"mtime"`
}

//TableName table name of client moni model
func (c ClientMoni) TableName() string {
	return "client_moni"
}
