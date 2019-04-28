package model

import "go-common/library/time"

// LimitUser limit user.
type LimitUser struct {
	ID      int64     `json:"id"`
	Mid     int64     `json:"mid"`     //用户mid
	Name    string    `json:"name"`    //用户名称
	Creator string    `json:"creator"` //创建者名称
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"-"`
}

// LimitRes Limit Res.
type LimitRes struct {
	ID        int64     `json:"id"`
	Oid       int64     `json:"oid"`
	Type      int64     `json:"type"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Operation int32     `json:"operation"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}
