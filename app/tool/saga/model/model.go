package model

import (
	"reflect"
	"regexp"
)

// reids锁
const (
	SagaTask        = "_SagaTask_%d"
	SagaRepoLockKey = "_SagaRepoLockKey_%d"
	SagaLockValue   = "_SagaLockValue"
)

// gitlab指令
const (
	SagaCommandPlusOne  = "+ok"
	SagaCommandMerge    = "+mr"
	SagaCommandPlusOne1 = "+1"
	SagaCommandMerge1   = "+merge"
)

// 任务状态
const (
	TaskStatusFailed  = 1 // 任务失败
	TaskStatusSuccess = 2 // 任务成功
	TaskStatusRunning = 3 // 任务运行中
	TaskStatusWaiting = 4 // 任务等待
)

// CONTRIBUTORS define
const (
	SagaContributorsName = "CONTRIBUTORS.md"
)

// RepoConfig def
type RepoConfig struct {
	URL                 string
	Group               string
	Name                string
	GName               string // gitlab仓库别名
	Language            string
	AuthBranches        []string // 鉴权分支
	TargetBranches      []string // 分支白名单
	TargetBranchRegexes []*regexp.Regexp
	LockTimeout         int32
	MinReviewer         int
	RelatePipeline      bool
	DelayMerge          bool
	LimitAuth           bool
	AllowLabel          string
	SuperAuthUsers      []string
}

// RequireReviewFolder ...
type RequireReviewFolder struct {
	Folder    string
	Owners    []string
	Reviewers []string
}

// AuthUsers ...
type AuthUsers struct {
	Owners    []string
	Reviewers []string
}

// ContactInfo def
type ContactInfo struct {
	ID          string `json:"id,omitempty" gorm:"column:id"`
	UserName    string `json:"english_name" gorm:"column:user_name"`
	UserID      string `json:"userid" gorm:"column:user_id"`
	NickName    string `json:"name" gorm:"column:nick_name"`
	VisibleSaga bool   `json:"visible_saga" gorm:"column:visible_saga"`
}

// RequireVisibleUser def
type RequireVisibleUser struct {
	UserName string
	NickName string
}

// AlmostEqual return the compare result with fields
func (contact *ContactInfo) AlmostEqual(other *ContactInfo) bool {
	if contact.UserID == other.UserID &&
		contact.UserName == other.UserName &&
		contact.NickName == other.NickName {
		return true
	}
	return false
}

// TaskInfo ...
type TaskInfo struct {
	NoteID int
	Event  *HookComment
	Repo   *Repo
}

// MergeInfo ...
type MergeInfo struct {
	PipelineID   int
	NoteID       int
	AuthorID     int
	UserName     string
	MRIID        int
	ProjID       int
	URL          string
	AuthBranches []string
	SourceBranch string
	TargetBranch string
	MinReviewer  int
	LockTimeout  int32
	Title        string
	Description  string
}

// Repo structure
type Repo struct {
	Config *RepoConfig
}

// Update if config is changed
func (r *Repo) Update(conf *RepoConfig) bool {
	if r.confEqual(conf) {
		return false
	}
	r.Config = conf
	return true
}

func (r *Repo) confEqual(conf *RepoConfig) bool {
	if r.Config.URL == conf.URL &&
		r.Config.Group == conf.Group &&
		r.Config.Name == conf.Name &&
		r.Config.GName == conf.GName &&
		r.Config.Language == conf.Language &&
		reflect.DeepEqual(r.Config.AuthBranches, conf.AuthBranches) &&
		reflect.DeepEqual(r.Config.TargetBranches, conf.TargetBranches) &&
		r.Config.LockTimeout == conf.LockTimeout &&
		r.Config.MinReviewer == conf.MinReviewer &&
		r.Config.RelatePipeline == conf.RelatePipeline &&
		r.Config.DelayMerge == conf.DelayMerge &&
		r.Config.LimitAuth == conf.LimitAuth &&
		r.Config.AllowLabel == conf.AllowLabel &&
		reflect.DeepEqual(r.Config.SuperAuthUsers, conf.SuperAuthUsers) {
		return true
	}
	return false
}

// AuthUpdate ...
func (r *Repo) AuthUpdate(conf *RepoConfig) bool {
	if r.Config.Group == conf.Group &&
		r.Config.Name == conf.Name &&
		reflect.DeepEqual(r.Config.AuthBranches, conf.AuthBranches) {
		return false
	}
	return true
}

// WebHookUpdate ...
func (r *Repo) WebHookUpdate(conf *RepoConfig) bool {
	return r.Config.URL != conf.URL
}
