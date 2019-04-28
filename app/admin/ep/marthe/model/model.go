package model

import (
	"time"
)

// BuglyIssue Issue Record.
type BuglyIssue struct {
	ID           int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	IssueNo      string    `json:"issue_no" gorm:"column:issue_no"`
	Title        string    `json:"title" gorm:"column:title"`
	ExceptionMsg string    `json:"exception_msg" gorm:"column:exception_msg"`
	KeyStack     string    `json:"key_stack" gorm:"column:key_stack"`
	Detail       string    `json:"detail" gorm:"column:detail"`
	Tags         string    `json:"tags" gorm:"column:tags"`
	LastTime     time.Time `json:"last_time" gorm:"column:last_time"`
	HappenTimes  int       `json:"happen_times" gorm:"column:happen_times"`
	UserTimes    int       `json:"user_times" gorm:"column:user_times"`
	Version      string    `json:"version" gorm:"column:version"`
	ProjectID    string    `json:"project_id" gorm:"column:project_id"`
	IssueLink    string    `json:"issue_link" gorm:"column:issue_link"`
	TapdBugID    string    `json:"tapd_bug_id" gorm:"column:tapd_bug_id"`
	CTime        time.Time `json:"ctime" gorm:"column:ctime"`
	MTime        time.Time `json:"mtime" gorm:"column:mtime"`
}

// TapdProperty TapdProperty.
type TapdProperty struct {
	Title            string `json:"title" gorm:"column:title"`
	Description      string `json:"description" gorm:"column:description"`
	CurrentOwner     string `json:"current_owner" gorm:"column:current_owner"`
	Platform         string `json:"platform" gorm:"column:platform"`
	Module           string `json:"module" gorm:"column:module"`
	IterationID      string `json:"iteration_id" gorm:"column:iteration_id"`
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
}

// TapdBugTemplate TapdBugTemplate.
type TapdBugTemplate struct {
	ID             int64  `json:"id" gorm:"auto_increment;primary_key;column:id"`
	WorkspaceID    string `json:"workspace_id" gorm:"column:workspace_id"`
	BuglyProjectId int64  `json:"bugly_project_id" gorm:"column:bugly_project_id"`

	TapdProperty

	IssueFilterSQL string    `json:"issue_filter_sql" gorm:"column:issue_filter_sql"`
	SeverityKey    string    `json:"severity_key" gorm:"column:severity_key"`
	CTime          time.Time `json:"ctime" gorm:"column:ctime"`
	MTime          time.Time `json:"mtime" gorm:"column:mtime"`
	UpdateBy       string    `json:"update_by" gorm:"column:update_by"`
}

// TapdBugVersionTemplate TapdBugVersionTemplate.
type TapdBugVersionTemplate struct {
	ID                int64  `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Version           string `json:"version" gorm:"column:version"`
	ProjectTemplateID int64  `json:"project_template_id" gorm:"column:project_template_id"`

	TapdProperty

	IssueFilterSQL string    `json:"issue_filter_sql" gorm:"column:issue_filter_sql"`
	SeverityKey    string    `json:"severity_key" gorm:"column:severity_key"`
	CTime          time.Time `json:"ctime" gorm:"column:ctime"`
	MTime          time.Time `json:"mtime" gorm:"column:mtime"`
	UpdateBy       string    `json:"update_by" gorm:"column:update_by"`
}

// BuglyVersion Bugly Version Record.
type BuglyVersion struct {
	ID             int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Version        string    `json:"version" gorm:"column:version"`
	BuglyProjectID int64     `json:"bugly_project_id" gorm:"column:bugly_project_id"`
	Action         int       `json:"action" gorm:"column:action"`
	TaskStatus     int       `json:"task_status" gorm:"column:task_status"`
	UpdateBy       string    `json:"update_by" gorm:"column:update_by"`
	CTime          time.Time `json:"ctime" gorm:"column:ctime"`
	MTime          time.Time `json:"mtime" gorm:"column:mtime"`
}

// BuglyBatchRun Bugly Batch Run.
type BuglyBatchRun struct {
	ID             int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	BuglyVersionID int64     `json:"bugly_version_id" gorm:"column:bugly_version_id"`
	Version        string    `json:"version" gorm:"column:version"`
	BatchID        string    `json:"batch_id" gorm:"column:batch_id"`
	RetryCount     int       `json:"retry_count" gorm:"retry_times:retry_count"`
	Status         int       `json:"status" gorm:"column:status"`
	ErrorMsg       string    `json:"error_msg" gorm:"column:error_msg"`
	CTime          time.Time `json:"ctime" gorm:"column:ctime"`
	MTime          time.Time `json:"mtime" gorm:"column:mtime"`
	EndTime        time.Time `json:"end_time" gorm:"column:end_time"`
}

// BuglyCookie Bugly Cookie.
type BuglyCookie struct {
	ID         int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	QQAccount  int       `json:"qq_account" gorm:"column:qq_account"`
	Cookie     string    `json:"cookie" gorm:"column:cookie"`
	Token      string    `json:"token" gorm:"column:token"`
	UsageCount int       `json:"usage_count" gorm:"column:usage_count"`
	Status     int       `json:"status" gorm:"column:status"`
	UpdateBy   string    `json:"update_by" gorm:"column:update_by"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
	MTime      time.Time `json:"mtime" gorm:"column:mtime"`
}

// User User.
type User struct {
	ID           int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Name         string    `json:"username" gorm:"column:name"`
	EMail        string    `json:"email" gorm:"column:email"`
	VisibleBugly bool      `json:"visible_bugly" gorm:"column:visible_bugly"`
	CTime        time.Time `gorm:"column:ctime;default:current_timestamp"`
	UTime        time.Time `gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// TapdBugRecord Tapd Bug Insert Log.
type TapdBugRecord struct {
	ID                int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	ProjectTemplateID int64     `json:"project_template_id" gorm:"column:project_template_id"`
	VersionTemplateID int64     `json:"version_template_id" gorm:"column:version_template_id"`
	Operator          string    `json:"operator" gorm:"column:operator"`
	Count             int       `json:"count" gorm:"column:count"`
	Status            int       `json:"status" gorm:"column:status"`
	IssueFilterSQL    string    `json:"issue_filter_sql" gorm:"column:issue_filter_sql"`
	CTime             time.Time `json:"ctime" gorm:"column:ctime"`
	MTime             time.Time `json:"mtime" gorm:"column:mtime"`
}

// ScheduleTask Schedule Task.
type ScheduleTask struct {
	ID     int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Name   string    `json:"username" gorm:"column:name"`
	Status int       `json:"status" gorm:"column:status"`
	CTime  time.Time `gorm:"column:ctime;default:current_timestamp"`
	MTime  time.Time `gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// TapdBugPriorityConf Tapd Bug Priority Conf.
type TapdBugPriorityConf struct {
	ID                int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	ProjectTemplateID int64     `json:"project_template_id" gorm:"column:project_template_id"`
	Urgent            int       `json:"urgent" gorm:"column:urgent"`
	High              int       `json:"high" gorm:"column:high"`
	Medium            int       `json:"medium" gorm:"column:medium"`
	StartTime         time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime           time.Time `json:"end_time" gorm:"column:end_time"`
	CTime             time.Time `json:"ctime" gorm:"column:ctime"`
	MTime             time.Time `json:"mtime" gorm:"column:mtime"`
	UpdateBy          string    `json:"update_by" gorm:"column:update_by"`
	Status            int       `json:"status" gorm:"column:status"`
}

// ContactInfo Contact Info
type ContactInfo struct {
	ID       int64     `json:"id" gorm:"column:id"`
	UserName string    `json:"username" gorm:"column:username"`
	UserID   string    `json:"user_id" gorm:"column:user_id"`
	NickName string    `json:"nick_name" gorm:"column:nick_name"`
	CTime    time.Time `json:"ctime" gorm:"column:ctime"`
	MTime    time.Time `json:"mtime" gorm:"column:mtime"`
}

// BuglyProject Bugly Project.
type BuglyProject struct {
	ID            int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	ProjectID     string    `json:"project_id" gorm:"column:project_id"`
	ProjectName   string    `json:"project_name" gorm:"column:project_name"`
	PlatformID    string    `json:"platform_id" gorm:"column:platform_id"`
	UpdateBy      string    `json:"update_by" gorm:"column:update_by"`
	ExceptionType string    `json:"exception_type" gorm:"column:exception_type"`
	CTime         time.Time `json:"ctime" gorm:"column:ctime"`
	MTime         time.Time `json:"mtime" gorm:"column:mtime"`
}
