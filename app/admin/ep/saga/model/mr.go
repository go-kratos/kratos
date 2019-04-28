package model

import "time"

// ProjectMrReportReq ...
type ProjectMrReportReq struct {
	ProjectID int    `form:"project_id"`
	Member    string `form:"member"`
	Username  string `form:"username"`
}

// ProjectMrReportResp ...
type ProjectMrReportResp struct {
	ChangeAdd       int      `json:"change_add"`
	ChangeDel       int      `json:"change_del"`
	MrCount         int      `json:"mr_count"`
	StateCount      int      `json:"merged_count"`
	Discussion      int      `json:"discussion"`
	Resolve         int      `json:"resolved_discussion"`
	AverageMerge    string   `json:"average_merge_time"`
	Reviewers       []string `json:"reviewer"`
	SpentTime       string   `json:"spent_time"`
	ReviewerOther   []string `json:"reviewer_other"`
	ReviewChangeAdd int      `json:"review_add"`
	ReviewChangeDel int      `json:"review_del"`
	ReviewTotalTime string   `json:"review_total_time"`
}

// MrReviewer ...
type MrReviewer struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	FinishedAt *time.Time `json:"finished_at"`
	UserType   string     `json:"type"`
	// SpentTime 其实是反应时间+review时间
	SpentTime int `json:"spent_time"`
}

// MrInfo ...
type MrInfo struct {
	ProjectID        int           `json:"project_id"`
	MrID             int           `json:"mr_id"`
	State            string        `json:"state"`
	SpentTime        int           `json:"spent_time"`
	Author           string        `json:"author"`
	ChangeAdd        int           `json:"change_add"`
	ChangeDel        int           `json:"change_del"`
	TotalDiscussion  int           `json:"total_discussion"`
	SolvedDiscussion int           `json:"solved_discussion"`
	Reviewers        []*MrReviewer `json:"reviewers"`
}
