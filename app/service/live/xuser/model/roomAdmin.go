package model

import xtime "go-common/library/time"

// TableName is used to identify table name in gorm
func (ra *RoomAdmin) TableName() string {
	return "ap_room_admin"
}

// RoomAdmin .
type RoomAdmin struct {
	Id     int64      `json:"id" gorm:"column:id"`
	Uid    int64      `json:"uid" gorm:"column:uid"`
	Roomid int64      `json:"roomid" gorm:"column:roomid"`
	Ctime  xtime.Time `json:"ctime" gorm:"comumn:ctime"`
}

// RoomAdmins multi RoomAdmin .
type RoomAdmins []*RoomAdmin

// Len returns length of RoomAdmins.
func (ras RoomAdmins) Len() int {
	return len(ras)
}

// Swap .
func (ras RoomAdmins) Swap(i, j int) {
	ras[i], ras[j] = ras[j], ras[i]
}

// Less returns sorting rule.
func (ras RoomAdmins) Less(i, j int) bool {
	return ras[i].Ctime < ras[j].Ctime
}
