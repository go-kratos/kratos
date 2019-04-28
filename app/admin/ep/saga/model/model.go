package model

// pagination.
const (
	DefaultPageSize = 10
	DefaultPageNum  = 1
)

// Pagination Pagination.
type Pagination struct {
	PageSize int `form:"page_size"`
	PageNum  int `form:"page_num"`
}

// TeamParam struct for organization Info.
type TeamParam struct {
	Department string `form:"department"`
	Business   string `form:"business"`
}

// EmptyReq params for request without params
type EmptyReq struct {
}

// EmptyResp resp for response without data
type EmptyResp struct {
}

// Task def
type Task struct {
	ID           int64  `json:"id,omitempty" gorm:"column:id"`
	ProjID       int    `json:"proj_id,omitempty" gorm:"column:proj_id"`
	EventType    string `json:"event_type,omitempty" gorm:"column:event_type"`
	Author       string `json:"author,omitempty" gorm:"author"`
	MRID         int    `json:"mr_id,omitempty" gorm:"column:mr_id"`
	URL          string `json:"url,omitempty" gorm:"column:url"`
	Status       int    `json:"status,omitempty" gorm:"status"`
	TaskDetails  string `json:"task_details,omitempty" gorm:"task_details"`
	SourceBranch string `json:"source_branch,omitempty" gorm:"source_branch"`
	TargetBranch string `json:"target_branch,omitempty" gorm:"target_branch"`
	Title        string `json:"title,omitempty" gorm:"title"`
}

// TasksReq params for tasks
type TasksReq struct {
	ProjID   int   `form:"proj_id" validate:"required"`
	Statuses []int `form:"statuses,split" default:"3,4"` // 3 - running, 4 - waiting, 默认查运行中和等待的任务
}

// TasksResp resp for tasks
type TasksResp struct {
	Tasks []*Task `json:"tasks,omitempty"`
}

// User User.
type User struct {
	Name  string `json:"username" gorm:"column:name"`
	EMail string `json:"email" gorm:"column:email"`
}

// RequireVisibleUser def
type RequireVisibleUser struct {
	UserName string
	NickName string
}

// ContactInfo def
type ContactInfo struct {
	ID          string `json:"id,omitempty" gorm:"column:id"`
	UserName    string `json:"english_name" gorm:"column:user_name"`
	UserID      string `json:"userid" gorm:"column:user_id"`
	NickName    string `json:"name" gorm:"column:nick_name"`
	VisibleSaga bool   `json:"visible_saga" gorm:"column:visible_saga"`
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
