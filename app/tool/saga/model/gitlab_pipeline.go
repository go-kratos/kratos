package model

const (
	// HookPipelineType ...
	HookPipelineType = "pipeline"
	// PipelineFailed ...
	PipelineFailed = "failed"
	// PipelineSuccess ...
	PipelineSuccess = "success"
	// PipelineSkipped ...
	PipelineSkipped = "skipped"
	// PipelineCanceled ...
	PipelineCanceled = "canceled"
	// PipelineRunning ...
	PipelineRunning = "running"
	// PipelinePending ...
	PipelinePending = "pending"
	// MergeStatusOk ...
	MergeStatusOk = "can_be_merged"
	// MergeStateOpened ...
	MergeStateOpened = "opened"
)

// QueryStatus ...
type QueryStatus int

// query pipeline type.
const (
	QueryProcessing QueryStatus = iota
	QuerySuccess
	QuerySuccessRmNote
	QueryID
)

// HookPipeline webhook for pipeline
type HookPipeline struct {
	ObjectKind       string    `json:"object_kind"`
	User             *User     `json:"user"`
	Project          *Project  `json:"project"`
	ObjectAttributes *Pipeline `json:"object_attributes"`
	Commit           *Commit   `json:"commit"`
}

// Pipeline object_attributes for pipeline
type Pipeline struct {
	ID         int64    `json:"id"`
	Ref        string   `json:"ref"`
	Tag        bool     `json:"tag"`
	Sha        string   `json:"sha"`
	BeforeSha  string   `json:"before_sha"`
	Status     string   `json:"status"`
	Stages     []string `json:"stages"`
	CreatedAt  string   `json:"created_at"`
	FinishedAt string   `json:"finished_at"`
	Duration   uint64   `json:"duration"`
}
