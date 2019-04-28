package model

import xtime "go-common/library/time"

// BusinessAttr will record business attributes
type BusinessAttr struct {
	ID           int64      `json:"id" gorm:"column:id"`
	BID          int64      `json:"bid" gorm:"column:bid"`
	BusinessName string     `json:"business_name" gorm:"business_name"`
	Name         string     `json:"name" gorm:"column:name"`
	DealType     int8       `json:"deal_type" gorm:"column:deal_type"`
	ExpireTime   int64      `json:"expire_time" gorm:"column:expire_time"`
	AssignType   int8       `json:"assign_type" gorm:"column:assign_type"`
	AssignMax    int8       `json:"assign_max" gorm:"column:assign_max"`
	GroupType    int8       `json:"group_type" gorm:"column:group_type"`
	Button       uint8      `json:"-" gorm:"button"`
	ButtonKey    string     `json:"-" gorm:"button_key"`
	CTime        xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime        xtime.Time `json:"mtime" gorm:"column:mtime"`
	Buttons      []*Button  `json:"button" gorm:"-"`
}

// TableName is used to identify chall table name in gorm
func (BusinessAttr) TableName() string {
	return "workflow_business_attr"
}

// Button .
type Button struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	State bool   `json:"state"`
	Key   string `json:"key"`
}
