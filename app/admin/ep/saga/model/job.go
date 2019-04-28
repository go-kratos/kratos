package model

import "time"

// ProjectJobRequest ...
type ProjectJobRequest struct {
	ProjectID      int    `form:"project_id"`
	Scope          string `form:"state"`
	User           string `form:"user"`
	Branch         string `form:"branch"`
	Machine        string `form:"machine"`
	StatisticsType int    `form:"statistics_type"`
	Username       string `form:"username"`
}

// ProjectJobResp ...
type ProjectJobResp struct {
	ProjectID        int            `json:"project_id"`
	QueryDescription string         `json:"query_description"`
	TotalItem        int            `json:"total"`
	State            string         `json:"state"`
	DataInfo         []*DateJobInfo `json:"data_info"`
}

// DateJobInfo ...
type DateJobInfo struct {
	Date              string        `json:"date"`
	JobTotal          int           `json:"total_num"`
	StatusNum         int           `json:"status_num"`
	PendingTime       float64       `json:"pending_time"`
	RunningTime       float64       `json:"running_time"`
	SlowestPendingJob []*ProjectJob `json:"slowest_pending_jobs"`
}

// ProjectJob ...
type ProjectJob struct {
	Status     string
	User       string
	Branch     string
	Machine    string
	CreatedAt  *time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
}
