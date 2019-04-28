package card

import (
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/show"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/log"
)

// ButtonStatus is
type ButtonStatus struct {
	Text    string
	Goto    model.Gt
	Param   string
	IsAtten int8
	Type    model.Type
	Event   model.Event
}

// AvatarStatus is
type AvatarStatus struct {
	Cover string
	Text  string
	Goto  model.Gt
	Param string
	Type  model.Type
}

// Base is
type Base struct {
	CardType   model.CardType           `json:"card_type,omitempty"`
	CardGoto   model.CardGt             `json:"card_goto,omitempty"`
	Goto       model.Gt                 `json:"goto,omitempty"`
	Param      string                   `json:"param,omitempty"`
	Cover      string                   `json:"cover,omitempty"`
	Title      string                   `json:"title,omitempty"`
	URI        string                   `json:"uri,omitempty"`
	DescButton *Button                  `json:"desc_button,omitempty"`
	ThreePoint *ThreePoint              `json:"three_point,omitempty"`
	Args       Args                     `json:"args,omitempty"`
	PlayerArgs *PlayerArgs              `json:"player_args,omitempty"`
	Idx        int64                    `json:"idx,omitempty"`
	AdInfo     *cm.AdInfo               `json:"ad_info,omitempty"`
	Right      bool                     `json:"-"`
	Rcmd       *ai.Item                 `json:"-"`
	Tagm       map[int64]*tag.Tag       `json:"-"`
	IsAttenm   map[int64]int8           `json:"-"`
	Statm      map[int64]*relation.Stat `json:"-"`
	Cardm      map[int64]*account.Card  `json:"-"`
	CardLen    int                      `json:"-"`
	Columnm    model.ColumnStatus       `json:"-"`
	FromType   string                   `json:"from_type,omitempty"`
}

// ThreePoint is
type ThreePoint struct {
	DislikeReasons []*DislikeReason `json:"dislike_reasons,omitempty"`
	Feedbacks      []*Feedback      `json:"feedbacks,omitempty"`
	WatchLater     int8             `json:"watch_later,omitempty"`
}

func (c *Base) from(param, cover, title string, gt model.Gt, uri string, f func(uri string) string) {
	c.URI = model.FillURI(gt, uri, f)
	c.Cover = cover
	c.Title = title
	if gt != "" {
		c.Goto = gt
	} else {
		c.Goto = model.Gt(c.CardGoto)
	}
	c.Param = param
}

// Handler is
type Handler interface {
	From(main interface{}, op *operate.Card)
	Get() *Base
}

// Handle is
func Handle(plat int8, cardGoto model.CardGt, cardType model.CardType, column model.ColumnStatus, rcmd *ai.Item, tagm map[int64]*tag.Tag, isAttenm map[int64]int8, statm map[int64]*relation.Stat, cardm map[int64]*account.Card) (hander Handler) {
	if model.IsIPad(plat) {
		return ipadHandle(cardGoto, cardType, rcmd, nil, isAttenm, statm, cardm)
	}
	switch model.Columnm[column] {
	case model.ColumnSvrSingle:
		return singleHandle(cardGoto, cardType, rcmd, tagm, isAttenm, statm, cardm)
	case model.ColumnSvrDouble:
		return doubleHandle(cardGoto, cardType, rcmd, tagm, isAttenm, statm, cardm)
	}
	return
}

// SwapTwoItem is
func SwapTwoItem(rs []Handler, i Handler) (is []Handler) {
	is = append(rs, rs[len(rs)-1])
	is[len(is)-2] = i
	return
}

func SwapThreeItem(rs []Handler, i Handler) (is []Handler) {
	is = append(rs, rs[len(rs)-1])
	is[len(is)-2] = i
	is[len(is)-3], is[len(is)-2] = is[len(is)-2], is[len(is)-3]
	return
}

func SwapFourItem(rs []Handler, i Handler) (is []Handler) {
	is = append(rs, rs[len(rs)-1])
	is[len(is)-2] = i
	is[len(is)-3], is[len(is)-2] = is[len(is)-2], is[len(is)-3]
	is[len(is)-4], is[len(is)-3] = is[len(is)-3], is[len(is)-4]
	return
}

// TopBottomRcmdReason is
func TopBottomRcmdReason(r *ai.RcmdReason, isAtten int8, cardm map[int64]*account.Card) (topRcmdReason, bottomRcomdReason string) {
	if r == nil {
		if isAtten == 1 {
			bottomRcomdReason = "已关注"
		}
		return
	}
	switch r.Style {
	case 3:
		if isAtten != 1 {
			return
		}
		bottomRcomdReason = r.Content
	case 4:
		_, ok := cardm[r.FollowedMid]
		if !ok {
			return
		}
		topRcmdReason = "关注的人赞过"
	default:
		topRcmdReason = r.Content
	}
	return
}

// Button is
type Button struct {
	Text     string      `json:"text,omitempty"`
	Param    string      `json:"param,omitempty"`
	URI      string      `json:"uri,omitempty"`
	Event    model.Event `json:"event,omitempty"`
	Selected int8        `json:"selected,omitempty"`
	Type     model.Type  `json:"type,omitempty"`
}

func buttonFrom(v interface{}, plat int8) (button *Button) {
	switch v.(type) {
	case *tag.Tag:
		t := v.(*tag.Tag)
		if t != nil {
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  t.Name,
				URI:   model.FillURI(model.GotoTag, strconv.FormatInt(t.ID, 10), nil),
				Event: model.EventChannelClick,
			}
		}
	case []*audio.Ctg:
		ctgs := v.([]*audio.Ctg)
		if len(ctgs) > 1 {
			var name string
			if ctgs[0] != nil {
				name = ctgs[0].ItemVal
				if ctgs[1] != nil {
					name += " · " + ctgs[1].ItemVal
				}
			}
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  name,
				URI:   model.FillURI(model.GotoAudioTag, "", model.AudioTagHandler(ctgs)),
				Event: model.EventChannelClick,
			}
		}
	case []*article.Category:
		ctgs := v.([]*article.Category)
		if len(ctgs) > 1 {
			var name string
			if ctgs[0] != nil {
				name = ctgs[0].Name
				if ctgs[1] != nil {
					name += " · " + ctgs[1].Name
				}
			}
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  name,
				URI:   model.FillURI(model.GotoArticleTag, "", model.ArticleTagHandler(ctgs, plat)),
				Event: model.EventChannelClick,
			}
		}
	case *live.Room:
		r := v.(*live.Room)
		if r != nil {
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  r.AreaV2Name,
				URI:   model.FillURI(model.GotoLiveTag, strconv.FormatInt(r.AreaV2ParentID, 10), model.LiveRoomTagHandler(r)),
				Event: model.EventChannelClick,
			}
		}
	case *live.Card:
		card := v.(*live.Card)
		if card != nil {
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  card.Uname,
				URI:   model.FillURI(model.GotoMid, strconv.FormatInt(card.UID, 10), nil),
				Event: model.EventUpClick,
			}
		}
	case *bplus.Picture:
		p := v.(*bplus.Picture)
		if p != nil {
			if len(p.Topics) == 0 {
				return
			}
			button = &Button{
				Type:  model.ButtonGrey,
				Text:  p.Topics[0],
				URI:   model.FillURI(model.GotoPictureTag, p.Topics[0], nil),
				Event: model.EventChannelClick,
			}
		}
	case *ButtonStatus:
		b := v.(*ButtonStatus)
		if b != nil {
			event, ok := model.ButtonEvent[b.Goto]
			if ok {
				button = &Button{
					Text:     model.ButtonText[b.Goto],
					Event:    event,
					Selected: b.IsAtten,
					Type:     model.ButtonTheme,
				}
			} else {
				button = &Button{
					Text:  b.Text,
					Param: b.Param,
					URI:   model.FillURI(b.Goto, b.Param, nil),
				}
				if b.Event != "" {
					button.Event = b.Event
				} else {
					button.Event = model.EventChannelClick
				}
				if b.Type != 0 {
					button.Type = b.Type
				} else {
					button.Type = model.ButtonGrey
				}
			}
		}
	case nil:
	default:
		log.Warn("buttonFrom: unexpected type %T", v)
	}
	return
}

// Avatar is
type Avatar struct {
	Cover string      `json:"cover,omitempty"`
	Text  string      `json:"text,omitempty"`
	URI   string      `json:"uri,omitempty"`
	Type  model.Type  `json:"type,omitempty"`
	Event model.Event `json:"event,omitempty"`
}

func avatarFrom(status *AvatarStatus) (avatar *Avatar) {
	if status == nil {
		return
	}
	avatar = &Avatar{
		Cover: status.Cover,
		Text:  status.Text,
		URI:   model.FillURI(status.Goto, status.Param, nil),
		Type:  status.Type,
		Event: model.AvatarEvent[status.Goto],
	}
	return
}

// DislikeReason is
type DislikeReason struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Feedback is
type Feedback struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ThreePointFrom is
func (c *Base) ThreePointFrom() {
	const (
		_noSeason = 1
		_region   = 2
		_channel  = 3
		_upper    = 4
	)
	if c.CardGoto == model.CardGotoLogin || c.CardGoto == model.CardGotoBanner || c.CardGoto == model.CardGotoRank || c.CardGoto == model.CardGotoConverge || c.CardGoto == model.CardGotoBangumiRcmd || c.CardGoto == model.CardGotoInterest || c.CardGoto == model.CardGotoFollowMode {
		return
	}
	c.ThreePoint = &ThreePoint{}
	if c.CardGoto == model.CardGotoAv || c.CardGoto == model.CardGotoPlayer || c.CardGoto == model.CardGotoUpRcmdAv || c.CardGoto == model.CardGotoChannelRcmd {
		dislikeReasons := make([]*DislikeReason, 0, 4)
		if c.Args.UpName != "" {
			dislikeReasons = append(dislikeReasons, &DislikeReason{ID: _upper, Name: "UP主:" + c.Args.UpName})
		}
		if c.Args.Rname != "" {
			dislikeReasons = append(dislikeReasons, &DislikeReason{ID: _region, Name: "分区:" + c.Args.Rname})
		}
		if c.Args.Tname != "" {
			dislikeReasons = append(dislikeReasons, &DislikeReason{ID: _channel, Name: "频道:" + c.Args.Tname})
		}
		c.ThreePoint.DislikeReasons = append(dislikeReasons, &DislikeReason{ID: _noSeason, Name: "不感兴趣"})
		c.ThreePoint.Feedbacks = []*Feedback{{ID: 1, Name: "恐怖血腥"}, {ID: 2, Name: "色情低俗"}, {ID: 3, Name: "封面恶心"}, {ID: 4, Name: "标题党/封面党"}}
		c.ThreePoint.WatchLater = 1
	} else {
		c.ThreePoint.DislikeReasons = []*DislikeReason{{ID: _noSeason, Name: "不感兴趣"}}
	}
}

// ThreePointChannel is
func (c *Base) ThreePointChannel() {
	const (
		_noSeason = 1
		_upper    = 4
	)
	if c.CardGoto == model.CardGotoAv || c.CardGoto == model.CardGotoPlayer || c.CardGoto == model.CardGotoUpRcmdAv {
		c.ThreePoint = &ThreePoint{}
		if c.Args.UpName != "" {
			c.ThreePoint.DislikeReasons = append(c.ThreePoint.DislikeReasons, &DislikeReason{ID: _upper, Name: "UP主:" + c.Args.UpName})
		}
		c.ThreePoint.DislikeReasons = append(c.ThreePoint.DislikeReasons, &DislikeReason{ID: _noSeason, Name: "不感兴趣"})
		c.ThreePoint.WatchLater = 1
	}
}

// ThreePointWatchLater is
func (c *Base) ThreePointWatchLater() {
	if c.CardGoto == model.CardGotoAv || c.CardGoto == model.CardGotoPlayer || c.CardGoto == model.CardGotoUpRcmdAv || c.Goto == model.GotoAv {
		c.ThreePoint = &ThreePoint{}
		c.ThreePoint.WatchLater = 1
	}
}

// TabThreePointWatchLater is
func (c *Base) TabThreePointWatchLater() {
	if c.Goto == model.GotoAv && c.CardGoto != model.CardGotoPlayer {
		c.ThreePoint = &ThreePoint{}
		c.ThreePoint.WatchLater = 1
	}
}

// Args is
type Args struct {
	Type   int8   `json:"type,omitempty"`
	UpID   int64  `json:"up_id,omitempty"`
	UpName string `json:"up_name,omitempty"`
	Rid    int32  `json:"rid,omitempty"`
	Rname  string `json:"rname,omitempty"`
	Tid    int64  `json:"tid,omitempty"`
	Tname  string `json:"tname,omitempty"`
}

func (c *Args) fromShopping(s *show.Shopping) {
	c.Type = s.Type
}

func (c *Args) fromArchive(a *archive.Archive3, t *tag.Tag) {
	if a != nil {
		c.UpID = a.Author.Mid
		c.UpName = a.Author.Name
		c.Rid = a.TypeID
		c.Rname = a.TypeName
	}
	if t != nil {
		c.Tid = t.ID
		c.Tname = t.Name
	}
}

func (c *Args) fromLiveRoom(r *live.Room) {
	if r == nil {
		return
	}
	c.UpID = r.UID
	c.UpName = r.Uname
	c.Rid = int32(r.AreaV2ParentID)
	c.Rname = r.AreaV2ParentName
	c.Tid = r.AreaV2ID
	c.Tname = r.AreaV2Name
}

func (c *Args) fromLiveUp(card *live.Card) {
	if card == nil {
		return
	}
	c.UpID = card.UID
	c.UpName = card.Uname
}

func (c *Args) fromAudio(a *audio.Audio) {
	if a == nil {
		return
	}
	c.Type = a.Type
	if len(a.Ctgs) != 0 {
		c.Rid = int32(a.Ctgs[0].ItemID)
		c.Rname = a.Ctgs[0].ItemVal
		if len(a.Ctgs) > 1 {
			c.Tid = a.Ctgs[1].ItemID
			c.Tname = a.Ctgs[1].ItemVal
		}
	}
}

func (c *Args) fromArticle(m *article.Meta) {
	if m == nil {
		return
	}
	if m.Author != nil {
		c.UpID = m.Author.Mid
		c.UpName = m.Author.Name
	}
	if len(m.Categories) != 0 {
		if m.Categories[0] != nil {
			c.Rid = int32(m.Categories[0].ID)
			c.Rname = m.Categories[0].Name
		}
		if len(m.Categories) > 1 {
			if m.Categories[1] != nil {
				c.Tid = m.Categories[1].ID
				c.Tname = m.Categories[1].Name
			}
		}
	}
}

// PlayerArgs is
type PlayerArgs struct {
	IsLive int8  `json:"is_live,omitempty"`
	Aid    int64 `json:"aid,omitempty"`
	Cid    int64 `json:"cid,omitempty"`
	RoomID int64 `json:"room_id,omitempty"`
}

func playerArgsFrom(v interface{}) (playerArgs *PlayerArgs) {
	switch v.(type) {
	case *archive.Archive3:
		a := v.(*archive.Archive3)
		if a == nil || (a.AttrVal(archive.AttrBitIsPGC) == archive.AttrNo && a.Rights.Autoplay != 1) || (a.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes && a.AttrVal(archive.AttrBitBadgepay) == archive.AttrYes) {
			return
		}
		playerArgs = &PlayerArgs{Aid: a.Aid, Cid: a.FirstCid}
	case *live.Room:
		r := v.(*live.Room)
		if r == nil || r.LiveStatus != 1 {
			return
		}
		playerArgs = &PlayerArgs{RoomID: r.RoomID, IsLive: 1}
	case nil:
	default:
		log.Warn("playerArgsFrom: unexpected type %T", v)
	}
	return
}

// rcmdReason
func rcmdReason(r *ai.RcmdReason, name string, isAtten int8, cardm map[int64]*account.Card) (rcmdReason, desc string) {
	// "rcmd_reason":{"content":"已关注","font":1,"grounding":"yellow","id":3,"position":"bottom","style":3}
	if r == nil {
		if isAtten == 1 {
			rcmdReason = "已关注"
			desc = name
		}
		return
	}
	switch r.Style {
	case 3:
		if isAtten != 1 {
			return
		}
		rcmdReason = r.Content
		desc = name
	case 4:
		_, ok := cardm[r.FollowedMid]
		if !ok {
			return
		}
		if r.Content == "" {
			r.Content = "关注的人赞过"
		}
		rcmdReason = r.Content
	default:
		rcmdReason = r.Content
	}
	return
}

// ReasonStyle reason style
type ReasonStyle struct {
	Text              string `json:"text,omitempty"`
	TextColor         string `json:"text_color,omitempty"`
	BgColor           string `json:"bg_color,omitempty"`
	BorderColor       string `json:"border_color,omitempty"`
	BgStyle           int8   `json:"bg_style,omitempty"`
	NightAlphaPercent int    `json:"night_alpha_percent,omitempty"`
}

func topReasonStyleFrom(rcmd *ai.Item, text string, gt model.Gt) (res *ReasonStyle) {
	if text == "" || rcmd == nil {
		return
	}
	var (
		style, bgstyle int8
	)
	if style = rcmd.CornerMark; style == 0 {
		if rcmd.RcmdReason != nil {
			if rcmd.RcmdReason.Content == "" {
				style = 0
			} else {
				style = rcmd.RcmdReason.CornerMark
			}
		}
	}
	switch style {
	case 0, 2:
		bgstyle = model.BgColorOrange
	case 1:
		bgstyle = model.BgColorTransparentOrange
	case 3:
		bgstyle = model.BgTransparentTextOrange
	case 4:
		bgstyle = model.BgColorRed
	default:
		bgstyle = model.BgColorOrange
	}
	res = reasonStyleFrom(bgstyle, text)
	return
}

func bottomReasonStyleFrom(rcmd *ai.Item, text string, gt model.Gt) (res *ReasonStyle) {
	if text == "" || rcmd == nil {
		return
	}
	var (
		style, bgstyle int8
	)
	if style = rcmd.CornerMark; style == 0 {
		if rcmd.RcmdReason != nil {
			if rcmd.RcmdReason.Content == "" {
				style = 0
			} else {
				style = rcmd.RcmdReason.CornerMark
			}
		}
	}
	switch style {
	case 1:
		bgstyle = model.BgColorTransparentOrange
	case 3:
		bgstyle = model.BgTransparentTextOrange
	default:
		bgstyle = model.BgColorOrange
	}
	res = reasonStyleFrom(bgstyle, text)
	return
}

func reasonStyleFrom(style int8, text string) (res *ReasonStyle) {
	res = &ReasonStyle{
		Text: text,
	}
	switch style {
	case model.BgColorOrange: //defalut
		res.TextColor = "#FFFFFFFF"
		res.BgColor = "#FFFB9E60"
		res.BorderColor = "#FFFB9E60"
		res.BgStyle = model.BgStyleFill
	case model.BgColorTransparentOrange:
		res.TextColor = "#FFFB9E60"
		res.BorderColor = "#FFFB9E60"
		res.BgStyle = model.BgStyleStroke
	case model.BgColorBlue:
		res.TextColor = "#FF23ADE5"
		res.BgColor = "#3323ADE5"
		res.BorderColor = "#3323ADE5"
		res.BgStyle = model.BgStyleFill
	case model.BgColorRed:
		res.TextColor = "#FFFFFFFF"
		res.BgColor = "#FFFB7299"
		res.BorderColor = "#FFFB7299"
		res.BgStyle = model.BgStyleFill
	case model.BgTransparentTextOrange:
		res.TextColor = "#FFFB9E60"
		res.BgStyle = model.BgStyleNoFillAndNoStroke
	case model.BgColorPurple:
		res.TextColor = "#FFFFFFFF"
		res.BgColor = "#FF7D75F2"
		res.BorderColor = "#FF7D75F2"
		res.BgStyle = model.BgStyleFill
	}
	return
}

func unionAuthor(a *archive.ArchiveWithPlayer) (name string) {
	if a.Rights.IsCooperation == 1 {
		name = a.Author.Name + " 等联合创作"
		return
	}
	name = a.Author.Name
	return
}
