package model

import "time"

// BranchDeleted ...
const BranchDeleted = true

// BranchDiffWithRequest ...
type BranchDiffWithRequest struct {
	ProjectID int    `form:"project_id"`
	Master    string `form:"comparator"`
	SortBy    string `form:"sort_by"`
	Branch    string `form:"branch"`
	Username  string `form:"username"`
}

// BranchDiffWithResponse ...
type BranchDiffWithResponse struct {
	Branch           string     `json:"branch"`
	Behind           int        `json:"behind"`
	Ahead            int        `json:"ahead"`
	LatestSyncTime   *time.Time `json:"latest_sync_time"`
	LatestUpdateTime *time.Time `json:"latest_update_time"`
}

// CommitTreeNode ...
type CommitTreeNode struct {
	CommitID  string     `json:"commit_id"`
	Parents   []string   `json:"parents"`
	CreatedAt *time.Time `json:"created_at"`
	Author    string     `json:"author"`
}

// StatisticsBranches ...
type StatisticsBranches struct {
	ID                 int    `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID          int    `json:"project_id"`
	ProjectName        string `json:"project_name"`
	CommitID           string `json:"commit_id"`
	BranchName         string `json:"branch_name"`
	Protected          bool   `json:"protected"`
	Merged             bool   `json:"merged"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
	IsDeleted          bool   `json:"is_deleted"`
}

// AggregateBranches ...
type AggregateBranches struct {
	ID               int        `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ProjectID        int        `json:"project_id"`
	ProjectName      string     `json:"project_name"`
	BranchName       string     `json:"branch_name"`
	BranchUserName   string     `json:"branch_user_name"`
	BranchMaster     string     `json:"branch_master"`
	Behind           int        `json:"behind"`
	Ahead            int        `json:"ahead"`
	LatestSyncTime   *time.Time `json:"latest_sync_time"`
	LatestUpdateTime *time.Time `json:"latest_update_time"`
	IsDeleted        bool       `json:"is_deleted"`
}
