package model

// ExpiredOneDay ...
const ExpiredOneDay = 86400

// Status ...
const (
	StatusCancel = "cancel"
	StatusMerged = "merged"
	StatusClosed = "closed"
)

// number per year.
const (
	MonthNumPerYear = 12
	DayNumPerYear   = 365
	DayNumPerWeek   = 7
	DayNumPerMonth  = 30
)

// query type.
const (
	LastYearPerMonth = iota
	LastMonthPerDay
	LastYearPerDay
	LastWeekPerDay
)

// query type note.
const (
	LastYearPerMonthNote = "最近一年每月数量"
	LastMonthPerDayNote  = "上一月每天数量"
	LastYearPerDayNote   = "最近一年每天数量"
)

// query object type.
const (
	ObjectMR     = "mr"
	ObjectCommit = "commit"
	ObjectSaga   = "saga"
	ObjectRunner = "runner"
)

// KeyTypeConst ...
var KeyTypeConst = map[int]string{
	0: "LastYearPerMonth",
	1: "LastMonthPerDay",
	2: "LastYearPerDay",
	3: "LastWeekPerDay",
}

// CommitRequest ...
type CommitRequest struct {
	TeamParam
	Since    string `form:"since"`
	Until    string `form:"until"`
	Username string `form:"username"`
}

// ProjectCommit ...
type ProjectCommit struct {
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	CommitNum int    `json:"commit_num"`
}

// CommitResp ...
type CommitResp struct {
	Total         int              `json:"total"`
	ProjectCommit []*ProjectCommit `json:"commit_per_project"`
}

// ProjectDataReq ...
type ProjectDataReq struct {
	ProjectID   int    `form:"project_id" validate:"required"`
	ProjectName string `form:"project_name"`
	QueryType   int    `form:"query_type"`
	Username    string `form:"username"`
}

// ProjectDataResp ...
type ProjectDataResp struct {
	ProjectName string          `json:"project_name"`
	QueryDes    string          `json:"query_description"`
	Total       int             `json:"total"`
	Data        []*DataWithTime `json:"data_info"`
}

// TeamDataRequest ...
type TeamDataRequest struct {
	TeamParam
	QueryType int    `form:"query_type"`
	Username  string `form:"username"`
}

// TeamDataResp ...
type TeamDataResp struct {
	Department string          `json:"department"`
	Business   string          `json:"business"`
	QueryDes   string          `json:"query_description"`
	Total      int             `json:"total"`
	Data       []*DataWithTime `json:"data_info"`
}

// DataWithTime ...
type DataWithTime struct {
	TotalItem int    `json:"total_item"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// PipelineDataTime ...
type PipelineDataTime struct {
	TotalItem   int    `json:"total_item"`
	SuccessItem int    `json:"success_item"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// PipelineDataResp ...
type PipelineDataResp struct {
	Department   string              `json:"department"`
	Business     string              `json:"business"`
	QueryDes     string              `json:"query_description"`
	Total        int                 `json:"total"`
	SuccessNum   int                 `json:"success_num"`
	SuccessScale int                 `json:"success_scale"`
	Data         []*PipelineDataTime `json:"data_info"`
}

// PipelineDataReq ...
type PipelineDataReq struct {
	ProjectID      int    `form:"project_id" validate:"required"`
	ProjectName    string `form:"project_name"`
	Branch         string `form:"branch"`
	State          string `form:"state"`
	User           string `form:"user"`
	Type           int    `form:"query_type"` //0 最近一年每月数量;1 上一月每天数量;2 最近一年每天数量
	StatisticsType int    `form:"statistics_type"`
	Username       string `form:"username"`
}

// PipelineDataAvgResp ...
type PipelineDataAvgResp struct {
	ProjectName     string             `json:"project_name"`
	QueryDes        string             `json:"query_description"`
	Status          string             `json:"status"`
	Total           int                `json:"total"`
	TotalStatus     int                `json:"total_status"`
	AvgDurationTime float64            `json:"avg_duration_time"`
	AvgPendingTime  float64            `json:"avg_pending_time"`
	AvgRunningTime  float64            `json:"avg_running_time"`
	Data            []*PipelineDataAvg `json:"data_info"`
}

// PipelineDataAvg ...
type PipelineDataAvg struct {
	TotalItem       int     `json:"total_item"`
	TotalStatusItem int     `json:"total_status_item"`
	AvgDurationTime float64 `json:"avg_total_time"`
	MaxDurationTime float64 `json:"max_duration_time"`
	MinDurationTime float64 `json:"min_duration_time"`
	AvgPendingTime  float64 `json:"avg_pending_time"`
	MaxPendingTime  float64 `json:"max_pending_time"`
	MinPendingTime  float64 `json:"min_pending_time"`
	AvgRunningTime  float64 `json:"avg_running_time"`
	MaxRunningTime  float64 `json:"max_running_time"`
	MinRunningTime  float64 `json:"min_running_time"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
}

// PipelineTime ...
type PipelineTime struct {
	PendingMax   float64
	PendingMin   float64
	RunningMax   float64
	RunningMin   float64
	DurationMax  float64
	DurationMin  float64
	PendingList  []float64
	RunningList  []float64
	DurationList []float64
}

// AlertPipeline ...
type AlertPipeline struct {
	ProjectName      string
	ProjectID        int
	RunningTimeout   int
	RunningRate      int
	RunningThreshold int
	PendingTimeout   int
	PendingThreshold int
}
