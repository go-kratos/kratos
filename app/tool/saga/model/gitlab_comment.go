package model

const (
	//HookCommentTypeMR ...
	HookCommentTypeMR = "MergeRequest"
)

const (
	// CommentTypeStandard iota
	CommentTypeStandard = iota
	// CommentTypeMisaka ...
	CommentTypeMisaka
	// CommentTypeMmerge ...
	CommentTypeMmerge
	// CommentTypeMerge ...
	CommentTypeMerge
	// CommentTypeRider ...
	CommentTypeRider
	// CommentTypeDeploy ...
	CommentTypeDeploy
	// CommentTypeAddOne ...
	CommentTypeAddOne
)

// HookComment struct
type HookComment struct {
	ObjectKind       string        `json:"object_kind"`
	User             *User         `json:"user"`
	ProjectID        int64         `json:"project_id"`
	Project          *Project      `json:"project"`
	Repository       *Repository   `json:"repository"`
	ObjectAttributes *Comment      `json:"object_attributes"`
	MergeRequest     *MergeRequest `json:"merge_request"`
	Commit           *Commit       `json:"commit"`
}

// Comment struct
type Comment struct {
	ID           int64  `json:"id"`
	Note         string `json:"note"`
	NoteableType string `json:"noteable_type"`
	AuthorID     int64  `json:"author_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	ProjectID    int64  `json:"project_id"`
	Attachment   string `json:"attachment"`
	LineCode     string `json:"line_code"`
	CommitID     string `json:"commit_id"`
	NoteableID   int64  `json:"noteable_id"`
	System       bool   `json:"system"`
	STDiff       string `json:"st_diff"`
	URL          string `json:"url"`
}
