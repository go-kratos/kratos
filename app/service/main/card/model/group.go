package model

import (
	"go-common/library/time"
)

// CardGroup card group info.
type CardGroup struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	State    int8      `json:"state"`
	Deleted  int8      `json:"deleted"`
	Operator string    `json:"operator"`
	OrderNum int64     `json:"order_num"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// AllGroupResp all group resp.
type AllGroupResp struct {
	List     []*GroupInfo `json:"list"`
	UserCard *UserCard    `json:"user_card,omitempty"`
}

// GroupInfo group info
type GroupInfo struct {
	GroupID   int64   `json:"group_id"`
	GroupName string  `json:"group_name"`
	Cards     []*Card `json:"cards"`
	OrderNum  int64   `json:"-"`
}
