package model

import "go-common/library/ecode"

// AddVersionRequest Add Version Request.
type AddVersionRequest struct {
	ID             int64  `json:"id"`
	BuglyProjectID int64  `json:"bugly_project_id"`
	Version        string `json:"version"`
	Action         int    `json:"action"`
}

// AddProjectRequest Add Project Request.
type AddProjectRequest struct {
	ID            int64  `json:"id"`
	ProjectName   string `json:"project_name"`
	ProjectID     string `json:"project_id"`
	PlatformID    string `json:"platform_id"`
	ExceptionType string `json:"exception_type"`
}

// AddCookieRequest Add Cookie Request.
type AddCookieRequest struct {
	ID        int64  `json:"id"`
	Cookie    string `json:"cookie"`
	Token     string `json:"token"`
	Status    int    `json:"status"`
	QQAccount int    `json:"qq_account"`
}

// BuglyIssueImportRequest Bugly Issue Import Request.
type BuglyIssueImportRequest struct {
	IssueImportInfo []*IssueImportInfo `json:"import_info"`
}

// IssueImportInfo Issue Import Info.
type IssueImportInfo struct {
	ProjectID  string   `json:"project_id"`
	PlatformID string   `json:"platform_id"`
	Version    []string `json:"version"`
}

// Pagination Pagination.
type Pagination struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num" json:"page_num"`
}

// PaginationRep Pagination Response.
type PaginationRep struct {
	PageSize int   `json:"page_size"`
	PageNum  int   `json:"page_num"`
	Total    int64 `json:"total"`
}

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return ecode.MerlinIllegalPageNumErr
	} else if p.PageNum == 0 {
		p.PageNum = DefaultPageNum
	}
	if p.PageSize < 0 {
		return ecode.MerlinIllegalPageSizeErr
	} else if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return nil
}

// QueryBuglyVersionRequest Query Bugly Version Request.
type QueryBuglyVersionRequest struct {
	Pagination
	Version     string `json:"version"`
	ProjectName string `json:"project_name"`
	Action      int    `json:"action"`
	TaskStatus  int    `json:"task_status"`
	UpdateBy    string `json:"update_by"`
}

// QueryBuglyBatchRunsRequest Query Bugly Batch Runs Request.
type QueryBuglyBatchRunsRequest struct {
	Pagination
	Version string `json:"version"`
	Status  int    `json:"status"`
	BatchID string `json:"batch_id"`
}

// PaginateBuglyBatchRuns Paginate Bugly Batch Runs.
type PaginateBuglyBatchRuns struct {
	PaginationRep
	BuglyBatchRuns []*BuglyBatchRun `json:"bugly_batch_runs"`
}

// QueryBugRecordsRequest Query Bug Records Request.
type QueryBugRecordsRequest struct {
	Pagination
	ProjectTemplateID int64  `json:"project_template_id"`
	VersionTemplateID int64  `json:"version_template_id"`
	Operator          string `json:"operator"`
	Status            int    `json:"status"`
}

// QueryTapdBugPriorityConfsRequest Query Tapd Bug Priority Confs Request.
type QueryTapdBugPriorityConfsRequest struct {
	Pagination
	ProjectTemplateID int64  `json:"project_template_id"`
	UpdateBy          string `json:"update_by"`
	Status            int    `json:"status"`
}

// PaginateTapdBugPriorityConfs Paginate Tapd Bug Priority Confs.
type PaginateTapdBugPriorityConfs struct {
	PaginationRep
	TapdBugPriorityConfs []*TapdBugPriorityConf `json:"tapd_bug_priority_confs"`
}

// PaginateBugRecords Paginate Bug Records.
type PaginateBugRecords struct {
	PaginationRep
	TapdBugRecords []*TapdBugRecord `json:"tapd_bug_records"`
}

// QueryBuglyCookiesRequest Query Bugly Batch Runs Request.
type QueryBuglyCookiesRequest struct {
	Pagination
	QQAccount int `json:"qq_account"`
	Status    int `json:"status"`
}

// PaginateBuglyCookies Paginate Bugly Cookies.
type PaginateBuglyCookies struct {
	PaginationRep
	BuglyCookies []*BuglyCookie `json:"bugly_cookies"`
}

// PaginateBuglyProjectVersions Paginate Bugly Project Versions.
type PaginateBuglyProjectVersions struct {
	PaginationRep
	BuglyProjectVersions []*BuglyProjectVersion `json:"bugly_project_versions"`
}

// QueryTapdBugTemplateRequest Query tapd Bug Template Request.
type QueryTapdBugTemplateRequest struct {
	Pagination
	ProjectName string `json:"project_name"`
	UpdateBy    string `json:"update_by"`
}

// QueryTapdBugVersionTemplateRequest Query Tapd Bug Version Template Request.
type QueryTapdBugVersionTemplateRequest struct {
	Pagination
	ProjectID int64  `json:"project_template_id"`
	Version   string `json:"version"`
	UpdateBy  string `json:"update_by"`
}

// PaginateTapdBugTemplates Paginate Tapd Bug Template.
type PaginateTapdBugTemplates struct {
	PaginationRep
	TapdBugTemplateWithProjectNames []*TapdBugTemplateWithProjectName `json:"tapd_bug_templates"`
}

// TapdBugTemplateWithProjectName Paginate Tapd Bug Template.
type TapdBugTemplateWithProjectName struct {
	*TapdBugTemplate
	ProjectName string `json:"project_name"`
}

// PaginateTapdBugVersionTemplates Paginate Tapd Bug Version Template.
type PaginateTapdBugVersionTemplates struct {
	PaginationRep
	TapdBugVersionTemplates []*TapdBugVersionTemplate `json:"tapd_bug_version_templates"`
}

// UpdateTapdBugTplRequest Update Tapd Bug Tpl Request.
type UpdateTapdBugTplRequest struct {
	ID             int64  `json:"id" `
	WorkspaceID    string `json:"workspace_id"`
	BuglyProjectId int64  `json:"bugly_project_id" `

	IssueFilterSQL string `json:"issue_filter_sql"`
	SeverityKey    string `json:"severity_key"`

	TapdProperty
}

// UpdateTapdBugVersionTplRequest Update Tapd Bug Tpl Request.
type UpdateTapdBugVersionTplRequest struct {
	ID                int64  `json:"id" `
	Version           string `json:"version" `
	ProjectTemplateID int64  `json:"project_template_id"`

	IssueFilterSQL string `json:"issue_filter_sql"`
	SeverityKey    string `json:"severity_key"`

	TapdProperty
}

// QueryBuglyIssueRequest Query Bugly Issue Request.
type QueryBuglyIssueRequest struct {
	Pagination
	IssueNo      string `json:"issue_no"`
	Title        string `json:"title"`
	ExceptionMsg string `json:"exception_msg" `
	KeyStack     string `json:"key_stack"`
	Detail       string `json:"detail"`
	Tags         string `json:"tags"`
	Version      string `json:"version" `
	ProjectID    string `json:"project_id"`
	TapdBugID    string `json:"tapd_bug_id"`
}

// PaginateBuglyIssues Paginate Bugly Issues.
type PaginateBuglyIssues struct {
	PaginationRep
	BuglyIssues []*BuglyIssue `json:"bugly_issues"`
}

// UpdateTapdBugPriorityConfRequest Update Tapd Bug Priority Conf Request.
type UpdateTapdBugPriorityConfRequest struct {
	ID                int64  `json:"id" `
	ProjectTemplateID int64  `json:"project_template_id"`
	Urgent            int    `json:"urgent"`
	High              int    `json:"high"`
	Medium            int    `json:"medium"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Status            int    `json:"status"`
}

// QueryBuglyProjectRequest Query Bugly Project Request.
type QueryBuglyProjectRequest struct {
	Pagination
	ProjectName string `json:"project_name"`
	ProjectID   string `json:"project_id"`
	PlatformID  string `json:"platform_id"`
	UpdateBy    string `json:"update_by"`
}

// PaginateBuglyProjects Paginate Bugly Projects.
type PaginateBuglyProjects struct {
	PaginationRep
	BuglyProjects []*BuglyProject `json:"bugly_projects"`
}

// BuglyProjectVersion Bugly Project Version.
type BuglyProjectVersion struct {
	BuglyVersion
	ProjectName   string `json:"project_name"`
	ExceptionType string `json:"exception_type"`
}

// TapdBugTemplateShortResponse Tapd Bug Template Short Response.
type TapdBugTemplateShortResponse struct {
	ID               int64  `json:"id"`
	WorkspaceID      string `json:"workspace_id"`
	BuglyProjectId   int64  `json:"bugly_project_id"`
	BuglyProjectName string `json:"project_name"`
}
