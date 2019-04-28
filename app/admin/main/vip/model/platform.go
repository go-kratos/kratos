package model

import "go-common/library/time"

// const .
var (
	// PlatformMap 平台
	PlatformMap = map[string]string{"android": "android", "ios": "ios", "pc": "pc", "public": "public"}
	// Device 对应设备
	DeviceMap = map[string]string{"pad": "pad", "phone": "phone"}
	// MobiAPPIDIosMap iOS
	MobiAPPIDIosMap = map[string]string{"iphone": "iphone", "ipad": "ipad", "iphone_b": "iphone_b"}
	// MobiAPPIDAndroidMap Android
	MobiAPPIDAndroidMap = map[string]string{"android": "android", "android_tv_yst": "android_tv_yst", "android_tv": "android_tv", "android_i": "android_i", "android_b": "android_b"}
)

// ConfPlatform struct .
type ConfPlatform struct {
	ID           int64     `gorm:"column:id" json:"id" form:"id"`
	PlatformName string    `gorm:"column:platform_name" json:"platform_name" form:"platform_name" validate:"required"`
	Platform     string    `gorm:"column:platform" json:"platform" form:"platform" validate:"required"`
	Device       string    `gorm:"column:device" json:"device" form:"device"`
	MobiApp      string    `gorm:"column:mobi_app" json:"mobi_app" form:"mobi_app"`
	PanelType    string    `gorm:"column:panel_type" json:"panel_type" form:"panel_type" default:"normal"`
	IsDel        int8      `gorm:"column:is_del" json:"is_del" form:"is_del"`
	Operator     string    `gorm:"column:operator" json:"operator" form:"operator"`
	Ctime        time.Time `gorm:"column:ctime" json:"ctime" form:"ctime"`
	Mtime        time.Time `gorm:"column:mtime" json:"mtime" form:"mtime"`
}

// TableName for grom.
func (s *ConfPlatform) TableName() string {
	return "vip_platform_config"
}

// TypePlatform struct .
type TypePlatform struct {
	ID           int64  `json:"id"`
	PlatformName string `json:"platform_name"`
}
