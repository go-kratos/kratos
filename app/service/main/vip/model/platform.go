package model

const (
	// PanelTypeDefault .
	PanelTypeDefault = "normal"
	// PlatformWeb .
	PlatformWeb = "web"
	// PlatformAndroid .
	PlatformAndroid = "android"
	// PlatformIos .
	PlatformIos = "ios"
	// DeviceIos .
	DeviceIos = "ios"
)

var (
	// PlatformMap 平台
	PlatformMap = map[string]string{"android": "android", "ios": "ios", "pc": "pc", "public": "public"}
	// DeviceMap 对应设备
	DeviceMap = map[string]string{"pad": "pad", "phone": "phone"}
	// MobiAPPIDIosMap iOS
	MobiAPPIDIosMap = map[string]string{"iphone": "iphone", "ipad": "ipad", "iphone_b": "iphone_b"}
	// MobiAPPIDAndroidMap Android
	MobiAPPIDAndroidMap = map[string]string{"android": "android", "android_tv_yst": "android_tv_yst", "android_tv": "android_tv", "android_i": "android_i", "android_b": "android_b"}
	// MobiAPPIDMap all
	MobiAPPIDMap = map[string]string{"iphone": "iphone", "ipad": "ipad", "iphone_b": "iphone_b", "android": "android", "android_tv_yst": "android_tv_yst", "android_tv": "android_tv", "android_i": "android_i", "android_b": "android_b"}
)

// ConfPlatform struct .
type ConfPlatform struct {
	ID           int64  `json:"id"`
	PlatformName string `json:"platform_name"`
	Platform     string `json:"platform"`
	Device       string `json:"device"`
	MobiApp      string `json:"mobi_app"`
	PanelType    string `json:"panel_type"`
}
