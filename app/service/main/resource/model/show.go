package model

import (
	"encoding/json"
	"strconv"

	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// Card audio card struct
type Card struct {
	ID        int    `json:"-"`
	Tab       int    `json:"-"`
	RegionID  int    `json:"-"`
	Type      int    `json:"-"`
	Title     string `json:"-"`
	Cover     string `json:"-"`
	Rtype     int    `json:"-"`
	Rvalue    string `json:"-"`
	PlatVer   string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	TypeStr   string `json:"-"`
	Goto      string `json:"-"`
	Param     string `json:"-"`
	URI       string `json:"-"`
	Desc      string `json:"-"`
	TagID     int    `json:"-"`
}

// Content audio content struct
type Content struct {
	ID     int    `json:"-"`
	Module int    `json:"-"`
	RecID  int    `json:"-"`
	Type   int8   `json:"-"`
	Value  string `json:"-"`
	Title  string `json:"-"`
	TagID  int    `json:"-"`
}

// PlatLimit audio plat limit struct
type PlatLimit struct {
	Plat      int8   `json:"plat"`
	Build     int    `json:"build"`
	Condition string `json:"conditions"`
}

// ShowItem audio show item struct
type ShowItem struct {
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	URI    string `json:"uri"`
	NewURI string `json:"-"`
	Param  string `json:"param"`
	Goto   string `json:"goto"`
	// up
	Mid            int64           `json:"mid,omitempty"`
	Name           string          `json:"name,omitempty"`
	Face           string          `json:"face,omitempty"`
	Follower       int             `json:"follower,omitempty"`
	Attribute      int             `json:"attribute,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
	// stat
	Play    int `json:"play,omitempty"`
	Danmaku int `json:"danmaku,omitempty"`
	Reply   int `json:"reply,omitempty"`
	Fav     int `json:"favourite,omitempty"`
	// movie and bangumi badge
	Status    int8   `json:"status,omitempty"`
	CoverMark string `json:"cover_mark,omitempty"`
	// ranking
	Pts      int64       `json:"pts,omitempty"`
	Children []*ShowItem `json:"children,omitempty"`
	// av
	PubDate xtime.Time `json:"pubdate"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// region
	Rid   int    `json:"rid,omitempty"`
	Rname string `json:"rname,omitempty"`
	Reid  int    `json:"reid,omitempty"`
	//new manager
	Desc  string `json:"desc,omitempty"`
	Stime string `json:"stime,omitempty"`
	Etime string `json:"etime,omitempty"`
	Like  int    `json:"like,omitempty"`
}

// OfficialVerify audio verify struct
type OfficialVerify struct {
	Type int    `json:"type"`
	Desc string `json:"desc"`
}

// Head audio struct
type Head struct {
	CardID    int         `json:"card_id,omitempty"`
	Title     string      `json:"title,omitempty"`
	Cover     string      `json:"cover,omitempty"`
	Type      string      `json:"type,omitempty"`
	Date      int64       `json:"date,omitempty"`
	Plat      int8        `json:"-"`
	Build     int         `json:"-"`
	Condition string      `json:"-"`
	URI       string      `json:"uri,omitempty"`
	Goto      string      `json:"goto,omitempty"`
	Param     string      `json:"param,omitempty"`
	Body      []*ShowItem `json:"body,omitempty"`
}

// CardPlatChange audio card change plat
func (c *Card) CardPlatChange() (platlinits []*PlatLimit) {
	platlinits = platJSONChange(c.PlatVer)
	return
}

// platJSONChange json change plat build condition
func platJSONChange(jsonStr string) (platlinits []*PlatLimit) {
	var tmp []struct {
		Plat      string `json:"plat"`
		Build     string `json:"build"`
		Condition string `json:"conditions"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &tmp); err == nil {
		for _, limit := range tmp {
			platlinit := &PlatLimit{}
			switch limit.Plat {
			case "0": // resource android
				platlinit.Plat = PlatAndroid
			case "1": // resource iphone
				platlinit.Plat = PlatIPhone
			case "2": // resource pad
				platlinit.Plat = PlatIPad
			case "5": // resource iphone_i
				platlinit.Plat = PlatIPhoneI
			case "8": // resource android_i
				platlinit.Plat = PlatAndroidI
			}
			platlinit.Build, _ = strconv.Atoi(limit.Build)
			platlinit.Condition = limit.Condition
			platlinits = append(platlinits, platlinit)
		}
	}
	return
}

// FromArchivePB from archive archive.
func (i *ShowItem) FromArchivePB(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = FillURI(GotoAv, i.Param)
	i.Goto = GotoAv
	i.Play = int(a.Stat.View)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Name = a.Author.Name
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.PubDate = a.PubDate
	i.Rid = int(a.TypeID)
	i.Rname = a.TypeName
	i.Duration = a.Duration
	i.Like = int(a.Stat.Like)
	if a.Access > 0 {
		i.Play = 0
	}
}

// FillBuildURI fill url by plat build
func (h *Head) FillBuildURI(plat int8, build int) {
	switch h.Goto {
	case GotoDaily:
		if (plat == PlatIPhone && build > 6670) || (plat == PlatAndroid && build > 5250000) {
			h.URI = "bilibili://pegasus/list/daily/" + h.Param
		}
	}
}

// SideBars for side bars
type SideBars struct {
	SideBar []*SideBar                `json:"sidebar,omitempty"`
	Limit   map[int64][]*SideBarLimit `json:"limit,omitempty"`
}

// SideBar for side bar
type SideBar struct {
	ID           int64      `json:"id,omitempty"`
	Tip          int        `json:"tip,omitempty"`
	Rank         int        `json:"rank,omitempty"`
	Logo         string     `json:"logo,omitempty"`
	LogoWhite    string     `json:"logo_white,omitempty"`
	Name         string     `json:"name,omitempty"`
	Param        string     `json:"param,omitempty"`
	Module       int        `json:"module,omitempty"`
	Plat         int8       `json:"-"`
	Build        int        `json:"-"`
	Conditions   string     `json:"-"`
	OnlineTime   xtime.Time `json:"online_time"`
	NeedLogin    int8       `json:"need_login,omitempty"`
	WhiteURL     string     `json:"white_url,omitempty"`
	Menu         int8       `json:"menu,omitempty"`
	LogoSelected string     `json:"logo_selected,omitempty"`
	TabID        string     `json:"tab_id,omitempty"`
	Red          string     `json:"red_dot_url,omitempty"`
	Language     string     `json:"language,omitempty"`
}

// SideBarLimit side bar limit
type SideBarLimit struct {
	ID        int64  `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}
