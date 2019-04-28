package model

// resource const
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
	// PlatAndroidB is int8 for Android Bule.
	PlatAndroidB = int8(9)
	// PlatWEB is int8 for web.
	PlatWEB = int8(99)

	// goto
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

	CardGotoAv       = int8(1)
	CardGotoTopic    = int8(2)
	CardGotoActivity = int8(3)
)

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
	case "android":
		return PlatAndroid
	case "win":
		return PlatWPhone
	case "android_G":
		return PlatAndroidG
	case "android_I":
		return PlatAndroidI
	case "iphone_I":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_I":
		return PlatIPadI
	case "android_tv":
		return PlatAndroidTV
	}
	return PlatIPhone
}

// FillURI deal app schema.
func FillURI(gt, param string) (uri string) {
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
	return
}
