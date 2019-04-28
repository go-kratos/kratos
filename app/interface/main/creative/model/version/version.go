package version

import (
	xtime "go-common/library/time"
)

// Version str.
type Version struct {
	ID       int64      `json:"id"`
	Ty       string     `json:"type"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Link     string     `json:"link"`
	Ctime    xtime.Time `json:"-"`
	Dateline xtime.Time `json:"dateline"`
}

// FullVersions fn 4=>Android;5:iPhone;6:PC
func FullVersions() (tys []int) {
	return []int{4, 5, 6}
}
