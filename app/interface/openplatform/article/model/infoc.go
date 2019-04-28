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
	// PlatH5 is int8 for H5
	PlatH5 = int8(9)
	// PlatPC is int8 for PC
	PlatPC = int8(10)
	//PlatOther is int8 for unknow plat
	PlatOther = int8(11)
)

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
	case "iphone_i":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_i":
		return PlatIPadI
	case "android_tv":
		return PlatAndroidTV
	case "h5":
		return PlatH5
	case "pc":
		return PlatPC
	}
	return PlatOther
}

// Client 成转换AI部门的client
func Client(plat int8) string {
	switch plat {
	case PlatIPad, PlatIPadI:
		return "ipad"
	case PlatIPhone, PlatIPhoneI:
		return "iphone"
	case PlatAndroid, PlatAndroidG, PlatAndroidI, PlatAndroidTV:
		return "android"
	default:
		return "web"
	}
}

// HistoryClient .
func HistoryClient(plat int8) (client int8) {
	switch plat {
	case PlatAndroid, PlatAndroidG, PlatAndroidI:
		client = 3
	case PlatIPhone, PlatIPhoneI:
		client = 1
	case PlatPC, PlatH5:
		client = 2
	case PlatAndroidTV:
		client = 33
	case PlatIPad, PlatIPadI:
		client = 4
	case PlatWPhone:
		client = 6
	}
	return
}
