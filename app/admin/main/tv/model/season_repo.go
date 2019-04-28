package model

import (
	"net/url"
	"reflect"
	"strconv"

	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

// TVEpContent reprensents the content table
type TVEpContent struct {
	ID        int64  `form:"id" params:"id" validate:"required"`
	CID       int    `form:"cid" params:"cid" gorm:"column:cid" validate:"required"`
	SeasonID  int64  `form:"season_id" params:"season_id" validate:"required"`
	Title     string `form:"title" params:"title"`
	LongTitle string `form:"long_title" params:"long_title"`
	Cover     string `form:"cover" params:"cover"`
	Length    int32  `form:"length" params:"length"`
	Order     int    `form:"order" params:"order" validate:"required"`
	PayStatus int    `form:"pay_status" validate:"required" gorm:"-"`
	Desc      string `form:"desc" gorm:"-"`
	IsDeleted int8
	Ctime     time.Time
	Mtime     time.Time
}

// TVEpSeason represents the season table
type TVEpSeason struct {
	ID          int64     `form:"id" json:"id" params:"id" validate:"required" gorm:"column:id"`
	OriginName  string    `form:"origin_name" json:"origin_name" params:"origin_name" validate:"required"`
	Title       string    `form:"title" json:"title" params:"title"`
	Alias       string    `form:"alias" json:"alias" params:"alias"`
	Category    int       `form:"category" json:"category" params:"category" validate:"required" gorm:"column:category"`
	Desc        string    `form:"desc" json:"desc" params:"desc"`
	Style       string    `form:"style" json:"style" params:"style"`
	Area        string    `form:"area" json:"area" params:"area"`
	PlayTime    time.Time `form:"play_time" json:"play_time" params:"play_time" validate:"required"`
	Info        int       `form:"info" json:"info" params:"info" validate:"required"`
	State       string    `form:"state" json:"state" validate:"required" params:"state"`
	TotalNum    string    `form:"total_num" json:"total_num" params:"total_num" validate:"required"`
	Upinfo      string    `form:"upinfo" json:"upinfo" params:"upinfo"`
	Staff       string    `form:"staff" json:"staff" params:"staff"`
	Role        string    `form:"role" json:"role" params:"role"`
	Copyright   string    `form:"copyright" json:"copyright" params:"copyright"`
	Cover       string    `form:"cover" json:"cover" params:"cover" gorm:"column:cover"`
	Check       int       `json:"check"`
	IsDeleted   int       `json:"is_deleted"`
	AuditTime   int       `json:"audit_time"`
	Valid       int       `json:"valid"`
	Reason      string    `json:"reason"`
	Version     string    `json:"version" form:"version"`   // v1.13 new fields, movie, OVA or normal
	Producer    string    `json:"producer" form:"producer"` // v1.13 new fields, BBC, CCTV etc
	AliasSearch string    `json:"alias_search" form:"alias_search"`
	Brief       string    `json:"brief" form:"brief"`
	Status      string    `json:"status" form:"status"`
}

// SeaRepoCore def.
type SeaRepoCore struct {
	ID         int64     `json:"id" params:"id"`
	OriginName string    `json:"origin_name" params:"origin_name"`
	Title      string    `json:"title" params:"title"`
	Alias      string    `json:"alias" params:"alias"`
	Category   int8      `json:"category" params:"category"`
	Desc       string    `json:"desc" params:"desc"`
	Style      string    `json:"style" params:"style"`
	Area       string    `json:"area" params:"area"`
	Info       int8      `json:"info" params:"info"`
	State      int8      `json:"state" params:"state"`
	TotalNum   int32     `json:"total_num" params:"total_num"`
	Upinfo     string    `json:"upinfo" params:"upinfo"`
	Staff      string    `json:"staff" params:"staff"`
	Role       string    `json:"role" params:"role"`
	Copyright  string    `json:"copyright" params:"copyright"`
	Cover      string    `json:"cover" params:"cover" gorm:"column:cover"`
	Check      int8      `json:"check"`
	IsDeleted  int8      `json:"is_deleted"`
	AuditTime  int       `json:"audit_time"`
	Ctime      time.Time `json:"ctime"`
	Valid      int8      `json:"valid"`
	InjectTime time.Time `json:"inject_time"`
	Reason     string    `json:"reason"`
}

// SeaRepoDB def.
type SeaRepoDB struct {
	SeaRepoCore
	PlayTime time.Time `gorm:"column:play_time"`
	Mtime    time.Time `json:"mtime"`
}

// SeaRepoList def.
type SeaRepoList struct {
	SeaRepoCore
	Mtime   string `json:"mtime"`
	Pubdate string `json:"pubdate"`
}

// ToList transforms a SeaRepoDB to list, time transformation
func (v *SeaRepoDB) ToList() (list *SeaRepoList) {
	list = &SeaRepoList{
		SeaRepoCore: v.SeaRepoCore,
	}
	list.Mtime = v.Mtime.Time().Format(mysql.TimeFormat)
	list.Pubdate = v.PlayTime.Time().Format(mysql.TimeFormat)
	return
}

// TableName gives the table name of content
func (*TVEpContent) TableName() string {
	return "tv_ep_content"
}

// TableName gives the table name of season
func (*TVEpSeason) TableName() string {
	return "tv_ep_season"
}

// TableName gives the table name of season
func (v *SeaRepoDB) TableName() string {
	return "tv_ep_season"
}

// SeasonRepoPager def.
type SeasonRepoPager struct {
	TotalCount int64          `json:"total_count"`
	Pn         int            `json:"pn"`
	Ps         int            `json:"ps"`
	Items      []*SeaRepoList `json:"items"`
}

// Updated picks value from request and compare with the struct to analyse the difference
func (sn TVEpSeason) Updated(req url.Values) (fields map[string]interface{}) {
	var (
		vl = reflect.ValueOf(sn)
		tp = reflect.TypeOf(sn)
	)
	fields = make(map[string]interface{})
	for i := 0; i < vl.NumField(); i++ {
		var (
			t    = tp.Field(i)
			v    = vl.Field(i)
			name = t.Tag.Get("json")
		)
		if reqV := req.Get(name); reqV != "" {
			if kind := t.Type.Kind(); kind == reflect.Int || kind == reflect.Int64 {
				if v.Int() != atoi(reqV) {
					fields[name] = atoi(reqV)
				}
			} else {
				if v.String() != reqV {
					fields[name] = reqV
				}
			}
		}
	}
	return
}

// ToContent transforms an ep to content object
func (epc *TVEpContent) ToContent(isInit bool) (res *Content) {
	res = &Content{
		Title:     epc.LongTitle,
		Subtitle:  epc.Title,
		Desc:      epc.Desc,
		Cover:     epc.Cover,
		SeasonID:  int(epc.SeasonID),
		CID:       epc.CID,
		EPID:      int(epc.ID),
		PayStatus: epc.PayStatus,
	}
	if isInit {
		res.State = 1
	}
	return
}

func atoi(value string) (intval int64) {
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		intval = 0
	}
	return intval
}
