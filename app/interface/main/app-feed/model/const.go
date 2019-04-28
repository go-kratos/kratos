package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/service/main/archive/model/archive"
)

const (
	// PlatAndroid is int8 for android.
	PlatAndroid = int8(0)
	// PlatIPhone is int8 for iphone.
	PlatIPhone = int8(1)
	// PlatIPad is int8 for ipad.
	PlatIPad = int8(2)
	// PlatWPhone is int8 for wphone.
	PlatWPhone = int8(3)
	// PlatAndroidG is int8 for Android Global.
	PlatAndroidG = int8(4)
	// PlatIPhoneI is int8 for Iphone Global.
	PlatIPhoneI = int8(5)
	// PlatIPadI is int8 for IPAD Global.
	PlatIPadI = int8(6)
	// PlatAndroidTV is int8 for AndroidTV Global.
	PlatAndroidTV = int8(7)
	// PlatAndroidI is int8 for Android Global.
	PlatAndroidI = int8(8)
	// PlatAndroidB is int8 for android_b
	PlatAndroidB = int8(10)
	// PlatIPhoneB is int8 for iphone_b
	PlatIPhoneB = int8(11)

	GotoAv              = "av"
	GotoWeb             = "web"
	GotoBangumi         = "bangumi"
	GotoPGC             = "pgc"
	GotoLive            = "live"
	GotoGame            = "game"
	GotoAdAv            = "ad_av"
	GotoAdWeb           = "ad_web"
	GotoRank            = "rank"
	GotoBangumiRcmd     = "bangumi_rcmd"
	GotoLogin           = "login"
	GotoUpBangumi       = "up_bangumi"
	GotoBanner          = "banner"
	GotoAdWebS          = "ad_web_s"
	GotoUpArticle       = "up_article"
	GotoConverge        = "converge"
	GotoSpecial         = "special"
	GotoArticle         = "article"
	GotoArticleS        = "article_s"
	GotoGameDownloadS   = "game_download_s"
	GotoShoppingS       = "shopping_s"
	GotoAudio           = "audio"
	GotoPlayer          = "player"
	GotoAdLarge         = "ad_large"
	GotoSpecialS        = "special_s"
	GotoPlayerLive      = "player_live"
	GotoSong            = "song"
	GotoLiveUpRcmd      = "live_up_rcmd"
	GotoUpRcmdAv        = "up_rcmd_av"
	GotoSubscribe       = "subscribe"
	GotoSearchSubscribe = "search_subscribe"
	GotoChannelRcmd     = "channel_rcmd"
	GotoMoe             = "moe"
	GotoPicture         = "picture"
	GotoInterest        = "interest"
	GotoFollowMode      = "follow_mode"

	// for fill uri
	GotoAudioTag = "audio_tag"
	GotoAlbum    = "album"
	GotoClip     = "clip"
	GotoDaily    = "daily"

	// extra tab
	GotoTabBackground  = "background"
	GotoTabEntrance    = "entrance"
	GotoTabContentRcmd = "content_rcmd"
	GotoTabTagRcmd     = "tag_rcmd"
	GotoTabSignIn      = "sign_in"
	GotoTabNews        = "news"
)

var (
	OperateType = map[int]string{
		0:  GotoWeb,
		1:  GotoGame,
		2:  GotoAv,
		3:  GotoBangumi,
		4:  GotoLive,
		6:  GotoArticleS,
		7:  GotoDaily,
		8:  GotoAudio,
		9:  GotoSong,
		10: GotoAlbum,
		11: GotoClip,
	}

	AudioHandler = func(uri string) string {
		return uri + "?from=tianma"
	}
	AvPlayHandler = func(a *archive.Archive3, ap *archive.PlayerInfo) func(uri string) string {
		var player string
		if ap != nil {
			bs, _ := json.Marshal(ap)
			player = url.QueryEscape(string(bs))
			if strings.IndexByte(player, '+') > -1 {
				player = strings.Replace(player, "+", "%20", -1)
			}
		}
		return func(uri string) string {
			if player != "" && (a.Dimension.Height != 0 || a.Dimension.Width != 0) {
				return fmt.Sprintf("%s?page=1&player_preload=%s&player_width=%d&player_height=%d&player_rotate=%d", uri, player, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			} else if player != "" {
				return fmt.Sprintf("%s?page=1&player_preload=%s", uri, player)
			} else if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				return fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			return uri
		}
	}
	AvGameHandler = func(uri string) string {
		return uri + "&sourceType=adPut"
	}
	LiveUpHandler = func(l *live.Card) func(uri string) string {
		return func(uri string) string {
			if l == nil {
				return uri
			}
			return fmt.Sprintf("%s?broadcast_type=%d", uri, l.BroadcastType)
		}
	}
	LiveRoomHandler = func(l *live.Room) func(uri string) string {
		return func(uri string) string {
			if l == nil {
				return uri
			}
			return fmt.Sprintf("%s?broadcast_type=%d", uri, l.BroadcastType)
		}
	}
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI || plat == PlatIPhoneB
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// IsIPhoneB check plat is ios but not iphone_b.
func IsIOSNormal(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// FillURI deal app schema.
func FillURI(gt, param string, plat int8, build int, f func(uri string) string) (uri string) {
	if param == "" {
		return
	}
	switch gt {
	case GotoAv, GotoAdAv, GotoChannelRcmd, "":
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "https://www.bilibili.com/bangumi/play/ep" + param
	case GotoUpBangumi:
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoUpArticle, GotoArticle, GotoArticleS:
		uri = "bilibili://article/" + param
	case GotoGame:
		const (
			_iPhoneGameCenter  = 6500
			_androidGameCenter = 519010
		)
		// TODO FUCK GAME
		if (plat == PlatAndroid && build >= _androidGameCenter) || (plat == PlatIPhone && build >= _iPhoneGameCenter) || plat == PlatIPhoneB {
			uri = "bilibili://game_center/detail?id=" + param + "&sourceType=adPut"
		} else {
			uri = "bilibili://game/" + param
		}
	case GotoAudio:
		uri = "bilibili://music/menu/detail/" + param
	case GotoSong:
		uri = "bilibili://music/detail/" + param
	case GotoAudioTag:
		uri = "bilibili://music/categorydetail/" + param
	case GotoDaily:
		uri = "bilibili://pegasus/list/daily/" + param
	case GotoAlbum:
		uri = "bilibili://album/" + param
	case GotoClip:
		uri = "bilibili://clip/" + param
	case GotoWeb, GotoAdWeb, GotoRank, GotoAdWebS, GotoShoppingS, GotoAdLarge:
		uri = param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

func FillRedirect(gt string, typ int) (redirect string) {
	switch gt {
	case GotoSpecial, GotoConverge, GotoGameDownloadS, GotoSpecialS:
		switch typ {
		case 7:
			redirect = "daily"
		case 6:
			redirect = "article"
		case 5:
			redirect = "category/65541"
		case 4:
			redirect = "live"
		case 3:
			redirect = ""
		case 2:
			redirect = "video"
		case 1:
			redirect = "game"
		case 0:
			redirect = ""
		}
	}
	return
}

// CoverURL convert cover url to full url.
func CoverURL(uri string) (cover string) {
	if uri == "" {
		cover = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	if strings.HasPrefix(uri, "http://i0.hdslb.com") || strings.HasPrefix(uri, "http://i1.hdslb.com") || strings.HasPrefix(uri, "http://i2.hdslb.com") {
		uri = uri[19:]
	} else if strings.HasPrefix(uri, "https://i0.hdslb.com") || strings.HasPrefix(uri, "https://i1.hdslb.com") || strings.HasPrefix(uri, "https://i2.hdslb.com") {
		uri = uri[20:]
	}
	cover = uri
	if strings.HasPrefix(uri, "/bfs") {
		cover = "http://i0.hdslb.com" + uri
		return
	}
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") || strings.HasPrefix(uri, "/group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i0.hdslb.com" + cover
	return
}

func CoverURLHTTPS(uri string) (cover string) {
	if strings.HasPrefix(uri, "http://") {
		cover = "https://" + uri[7:]
	} else {
		cover = uri
	}
	return
}

// InvalidBuild check source build is not allow by config build and condition.
// eg: when condition is gt, means srcBuild must gt cfgBuild, otherwise is invalid srcBuild.
func InvalidBuild(srcBuild, cfgBuild int, cfgCond string) bool {
	if cfgBuild != 0 && cfgCond != "" {
		switch cfgCond {
		case "gt":
			if cfgBuild >= srcBuild {
				return true
			}
		case "lt":
			if cfgBuild <= srcBuild {
				return true
			}
		case "eq":
			if cfgBuild != srcBuild {
				return true
			}
		case "ne":
			if cfgBuild == srcBuild {
				return true
			}
		}
	}
	return false
}

// InvalidChannel check source channel is not allow by config channel.
func InvalidChannel(plat int8, srcCh, cfgCh string) bool {
	return plat == PlatAndroid && cfgCh != "*" && cfgCh != srcCh
}

// Plat return plat by platStr or mobiApp
func Plat(mobiApp, device string) int8 {
	switch mobiApp {
	case "iphone":
		if device == "pad" {
			return PlatIPad
		}
		return PlatIPhone
	case "white":
		return PlatIPhone
	case "ipad":
		return PlatIPad
	case "android", "android_b":
		return PlatAndroid
	case "win":
		return PlatWPhone
	case "android_G":
		return PlatAndroidG
	case "android_i":
		return PlatAndroidI
	case "iphone_i":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_i":
		return PlatIPadI
	case "android_tv":
		return PlatAndroidTV
	case "iphone_b":
		return PlatIPhoneB
	}
	return PlatIPhone
}

type SortInt64 []int64

func (is SortInt64) Len() int           { return len(is) }
func (is SortInt64) Less(i, j int) bool { return is[i] > is[j] }
func (is SortInt64) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

// AdAvIsNormal check advert archive normal.
func AdAvIsNormal(a *archive.ArchiveWithPlayer) bool {
	if a == nil || a.Archive3 == nil {
		return false
	}
	return a.State >= 0 || a.State == -6 || a.State == -40
}

func Rounding(number, divisor int64) string {
	if divisor > 0 {
		tmp := float64(number) / float64(divisor)
		tmpStr := fmt.Sprintf("%0.1f", tmp)
		parts := strings.Split(tmpStr, ".")
		if len(parts) > 1 && parts[1] == "0" {
			return parts[0]
		}
		return tmpStr
	}
	return strconv.FormatInt(number, 10)
}

func PlatAPPBuleChange(plat int8) int8 {
	switch plat {
	case PlatAndroidB:
		return PlatAndroid
	case PlatIPhoneB:
		return PlatIPhone
	}
	return plat
}
