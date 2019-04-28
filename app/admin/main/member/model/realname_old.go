package model

import (
	"time"
)

// DeDeIdentificationCardApplyImg is.
type DeDeIdentificationCardApplyImg struct {
	ID      int64     `gorm:"column:id"`
	IMGData string    `gorm:"column:img_data"`
	AddTime time.Time `gorm:"column:add_time"`
}

// TableName is
func (d *DeDeIdentificationCardApplyImg) TableName() string {
	return "dede_identification_card_apply_img"
}

// DeDeIdentificationCardApply is
type DeDeIdentificationCardApply struct {
	ID            int64  `gorm:"column:id"`
	MID           int64  `gorm:"column:mid"`
	Realname      string `gorm:"column:realname"`
	Type          int8   `gorm:"column:type"`
	CardData      string `gorm:"column:card_data"`
	CardForSearch string `gorm:"column:card_for_search"`
	FrontImg      int64  `gorm:"column:front_img"`
	FrontImg2     int64  `gorm:"column:front_img2"`
	BackImg       int64  `gorm:"column:back_img"`
	ApplyTime     int32  `gorm:"column:apply_time"`
	Operator      string `gorm:"column:operater"`
	OperatorTime  int32  `gorm:"column:operater_time"`
	Status        int8   `gorm:"column:status"`
	Remark        string `gorm:"column:remark"`
	RemarkStatus  int8   `gorm:"column:remark_status"`
}

// TableName is
func (d *DeDeIdentificationCardApply) TableName() string {
	return "dede_identification_card_apply"
}
