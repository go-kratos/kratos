package model

import (
	"fmt"
	"strings"

	"go-common/app/interface/main/app-resource/model/tab"
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
	// PlatAndroidB is int8 for Android Blue.
	PlatAndroidB = int8(9)
	// PlatIPhoneB is int8 for Ios Blue
	PlatIPhoneB = int8(10)
	// PlatBilistudio is int8 for bilistudio
	PlatBilistudio = int8(11)
	// PlatAndroidTVYST is int8 for AndroidTV_YST Global.
	PlatAndroidTVYST = int8(12)

	GotoAv         = "av"
	GotoWeb        = "web"
	GotoBangumi    = "bangumi"
	GotoSp         = "sp"
	GotoLive       = "live"
	GotoGame       = "game"
	GotoPegasusTab = "pegasus"
)

var (
	PegasusHandler = func(m *tab.Menu) func(uri string) string {
		return func(uri string) string {
			if m == nil {
				return uri
			}
			if m.Name != "" {
				return fmt.Sprintf("%s?name=%s", uri, m.Name)
			}
			return uri
		}
	}
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
	case "android_b":
		return PlatAndroidB
	case "iphone_i":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_i":
		return PlatIPadI
	case "iphone_b":
		return PlatIPhoneB
	case "android_tv":
		return PlatAndroidTV
	case "android_tv_yst":
		return PlatAndroidTVYST
	case "bilistudio":
		return PlatBilistudio
	case "biliLink":
		return PlatIPhone
	}
	return PlatIPhone
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

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidB || plat == PlatAndroidI ||
		plat == PlatBilistudio || plat == PlatAndroidTV || plat == PlatAndroidTVYST
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI || plat == PlatIPhoneB
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
		uri = "bilibili://game/" + param
	case GotoSp:
		uri = "bilibili://splist/" + param
	case GotoWeb:
		uri = param
	case GotoPegasusTab:
		uri = "bilibili://pegasus/op/" + param
	}
	if f != nil {
		uri = f(uri)
	}
	return
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

func URLHTTPS(uri string) (url string) {
	if strings.HasPrefix(uri, "http://") {
		url = "https://" + uri[7:]
	} else {
		url = uri
	}
	return
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
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
