package card

import (
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/show"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	season "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

func doubleHandle(cardGoto model.CardGt, cardType model.CardType, rcmd *ai.Item, tagm map[int64]*tag.Tag, isAttenm map[int64]int8, statm map[int64]*relation.Stat, cardm map[int64]*account.Card) (hander Handler) {
	base := &Base{CardType: cardType, CardGoto: cardGoto, Rcmd: rcmd, Tagm: tagm, IsAttenm: isAttenm, Statm: statm, Cardm: cardm, Columnm: model.ColumnSvrDouble}
	switch cardType {
	case model.ThreePicV2:
		base.CardLen = 1
		hander = &ThreePicV2{Base: base}
	case model.SmallCoverV2:
		base.CardLen = 1
		hander = &SmallCoverV2{Base: base}
	case model.OptionsV2:
		hander = &Option{Base: base}
	case model.OnePicV2:
		base.CardLen = 1
		hander = &OnePicV2{Base: base}
	case model.Select:
		hander = &Select{Base: base}
	default:
		switch cardGoto {
		case model.CardGotoAv, model.CardGotoLive, model.CardGotoArticleS, model.CardGotoSpecialS, model.CardGotoShoppingS, model.CardGotoAudio, model.CardGotoGameDownloadS, model.CardGotoBangumi, model.CardGotoMoe, model.CardGotoPGC:
			base.CardType = model.SmallCoverV2
			base.CardLen = 1
			hander = &SmallCoverV2{Base: base}
		case model.CardGotoAdAv:
			base.CardType = model.CmV2
			base.CardLen = 1
			hander = &SmallCoverV2{Base: base}
		case model.CardGotoChannelRcmd, model.CardGotoUpRcmdAv:
			base.CardType = model.SmallCoverV3
			base.CardLen = 1
			hander = &SmallCoverV3{Base: base}
		case model.CardGotoSpecial:
			base.CardType = model.MiddleCoverV2
			hander = &MiddleCover{Base: base}
		case model.CardGotoPlayer, model.CardGotoPlayerLive:
			base.CardType = model.LargeCoverV2
			hander = &LargeCoverV2{Base: base}
		case model.CardGotoSubscribe, model.CardGotoSearchSubscribe:
			base.CardType = model.ThreeItemHV2
			hander = &ThreeItemH{Base: base}
		case model.CardGotoLiveUpRcmd:
			base.CardType = model.TwoItemV2
			return &TwoItemV2{Base: base}
		case model.CardGotoConverge, model.CardGotoRank:
			base.CardType = model.ThreeItemV2
			hander = &ThreeItemV2{Base: base}
		case model.CardGotoBangumiRcmd:
			base.CardType = model.SmallCoverV4
			hander = &SmallCoverV4{Base: base}
		case model.CardGotoLogin:
			base.CardType = model.CoverOnlyV2
			base.CardLen = 1
			return &CoverOnly{Base: base}
		case model.CardGotoBanner:
			base.CardType = model.BannerV2
			return &Banner{Base: base}
		case model.CardGotoAdWebS:
			base.CardType = model.CmV2
			base.CardLen = 1
			hander = &SmallCoverV2{Base: base}
		case model.CardGotoAdWeb:
			base.CardType = model.CmV2
			hander = &MiddleCover{Base: base}
		case model.CardGotoNews:
			base.CardType = model.News
			hander = &Text{Base: base}
		case model.CardGotoEntrance:
			base.CardType = model.MultiItemH
			hander = &MultiItem{Base: base}
		case model.CardGotoTagRcmd, model.CardGotoContentRcmd:
			base.CardType = model.MultiItem
			hander = &MultiItem{Base: base}
		}
	}
	return
}

type SmallCoverV2 struct {
	*Base
	CoverBlur                 model.BlurStatus `json:"cover_blur,omitempty"`
	CoverLeftText1            string           `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1            model.Icon       `json:"cover_left_icon_1,omitempty"`
	CoverLeftText2            string           `json:"cover_left_text_2,omitempty"`
	CoverLeftIcon2            model.Icon       `json:"cover_left_icon_2,omitempty"`
	CoverRightText            string           `json:"cover_right_text,omitempty"`
	CoverRightIcon            model.Icon       `json:"cover_right_icon,omitempty"`
	CoverRightBackgroundColor string           `json:"cover_right_background_color,omitempty"`
	Subtitle                  string           `json:"subtitle,omitempty"`
	Badge                     string           `json:"badge,omitempty"`
	RcmdReason                string           `json:"rcmd_reason,omitempty"`
	DescButton                *Button          `json:"desc_button,omitempty"`
	Desc                      string           `json:"desc,omitempty"`
	Avatar                    *Avatar          `json:"avatar,omitempty"`
	OfficialIcon              model.Icon       `json:"official_icon,omitempty"`
	CanPlay                   int32            `json:"can_play,omitempty"`
	RcmdReasonStyle           *ReasonStyle     `json:"rcmd_reason_style,omitempty"`
}

func (c *SmallCoverV2) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		upID   int64
		button interface{}
		avatar *AvatarStatus
	)
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		am := main.(map[int64]*archive.ArchiveWithPlayer)
		a, ok := am[op.ID]
		if !ok {
			return
		}
		switch op.CardGoto {
		case model.CardGotoAdAv:
			if !model.AdAvIsNormal(a) {
				return
			}
			c.AdInfo = op.AdInfo
		default:
			if !model.AvIsNormal(a) {
				return
			}
		}
		c.Base.from(op.Param, a.Pic, a.Title, model.GotoAv, strconv.FormatInt(a.Aid, 10), model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
		// c.CoverLeftText1 = model.RecommendString(a.Stat.Like, a.Stat.DisLike)
		c.CoverLeftText1 = model.StatString(a.Stat.View, "")
		c.CoverLeftIcon1 = model.IconPlay
		c.CoverLeftText2 = model.StatString(a.Stat.Danmaku, "")
		c.CoverLeftIcon2 = model.IconDanmaku
		if op.SwitchLike == model.SwitchFeedIndexLike {
			c.CoverLeftText1 = model.StatString(a.Stat.Like, "")
			c.CoverLeftIcon1 = model.IconLike
			c.CoverLeftText2 = model.StatString(a.Stat.View, "")
			c.CoverLeftIcon2 = model.IconPlay
		}
		c.CoverRightText = model.DurationString(a.Duration)
		if c.Rcmd != nil {
			c.RcmdReason, c.Desc = rcmdReason(c.Rcmd.RcmdReason, a.Author.Name, c.IsAttenm[a.Author.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		}
		if c.RcmdReason == "" {
			if t, ok := c.Tagm[op.Tid]; ok {
				tag := &tag.Tag{}
				*tag = *t
				tag.Name = a.TypeName + " · " + tag.Name
				button = tag
			} else {
				button = &ButtonStatus{Text: a.TypeName}
			}
		}
		c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
		c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
		c.CanPlay = a.Rights.Autoplay
		upID = a.Author.Mid
		switch op.CardGoto {
		case model.CardGotoAdAv:
			c.AdInfo = op.AdInfo
		}
	case map[int64]*bangumi.Season:
		sm := main.(map[int64]*bangumi.Season)
		s, ok := sm[op.ID]
		if !ok {
			return
		}
		c.Base.from(s.EpisodeID, s.Cover, s.Title, model.GotoBangumi, s.EpisodeID, nil)
		c.CoverLeftText1 = model.StatString(s.PlayCount, "")
		c.CoverLeftIcon1 = model.IconPlay
		c.CoverLeftText2 = model.StatString(s.Favorites, "")
		c.CoverLeftIcon2 = model.BangumiIcon(s.SeasonType)
		c.Badge = s.TypeBadge
		c.Subtitle = s.UpdateDesc
	case map[int32]*season.CardInfoProto:
		sm := main.(map[int32]*season.CardInfoProto)
		s, ok := sm[int32(op.ID)]
		if !ok {
			return
		}
		c.Base.from(op.Param, s.Cover, s.Title, model.GotoPGC, op.URI, nil)
		c.CoverLeftText1 = model.StatString(int32(s.Stat.View), "")
		c.CoverLeftIcon1 = model.IconPlay
		if s.Stat != nil {
			c.CoverLeftText2 = model.StatString(int32(s.Stat.Follow), "")
		}
		c.CoverLeftIcon2 = model.BangumiIcon(int8(s.SeasonType))
		c.Badge = s.SeasonTypeName
		if s.NewEp != nil {
			c.Subtitle = s.NewEp.IndexShow
		}
	case map[int32]*episodegrpc.EpisodeCardsProto:
		sm := main.(map[int32]*episodegrpc.EpisodeCardsProto)
		s, ok := sm[int32(op.ID)]
		if !ok {
			return
		}
		title := s.Season.Title + "：" + s.ShowTitle
		c.Base.from(op.Param, s.Cover, title, model.GotoBangumi, op.URI, nil)
		c.Goto = model.GotoPGC
		c.CoverLeftText1 = model.StatString(int32(s.Season.Stat.View), "")
		c.CoverLeftIcon1 = model.IconPlay
		if s.Season.Stat != nil {
			c.CoverLeftText2 = model.StatString(int32(s.Season.Stat.Follow), "")
		}
		c.CoverLeftIcon2 = model.BangumiIcon(int8(s.Season.SeasonType))
		c.Badge = s.Season.SeasonTypeName
		if s.Season != nil {
			c.Subtitle = s.Season.NewEpShow
		}
	case map[int64]*live.Room:
		rm := main.(map[int64]*live.Room)
		r, ok := rm[op.ID]
		if !ok || r.LiveStatus != 1 {
			return
		}
		c.Base.from(op.Param, r.Cover, r.Title, model.GotoLive, strconv.FormatInt(r.RoomID, 10), model.LiveRoomHandler(r))
		c.CoverLeftText1 = model.StatString(r.Online, "")
		c.CoverLeftIcon1 = model.IconOnline
		c.CoverRightText = r.Uname
		c.Badge = "直播"
		c.Base.PlayerArgs = playerArgsFrom(r)
		c.Args.fromLiveRoom(r)
		if c.Rcmd != nil && (c.Rcmd.RcmdReason != nil || c.IsAttenm[r.UID] == 1) {
			c.RcmdReason, c.Desc = rcmdReason(c.Rcmd.RcmdReason, "", c.IsAttenm[r.UID], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else {
			button = r
		}
		upID = r.UID
		c.CanPlay = 1
	case map[int64]*show.Shopping:
		sm := main.(map[int64]*show.Shopping)
		s, ok := sm[op.ID]
		if !ok {
			return
		}
		c.Base.from(op.Param, model.ShoppingCover(s.PerformanceImage), s.Name, model.GotoWeb, s.URL, nil)
		if s.Type == 1 {
			c.CoverLeftText1 = model.ShoppingDuration(s.STime, s.ETime)
			c.CoverRightText = s.CityName
			c.CoverRightIcon = model.IconLocation
			if len(s.Tags) != 0 {
				c.Desc = s.Tags[0].TagName
			}
		} else if s.Type == 2 {
			c.CoverLeftText1 = s.Want
			c.Desc = s.Subname
		}
		c.Badge = "会员购"
		c.Args.fromShopping(s)
	case map[int64]*audio.Audio:
		am := main.(map[int64]*audio.Audio)
		a, ok := am[op.ID]
		if !ok {
			return
		}
		c.Base.from(op.Param, a.CoverURL, a.Title, model.GotoAudio, strconv.FormatInt(a.MenuID, 10), nil)
		c.CoverBlur = model.BlurYes
		c.CoverLeftText1 = model.StatString(a.PlayNum, "")
		c.CoverLeftIcon1 = model.IconHeadphone
		c.CoverRightText = model.AudioTotalStirng(a.RecordNum)
		c.Badge = model.AudioBadgeString(a.Type)
		button = a.Ctgs
		c.Args.fromAudio(a)
	case map[int64]*article.Meta:
		mm := main.(map[int64]*article.Meta)
		m, ok := mm[op.ID]
		if !ok {
			return
		}
		if len(m.ImageURLs) == 0 {
			return
		}
		c.Base.from(op.Param, m.ImageURLs[0], m.Title, model.GotoArticle, strconv.FormatInt(m.ID, 10), nil)
		if m.Stats != nil {
			c.CoverLeftText1 = model.StatString(int32(m.Stats.View), "")
			c.CoverLeftIcon1 = model.IconRead
			c.CoverLeftText2 = model.StatString(int32(m.Stats.Reply), "")
			c.CoverLeftIcon2 = model.IconComment
		}
		button = m.Categories
		c.Badge = "文章"
		c.Args.fromArticle(m)
	case map[int64]*bplus.Picture:
		pm := main.(map[int64]*bplus.Picture)
		p, ok := pm[op.ID]
		if !ok || len(p.Imgs) == 0 || p.ViewCount == 0 {
			return
		}
		c.Base.from(op.Param, p.Imgs[0], p.DynamicText, model.GotoPicture, strconv.FormatInt(p.DynamicID, 10), nil)
		c.CoverLeftText1 = model.StatString(int32(p.ViewCount), "")
		c.CoverLeftIcon1 = model.IconRead
		if p.ImgCount > 1 {
			c.CoverRightText = model.PictureCountString(p.ImgCount)
			c.CoverRightBackgroundColor = "#66666666"
		}
		if c.Rcmd != nil && c.Rcmd.RcmdReason != nil {
			c.RcmdReason, _ = rcmdReason(c.Rcmd.RcmdReason, p.NickName, c.IsAttenm[p.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else {
			button = p
			c.Badge = "动态"
		}
	case *cm.AdInfo:
		ad := main.(*cm.AdInfo)
		c.AdInfo = ad
	case *bangumi.Moe:
		m := main.(*bangumi.Moe)
		if m == nil {
			return
		}
		c.Base.from(strconv.FormatInt(m.ID, 10), m.Square, m.Title, model.GotoWeb, m.Link, nil)
		c.Desc = m.Desc
		c.Badge = m.Badge
	case nil:
		if op == nil {
			return
		}
		c.Base.from(op.Param, op.Coverm[c.Columnm], op.Title, op.Goto, op.URI, nil)
		switch op.CardGoto {
		case model.CardGotoDownload:
			c.CoverLeftText1 = model.DownloadString(op.Download)
			avatar = &AvatarStatus{Cover: op.Avatar, Goto: op.Goto, Param: op.URI, Type: model.AvatarSquare}
			c.Desc = op.Desc
		case model.CardGotoSpecial:
			c.Desc = op.Desc
			c.Badge = op.Badge
		default:
			log.Warn("SmallCoverV2 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("SmallCoverV2 From: unexpected type %T", main)
		return
	}
	c.OfficialIcon = model.OfficialIcon(c.Cardm[upID])
	c.Avatar = avatarFrom(avatar)
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *SmallCoverV2) Get() *Base {
	return c.Base
}

type SmallCoverV3 struct {
	*Base
	Avatar           *Avatar      `json:"avatar,omitempty"`
	CoverLeftText    string       `json:"cover_left_text,omitempty"`
	CoverRightButton *Button      `json:"cover_right_button,omitempty"`
	RcmdReason       string       `json:"rcmd_reason,omitempty"`
	Desc             string       `json:"desc,omitempty"`
	DescButton       *Button      `json:"desc_button,omitempty"`
	OfficialIcon     model.Icon   `json:"official_icon,omitempty"`
	CanPlay          int32        `json:"can_play,omitempty"`
	RcmdReasonStyle  *ReasonStyle `json:"rcmd_reason_style,omitempty"`
}

func (c *SmallCoverV3) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button     interface{}
		descButton interface{}
	)
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		am := main.(map[int64]*archive.ArchiveWithPlayer)
		a, ok := am[op.ID]
		if !ok || !model.AvIsNormal(a) {
			return
		}
		c.Base.from(op.Param, a.Pic, a.Title, model.GotoAv, op.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
		switch op.CardGoto {
		case model.CardGotoUpRcmdAv:
			c.Avatar = avatarFrom(&AvatarStatus{Cover: a.Author.Face, Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), Type: model.AvatarRound})
			c.CoverLeftText = a.Author.Name
			if c.Rcmd != nil && c.Rcmd.RcmdReason != nil {
				c.RcmdReason, _ = rcmdReason(c.Rcmd.RcmdReason, "", c.IsAttenm[a.Author.Mid], c.Cardm)
				c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
			} else {
				descButton = c.Tagm[op.Tid]
			}
			button = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), IsAtten: c.IsAttenm[a.Author.Mid]}
			c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
			c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
		case model.CardGotoChannelRcmd:
			t, ok := c.Tagm[op.Tid]
			if !ok {
				return
			}
			c.Avatar = avatarFrom(&AvatarStatus{Cover: t.Cover, Goto: model.GotoTag, Param: strconv.FormatInt(t.ID, 10), Type: model.AvatarSquare})
			c.CoverLeftText = t.Name
			c.Desc = model.SubscribeString(int32(t.Count.Atten))
			button = &ButtonStatus{Goto: model.GotoTag, Param: strconv.FormatInt(t.ID, 10), IsAtten: t.IsAtten}
			c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
			c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
		default:
			log.Warn("SmallCoverV3 From: unexpected card_goto %s", op.CardGoto)
			return
		}
		c.CanPlay = a.Rights.Autoplay
	default:
		log.Warn("SmallCoverV3 From: unexpected type %T", main)
		return
	}
	c.CoverRightButton = buttonFrom(button, op.Plat)
	c.DescButton = buttonFrom(descButton, op.Plat)
	c.Right = true
}

func (c *SmallCoverV3) Get() *Base {
	return c.Base
}

type MiddleCoverV2 struct {
	*Base
	Ratio int    `json:"ratio,omitempty"`
	Desc  string `json:"desc,omitempty"`
	Badge string `json:"badge,omitempty"`
}

func (c *MiddleCoverV2) Get() *Base {
	return c.Base
}

type LargeCoverV2 struct {
	*Base
	Avatar           *Avatar      `json:"avatar,omitempty"`
	Badge            string       `json:"badge,omitempty"`
	CoverRightButton *Button      `json:"cover_right_button,omitempty"`
	CoverLeftText1   string       `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1   model.Icon   `json:"cover_left_icon_1,omitempty"`
	CoverLeftText2   string       `json:"cover_left_text_2,omitempty"`
	CoverLeftIcon2   model.Icon   `json:"cover_left_icon_2,omitempty"`
	RcmdReason       string       `json:"rcmd_reason,omitempty"`
	DescButton       *Button      `json:"desc_button,omitempty"`
	OfficialIcon     model.Icon   `json:"official_icon,omitempty"`
	CanPlay          int32        `json:"can_play,omitempty"`
	RcmdReasonStyle  *ReasonStyle `json:"rcmd_reason_style,omitempty"`
}

func (c *LargeCoverV2) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button      interface{}
		coverButton interface{}
		upID        int64
	)
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		am := main.(map[int64]*archive.ArchiveWithPlayer)
		a, ok := am[op.ID]
		if !ok || !model.AvIsNormal(a) {
			return
		}
		c.Base.from(op.Param, a.Pic, a.Title, model.GotoAv, op.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
		c.Avatar = avatarFrom(&AvatarStatus{Cover: a.Author.Face, Text: a.Author.Name, Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), Type: model.AvatarRound})
		c.CoverLeftText1 = model.StatString(a.Stat.View, "")
		c.CoverLeftIcon1 = model.IconPlay
		c.CoverLeftText2 = model.StatString(a.Stat.Danmaku, "")
		c.CoverLeftIcon2 = model.IconDanmaku
		if op.SwitchLike == model.SwitchFeedIndexLike {
			c.CoverLeftText1 = model.StatString(a.Stat.Like, "")
			c.CoverLeftIcon1 = model.IconLike
			c.CoverLeftText2 = model.StatString(a.Stat.View, "")
			c.CoverLeftIcon2 = model.IconPlay
		}
		coverButton = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), IsAtten: c.IsAttenm[a.Author.Mid]}
		if c.Rcmd != nil && c.Rcmd.RcmdReason != nil {
			c.RcmdReason, _ = rcmdReason(c.Rcmd.RcmdReason, "", c.IsAttenm[a.Author.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else if t, ok := c.Tagm[op.Tid]; ok {
			tag := &tag.Tag{}
			*tag = *t
			tag.Name = a.TypeName + " · " + tag.Name
			button = tag
		} else {
			button = &ButtonStatus{Text: a.TypeName}
		}
		c.CanPlay = a.Rights.Autoplay
		c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
		if op.CardGoto == model.CardGotoPlayer && c.Base.PlayerArgs == nil {
			log.Warn("player card aid(%d) can't auto player", a.Aid)
			return
		}
		c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
		upID = a.Author.Mid
	case map[int64]*live.Room:
		rm := main.(map[int64]*live.Room)
		r, ok := rm[op.ID]
		if !ok || r.LiveStatus != 1 {
			return
		}
		c.Base.from(op.Param, r.Cover, r.Title, model.GotoLive, op.URI, model.LiveRoomHandler(r))
		c.Avatar = avatarFrom(&AvatarStatus{Cover: r.Face, Text: r.Uname, Goto: model.GotoMid, Param: strconv.FormatInt(r.UID, 10), Type: model.AvatarRound})
		c.CoverLeftText1 = model.StatString(r.Online, "")
		c.CoverLeftIcon1 = model.IconOnline
		coverButton = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(r.UID, 10), IsAtten: c.IsAttenm[r.UID]}
		if c.Rcmd != nil && (c.Rcmd.RcmdReason != nil || c.IsAttenm[r.UID] == 1) {
			c.RcmdReason, _ = rcmdReason(c.Rcmd.RcmdReason, r.Uname, c.IsAttenm[r.UID], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else {
			button = r
		}
		c.Badge = "直播"
		c.CanPlay = 1
		c.Base.PlayerArgs = playerArgsFrom(r)
		c.Args.fromLiveRoom(r)
		upID = r.UID
	default:
		log.Warn("MiddleCoverV2 From: unexpected type %T", main)
		return
	}
	c.DescButton = buttonFrom(button, op.Plat)
	c.CoverRightButton = buttonFrom(coverButton, op.Plat)
	c.OfficialIcon = model.OfficialIcon(c.Cardm[upID])
	c.Right = true
}

func (c *LargeCoverV2) Get() *Base {
	return c.Base
}

type ThreeItemV2 struct {
	*Base
	TitleIcon model.Icon         `json:"title_icon,omitempty"`
	MoreURI   string             `json:"more_uri,omitempty"`
	MoreText  string             `json:"more_text,omitempty"`
	Items     []*ThreeItemV2Item `json:"items,omitempty"`
}

type ThreeItemV2Item struct {
	Base
	CoverLeftIcon model.Icon `json:"cover_left_icon,omitempty"`
	DescText1     string     `json:"desc_text_1,omitempty"`
	DescIcon1     model.Icon `json:"desc_icon_1,omitempty"`
	DescText2     string     `json:"desc_text_2,omitempty"`
	DescIcon2     model.Icon `json:"desc_icon_2,omitempty"`
	Badge         string     `json:"badge,omitempty"`
}

func (c *ThreeItemV2) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case map[model.Gt]interface{}:
		intfcm := main.(map[model.Gt]interface{})
		if op == nil {
			return
		}
		switch op.CardGoto {
		case model.CardGotoRank:
			const (
				_title = "全站排行榜"
				_limit = 3
			)
			c.Base.from("0", "", _title, "", "", nil)
			c.TitleIcon = model.IconRank
			c.MoreURI = model.FillURI(op.Goto, op.URI, nil)
			c.MoreText = "查看更多"
			c.Items = make([]*ThreeItemV2Item, 0, _limit)
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				var item *ThreeItemV2Item
				switch intfc.(type) {
				case map[int64]*archive.ArchiveWithPlayer:
					am := intfc.(map[int64]*archive.ArchiveWithPlayer)
					a, ok := am[v.ID]
					if !ok || !model.AvIsNormal(a) {
						continue
					}
					item = &ThreeItemV2Item{
						DescText1: model.ScoreString(v.Score),
					}
					item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, v.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
					item.Args.fromArchive(a.Archive3, nil)
				default:
					log.Warn("ThreeItemV2 From: unexpected type %T", intfc)
					continue
				}
				c.Items = append(c.Items, item)
				if len(c.Items) == _limit {
					break
				}
			}
			if len(c.Items) < _limit {
				return
			}
			c.Items[0].CoverLeftIcon = model.IconGoldMedal
			c.Items[1].CoverLeftIcon = model.IconSilverMedal
			c.Items[2].CoverLeftIcon = model.IconBronzeMedal
		case model.CardGotoConverge:
			limit := 3
			if op.Coverm[c.Columnm] != "" {
				limit = 2
			}
			c.Base.from(op.Param, op.Coverm[c.Columnm], op.Title, op.Goto, op.URI, nil)
			c.MoreURI = model.FillURI(model.GotoConverge, op.Param, nil)
			c.MoreText = "查看更多"
			c.Items = make([]*ThreeItemV2Item, 0, len(op.Items))
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				var item *ThreeItemV2Item
				switch intfc.(type) {
				case map[int64]*archive.ArchiveWithPlayer:
					am := intfc.(map[int64]*archive.ArchiveWithPlayer)
					a, ok := am[v.ID]
					if !ok || !model.AvIsNormal(a) {
						continue
					}
					item = &ThreeItemV2Item{
						DescText1: model.StatString(a.Stat.View, ""),
						DescIcon1: model.IconPlay,
						DescText2: model.StatString(a.Stat.Danmaku, ""),
						DescIcon2: model.IconDanmaku,
					}
					if op.SwitchLike == model.SwitchFeedIndexLike {
						item.DescText1 = model.StatString(a.Stat.Like, "")
						item.DescIcon1 = model.IconLike
						item.DescText2 = model.StatString(a.Stat.View, "")
						item.DescIcon2 = model.IconPlay
					}
					item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, v.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
					item.Args.fromArchive(a.Archive3, nil)
				case map[int64]*live.Room:
					rm := intfc.(map[int64]*live.Room)
					r, ok := rm[v.ID]
					if !ok || r.LiveStatus != 1 {
						continue
					}
					item = &ThreeItemV2Item{
						DescText1: model.StatString(r.Online, ""),
						DescIcon1: model.IconOnline,
						Badge:     "直播",
					}
					item.Base.from(v.Param, r.Cover, r.Title, model.GotoLive, v.URI, model.LiveRoomHandler(r))
					item.Args.fromLiveRoom(r)
				case map[int64]*article.Meta:
					mm := intfc.(map[int64]*article.Meta)
					m, ok := mm[v.ID]
					if !ok {
						continue
					}
					if len(m.ImageURLs) == 0 {
						continue
					}
					item = &ThreeItemV2Item{
						Badge: "文章",
					}
					item.Base.from(v.Param, m.ImageURLs[0], m.Title, model.GotoArticle, v.URI, nil)
					if m.Stats != nil {
						item.DescText1 = model.StatString(int32(m.Stats.View), "")
						item.DescIcon1 = model.IconRead
						item.DescText2 = model.StatString(int32(m.Stats.Reply), "")
						item.DescIcon2 = model.IconComment
					}
					item.Args.fromArticle(m)
				default:
					log.Warn("ThreeItemV2 From: unexpected type %T", intfc)
					continue
				}
				c.Items = append(c.Items, item)
				if len(c.Items) == limit {
					break
				}
			}
			if len(c.Items) < limit {
				return
			}
		default:
			log.Warn("ThreeItemV2 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("ThreeItemV2 From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *ThreeItemV2) Get() *Base {
	return c.Base
}

type SmallCoverV4 struct {
	*Base
	CoverBadge     string     `json:"cover_badge,omitempty"`
	Desc           string     `json:"desc,omitempty"`
	TitleRightText string     `json:"title_right_text,omitempty"`
	TitleRightPic  model.Icon `json:"title_right_pic,omitempty"`
}

func (c *SmallCoverV4) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case *bangumi.Update:
		title := "你的追番更新啦"
		const (
			_updates = 99
		)
		u := main.(*bangumi.Update)
		if u == nil || u.Updates == 0 {
			return
		}
		emojim := map[string]struct{}{
			"(´∀｀*)ｳﾌﾌ": struct{}{},
			"ヾ( ・∀・)ﾉ":  struct{}{},
			"(｀･ω･´)ゞ":  struct{}{},
			"(・∀・)ｲｲ!!": struct{}{},
		}
		for emoji := range emojim {
			title = title + emoji
			break
		}
		c.Base.from("", u.SquareCover, title, "", "", nil)
		updates := u.Updates
		if updates > _updates {
			updates = _updates
			c.TitleRightPic = model.IconBomb
		} else {
			c.TitleRightPic = model.IconTV
		}
		c.Desc = u.Title
		c.TitleRightText = strconv.Itoa(updates)
	default:
		log.Warn("SmallCoverV4 From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *SmallCoverV4) Get() *Base {
	return c.Base
}

type TwoItemV2 struct {
	*Base
	Items []*TwoItemV2Item `json:"items,omitempty"`
}

type TwoItemV2Item struct {
	Base
	Badge          string     `json:"badge,omitempty"`
	CoverLeftText1 string     `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1 model.Icon `json:"cover_left_icon_1,omitempty"`
	DescButton     *Button    `json:"desc_button,omitempty"`
}

func (c *TwoItemV2) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case map[int64][]*live.Card:
		const _limit = 2
		csm := main.(map[int64][]*live.Card)
		cs, ok := csm[op.ID]
		if !ok {
			return
		}
		c.Base.from(op.Param, "", "", "", "", nil)
		c.Items = make([]*TwoItemV2Item, 0, _limit)
		for _, card := range cs {
			if card == nil || card.LiveStatus != 1 {
				continue
			}
			item := &TwoItemV2Item{
				Badge:          "直播",
				CoverLeftText1: model.StatString(card.Online, ""),
				CoverLeftIcon1: model.IconOnline,
			}
			item.DescButton = buttonFrom(card, op.Plat)
			item.Base.from(strconv.FormatInt(card.RoomID, 10), card.ShowCover, card.Title, model.GotoLive, strconv.FormatInt(card.RoomID, 10), model.LiveUpHandler(card))
			item.Args.fromLiveUp(card)
			c.Items = append(c.Items, item)
			if len(c.Items) == _limit {
				break
			}
		}
	}
	c.Right = true
}

func (c *TwoItemV2) Get() *Base {
	return c.Base
}

type MultiItem struct {
	*Base
	MoreURI  string    `json:"more_uri,omitempty"`
	MoreText string    `json:"more_text,omitempty"`
	Items    []Handler `json:"items,omitempty"`
}

func (c *MultiItem) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case map[model.Gt]interface{}:
		intfcm := main.(map[model.Gt]interface{})
		switch op.CardGoto {
		case model.CardGotoTagRcmd, model.CardGotoContentRcmd:
			items := make([]Handler, 0, len(op.Items))
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				var hander Handler
				switch intfc.(type) {
				case map[int64]*archive.ArchiveWithPlayer:
					am := intfc.(map[int64]*archive.ArchiveWithPlayer)
					a, ok := am[v.ID]
					if !ok || !model.AvIsNormal(a) {
						continue
					}
					item := &SmallCoverV2{
						CoverLeftText1: model.StatString(a.Stat.View, ""),
						CoverLeftIcon1: model.IconPlay,
						CoverLeftText2: model.StatString(a.Stat.Danmaku, ""),
						CoverLeftIcon2: model.IconDanmaku,
						CoverRightText: model.DurationString(a.Duration),
						Base:           &Base{CardType: model.SmallCoverV2},
					}
					if op.SwitchLike == model.SwitchFeedIndexLike {
						item.CoverLeftText1 = model.StatString(a.Stat.Like, "")
						item.CoverLeftIcon1 = model.IconLike
						item.CoverLeftText2 = model.StatString(a.Stat.View, "")
						item.CoverLeftIcon2 = model.IconPlay
					}
					item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, strconv.FormatInt(a.Aid, 10), model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
					item.Args.fromArchive(a.Archive3, nil)
					if op.Switch == model.SwitchFeedIndexTabThreePoint {
						item.TabThreePointWatchLater()
					}
					item.DescButton = buttonFrom(&ButtonStatus{Text: a.TypeName}, op.Plat)
					hander = item
				case map[int64]*live.Room:
					rm := intfc.(map[int64]*live.Room)
					r, ok := rm[v.ID]
					if !ok || r.LiveStatus != 1 {
						continue
					}
					item := &SmallCoverV2{
						CoverLeftText1: model.StatString(r.Online, ""),
						CoverLeftIcon1: model.IconOnline,
						Badge:          "直播",
						Base:           &Base{CardType: model.SmallCoverV2},
					}
					item.Base.from(v.Param, r.Cover, r.Title, model.GotoLive, strconv.FormatInt(r.RoomID, 10), model.LiveRoomHandler(r))
					item.Args.fromLiveRoom(r)
					item.DescButton = buttonFrom(r, op.Plat)
					hander = item
				case map[int64]*article.Meta:
					mm := intfc.(map[int64]*article.Meta)
					m, ok := mm[v.ID]
					if !ok {
						continue
					}
					if len(m.ImageURLs) == 0 {
						continue
					}
					item := &SmallCoverV2{
						Badge: "文章",
						Base:  &Base{CardType: model.SmallCoverV2},
					}
					item.Base.from(v.Param, m.ImageURLs[0], m.Title, model.GotoArticle, strconv.FormatInt(m.ID, 10), nil)
					if m.Stats != nil {
						item.CoverLeftText1 = model.StatString(int32(m.Stats.View), "")
						item.CoverLeftIcon1 = model.IconRead
						item.CoverLeftText2 = model.StatString(int32(m.Stats.Reply), "")
						item.CoverLeftIcon2 = model.IconComment
					}
					item.Args.fromArticle(m)
					item.DescButton = buttonFrom(m.Categories, op.Plat)
					hander = item
				case map[int64]*operate.Card:
					dm := intfc.(map[int64]*operate.Card)
					d, ok := dm[v.ID]
					if !ok {
						continue
					}
					item := &SmallCoverV2{
						CoverLeftText1: model.DownloadString(d.Download),
						Base:           &Base{CardType: model.SmallCoverV2},
					}
					item.Base.from(v.Param, d.Coverm[c.Columnm], d.Title, d.Goto, d.URI, nil)
					hander = item
				case map[int64]*bangumi.Season:
					sm := intfc.(map[int64]*bangumi.Season)
					s, ok := sm[v.ID]
					if !ok {
						continue
					}
					item := &SmallCoverV2{
						CoverLeftText1: model.StatString(s.PlayCount, ""),
						CoverLeftIcon1: model.IconPlay,
						CoverLeftText2: model.StatString(s.Favorites, ""),
						CoverLeftIcon2: model.BangumiIcon(s.SeasonType),
						Badge:          s.TypeBadge,
						Desc:           s.UpdateDesc,
						Base:           &Base{CardType: model.SmallCoverV2},
					}
					item.Base.from(s.EpisodeID, s.Cover, s.Title, model.GotoBangumi, s.EpisodeID, nil)
					hander = item
				case map[int64]*bplus.Picture:
					pm := intfc.(map[int64]*bplus.Picture)
					p, ok := pm[v.ID]
					if !ok {
						continue
					}
					if len(p.Imgs) < 3 {
						hander = &OnePicV2{Base: &Base{CardType: model.OnePicV2}}
					} else {
						hander = &ThreePicV2{Base: &Base{CardType: model.ThreePicV2}}
					}
					hander.From(pm, v)
					if !hander.Get().Right {
						continue
					}
				default:
					log.Warn("MultiItem From: unexpected type %T", intfc)
					continue
				}
				if hander != nil {
					items = append(items, hander)
				}
			}
			if len(items) < 2 {
				return
			}
			if len(items)%2 != 0 {
				c.Items = items[:len(items)-1]
			} else {
				c.Items = items
			}
			var title string
			switch op.Goto {
			case model.GotoTag:
				if t, ok := c.Tagm[op.ID]; ok {
					title = t.Name
				}
			default:
				title = op.Title
			}
			c.Base.from(op.Param, "", title, "", "", nil)
			c.MoreURI = model.FillURI(op.Goto, op.URI, nil)
			c.MoreText = op.Subtitle
		default:
			log.Warn("MultiItem From: unexpected card_goto %s", op.CardGoto)
			return
		}
	case nil:
		switch op.CardGoto {
		case model.CardGotoEntrance:
			c.Items = make([]Handler, 0, len(op.Items))
			for _, v := range op.Items {
				item := &SmallCoverV2{Base: &Base{CardType: model.SmallCoverV2}}
				item.Base.from(v.Param, v.Cover, v.Title, v.Goto, v.URI, nil)
				c.Items = append(c.Items, item)
			}
		}
	default:
		log.Warn("MultiItem From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *MultiItem) Get() *Base {
	return c.Base
}

type ThreePicV2 struct {
	*Base
	LeftCover                 string       `json:"left_cover,omitempty"`
	RightCover1               string       `json:"right_cover_1,omitempty"`
	RightCover2               string       `json:"right_cover_2,omitempty"`
	CoverLeftText1            string       `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1            model.Icon   `json:"cover_left_icon_1,omitempty"`
	CoverLeftText2            string       `json:"cover_left_text_2,omitempty"`
	CoverLeftIcon2            model.Icon   `json:"cover_left_icon_2,omitempty"`
	CoverRightText            string       `json:"cover_right_text,omitempty"`
	CoverRightIcon            model.Icon   `json:"cover_right_icon,omitempty"`
	CoverRightBackgroundColor string       `json:"cover_right_background_color,omitempty"`
	Badge                     string       `json:"badge,omitempty"`
	RcmdReason                string       `json:"rcmd_reason,omitempty"`
	DescButton                *Button      `json:"desc_button,omitempty"`
	Desc                      string       `json:"desc,omitempty"`
	Avatar                    *Avatar      `json:"avatar,omitempty"`
	RcmdReasonStyle           *ReasonStyle `json:"rcmd_reason_style,omitempty"`
}

func (c *ThreePicV2) From(main interface{}, op *operate.Card) {
	var (
		button interface{}
	)
	if op == nil {
		return
	}
	switch main.(type) {
	case map[int64]*bplus.Picture:
		pm := main.(map[int64]*bplus.Picture)
		p, ok := pm[op.ID]
		if !ok || len(p.Imgs) < 3 || p.ViewCount == 0 {
			return
		}
		c.Base.from(op.Param, "", p.DynamicText, model.GotoPicture, strconv.FormatInt(p.DynamicID, 10), nil)
		c.LeftCover = p.Imgs[0]
		c.RightCover1 = p.Imgs[1]
		c.RightCover2 = p.Imgs[2]
		c.CoverLeftText1 = model.StatString(int32(p.ViewCount), "")
		c.CoverLeftIcon1 = model.IconRead
		if p.ImgCount > 3 {
			c.CoverRightText = model.PictureCountString(p.ImgCount)
			c.CoverRightBackgroundColor = "#66666666"
		}
		if c.Rcmd != nil && c.Rcmd.RcmdReason != nil {
			c.RcmdReason, c.Desc = rcmdReason(c.Rcmd.RcmdReason, p.NickName, c.IsAttenm[p.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else {
			button = p
			c.Badge = "动态"
		}
		c.Avatar = avatarFrom(&AvatarStatus{Cover: p.FaceImg, Text: p.NickName, Goto: model.GotoDynamicMid, Param: strconv.FormatInt(p.Mid, 10), Type: model.AvatarRound})
	default:
		log.Warn("ThreePicV2 From: unexpected type %T", main)
		return
	}
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *ThreePicV2) Get() *Base {
	return c.Base
}

type OnePicV2 struct {
	*Base
	CoverLeftText1            string       `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1            model.Icon   `json:"cover_left_icon_1,omitempty"`
	CoverRightText            string       `json:"cover_right_text,omitempty"`
	CoverRightIcon            model.Icon   `json:"cover_right_icon,omitempty"`
	CoverRightBackgroundColor string       `json:"cover_right_background_color,omitempty"`
	Badge                     string       `json:"badge,omitempty"`
	RcmdReason                string       `json:"rcmd_reason,omitempty"`
	Avatar                    *Avatar      `json:"avatar,omitempty"`
	RcmdReasonStyle           *ReasonStyle `json:"rcmd_reason_style,omitempty"`
}

func (c *OnePicV2) From(main interface{}, op *operate.Card) {
	var (
		button interface{}
	)
	if op == nil {
		return
	}
	switch main.(type) {
	case map[int64]*bplus.Picture:
		pm := main.(map[int64]*bplus.Picture)
		p, ok := pm[op.ID]
		if !ok || len(p.Imgs) == 0 || p.ViewCount == 0 {
			return
		}
		c.Base.from(op.Param, p.Imgs[0], p.DynamicText, model.GotoPicture, strconv.FormatInt(p.DynamicID, 10), nil)
		c.CoverLeftText1 = model.StatString(int32(p.ViewCount), "")
		c.CoverLeftIcon1 = model.IconRead
		if p.ImgCount > 1 {
			c.CoverRightText = model.PictureCountString(p.ImgCount)
			c.CoverRightBackgroundColor = "#66666666"
		}
		if c.Rcmd != nil && c.Rcmd.RcmdReason != nil {
			c.RcmdReason, _ = rcmdReason(c.Rcmd.RcmdReason, p.NickName, c.IsAttenm[p.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.RcmdReason, c.Base.Goto)
		} else {
			button = p
			c.Badge = "动态"
		}
		c.Avatar = avatarFrom(&AvatarStatus{Cover: p.FaceImg, Text: p.NickName, Goto: model.GotoDynamicMid, Param: strconv.FormatInt(p.Mid, 10), Type: model.AvatarRound})
	default:
		log.Warn("OnePicV2 From: unexpected type %T", main)
		return
	}
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *OnePicV2) Get() *Base {
	return c.Base
}
