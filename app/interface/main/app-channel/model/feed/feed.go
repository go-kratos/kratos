package feed

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	cardlive "go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	shopping "go-common/app/interface/main/app-card/model/card/show"
	"go-common/app/interface/main/app-channel/model"
	"go-common/app/interface/main/app-channel/model/activity"
	"go-common/app/interface/main/app-channel/model/card"
	"go-common/app/interface/main/app-channel/model/channel"
	"go-common/app/interface/main/app-channel/model/dislike"
	"go-common/app/interface/main/app-channel/model/recommend"
	bustag "go-common/app/interface/main/tag/model"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	xtime "go-common/library/time"
)

const (
	_seasonNoSeason   = 1
	_seasonUpper      = 4
	_activityForm     = "2006-01-02 15:04:05"
	_convergeMinCount = 2
)

type Show struct {
	Topic *Item   `json:"topic,omitempty"`
	Feed  []*Item `json:"feed"`
}

// Item is feed item,
type Item struct {
	Title       string      `json:"title,omitempty"`
	Subtitle    string      `json:"subtitle,omitempty"`
	Cover       string      `json:"cover,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Redirect    string      `json:"redirect,omitempty"`
	RedirectURI string      `json:"redirect_uri,omitempty"`
	Param       string      `json:"param,omitempty"`
	Goto        string      `json:"goto,omitempty"`
	ViewType    string      `json:"view_type,omitempty"`
	Kind        string      `json:"kind,omitempty"`
	Desc        string      `json:"desc,omitempty"`
	Play        int         `json:"play,omitempty"`
	Danmaku     int         `json:"danmaku,omitempty"`
	Reply       int         `json:"reply,omitempty"`
	Fav         int         `json:"favorite,omitempty"`
	Coin        int         `json:"coin,omitempty"`
	Share       int         `json:"share,omitempty"`
	Like        int         `json:"like,omitempty"`
	Count       int         `json:"count,omitempty"`
	Status      int8        `json:"status,omitempty"`
	Type        int8        `json:"type,omitempty"`
	Badge       string      `json:"badge,omitempty"`
	StatType    int8        `json:"stat_type,omitempty"`
	RcmdReason  *RcmdReason `json:"rcmd_reason,omitempty"`
	Item        []*Item     `json:"item,omitempty"`
	// sortedset index
	Idx int64 `json:"idx,omitempty"`
	// av info
	Cid             int64                     `json:"cid,omitempty"`
	Rid             int32                     `json:"tid,omitempty"`
	TName           string                    `json:"tname,omitempty"`
	Tag             *Tag                      `json:"tag,omitempty"`
	DisklikeReasons []*dislike.DisklikeReason `json:"dislike_reasons,omitempty"`
	PTime           xtime.Time                `json:"ctime,omitempty"`
	Autoplay        int32                     `json:"autoplay,omitempty"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// upper
	Mid      int64         `json:"mid,omitempty"`
	Name     string        `json:"name,omitempty"`
	Face     string        `json:"face,omitempty"`
	IsAtten  int8          `json:"is_atten,omitempty"`
	Fans     int64         `json:"fans,omitempty"`
	RecCnt   int           `json:"recent_count,omitempty"`
	Recent   []*Item       `json:"recent,omitempty"`
	Official *OfficialInfo `json:"official,omitempty"`
	// live
	Online int    `json:"online,omitempty"`
	Area   string `json:"area,omitempty"`
	AreaID int    `json:"area_id,omitempty"`
	Area2  *Area2 `json:"area2,omitempty"`
	// bangumi
	Index       string `json:"index,omitempty"`
	IndexTitle  string `json:"index_title,omitempty"`
	CoverMark   string `json:"cover_mark,omitempty"`
	Finish      bool   `json:"finish,omitempty"`
	LatestIndex string `json:"last_index,omitempty"`
	// bangumi recommend
	Updates int `json:"updates,omitempty"`
	// live or bangumi
	From int8 `json:"from,omitempty"`
	// adviertisement
	RequestID  string          `json:"request_id,omitempty"`
	CreativeID int64           `json:"creative_id,omitempty"`
	SrcID      int             `json:"src_id,omitempty"`
	IsAd       bool            `json:"is_ad,omitempty"`
	IsAdLoc    bool            `json:"is_ad_loc,omitempty"`
	AdCb       string          `json:"ad_cb,omitempty"`
	ShowURL    string          `json:"show_url,omitempty"`
	ClickURL   string          `json:"click_url,omitempty"`
	ClientIP   string          `json:"client_ip,omitempty"`
	CmMark     int64           `json:"cm_mark,omitempty"`
	AdIndex    int             `json:"ad_index,omitempty"`
	Extra      json.RawMessage `json:"extra,omitempty"`
	CardIndex  int             `json:"card_index,omitempty"`
	// activity
	STime string `json:"stime,omitempty"`
	ETime string `json:"etime,omitempty"`
	// tag
	Tags []*channel.Tag `json:"tags,omitempty"`
	// rank
	Cover1 string `json:"cover1,omitempty"`
	Cover2 string `json:"cover2,omitempty"`
	Cover3 string `json:"cover3,omitempty"`
	// banner`
	Hash string `json:"hash,omitempty"`
	// upper article
	Covers    []string  `json:"covers,omitempty"`
	Temple    int       `json:"temple,omitempty"`
	Template  int       `json:"template,omitempty"`
	Category  *Category `json:"category,omitempty"`
	BannerURL string    `json:"banner_url,omitempty"`
	// game download
	Button   *Button `json:"button,omitempty"`
	Download int32   `json:"download,omitempty"`
	BigCover string  `json:"big_cover,omitempty"`
	// special
	HideBadge bool    `json:"hide_badge,omitempty"`
	Ratio     float64 `json:"ratio,omitempty"`
	// shopping
	City   string `json:"city,omitempty"`
	PType  string `json:"ptype,omitempty"`
	Price  string `json:"price,omitempty"`
	Square string `json:"square,omitempty"`
	// news
	Content string `json:"content,omitempty"`
	// bigdata source
	Source    string          `json:"-"`
	AvFeature json.RawMessage `json:"-"`
	// common
	GotoOrg  string `json:"-"`
	FromType string `json:"from_type,omitempty"`
	Pos      int    `json:"-"`
	// audio
	SongTitle string `json:"song_title,omitempty"`
}

type Tag struct {
	TagID   int64     `json:"tag_id,omitempty"`
	TagName string    `json:"tag_name,omitempty"`
	IsAtten int8      `json:"is_atten,omitempty"`
	Count   *TagCount `json:"count,omitempty"`
	Name    string    `json:"name,omitempty"`
	URI     string    `json:"uri,omitempty"`
	//channel
	ID   int64  `json:"id,omitempty"`
	Face string `json:"face,omitempty"`
	Fans int    `json:"fans,omitempty"`
}

type RcmdReason struct {
	ID           int    `json:"id,omitempty"`
	Content      string `json:"content,omitempty"`
	BgColor      string `json:"bg_color,omitempty"`
	IconLocation string `json:"icon_location,omitempty"`
	Message      string `json:"message,omitempty"`
}

type TagCount struct {
	Atten int `json:"atten,omitempty"`
}

type Category struct {
	ID       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Children *Category `json:"children,omitempty"`
}

type Button struct {
	Name        string `json:"name,omitempty"`
	URI         string `json:"uri,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type Area2 struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Children *Area2 `json:"children,omitempty"`
}

type OfficialInfo struct {
	Role  int8   `json:"role,omitempty"`
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

func (i *Item) fillArcStat(a *archive.Archive3) {
	if a.Access == 0 {
		i.Play = int(a.Stat.View)
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.Coin = int(a.Stat.Coin)
	i.Share = int(a.Stat.Share)
	i.Like = int(a.Stat.Like)
}

func (i *Item) FromPlayerAv(a *archive.ArchiveWithPlayer) {
	if i.Title == "" {
		i.Title = a.Title
	}
	if i.Cover == "" {
		i.Cover = model.CoverURLHTTPS(a.Pic)
	} else {
		i.Cover = model.CoverURLHTTPS(i.Cover)
	}
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.Goto = model.GotoAv
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, model.AvPlayHandler(a.Archive3, a.PlayerInfo))
	i.Cid = a.FirstCid
	i.Rid = a.TypeID
	i.TName = a.TypeName
	i.Desc = strconv.Itoa(int(a.Stat.Danmaku)) + "弹幕"
	i.fillArcStat(a.Archive3)
	i.Duration = a.Duration
	i.Mid = a.Author.Mid
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.PTime = a.PubDate
	i.Cid = a.FirstCid
	i.Autoplay = a.Rights.Autoplay
}

func (i *Item) FromDislikeReason() {
	i.DisklikeReasons = []*dislike.DisklikeReason{
		&dislike.DisklikeReason{ReasonID: _seasonUpper, ReasonName: "UP主:" + i.Name},
		&dislike.DisklikeReason{ReasonID: _seasonNoSeason, ReasonName: "不感兴趣"},
	}
}

func (i *Item) FromRcmdReason(c *card.Card) {
	var content string
	switch c.ReasonType {
	case 0:
		content = ""
	case 1:
		content = "编辑精选"
	case 2:
		content = "热门推荐"
	case 3:
		content = c.Reason
	}
	if content != "" {
		i.RcmdReason = &RcmdReason{ID: 1, Content: content, BgColor: "yellow", IconLocation: "left_top"}
	}
}

func (i *Item) FromLive(r *cardlive.Room) {
	if r.LiveStatus != 1 || r.Title == "" || r.Cover == "" {
		return
	}
	i.Title = r.Title
	i.Cover = r.Cover
	i.Goto = model.GotoLive
	i.Param = strconv.FormatInt(r.RoomID, 10)
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, model.LiveRoomHandler(r))
	i.Name = r.Uname
	i.Mid = r.UID
	i.Face = r.Face
	i.Online = int(r.Online)
	i.Autoplay = 1
	// i.Area = r.Area
	// i.AreaID = r.AreaID
	i.Area2 = &Area2{ID: r.AreaV2ParentID, Name: r.AreaV2ParentName, Children: &Area2{ID: r.AreaV2ID, Name: r.AreaV2Name}}
	i.Autoplay = 1
}

func (i *Item) FromSeason(b *bangumi.Season) {
	if i.Title == "" {
		i.Title = b.Title
	}
	if i.Cover == "" {
		i.Cover = b.Cover
	}
	i.Goto = model.GotoUpBangumi
	i.Param = strconv.FormatInt(b.SeasonID, 10)
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, nil)
	i.Play = int(b.PlayCount)
	i.Fav = int(b.Favorites)
	i.Type = b.SeasonType
	i.Badge = b.TypeBadge
	i.Desc = b.UpdateDesc
	i.Face = b.SeasonCover
	i.Square = b.SeasonCover
}

func (i *Item) FromPGCSeason(s *episodegrpc.EpisodeCardsProto) {
	if i.Title == "" {
		i.Title = s.Season.Title
	}
	if i.Cover == "" {
		i.Cover = s.Cover
	}
	i.Goto = model.GotoBangumi
	i.Param = strconv.Itoa(int(s.EpisodeId))
	i.URI = model.FillURI(model.GotoBangumi, i.Param, 0, 0, 0, nil)
	i.Index = s.Title
	i.IndexTitle = s.LongTitle
	i.Status = int8(s.Season.SeasonStatus)
	i.CoverMark = s.Season.Badge
	i.Play = int(s.Season.Stat.View)
	i.Fav = int(s.Season.Stat.Follow)
	i.Type = int8(s.Season.SeasonType)
	i.Badge = s.Season.SeasonTypeName
	if s.Season.IsFinish == 1 {
		i.Finish = true
	}
	i.Count = int(s.Season.TotalCount)
	i.LatestIndex = s.Title
	i.Desc = s.Season.NewEpShow
	i.Face = s.Season.Cover
	i.Square = s.Season.Cover
}

func (i *Item) FromActivity(a *activity.Activity, now time.Time) {
	stime, err := time.ParseInLocation(_activityForm, a.STime, time.Local)
	if err != nil {
		return
	}
	etime, err := time.ParseInLocation(_activityForm, a.ETime, time.Local)
	if err != nil {
		return
	}
	if now.After(etime) {
		i.Status = 1
	} else if now.Before(stime) {
		i.Status = 2
	}
	i.Title = a.Name
	i.Cover = a.H5Cover
	i.Goto = model.GotoActivity
	i.URI = model.FillURI(i.Goto, a.H5URL, 0, 0, 0, nil)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, 0)
	i.Desc = a.Desc
	i.STime = a.STime
	i.ETime = a.ETime
	i.Param = strconv.FormatInt(a.ID, 10)
}

func (i *Item) FromTopic(a *activity.Activity) {
	i.Title = a.Name
	i.Cover = a.H5Cover
	i.Goto = model.GotoTopic
	i.URI = model.FillURI(i.Goto, a.H5URL, 0, 0, 0, nil)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, 0)
	i.Desc = a.Desc
	i.Param = strconv.FormatInt(a.ID, 10)
}

func (i *Item) FromSpecial(id int64, title, cover, desc, url string, typ int, badge string, size string) {
	if title == "" || cover == "" {
		return
	}
	i.Title = title
	i.Cover = cover
	i.Goto = model.GotoSpecial
	i.URI = model.FillURI(i.Goto, url, typ, 0, 0, nil)
	i.Redirect = model.FillRedirect(i.Goto, typ)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, typ)
	i.Desc = desc
	i.Param = strconv.FormatInt(id, 10)
	i.HideBadge = true
	i.Badge = badge
	var ratio float64
	if size == "1020x300" {
		ratio = 34
	} else if size == "1020x378" {
		ratio = 27
	}
	i.Ratio = ratio
}

func (i *Item) FromTopstick(id int64, title, cover, desc, url string, typ int) {
	if title == "" {
		return
	}
	i.Title = title
	i.Goto = model.GotoTopstick
	i.URI = model.FillURI(i.Goto, url, typ, 0, 0, nil)
	i.Redirect = model.FillRedirect(i.Goto, typ)
	if desc == "" {
		i.Desc = "立即查看"
	} else {
		i.Desc = desc
	}
	i.Param = strconv.FormatInt(id, 10)
}

func (i *Item) FromSpecialS(id int64, title, cover, desc, url string, typ int, badge string) {
	if title == "" || cover == "" {
		return
	}
	i.Title = title
	i.Cover = cover
	i.Goto = model.GotoSpecialS
	i.URI = model.FillURI(i.Goto, url, typ, 0, 0, nil)
	i.Redirect = model.FillRedirect(i.Goto, typ)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, typ)
	i.Desc = desc
	i.Param = strconv.FormatInt(id, 10)
	i.HideBadge = true
	i.Badge = badge
}

func (i *Item) FromConverge(c *operate.Converge, am map[int64]*archive.ArchiveWithPlayer, rm map[int64]*cardlive.Room, artm map[int64]*article.Meta) {
	if len(c.Items) < _convergeMinCount {
		return
	}
	cis := []*Item{}
	for _, content := range c.Items {
		ci := &Item{Title: content.Title}
		switch content.Goto {
		case model.GotoAv:
			if a, ok := am[content.Pid]; ok && a.IsNormal() {
				if ci.Title == "" {
					ci.Title = a.Title
				}
				ci.Cover = a.Pic
				ci.Goto = model.GotoAv
				ci.Param = strconv.FormatInt(a.Aid, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, 0, 0, 0, model.AvPlayHandler(a.Archive3, a.PlayerInfo))
				ci.fillArcStat(a.Archive3)
				ci.Duration = a.Duration
				cis = append(cis, ci)
			}
		case model.GotoLive:
			if r, ok := rm[content.ID]; ok {
				if r.LiveStatus == 0 || r.Title == "" || r.Cover == "" {
					continue
				}
				if ci.Title == "" {
					ci.Title = r.Title
				}
				ci.Cover = r.Cover
				ci.Goto = model.GotoLive
				ci.Param = strconv.FormatInt(r.RoomID, 10)
				ci.Online = int(r.Online)
				ci.URI = model.FillURI(ci.Goto, ci.Param, 0, 0, 0, model.LiveRoomHandler(r))
				// ci.RedirectURI = model.FillRedirectURI(ci.Goto, ci.URI, 0)
				ci.Badge = "直播"
				cis = append(cis, ci)
			}
		case model.GotoArticle:
			if art, ok := artm[content.ID]; ok {
				ci.Title = art.Title
				ci.Desc = art.Summary
				if len(art.ImageURLs) != 0 {
					ci.Cover = art.ImageURLs[0]
				}
				ci.Goto = model.GotoArticle
				ci.Param = strconv.FormatInt(art.ID, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, 0, 0, 0, nil)
				if art.Stats != nil {
					ci.fillArtStat(art)
				}
				ci.Badge = "文章"
				cis = append(cis, ci)
			}
		}
	}
	if len(cis) < _convergeMinCount {
		return
	}
	i.Item = cis
	i.Goto = model.GotoConverge
	i.URI = model.FillURI(i.Goto, c.ReValue, c.ReType, 0, 0, nil)
	i.Redirect = model.FillRedirect(i.Goto, c.ReType)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, c.ReType)
	i.Title = c.Title
	i.Cover = c.Cover
	i.Param = strconv.FormatInt(c.ID, 10)
}

func (i *Item) fillArtStat(m *article.Meta) {
	i.Play = int(m.Stats.View)
	i.Reply = int(m.Stats.Reply)
}

func (i *Item) FromGameDownloadS(d *operate.Download, plat int8, build int) {
	i.Title = d.Title
	i.Cover = d.DoubleCover
	i.BigCover = d.Cover
	i.Goto = model.GotoGameDownloadS
	i.ViewType = "bili_game_download_layout"
	i.Desc = d.Desc
	i.URI = model.FillURI(i.Goto, d.URLValue, d.URLType, plat, build, nil)
	i.Redirect = model.FillRedirect(i.Goto, d.URLType)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, d.URLType)
	i.Face = d.Icon
	i.Button = &Button{Name: d.ButtonText, URI: model.FillURI(i.Goto, d.ReValue, d.ReType, plat, build, nil)}
	i.Param = strconv.FormatInt(d.ID, 10)
	i.Download = d.Number
}

func (i *Item) FromGameDownload(d *operate.Download, plat int8, build int) {
	if d.URLValue == "" || d.ReValue == "" {
		return
	}
	i.Title = d.Title
	i.Cover = d.Cover
	i.Goto = model.GotoGameDownload
	i.ViewType = "bili_game_download_layout"
	i.Desc = d.Desc
	i.URI = model.FillURI(i.Goto, d.URLValue, d.URLType, plat, build, nil)
	i.Redirect = model.FillRedirect(i.Goto, d.URLType)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, d.URLType)
	i.Face = d.Icon
	i.Button = &Button{Name: d.ButtonText, URI: model.FillURI(i.Goto, d.ReValue, d.ReType, plat, build, nil)}
	i.Param = strconv.FormatInt(d.ID, 10)
	i.Download = d.Number
}

func (i *Item) FromArticle(m *article.Meta) {
	if m.State < 0 || (m.TemplateID != 3 && m.TemplateID != 4) {
		return
	}
	i.Title = m.Title
	i.Desc = m.Summary
	i.Covers = m.ImageURLs
	i.Goto = model.GotoArticle
	i.Param = strconv.FormatInt(m.ID, 10)
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, nil)
	if m.Author != nil {
		i.Mid = m.Author.Mid
		i.Name = m.Author.Name
		i.Face = m.Author.Face
	}
	if len(m.Categories) >= 2 && m.Categories[0] != nil && m.Categories[1] != nil {
		i.Category = &Category{ID: m.Categories[0].ID, Name: m.Categories[0].Name}
		i.Category.Children = &Category{ID: m.Categories[1].ID, Name: m.Categories[1].Name}
	}
	if m.Stats != nil {
		i.fillArtStat(m)
	}
	// i.Temple = int(m.TemplateID)
	// if i.Temple == 4 {
	// 	i.Temple = 1
	// }
	i.Template = int(m.TemplateID)
	i.BannerURL = m.BannerURL
	i.PTime = m.PublishTime
}

func (i *Item) FromArticleS(m *article.Meta) {
	if m.State < 0 {
		return
	}
	i.Title = m.Title
	i.Desc = m.Summary
	i.Covers = m.ImageURLs
	i.Goto = model.GotoArticleS
	i.Param = strconv.FormatInt(m.ID, 10)
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, nil)
	if m.Author != nil {
		i.Mid = m.Author.Mid
		i.Name = m.Author.Name
		i.Face = m.Author.Face
	}
	if len(m.Categories) >= 2 && m.Categories[0] != nil && m.Categories[1] != nil {
		i.Category = &Category{ID: m.Categories[0].ID, Name: m.Categories[0].Name}
		i.Category.Children = &Category{ID: m.Categories[1].ID, Name: m.Categories[1].Name}
	}
	if m.Stats != nil {
		i.fillArtStat(m)
	}
	// i.Temple = int(m.TemplateID)
	i.Template = int(m.TemplateID)
	i.BannerURL = m.BannerURL
	i.PTime = m.PublishTime
}

func (i *Item) FromShoppingS(c *shopping.Shopping) {
	if c.Name == "" || c.PerformanceImage == "" || c.URL == "" {
		return
	}
	i.Title = c.Name
	if strings.HasPrefix(c.PerformanceImage, "http:") || strings.HasPrefix(c.PerformanceImage, "https:") {
		i.Cover = c.PerformanceImage
	} else {
		i.Cover = "http:" + c.PerformanceImage
	}
	i.Goto = model.GotoShoppingS
	i.URI = model.FillURI(i.Goto, c.URL, 0, 0, 0, nil)
	// i.RedirectURI = model.FillRedirectURI(i.Goto, i.URI, 0)
	i.STime = c.STime
	i.ETime = c.ETime
	i.City = c.CityName
	if len(c.Tags) != 0 {
		i.PType = c.Tags[0].TagName
	}
	i.Param = strconv.FormatInt(c.ID, 10)
	//  竖图
	i.Subtitle = c.Subname
	i.Price = c.Pricelt
	i.Desc = c.Want
	i.Type = c.Type
}

func (i *Item) FromAudio(a *audio.Audio) {
	i.Title = a.Title
	i.Cover = a.CoverURL
	i.Param = strconv.FormatInt(a.MenuID, 10)
	i.Goto = model.GotoAudio
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, 0, nil) + "?from=tianma"
	i.Play = int(a.PlayNum)
	i.Count = a.RecordNum
	i.Fav = int(a.FavoriteNum)
	i.Face = a.Face
	titles := make([]string, 0, len(a.Songs))
	for index, song := range a.Songs {
		if song != nil || song.Title != "" {
			titles = append(titles, song.Title)
			if index == 0 {
				i.SongTitle = song.Title
			}
		}
	}
	i.Desc = strings.Join(titles, ",")
	for _, ctg := range a.Ctgs {
		tag := &channel.Tag{ID: ctg.ItemID, Name: ctg.ItemVal}
		i.Tags = append(i.Tags, tag)
		if len(i.Tags) == 2 {
			break
		}
	}
	if len(a.Ctgs) != 0 {
		id := a.Ctgs[0].ItemID
		name := a.Ctgs[0].ItemVal
		if len(a.Ctgs) > 1 {
			id = a.Ctgs[1].ItemID
			name += "·" + a.Ctgs[1].ItemVal
		}
		i.Tag = &Tag{Name: name, URI: model.FillURI(model.GotoAudioTag, strconv.FormatInt(id, 10), 0, 0, 0, nil) + "?from=tianma"}
	}
	if a.Type == 5 {
		i.Badge = "专辑"
		i.Type = 2
	} else {
		i.Badge = "歌单"
		i.Type = 1
	}
	i.PTime = xtime.Time(a.PaTime)
}

func (i *Item) FromPlayer(a *archive.ArchiveWithPlayer) {
	if !a.IsNormal() {
		return
	}
	title := i.Title
	if title == "" {
		title = a.Title
	}
	cover := i.Cover
	if cover == "" {
		cover = a.Pic
	}
	item := &Item{Title: title, Cover: cover, Param: strconv.FormatInt(a.Aid, 10), Goto: model.GotoAv}
	item.URI = model.FillURI(item.Goto, item.Param, 0, 0, 0, model.AvPlayHandler(a.Archive3, a.PlayerInfo))
	item.fillArcStat(a.Archive3)
	i.Item = []*Item{item}
	i.Cid = a.FirstCid
	i.Rid = a.TypeID
	i.TName = a.TypeName
	i.Mid = a.Author.Mid
	i.Goto = model.GotoPlayer
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.Duration = a.Duration
	i.Autoplay = a.Rights.Autoplay
}

func (i *Item) FromPlayerLive(r *cardlive.Room) {
	if r.LiveStatus == 0 || r.Title == "" || r.Cover == "" {
		return
	}
	i.Name = r.Uname
	i.Mid = r.UID
	i.Face = r.Face
	item := &Item{Title: r.Title, Cover: r.Cover, Param: strconv.FormatInt(r.RoomID, 10), Goto: model.GotoLive, URI: model.FillURI(i.Goto, i.Param, 0, 0, 0, model.LiveRoomHandler(r))}
	item.Online = int(r.Online)
	item.Area2 = &Area2{ID: r.AreaV2ParentID, Name: r.AreaV2ParentName, Children: &Area2{ID: r.AreaV2ID, Name: r.AreaV2Name}}
	i.Item = []*Item{item}
	i.Goto = model.GotoPlayer
	i.Autoplay = 1
}

func (i *Item) FromLiveUpRcmd(id int64, cs []*cardlive.Card, card map[int64]*account.Card) {
	if len(cs) < 2 {
		return
	}
	is := make([]*Item, 0, 2)
	for _, c := range cs[:2] {
		if c.LiveStatus != 1 {
			return
		}
		it := &Item{}
		it.Title = c.Title
		it.Cover = c.ShowCover
		it.Goto = model.GotoLive
		it.Param = strconv.FormatInt(c.RoomID, 10)
		it.URI = model.FillURI(it.Goto, it.Param, 0, 0, 0, model.LiveUpHandler(c))
		it.Fans = int64(c.Online)
		it.Mid = c.UID
		it.Name = c.Uname
		it.Badge = "直播"
		if card, ok := card[it.Mid]; ok {
			if card.Official.Role != 0 {
				it.Official = &OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
			}
		}
		is = append(is, it)
	}
	i.Item = is
	i.Goto = model.GotoLiveUpRcmd
	i.Param = strconv.FormatInt(id, 10)
}

func (i *Item) FromSubscribe(r *operate.Follow, card map[int64]*account.Card, follow map[int64]bool, upStatm map[int64]*relation.Stat, tagm map[int64]*bustag.Tag) {
	if r == nil {
		return
	}
	is := make([]*Item, 0, 3)
	switch r.Type {
	case "upper":
		for _, r := range r.Items {
			item := &Item{}
			if card, ok := card[r.Pid]; ok {
				if f, ok := follow[r.Pid]; ok && f {
					continue
				}
				item.Name = card.Name
				item.Face = card.Face
				item.Mid = card.Mid
				if card.Official.Role != 0 {
					item.Official = &OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
				}
				item.IsAtten = 0
				if stat, ok := upStatm[r.Pid]; ok {
					item.Fans = stat.Follower
				}
				is = append(is, item)
			}
		}
		i.Kind = "upper"
	case "channel_three":
		for _, r := range r.Items {
			item := &Item{}
			tg, ok := tagm[r.Pid]
			if !ok || tg.IsAtten == 1 {
				continue
			}
			item.Name = tg.Name
			item.Face = tg.Cover
			item.Fans = int64(tg.Count.Atten)
			item.IsAtten = tg.IsAtten
			item.Param = strconv.FormatInt(tg.ID, 10)
			if item.Face != "" {
				is = append(is, item)
			}
		}
		i.Kind = "channel"
	}
	if len(is) < 3 {
		return
	}
	i.Item = is[:3]
	i.Title = r.Title
	i.Param = strconv.FormatInt(r.ID, 10)
	i.Goto = model.GotoSubscribe
}

func (i *Item) FromSubscribeChannel(r *recommend.Item, tags map[int64]*tag.Tag) {
	if len(tags) == 0 {
		return
	}
	is := []*Item{}
	for _, item := range r.Items {
		if t, ok := tags[item.ID]; ok {
			if t.IsAtten == 1 {
				continue
			}
			tmp := &Item{
				Name:    t.Name,
				Face:    t.Cover,
				Param:   strconv.FormatInt(t.ID, 10),
				Fans:    int64(t.Count.Atten),
				IsAtten: t.IsAtten,
			}
			is = append(is, tmp)
		}
	}
	if len(is) < 3 {
		return
	}
	i.Item = is[:3]
	if r.Config != nil {
		i.Title = r.Config.Title
	}
	i.Param = strconv.FormatInt(r.ID, 10)
	i.Kind = "channel"
	i.Goto = model.GotoSubscribe
}

func (i *Item) FromChannelRcmd(r *operate.Follow, am map[int64]*archive.ArchiveWithPlayer, tagm map[int64]*bustag.Tag) {
	if r == nil {
		return
	}
	if a, ok := am[r.Pid]; ok {
		i.Goto = model.GotoChannelRcmd
		i.URI = model.FillURI(model.GotoAv, strconv.FormatInt(a.Aid, 10), 0, 0, 0, model.AvPlayHandler(a.Archive3, a.PlayerInfo))
		i.Title = a.Title
		i.Cover = a.Pic
		if tag, ok := tagm[r.Tid]; ok {
			i.Tag = &Tag{ID: tag.ID, Name: tag.Name, Face: tag.Cover, Fans: tag.Count.Atten, IsAtten: tag.IsAtten}
		}
		i.Cid = a.FirstCid
		i.Autoplay = a.Rights.Autoplay
		i.fillArcStat(a.Archive3)
		i.Duration = a.Duration
		// TODO 等待开启
		// percent := i.Like / (i.Like + i.Dislike) * 100
		// if percent != 0 {
		// 	i.Desc = strconv.Itoa(percent) + "%的人推荐"
		// }
		i.Param = strconv.FormatInt(r.ID, 10)
	}
}
