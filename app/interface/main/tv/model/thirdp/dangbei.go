package thirdp

import (
	"strconv"

	"go-common/app/interface/main/tv/model"
)

// DBeiPage is the dangbei page struct
type DBeiPage struct {
	List  []*DBeiSeason   `json:"list"`
	Pager *model.IdxPager `json:"pager"`
}

// DBeiSeason is the dangbei season struct
type DBeiSeason struct {
	SeasonID    *int64 `json:"cid,omitempty"`
	Cover       string `json:"cover"`
	Desc        string `json:"desc"`
	Title       string `json:"title"`
	UpInfo      string `json:"upinfo"`
	Category    string `json:"category"` // - cn, jp, movie, tv, documentary
	Area        string `json:"area"`     // - cn, jp, others
	Playtime    string `json:"play_time"`
	Role        string `json:"role"`
	Staff       string `json:"staff"`
	NewestOrder int    `json:"newest_order"` // the newest passed ep's order
	NewestNB    int    `json:"newest_nb"`
	TotalNum    int    `json:"total_num"`
	Style       string `json:"style"`
	Paystatus   string `json:"pay_status"` // paid or not
	Official    string `json:"official"`   // is official or preview
}

// VideoCMS def.
type VideoCMS struct {
	// Media Info
	CID        int
	Title      string
	AID        int
	IndexOrder int
	// Auth Info
	Valid   int
	Deleted int
	Result  int
}

// PgcCat transforms the pgc category to string
func PgcCat(cat int) (category string) {
	switch cat {
	case 1:
		category = "番剧"
	case 2:
		category = "电影"
	case 3:
		category = "纪录片"
	case 4:
		category = "国创"
	case 5:
		category = "电视剧"
	default:
		category = "其他"
	}
	return
}

// DBeiSn transform the object
func DBeiSn(s *model.SeasonCMS) (dbei *DBeiSeason) {
	var (
		category, area string
		official       = "正片"
	)
	category = PgcCat(int(s.Category))
	if s.NewestNb == 0 { // in case of job calculation fail
		if s.NewestOrder == 0 {
			s.NewestNb = s.TotalNum
		} else {
			s.NewestNb = min(s.TotalNum, s.NewestOrder)
		}
	}
	areaInt, _ := strconv.ParseInt(s.Area, 10, 64)
	if areaInt != 0 {
		switch areaInt {
		case 1:
			area = "中国"
		case 2:
			area = "日本"
		default:
			area = "其他"
		}
	} else {
		area = s.Area
	}
	dbei = &DBeiSeason{
		SeasonID:    &s.SeasonID,
		Cover:       s.Cover,
		Desc:        s.Desc,
		Title:       s.Title,
		UpInfo:      s.UpInfo,
		Category:    category,
		Area:        area,
		Playtime:    s.Playtime.Time().Format("2006-01-02"),
		Role:        s.Role,
		Staff:       s.Staff,
		NewestOrder: s.NewestOrder,
		NewestNB:    s.NewestNb,
		TotalNum:    s.TotalNum,
		Style:       s.Style,
		Paystatus:   "",
		Official:    official,
	}
	return
}

// DbeiArc transforms an arc cms to dangbei season structure
func DbeiArc(s *model.ArcCMS, first, second string) (dbei *DBeiSeason) {
	official := "正片"
	dbei = &DBeiSeason{
		SeasonID:  &s.AID,
		Cover:     s.Cover,
		Desc:      s.Content,
		Title:     s.Title,
		Category:  first,
		Playtime:  s.Pubtime.Time().Format("2006-01-02"),
		TotalNum:  s.Videos,
		Style:     second,
		Paystatus: "",
		Official:  official,
	}
	return
}

// return min value
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// ReqDBeiPages is request for dangbei pages
type ReqDBeiPages struct {
	Page   int64
	LastID int64
	Ps     int64
	TypeC  string
}

// ReqPageID is request for page ID
type ReqPageID struct {
	Page  int64
	ID    int64
	TypeC string
}
