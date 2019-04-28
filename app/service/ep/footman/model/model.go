package model

import "time"

// IssueRecord Issue Record.
type IssueRecord struct {
	ID           int64     `json:"id" gorm:"column:id"`
	IssueNo      string    `json:"issue_no" gorm:"column:issue_no"`
	Title        string    `json:"title" gorm:"column:title"`
	ExceptionMsg string    `json:"exception_msg" gorm:"column:exception_msg"`
	KeyStack     string    `json:"key_stack" gorm:"column:key_stack"`
	Detail       string    `json:"detail" gorm:"column:detail"`
	Tags         string    `json:"tags" gorm:"column:tags"`
	LastTime     time.Time `json:"last_time" gorm:"column:last_time"`
	HappenTimes  int64     `json:"happen_times" gorm:"column:happen_times"`
	UserTimes    int64     `json:"user_times" gorm:"column:user_times"`
	Version      string    `json:"version" gorm:"column:version"`
	ProjectID    string    `json:"project_id" gorm:"column:project_id"`
	IssueLink    string    `json:"issue_link" gorm:"column:issue_link"`
	TapdBugID    string    `json:"tapd_bug_id" gorm:"column:tapd_bug_id"`
}

// IssueLastTime Issue Last Time.
type IssueLastTime struct {
	ID       int64     `json:"id" gorm:"column:id"`
	LastTime time.Time `json:"last_time" gorm:"column:last_time"`
	Version  string    `json:"version" gorm:"column:version"`
	//1-正在执行中，0未执行或已执行完
	TaskStatus int    `json:"task_status" gorm:"column:task_status"`
	LastIssue  string `json:"last_issue" gorm:"column:last_issue"`
}

// BugTemplate BugTemplate.
type BugTemplate struct {
	ID               int64  `json:"id" gorm:"column:id"`
	WorkspaceID      string `json:"workspace_id" gorm:"column:workspace_id"`
	ProjectID        string `json:"project_id" gorm:"column:project_id"`
	PlatformID       string `json:"platform_id" gorm:"column:platform_id"`
	Title            string `json:"title" gorm:"column:title"`
	Description      string `json:"description" gorm:"column:description"`
	CurrentOwner     string `json:"current_owner" gorm:"column:current_owner"`
	Platform         string `json:"platform" gorm:"column:platform"`
	Module           string `json:"module" gorm:"column:module"`
	ReleaseID        string `json:"release_id" gorm:"column:release_id"`
	Priority         string `json:"priority" gorm:"column:priority"`
	Severity         string `json:"severity" gorm:"column:severity"`
	Source           string `json:"source" gorm:"column:source"`
	CustomFieldFour  string `json:"custom_field_four" gorm:"column:custom_field_four"`
	BugType          string `json:"bugtype" gorm:"column:bugtype"`
	OriginPhase      string `json:"originphase" gorm:"column:originphase"`
	CustomFieldThree string `json:"custom_field_three" gorm:"column:custom_field_three"`
	Reporter         string `json:"reporter" gorm:"column:reporter"`
	Status           string `json:"status" gorm:"column:status"`
	IssueFilterSQL   string `json:"issue_filter_sql" gorm:"column:issue_filter_sql"`
	SeverityKey      string `json:"severity_key" gorm:"column:severity_key"`
}

// StoryWallTimeModel Story Wall Time Model.
type StoryWallTimeModel struct {
	StepStartTime time.Time
	StepEndTime   time.Time
}
