package audit

import (
	"time"
)

// Audit audit
type Audit struct {
	ID      int64     `json:"id"`
	Build   int       `json:"build"`
	Remark  string    `json:"remark"`
	MobiApp string    `json:"mobi_app"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// Param param
type Param struct {
	ID      int64  `form:"id"`
	Build   int    `form:"build"`
	Remark  string `form:"remark"`
	MobiApp string `form:"mobi_app"`
}

// TableName return table name
func (*Audit) TableName() string {
	return "audit"
}
