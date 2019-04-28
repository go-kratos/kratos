package model

import "time"

// Card info.
type Card struct {
	ID         int64     `json:"id" gorm:"primary_key"`
	Name       string    `json:"name" gorm:"column:name"`
	State      int32     `json:"state" gorm:"column:state"`
	Deleted    int32     `json:"deleted" gorm:"column:deleted"`
	IsHot      int32     `json:"is_hot" gorm:"column:is_hot"`
	CardURL    string    `json:"card_url" gorm:"column:card_url"`
	BigCradURL string    `json:"big_crad_url" gorm:"column:big_crad_url"`
	CardType   int32     `json:"card_type" gorm:"column:card_type"`
	OrderNum   int64     `json:"order_num" gorm:"column:order_num"`
	GroupID    int64     `json:"group_id" gorm:"column:group_id"`
	Operator   string    `json:"operator" gorm:"column:operator"`
	Ctime      time.Time `json:"-" gorm:"-"`
	Mtime      time.Time `json:"-" gorm:"-"`
}

// CardGroup card group info.
type CardGroup struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	Name     string    `json:"name" gorm:"column:name"`
	State    int8      `json:"state" gorm:"column:state"`
	Deleted  int8      `json:"deleted" gorm:"column:deleted"`
	Operator string    `json:"operator" gorm:"column:operator"`
	OrderNum int64     `json:"order_num" gorm:"column:order_num"`
	Ctime    time.Time `json:"-" gorm:"-"`
	Mtime    time.Time `json:"-" gorm:"-"`
	Cards    []*Card   `json:"cards,omitempty" gorm:"-"`
}
