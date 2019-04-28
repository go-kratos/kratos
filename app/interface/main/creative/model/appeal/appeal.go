package appeal

import (
	"go-common/app/service/main/archive/api"
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
	// Business appeal business
	Business = 2
	// ReplyMsg appeal auto reply msg
	ReplyMsg = "您好，您的反馈我们已收到，会尽快核实处理，请您稍等。"
	//ReplyEvent 1：管理员回复；2：管理员备注；3：用户回复；4：系统回复
	ReplyUserEvent   = 3
	ReplySystemEvent = 4
)

// AppealMeta for appeal detail.
type AppealMeta struct {
	ID            int64         `gorm:"column:id" json:"id"`
	Tid           int32         `gorm:"column:tid" json:"tid"`
	Gid           int32         `gorm:"column:gid" json:"gid"`
	Oid           int64         `gorm:"column:oid" json:"oid"`
	Mid           int64         `gorm:"column:mid" json:"mid"`
	State         int8          `gorm:"column:state" json:"state"`
	Business      int8          `gorm:"column:business" json:"business"`
	BusinessState int8          `gorm:"column:business_state" json:"business_state"`
	Assignee      int32         `gorm:"column:assignee_adminid" json:"assignee_adminid"`
	Adminid       int32         `gorm:"column:adminid" json:"adminid"`
	MetaData      string        `gorm:"column:metadata" json:"metadata"`
	Desc          string        `gorm:"column:description" json:"description"`
	Attachments   []*Attachment `gorm:"-" json:"attachments"`
	Events        []*EventNew   `gorm:"-" json:"events"`
	CTime         time.Time     `json:"ctime"`
	MTime         time.Time     `json:"mtime"`
}

//EventNew for new.
type EventNew struct {
	ID          int64     `gorm:"column:id" json:"id"`
	Cid         int64     `gorm:"column:cid" json:"cid"`
	Event       int64     `gorm:"column:event" json:"event"`
	Adminid     int64     `gorm:"column:adminid" json:"adminid"`
	Content     string    `gorm:"column:content" json:"content"`
	Attachments string    `gorm:"column:attachments" json:"attachments"`
	CTime       time.Time `gorm:"column:ctime" json:"ctime"`
	MTime       time.Time `gorm:"column:mtime" json:"mtime"`
}

// Appeal info.
type Appeal struct {
	ID          int64         `json:"id"`
	Oid         int64         `json:"oid"`
	Cid         int64         `json:"cid"`
	Mid         int64         `json:"mid"`
	Aid         int64         `json:"aid"`
	Tid         int8          `json:"tid"`
	Title       string        `json:"title"`
	State       int8          `json:"state"`
	Visit       int8          `json:"visit"`
	QQ          string        `json:"qq"`
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	Pics        string        `json:"pics"`
	Content     string        `json:"content"`
	Description string        `json:"description"`
	Star        int8          `json:"star"`
	CTime       time.Time     `json:"ctime"`
	MTime       time.Time     `json:"mtime"`
	Attachments []*Attachment `json:"attachments"`
	// event
	Events []*Event `json:"events"`
	// archive
	Archive  *api.Arc  `json:"archive,omitempty"`
	UserInfo *UserInfo `json:"userinfo"`
}

type UserInfo struct {
	MID   int64  `json:"mid"`
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Face  string `json:"face"`
	Rank  int32  `json:"rank"`
	Level int32  `json:"level"`
}

// Event appeal work order deal.
type Event struct {
	ID          int64     `json:"id"`
	AdminID     int64     `json:"adminid"`
	Content     string    `json:"content"`
	ApID        int64     `json:"apid"`
	Pics        string    `json:"pics"`
	Event       int64     `json:"event"`
	Attachments string    `json:"attachments"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// Attachment is appeal attachment.
type Attachment struct {
	ID   int64  `json:"id"`
	Cid  int64  `json:"cid"`
	Path string `json:"path"`
}

// Contact user contacts.
type Contact struct {
	Mid      int64  `json:"mid"`
	Uname    string `json:"uname"`
	TelPhone string `json:"telPhone"`
	Email    string `json:"email"`
}

// BusinessAppeal for new arc add appeal.
type BusinessAppeal struct {
	BusinessTypeID  int64  `json:"business_typeid"`
	BusinessMID     int64  `json:"business_mid"`
	BusinessTitle   string `json:"business_title"`
	BusinessContent string `json:"business_content	"`
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

// IsClosed appeal is close.
func IsClosed(state int8) (is bool) {
	if state == StateAdminClose || state == StateUserFinished || state == StateTimeoutClose || state == StateUserClosed || state == StateAdminFinished {
		is = true
	}
	return
}

// Allow archive state in (-2,-4) can add appeal.
func Allow(state int8) bool {
	return state == -2 || state == -4
}
