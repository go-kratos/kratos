package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/live"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// CardGt is
type CardGt string

// CardType is
type CardType string

// ColumnStatus is
type ColumnStatus int8

// Gt is
type Gt string

// Icon is
type Icon int8

// Type is
type Type int8

// BlurStatus is
type BlurStatus int8

// Event is
type Event string

// CoverColor is
type CoverColor string

// Switch is
type Switch string

const (
	// PlatAndroid is int8 for android.
	PlatAndroid = int8(0)
	// PlatIPhone is int8 for iphone.
	PlatIPhone = int8(1)
	// PlatIPad is int8 for ipad.
	PlatIPad = int8(2)
	// PlatWPhone is int8 for wphone.
	PlatWPhone = int8(3)
	// PlatAndroidG is int8 for Android Googleplay.
	PlatAndroidG = int8(4)
	// PlatIPhoneI is int8 for Iphone Global.
	PlatIPhoneI = int8(5)
	// PlatIPadI is int8 for IPAD Global.
	PlatIPadI = int8(6)
	// PlatAndroidTV is int8 for AndroidTV Global.
	PlatAndroidTV = int8(7)
	// PlatAndroidI is int8 for Android Global.
	PlatAndroidI = int8(8)
	// PlatIpadHD is int8 for IpadHD
	PlatIpadHD = int8(9)
	// PlatAndroidTVYST is int8 for AndroidTV_YST Global.
	PlatAndroidTVYST = int8(12)

	CardGotoAv              = CardGt("av")
	CardGotoBangumi         = CardGt("bangumi")
	CardGotoLive            = CardGt("live")
	CardGotoArticle         = CardGt("article")
	CardGotoAudio           = CardGt("audio")
	CardGotoRank            = CardGt("rank")
	CardGotoConverge        = CardGt("converge")
	CardGotoDownload        = CardGt("download")
	CardGotoMoe             = CardGt("moe")
	CardGotoBanner          = CardGt("banner")
	CardGotoLogin           = CardGt("login")
	CardGotoPGC             = CardGt("pgc")
	CardGotoSpecial         = CardGt("special")
	CardGotoSubscribe       = CardGt("subscribe")
	CardGotoBangumiRcmd     = CardGt("bangumi_rcmd")
	CardGotoUpRcmdAv        = CardGt("up_rcmd_av")
	CardGotoChannelRcmd     = CardGt("channel_rcmd")
	CardGotoLiveUpRcmd      = CardGt("live_up_rcmd")
	CardGotoAdAv            = CardGt("ad_av")
	CardGotoAdWeb           = CardGt("ad_web")
	CardGotoAdWebS          = CardGt("ad_web_s")
	CardGotoPlayer          = CardGt("player")
	CardGotoPlayerLive      = CardGt("player_live")
	CardGotoArticleS        = CardGt("article_s")
	CardGotoSpecialS        = CardGt("special_s")
	CardGotoShoppingS       = CardGt("shopping_s")
	CardGotoGameDownloadS   = CardGt("game_download_s")
	CardGotoHotTopic        = CardGt("hottopic")
	CardGotoTopstick        = CardGt("topstick")
	CardGotoSearchSubscribe = CardGt("search_subscribe")
	CardGotoPicture         = CardGt("picture")
	CardGotoInterest        = CardGt("interest")
	CardGotoFollowMode      = CardGt("follow_mode")
	// operate tab
	CardGotoEntrance      = CardGt("entrance")
	CardGotoContentRcmd   = CardGt("content_rcmd")
	CardGotoTagRcmd       = CardGt("tag_rcmd")
	CardGotoNews          = CardGt("news")
	CardGotoChannelSquare = CardGt("channel_square")
	CardGotoPgcsRcmd      = CardGt("pgcs_rcmd")
	CardGotoUpRcmdS       = CardGt("up_rcmd_s")
	CardGotoSearchUpper   = CardGt("search_upper")
	CardGotoUpRcmdNew     = CardGt("up_rcmd_new")
	CardGotoDynamicHot    = CardGt("hot_dynamic")
	CardGotoUpRcmdNewV2   = CardGt("up_rcmd_new_v2")
	CardGotoEventTopic    = CardGt("event_topic")
	// single card
	LargeCoverV1   = CardType("large_cover_v1")
	SmallCoverV1   = CardType("small_cover_v1")
	MiddleCoverV1  = CardType("middle_cover_v1")
	ThreeItemV1    = CardType("three_item_v1")
	ThreeItemHV1   = CardType("three_item_h_v1")
	ThreeItemHV3   = CardType("three_item_h_v3")
	TwoItemV1      = CardType("two_item_v1")
	CoverOnlyV1    = CardType("cover_only_v1")
	BannerV1       = CardType("banner_v1")
	CmV1           = CardType("cm_v1")
	HotTopic       = CardType("hot_topic")
	TopStick       = CardType("top_stick")
	ChannelSquare  = CardType("channel_square")
	ThreeItemHV4   = CardType("three_item_h_v4")
	UpRcmdCover    = CardType("up_rcmd_cover")
	ThreeItemAll   = CardType("three_item_all")
	TwoItemHV1     = CardType("two_item_h_v1")
	OnePicV1       = CardType("one_pic_v1")
	ThreePicV1     = CardType("three_pic_v1")
	SmallCoverV5   = CardType("small_cover_v5")
	OptionsV1      = CardType("options_v1")
	HotDynamic     = CardType("hot_dynamic")
	ThreeItemAllV2 = CardType("three_item_all_v2")
	ThreeItemHV5   = CardType("three_item_h_v5")
	MiddleCoverV3  = CardType("middle_cover_v3")
	Select         = CardType("select")
	// double card
	SmallCoverV2  = CardType("small_cover_v2")
	SmallCoverV3  = CardType("small_cover_v3")
	MiddleCoverV2 = CardType("middle_cover_v2")
	LargeCoverV2  = CardType("large_cover_v2")
	ThreeItemHV2  = CardType("three_item_h_v2")
	ThreeItemV2   = CardType("three_item_v2")
	TwoItemV2     = CardType("two_item_v2")
	SmallCoverV4  = CardType("small_cover_v4")
	CoverOnlyV2   = CardType("cover_only_v2")
	BannerV2      = CardType("banner_v2")
	CmV2          = CardType("cm_v2")
	News          = CardType("news")
	MultiItem     = CardType("multi_item")
	MultiItemH    = CardType("multi_item_h")
	ThreePicV2    = CardType("three_pic_v2")
	OptionsV2     = CardType("options_v2")
	OnePicV2      = CardType("one_pic_v2")
	// ipad card
	BannerV3    = CardType("banner_v3")
	CoverOnlyV3 = CardType("cover_only_v3")
	FourItemHV3 = CardType("four_item_h_v3")

	ColumnDefault    = ColumnStatus(0)
	ColumnSvrSingle  = ColumnStatus(1)
	ColumnSvrDouble  = ColumnStatus(2)
	ColumnUserSingle = ColumnStatus(3)
	ColumnUserDouble = ColumnStatus(4)

	GotoWeb        = Gt("web")
	GotoAv         = Gt("av")
	GotoBangumi    = Gt("bangumi")
	GotoLive       = Gt("live")
	GotoGame       = Gt("game")
	GotoArticle    = Gt("article")
	GotoArticleTag = Gt("article_tag")
	GotoAudio      = Gt("audio")
	GotoAudioTag   = Gt("audio_tag")
	GotoSong       = Gt("song")
	GotoAlbum      = Gt("album")
	GotoClip       = Gt("clip")
	GotoDaily      = Gt("daily")
	GotoTag        = Gt("tag")
	GotoMid        = Gt("mid")
	GotoDynamicMid = Gt("dynamic_mid")
	GotoConverge   = Gt("converge")
	GotoRank       = Gt("rank")
	GotoLiveTag    = Gt("live_tag")
	GotoPGC        = Gt("pgc")
	GotoHotTopic   = Gt("hottopic")
	GotoTopstick   = Gt("topstick")
	GotoSpecial    = Gt("special")
	GotoSubscribe  = Gt("subscribe")
	GotoPicture    = Gt("picture")
	GotoPictureTag = Gt("picture_tag")
	GotoHotDynamic = Gt("hot_dynamic")

	IconPlay           = Icon(1)
	IconOnline         = Icon(2)
	IconDanmaku        = Icon(3)
	IconFavorite       = Icon(4)
	IconStar           = Icon(5)
	IconRead           = Icon(6)
	IconComment        = Icon(7)
	IconLocation       = Icon(8)
	IconHeadphone      = Icon(9)
	IconRank           = Icon(10)
	IconGoldMedal      = Icon(11)
	IconSilverMedal    = Icon(12)
	IconBronzeMedal    = Icon(13)
	IconTV             = Icon(14)
	IconBomb           = Icon(15)
	IconRoleYellow     = Icon(16)
	IconRoleBlue       = Icon(17)
	IconRoleVipRed     = Icon(18)
	IconRoleYearVipRed = Icon(19)
	IconLike           = Icon(20)

	AvatarRound  = Type(0)
	AvatarSquare = Type(1)

	ButtonGrey  = Type(1)
	ButtonTheme = Type(2)

	BlurNo  = BlurStatus(0)
	BlurYes = BlurStatus(1)

	EventUpFollow         = Event("up_follow")
	EventChannelSubscribe = Event("channel_subscribe")
	EventUpClick          = Event("up_click")
	EventChannelClick     = Event("channel_click")
	EventButtonClick      = Event("button_click")
	EventGameClick        = Event("game_click")

	PurpleCoverBadge = CoverColor("purple")

	BgColorOrange            = int8(0)
	BgColorTransparentOrange = int8(1)
	BgColorBlue              = int8(2)
	BgColorRed               = int8(3)
	BgTransparentTextOrange  = int8(4)
	BgColorPurple            = int8(5)

	BgStyleFill              = int8(1)
	BgStyleStroke            = int8(2)
	BgStyleFillAndStroke     = int8(3)
	BgStyleNoFillAndNoStroke = int8(4)

	SwitchFeedIndexLike          = Switch("天马卡片好评数替换弹幕数")
	SwitchFeedIndexTabThreePoint = Switch("运营tab稿件卡片三点稍后再看")
	SwitchCooperationHide        = Switch("cooperation_hide")
	SwitchCooperationShow        = Switch("cooperation_show")

	// 热门显示up主信息abtest
	HotCardStyleOld    = int8(0)
	HotCardStyleShowUp = int8(1)
	HotCardStyleHideUp = int8(2)
)

var (
	OperateType = map[int]Gt{
		0:  GotoWeb,
		1:  GotoGame,
		2:  GotoAv,
		3:  GotoBangumi,
		4:  GotoLive,
		6:  GotoArticle,
		7:  GotoDaily,
		8:  GotoAudio,
		9:  GotoSong,
		10: GotoAlbum,
		11: GotoClip,
		12: GotoSpecial,
		13: GotoPicture,
	}

	Columnm = map[ColumnStatus]ColumnStatus{
		ColumnDefault:    ColumnSvrDouble,
		ColumnSvrSingle:  ColumnSvrSingle,
		ColumnSvrDouble:  ColumnSvrDouble,
		ColumnUserSingle: ColumnSvrSingle,
		ColumnUserDouble: ColumnSvrDouble,
	}

	AvatarEvent = map[Gt]Event{
		GotoMid:        EventUpClick,
		GotoTag:        EventChannelClick,
		GotoDynamicMid: EventUpClick,
	}

	ButtonEvent = map[Gt]Event{
		GotoMid: EventUpFollow,
		GotoTag: EventChannelSubscribe,
	}

	ButtonText = map[Gt]string{
		GotoMid: "+ 关注",
		GotoTag: "订阅",
	}

	LiveRoomTagHandler = func(r *live.Room) func(uri string) string {
		return func(uri string) string {
			if r == nil {
				return ""
			}
			return fmt.Sprintf("%s?parent_area_id=%d&parent_area_name=%s&area_id=%d&area_name=%s", uri, r.AreaV2ParentID, url.QueryEscape(r.AreaV2ParentName), r.AreaV2ID, url.QueryEscape(r.AreaV2Name))
		}
	}
	AudioTagHandler = func(c []*audio.Ctg) func(uri string) string {
		return func(uri string) string {
			var schema string
			if len(c) != 0 {
				schema = c[0].Schema
				if len(c) > 1 {
					schema = c[1].Schema
				}
			}
			return schema
		}
	}
	LiveUpHandler = func(card *live.Card) func(uri string) string {
		return func(uri string) string {
			if card == nil {
				return uri
			}
			return fmt.Sprintf("%s?broadcast_type=%d", uri, card.BroadcastType)
		}
	}
	LiveRoomHandler = func(r *live.Room) func(uri string) string {
		return func(uri string) string {
			if r == nil {
				return uri
			}
			return fmt.Sprintf("%s?broadcast_type=%d", uri, r.BroadcastType)
		}
	}
	AvPlayHandler = func(a *archive.Archive3, ap *archive.PlayerInfo, trackID string) func(uri string) string {
		var player string
		if ap != nil {
			bs, _ := json.Marshal(ap)
			player = url.QueryEscape(string(bs))
			if strings.IndexByte(player, '+') > -1 {
				player = strings.Replace(player, "+", "%20", -1)
			}
		}
		return func(uri string) string {
			var uriStr string
			if player != "" && (a.Dimension.Height != 0 || a.Dimension.Width != 0) {
				uriStr = fmt.Sprintf("%s?page=1&player_preload=%s&player_width=%d&player_height=%d&player_rotate=%d", uri, player, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			} else if player != "" {
				uriStr = fmt.Sprintf("%s?page=1&player_preload=%s", uri, player)
			} else if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				uriStr = fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			if trackID != "" {
				if uriStr == "" {
					uriStr = fmt.Sprintf("%s?trackid=%s", uri, trackID)
				} else {
					uriStr = fmt.Sprintf("%s&trackid=%s", uriStr, trackID)
				}
			}
			if uriStr != "" {
				return uriStr
			}
			return uri
		}
	}

	HottopicHandler = func(l *live.TopicHot) func(uri string) string {
		return func(uri string) string {
			return fmt.Sprintf("%s?type=topic", uri)
		}
	}

	ArticleTagHandler = func(c []*article.Category, plat int8) func(uri string) string {
		return func(uri string) string {
			var (
				rid int64
				tid int64
			)
			if len(c) > 1 {
				if c[0] != nil {
					rid = c[0].ID
				}
				if c[1] != nil {
					tid = c[1].ID
				}
			}
			if rid != 0 && tid != 0 {
				return fmt.Sprintf("bilibili://article/category/%d?sec_cid=%d", rid, tid)
			}
			return ""
		}
	}
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidI
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// IsIPhone check plat is iphone.
func IsIPhone(plat int8) bool {
	return plat == PlatIPhone || plat == PlatIPhoneI
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI || plat == PlatIpadHD
}

// AdAvIsNormal check advert archive normal.
func AdAvIsNormal(a *archive.ArchiveWithPlayer) bool {
	if a == nil || a.Archive3 == nil {
		return false
	}
	return a.State >= 0 || a.State == -6 || a.State == -40
}

func AvIsNormal(a *archive.ArchiveWithPlayer) bool {
	if a == nil || a.Archive3 == nil {
		return false
	}
	return a.IsNormal()
}

// FillURI deal app schema.
func FillURI(gt Gt, param string, f func(uri string) string) (uri string) {
	switch gt {
	case GotoAv:
		if param != "" {
			uri = "bilibili://video/" + param
		}
	case GotoLive:
		if param != "" {
			uri = "bilibili://live/" + param
		}
	case GotoBangumi:
		if param != "" {
			uri = "https://www.bilibili.com/bangumi/play/ep" + param
		}
	case GotoPGC:
		if param != "" {
			uri = "https://www.bilibili.com/bangumi/play/ss" + param
		}
	case GotoArticle:
		if param != "" {
			uri = "bilibili://article/" + param
		}
	case GotoArticleTag:
		// TODO fuck article
	case GotoGame:
		// TODO fuck game
		if param != "" {
			uri = "bilibili://game_center/detail?id=" + param + "&sourceType=adPut"
		}
	case GotoAudio:
		if param != "" {
			uri = "bilibili://music/menu/detail/" + param
		}
	case GotoSong:
		if param != "" {
			uri = "bilibili://music/detail/" + param
		}
	case GotoAudioTag:
		// uri = "bilibili://music/menus/menu?itemId=(请求所需参数)&cateId=(请求所需参数)&itemVal=(分类的标题value)"
	case GotoDaily:
		if param != "" {
			uri = "bilibili://pegasus/list/daily/" + param
		}
	case GotoAlbum:
		if param != "" {
			uri = "bilibili://album/" + param
		}
	case GotoClip:
		if param != "" {
			uri = "bilibili://clip/" + param
		}
	case GotoTag:
		if param != "" {
			uri = "bilibili://pegasus/channel/" + param
		}
	case GotoMid:
		if param != "" {
			uri = "bilibili://space/" + param
		}
	case GotoDynamicMid:
		if param != "" {
			uri = "bilibili://space/" + param + "?defaultTab=dynamic"
		}
	case GotoRank:
		uri = "bilibili://rank/"
	case GotoConverge:
		if param != "" {
			uri = "bilibili://pegasus/converge/" + param
		}
	case GotoLiveTag:
		uri = "https://live.bilibili.com/app/area"
	case GotoHotTopic:
		uri = "bilibili://pegasus/channel/" + param
	case GotoWeb:
		uri = param
	case GotoPicture:
		uri = "bilibili://following/detail/" + param
	case GotoPictureTag:
		uri = "bilibili://pegasus/channel/0/?name=" + param + "&type=topic"
	case GotoHotDynamic:
		uri = "bilibili://following/detail/" + param
	default:
		uri = param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

// DurationString duration to string
func DurationString(second int64) (s string) {
	var hour, min, sec int
	if second < 1 {
		return
	}
	d, err := time.ParseDuration(strconv.FormatInt(second, 10) + "s")
	if err != nil {
		log.Error("%+v", err)
		return
	}
	r := strings.NewReplacer("h", ":", "m", ":", "s", ":")
	ts := strings.Split(strings.TrimSuffix(r.Replace(d.String()), ":"), ":")
	if len(ts) == 1 {
		sec, _ = strconv.Atoi(ts[0])
	} else if len(ts) == 2 {
		min, _ = strconv.Atoi(ts[0])
		sec, _ = strconv.Atoi(ts[1])
	} else if len(ts) == 3 {
		hour, _ = strconv.Atoi(ts[0])
		min, _ = strconv.Atoi(ts[1])
		sec, _ = strconv.Atoi(ts[2])
	}
	if hour == 0 {
		s = fmt.Sprintf("%d:%02d", min, sec)
		return
	}
	s = fmt.Sprintf("%d:%02d:%02d", hour, min, sec)
	return
}

// StatString Stat to string
func StatString(number int32, suffix string) (s string) {
	if number == 0 {
		s = "-" + suffix
		return
	}
	if number < 10000 {
		s = strconv.FormatInt(int64(number), 10) + suffix
		return
	}
	if number < 100000000 {
		s = strconv.FormatFloat(float64(number)/10000, 'f', 1, 64)
		return strings.TrimSuffix(s, ".0") + "万" + suffix
	}
	s = strconv.FormatFloat(float64(number)/100000000, 'f', 1, 64)
	return strings.TrimSuffix(s, ".0") + "亿" + suffix
}

// ArchiveViewString ArchiveView to string
func ArchiveViewString(number int32) string {
	const _suffix = "观看"
	return StatString(number, _suffix)
}

// DanmakuString Danmaku to string
func DanmakuString(number int32) string {
	const _suffix = "弹幕"
	return StatString(number, _suffix)
}

// LikeString Danmaku to string
func LikeString(number int32) string {
	const _suffix = "点赞"
	return StatString(number, _suffix)
}

// BangumiFavString BangumiFav to string
func BangumiFavString(number int32) string {
	const _suffix = "追番"
	return StatString(number, _suffix)
}

// LiveOnlineString online to string
func LiveOnlineString(number int32) string {
	const _suffix = "人气"
	return StatString(number, _suffix)
}

// FanString fan to string
func FanString(number int32) string {
	const _suffix = "粉丝"
	return StatString(number, _suffix)
}

// AttentionString fan to string
func AttentionString(number int32) string {
	const _suffix = "人关注"
	return StatString(number, _suffix)
}

// AudioDescString audio to string
func AudioDescString(firstSong string, total int) (desc1, desc2 string) {
	desc1 = firstSong
	if total == 1 {
		desc2 = "共1首歌曲"
		return
	}
	desc2 = "...共" + strconv.Itoa(total) + "首歌曲"
	return
}

// AudioTotalStirng audioTotal to string
func AudioTotalStirng(total int) string {
	if total == 0 {
		return ""
	}
	return strconv.Itoa(total) + "首歌曲"
}

// AudioBadgeString audioBadge to string
func AudioBadgeString(number int8) string {
	if number == 5 {
		return "专辑"
	}
	return "歌单"
}

// AudioPlayString audioPlay to string
func AudioPlayString(number int32) string {
	const _suffix = "收听"
	return StatString(number, _suffix)
}

// AudioFavString audioFav to string
func AudioFavString(numbber int32) string {
	const _suffix = "收藏"
	return StatString(numbber, _suffix)
}

// DownloadString download to string
func DownloadString(number int32) string {
	if number == 0 {
		return ""
	}
	const _suffix = "下载"
	return StatString(number, _suffix)
}

// ArticleViewString articleView to string
func ArticleViewString(number int64) string {
	const _suffix = "阅读"
	return StatString(int32(number), _suffix)
}

// PictureViewString pictureView to string
func PictureViewString(number int64) string {
	const _suffix = "浏览"
	return StatString(int32(number), _suffix)
}

// ArticleReplyString articleReply to string
func ArticleReplyString(number int64) string {
	const _suffix = "评论"
	return StatString(int32(number), _suffix)
}

// SubscribeString subscribe to string
func SubscribeString(number int32) string {
	const _suffix = "人已订阅"
	return StatString(number, _suffix)
}

// RecommendString recommend to string
func RecommendString(like, dislike int32) string {
	rcmd := like / (like + dislike) * 100
	if rcmd != 0 {
		return strconv.Itoa(int(rcmd)) + "%的人推荐"
	}
	return ""
}

// ShoppingDuration shopping duration
func ShoppingDuration(stime, etime string) string {
	if stime == "" && etime == "" {
		return ""
	}
	return stime + " - " + etime
}

// ScoreString is
func ScoreString(number int32) string {
	const _prefix = "综合评分："
	score := StatString(number, "")
	if score != "" {
		return _prefix + score
	}
	return _prefix + "-"
}

// ShoppingCover is
func ShoppingCover(cover string) string {
	if strings.HasPrefix(cover, "http:") || strings.HasPrefix(cover, "https:") {
		return cover
	}
	return "http:" + cover
}

// BangumiIcon is.
func BangumiIcon(typ int8) (icon Icon) {
	switch typ {
	case 1, 4:
		icon = IconFavorite
	case 2, 3, 5:
		icon = IconStar
	}
	return icon
}

// PubDataString is.
func PubDataString(t time.Time) (s string) {
	if t.IsZero() {
		return
	}
	now := time.Now()
	sub := now.Sub(t)
	if sub < time.Minute {
		s = "刚刚"
		return
	}
	if sub < time.Hour {
		s = strconv.FormatFloat(sub.Minutes(), 'f', 0, 64) + "分钟前"
		return
	}
	if sub < 24*time.Hour {
		s = strconv.FormatFloat(sub.Hours(), 'f', 0, 64) + "小时前"
		return
	}
	if now.Year() == t.Year() {
		if now.YearDay()-t.YearDay() == 1 {
			s = "昨天"
			return
		}
		s = t.Format("01-02")
		return
	}
	s = t.Format("2006-01-02")
	return
}

// PictureCountString is.
func PictureCountString(count int) string {
	return strconv.Itoa(count) + "P"
}

// OfficialIcon return 认证图标（1 UP 主认证，2 身份认证）黄标，（3 企业认证，4 政府认证，5 媒体认证，6 其他认证）蓝标
func OfficialIcon(cd *account.Card) (icon Icon) {
	if cd == nil {
		return
	}
	switch cd.Official.Role {
	case 1, 2:
		icon = IconRoleYellow
	case 3, 4, 5, 6:
		icon = IconRoleBlue
	}
	return
}
