package model

import (
	"time"
)

// ActBws def.
type ActBws struct {
	ID    int64     `json:"id" form:"id"`
	Name  string    `json:"name" form:"name"`
	Image string    `json:"image" form:"image"`
	Dic   string    `json:"dic" form:"dic"`
	Del   int8      `json:"del" form:"del"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// ActBwsAchievement def.
type ActBwsAchievement struct {
	ID            int64     `json:"id" form:"id"`
	Name          string    `json:"name" form:"name"`
	Icon          string    `json:"icon" form:"icon"`
	Dic           string    `json:"dic" form:"dic"`
	Image         string    `json:"image" form:"image"`
	LinkType      int64     `json:"link_type" form:"link_type"`
	Unlock        int64     `json:"unlock" form:"unlock"`
	BID           int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	IconBig       string    `json:"icon_big" form:"icon_big"`
	IconActive    string    `json:"icon_active" form:"icon_active"`
	IconActiveBig string    `json:"icon_active_big" form:"icon_active_big"`
	Award         int8      `json:"award" form:"award"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
	Del           int8      `json:"del"  form:"del"`
	SuitID        int64     `json:"suit_id" gorm:"column:suit_id"  form:"suit_id"`
}

// ActBwsField def.
type ActBwsField struct {
	ID    int64     `json:"id" form:"id"`
	Name  string    `json:"name" form:"name"`
	Area  string    `json:"area" form:"area"`
	BID   int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	Del   int8      `json:"del" form:"del"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// ActBwsPoint def.
type ActBwsPoint struct {
	ID           int64     `json:"id" form:"id"`
	Name         string    `json:"name" form:"name"`
	Icon         string    `json:"icon" form:"icon"`
	FID          int64     `json:"fid" gorm:"column:fid"  form:"fid"`
	Ower         int64     `json:"ower" gorm:"column:ower"  form:"ower"`
	Image        string    `json:"image" form:"image"`
	Unlocked     int64     `json:"unlocked" form:"unlocked"`
	LoseUnlocked int64     `json:"lose_unlocked" form:"lose_unlocked"`
	LockType     int64     `json:"lock_type" form:"lock_type"`
	Dic          string    `json:"dic" form:"dic"`
	Rule         string    `json:"rule" form:"rule"`
	BID          int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	OtherIP      string    `json:"other_ip" gorm:"column:other_ip" form:"other_ip"`
	Del          int8      `json:"del" form:"del"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// ActBwsUserAchievement def.
type ActBwsUserAchievement struct {
	ID    int64     `json:"id" form:"id"`
	MID   int64     `json:"mid" gorm:"column:mid"  form:"mid"`
	AID   int64     `json:"aid" gorm:"column:aid"  form:"aid"`
	BID   int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	Key   string    `json:"key" form:"key"`
	Del   int8      `json:"del" form:"del"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// ActBwsUserPoint def.
type ActBwsUserPoint struct {
	ID     int64     `json:"id" form:"id"`
	MID    int64     `json:"mid" gorm:"column:mid"  form:"mid"`
	PID    int64     `json:"pid" gorm:"column:pid"  form:"pid"`
	BID    int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	Key    string    `json:"key" form:"key"`
	Points int64     `json:"points" form:"points"`
	Del    int8      `json:"del" form:"del"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// ActBwsUser def.
type ActBwsUser struct {
	ID    int64     `json:"id" form:"id"`
	MID   int64     `json:"mid" gorm:"column:mid"  form:"mid"`
	BID   int64     `json:"bid" gorm:"column:bid"  form:"bid"`
	Key   string    `json:"key" form:"key"`
	Del   int8      `json:"del" form:"del"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// TableName ActBws def.
func (ActBws) TableName() string {
	return "act_bws"
}
