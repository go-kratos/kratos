package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"go-common/app/service/main/archive/api"
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
	// PlatAndroidB is int8 for android_b
	PlatAndroidB = int8(10)
	// PlatIPhoneB is int8 for iphone_b
	PlatIPhoneB = int8(11)
	// PlatAndroidTVYST is int8 for AndroidTV_YST Global.
	PlatAndroidTVYST = int8(12)

	GotoAv          = "av"
	GotoWeb         = "web"
	GotoBangumi     = "bangumi"
	GotoLive        = "live"
	GotoGame        = "game"
	GotoArticle     = "article"
	GotoSpecial     = "special"
	GotoCm          = "cm"
	GotoSearchUpper = "search_upper"

	// for fill uri
	GotoAudio    = "audio"
	GotoSong     = "song"
	GotoAudioTag = "audio_tag"
	GotoAlbum    = "album"
	GotoClip     = "clip"
	GotoDaily    = "daily"

	// EnvPro is pro.
	EnvPro = "pro"
	EnvHK  = "hk"
	// EnvTest is env.
	EnvTest = "test"
	// EnvDev is env.
	EnvDev = "dev"
	// ForbidCode is forbid by law
	ForbidCode = -110

	StatusIng      = 0
	StatusPay      = 1
	StatusFree     = 2
	StatusVipFree  = 3
	StatusVipOnly  = 4
	StatusVipFirst = 5
	CoverIng       = "即将上映"
	CoverPay       = "付费观看"
	CoverFree      = "免费观看"
	CoverVipFree   = "付费观看"
	CoverVipOnly   = "专享"
	CoverVipFirst  = "抢先"

	Hans = "hans"
	Hant = "hant"

	FromOrder     = "order"
	FromOperation = "operation"
	FromRcmd      = "recommend"
)

var (
	OperateType = map[int]string{
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
	}

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
	LiveRoomHandler = func(broadcastType int) func(uri string) string {
		return func(uri string) string {
			return fmt.Sprintf("%s?broadcast_type=%d", uri, broadcastType)
		}
	}
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidI
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI || plat == PlatIPhoneB
}

// IsIPhone check plat is iphone.
func IsIPhone(plat int8) bool {
	return plat == PlatIPhone || plat == PlatIPhoneI || plat == PlatIPhoneB
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI || plat == PlatIpadHD
}

// IsIOSNormal check plat is ios except iphone_b
func IsIOSNormal(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// IsIPhoneB check plat is iphone_b
func IsIPhoneB(plat int8) bool {
	return plat == PlatIPhoneB
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
		return PlatIpadHD
	case "android", "android_b":
		return PlatAndroid
	case "win", "winphone":
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
	case "android_tv_yst":
		return PlatAndroidTVYST
	case "iphone_b":
		return PlatIPhoneB
	}
	return PlatIPhone
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
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	case GotoGame:
		uri = param
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
	case GotoWeb:
		uri = param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

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

// Platform plat to platform
func Platform(plat int8) string {
	if IsAndroid(plat) {
		return "android"
	} else {
		return "ios"
	}
}
