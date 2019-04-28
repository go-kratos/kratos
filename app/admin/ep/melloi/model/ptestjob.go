package model

import "time"

//PtestJob performance test job
type PtestJob struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ScriptID   int       `json:"script_id" form:"script_id"`
	ReportSuID int       `json:"report_su_id" form:"report_su_id"`
	JobName    string    `json:"job_name" form:"job_name"`
	Active     int       `json:"active" form:"active"`
	ExecuteID  string    `json:"execute_id" form:"execute_id"`
	HostIP     string    `json:"host_ip"`
	JobIP      string    `json:"job_ip"`
	JobID      string    `json:"job_id"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

//PtestAdd model for adding performance test job
type PtestAdd struct {
	ReportSuID int    `json:"report_su_id" form:"report_su_id"`
	ScriptID   int    `json:"script_id" form:"script_id"`
	JmeterLog  string `json:"jmeter_log" form:"jmeter_log"`
	ResJtl     string `json:"res_jtl" form:"res_jtl"`
	JobName    string `json:"job_name" form:"job_name"`
	DockerSum  int    `json:"docker_sum" form:"docker_sum"`
	ScriptType int    `json:"script_type" form:"script_type"`
	ExecuteID  string `json:"execute_id" form:"execute_id"`
	SceneId    int    `json:"scene_id" form:"scene_id"`
	UserName   string `json:"user_name" form:"user_name"`
	DockerNum  int    `json:"docker_num" form:"docker_num"`
	SleepTime  int    `json:"sleep_time" form:"sleep_time"`
}

//AddReGraphTimer model for report graph timer
type AddReGraphTimer struct {
	ScriptID            int      `json:"script_id" form:"script_id"`
	ReportSuID          int      `json:"report_su_id" form:"report_su_id"`
	JobName             string   `json:"job_name" form:"job_name"`
	BeginTime           string   `json:"begin_time" form:"begin_time"`
	Token               string   `json:"token" form:"token"`
	TestNames           []string `json:"test_names" form:"test_names"`
	TestNameNicks       []string `json:"test_name_nicks" form:"test_name_nicks"`
	Fusing              int      `json:"fusing"`
	FusingList          []int    `json:"fusing_list"`
	TestType            int      `json:"test_type"`
	UseBusinessStop     bool     `json:"use_business_stop"`
	BusinessStopPercent int      `json:"business_stop_percent"`
	UseBusiStopList     []bool   `json:"use_busi_stop_list"`
	BusiStopPercentList []int    `json:"busi_stop_percent_list"`
}

//DoPtestResp doptest response
type DoPtestResp struct {
	BeginTime     string `json:"begin_time"`
	JobName       string `json:"job_name"`
	ReportSuID    int    `json:"report_su_id"`
	ScriptSnapIDs []int  `json:"script_snap_ids"`
	ScriptID      int    `json:"script_id"`
	Message       string `json:"message"`
	ScriptSnapID  int    `json:"script_snap_id"`
	JmeterLog     string `json:"jmeter_log"`
	JtlLog        string `json:"jtl_log"`
	JmxFile       string `json:"jmx_file"`
	GroupID       int    `json:"group_id"`
	RunOrder      int    `json:"run_order"`
	LoadTime      int    `json:"load_time"`
	HostIP        string `json:"host_ip"`
	SOS           string `json:"sos"`
}

//AddScene add scene
type AddScene struct {
	SceneID  int    `json:"scene_id" form:"scene_id"`
	UserName string `json:"user_name" form:"user_name"`
}

//TableName tablename
func (r PtestJob) TableName() string {
	return "ptest_job"
}

// JobInfo Job Info
type JobInfo struct {
	HostIp  string `json:"host_ip" form:"host_ip"`
	JobName string `json:"job_name" form:"job_name"`
}

// JobInfoList Job Info List
type JobInfoList struct {
	JobList []JobInfo `json:"job_list"`
}
