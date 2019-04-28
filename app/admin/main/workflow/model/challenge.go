package model

import (
	"net/url"

	xtime "go-common/library/time"
)

// Chall is the row view for every challenge
type Chall struct {
	Cid               int64      `json:"cid" gorm:"column:id"`
	Gid               int64      `json:"gid" gorm:"column:gid"`
	Oid               int64      `json:"oid" gorm:"column:oid"`
	OidStr            string     `json:"oid_str" gorm:"-"`
	Business          int8       `json:"business" gorm:"column:business"`
	Mid               int64      `json:"mid" gorm:"column:mid"`
	MName             string     `json:"m_name" gorm:"-"`
	Tid               int64      `json:"tid" gorm:"column:tid"`
	State             int8       `json:"state"`
	BusinessState     int8       `json:"business_state"`
	DispatchState     uint32     `json:"-" gorm:"column:dispatch_state"`
	DispatchTime      xtime.Time `json:"dispatch_time" gorm:"column:dispatch_time"`
	Description       string     `json:"description" gorm:"column:description"`
	Metadata          string     `json:"metadata" gorm:"column:metadata"`
	CTime             xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime             xtime.Time `json:"mtime" gorm:"column:mtime"`
	BusinessObject    *Business  `json:"business_object,omitempty" gorm:"-"`
	AssigneeAdminID   int64      `json:"assignee_adminid" gorm:"column:assignee_adminid"`
	AdminID           int64      `json:"adminid" gorm:"column:adminid"`
	AssigneeAdminName string     `json:"assignee_admin_name" gorm:"-"`
	AdminName         string     `json:"admin_name" gorm:"-"`
	TotalStates       string     `json:"total_states" gorm:"-"`

	// tag related fields
	Tag   string `json:"tag" gorm:"-"`
	Round int8   `json:"round" gorm:"-"`

	// log related
	LastLog string `json:"last_log" gorm:"-"`
	// event related
	LastEvent *Event `json:"last_event" gorm:"-"`

	// Attachments
	Attachments []string `json:"attachments" gorm:"-"`

	// linked group object
	Group *Group `json:"group,omitempty" gorm:"-"`

	Meta     interface{} `json:"meta" gorm:"-"`
	AuditLog interface{} `json:"audit_log" gorm:"-"`
	Producer *Account    `json:"producer" gorm:"-"` //举报人

	// business table
	Title  string `json:"title,omitempty" gorm:"-"`
	TypeID int64  `json:"type_id,omitempty" gorm:"-"`
}

// TinyChall is the tiny row view for every challenge
type TinyChall struct {
	Cid   int64      `json:"cid" gorm:"column:id"`
	Gid   int64      `json:"gid" gorm:"column:gid"`
	Mid   int64      `json:"mid" gorm:"column:mid"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	State int8       `json:"state" gorm:"-"`
	Title string     `json:"title" gorm:"-"`
}

// ChallTagSlice is the slice to ChallTag
type ChallTagSlice []*ChallTag

// ChallTag is the model to retrive user submitted tags in group view
type ChallTag struct {
	ID      int64  `json:"id"`
	Tag     string `json:"tag"`
	Round   int8   `json:"round"`
	Count   int64  `json:"count"`
	Percent int32  `json:"percent"`
}

// TableName is used to identify chall table name in gorm
func (Chall) TableName() string {
	return "workflow_chall"
}

func (c ChallTagSlice) Len() int {
	return len(c)
}

func (c ChallTagSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ChallTagSlice) Less(i, j int) bool {
	return c[i].Percent < c[j].Percent
}

// FixAttachments will fix attachments url as user friendly
// ignore https case
// FIXME: this should be removed after attachment url is be normed
func (c *Chall) FixAttachments() {
	if c.Attachments == nil {
		return
	}
	fixed := make([]string, 0, len(c.Attachments))
	for _, a := range c.Attachments {
		u, err := url.Parse(a)
		if err != nil {
			continue
		}
		u.Scheme = "http"
		fixed = append(fixed, u.String())
	}
	c.Attachments = fixed
}

// SetState update state of a role
// ex. oldState=0x3a4b5c6d, state=15, role=1 then result is 0x3a4b5cfd
func (c *Chall) SetState(state uint32, role uint8) {
	oldState := c.DispatchState
	mod := uint32(^(0xf << (4 * role)))
	oldState = oldState & mod // all bit keep unchanged and bits you want update change to 0
	c.DispatchState = oldState + state<<(4*role)
}

// getState return state of a role from dispatchState field
// ex. dispatchState=0x3a4b5c6d, role=1 then result is 0x6
func (c *Chall) getState(role uint8) (result int8) {
	dispatchState := c.DispatchState
	mod := uint32(0xf << (4 * role))
	dispatchState &= mod
	result = int8(dispatchState >> (4 * role))
	return
}

// FromState set State and BusinessState field from DispatchState field
func (c *Chall) FromState() {
	c.State = c.getState(uint8(0))
	c.BusinessState = c.getState(uint8(1))
}

// FormatState transform state in queue
func (c *Chall) FormatState() {
	if c.State == QueueState {
		c.State = QueueStateBefore
	}

	if c.BusinessState == QueueBusinessState {
		c.BusinessState = QueueBusinessStateBefore
	}
}
