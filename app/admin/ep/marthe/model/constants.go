package model

// TimeLayout
const (
	TimeLayout = "2006-01-02 15:04:05"
)

// page
const (
	DefaultPageSize = 5
	DefaultPageNum  = 1
)

// BuglyVersion task status
const (
	BuglyVersionTaskStatusReady   = 1
	BuglyVersionTaskStatusRunning = 2
)

// BuglyVersion action
const (
	BuglyVersionActionEnable  = 1
	BuglyVersionActionDisable = 2
)

// BuglyCookie status
const (
	BuglyCookieStatusEnable  = 1
	BuglyCookieStatusDisable = 2
)

// BuglyBatchRun status
const (
	BuglyBatchRunStatusRunning = 1
	BuglyBatchRunStatusDone    = 2
	BuglyBatchRunStatusFailed  = 3
)

// Insert bug status
const (
	InsertBugStatusRunning = 1
	InsertBugStatusDone    = 2
	InsertBugStatusFailed  = 3
)

// Tapd Bug Priority Conf Enable
const (
	TapdBugPriorityConfEnable  = 1
	TapdBugPriorityConfDisable = 2
)

// Task Status
const (
	TaskStatusRunning = 1
	TaskStatusDone    = 2
	TaskStatusFailed  = 3
)

// Task Type
const (
	TaskBatchRunVersions        = "BatchRunVersions"
	TaskDisableBatchRunOverTime = "DisableBatchRunOverTime"
	TaskBatchRunUpdateBugInTapd = "BatchRunUpdateBugInTapd"
	TaskSyncWechatContact       = "SyncWechatContact"
)
