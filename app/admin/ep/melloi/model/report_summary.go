package model

import (
	"time"
)

//ReportSummary report summary
type ReportSummary struct {
	ID                  int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ScriptID            int       `json:"script_id" form:"script_id"`
	ScriptSnapID        int       `json:"script_snap_id" form:"script_snap_id"`
	ExecuteID           string    `json:"execute_id" form:"execute_id"`
	Department          string    `json:"department" form:"department"`
	Project             string    `json:"project" form:"project"`
	APP                 string    `json:"app" form:"app"`
	TestName            string    `json:"test_name" form:"test_name" param:"test_name"`
	TestNameNick        string    `json:"test_name_nick" form:"test_name_nick"`
	JobName             string    `json:"job_name" form:"job_name"`
	Count               int       `json:"count"`
	QPS                 int       `json:"qps"`
	AvgTime             int       `json:"avg_time"`
	Min                 int       `json:"min"`
	Max                 int       `json:"max"`
	Error               int       `json:"error"`
	FailPercent         string    `json:"fail_percent"`
	NinetyTime          int       `json:"ninety_time"`
	NinetyFiveTime      int       `json:"ninety_five_time"`
	NinetyNineTime      int       `json:"ninety_nine_time"`
	NetIo               int       `json:"net_io"`
	ElapsdTime          int       `json:"elapsd_time"`
	TestStatus          int       `json:"test_status"`
	UserName            string    `json:"user_name" form:"user_name"`
	ResJtl              string    `json:"res_jtl"`
	JmeterLog           string    `json:"jmeter_log"`
	DockerSum           int       `json:"docker_sum"`
	Ctime               time.Time `json:"ctime"`
	Mtime               time.Time `json:"mtime"`
	Debug               int       `json:"debug"`
	Active              int       `json:"active" form:"active"`
	SceneID             int       `json:"scene_id" form:"scene_id"`
	Type                int       `json:"type" form:"type"` // 0.http单接口  1.grpc报告  2.场景报告  3.全链路
	LoadTime            int       `json:"load_time"`        //执行时间
	FiftyTime           int       `json:"fifty_time"`
	IsFusing            bool      `json:"is_fusing"`             //是否熔断
	FusingTestName      string    `json:"fusing_test_name"`      //被熔断接口
	SuccessCodeRate     int       `json:"success_code_rate"`     //熔断时接口的httpcode
	SuccessBusinessRate int       `json:"success_business_rate"` //熔断时接口的成功率
	FusingValue         int       `json:"fusing_value"`          //熔断阈值
	BusinessValue       int       `json:"business_value"`        //业务熔断阈值
	UseBusinessStop     bool      `json:"use_business_stop"`     //是否使用业务熔断
}

//QueryReportSuRequest query report summary request
type QueryReportSuRequest struct {
	ReportSummary
	//Script
	Pagination
	Executor  string `json:"executor" form:"executor"`
	SearchAll bool   `json:"search_all" form:"search_all"`
}

//QueryReportSuResponse query report summary response
type QueryReportSuResponse struct {
	ReportSummarys []*ReportLabels `json:"reports"`
	Pagination
}

//ReportLabels report labels
type ReportLabels struct {
	ReportSummary
	Labels []*LabelRelation `json:"labels"`
}

//TableName tablename
func (r ReportSummary) TableName() string {
	return "report_summary"
}
