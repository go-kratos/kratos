package search

import (
	"bytes"
	"fmt"

	// "hash/crc32"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/interface/main/app-intl/model"
	bgmmdl "go-common/app/interface/main/app-intl/model/bangumi"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// search const
var getHightLight = regexp.MustCompile(`<em.*?em>`)

// Result struct
type Result struct {
	Trackid   string     `json:"trackid,omitempty"`
	Page      int        `json:"page,omitempty"`
	NavInfo   []*NavInfo `json:"nav,omitempty"`
	Item      []*Item    `json:"item,omitempty"`
	Array     int        `json:"array,omitempty"`
	Attribute int32      `json:"attribute"`
	EasterEgg *EasterEgg `json:"easter_egg,omitempty"`
}

// NavInfo struct
type NavInfo struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Type  int    `json:"type"`
	Show  int    `json:"show_more,omitempty"`
}

// TypeSearch struct
type TypeSearch struct {
	TrackID string  `json:"trackid"`
	Pages   int     `json:"pages"`
	Total   int     `json:"total"`
	Items   []*Item `json:"items,omitempty"`
}

// Suggestion struct
type Suggestion struct {
	TrackID string      `json:"trackid"`
	UpUser  interface{} `json:"upuser,omitempty"`
	Bangumi interface{} `json:"bangumi,omitempty"`
	Suggest []string    `json:"suggest,omitempty"`
}

// SuggestionResult3 struct
type SuggestionResult3 struct {
	TrackID string  `json:"trackid"`
	List    []*Item `json:"list,omitempty"`
}

// NoResultRcndResult struct
type NoResultRcndResult struct {
	TrackID string  `json:"trackid"`
	Title   string  `json:"title,omitempty"`
	Pages   int     `json:"pages"`
	Items   []*Item `json:"items,omitempty"`
}

// EasterEgg struct
type EasterEgg struct {
	ID        int64 `json:"id,omitempty"`
	ShowCount int   `json:"show_count,omitempty"`
}

// Item struct
type Item struct {
	TrackID        string `json:"trackid,omitempty"`
	LinkType       string `json:"linktype,omitempty"`
	Position       int    `json:"position,omitempty"`
	SuggestKeyword string `json:"suggest_keyword,omitempty"`
	Title          string `json:"title,omitempty"`
	Name           string `json:"name,omitempty"`
	Cover          string `json:"cover,omitempty"`
	URI            string `json:"uri,omitempty"`
	Param          string `json:"param,omitempty"`
	Goto           string `json:"goto,omitempty"`
	// av
	Play     int        `json:"play,omitempty"`
	Danmaku  int        `json:"danmaku,omitempty"`
	Author   string     `json:"author,omitempty"`
	ViewType string     `json:"view_type,omitempty"`
	PTime    xtime.Time `json:"ptime,omitempty"`
	RecTags  []string   `json:"rec_tags,omitempty"`
	// bangumi season
	SeasonID     int64   `json:"season_id,omitempty"`
	SeasonType   int     `json:"season_type,omitempty"`
	Finish       int8    `json:"finish,omitempty"`
	Started      int8    `json:"started,omitempty"`
	Index        string  `json:"index,omitempty"`
	NewestCat    string  `json:"newest_cat,omitempty"`
	NewestSeason string  `json:"newest_season,omitempty"`
	CatDesc      string  `json:"cat_desc,omitempty"`
	TotalCount   int     `json:"total_count,omitempty"`
	MediaType    int     `json:"media_type,omitempty"`
	PlayState    int     `json:"play_state,omitempty"`
	Style        string  `json:"style,omitempty"`
	CV           string  `json:"cv,omitempty"`
	Rating       float64 `json:"rating,omitempty"`
	Vote         int     `json:"vote,omitempty"`
	RatingCount  int     `json:"rating_count,omitempty"`
	BadgeType    int     `json:"badge_type,omitempty"`
	// upper
	Sign           string          `json:"sign,omitempty"`
	Fans           int             `json:"fans,omitempty"`
	Level          int             `json:"level,omitempty"`
	Desc           string          `json:"desc,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
	AvItems        []*Item         `json:"av_items,omitempty"`
	Item           []*Item         `json:"item,omitempty"`
	CTime          int64           `json:"ctime,omitempty"`
	IsUp           bool            `json:"is_up,omitempty"`
	LiveURI        string          `json:"live_uri,omitempty"`
	// movie
	ScreenDate string `json:"screen_date,omitempty"`
	Area       string `json:"area,omitempty"`
	CoverMark  string `json:"cover_mark,omitempty"`
	// arc and sp
	Arcs int `json:"archives,omitempty"`
	// arc and movie
	Duration    string `json:"duration,omitempty"`
	DurationInt int64  `json:"duration_int,omitempty"`
	Actors      string `json:"actors,omitempty"`
	Staff       string `json:"staff,omitempty"`
	Length      int    `json:"length,omitempty"`
	Status      int    `json:"status,omitempty"`
	// live
	RoomID      int64  `json:"roomid,omitempty"`
	Mid         int64  `json:"mid,omitempty"`
	Type        string `json:"type,omitempty"`
	Attentions  int    `json:"attentions,omitempty"`
	LiveStatus  int    `json:"live_status,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Region      int    `json:"region,omitempty"`
	Online      int    `json:"online,omitempty"`
	ShortID     int    `json:"short_id,omitempty"`
	CateName    string `json:"area_v2_name,omitempty"`
	IsSelection int    `json:"is_selection,omitempty"`
	// article
	ID         int64    `json:"id,omitempty"`
	TemplateID int      `json:"template_id,omitempty"`
	ImageUrls  []string `json:"image_urls,omitempty"`
	View       int      `json:"view,omitempty"`
	Like       int      `json:"like,omitempty"`
	Reply      int      `json:"reply,omitempty"`
	// special
	Badge      string      `json:"badge,omitempty"`
	RcmdReason *RcmdReason `json:"rcmd_reason,omitempty"`
	// media bangumi and mdeia ft
	Prompt   string  `json:"prompt,omitempty"`
	Episodes []*Item `json:"episodes,omitempty"`
	Label    string  `json:"label,omitempty"`
	// game
	Reserve string `json:"reserve,omitempty"`
	// user
	Face string `json:"face,omitempty"`
	// suggest
	From      string  `json:"from,omitempty"`
	KeyWord   string  `json:"keyword,omitempty"`
	CoverSize float64 `json:"cover_size,omitempty"`
	SugType   string  `json:"sug_type,omitempty"`
	TermType  int     `json:"term_type,omitempty"`
	// rcmd query
	List       []*Item `json:"list,omitempty"`
	FromSource string  `json:"from_source,omitempty"`
	// live master
	UCover         string `json:"ucover,omitempty"`
	VerifyType     int    `json:"verify_type,omitempty"`
	VerifyDesc     string `json:"verify_desc,omitempty"`
	LevelColor     int64  `json:"level_color,omitempty"`
	IsAttention    int    `json:"is_atten,omitempty"`
	CateParentName string `json:"cate_parent_name,omitempty"`
	CateNameNew    string `json:"cate_name,omitempty"`
	Glory          *Glory `json:"glory_info,omitempty"`
	// twitter
	Covers     []string `json:"covers,omitempty"`
	CoverCount int      `json:"cover_count,omitempty"`
}

// Glory live struct
type Glory struct {
	Title string  `json:"title,omitempty"`
	Total int     `json:"total"`
	Items []*Item `json:"items,omitempty"`
}

// RcmdReason struct
type RcmdReason struct {
	Content string `json:"content,omitempty"`
}

// UserResult struct
type UserResult struct {
	Items []*Item `json:"items,omitempty"`
}

// DefaultWords struct
type DefaultWords struct {
	Trackid string `json:"trackid,omitempty"`
	Param   string `json:"param,omitempty"`
	Show    string `json:"show,omitempty"`
	Word    string `json:"word,omitempty"`
}

// FromSeason form func
func (i *Item) FromSeason(b *Bangumi, bangumi string) {
	i.Title = b.Title
	i.Cover = b.Cover
	i.Goto = model.GotoBangumi
	i.Param = strconv.Itoa(int(b.SeasonID))
	i.URI = model.FillURI(bangumi, i.Param, nil)
	i.Finish = int8(b.IsFinish)
	i.Started = int8(b.IsStarted)
	i.Index = b.NewestEpIndex
	i.NewestCat = b.NewestCat
	i.NewestSeason = b.NewestSeason
	i.TotalCount = b.TotalCount
	var buf bytes.Buffer
	if b.CatList.TV != 0 {
		buf.WriteString(`TV(`)
		buf.WriteString(strconv.Itoa(b.CatList.TV))
		buf.WriteString(`) `)
	}
	if b.CatList.Movie != 0 {
		buf.WriteString(`剧场版(`)
		buf.WriteString(strconv.Itoa(b.CatList.Movie))
		buf.WriteString(`) `)
	}
	if b.CatList.Ova != 0 {
		buf.WriteString(`OVA/OAD/SP(`)
		buf.WriteString(strconv.Itoa(b.CatList.Ova))
		buf.WriteString(`)`)
	}
	i.CatDesc = buf.String()
}

// FromUser form func
func (i *Item) FromUser(u *User, as map[int64]*api.Arc) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = model.GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Mid = u.Mid
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	i.RoomID = u.RoomID
	if u.IsUpuser == 1 {
		for _, v := range u.Res {
			vi := &Item{}
			vi.Title = v.Title
			vi.Cover = v.Pic
			vi.Goto = model.GotoAv
			vi.Param = strconv.Itoa(int(v.Aid))
			vi.URI = model.FillURI(vi.Goto, vi.Param, model.AvHandler(as[v.Aid], "", nil))
			a, ok := as[v.Aid]
			if ok {
				vi.Play = int(a.Stat.View)
				vi.Danmaku = int(a.Stat.Danmaku)
			} else {
				switch play := v.Play.(type) {
				case float64:
					vi.Play = int(play)
				case string:
					vi.Play, _ = strconv.Atoi(play)
				}
				vi.Danmaku = v.Danmaku
			}
			vi.CTime = v.Pubdate
			vi.Duration = v.Duration
			i.AvItems = append(i.AvItems, vi)
		}
		i.IsUp = true
	}
}

// FromUpUser form func
func (i *Item) FromUpUser(u *User, as map[int64]*api.Arc) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = model.GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Mid = u.Mid
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	for _, v := range u.Res {
		vi := &Item{}
		vi.Title = v.Title
		vi.Cover = v.Pic
		vi.Goto = model.GotoAv
		vi.Param = strconv.Itoa(int(v.Aid))
		vi.URI = model.FillURI(vi.Goto, vi.Param, model.AvHandler(as[v.Aid], "", nil))
		a, ok := as[v.Aid]
		if ok {
			vi.Play = int(a.Stat.View)
			vi.Danmaku = int(a.Stat.Danmaku)
		} else {
			switch play := v.Play.(type) {
			case float64:
				vi.Play = int(play)
			case string:
				vi.Play, _ = strconv.Atoi(play)
			}
			vi.Danmaku = v.Danmaku
		}
		vi.CTime = v.Pubdate
		vi.Duration = v.Duration
		i.AvItems = append(i.AvItems, vi)
	}
	i.RoomID = u.RoomID
	i.IsUp = u.IsUpuser == 1
}

// FromMovie form func
func (i *Item) FromMovie(m *Movie, as map[int64]*api.Arc) {
	i.Title = m.Title
	i.Desc = m.Desc
	if m.Type == "movie" {
		i.Cover = m.Cover
		i.Param = strconv.Itoa(int(m.Aid))
		i.Goto = model.GotoAv
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(as[m.Aid], "", nil))
		i.CoverMark = model.StatusMark(m.Status)
	} else if m.Type == "special" {
		i.Param = m.SpID
		i.Goto = model.GotoSp
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		i.Cover = m.Pic
	}
	i.Staff = m.Staff
	i.Actors = m.Actors
	i.Area = m.Area
	i.Length = m.Length
	i.Status = m.Status
	i.ScreenDate = m.ScreenDate
}

// FromVideo form func
func (i *Item) FromVideo(v *Video, a *api.Arc) {
	i.Title = v.Title
	i.Cover = v.Pic
	i.Author = v.Author
	i.Param = strconv.Itoa(int(v.ID))
	i.Goto = model.GotoAv
	if a != nil {
		i.Face = a.Author.Face
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(a, "", nil))
		i.Play = int(a.Stat.View)
		i.Danmaku = int(a.Stat.Danmaku)
	} else {
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		switch play := v.Play.(type) {
		case float64:
			i.Play = int(play)
		case string:
			i.Play, _ = strconv.Atoi(play)
		}
		i.Danmaku = v.Danmaku
	}
	i.Desc = v.Desc
	i.Duration = v.Duration
	i.ViewType = v.ViewType
	i.RecTags = v.RecTags
}

// FromOperate form func
func (i *Item) FromOperate(o *Operate, gt string) {
	i.Title = o.Title
	i.Cover = o.Cover
	i.URI = o.RedirectURL
	i.Param = strconv.FormatInt(o.ID, 10)
	i.Desc = o.Desc
	i.Badge = o.Corner
	i.Goto = gt
	if o.RecReason != "" {
		i.RcmdReason = &RcmdReason{Content: o.RecReason}
	}
}

// FromConverge form func
func (i *Item) FromConverge(o *Operate, am map[int64]*api.Arc, artm map[int64]*article.Meta) {
	const _convergeMinCount = 2
	cis := make([]*Item, 0, len(o.ContentList))
	for _, c := range o.ContentList {
		ci := &Item{}
		switch c.Type {
		case 0:
			if a, ok := am[c.ID]; ok && a.IsNormal() {
				ci.Title = a.Title
				ci.Cover = a.Pic
				ci.Goto = model.GotoAv
				ci.Param = strconv.FormatInt(a.Aid, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, model.AvHandler(a, "", nil))
				ci.fillArcStat(a)
				cis = append(cis, ci)
			}
		case 2:
			if art, ok := artm[c.ID]; ok {
				ci.Title = art.Title
				ci.Desc = art.Summary
				if len(art.ImageURLs) != 0 {
					ci.Cover = art.ImageURLs[0]
				}
				ci.Goto = model.GotoArticle
				ci.Param = strconv.FormatInt(art.ID, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, nil)
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
	i.Title = o.Title
	i.Cover = o.Cover
	i.URI = o.RedirectURL
	i.Param = strconv.FormatInt(o.ID, 10)
	i.Desc = o.Desc
	i.Badge = o.Corner
	i.Goto = model.GotoConverge
	if o.RecReason != "" {
		i.RcmdReason = &RcmdReason{Content: o.RecReason}
	}
}

// FromMedia form func
func (i *Item) FromMedia(m *Media, prompt string, gt string, bangumis map[string]*bgmmdl.Card) {
	i.Title = m.Title
	if i.Title == "" {
		i.Title = m.OrgTitle
	}
	i.Cover = m.Cover
	i.Goto = gt
	i.Param = strconv.Itoa(int(m.MediaID))
	i.URI = m.GotoURL
	i.MediaType = m.MediaType
	i.PlayState = m.PlayState
	i.Style = m.Styles
	i.CV = m.CV
	i.Staff = m.Staff
	if m.MediaScore != nil {
		i.Rating = m.MediaScore.Score
		i.Vote = m.MediaScore.UserCount
	}
	i.PTime = m.Pubtime
	areas := strings.Split(m.Areas, "、")
	if len(areas) != 0 {
		i.Area = areas[0]
	}
	i.Prompt = prompt
	var hit string
	for _, v := range m.HitColumns {
		if v == "cv" {
			hit = v
			break
		} else if v == "staff" {
			hit = v
		}
	}
	if hit == "cv" {
		for _, v := range getHightLight.FindAllStringSubmatch(m.CV, -1) {
			if gt == model.GotoBangumi {
				i.Label = fmt.Sprintf("声优: %v...", v[0])
				break
			} else if gt == model.GotoMovie {
				i.Label = fmt.Sprintf("演员: %v...", v[0])
				break
			}
		}
	} else if hit == "staff" {
		for _, v := range getHightLight.FindAllStringSubmatch(m.Staff, -1) {
			i.Label = fmt.Sprintf("制作人员: %v...", v[0])
			break
		}
	}
	// get from PGC API.
	i.SeasonID = m.SeasonID
	ssID := strconv.Itoa(int(m.SeasonID))
	if bgm, ok := bangumis[ssID]; ok {
		i.IsAttention = bgm.IsFollow
		i.IsSelection = bgm.IsSelection
		i.SeasonType = bgm.SeasonType
		for _, v := range bgm.Episodes {
			tmp := &Item{
				Param: strconv.Itoa(int(v.ID)),
				Index: v.Index,
			}
			tmp.URI = model.FillURI(model.GotoEP, tmp.Param, nil)
			i.Episodes = append(i.Episodes, tmp)
		}
	}
}

// FromArticle form func
func (i *Item) FromArticle(a *Article) {
	i.ID = a.ID
	i.Mid = a.Mid
	i.TemplateID = a.TemplateID
	i.Title = a.Title
	i.Desc = a.Desc
	i.ImageUrls = a.ImageUrls
	i.View = a.View
	i.Play = a.View
	i.Like = a.Like
	i.Reply = a.Reply
	i.Goto = model.GotoArticle
	i.Param = strconv.Itoa(int(a.ID))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

// FromChannel form func
func (i *Item) FromChannel(c *Channel) {
	i.ID = c.TagID
	i.Title = c.TagName
	i.Cover = c.Cover
	i.Param = strconv.FormatInt(c.TagID, 10)
	i.Goto = model.GotoChannel
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Type = c.Type
	i.Attentions = c.AttenCount
}

// FromQuery form func
func (i *Item) FromQuery(qs []*Query) {
	i.Goto = model.GotoRecommendWord
	for _, q := range qs {
		i.List = append(i.List, &Item{Param: strconv.FormatInt(q.ID, 10), Title: q.Name, Type: q.Type, FromSource: q.FromSource})
	}
}

// FromTwitter form twitter
func (i *Item) FromTwitter(t *Twitter) {
	i.Title = t.Content
	i.Covers = t.Cover
	i.CoverCount = t.CoverCount
	i.Param = strconv.FormatInt(t.ID, 10)
	i.Goto = model.GotoTwitter
	i.URI = model.FillURI(i.Goto, strconv.FormatInt(t.PicID, 10), nil)
}

// fillArcStat fill func
func (i *Item) fillArcStat(a *api.Arc) {
	if a.Access == 0 {
		i.Play = int(a.Stat.View)
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Reply = int(a.Stat.Reply)
	i.Like = int(a.Stat.Like)
}

// fillArtStat fill func
func (i *Item) fillArtStat(m *article.Meta) {
	i.Play = int(m.Stats.View)
	i.Reply = int(m.Stats.Reply)
}

// FromSuggest3 form func
func (i *Item) FromSuggest3(st *Sug, as map[int64]*api.Arc) {
	i.From = "search"
	i.Title = st.ShowName
	i.KeyWord = st.Term
	i.Position = st.Pos
	i.Cover = st.Cover
	i.CoverSize = st.CoverSize
	i.SugType = st.SubType
	i.TermType = st.TermType
	if st.TermType == SuggestionJump {
		switch st.SubType {
		case SuggestionAV:
			i.Goto = model.GotoAv
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), model.AvHandler(as[st.Ref], "", nil))
			i.SugType = "视频"
		}
	}
}
