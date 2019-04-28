package model

import xtime "go-common/library/time"

// Record .
type Record struct {
	ID       int        `json:"id" gorm:"column:id"`
	Bucket   string     `json:"bucket" gorm:"column:bucket"`
	FileName string     `json:"filename" gorm:"column:filename"`
	AdminID  int        `json:"admin_id" gorm:"column:adminid"`
	State    int        `json:"state" gorm:"column:state"`
	CTime    xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime    xtime.Time `json:"mtime" gorm:"column:mtime"`
	URL      string     `json:"url" gorm:"url"`
	Sex      int        `json:"sex" gorm:"sex"`
	Politics int        `json:"politics" gorm:"politics"`
}

// TableName .
func (Record) TableName() string {
	return "upload_yellowing"
}

// TinyRecord .
type TinyRecord struct {
	Rid int `gorm:"column:id"`
}

// TableName .
func (TinyRecord) TableName() string {
	return "upload_yellowing"
}
