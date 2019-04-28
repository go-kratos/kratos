package anchorTask

import (
	xtime "go-common/library/time"
)

// TableName is used to identify table name in gorm
func (ar *AnchorReward) TableName() string {
	return "ap_anchor_task_reward_list"
}

// consts .
const (
	RewardUnUsed  = int64(1)
	RewardUsed    = int64(3)
	RewardExpired = int64(5)

	CountExpireLockKey = "xrc:lock:count_expire:lock:v1"
	CountExpireUserKey = "xrc:lock:count_expire:user:v1:%d"
	SetExpireLockKey   = "xrc:lock:set_expire:lock:v1"
	ExpireCountTime    = 86400 * 30 * 3

	RewardExists = int64(1)
)

// AnchorReward .
type AnchorReward struct {
	Id          int64      `json:"id" gorm:"comumn:id"`
	Uid         int64      `json:"uid" gorm:"comumn:uid"`
	RewardId    int64      `json:"reward_id" gorm:"comumn:reward_id"`
	Roomid      int64      `json:"roomid" gorm:"comumn:roomid"`
	Lid         int64      `json:"lid" gorm:"comumn:lid"`
	Source      int64      `json:"source" gorm:"comumn:source"`
	UsePlat     int64      `json:"use_plat" gorm:"comumn:use_plat"`
	AchieveTime xtime.Time `json:"achieve_time" gorm:"comumn:achieve_time"`
	UseTime     xtime.Time `json:"use_time" gorm:"comumn:use_time;not null;"`
	StartTime   xtime.Time `json:"start_time" gorm:"comumn:start_time"`
	EndTime     xtime.Time `json:"end_time" gorm:"comumn:end_time"`
	ExpireTime  xtime.Time `json:"expire_time" gorm:"comumnw:expire_time"`
	Status      int64      `json:"status" gorm:"comumn:status"`
	Reserved1   int64      `json:"reserved1" gorm:"comumn:reserved1"`
	Reserved2   string     `json:"reserved2" gorm:"comumn:reserved2"`
	Ctime       xtime.Time `json:"ctime" gorm:"comumn:ctime"`
	Mtime       xtime.Time `json:"mtime" gorm:"comumn:mtime"`
}

// AnchorRewardObject .
type AnchorRewardObject struct {
	// id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// 奖励类型 1:ss推荐卡 2:s推荐卡、任意门
	RewardType int64 `protobuf:"varint,2,opt,name=reward_type,json=rewardType,proto3" json:"reward_type,omitempty"`
	// 1:未使用,3:已使用,5:已过期
	Status int64 `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	// 奖励id
	RewardId int64 `protobuf:"varint,4,opt,name=reward_id,json=rewardId,proto3" json:"reward_id,omitempty"`
	// 奖励名称
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// 奖励图标
	Icon string `protobuf:"bytes,6,opt,name=icon,proto3" json:"icon,omitempty"`
	// 获得时间datetime
	AchieveTime string `protobuf:"bytes,7,opt,name=achieve_time,json=achieveTime,proto3" json:"achieve_time,omitempty"`
	// 过期时间datetime
	ExpireTime string `protobuf:"bytes,8,opt,name=expire_time,json=expireTime,proto3" json:"expire_time,omitempty"`
	// 过期时间datetime
	UseTime string `protobuf:"bytes,8,opt,name=expire_time,json=useTime,proto3" json:"use_time,omitempty"`
	// 来源,1:主播任务,2:小时榜
	Source int64 `protobuf:"varint,9,opt,name=source,proto3" json:"source,omitempty"`
	// 奖励简介
	RewardIntro string `protobuf:"bytes,10,opt,name=reward_intro,json=rewardIntro,proto3" json:"reward_intro,omitempty"`
}

// AnchorRewardPager .
type AnchorRewardPager struct {
	// 当前页码
	Page int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// 每页大小
	PageSize int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// 总页数
	TotalPage int64 `protobuf:"varint,3,opt,name=total_page,json=totalPage,proto3" json:"total_page,omitempty"`
	// 总记录数
	TotalCount int64 `protobuf:"varint,4,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}
