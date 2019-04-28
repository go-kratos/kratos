package pgc

import "go-common/library/time"

// Content content def.
type Content struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Desc      string    `json:"desc"`
	Cover     string    `json:"cover"`
	EPID      int       `json:"epid"`
	CID       int       `json:"cid"`
	MenuID    int       `json:"menu_id"`
	SeasonID  int       `json:"season_id"`
	State     int       `json:"state"`
	Valid     int       `json:"valid"`
	PayStatus int       `json:"pay_status"`
	IsDeleted int       `json:"is_deleted"`
	AuditTime int       `json:"audit_time"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TableName tv_content
func (c Content) TableName() string {
	return "tv_content"
}
