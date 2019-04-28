package model

// BugIssueRequest Bug Issue Request.
type BugIssueRequest struct {
	StartNum      int
	Version       string
	Rows          int
	PlatformID    string
	ProjectID     string
	ExceptionType string
}

// BugIssueResponse Bug Issue Response.
type BugIssueResponse struct {
	Status int     `json:"status"`
	Ret    *BugRet `json:"ret"`
}

// BugRet Bug Ret.
type BugRet struct {
	NumFound  int          `json:"numFound"`
	BugIssues []*BugIssues `json:"issueList"`
}

// BugIssues Bug Issues.
type BugIssues struct {
	IssueID      string    `json:"issueId"`
	Title        string    `json:"exceptionName"`
	ExceptionMsg string    `json:"exceptionMessage"`
	KeyStack     string    `json:"keyStack"`
	LastTime     string    `json:"lastestUploadTime"`
	Count        int64     `json:"count"`
	Tags         []*BugTag `json:"tagInfoList"`
	UserCount    int64     `json:"imeiCount"`
	Version      string    `json:"version"`
}

// BugTag Bug Tag.
type BugTag struct {
	TagName string `json:"tagName"`
}

// BugIssueDetailResponse Bug Issue Detail Response.
type BugIssueDetailResponse struct {
	Code int             `json:"code"`
	Data *BugIssueDetail `json:"data"`
}

// BugIssueDetail Bug Issue Detail.
type BugIssueDetail struct {
	CallStack string `json:"callStack"`
}

// BugVersionResponse Bug Version Response.
type BugVersionResponse struct {
	Status int                   `json:"status"`
	Ret    *SelectorPropertyList `json:"ret"`
}

// SelectorPropertyList SelectorPropertyList.
type SelectorPropertyList struct {
	BugVersionList []*BugVersion `json:"versionList"`
}

// BugVersion BugVersion.
type BugVersion struct {
	Name       string `json:"name"`
	Enable     int    `json:"enable"`
	SDKVersion string `json:"sdkVersion"`
}

// BugIssueExceptionListResponse Bug Issue Exception List Response.
type BugIssueExceptionListResponse struct {
	Status int                 `json:"status"`
	Ret    *IssueExceptionList `json:"ret"`
}

// IssueExceptionList IssueExceptionList.
type IssueExceptionList struct {
	IssueException []*IssueException `json:"issueList"`
}

// IssueException IssueException.
type IssueException struct {
	IssueID string `json:"issueId"`
	Status  int    `json:"status"`
}
