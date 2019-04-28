package card

import (
	"fmt"
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/show"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	season "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

func singleHandle(cardGoto model.CardGt, cardType model.CardType, rcmd *ai.Item, tagm map[int64]*tag.Tag, isAttenm map[int64]int8, statm map[int64]*relation.Stat, cardm map[int64]*account.Card) (hander Handler) {
	base := &Base{CardType: cardType, CardGoto: cardGoto, Rcmd: rcmd, Tagm: tagm, IsAttenm: isAttenm, Statm: statm, Cardm: cardm, Columnm: model.ColumnSvrSingle}
	switch cardType {
	case model.LargeCoverV1:
		hander = &LargeCoverV1{Base: base}
	case model.OnePicV1:
		hander = &OnePicV1{Base: base}
	case model.ThreePicV1:
		hander = &ThreePicV1{Base: base}
	case model.SmallCoverV5:
		hander = &SmallCoverV5{Base: base}
	case model.OptionsV1:
		hander = &Option{Base: base}
	case model.Select:
		hander = &Select{Base: base}
	default:
		switch cardGoto {
		case model.CardGotoAv, model.CardGotoBangumi, model.CardGotoLive, model.CardGotoPlayer, model.CardGotoPlayerLive, model.CardGotoChannelRcmd, model.CardGotoUpRcmdAv, model.CardGotoPGC:
			base.CardType = model.LargeCoverV1
			hander = &LargeCoverV1{Base: base}
		case model.CardGotoAudio, model.CardGotoBangumiRcmd, model.CardGotoGameDownloadS, model.CardGotoShoppingS, model.CardGotoSpecialS, model.CardGotoMoe:
			base.CardType = model.SmallCoverV1
			hander = &SmallCoverV1{Base: base}
		case model.CardGotoSpecial:
			base.CardType = model.MiddleCoverV1
			hander = &MiddleCover{Base: base}
		case model.CardGotoConverge, model.CardGotoRank:
			base.CardType = model.ThreeItemV1
			hander = &ThreeItemV1{Base: base}
		case model.CardGotoSubscribe, model.CardGotoSearchSubscribe:
			base.CardType = model.ThreeItemHV1
			hander = &ThreeItemH{Base: base}
		case model.CardGotoArticleS:
			base.CardType = model.ThreeItemHV3
			hander = &ThreeItemHV3{Base: base}
		case model.CardGotoLiveUpRcmd:
			base.CardType = model.TwoItemV1
			hander = &TwoItemV1{Base: base}
		case model.CardGotoLogin:
			base.CardType = model.CoverOnlyV1
			hander = &CoverOnly{Base: base}
		case model.CardGotoBanner:
			base.CardType = model.BannerV1
			hander = &Banner{Base: base}
		case model.CardGotoAdAv:
			base.CardType = model.CmV1
			hander = &LargeCoverV1{Base: base}
		case model.CardGotoAdWebS, model.CardGotoAdWeb:
			base.CardType = model.CmV1
			hander = &SmallCoverV1{Base: base}
		case model.CardGotoHotTopic:
			base.CardType = model.HotTopic
			hander = &HotTopic{Base: base}
		case model.CardGotoTopstick:
			base.CardType = model.TopStick
			hander = &Topstick{Base: base}
		case model.CardGotoChannelSquare:
			base.CardType = model.ChannelSquare
			hander = &ChannelSquare{Base: base}
		case model.CardGotoPgcsRcmd:
			base.CardType = model.ThreeItemHV4
			hander = &ThreeItemHV4{Base: base}
		case model.CardGotoUpRcmdS:
			base.CardType = model.UpRcmdCover
			hander = &UpRcmdCover{Base: base}
		case model.CardGotoSearchUpper:
			base.CardType = model.ThreeItemAll
			hander = &ThreeItemAll{Base: base}
		case model.CardGotoUpRcmdNew:
			base.CardType = model.TwoItemHV1
			hander = &TwoItemHV1{Base: base}
		case model.CardGotoUpRcmdNewV2:
			base.CardType = model.ThreeItemAllV2
			hander = &ThreeItemAllV2{Base: base}
		case model.CardGotoDynamicHot:
			base.CardType = model.ThreeItemHV5
			hander = &DynamicHot{Base: base}
		case model.CardGotoEventTopic:
			base.CardType = model.MiddleCoverV3
			hander = &MiddleCoverV3{Base: base}
		}
	}
	return
}

type LargeCoverV1 struct {
	*Base
	Avatar                *Avatar          `json:"avatar,omitempty"`
	CoverLeftText1        string           `json:"cover_left_text_1,omitempty"`
	CoverLeftText2        string           `json:"cover_left_text_2,omitempty"`
	CoverLeftText3        string           `json:"cover_left_text_3,omitempty"`
	CoverBadge            string           `json:"cover_badge,omitempty"`
	TopRcmdReason         string           `json:"top_rcmd_reason,omitempty"`
	BottomRcmdReason      string           `json:"bottom_rcmd_reason,omitempty"`
	Desc                  string           `json:"desc,omitempty"`
	OfficialIcon          model.Icon       `json:"official_icon,omitempty"`
	CanPlay               int32            `json:"can_play,omitempty"`
	CoverBadgeColor       model.CoverColor `json:"cover_badge_color,omitempty"`
	TopRcmdReasonStyle    *ReasonStyle     `json:"top_rcmd_reason_style,omitempty"`
	BottomRcmdReasonStyle *ReasonStyle     `json:"bottom_rcmd_reason_style,omitempty"`
	CoverBadge2           string           `json:"cover_badge_2,omitempty"`
}

func (c *LargeCoverV1) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button interface{}
		avatar *AvatarStatus
		upID   int64
	)
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		am := main.(map[int64]*archive.ArchiveWithPlayer)
		a, ok := am[op.ID]
		if !ok || !model.AvIsNormal(a) {
			return
		}
		c.Base.from(op.Param, a.Pic, a.Title, model.GotoAv, op.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
		c.CoverLeftText1 = model.DurationString(a.Duration)
		c.CoverLeftText2 = model.ArchiveViewString(a.Stat.View)
		c.CoverLeftText3 = model.DanmakuString(a.Stat.Danmaku)
		if op.SwitchLike == model.SwitchFeedIndexLike {
			c.CoverLeftText2 = model.LikeString(a.Stat.Like)
			c.CoverLeftText3 = model.ArchiveViewString(a.Stat.View)
		}
		switch op.CardGoto {
		case model.CardGotoAv, model.CardGotoUpRcmdAv, model.CardGotoPlayer:
			var (
				authorface = a.Author.Face
				authorname = a.Author.Name
			)
			if a.Author.Name != "" {
				if op.Switch != model.SwitchCooperationHide {
					authorname = unionAuthor(a)
				}
			}
			if (authorface == "" || authorname == "") && c.Cardm != nil {
				if au, ok := c.Cardm[a.Author.Mid]; ok {
					authorface = au.Face
					authorname = au.Name
				}
			}
			avatar = &AvatarStatus{Cover: authorface, Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), Type: model.AvatarRound}
			if c.Rcmd != nil && c.Rcmd.RcmdReason != nil && c.Rcmd.RcmdReason.Style == 3 && c.IsAttenm[a.Author.Mid] == 1 {
				c.Desc = authorname
			} else {
				c.Desc = authorname + " · " + model.PubDataString(a.PubDate.Time())
			}
			if op.CardGoto == model.CardGotoUpRcmdAv {
				button = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), IsAtten: c.IsAttenm[a.Author.Mid]}
			} else {
				if t, ok := c.Tagm[op.Tid]; ok {
					button = t
				} else {
					button = &ButtonStatus{Text: a.TypeName}
				}
			}
			c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
			if op.CardGoto == model.CardGotoPlayer && c.Base.PlayerArgs == nil {
				log.Warn("player card aid(%d) can't auto player", a.Aid)
				return
			}
			c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
			upID = a.Author.Mid
		case model.CardGotoChannelRcmd:
			t, ok := c.Tagm[op.Tid]
			if !ok {
				return
			}
			avatar = &AvatarStatus{Cover: t.Cover, Goto: model.GotoTag, Param: strconv.FormatInt(t.ID, 10), Type: model.AvatarSquare}
			c.Desc = model.SubscribeString(int32(t.Count.Atten))
			button = &ButtonStatus{Goto: model.GotoTag, Param: strconv.FormatInt(t.ID, 10), IsAtten: t.IsAtten}
			c.Base.PlayerArgs = playerArgsFrom(a.Archive3)
			c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
		case model.CardGotoAdAv:
			c.AdInfo = c.Rcmd.Ad
			avatar = &AvatarStatus{Cover: a.Author.Face, Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), Type: model.AvatarRound}
			c.Desc = a.Author.Name + " · " + model.PubDataString(a.PubDate.Time())
			button = c.Tagm[op.Tid]
			c.Args.fromArchive(a.Archive3, c.Tagm[op.Tid])
			upID = a.Author.Mid
		default:
			log.Warn("LargeCoverV1 From: unexpected card_goto %s", op.CardGoto)
			return
		}
		c.CanPlay = a.Rights.Autoplay
		if a.Rights.UGCPay == 1 && op.ShowUGCPay {
			c.CoverBadge2 = "付费"
		}
	case map[int64]*bangumi.Season:
		sm := main.(map[int64]*bangumi.Season)
		s, ok := sm[op.ID]
		if !ok {
			return
		}
		c.Base.from(s.EpisodeID, s.Cover, s.Title, model.GotoBangumi, s.EpisodeID, nil)
		c.CoverLeftText2 = model.ArchiveViewString(s.PlayCount)
		c.CoverLeftText3 = model.BangumiFavString(s.Favorites)
		avatar = &AvatarStatus{Cover: s.SeasonCover, Type: model.AvatarSquare}
		c.CoverBadge = s.TypeBadge
		c.Desc = s.UpdateDesc
	case map[int32]*season.CardInfoProto:
		sm := main.(map[int32]*season.CardInfoProto)
		s, ok := sm[int32(op.ID)]
		if !ok {
			return
		}
		c.Base.from(op.Param, s.Cover, s.Title, model.GotoPGC, op.URI, nil)
		if s.Stat != nil {
			c.CoverLeftText2 = model.ArchiveViewString(int32(s.Stat.View))
			c.CoverLeftText3 = model.BangumiFavString(int32(s.Stat.Follow))
		}
		avatar = &AvatarStatus{Cover: s.Cover, Type: model.AvatarSquare}
		c.CoverBadge = s.SeasonTypeName
		if s.NewEp != nil {
			c.Desc = s.NewEp.IndexShow
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
		if s.Season.Stat != nil {
			c.CoverLeftText2 = model.ArchiveViewString(int32(s.Season.Stat.View))
			c.CoverLeftText3 = model.BangumiFavString(int32(s.Season.Stat.Follow))
		}
		avatar = &AvatarStatus{Cover: s.Season.Cover, Type: model.AvatarSquare}
		c.CoverBadge = s.Season.SeasonTypeName
		if s.Season != nil {
			c.Desc = s.Season.NewEpShow
		}
	case map[int64]*live.Room:
		rm := main.(map[int64]*live.Room)
		r, ok := rm[op.ID]
		if !ok || r.LiveStatus != 1 {
			return
		}
		c.Base.from(strconv.FormatInt(op.ID, 10), r.Cover, r.Uname, model.GotoLive, strconv.FormatInt(r.RoomID, 10), model.LiveRoomHandler(r))
		c.CoverLeftText2 = model.LiveOnlineString(r.Online)
		avatar = &AvatarStatus{Cover: r.Cover, Goto: model.GotoMid, Param: strconv.FormatInt(r.UID, 10), Type: model.AvatarRound}
		c.CoverBadge = "直播"
		c.Desc = r.Title
		c.Base.PlayerArgs = playerArgsFrom(r)
		c.Args.fromLiveRoom(r)
		upID = r.UID
		button = r
		c.CanPlay = 1
		// SmallCoverV1
	case map[int64]*show.Shopping:
		const _buttonText = "进入"
		sm := main.(map[int64]*show.Shopping)
		s, ok := sm[op.ID]
		if !ok {
			return
		}
		c.Base.from(strconv.FormatInt(op.ID, 10), model.ShoppingCover(s.PerformanceImageP), s.Name, model.GotoWeb, s.URL, nil)
		if s.Type == 1 {
			c.CoverLeftText2 = s.Want
			c.CoverLeftText3 = s.CityName
			c.Desc = s.STime + " - " + s.ETime
		} else if s.Type == 2 {
			c.CoverLeftText2 = s.Want
			c.CoverLeftText3 = s.Subname
			c.Desc = s.Pricelt
		}
		button = &ButtonStatus{Text: _buttonText, Goto: model.GotoWeb, Param: s.URL, Type: model.ButtonTheme, Event: model.EventButtonClick}
		c.Args.fromShopping(s)
		c.CoverBadgeColor = model.PurpleCoverBadge
	case nil:
		c.Base.from(op.Param, op.Coverm[model.ColumnSvrDouble], op.Title, op.Goto, op.URI, nil)
		switch op.CardGoto {
		case model.CardGotoDownload:
			const _buttonText = "进入"
			c.Desc = op.Desc
			c.CoverLeftText2 = model.DownloadString(op.Download)
			if (op.Plat == model.PlatIPhone && op.Build > 8220) || (op.Plat == model.PlatAndroid && op.Build > 5335001) {
				button = &ButtonStatus{Text: _buttonText, Goto: op.Goto, Param: op.URI, Type: model.ButtonTheme, Event: model.EventGameClick}
			} else {
				button = &ButtonStatus{Text: _buttonText, Goto: op.Goto, Param: op.URI, Type: model.ButtonTheme, Event: model.EventButtonClick}
			}
			c.CoverBadgeColor = model.PurpleCoverBadge
		case model.CardGotoSpecial:
			c.Desc = op.Desc
			c.CoverBadge = op.Badge
			c.CoverBadgeColor = model.PurpleCoverBadge
		default:
			log.Warn("LargeCoverV1 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("LargeCoverV1 From: unexpected type %T", main)
		return
	}
	if c.Rcmd != nil {
		c.TopRcmdReason, c.BottomRcmdReason = TopBottomRcmdReason(c.Rcmd.RcmdReason, c.IsAttenm[upID], c.Cardm)
		c.TopRcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.TopRcmdReason, c.Base.Goto)
		c.BottomRcmdReasonStyle = bottomReasonStyleFrom(c.Rcmd, c.BottomRcmdReason, c.Base.Goto)
	}
	c.OfficialIcon = model.OfficialIcon(c.Cardm[upID])
	c.Avatar = avatarFrom(avatar)
	if c.Rcmd == nil || !c.Rcmd.HideButton {
		c.DescButton = buttonFrom(button, op.Plat)
	}
	c.Right = true
}

func (c *LargeCoverV1) Get() *Base {
	return c.Base
}

type SmallCoverV1 struct {
	*Base
	CoverBadge     string     `json:"cover_badge,omitempty"`
	Desc1          string     `json:"desc_1,omitempty"`
	Desc2          string     `json:"desc_2,omitempty"`
	Desc3          string     `json:"desc_3,omitempty"`
	TitleRightText string     `json:"title_right_text,omitempty"`
	TitleRightPic  model.Icon `json:"title_right_pic,omitempty"`
}

func (c *SmallCoverV1) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var button interface{}
	switch main.(type) {
	case map[int64]*audio.Audio:
		var firstSong string
		am := main.(map[int64]*audio.Audio)
		a, ok := am[op.ID]
		if !ok {
			return
		}
		if len(a.Songs) != 0 {
			firstSong = a.Songs[0].Title
		}
		c.Base.from(op.Param, a.CoverURL, a.Title, model.GotoAudio, op.URI, nil)
		c.Desc1, c.Desc2 = model.AudioDescString(firstSong, a.RecordNum)
		c.Desc3 = model.AudioPlayString(a.PlayNum) + "  " + model.AudioFavString(a.FavoriteNum)
		c.Args.fromAudio(a)
		button = a.Ctgs
	case *bangumi.Update:
		const (
			_title   = "你的追番更新啦"
			_updates = 99
		)
		u := main.(*bangumi.Update)
		if u == nil || u.Updates == 0 {
			return
		}
		c.Base.from("", u.SquareCover, _title, "", "", nil)
		updates := u.Updates
		if updates > _updates {
			updates = _updates
			c.TitleRightPic = model.IconBomb
		} else {
			c.TitleRightPic = model.IconTV
		}
		c.Desc1 = u.Title
		c.TitleRightText = strconv.Itoa(updates)
	case map[int64]*show.Shopping:
		const _buttonText = "进入"
		sm := main.(map[int64]*show.Shopping)
		s, ok := sm[op.ID]
		if !ok {
			return
		}
		c.Base.from(strconv.FormatInt(op.ID, 10), model.ShoppingCover(s.PerformanceImageP), s.Name, model.GotoWeb, s.URL, nil)
		if s.Type == 1 {
			c.Desc1 = s.STime + " - " + s.ETime
			c.Desc2 = s.CityName
			c.Desc3 = "￥" + s.Pricelt
		} else if s.Type == 2 {
			c.Desc1 = s.Subname
			c.Desc2 = s.Want
			c.Desc3 = s.Pricelt
		}
		button = &ButtonStatus{Text: _buttonText, Goto: model.GotoWeb, Param: s.URL, Type: model.ButtonTheme, Event: model.EventButtonClick}
		c.Args.fromShopping(s)
	case *cm.AdInfo:
		ad := main.(*cm.AdInfo)
		c.AdInfo = ad
	case *bangumi.Moe:
		m := main.(*bangumi.Moe)
		if m == nil {
			return
		}
		c.Base.from(strconv.FormatInt(m.ID, 10), m.Square, m.Title, model.GotoWeb, m.Link, nil)
		c.Desc1 = m.Desc
		c.CoverBadge = m.Badge
	case nil:
		c.Base.from(op.Param, op.Coverm[c.Columnm], op.Title, op.Goto, op.URI, nil)
		switch op.CardGoto {
		case model.CardGotoDownload:
			const _buttonText = "进入"
			c.Desc1 = op.Desc
			c.Desc2 = model.DownloadString(op.Download)
			if (op.Plat == model.PlatIPhone && op.Build > 8220) || (op.Plat == model.PlatAndroid && op.Build > 5335001) {
				button = &ButtonStatus{Text: _buttonText, Goto: op.Goto, Param: op.URI, Type: model.ButtonTheme, Event: model.EventGameClick}
			} else {
				button = &ButtonStatus{Text: _buttonText, Goto: op.Goto, Param: op.URI, Type: model.ButtonTheme, Event: model.EventButtonClick}
			}
		case model.CardGotoSpecial:
			c.Desc1 = op.Desc
			c.CoverBadge = op.Badge
		default:
			log.Warn("SmallCoverV1 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("SmallCoverV1 From: unexpected type %T", main)
		return
	}
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *SmallCoverV1) Get() *Base {
	return c.Base
}

type MiddleCover struct {
	*Base
	Ratio int    `json:"ratio,omitempty"`
	Badge string `json:"badge,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

func (c *MiddleCover) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case *cm.AdInfo:
		ad := main.(*cm.AdInfo)
		c.AdInfo = ad
	case nil:
		if op == nil {
			return
		}
		c.Base.from(op.Param, op.Coverm[c.Columnm], op.Title, op.Goto, op.URI, nil)
		switch op.CardGoto {
		case model.CardGotoSpecial:
			c.Desc = op.Desc
			c.Badge = op.Badge
			c.Ratio = op.Ratio
		default:
			log.Warn("MiddleCover From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("MiddleCover From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *MiddleCover) Get() *Base {
	return c.Base
}

type Topstick struct {
	*Base
	Desc string `json:"desc,omitempty"`
}

func (c *Topstick) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case nil:
		if op == nil {
			return
		}
		c.Base.from(op.Param, op.Coverm[c.Columnm], op.Title, op.Goto, op.URI, nil)
		switch op.CardGoto {
		case model.CardGotoTopstick:
			c.Desc = op.Desc
		default:
			return
		}
	}
	c.Right = true
}

func (c *Topstick) Get() *Base {
	return c.Base
}

type ThreeItemV1 struct {
	*Base
	TitleIcon   model.Icon         `json:"title_icon,omitempty"`
	BannerCover string             `json:"banner_cover,omitempty"`
	BannerURI   string             `json:"banner_uri,omitempty"`
	MoreURI     string             `json:"more_uri,omitempty"`
	MoreText    string             `json:"more_text,omitempty"`
	Items       []*ThreeItemV1Item `json:"items,omitempty"`
}

type ThreeItemV1Item struct {
	Base
	CoverLeftText string     `json:"cover_left_text,omitempty"`
	CoverLeftIcon model.Icon `json:"cover_left_icon,omitempty"`
	Desc1         string     `json:"desc_1,omitempty"`
	Desc2         string     `json:"desc_2,omitempty"`
	Badge         string     `json:"badge,omitempty"`
}

func (c *ThreeItemV1) From(main interface{}, op *operate.Card) {
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
			// c.TitleIcon = model.IconRank
			c.MoreURI = model.FillURI(op.Goto, op.URI, nil)
			c.MoreText = "查看更多"
			c.Items = make([]*ThreeItemV1Item, 0, _limit)
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				am := intfc.(map[int64]*archive.ArchiveWithPlayer)
				a, ok := am[v.ID]
				if !ok || !model.AvIsNormal(a) {
					continue
				}
				item := &ThreeItemV1Item{
					CoverLeftText: model.DurationString(a.Duration),
					Desc1:         model.ScoreString(v.Score),
				}
				item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, v.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
				item.Args.fromArchive(a.Archive3, nil)
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
			c.Items = make([]*ThreeItemV1Item, 0, len(op.Items))
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				var item *ThreeItemV1Item
				switch intfc.(type) {
				case map[int64]*archive.ArchiveWithPlayer:
					am := intfc.(map[int64]*archive.ArchiveWithPlayer)
					a, ok := am[v.ID]
					if !ok || !model.AvIsNormal(a) {
						continue
					}
					item = &ThreeItemV1Item{
						CoverLeftText: model.DurationString(a.Duration),
						Desc1:         model.ArchiveViewString(a.Stat.View),
						Desc2:         model.DanmakuString(a.Stat.Danmaku),
					}
					if op.SwitchLike == model.SwitchFeedIndexLike {
						item.Desc1 = model.LikeString(a.Stat.Like)
						item.Desc2 = model.ArchiveViewString(a.Stat.View)
					}
					item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, v.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
					item.Args.fromArchive(a.Archive3, nil)
				case map[int64]*live.Room:
					rm := intfc.(map[int64]*live.Room)
					r, ok := rm[v.ID]
					if !ok || r.LiveStatus != 1 {
						continue
					}
					item = &ThreeItemV1Item{
						Desc1: model.LiveOnlineString(r.Online),
						Badge: "直播",
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
					item = &ThreeItemV1Item{
						Badge: "文章",
					}
					item.Base.from(v.Param, m.ImageURLs[0], m.Title, model.GotoArticle, v.URI, nil)
					if m.Stats != nil {
						item.Desc1 = model.ArticleViewString(m.Stats.View)
						item.Desc2 = model.ArticleReplyString(m.Stats.Reply)
					}
					item.Args.fromArticle(m)
				default:
					log.Warn("ThreeItemV1 From: unexpected type %T", intfc)
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
			log.Warn("ThreeItemV1 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("ThreeItemV1 From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *ThreeItemV1) Get() *Base {
	return c.Base
}

type ThreeItemH struct {
	*Base
	Items []*ThreeItemHItem `json:"items,omitempty"`
}

type ThreeItemHItem struct {
	Base
	CoverType    model.Type `json:"cover_type,omitempty"`
	Desc         string     `json:"desc,omitempty"`
	OfficialIcon model.Icon `json:"official_icon,omitempty"`
}

func (c *ThreeItemH) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case nil:
		if op == nil {
			return
		}
		switch op.CardGoto {
		case model.CardGotoSubscribe, model.CardGotoSearchSubscribe:
			const _limit = 3
			c.Base.from(op.Param, "", op.Title, "", "", nil)
			c.Items = make([]*ThreeItemHItem, 0, _limit)
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				var (
					item   *ThreeItemHItem
					button interface{}
				)
				switch v.Goto {
				case model.GotoTag:
					t, ok := c.Tagm[v.ID]
					if !ok || t.IsAtten == 1 {
						continue
					}
					item = &ThreeItemHItem{
						CoverType: model.AvatarSquare,
						Desc:      model.SubscribeString(int32(t.Count.Atten)),
					}
					item.Base.from(v.Param, t.Cover, t.Name, v.Goto, v.URI, nil)
					button = &ButtonStatus{Goto: model.GotoTag, Param: strconv.FormatInt(t.ID, 10)}
				case model.GotoMid:
					cd, ok := c.Cardm[v.ID]
					if !ok || c.IsAttenm[v.ID] == 1 {
						continue
					}
					item = &ThreeItemHItem{
						CoverType: model.AvatarRound,
					}
					item.Base.from(v.Param, cd.Face, cd.Name, v.Goto, v.URI, nil)
					button = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(cd.Mid, 10)}
					if v.Desc != "" {
						item.Desc = v.Desc
					} else if stat, ok := c.Statm[cd.Mid]; ok {
						item.Desc = model.FanString(int32(stat.Follower))
					}
					item.OfficialIcon = model.OfficialIcon(cd)
				default:
					log.Warn("ThreeItemH From: unexpected type %T", v.Goto)
					continue
				}
				item.DescButton = buttonFrom(button, op.Plat)
				c.Items = append(c.Items, item)
				if len(c.Items) == _limit {
					break
				}
			}
			if len(c.Items) < _limit {
				return
			}
		default:
			log.Warn("ThreeItemH From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("ThreeItemH From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *ThreeItemH) Get() *Base {
	return c.Base
}

type ThreeItemHV3 struct {
	*Base
	Covers        []string   `json:"covers,omitempty"`
	CoverTopText1 string     `json:"cover_top_text_1,omitempty"`
	CoverTopText2 string     `json:"cover_top_text_2,omitempty"`
	Desc          string     `json:"desc,omitempty"`
	Avatar        *Avatar    `json:"avatar,omitempty"`
	OfficialIcon  model.Icon `json:"official_icon,omitempty"`
}

func (c *ThreeItemHV3) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		upID int64
	)
	switch main.(type) {
	case map[int64]*article.Meta:
		mm := main.(map[int64]*article.Meta)
		m, ok := mm[op.ID]
		if !ok {
			return
		}
		c.Base.from(op.Param, "", m.Title, model.GotoArticle, op.URI, nil)
		c.Covers = m.ImageURLs
		c.CoverTopText1 = model.ArticleViewString(m.Stats.View)
		c.CoverTopText2 = model.ArticleReplyString(m.Stats.Reply)
		c.Desc = m.Summary
		if m.Author != nil {
			c.Avatar = avatarFrom(&AvatarStatus{Cover: m.Author.Face, Text: m.Author.Name + "·" + model.PubDataString(m.PublishTime.Time()), Goto: model.GotoMid, Param: strconv.FormatInt(m.Author.Mid, 10), Type: model.AvatarRound})
			upID = m.Author.Mid
		}
		c.Args.fromArticle(m)
	default:
		log.Warn("ThreeItemHV3 From: unexpected type %T", main)
		return
	}
	c.OfficialIcon = model.OfficialIcon(c.Cardm[upID])
	c.Right = true
}

func (c *ThreeItemHV3) Get() *Base {
	return c.Base
}

type TwoItemV1 struct {
	*Base
	Items []*TwoItemV1Item `json:"items,omitempty"`
}

type TwoItemV1Item struct {
	Base
	CoverBadge     string `json:"cover_badge,omitempty"`
	CoverLeftText1 string `json:"cover_left_text_1,omitempty"`
}

func (c *TwoItemV1) From(main interface{}, op *operate.Card) {
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
		c.Items = make([]*TwoItemV1Item, 0, _limit)
		for _, card := range cs {
			if card == nil || card.LiveStatus != 1 {
				continue
			}
			item := &TwoItemV1Item{
				CoverBadge:     "直播",
				CoverLeftText1: model.LiveOnlineString(card.Online),
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

func (c *TwoItemV1) Get() *Base {
	return c.Base
}

type CoverOnly struct {
	*Base
}

func (c *CoverOnly) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case nil:
		switch op.CardGoto {
		case model.CardGotoLogin:
			c.Base.from(op.Param, "", "", "", "", nil)
		}
	default:
		log.Warn("CoverOnly From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *CoverOnly) Get() *Base {
	return c.Base
}

type Banner struct {
	*Base
	Hash       string           `json:"hash,omitempty"`
	BannerItem []*banner.Banner `json:"banner_item,omitempty"`
}

func (c *Banner) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case nil:
		switch op.CardGoto {
		case model.CardGotoBanner:
			if len(op.Banner) == 0 {
				log.Warn("Banner len is null")
				return
			}
			c.BannerItem = op.Banner
			c.Hash = op.Hash
		default:
			log.Warn("Banner From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("Banner From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *Banner) Get() *Base {
	return c.Base
}

type HotTopic struct {
	*Base
	Desc  string          `json:"desc,omitempty"`
	Items []*HotTopicItem `json:"items,omitempty"`
}

type HotTopicItem struct {
	Cover string `json:"cover,omitempty"`
	URI   string `json:"uri,omitempty"`
	Param string `json:"param,omitempty"`
	Name  string `json:"name,omitempty"`
}

func (c *HotTopic) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case []*live.TopicHot:
		th := main.([]*live.TopicHot)
		if len(th) == 0 {
			return
		}
		items := make([]*HotTopicItem, 0, len(th))
		for _, t := range th {
			it := &HotTopicItem{
				Name:  t.TName,
				Param: strconv.Itoa(t.TID),
				Cover: t.ImageURL,
				URI:   model.FillURI(model.GotoHotTopic, strconv.Itoa(t.TID), model.HottopicHandler(t)),
			}
			items = append(items, it)
		}
		c.Items = items
		c.Base.from("0", "", "热门话题", model.GotoWeb, "bilibili://following/hot_topic_list", nil)
		c.Desc = "更多热门话题"
	default:
		log.Warn("HotTopic From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *HotTopic) Get() *Base {
	return c.Base
}

type Text struct {
	*Base
	Content string `json:"content,omitempty"`
}

func (c *Text) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case nil:
		switch op.CardGoto {
		case model.CardGotoNews:
			c.Base.from(op.Param, "", op.Title, model.GotoWeb, op.URI, nil)
			c.Content = op.Desc
		default:
			log.Warn("Text From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("Text From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *Text) Get() *Base {
	return c.Base
}

type ThreeItemHV4 struct {
	*Base
	MoreURI  string              `json:"more_uri,omitempty"`
	MoreText string              `json:"more_text,omitempty"`
	Items    []*ThreeItemHV4Item `json:"items,omitempty"`
}

type ThreeItemHV4Item struct {
	Cover           string           `json:"cover,omitempty"`
	Title           string           `json:"title,omitempty"`
	Desc            string           `json:"desc,omitempty"`
	Goto            model.Gt         `json:"goto,omitempty"`
	Param           string           `json:"param,omitempty"`
	URI             string           `json:"uri,omitempty"`
	CoverBadge      string           `json:"cover_badge,omitempty"`
	CoverBadgeColor model.CoverColor `json:"cover_badge_color,omitempty"`
}

func (c *ThreeItemHV4) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case map[int32]*season.CardInfoProto:
		const _limit = 3
		c.Base.from(op.Param, "", op.Title, op.Goto, "", nil)
		c.Items = make([]*ThreeItemHV4Item, 0, _limit)
		for _, v := range op.Items {
			if v == nil {
				continue
			}
			var (
				item *ThreeItemHV4Item
			)
			sm := main.(map[int32]*season.CardInfoProto)
			s, ok := sm[int32(v.ID)]
			if !ok {
				return
			}
			item = &ThreeItemHV4Item{
				Title:      s.Title,
				Cover:      s.Cover,
				Goto:       model.GotoPGC,
				URI:        model.FillURI(model.GotoPGC, strconv.FormatInt(int64(s.SeasonId), 10), nil),
				Param:      strconv.FormatInt(int64(s.SeasonId), 10),
				CoverBadge: s.Badge,
				// CoverBadgeColor: model.PurpleCoverBadge,
				// Desc:SeasonTypeName + " · " +
			}
			if s.Rating != nil && s.Rating.Score > 0 {
				item.Desc = fmt.Sprintf("%s · %.1f分", s.SeasonTypeName, s.Rating.Score)
			}
			c.Items = append(c.Items, item)
			if len(c.Items) == _limit {
				break
			}
		}
		if len(c.Items) > _limit {
			// c.MoreText = "查看更多"
			// c.MoreURI = model.FillURI(op.Goto, op.URI, nil)
			c.Items = c.Items[:_limit]
		}
		if len(c.Items) < _limit {
			return
		}
	default:
		log.Warn("ThreeItemHV4Item From: unexpected card_goto %s", op.CardGoto)
		return
	}
	c.Right = true
	return
}

func (c *ThreeItemHV4) Get() *Base {
	return c.Base
}

type UpRcmdCover struct {
	*Base
	CoverType    model.Type `json:"cover_type,omitempty"`
	Level        int32      `json:"level,omitempty"`
	OfficialIcon model.Icon `json:"official_icon,omitempty"`
	DescButton   *Button    `json:"desc_button,omitempty"`
	Desc1        string     `json:"desc_1,omitempty"`
	Desc2        string     `json:"desc_2,omitempty"`
	Desc3        string     `json:"desc_3,omitempty"`
}

func (c *UpRcmdCover) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case nil:
		switch op.CardGoto {
		case model.CardGotoUpRcmdS:
			c.Base.from(strconv.FormatInt(op.ID, 10), "", "", model.GotoMid, strconv.FormatInt(op.ID, 10), nil)
			var (
				button interface{}
			)
			cd, ok := c.Cardm[op.ID]
			if !ok {
				return
			}
			c.Cover = cd.Face
			c.CoverType = model.AvatarRound
			c.Title = cd.Name
			c.Level = cd.Level
			c.OfficialIcon = model.OfficialIcon(cd)
			if stat, ok := c.Statm[cd.Mid]; ok {
				c.Desc1 = "粉丝: " + model.StatString(int32(stat.Follower), "")
			}
			c.Desc2 = "视频: " + strconv.Itoa(op.Limit)
			c.Desc3 = cd.Sign
			button = &ButtonStatus{
				Goto:    model.GotoMid,
				Param:   strconv.FormatInt(cd.Mid, 10),
				IsAtten: c.IsAttenm[op.ID],
				Event:   model.EventUpClick,
			}
			c.DescButton = buttonFrom(button, op.Plat)
		default:
			log.Warn("UpRcmdCover From: unexpected card_goto %s", op.CardGoto)
			return
		}
		c.Right = true
	}
}

func (c *UpRcmdCover) Get() *Base {
	return c.Base
}

type ThreeItemAll struct {
	*Base
	Items []*ThreeItemAllItem `json:"items,omitempty"`
}

type ThreeItemAllItem struct {
	Base
	CoverType    model.Type `json:"cover_type,omitempty"`
	Desc         string     `json:"desc,omitempty"`
	DescButton   *Button    `json:"desc_button,omitempty"`
	OfficialIcon model.Icon `json:"official_icon,omitempty"`
}

func (c *ThreeItemAll) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case nil:
		if op == nil {
			return
		}
		switch op.CardGoto {
		case model.CardGotoSearchUpper:
			const _limit = 3
			c.Base.from(op.Param, "", op.Title, "", "", nil)
			c.Items = []*ThreeItemAllItem{}
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				var (
					item   *ThreeItemAllItem
					button interface{}
				)
				switch v.Goto {
				case model.GotoMid:
					cd, ok := c.Cardm[v.ID]
					if !ok || c.IsAttenm[v.ID] == 1 {
						continue
					}
					item = &ThreeItemAllItem{
						CoverType: model.AvatarRound,
					}
					item.Base.from(v.Param, cd.Face, cd.Name, v.Goto, v.URI, nil)
					button = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(cd.Mid, 10)}
					if v.Desc != "" {
						item.Desc = v.Desc
					} else if stat, ok := c.Statm[cd.Mid]; ok {
						item.Desc = model.FanString(int32(stat.Follower))
					}
					item.OfficialIcon = model.OfficialIcon(cd)
					if item.OfficialIcon == 0 {
						switch cd.Vip.Type {
						case 1:
							item.OfficialIcon = model.IconRoleVipRed
						case 2:
							item.OfficialIcon = model.IconRoleYearVipRed
						}
					}
				default:
					log.Warn("ThreeItemAll From: unexpected type %T", v.Goto)
					continue
				}
				item.DescButton = buttonFrom(button, op.Plat)
				c.Items = append(c.Items, item)
			}
			if len(c.Items) < _limit {
				return
			}
		default:
			log.Warn("ThreeItemAll From: unexpected card_goto %s", op.CardGoto)
			return
		}
		c.Right = true
	}
}

func (c *ThreeItemAll) Get() *Base {
	return c.Base
}

type ChannelSquare struct {
	*Base
	Desc1 string               `json:"desc_1,omitempty"`
	Desc2 string               `json:"desc_2,omitempty"`
	Item  []*ChannelSquareItem `json:"item,omitempty"`
}

type ChannelSquareItem struct {
	Title          string     `json:"title,omitempty"`
	Cover          string     `json:"cover,omitempty"`
	URI            string     `json:"uri,omitempty"`
	Param          string     `json:"param,omitempty"`
	Goto           string     `json:"goto,omitempty"`
	CoverLeftText1 string     `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1 model.Icon `json:"cover_left_icon_1,omitempty"`
	CoverLeftText2 string     `json:"cover_left_text_2,omitempty"`
	CoverLeftIcon2 model.Icon `json:"cover_left_icon_2,omitempty"`
	CoverLeftText3 string     `json:"cover_left_text_3,omitempty"`
	FromType       string     `json:"from_type"`
}

//From ChannelSquare op:channel--av对应关系, main:av map, c.base.tagm:tag map
func (c *ChannelSquare) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case map[int64]*api.Arc:
		t := c.Base.Tagm[op.ID]
		c.Base.from(op.Param, t.Cover, t.Name, model.GotoTag, op.Param, nil)
		button := &ButtonStatus{Goto: model.GotoTag, IsAtten: t.IsAtten}
		c.DescButton = buttonFrom(button, op.Plat)
		c.Desc1 = t.Content
		c.Desc2 = model.SubscribeString(int32(t.Count.Atten))
		for _, item := range op.Items {
			am := main.(map[int64]*api.Arc)
			av := am[item.ID]
			c.Item = append(c.Item, &ChannelSquareItem{
				Title:          av.Title,
				Cover:          av.Pic,
				URI:            model.FillURI(model.GotoAv, strconv.FormatInt(item.ID, 10), model.AvPlayHandler(archive.BuildArchive3(av), nil, "")),
				Goto:           string(model.GotoAv),
				Param:          strconv.FormatInt(item.ID, 10),
				CoverLeftText1: model.StatString(av.Stat.View, ""),
				CoverLeftIcon1: model.IconPlay,
				CoverLeftText2: model.StatString(av.Stat.Danmaku, ""),
				CoverLeftIcon2: model.IconDanmaku,
				CoverLeftText3: model.DurationString(av.Duration),
				FromType:       item.FromType,
			})
		}
	}
	c.Right = true
}

func (c *ChannelSquare) Get() *Base {
	return c.Base
}

type TwoItemHV1 struct {
	*Base
	Desc       string            `json:"desc,omitempty"`
	DescButton *Button           `json:"desc_button,omitempty"`
	Items      []*TwoItemHV1Item `json:"item,omitempty"`
}

type TwoItemHV1Item struct {
	Title          string     `json:"title,omitempty"`
	Cover          string     `json:"cover,omitempty"`
	URI            string     `json:"uri,omitempty"`
	Param          string     `json:"param,omitempty"`
	Args           Args       `json:"args,omitempty"`
	Goto           string     `json:"goto,omitempty"`
	CoverLeftText1 string     `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1 model.Icon `json:"cover_left_icon_1,omitempty"`
	CoverRightText string     `json:"cover_right_text,omitempty"`
}

func (c *TwoItemHV1) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	au, ok := c.Cardm[op.ID]
	if !ok {
		return
	}
	c.Base.from(op.Param, au.Face, au.Name, model.GotoMid, op.Param, nil)
	button := &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(op.ID, 10), IsAtten: c.IsAttenm[op.ID]}
	c.DescButton = buttonFrom(button, op.Plat)
	if op.Desc != "" {
		c.Desc = op.Desc
	} else {
		c.Desc = au.Sign
	}
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		for _, item := range op.Items {
			am := main.(map[int64]*archive.ArchiveWithPlayer)
			var (
				a  *archive.ArchiveWithPlayer
				ok bool
			)
			if a, ok = am[item.ID]; !ok {
				continue
			}
			args := Args{}
			args.fromArchive(a.Archive3, c.Tagm[op.Tid])
			c.Items = append(c.Items, &TwoItemHV1Item{
				Title:          a.Title,
				Cover:          a.Pic,
				URI:            model.FillURI(model.GotoAv, strconv.FormatInt(item.ID, 10), model.AvPlayHandler(a.Archive3, nil, "")),
				Goto:           string(model.GotoAv),
				Param:          strconv.FormatInt(item.ID, 10),
				CoverLeftText1: model.StatString(a.Stat.View, ""),
				CoverLeftIcon1: model.IconPlay,
				CoverRightText: model.DurationString(a.Duration),
				Args:           args,
			})
			if len(c.Items) >= 2 {
				break
			}
		}
		if len(c.Items) < 2 {
			return
		}
	}
	c.Right = true
}

func (c *TwoItemHV1) Get() *Base {
	return c.Base
}

type OnePicV1 struct {
	*Base
	Desc1                     string
	Desc2                     string
	Avatar                    *Avatar          `json:"avatar,omitempty"`
	CoverLeftText1            string           `json:"cover_left_text_1,omitempty"`
	CoverLeftText2            string           `json:"cover_left_text_2,omitempty"`
	CoverRightText            string           `json:"cover_right_text,omitempty"`
	CoverRightBackgroundColor string           `json:"cover_right_background_color,omitempty"`
	CoverBadge                string           `json:"cover_badge,omitempty"`
	TopRcmdReason             string           `json:"top_rcmd_reason,omitempty"`
	BottomRcmdReason          string           `json:"bottom_rcmd_reason,omitempty"`
	Desc                      string           `json:"desc,omitempty"`
	OfficialIcon              model.Icon       `json:"official_icon,omitempty"`
	CanPlay                   int32            `json:"can_play,omitempty"`
	CoverBadgeColor           model.CoverColor `json:"cover_badge_color,omitempty"`
	TopRcmdReasonStyle        *ReasonStyle     `json:"top_rcmd_reason_style,omitempty"`
	BottomRcmdReasonStyle     *ReasonStyle     `json:"bottom_rcmd_reason_style,omitempty"`
}

func (c *OnePicV1) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button interface{}
		avatar *AvatarStatus
		upID   int64
	)
	switch main.(type) {
	case map[int64]*bplus.Picture:
		pm := main.(map[int64]*bplus.Picture)
		p, ok := pm[op.ID]
		if !ok || len(p.Imgs) == 0 || p.ViewCount == 0 {
			return
		}
		c.Base.from(op.Param, p.Imgs[0], p.DynamicText, model.GotoPicture, strconv.FormatInt(p.DynamicID, 10), nil)
		c.CoverLeftText1 = model.PictureViewString(p.ViewCount)
		c.CoverLeftText2 = model.ArticleReplyString(p.CommentCount)
		if p.ImgCount > 1 {
			c.CoverRightText = model.PictureCountString(p.ImgCount)
			c.CoverRightBackgroundColor = "#66666666"
		}
		c.Desc1 = p.NickName
		c.Desc2 = model.PubDataString(p.PublishTime.Time())
		avatar = &AvatarStatus{Cover: p.FaceImg, Goto: model.GotoDynamicMid, Param: strconv.FormatInt(p.Mid, 10), Type: model.AvatarRound}
		button = p
		upID = p.Mid
	default:
		log.Warn("OnePicV1 From: unexpected type %T", main)
	}
	if c.Rcmd != nil {
		c.TopRcmdReason, c.BottomRcmdReason = TopBottomRcmdReason(c.Rcmd.RcmdReason, c.IsAttenm[upID], c.Cardm)
		c.TopRcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.TopRcmdReason, c.Base.Goto)
		c.BottomRcmdReasonStyle = bottomReasonStyleFrom(c.Rcmd, c.BottomRcmdReason, c.Base.Goto)
	}
	c.OfficialIcon = model.OfficialIcon(c.Cardm[upID])
	c.Avatar = avatarFrom(avatar)
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *OnePicV1) Get() *Base {
	return c.Base
}

type ThreePicV1 struct {
	*Base
	Covers                    []string     `json:"covers,omitempty"`
	Desc1                     string       `json:"desc_1,omitempty"`
	Desc2                     string       `json:"desc_2,omitempty"`
	Avatar                    *Avatar      `json:"avatar,omitempty"`
	TitleLeftText1            string       `json:"title_left_text_1,omitempty"`
	TitleLeftText2            string       `json:"title_left_text_2,omitempty"`
	CoverRightText            string       `json:"cover_right_text,omitempty"`
	CoverRightBackgroundColor string       `json:"cover_right_background_color,omitempty"`
	TopRcmdReason             string       `json:"top_rcmd_reason,omitempty"`
	BottomRcmdReason          string       `json:"bottom_rcmd_reason,omitempty"`
	TopRcmdReasonStyle        *ReasonStyle `json:"top_rcmd_reason_style,omitempty"`
	BottomRcmdReasonStyle     *ReasonStyle `json:"bottom_rcmd_reason_style,omitempty"`
}

func (c *ThreePicV1) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button interface{}
		avatar *AvatarStatus
		upID   int64
	)
	switch main.(type) {
	case map[int64]*bplus.Picture:
		pm := main.(map[int64]*bplus.Picture)
		p, ok := pm[op.ID]
		if !ok || len(p.Imgs) < 3 || p.ViewCount == 0 {
			return
		}
		c.Base.from(op.Param, "", p.DynamicText, model.GotoPicture, strconv.FormatInt(p.DynamicID, 10), nil)
		c.Covers = p.Imgs[:3]
		c.TitleLeftText1 = model.PictureViewString(p.ViewCount)
		c.TitleLeftText2 = model.ArticleReplyString(p.CommentCount)
		if p.ImgCount > 3 {
			c.CoverRightText = model.PictureCountString(p.ImgCount)
			c.CoverRightBackgroundColor = "#66666666"
		}
		c.Desc1 = p.NickName
		c.Desc2 = model.PubDataString(p.PublishTime.Time())
		avatar = &AvatarStatus{Cover: p.FaceImg, Goto: model.GotoDynamicMid, Param: strconv.FormatInt(p.Mid, 10), Type: model.AvatarRound}
		button = p
		upID = p.Mid
	default:
		log.Warn("ThreePicV1 From: unexpected type %T", main)
	}
	if c.Rcmd != nil {
		c.TopRcmdReason, c.BottomRcmdReason = TopBottomRcmdReason(c.Rcmd.RcmdReason, c.IsAttenm[upID], c.Cardm)
		c.TopRcmdReasonStyle = topReasonStyleFrom(c.Rcmd, c.TopRcmdReason, c.Base.Goto)
		c.BottomRcmdReasonStyle = bottomReasonStyleFrom(c.Rcmd, c.BottomRcmdReason, c.Base.Goto)
	}
	c.Avatar = avatarFrom(avatar)
	c.DescButton = buttonFrom(button, op.Plat)
	c.Right = true
}

func (c *ThreePicV1) Get() *Base {
	return c.Base
}

type SmallCoverV5 struct {
	*Base
	Up              *Up          `json:"up,omitempty"`
	CoverRightText1 string       `json:"cover_right_text_1,omitempty"`
	RightDesc1      string       `json:"right_desc_1,omitempty"`
	RightDesc2      string       `json:"right_desc_2,omitempty"`
	CanPlay         int32        `json:"can_play,omitempty"`
	RcmdReasonStyle *ReasonStyle `json:"rcmd_reason_style,omitempty"`
}

type Up struct {
	ID           int64      `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	Desc         string     `json:"desc,omitempty"`
	Avatar       *Avatar    `json:"avatar,omitempty"`
	OfficialIcon model.Icon `json:"official_icon,omitempty"`
	DescButton   *Button    `json:"desc_button,omitempty"`
	Cooperation  string     `json:"cooperation,omitempty"`
}

func (c *SmallCoverV5) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	var (
		button     interface{}
		avatar     *AvatarStatus
		rcmdReason string
	)
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		am := main.(map[int64]*archive.ArchiveWithPlayer)
		a, ok := am[op.ID]
		if !ok || !model.AvIsNormal(a) {
			return
		}
		c.Base.from(op.Param, a.Pic, a.Title, model.GotoAv, op.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
		c.CoverRightText1 = model.DurationString(a.Duration)
		if c.Rcmd != nil {
			rcmdReason, _ = TopBottomRcmdReason(c.Rcmd.RcmdReason, c.IsAttenm[a.Author.Mid], c.Cardm)
			c.RcmdReasonStyle = topReasonStyleFrom(c.Rcmd, rcmdReason, c.Base.Goto)
		}
		switch op.CardGoto {
		case model.CardGotoAv:
			var (
				authorface = a.Author.Face
				authorname = a.Author.Name
			)
			if (authorface == "" || authorname == "") && c.Cardm != nil {
				if au, ok := c.Cardm[a.Author.Mid]; ok {
					authorface = au.Face
					authorname = au.Name
				}
			}
			switch c.Rcmd.Style {
			case model.HotCardStyleShowUp:
				c.Up = &Up{
					ID:   a.Author.Mid,
					Name: authorname,
				}
				if stat, ok := c.Statm[a.Author.Mid]; ok {
					c.Up.Desc = model.AttentionString(int32(stat.Follower))
				}
				avatar = &AvatarStatus{Cover: authorface, Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), Type: model.AvatarRound}
				c.Up.Avatar = avatarFrom(avatar)
				c.Up.OfficialIcon = model.OfficialIcon(c.Cardm[a.Author.Mid])
				c.RightDesc1 = model.ArchiveViewString(a.Stat.View) + " · " + model.PubDataString(a.PubDate.Time())
				button = &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(a.Author.Mid, 10), IsAtten: c.IsAttenm[a.Author.Mid]}
				c.Up.DescButton = buttonFrom(button, op.Plat)
				if a.Rights.IsCooperation > 0 {
					c.Up.Cooperation = "等联合创作"
				}
			default:
				if op.Switch != model.SwitchCooperationHide {
					c.RightDesc1 = unionAuthor(a)
				} else {
					c.RightDesc1 = authorname
				}
				c.RightDesc2 = model.ArchiveViewString(a.Stat.View) + " · " + model.PubDataString(a.PubDate.Time())
			}
			// c.CanPlay = a.Rights.Autoplay
		default:
			log.Warn("SmallCoverV5 From: unexpected type %T", main)
			return
		}
	}
	c.Right = true
}

func (c *SmallCoverV5) Get() *Base {
	return c.Base
}

// Option struct.
type Option struct {
	*Base
	Option []string `json:"option,omitempty"`
}

// From is.
func (c *Option) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case []string:
		os := main.([]string)
		if len(os) == 0 {
			return
		}
		c.Base.from(op.Param, "", "选择感兴趣的内容", "", "", nil)
		c.Option = os
		c.DescButton = &Button{Text: "选好啦，刷新首页"}
	default:
		log.Warn("Option From: unexpected type %T", main)
		return
	}
	c.Right = true
}

// Get is.
func (c *Option) Get() *Base {
	return c.Base
}

type DynamicHot struct {
	*Base
	TopLeftTitle       string       `json:"top_left_title,omitempty"`
	Desc1              string       `json:"desc1,omitempty"`
	Desc2              string       `json:"desc2,omitempty"`
	MoreURI            string       `json:"more_uri,omitempty"`
	MoreText           string       `json:"more_text,omitempty"`
	Covers             []string     `json:"covers,omitempty"`
	CoverRightText     string       `json:"cover_right_text,omitempty"`
	TopRcmdReasonStyle *ReasonStyle `json:"top_rcmd_reason_style,omitempty"`
}

func (c *DynamicHot) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case []*live.DynamicHot:
		lds := main.([]*live.DynamicHot)
		if len(lds) == 0 {
			return
		}
		ld := lds[0]
		c.Base.from("", "", "", model.GotoHotDynamic, strconv.FormatInt(ld.ID, 10), nil)
		c.TopLeftTitle = "热门动态"
		c.MoreURI = "bilibili://following/recommend"
		c.MoreText = "查看更多"
		c.Title = ld.DynamicText
		if len(ld.Imgs) < 3 {
			return
		}
		c.Covers = ld.Imgs[:3]
		c.CoverRightText = strconv.Itoa(ld.ImgCount) + "P"
		c.Desc1 = ld.NickName
		var tmpdesc string
		if ld.ViewCount > 0 {
			tmpdesc = model.PictureViewString(ld.ViewCount)
		}
		if tmpdesc != "" && ld.CommentCount > 0 {
			tmpdesc = tmpdesc + " " + model.ArticleReplyString(ld.CommentCount)
		} else if ld.CommentCount > 0 {
			tmpdesc = model.ArticleReplyString(ld.CommentCount)
		}
		if tmpdesc != "" {
			c.Desc2 = model.PictureViewString(ld.ViewCount) + " " + model.ArticleReplyString(ld.CommentCount)
		}
		if ld.RcmdReason != "" {
			c.TopRcmdReasonStyle = reasonStyleFrom(model.BgColorOrange, ld.RcmdReason)
		}
	default:
		log.Warn("DynamicHot From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *DynamicHot) Get() *Base {
	return c.Base
}

type ThreeItemAllV2 struct {
	*Base
	DescButton         *Button           `json:"desc_button,omitempty"`
	TopRcmdReasonStyle *ReasonStyle      `json:"top_rcmd_reason_style,omitempty"`
	Items              []*TwoItemHV1Item `json:"item,omitempty"`
}

type ThreeItemAllV2Item struct {
	Title          string     `json:"title,omitempty"`
	Cover          string     `json:"cover,omitempty"`
	URI            string     `json:"uri,omitempty"`
	Param          string     `json:"param,omitempty"`
	Args           Args       `json:"args,omitempty"`
	Goto           string     `json:"goto,omitempty"`
	CoverLeftText1 string     `json:"cover_left_text_1,omitempty"`
	CoverLeftIcon1 model.Icon `json:"cover_left_icon_1,omitempty"`
	CoverRightText string     `json:"cover_right_text,omitempty"`
}

func (c *ThreeItemAllV2) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	au, ok := c.Cardm[op.ID]
	if !ok {
		return
	}
	c.Base.from(op.Param, au.Face, au.Name, model.GotoMid, op.Param, nil)
	button := &ButtonStatus{Goto: model.GotoMid, Param: strconv.FormatInt(op.ID, 10), IsAtten: c.IsAttenm[op.ID]}
	c.DescButton = buttonFrom(button, op.Plat)
	if op.Desc != "" {
		c.TopRcmdReasonStyle = reasonStyleFrom(model.BgColorOrange, op.Desc)
	}
	switch main.(type) {
	case map[int64]*archive.ArchiveWithPlayer:
		for _, item := range op.Items {
			am := main.(map[int64]*archive.ArchiveWithPlayer)
			var (
				a  *archive.ArchiveWithPlayer
				ok bool
			)
			if a, ok = am[item.ID]; !ok {
				continue
			}
			args := Args{}
			args.fromArchive(a.Archive3, c.Tagm[op.Tid])
			c.Items = append(c.Items, &TwoItemHV1Item{
				Title:          a.Title,
				Cover:          a.Pic,
				URI:            model.FillURI(model.GotoAv, strconv.FormatInt(item.ID, 10), model.AvPlayHandler(a.Archive3, nil, "")),
				Goto:           string(model.GotoAv),
				Param:          strconv.FormatInt(item.ID, 10),
				CoverLeftText1: model.StatString(a.Stat.View, ""),
				CoverLeftIcon1: model.IconPlay,
				CoverRightText: model.DurationString(a.Duration),
				Args:           args,
			})
		}
		if len(c.Items) < 3 {
			return
		}
	}
	c.Right = true
}

func (c *ThreeItemAllV2) Get() *Base {
	return c.Base
}

type MiddleCoverV3 struct {
	*Base
	Desc1      string       `json:"desc1,omitempty"`
	Desc2      string       `json:"desc2,omitempty"`
	CoverBadge *ReasonStyle `json:"cover_badge_style,omitempty"`
}

func (c *MiddleCoverV3) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	c.Base.from(op.Param, op.Cover, op.Title, model.GotoWeb, op.URI, nil)
	c.Goto = op.Goto
	if op.Badge != "" {
		c.CoverBadge = reasonStyleFrom(model.BgColorPurple, op.Badge)
	}
	c.Desc1 = op.Desc
	c.Right = true
}

func (c *MiddleCoverV3) Get() *Base {
	return c.Base
}

type Select struct {
	*Base
	Desc        string  `json:"desc,omitempty"`
	LeftButton  *Button `json:"left_button,omitempty"`
	RightButton *Button `json:"right_button,omitempty"`
}

func (c *Select) From(main interface{}, op *operate.Card) {
	if op == nil {
		return
	}
	switch main.(type) {
	case nil:
		switch op.CardGoto {
		case model.CardGotoFollowMode:
			if len(op.Buttons) < 2 {
				return
			}
			c.Base.from(op.Param, "", op.Title, "", "", nil)
			c.Desc = op.Desc
			c.LeftButton = buttonFrom(&ButtonStatus{Text: op.Buttons[0].Text, Event: model.Event(op.Buttons[0].Event)}, op.Plat)
			c.RightButton = buttonFrom(&ButtonStatus{Text: op.Buttons[1].Text, Event: model.Event(op.Buttons[1].Event)}, op.Plat)
		default:
			log.Warn("Select From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("Select From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *Select) Get() *Base {
	return c.Base
}
