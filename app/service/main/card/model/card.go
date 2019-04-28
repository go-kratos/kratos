package model

import "go-common/library/time"

// Card info.
type Card struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	State        int32     `json:"state"`
	Deleted      int32     `json:"deleted"`
	IsHot        int32     `json:"is_hot"`
	CardURL      string    `json:"card_url"`
	BigCradURL   string    `json:"big_card_url"`
	CardType     int32     `json:"card_type"`
	CardTypeName string    `json:"card_type_name"`
	OrderNum     int64     `json:"order_num"`
	GroupID      int64     `json:"group_id"`
	Operator     string    `json:"operator"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// UserCard user card info.
type UserCard struct {
	Mid          int64  `json:"mid"`
	ID           int64  `json:"id"`
	CardURL      string `json:"card_url"`
	BigCradURL   string `json:"big_card_url"`
	CardType     int32  `json:"card_type"`
	Name         string `json:"name"`
	ExpireTime   int64  `json:"expire_time"`
	CardTypeName string `json:"card_type_name"`
}
