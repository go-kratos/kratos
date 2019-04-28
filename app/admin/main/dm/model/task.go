package model

// TaskStateJump dm task jump a queue
const (
	TaskRegexLen = 250 // dm task regex max length

	// dm task state
	TaskReviewPass  = int32(2)
	TaskStateRun    = int32(3)
	TaskStateFailed = int32(4)
)

// TaskList dm task info list
type TaskList struct {
	Page   *PageInfo
	Result []*TaskInfo `json:"result"`
}

// TaskView .
type TaskView struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Creator  string   `json:"creator"`
	Reviewer string   `json:"reviewer"`
	Regex    string   `json:"regex"`
	KeyWords string   `json:"keywords"`
	IPs      string   `json:"ips"`
	Mids     string   `json:"mids"`
	Cids     string   `json:"cids"`
	Start    string   `json:"start"`
	End      string   `json:"end"`
	QCount   int64    `json:"qcount"` //查询总数
	Tcount   int64    `json:"tcount"` //删除总数
	State    int32    `json:"state"`
	Result   string   `json:"-"`
	Ctime    string   `json:"ctime"`
	Mtime    string   `json:"mtime"`
	SubTask  *SubTask `json:"sub,omitempty"`
}

// SubTask .
type SubTask struct {
	ID        int64  `json:"id"`
	Operation int32  `json:"operation"`
	Rate      int32  `json:"rate"`
	Tcount    int64  `json:"tcount"` //删除总数
	Start     string `json:"start"`
	End       string `json:"end"`
}

// TaskInfo dm task info
type TaskInfo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Creator  string `json:"creator"`
	Reviewer string `json:"reviewer"`
	State    int32  `json:"state"`
	Result   string `json:"result"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`
}

// TaskListArg .
type TaskListArg struct {
	Creator  string `form:"creator"`
	Reviewer string `form:"reviewer"`
	Title    string `form:"title"`
	State    int32  `form:"state" default:"-1"`
	Ctime    string `form:"ctime"`
	Pn       int64  `form:"pn" default:"1" validate:"gt=0"`
	Ps       int64  `form:"ps" default:"50" validate:"gt=0"`
}

// AddTaskArg .
type AddTaskArg struct {
	Creator   string
	Title     string `form:"title" validate:"required"`
	Regex     string `form:"regex"`
	KeyWords  string `form:"keywords"`
	IPs       string `form:"ips"`
	Mids      string `form:"mids"`
	Cids      string `form:"cids"`
	Start     string `form:"start" validate:"required"`
	End       string `form:"end" validate:"required"`
	State     int32  `form:"state" default:"0" validate:"gte=0"`
	Operation int32  `form:"operation" default:"-1" `
	OpTime    string `form:"operation_time"`
	OpRate    int32  `form:"operation_rate" default:"100" validate:"gt=0"`
}

// ReviewTaskArg .
type ReviewTaskArg struct {
	ID       int64 `form:"id" validate:"required,gte=0"`
	State    int32 `form:"state" validate:"required,gte=0"`
	Reviewer string
	Topic    string
}

// EditTasksStateArg .
type EditTasksStateArg struct {
	IDs   string `form:"ids" validate:"required"`
	State int32  `form:"state" validate:"required,gte=0"`
}

// TaskViewArg .
type TaskViewArg struct {
	ID int64 `form:"id" validate:"required,gte=0"`
}

// TaskCsvArg .
type TaskCsvArg struct {
	ID int64 `form:"id" validate:"required,gte=0"`
}
