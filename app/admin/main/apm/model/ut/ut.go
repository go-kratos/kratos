package ut

import (
	"go-common/library/time"
)

// TableName .
func (*Merge) TableName() string {
	return "ut_merge"
}

// TableName .
func (*Commit) TableName() string {
	return "ut_commit"
}

// TableName .
func (*PkgAnls) TableName() string {
	return "ut_pkganls"
}

// TableName .
func (*File) TableName() string {
	return "ut_file"
}

// Merge ut_merge table from db.
type Merge struct {
	ID       int64     `gorm:"column:id" json:"id"`
	MergeID  int64     `gorm:"column:merge_id" json:"merge_id"`
	IsMerged int8      `gorm:"column:is_merged" json:"is_merged"`
	UserName string    `gorm:"column:username" json:"username"`
	Commit   *Commit   `gorm:"-" json:"commit"`
	NoteID   int       `gorm:"column:note_id"`
	CTime    time.Time `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time `gorm:"column:mtime" json:"mtime"`
}

// Commit ut_commit table from db.
type Commit struct {
	ID       int64      `gorm:"column:id" json:"id"`
	MergeID  int64      `gorm:"column:merge_id" json:"merge_id"`
	CommitID string     `gorm:"column:commit_id" json:"commit_id"`
	UserName string     `gorm:"column:username" json:"username"`
	PkgAnls  []*PkgAnls `gorm:"-" json:"pkganls"`
	CTime    time.Time  `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time  `gorm:"column:mtime" json:"mtime"`
}

// PkgAnls ut_pkganls table from db.
type PkgAnls struct {
	ID         int64     `gorm:"column:id" json:"id"`
	MergeID    int64     `gorm:"column:merge_id" json:"merge_id"`
	CommitID   string    `gorm:"column:commit_id" json:"commit_id"`
	PKG        string    `gorm:"column:pkg" json:"pkg"`
	Assertions int64     `gorm:"column:assertions" json:"assertions"`
	Passed     int64     `gorm:"column:passed" json:"passed"`
	Skipped    int64     `gorm:"column:skipped" json:"skipped"`
	Failures   int64     `gorm:"column:failures" json:"failures"`
	Panics     int64     `gorm:"column:panics" json:"panics"`
	Coverage   float64   `gorm:"column:coverage" json:"coverage"`
	Coverages  string    `gorm:"-" json:"coverages"`
	CovChange  float64   `gorm:"-" json:"cov_change"`
	PassRate   float64   `gorm:"-" json:"pass_rate"`
	PassRates  string    `gorm:"-" json:"pass_rates"`
	Score      float64   `gorm:"-" json:"score"`
	HTMLURL    string    `gorm:"column:html_url" json:"html_url"`
	ReportURL  string    `gorm:"column:report_url" json:"report_url"`
	DataURL    string    `gorm:"column:data_url" json:"data_url"`
	Files      []*File   `gorm:"-" json:"files"`
	CTime      time.Time `gorm:"column:ctime" json:"ctime"`
	MTime      time.Time `gorm:"column:mtime" json:"mtime"`
	Cids       string    `gorm:"-" json:"-"`
}

// MergeReq merge list req struct.
type MergeReq struct {
	MergeID  int64  `form:"merge_id" default:"0"`
	UserName string `form:"username" default:""`
	IsMerged int8   `form:"is_merged"`
	Pn       int    `form:"pn" default:"1"`
	Ps       int    `form:"ps" default:"20"`
}

// DetailReq .
type DetailReq struct {
	CommitID string `form:"commit_id"`
	PKG      string `form:"pkg"`
}

//HistoryCommitReq struct
type HistoryCommitReq struct {
	MergeID int64 `form:"merge_id" validate:"required"`
	//CommitID string `form:"commit_id"`
	Pn int `form:"pn" default:"1"`
	Ps int `form:"ps" default:"20"`
}

// Tyrant .
type Tyrant struct {
	Package  string  `json:"package"`
	Coverage float64 `json:"coverage"`
	PassRate float64 `json:"pass_rate"`
	Increase float64 `json:"increase"`
	LastCID  string  `json:"last_cid"`
	Standard int     `json:"standard"`
	Tyrant   bool    `json:"tyrant"`
}

// UploadRes .
type UploadRes struct {
	MergeID  int64  `form:"merge_id" validate:"required"`
	CommitID string `form:"commit_id" validate:"required"`
	UserName string `form:"username" validate:"required"`
	Author   string `form:"author"`
	PKG      string `form:"pkg" validate:"required"`
}

// SAGAResponse .
type SAGAResponse struct {
	Coverage float64 `json:"coverage"`
	PKG      string  `json:"pkg"`
}

//QATrendReq is
type QATrendReq struct {
	User      string `form:"user"`
	Period    string `form:"period" default:"day"`
	LastTime  int    `form:"last_time" default:"30"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
}

//QATrendResp is
type QATrendResp struct {
	Dates     []string  `json:"dates"`
	CommitIDs []string  `json:"commit_ids"`
	Coverages []float64 `json:"coverages"`
	PassRates []float64 `json:"pass_rates"`
	Scores    []float64 `json:"scores"`
	BaseLine  int       `json:"baseline"`
}

//CommitInfo is
type CommitInfo struct {
	MergeID      int64         `gorm:"column:merge_id" json:"merge_id"`
	CommitID     string        `gorm:"column:commit_id" json:"-"`
	MTime        time.Time     `gorm:"column:mtime" json:"mtime"`
	Coverage     float64       `gorm:"-" json:"coverage"`
	PassRate     float64       `gorm:"-" json:"pass_rate"`
	GitlabCommit *GitlabCommit `gorm:"-" json:"gitlab_commit"`
}

//GitlabCommit is
type GitlabCommit struct {
	ID         string `json:"id"`
	ShortID    string `json:"short_id"`
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
	Status     string `json:"status"`
	ProjectID  int64  `json:"project_id"`
}

// WechatUsersMsg is used for sending wechat msg for users
type WechatUsersMsg struct {
	ToUser  []string `json:"touser"`
	Content string   `json:"content"`
}

// WechatGroupMsg is used for sending wechat msg for group
type WechatGroupMsg struct {
	ChatID  string       `json:"chatid"`
	MsgType string       `json:"msgtype"`
	Text    *TextContent `json:"text"`
	Safe    int          `json:"safe"`
}

// TextContent textContent
type TextContent struct {
	Content string `json:"content"`
}

// File file
type File struct {
	ID                int64     `gorm:"column:id"`
	CommitID          string    `gorm:"column:commit_id"`
	PKG               string    `gorm:"column:pkg"`
	Name              string    `gorm:"column:name"`
	Statements        int64     `gorm:"colum:statements"`
	CoveredStatements int64     `gorm:"column:covered_statements"`
	Coverage          float64   `gorm:"-"`
	CTime             time.Time `gorm:"column:ctime"`
	MTime             time.Time `gorm:"column:mtime"`
}

// Block block
type Block struct {
	Start      int
	End        int
	Statements int
	Count      int
}
