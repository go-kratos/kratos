package model

import xtime "go-common/library/time"

const (
	// BusinessOpenState .
	BusinessOpenState = 1
	// BusinessOpenType .
	BusinessOpenType = 1
	// UserRoleDefaultVal .
	UserRoleDefaultVal = -1
	// UserOnState .
	UserOnState = 1
	// UserOffState .
	UserOffState = 0
)

// Business .
type Business struct {
	ID        int64      `json:"id" gorm:"primary_key" form:"id"`
	PID       int64      `json:"pid" gorm:"column:pid" form:"pid"`
	BID       int64      `json:"bid" gorm:"column:bid" form:"bid"`
	Name      string     `json:"name" gorm:"column:name" form:"name"`
	Names     []string   `json:"names" gorm:"-" form:"names,split"`
	Flow      int64      `json:"flow" gorm:"column:flow" form:"flow" default:"0"`
	FlowState int64      `json:"flow_state" gorm:"column:flow_state" form:"flow_state"`
	State     int64      `json:"state" gorm:"column:state" form:"state" default:"1"`
	Ctime     xtime.Time `json:"ctime" gorm:"-"`
	Mtime     xtime.Time `json:"mtime" gorm:"-"`
}

// BusinessList .
type BusinessList struct {
	ID        int64           `json:"id" gorm:"primary_key"`
	PID       int64           `json:"pid" gorm:"column:pid"`
	BID       int64           `json:"bid" gorm:"column:bid"`
	Name      string          `json:"name" gorm:"column:name"`
	Flow      int64           `json:"flow" gorm:"column:flow"`
	FlowState int64           `json:"flow_state" gorm:"column:flow_state"`
	State     int64           `json:"state" gorm:"column:state"`
	Ctime     xtime.Time      `json:"ctime" gorm:"-"`
	Mtime     xtime.Time      `json:"mtime" gorm:"-"`
	FlowChild []*FlowBusiness `json:"flowchild" gorm:"-"`
}

// TableName .
func (bl *BusinessList) TableName() string {
	return "manager_business"
}

// FlowBusiness .
type FlowBusiness struct {
	FlowState int64           `json:"flow_state"`
	Child     []*BusinessList `json:"child"`
}

// BusinessRole .
type BusinessRole struct {
	ID    int64      `json:"id" gorm:"primary_key" form:"id"`
	BID   int64      `json:"bid" gorm:"column:bid" form:"bid"`
	RID   int64      `json:"rid" gorm:"column:rid" form:"rid"`
	Name  string     `json:"name" gorm:"column:name" form:"name"`
	Type  int64      `json:"type" gorm:"column:type" form:"type" default:"-1"`
	State int64      `json:"state" gorm:"column:state" form:"state" default:"1"`
	Ctime xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime xtime.Time `json:"mtime" gorm:"column:mtime"`
	Auth  int64      `json:"auth" gorm:"-" form:"auth"`
	UID   int64      `json:"uid" gorm:"-"`
}

// TableName .
func (br *BusinessRole) TableName() string {
	return "manager_business_role"
}

// BusinessUserRole .
type BusinessUserRole struct {
	ID    int64      `json:"id" gorm:"primary_key" form:"id"`
	UID   int64      `json:"uid" gorm:"column:uid" form:"uid"`
	UIDs  []int64    `json:"uids" gorm:"-" form:"uids,split"`
	CUID  int64      `json:"cuid" gorm:"column:cuid" form:"cuid"`
	BID   int64      `json:"bid" gorm:"column:bid" form:"bid"`
	Role  string     `json:"role" gorm:"column:role" form:"role"`
	Ctime xtime.Time `json:"ctime" gorm:"-"`
	Mtime xtime.Time `json:"mtime" gorm:"-"`
}

// BusinessUserRoleList .
type BusinessUserRoleList struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	UID       int64      `json:"uid" gorm:"column:uid" form:"uid"`
	UName     string     `json:"uname" gorm:"-"`
	UNickname string     `json:"unickname" gorm:"-"`
	CUID      int64      `json:"cuid" gorm:"column:cuid" form:"cuid"`
	CUName    string     `json:"cuname" gorm:"-"`
	BID       int64      `json:"bid" gorm:"column:bid" form:"bid"`
	Role      string     `json:"role" gorm:"column:role" form:"role"`
	RoleName  []string   `json:"rolename" gorm:"-"`
	Ctime     xtime.Time `json:"ctime" gorm:"-"`
	Mtime     xtime.Time `json:"mtime" gorm:"-"`
}

// StateUpdate .
type StateUpdate struct {
	ID    int64 `json:"id" form:"id" validate:"required"`
	Type  int64 `json:"type" form:"type" validate:"required"` // 1:业务 2:角色
	State int64 `json:"state" form:"state"`
}

// UserListParams .
type UserListParams struct {
	BID   int64  `form:"bid" validate:"required"`
	Role  int64  `form:"role" default:"-1"`
	UName string `form:"uname"`
	UID   int64  `form:"uid"`
	PS    int64  `form:"ps" default:"50"`
	PN    int64  `form:"pn" default:"1"`
}

// UserRole .
type UserRole struct {
	ID   int64  `json:"id"`
	BID  int64  `json:"bid"`
	RID  int64  `json:"rid"`
	Type int64  `json:"type"`
	Name string `json:"name"`
}

// BusinessListParams .
type BusinessListParams struct {
	State int64 `json:"state" form:"state" default:"-1"`
	Check int64 `json:"check" form:"check" default:"0"`
	Flow  int64 `json:"flow" form:"flow"`
	UID   int64 `json:"uid"`
	Auth  int64 `json:"auth" form:"auth"`
}

// UserStateUp .
type UserStateUp struct {
	BID     int64 `json:"bid" form:"bid" validate:"required"`
	AdminID int64 `json:"adminid" form:"adminid" validate:"required"`
	State   int   `json:"state" form:"state" validate:"state"`
}
