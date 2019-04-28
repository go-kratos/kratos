package space

import (
	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/app/interface/main/app-interface/model/bplus"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"

	"strconv"
)

const (
	_gotoAv      = 0
	_gotoArticle = 1
	_gotoClip    = 2
	_gotoAlbum   = 3
	_gotoAudio   = 4
)

// Contributes struct
type Contributes struct {
	Tab   *Tab    `json:"tab,omitempty"`
	Items []*Item `json:"items,omitempty"`
	Links *Links  `json:"links,omitempty"`
}

// Tab struct
type Tab struct {
	Archive   bool `json:"archive"`
	Article   bool `json:"article"`
	Clip      bool `json:"clip"`
	Album     bool `json:"album"`
	Favorite  bool `json:"favorite"`
	Bangumi   bool `json:"bangumi"`
	Coin      bool `json:"coin"`
	Like      bool `json:"like"`
	Community bool `json:"community"`
	Dynamic   bool `json:"dynamic"`
	Audios    bool `json:"audios"`
	Shop      bool `json:"shop"`
}

// Item struct
type Item struct {
	ID        int64             `json:"id,omitempty"`
	TypeName  string            `json:"tname,omitempty"`
	Category  *article.Category `json:"category,omitempty"`
	Title     string            `json:"title,omitempty"`
	Cover     string            `json:"cover,omitempty"`
	Tag       string            `json:"tag,omitempty"`
	Tags      []*article.Tag    `json:"tags,omitempty"`
	Desc      string            `json:"description,omitempty"`
	URI       string            `json:"uri,omitempty"`
	Param     string            `json:"param,omitempty"`
	Goto      string            `json:"goto,omitempty"`
	Length    string            `json:"length,omitempty"`
	Duration  int64             `json:"duration,omitempty"`
	Banner    string            `json:"banner,omitempty"`
	Play      int               `json:"play,omitempty"`
	Comment   int               `json:"comment,omitempty"`
	Danmaku   int               `json:"danmaku,omitempty"`
	Count     int               `json:"count,omitempty"`
	Reply     int               `json:"reply,omitempty"`
	CTime     xtime.Time        `json:"ctime,omitempty"`
	MTime     xtime.Time        `json:"mtime,omitempty"`
	ImageURLs []string          `json:"image_urls,omitempty"`
	Pictures  []*bplus.Pictures `json:"pictures,omitempty"`
	Words     int64             `json:"words,omitempty"`
	Stats     *article.Stats    `json:"stats,omitempty"`
	AuthType  int               `json:"authType,omitempty"`
	Member    int64             `json:"member,omitempty"`
}

// Links struct
type Links struct {
	Previous int64 `json:"previous,omitempty"`
	Next     int64 `json:"next,omitempty"`
}

// Link func
func (l *Links) Link(sinceID, untilID int64) {
	if sinceID < 0 || untilID < 0 {
		return
	}
	l.Previous = sinceID
	l.Next = untilID
}

// Items struct
type Items []*Item

//Len()
func (is Items) Len() int { return len(is) }

//Less()
func (is Items) Less(i, j int) bool {
	var it, jt xtime.Time
	if is[i] != nil {
		it = is[i].CTime
	}
	if is[j] != nil {
		jt = is[j].CTime
	}
	return it > jt
}

//Swap()
func (is Items) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

// Clip struct
type Clip struct {
	ID       int64      `json:"id"`
	Duration int64      `json:"duration"`
	CTime    xtime.Time `json:"ctime"`
	View     int        `json:"view"`
	Damaku   int        `json:"damaku"`
	Title    string     `json:"title"`
	Cover    string     `json:"cover"`
	Tag      string     `json:"tag"`
}

// Album struct
type Album struct {
	ID       int64       `json:"doc_id"`
	CTime    xtime.Time  `json:"ctime"`
	Count    int         `json:"count"`
	View     int         `json:"view"`
	Comment  int         `json:"comment"`
	Title    string      `json:"title"`
	Desc     string      `json:"description"`
	Pictures []*Pictures `json:"pictures"`
}

// Pictures struct
type Pictures struct {
	ImgSrc    string `json:"img_src"`
	ImgWidth  string `json:"img_width"`
	ImgHeight string `json:"img_height"`
}

// Tag tag.
type Tag struct {
	Tid  int64  `json:"tid"`
	Name string `json:"name"`
}

// FromArc3 func
func (i *Item) FromArc3(a *api.Arc) {
	i.ID = a.Aid
	i.Title = a.Title
	i.Cover = a.Pic
	i.TypeName = a.TypeName
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.Goto = model.GotoAv
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Duration = a.Duration
	i.CTime = a.PubDate
	i.Play = int(a.Stat.View)
}

// FromArticle func
func (i *Item) FromArticle(a *article.Meta) {
	i.ID = a.ID
	i.Title = a.Title
	i.Category = a.Category
	i.Desc = a.Summary
	i.ImageURLs = a.ImageURLs
	i.CTime = a.Ctime
	i.Tags = a.Tags
	i.Banner = a.BannerURL
	i.Param = strconv.FormatInt(a.ID, 10)
	i.Goto = model.GotoArticle
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Stats = a.Stats
}

// FromClip func
func (i *Item) FromClip(c *bplus.Clip) {
	i.ID = c.ID
	i.Duration = c.Duration
	i.CTime = c.CTime
	i.Play = c.View
	i.Danmaku = c.Damaku
	i.Param = strconv.FormatInt(c.ID, 10)
	i.Goto = model.GotoClip
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Title = c.Title
	i.Cover = c.Cover
	i.Tag = c.Tag
}

// FromAlbum func
func (i *Item) FromAlbum(a *bplus.Album) {
	i.ID = a.ID
	i.CTime = a.CTime
	i.Count = a.Count
	i.Play = a.View
	i.Comment = a.Comment
	i.Param = strconv.FormatInt(a.ID, 10)
	i.Goto = model.GotoAlbum
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Title = a.Title
	i.Desc = a.Desc
	i.Pictures = a.Pictures
}

// FromAudio func
func (i *Item) FromAudio(a *audio.Audio) {
	i.ID = a.ID
	i.CTime = a.CTime
	i.Play = a.Play
	i.Reply = a.Reply
	i.Param = strconv.FormatInt(a.ID, 10)
	i.Goto = model.GotoAudio
	i.URI = a.Schema
	i.Cover = a.Cover
	i.Title = a.Title
	i.AuthType = a.AuthType
	i.Duration = a.Duration
}

// FormatKey func
func (i *Item) FormatKey() {
	switch i.Goto {
	case model.GotoAv:
		i.Member = i.ID<<6 | _gotoAv
	case model.GotoArticle:
		i.Member = i.ID<<6 | _gotoArticle
	case model.GotoClip:
		i.Member = i.ID<<6 | _gotoClip
	case model.GotoAlbum:
		i.Member = i.ID<<6 | _gotoAlbum
	case model.GotoAudio:
		i.Member = i.ID<<6 | _gotoAudio
	default:
		i.Member = i.ID
	}
}

// ParseKey func
func (i *Item) ParseKey() {
	i.ID = i.Member >> 6
	switch int(i.Member & 0x3f) {
	case _gotoAv:
		i.Goto = model.GotoAv
	case _gotoArticle:
		i.Goto = model.GotoArticle
	case _gotoClip:
		i.Goto = model.GotoClip
	case _gotoAlbum:
		i.Goto = model.GotoAlbum
	case _gotoAudio:
		i.Goto = model.GotoAudio
	}
}

// Attrs struct
type Attrs struct {
	Archive bool `json:"archive,omitempty"`
	Article bool `json:"article,omitempty"`
	Clip    bool `json:"clip,omitempty"`
	Album   bool `json:"album,omitempty"`
	Audio   bool `json:"audio,omitempty"`
}
