package model

import xtime "go-common/library/time"

const (
	StateNormal = 0
	StateDelted = 1
)

type ShortUrl struct {
	ID         int64      `json:"id"`
	Mid        int64      `json:"mid"`
	Short      string     `json:"short"`
	Long       string     `json:"long"`
	State      int8       `json:"state"`
	CTime      xtime.Time `json:"-"`
	MTime      xtime.Time `json:"-"`
	CreateTime string     `json:"ctime"`
}

type Param struct {
	ID  int64  `form:"id"`
	Mid int64  `form:"mid"`
	Uri string `form:"url"`
	Pn  string `form:"pn"`
	Ps  string `form:"ps"`
}

func (s *ShortUrl) FormatDate() {
	s.CreateTime = s.CTime.Time().Format("2006-01-02 15:04:05")
}
