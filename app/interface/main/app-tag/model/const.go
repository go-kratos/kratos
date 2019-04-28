package model

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

	GotoAv      = "av"
	GotoBangumi = "bangumi"
	GotoLive    = "live"
	GotoWeb     = "web"
	GotoGame    = "game"
)

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

// FillURI deal app schema.
func FillURI(gt, param string) (uri string) {
	if param == "" {
		return
	}
	switch gt {
	case GotoAv:
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "bilibili://bangumi/season/" + param
	case GotoGame:
		uri = "bilibili://game/" + param
	case GotoWeb:
		uri = param
	}
	return
}

func FillURIBangumi(gt, seasonID, episodeID string, episodeType int) (uri string) {
	var typeStr string
	switch episodeType {
	case 1, 4:
		typeStr = "anime"
	}
	switch gt {
	case GotoBangumi:
		uri = "https://bangumi.bilibili.com/" + typeStr + "/" + seasonID + "/play#/" + episodeID
	}
	return
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// MobiAPPBuleChange
func MobiAPPBuleChange(mobiApp string) string {
	switch mobiApp {
	case "android_b":
		return "android"
	case "iphone_b":
		return "iphone"
	}
	return mobiApp
}
