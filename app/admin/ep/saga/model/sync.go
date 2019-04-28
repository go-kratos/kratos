package model

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

// DatabaseErrorText ...
const (
	DatabaseErrorText         = "Incorrect string value"
	DatabaseMaxLenthErrorText = "Data too long for column"
	MessageMaxLen             = 2048
	JsonMarshalErrorText      = "XXXXX"
)

// DataType ...
const (
	DataTypePipeline = "pipeline"
	DataTypeJob      = "job"
	DataTypeCommit   = "commit"
	DataTypeMR       = "MR"
	DataTypeBranch   = "Branch"
)

// FailData ...
type FailData struct {
	ChildID    int
	ChildIDStr string
	SunID      int
}

// SyncResult ...
type SyncResult struct {
	TotalPage int
	TotalNum  int
	FailData  []*FailData
}

// StatisticsCommits ...
type StatisticsCommits struct {
	ID             int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	CommitID       string     `json:"commit_id"`
	ProjectID      int        `json:"project_id" gorm:"column:project_id"`
	ProjectName    string     `json:"project_name"`
	ShortID        string     `json:"short_id"`
	Title          string     `json:"title"`
	AuthorName     string     `json:"author_name"`
	AuthoredDate   *time.Time `json:"authored_date"`
	CommitterName  string     `json:"committer_name"`
	CommittedDate  *time.Time `json:"committed_date"`
	CreatedAt      *time.Time `json:"created_at"`
	Message        string     `json:"message"`
	ParentIDs      string     `json:"parent_ids"`
	StatsAdditions int        `json:"stats_additions"`
	StatsDeletions int        `json:"stats_deletions"`
	Status         string     `json:"status" default:""`
}

// StatisticsIssues ...
type StatisticsIssues struct {
	ID               int             `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID        int             `json:"project_id"`
	IssueID          int             `json:"issue_id"`
	IssueIID         int             `json:"issue_iid" gorm:"column:issue_iid"`
	MilestoneID      int             `json:"milestone_id"`
	AuthorID         int             `json:"author_id"`
	AuthorName       string          `json:"author_name"`
	Description      string          `json:"description"`
	State            string          `json:"state"`
	Assignees        string          `json:"assignees"`
	AssigneeID       int             `json:"assignee_id"`
	AssigneeName     string          `json:"assignee_name"`
	Upvotes          int             `json:"upvotes"`
	Downvotes        int             `json:"downvotes"`
	Labels           string          `json:"labels"`
	Title            string          `json:"title"`
	UpdatedAt        *time.Time      `json:"updated_at"`
	CreatedAt        *time.Time      `json:"created_at"`
	ClosedAt         *time.Time      `json:"closed_at"`
	Subscribed       bool            `json:"subscribed"`
	UserNotesCount   int             `json:"user_notes_count"`
	DueDate          *gitlab.ISOTime `json:"due_date"`
	WebURL           string          `json:"web_url"`
	TimeStats        string          `json:"time_stats"`
	Confidential     bool            `json:"confidential"`
	Weight           int             `json:"weight"`
	DiscussionLocked bool            `json:"discussion_locked"`
	IssueLinkID      int             `json:"issue_link_id"`
}

// StatisticsRunners ...
type StatisticsRunners struct {
	ID          int    `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"project_name"`
	RunnerID    int    `json:"runner_id"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	IsShared    bool   `json:"is_shared"`
	IPAddress   string `json:"ip_address"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	Status      string `json:"status"`
	Token       string `json:"token"`
}

// StatisticsJobs ...
type StatisticsJobs struct {
	ID                int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID         int        `json:"project_id"`
	ProjectName       string     `json:"project_name"`
	CommitID          string     `json:"commit_id"`
	CreatedAt         *time.Time `json:"created_at"`
	Coverage          float64    `json:"coverage"`
	ArtifactsFile     string     `json:"artifacts_file"`
	FinishedAt        *time.Time `json:"finished_at"`
	JobID             int        `json:"job_id"`
	Name              string     `json:"name"`
	Ref               string     `json:"ref"`
	RunnerID          int        `json:"runner_id"`
	RunnerDescription string     `json:"runner_description"`
	Stage             string     `json:"stage"`
	StartedAt         *time.Time `json:"started_at"`
	Status            string     `json:"status"`
	Tag               bool       `json:"tag"`
	UserID            int        `json:"user_id"`
	UserName          string     `json:"user_name"`
	WebURL            string     `json:"web_url"`
}

// StatisticsMrs ...
type StatisticsMrs struct {
	ID                           int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	MRID                         int        `json:"mr_id"`
	MRIID                        int        `json:"mr_iid" gorm:"column:mr_iid"`
	TargetBranch                 string     `json:"target_branch"`
	SourceBranch                 string     `json:"source_branch"`
	ProjectID                    int        `json:"project_id"`
	ProjectName                  string     `json:"project_name"`
	Title                        string     `json:"title"`
	State                        string     `json:"state"`
	CreatedAt                    *time.Time `json:"created_at"`
	UpdatedAt                    *time.Time `json:"updated_at"`
	Upvotes                      int        `json:"upvotes"`
	Downvotes                    int        `json:"downvotes"`
	AuthorID                     int        `json:"author_id"`
	AuthorName                   string     `json:"author_name"`
	AssigneeID                   int        `json:"assignee_id"`
	AssigneeName                 string     `json:"assignee_name"`
	SourceProjectID              int        `json:"source_project_id"`
	TargetProjectID              int        `json:"target_project_id"`
	Labels                       string     `json:"labels"`
	Description                  string     `json:"description"`
	WorkInProgress               bool       `json:"work_in_progress"`
	MilestoneID                  int        `json:"milestone_id"`
	MergeWhenPipelineSucceeds    bool       `json:"merge_when_pipeline_succeeds"`
	MergeStatus                  string     `json:"merge_status"`
	MergedByID                   int        `json:"merged_by_id"`
	MergedByName                 string     `json:"merged_by_name"`
	MergedAt                     *time.Time `json:"merged_at"`
	ClosedByID                   int        `json:"closed_by_id"`
	ClosedAt                     *time.Time `json:"closed_at"`
	Subscribed                   bool       `json:"subscribed"`
	SHA                          string     `json:"sha"`
	MergeCommitSHA               string     `json:"merge_commit_sha"`
	UserNotesCount               int        `json:"user_notes_count"`
	ChangesCount                 string     `json:"changes_count"`
	ShouldRemoveSourceBranch     bool       `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch      bool       `json:"force_remove_source_branch"`
	WebURL                       string     `json:"web_url"`
	DiscussionLocked             bool       `json:"discussion_locked"`
	Changes                      string     `json:"changes"`
	TimeStatsHumanTimeEstimate   string     `json:"time_stats_human_time_estimate"`
	TimeStatsHumanTotalTimeSpent string     `json:"time_stats_human_total_time_spent"`
	TimeStatsTimeEstimate        int        `json:"time_stats_time_estimate"`
	TimeStatsTotalTimeSpent      int        `json:"time_stats_total_time_spent"`
	Squash                       bool       `json:"squash"`
	PipelineID                   int        `json:"pipeline_id"`
	ChangeAdd                    int        `json:"change_add"`
	ChangeDel                    int        `json:"change_del"`
	TotalDiscussion              int        `json:"total_discussion"`
	SolvedDiscussion             int        `json:"solved_discussion"`
}

// AggregateMrReviewer ...
type AggregateMrReviewer struct {
	ID            int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID     int        `json:"project_id"`
	ProjectName   string     `json:"project_name"`
	MrIID         int        `json:"mr_iid" gorm:"column:mr_iid"`
	Title         string     `json:"title"`
	WebUrl        string     `json:"web_url"`
	AuthorName    string     `json:"author_name"`
	ReviewerID    int        `json:"reviewer_id"`
	ReviewerName  string     `json:"reviewer_name"`
	ReviewType    string     `json:"review_type"`
	ReviewID      int        `json:"review_id"`
	ReviewCommand string     `json:"review_command"`
	CreatedAt     *time.Time `json:"created_at"`
	UserType      string     `json:"type"`
	ApproveTime   int        `json:"approve_time"` // SpentTime 其实是反应时间+review时间
	MergeTime     int        `json:"merge_time"`
}

// StatisticsPipeline ...
type StatisticsPipeline struct {
	ID           int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	PipelineID   int        `json:"pipeline_id" gorm:"column:pipeline_id"`
	ProjectName  string     `json:"project_name"`
	ProjectID    int        `json:"project_id" gorm:"column:project_id"`
	Status       string     `json:"status" gorm:"column:status" default:""`
	Ref          string     `json:"ref" gorm:"column:ref"`
	Tag          bool       `json:"tag" gorm:"column:tag"`
	User         string     `json:"user" gorm:"column:user"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt    *time.Time `json:"created_at" gorm:"column:created_at"`
	StartedAt    *time.Time `json:"started_at" gorm:"column:started_at"`
	FinishedAt   *time.Time `json:"finished_at" gorm:"column:finished_at"`
	CommittedAt  *time.Time `json:"committed_at" gorm:"column:committed_at"`
	Duration     int        `json:"duration" gorm:"column:duration"`
	Coverage     string     `json:"coverage" gorm:"column:coverage"`
	DurationTime int        `json:"duration_time"`
}

// StatisticsNotes ...
type StatisticsNotes struct {
	ID             int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID      int        `json:"project_id"`
	ProjectName    string     `json:"project_name"`
	MrIID          int        `json:"mr_iid"  gorm:"column:mr_iid"`
	IssueIID       int        `json:"issue_iid"  gorm:"column:issue_iid"`
	NoteID         int        `json:"note_id"`
	Body           string     `json:"body"`
	Attachment     string     `json:"attachment"`
	Title          string     `json:"title"`
	FileName       string     `json:"file_name"`
	AuthorID       int        `json:"author_id"`
	AuthorName     string     `json:"author_name"`
	System         bool       `json:"system"`
	ExpiresAt      *time.Time `json:"expires_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	CreatedAt      *time.Time `json:"created_at"`
	NoteableID     int        `json:"noteable_id"`
	NoteableType   string     `json:"noteable_type"`
	Position       string     `json:"position"`
	Resolvable     bool       `json:"resolvable"`
	Resolved       bool       `json:"resolved"`
	ResolvedByID   int        `json:"resolved_by_id"`
	ResolvedByName string     `json:"resolved_by_name"`
	NoteableIID    int        `json:"noteable_iid" gorm:"column:noteable_iid"`
}

// StatisticsMembers ...
type StatisticsMembers struct {
	ID          int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID   int        `json:"project_id"`
	ProjectName string     `json:"project_name"`
	MemberID    int        `json:"member_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	State       string     `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	AccessLevel int        `json:"access_level"`
}

// StatisticsMRAwardEmojis ...
type StatisticsMRAwardEmojis struct {
	ID            int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID     int        `json:"project_id"`
	ProjectName   string     `json:"project_name"`
	MrIID         int        `json:"mr_iid" gorm:"column:mr_iid"`
	AwardEmojiID  int        `json:"award_emoji_id"`
	Name          string     `json:"name"`
	UserID        int        `json:"user_id"`
	UserName      string     `json:"user_name"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	AwardableID   int        `json:"awardable_id"`
	AwardableType string     `json:"awardable_type"`
}

// StatisticsDiscussions ...
type StatisticsDiscussions struct {
	ID             int    `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name"`
	MrIID          int    `json:"mr_iid" gorm:"column:mr_iid"`
	DiscussionID   string `json:"discussion_id"`
	IndividualNote bool   `json:"individual_note"`
	Notes          string `json:"notes"`
}
