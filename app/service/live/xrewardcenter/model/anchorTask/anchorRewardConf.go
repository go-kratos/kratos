package anchorTask

import (
	xtime "go-common/library/time"
)

// TableName is used to identify table name in gorm
func (arc *AnchorRewardConf) TableName() string {
	return "ap_anchor_task_reward_conf"
}

// AnchorRewardConf .
type AnchorRewardConf struct {
	ID           int64      `json:"id" gorm:"comumn:id"`
	Name         string     `json:"name" gorm:"comumn:name"`
	Icon         string     `json:"icon" gorm:"comumn:icon"`
	RewardIntro  string     `json:"reward_intro" gorm:"comumn:reward_intro"`
	RewardType   int64      `json:"reward_type" gorm:"comumn:reward_type"`
	Detail       string     `json:"detail" gorm:"comumn:detail"`
	Instructions string     `json:"instructions" gorm:"comumn:instructions"`
	Reserved1    int64      `json:"reserved1" gorm:"comumn:reserved1"`
	Reserved2    string     `json:"reserved2" gorm:"comumn:reserved2"`
	Ctime        xtime.Time `json:"ctime" gorm:"comumn:ctime"`
	Mtime        xtime.Time `json:"mtime" gorm:"comumn:mtime"`
}
