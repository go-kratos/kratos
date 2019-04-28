package codes

import (
	xtime "go-common/library/time"
)

// code status
const (
	StatusOpen = int8(1)
	Type       = int8(1)
)

// TableName case tablename
func (*Codes) TableName() string {
	return "codes"
}

// Codes codes
type Codes struct {
	ID          int64      `gorm:"column:id" json:"id"`
	Code        int32      `gorm:"column:code" json:"code"`
	Message     string     `gorm:"column:message" json:"message"`
	Operator    string     `gorm:"column:operator" json:"operator"`
	CTime       xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime       xtime.Time `gorm:"column:mtime" json:"mtime"`
	HantMessage string     `gorm:"column:hant_message" json:"hant_message"`
	Level       int8       `gorm:"column:level" json:"level"`
}

// ResultCodes ...
type ResultCodes struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    []*Codes
}

// NewCodes ...
type NewCodes struct {
	*Codes
	List []*CodeMsg `json:"list"`
}
