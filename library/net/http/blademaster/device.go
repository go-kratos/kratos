package blademaster

import (
	"strconv"

	"go-common/library/net/metadata"
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
)

// Device is the mobile device model
type Device struct {
	Build       int64
	Buvid       string
	Buvid3      string
	Channel     string
	Device      string
	Sid         string
	RawPlatform string
	RawMobiApp  string
}

// Mobile is the default handler
func Mobile() HandlerFunc {
	return func(ctx *Context) {
		req := ctx.Request
		dev := new(Device)
		dev.Buvid = req.Header.Get("Buvid")
		if buvid3, err := req.Cookie("buvid3"); err == nil && buvid3 != nil {
			dev.Buvid3 = buvid3.Value
		}
		if sid, err := req.Cookie("sid"); err == nil && sid != nil {
			dev.Sid = sid.Value
		}
		if build, err := strconv.ParseInt(req.Form.Get("build"), 10, 64); err == nil {
			dev.Build = build
		}
		dev.Channel = req.Form.Get("channel")
		dev.Device = req.Form.Get("device")
		dev.RawMobiApp = req.Form.Get("mobi_app")
		dev.RawPlatform = req.Form.Get("platform")
		ctx.Set("device", dev)
		if md, ok := metadata.FromContext(ctx); ok {
			md[metadata.Device] = dev
		}
	}
}

// Plat return platform from raw platform and mobiApp
func (d *Device) Plat() int8 {
	switch d.RawMobiApp {
	case "iphone":
		if d.Device == "pad" {
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
		if d.Device == "pad" {
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
	}
	return PlatIPhone
}

// IsAndroid check plat is android or ipad.
func (d *Device) IsAndroid() bool {
	plat := d.Plat()
	return plat == PlatAndroid ||
		plat == PlatAndroidG ||
		plat == PlatAndroidB ||
		plat == PlatAndroidI ||
		plat == PlatBilistudio ||
		plat == PlatAndroidTV ||
		plat == PlatAndroidTVYST
}

// IsIOS check plat is iphone or ipad.
func (d *Device) IsIOS() bool {
	plat := d.Plat()
	return plat == PlatIPad ||
		plat == PlatIPhone ||
		plat == PlatIPadI ||
		plat == PlatIPhoneI ||
		plat == PlatIPhoneB
}

// IsOverseas is overseas
func (d *Device) IsOverseas() bool {
	plat := d.Plat()
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// InvalidChannel check source channel is not allow by config channel.
func (d *Device) InvalidChannel(cfgCh string) bool {
	plat := d.Plat()
	return plat == PlatAndroid && cfgCh != "*" && cfgCh != d.Channel
}

// MobiApp by plat
func (d *Device) MobiApp() string {
	plat := d.Plat()
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

// MobiAPPBuleChange is app blue change.
func (d *Device) MobiAPPBuleChange() string {
	switch d.RawMobiApp {
	case "android_b":
		return "android"
	case "iphone_b":
		return "iphone"
	}
	return d.RawMobiApp
}
