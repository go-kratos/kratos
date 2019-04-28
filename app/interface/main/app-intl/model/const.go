package model

import (
	"encoding/json"
	"fmt"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"net/url"
	"strings"
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
	GotoBanner          = "banner"
	GotoAdWebS          = "ad_web_s"
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

	// GotoAuthor is search
	GotoAuthor         = "author"
	GotoSp             = "sp"
	GotoMovie          = "movie"
	GotoEP             = "ep"
	GotoSuggestKeyWord = "suggest_keyword"
	GotoRecommendWord  = "recommend_word"
	GotoTwitter        = "twitter"
	GotoChannel        = "channel"

	FromOrder     = "order"
	FromOperation = "operation"
	FromRcmd      = "recommend"

	CoverIng      = "即将上映"
	CoverPay      = "付费观看"
	CoverFree     = "免费观看"
	CoverVipFree  = "付费观看"
	CoverVipOnly  = "专享"
	CoverVipFirst = "抢先"

	Hans = "hans"
	Hant = "hant"

	// ForbidCode is forbid by law
	ForbidCode   = -110
	NoResultCode = -111
)

var (
	// AvHandler is handler
	AvHandler = func(a *api.Arc, trackid string, ap *archive.PlayerInfo) func(uri string) string {
		var player string
		if ap != nil {
			bs, _ := json.Marshal(ap)
			player = url.QueryEscape(string(bs))
			if strings.IndexByte(player, '+') > -1 {
				player = strings.Replace(player, "+", "%20", -1)
			}
		}
		return func(uri string) string {
			if a == nil {
				return uri
			}
			var uriStr string
			if player != "" && (a.Dimension.Height != 0 || a.Dimension.Width != 0) {
				uriStr = fmt.Sprintf("%s?page=1&player_preload=%s&player_width=%d&player_height=%d&player_rotate=%d", uri, player, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			} else if player != "" {
				uriStr = fmt.Sprintf("%s?page=1&player_preload=%s", uri, player)
			} else if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				uriStr = fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			if trackid != "" {
				if uriStr == "" {
					uriStr = fmt.Sprintf("%s?trackid=%s", uri, trackid)
				} else {
					uriStr = fmt.Sprintf("%s&trackid=%s", uriStr, trackid)
				}
			}
			if uriStr != "" {
				return uriStr
			}
			return uri
		}
	}
)

// TWLocale is taiwan locale
func TWLocale(locale string) bool {
	var twLocalem = map[string]struct{}{
		"zh_hk": struct{}{},
		"zh_mo": struct{}{},
		"zh_tw": struct{}{},
	}
	_, ok := twLocalem[strings.ToLower(locale)]
	return ok
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
	}
	return PlatIPhone
}

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG
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
	return plat == PlatIPad || plat == PlatIPadI
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// FillURI deal app schema.
func FillURI(gt, param string, f func(uri string) string) (uri string) {
	if param == "" {
		return
	}
	switch gt {
	case GotoAv, "":
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "bilibili://bangumi/season/" + param
	case GotoGame:
		uri = "bilibili://game_center/detail?id=" + param + "&sourceType=adPut"
	case GotoSp:
		uri = "bilibili://splist/" + param
	case GotoAuthor:
		uri = "bilibili://author/" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	case GotoWeb:
		uri = param
	case GotoPGC:
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoChannel:
		uri = "bilibili://pegasus/channel/" + param + "/"
	case GotoEP:
		uri = "https://www.bilibili.com/bangumi/play/ep" + param
	case GotoTwitter:
		uri = "bilibili://pictureshow/detail/" + param
	default:
		return
	}
	if f != nil {
		uri = f(uri)
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

// StatusMark cover status mark
func StatusMark(status int) string {
	if status == 0 {
		return CoverIng
	} else if status == 1 {
		return CoverPay
	} else if status == 2 {
		return CoverFree
	} else if status == 3 {
		return CoverVipFree
	} else if status == 4 {
		return CoverVipOnly
	} else if status == 5 {
		return CoverVipFirst
	}
	return ""
}
