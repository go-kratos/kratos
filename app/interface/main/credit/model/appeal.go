package model

import (
	"go-common/library/time"
)

// Appeal state
const (
	// StateCreate 用户刚创建申诉
	StateCreate = 1
	// StateReply 管理员回复,并且用户已读
	StateReply = 2
	// StateAdminClose 管理员关闭申诉
	StateAdminClose = 3
	// StateUserFinished 用户已解决申诉(评分)
	StateUserFinished = 4
	// StateTimeoutClose 超时关闭申诉
	StateTimeoutClose = 5
	// StateNoRead 管理员回复,用户未读
	StateNoRead = 6
	// StateUserClosed 用户直接关闭申诉
	StateUserClosed = 7
	// StateAdminFinished 管理员已通过申诉
	StateAdminFinished = 8

	// EventStateAdminReply  管理员回复
	EventStateAdminReply = 1
	// EventStateAdminNote 管理员回复并记录
	EventStateAdminNote = 2
	// EventStateUserReply 用户回复
	EventStateUserReply = 3
	// EventStateSystem  系统回复
	EventStateSystem = 4
	// appeal business
	Business = 5
)

// Appeal info.
type Appeal struct {
	ID          int64     `json:"id"`
	Oid         int64     `json:"oid"`
	Cid         int64     `json:"cid"`
	Mid         int64     `json:"mid"`
	Aid         int64     `json:"aid"`
	Tid         int8      `json:"tid"`
	Title       string    `json:"title"`
	State       int8      `json:"state"`
	Visit       int8      `json:"visit"`
	QQ          string    `json:"qq"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Pics        string    `json:"pics"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Star        int8      `json:"star"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// IsOpen appeal open state.
func IsOpen(state int8) bool {
	return state == StateCreate || state == StateReply || state == StateNoRead
}

// OpenedStates open get appeal
func OpenedStates() (states []int64) {
	return []int64{StateCreate, StateReply, StateNoRead}
}

// ClosedStates  get appeal
func ClosedStates() (states []int64) {
	return []int64{StateAdminClose, StateUserFinished, StateTimeoutClose, StateUserClosed, StateAdminFinished}
}
