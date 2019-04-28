package like

import (
	xtime "go-common/library/time"
)

// Subject type.
const (
	VIDEO         = 1
	PICTURE       = 2
	DRAWYOO       = 3
	VIDEOLIKE     = 4
	PICTURELIKE   = 5
	DRAWYOOLIKE   = 6
	TEXT          = 7
	TEXTLIKE      = 8
	ONLINEVOTE    = 9
	QUESTION      = 10
	LOTTERY       = 11
	ARTICLE       = 12
	VIDEO2        = 13
	MUSIC         = 15
	PHONEVIDEO    = 16
	SMALLVIDEO    = 17
	RESERVATION   = 18
	MISSIONGROUP  = 19
	STORYKING     = 20
	FLAGFIRST     = 1
	FLAGSPY       = 2
	FLAGUSTIME    = 4
	FLAGUETIME    = 8
	FLAGLEVEL     = 16
	FLAGIP        = 32
	FLAGRANKCLOSE = 64
	FLAGPHONEBIND = 128
)

// Subject group type
var (
	VIDEOALL = []int64{VIDEO, VIDEOLIKE, VIDEO2, PHONEVIDEO, SMALLVIDEO}
	LIKETYPE = []int64{VIDEOLIKE, PICTURELIKE, DRAWYOOLIKE, TEXTLIKE}
)

// Subject struct
type Subject struct {
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	Dic      string     `json:"dic"`
	Cover    string     `json:"cover"`
	Stime    xtime.Time `json:"stime"`
	Interval int32      `json:"interval"`
	Tlimit   int32      `json:"tlimit"`
	Ltime    int32      `json:"ltime"`
	List     []*Like    `json:"list"`
}

// SubItem .
type SubItem struct {
	ID    int64      `json:"id"`
	Ctime xtime.Time `json:"ctime"`
}

// SubjectStat .
type SubjectStat struct {
	Sid   int64 `json:"sid" form:"sid" validate:"min=1"`
	Count int64 `json:"count" form:"count"`
	View  int64 `form:"view" form:"view"`
	Like  int64 `form:"like" form:"like"`
	Fav   int64 `form:"fav" form:"fav"`
	Coin  int64 `form:"coin" form:"coin"`
}

// SubjectScore .
type SubjectScore struct {
	Score int64 `json:"score"`
}

// Page .
type Page struct {
	Num   int   `json:"num"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// SubProtocol .
type SubProtocol struct {
	*SubjectItem
	*ActSubjectProtocol
}
