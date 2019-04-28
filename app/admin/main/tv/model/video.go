package model

import (
	arccli "go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// Video is used from PGC video
type Video struct {
	ID          int       `gorm:"column:id" json:"id"`
	AID         int       `gorm:"column:aid" json:"aid"`
	Eptitle     string    `gorm:"column:eptitle" json:"eptitle"`
	Description string    `gorm:"column:description" json:"description"`
	CID         int64     `gorm:"column:cid" json:"cid"`
	Duration    int       `gorm:"column:duration" json:"duration"`
	IndexOrder  int       `gorm:"column:duration" json:"index_order"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
	InjectTime  time.Time `gorm:"column:inject_time" json:"inject_time"`
	Valid       uint8     `gorm:"column:valid" json:"valid"`
	Submit      uint8     `gorm:"column:submit" json:"submit"`
	Retry       int       `gorm:"column:retry" json:"retry"`
	Result      int       `gorm:"column:result" json:"result"`
	Deleted     uint8     `gorm:"column:deleted" json:"deleted"`
	State       int       `gorm:"column:state" json:"state"`
	Reason      string    `gorm:"column:reason" json:"reason"`
	Manual      int       `gorm:"column:manual" json:"manual"`
}

// VideoListParam is used for vlideolist funtion param valid
type VideoListParam struct {
	CID    string `form:"cid" json:"cid"`
	VID    string `form:"vid" json:"vid"`
	Typeid int16  `form:"typeid" json:"typeid"`
	Pid    int32  `form:"pid" json:"-"`
	Valid  string `form:"valid" json:"valid"`
	Order  int    `form:"order" json:"order" default:"2"`
	Pn     int    `form:"pn" json:"pn"  default:"1"`
	Ps     int    `form:"ps" json:"ps"  default:"20"`
}

// VideoListQuery is used for selecting the field of pgc video
type VideoListQuery struct {
	ID          string    `json:"id"`
	VID         string    `json:"vid" gorm:"column:cid"`
	CID         string    `json:"cid" gorm:"column:aid"`
	Eptitle     string    `json:"eptitle"`
	Valid       string    `json:"valid" gorm:"column:valid"`
	Mtime       time.Time `json:"mtime"`
	SeasonTitle string    `json:"season_title" gorm:"column:title"`
	TypeID      int32     `json:"typeid" gorm:"column:typeid"`
	PTypeID     int32     `json:"parent_typeid"`
	Page        int       `json:"page" gorm:"column:index_order"`
}

// VideoListPager is used by video list function to return result and page info
type VideoListPager struct {
	Items []*VideoListQuery `json:"items"`
	Page  *Page             `json:"page"`
}

// TableName ugc_video
func (a VideoListQuery) TableName() string {
	return "ugc_video"
}

// TableName ugc_video
func (video Video) TableName() string {
	return "ugc_video"
}

// ConsultRes transforms an archive to ArcRes
func (arc *Archive) ConsultRes(dict map[int32]*arccli.Tp) (res *ArcRes) {
	var pid int32
	res = &ArcRes{}
	if cat, ok := dict[arc.TypeID]; ok {
		pid = cat.Pid
		res.SecondCat = cat.Name
	}
	if pid != 0 {
		if pcat, ok := dict[pid]; ok {
			res.FirstCat = pcat.Name
		}
	}
	res.Status = int(arc.Result)
	res.AVID = arc.AID
	res.Title = arc.Title
	res.PubTime = arc.Pubtime.Time().Format("2006-01-02 15:04:05")
	if arc.InjectTime >= 0 {
		res.InjectTime = arc.InjectTime.Time().Format("2006-01-02 15:04:05")
	}
	res.Reason = arc.Reason
	return
}

// ConsultRes transforms an video to VideoRes
func (video *Video) ConsultRes() (res *VideoRes) {
	res = &VideoRes{
		CID:    video.CID,
		Title:  video.Eptitle,
		Page:   video.IndexOrder,
		Status: video.Result,
		Ctime:  video.Ctime.Time().Format("2006-01-02 15:04:05"),
		Reason: video.Reason,
	}
	if video.InjectTime >= 0 {
		res.InjectTime = video.InjectTime.Time().Format("2006-01-02 15:04:05")
	}
	return
}
