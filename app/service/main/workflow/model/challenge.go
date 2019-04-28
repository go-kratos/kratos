package model

import (
	"errors"

	"go-common/library/time"
)

const (
	// AuditRole 审核
	AuditRole = int8(0)
	// CustomerServiceRole 客服
	CustomerServiceRole = int8(1)

	// StateUntreated .
	StateUntreated = int8(0)
	// StatePassed .
	StatePassed = int8(1)
	// StateReject .
	StateReject = int8(2)
	// StateClose .
	StateClose = int8(3)

	// DispatchStateAuditMask 1111
	DispatchStateAuditMask = int32(0xf)
	// DispatchStateCustomerServiceMask 11110000
	DispatchStateCustomerServiceMask = int32(0xf0)

	// QueueState 队列中审核状态
	QueueState = 15
	// QueueBusinessState 队列中客服状态
	QueueBusinessState = 15

	// QueueStateBefore 默认审核状态
	QueueStateBefore = 0
	// QueueBusinessStateBefore 默认客服状态
	QueueBusinessStateBefore = 1
)

// Challenge struct
type Challenge struct {
	ID            int32         `gorm:"column:id" json:"id"`
	Tid           int32         `gorm:"column:tid" json:"tid"`
	Gid           int32         `gorm:"column:gid" json:"gid"`
	Oid           int64         `gorm:"column:oid" json:"oid"`
	Mid           int64         `gorm:"column:mid" json:"mid"`
	Eid           int64         `gorm:"column:eid" json:"eid"`
	State         int8          `gorm:"-" json:"state"`
	Business      int8          `gorm:"column:business" json:"business"`
	BusinessState int8          `gorm:"-" json:"business_state"`
	Assignee      int32         `gorm:"column:assignee_adminid" json:"assignee_adminid"`
	Adminid       int32         `gorm:"column:adminid" json:"adminid"`
	MetaData      string        `gorm:"column:metadata" json:"metadata"`
	Desc          string        `gorm:"column:description" json:"description"`
	Attachments   []*Attachment `gorm:"-" json:"attachments"`
	Events        []*Event      `gorm:"-" json:"events"`
	Ctime         time.Time     `gorm:"ctime" json:"ctime"`
	Mtime         time.Time     `gorm:"mtime" json:"mtime"`
	DispatchState uint32        `gorm:"dispatch_state"`
	BusinessInfo  Business      `gorm:"-" json:"business_info"`
}

// ChallengeParam appeal param
type ChallengeParam struct {
	ID              int32    `form:"id"`
	Tid             int32    `form:"tid"`
	Oid             int64    `form:"oid"`
	Mid             int64    `form:"mid"`
	Desc            string   `form:"description"`
	AdminID         int32    `form:"admin_id"`
	AssigneeID      int32    `form:"assignee_id"`
	AttachmentsStr  string   `form:"attachments"`
	Attachments     []string `form:"attachments[]"`
	Business        int8     `form:"business"`
	BusinessState   int8     `form:"business_state"`
	MetaData        string   `form:"metadata"`
	BusinessTypeid  int32    `form:"business_typeid"`
	BusinessTitle   string   `form:"business_title"`
	BusinessContent string   `form:"business_content"`
	BusinessMid     int64    `form:"business_mid"`
	BusinessExtra   string   `form:"business_extra"`
	Role            uint8    `form:"role"`
}

// TableName by Challenge
func (*Challenge) TableName() string {
	return "workflow_chall"
}

// SetDispatchState set DispatchState
func SetDispatchState(dispatchState int32, role, state int8) (result int32, err error) {
	switch role {
	case AuditRole:
		result = dispatchState&(^DispatchStateAuditMask) + int32(state)
	case CustomerServiceRole:
		result = dispatchState&(^DispatchStateCustomerServiceMask) + (int32(state) << 4)
	default:
		err = errors.New("changeDispatchState Unknown Role")
	}
	return result, err
}

// DispatchState get DispatchState
func DispatchState(dispatchState int32, role int8) (result int32, err error) {
	switch role {
	case AuditRole:
		result = dispatchState & DispatchStateAuditMask
	case CustomerServiceRole:
		result = (dispatchState & DispatchStateCustomerServiceMask) >> 4
	default:
		err = errors.New("changeDispatchState Unknown Role")
	}
	return result, err
}

// SetState update state of a role
// ex. oldState=0x3a4b5c6d, state=15, role=1 then result is 0x3a4b5cfd
func (c *Challenge) SetState(state uint32, role uint8) {
	oldState := c.DispatchState
	mod := uint32(^(0xf << (4 * role)))
	oldState = oldState & mod // all bit keep unchanged and bits you want update change to 0
	c.DispatchState = oldState + state<<(4*role)
}

// GetState return state of a role from dispatchState field
// ex. dispatchState=0x3a4b5c6d, role=1 then result is 0x6
func (c *Challenge) GetState(role uint8) (result int8) {
	dispatchState := c.DispatchState
	mod := uint32(0xf << (4 * role))
	dispatchState &= mod
	result = int8(dispatchState >> (4 * role))
	return
}

// FromState set State and BusinessState field from DispatchState field
func (c *Challenge) FromState() {
	c.State = c.GetState(uint8(0))
	if c.State == QueueState {
		c.State = QueueStateBefore
	}
	c.BusinessState = c.GetState(uint8(1))
	if c.BusinessState == QueueBusinessState {
		c.BusinessState = QueueBusinessStateBefore
	}
}

// CheckAdd check add challenge by params
func (ap *ChallengeParam) CheckAdd() bool {
	return !(ap.Oid == 0 || ap.Mid == 0 || ap.Business == 0 || ap.Tid == 0 || ap.Desc == "")
}

// CheckList check get list challenge by params
func (ap *ChallengeParam) CheckList() bool {
	return !(ap.Mid == 0 || ap.Business == 0)
}

// CheckInfo check get challenge info by params
func (ap *ChallengeParam) CheckInfo() bool {
	return !(ap.ID == 0 || ap.Mid == 0 || ap.Business == 0)
}

// CheckBusiness check challenge business field by params
func (ap *ChallengeParam) CheckBusiness() bool {
	return !(ap.BusinessTypeid == 0 && ap.BusinessMid == 0 && ap.BusinessTitle == "" && ap.BusinessContent == "" && ap.BusinessExtra == "")
}

// ChallengeParam3 .
type ChallengeParam3 struct {
	Business        int8     `form:"business" validate:"required"`
	Fid             int64    `form:"fid"`
	Rid             int64    `form:"rid"`
	Eid             int64    `form:"eid"`
	Score           int64    `form:"score"`
	Tid             int32    `form:"tid"`
	Oid             int64    `form:"oid" validate:"required"`
	Aid             int64    `form:"aid"`
	Mid             int64    `form:"mid"`
	Desc            string   `form:"description"`
	AdminID         int32    `form:"admin_id"`
	AssigneeID      int32    `form:"assignee_id"`
	Attachments     []string `form:"attachments,split"`
	BusinessState   int8     `form:"business_state"`
	MetaData        string   `form:"metadata"`
	BusinessTypeid  int16    `form:"business_typeid"`
	BusinessTitle   string   `form:"business_title"`
	BusinessContent string   `form:"business_content"`
	BusinessMid     int64    `form:"business_mid"`
	BusinessExtra   string   `form:"business_extra"`
	Role            uint8    `form:"role"`
}

// CheckBusiness .
func (cp3 *ChallengeParam3) CheckBusiness() bool {
	return !(cp3.BusinessTypeid == 0 && cp3.BusinessMid == 0 && cp3.BusinessTitle == "" && cp3.BusinessContent == "" && cp3.BusinessExtra == "")
}

// Challenge3 .
type Challenge3 struct {
	ID            int64     `json:"id" gorm:"column:id"`
	Gid           int64     `json:"gid" gorm:"column:gid"`
	Mid           int64     `json:"mid" gorm:"column:mid"`
	Tid           int64     `json:"tid" gorm:"column:tid"`
	Eid           int64     `json:"eid" gorm:"column:eid"`
	Oid           int64     `json:"oid" gorm:"column:oid"`
	Business      int64     `json:"business" gorm:"column:business"`
	Desc          string    `json:"description" gorm:"column:description"`
	MetaData      string    `json:"metadata" gorm:"column:metadata"`
	DispatchState int64     `json:"dispatch_state" gorm:"column:dispatch_state"`
	Ctime         time.Time `json:"ctime" gorm:"column:ctime"`
	Mtime         time.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (*Challenge3) TableName() string {
	return "workflow_chall"
}
