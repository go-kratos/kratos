package model

const (
	// MRActionOpen ...
	MRActionOpen = "open"
	// MRActionReopen ...
	MRActionReopen = "reopen"
	// MRActionMerge ...
	MRActionMerge = "merge"
)

const (
	// MRStateOpened ...
	MRStateOpened = "opened"
	// MRStateClosed ...
	MRStateClosed = "closed"
	// MRStateMerged ...
	MRStateMerged = "merged"
)

const (
	// MRMergeOK ...
	MRMergeOK = "can_be_merged"
	// MRMergeFailed ...
	MRMergeFailed = "cannot_be_merged"
	// MRMergeUnchecked ...
	MRMergeUnchecked = "unchecked"
)

// HookMR def
type HookMR struct {
	ObjectKind       string        `json:"object_kind"`
	Project          *Project      `json:"project"`
	User             *User         `json:"user"`
	ObjectAttributes *MergeRequest `json:"object_attributes"`
	Assignee         *User         `json:"assignee"`
}

// MergeRequest struct
type MergeRequest struct {
	ID              int64    `json:"id"`
	TargetBranch    string   `json:"target_branch"`
	SourceBranch    string   `json:"source_branch"`
	SourceProjectID int64    `json:"source_project_id"`
	AuthorID        int64    `json:"author_id"`
	AssigneeID      int64    `json:"assignee_id"`
	Title           string   `json:"title"`
	CreateAt        string   `json:"created_at"`
	UpdateAt        string   `json:"updated_at"`
	STCommits       int64    `json:"st_commits"`
	STDiffs         int64    `json:"st_diffs"`
	MilestoneID     int64    `json:"milestone_id"`
	State           string   `json:"state"`
	MergeStatus     string   `json:"merge_status"`
	TargetProjectID int64    `json:"target_project_id"`
	IID             int64    `json:"iid"`
	Description     string   `json:"description"`
	Source          *Project `json:"source"`
	Target          *Project `json:"target"`
	LastCommit      *Commit  `json:"last_commit"`
	WorkInProgress  bool     `json:"work_in_progress"`
	URL             string   `json:"url"`
	Action          string   `json:"action"` // "open","update","close"
	Sha             string   `json:"sha"`
}

// MRRecord def
type MRRecord struct {
	ProjectID  int    `json:"pid"`
	MRID       int    `json:"mrid"`
	LastCommit string `json:"lc"`
	Mail       bool   `json:"mail"` // 是否发送过邮件
	NoteID     int    `json:"note"`
	Report     struct {
		TimeSpend       int64 `json:"rts"`
		MergeFlag       bool  `json:"rmf"`
		BuildFlag       bool  `json:"rbf"`
		StaticCheckFlag bool  `json:"rsf"`
		VetFlag         bool  `json:"rvf"`
		LintFlag        bool  `json:"rlf"`
		RuleFlag        bool  `json:"rrf"`
	} `json:"report"`
	Rider struct {
		BuildID      int64  `json:"ribi"`
		BuildFlag    bool   `json:"ribf"`
		BuildCommit  string `json:"ribc"`
		DeployID     int64  `json:"ridi"`
		DeployFlag   bool   `json:"ridf"`
		DeployCommit string `json:"ridc"`
	} `json:"rider"`
	Reviwers     []Reviewer `json:"mus"`
	ReviewNotify struct {
		Reviewer []string `json:"rnr"`
		Assign   string   `json:"rna"`
	} `json:"rn"`
}

// Reviewer struct
type Reviewer struct {
	Name     string `json:"mun"`
	CommitID string `json:"muci"`
}

const (
	// MRTypeCommon iota
	MRTypeCommon = iota
	// MRTypeBiz ...
	MRTypeBiz
	// MRTypeRevert ...
	MRTypeRevert
	// MRTypeInvalid ...
	MRTypeInvalid
)
