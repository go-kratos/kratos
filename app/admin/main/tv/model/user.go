package model

import "go-common/library/time"

// TvUserInfo is table struct
type TvUserInfo struct {
	ID            int64     `json:"id"`
	MID           int64     `json:"mid" gorm:"column:mid"`
	Ver           int64     `json:"ver"`
	VipType       int8      `json:"vip_type"`
	PayType       int8      `json:"pay_type"`
	PayChannelID  string    `json:"pay_channel_id"`
	Status        int8      `json:"status"`
	OverdueTime   time.Time `json:"overdue_time"`
	RecentPayTime time.Time `json:"recent_pay_time"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

// TvUserInfoResp is used to user info
type TvUserInfoResp struct {
	ID            int64     `json:"id"`
	MID           int64     `json:"mid" gorm:"column:mid"`
	VipType       int8      `json:"vip_type"`
	PayType       int8      `json:"pay_type"`
	PayChannelID  string    `json:"pay_channel_id"`
	Status        int8      `json:"status"`
	OverdueTime   time.Time `json:"overdue_time"`
	RecentPayTime time.Time `json:"recent_pay_time"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

// TvUserChangeHistory is table struct
type TvUserChangeHistory struct {
	ID         int64     `json:"id"`
	MID        int64     `json:"mid"`
	ChangeType int8      `json:"change_type"`
	ChangeTime time.Time `json:"change_time"`
	Days       int64     `json:"days"`
	OperatorId string    `json:"operator_id"`
	Remark     string    `json:"remark"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// TableName tv_user_info
func (t *TvUserInfo) TableName() string {
	return "tv_user_info"
}

// TableName tv_user_info
func (t *TvUserInfoResp) TableName() string {
	return "tv_user_info"
}
