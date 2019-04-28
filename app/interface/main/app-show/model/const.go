package model

import (
	"fmt"
	"go-common/app/service/main/archive/api"
	"strconv"
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
	// PlatAndroidB is int8 for Android Bule.
	PlatAndroidB = int8(9)

	GotoAv          = "av"
	GotoWeb         = "web"
	GotoBangumi     = "bangumi"
	GotoBangumiWeb  = "bangumi_web"
	GotoSp          = "sp"
	GotoLive        = "live"
	GotoGame        = "game"
	GotoArticle     = "article"
	GotoActivity    = "activity_new"
	GotoTopic       = "topic_new"
	GotoDaily       = "daily"
	GotoRank        = "rank"
	GotoCard        = "card"
	GotoVeidoCard   = "video_card"
	GotoSpecialCard = "special_card"
	GotoTagCard     = "tag_card"
	GotoColumn      = "column"
	GotoColumnStage = "column_stage"
	GotoTagID       = "tag_id"
	GotoHotTopic    = "hottopic"
	GotoUpRcmdNew   = "up_rcmd_new"
	GotoUpRcmdNewV2 = "up_rcmd_new_v2"
	GotoHotDynamic  = "hot_dynamic"
	GotoEventTopic  = "event_topic"

	CardGotoAv       = int8(1)
	CardGotoTopic    = int8(2)
	CardGotoActivity = int8(3)

	// EnvPro is pro.
	EnvPro = "pro"
	// EnvTest is env.
	EnvTest = "test"
	// EnvDev is env.
	EnvDev = "dev"

	// movie copywriting
	CoverIng      = "即将上映"
	CoverPay      = "付费观看"
	CoverFree     = "免费观看"
	CoverVipFree  = "付费观看"
	CoverVipOnly  = "专享"
	CoverVipFirst = "抢先"
)

var (
	AvHandler = func(a *api.Arc) func(uri string) string {
		return func(uri string) string {
			if a == nil {
				return uri
			}
			if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				return fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			return uri
		}
	}
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidB || plat == PlatAndroidI
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI
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
	case GotoBangumiWeb:
		uri = "http://bangumi.bilibili.com/anime/" + param
	case GotoGame:
		uri = "bilibili://game/" + param
	case GotoSp:
		uri = "bilibili://splist/" + param
	case GotoWeb:
		uri = param
	case GotoDaily:
		uri = "bilibili://daily/" + param
	case GotoColumn:
		uri = "bilibili://pegasus/list/column/" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

func FillURIHotTopic(tagID, tagName string) (uri string) {
	return "bilibili://tag/" + tagID + "?name=" + tagName + "&show=dynamic"
}

func FillURIBangumi(gt, seasonID, episodeID string, episodeType int) (uri string) {
	var typeStr string
	switch episodeType {
	case 1, 4:
		typeStr = "anime"
	}
	switch gt {
	case GotoBangumi:
		uri = "http://bangumi.bilibili.com/" + typeStr + "/" + seasonID + "/play#" + episodeID
	}
	return
}

// FillURICategory deal app schema.
func FillURICategory(gt, columnID, sectionID string) (uri string) {
	if columnID == "" || sectionID == "" {
		return
	}
	switch gt {
	case GotoColumnStage:
		uri = "bilibili://pegasus/list/column/" + columnID + "/?sectionId=" + sectionID
	}
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
	case "iphone", "iphone_b":
		if device == "pad" {
			return PlatIPad
		}
		return PlatIPhone
	case "white":
		return PlatIPhone
	case "ipad":
		return PlatIPad
	case "android":
		return PlatAndroid
	case "win":
		return PlatWPhone
	case "android_G":
		return PlatAndroidG
	case "android_i":
		return PlatAndroidI
	case "android_b":
		return PlatAndroid
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

// MobiApp plat by mobi_app
func MobiApp(plat int8) string {
	switch plat {
	case PlatAndroid:
		return "android"
	case PlatIPhone:
		return "iphone"
	case PlatIPad:
		return "ipad"
	case PlatAndroidI:
		return "android_i"
	case PlatIPhoneI:
		return "iphone_i"
	case PlatIPadI:
		return "ipad_i"
	case PlatAndroidG:
		return "android_G"
	}
	return "iphone"
}

func StatusMark(status int8) string {
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

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

func IsGoto(gt string) bool {
	return gt == GotoAv || gt == GotoWeb || gt == GotoBangumi || gt == GotoSp || gt == GotoLive || gt == GotoGame
}

func MobiAPPBuleChange(mobiApp string) string {
	switch mobiApp {
	case "android_b":
		return "android"
	case "iphone_b":
		return "iphone"
	}
	return mobiApp
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
