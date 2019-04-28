package model

import "net"

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
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidI
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
	return plat == PlatIPad || plat == PlatIPadI || plat == PlatIpadHD
}

// IsPlayerScreen check plat is iphone or ipad.
func IsPlayerScreen(tagID int64) bool {
	return tagID == AndroidPlayerScreen || tagID == AndroidPlayerScreenNothing || tagID == AndroidPlayerScreenDlna || tagID == AndroidPlayerScreenTV || tagID == IOSPlayerScreen || tagID == IOSPlayerScreenNothing || tagID == IOSPlayerScreenDlna || tagID == IOSPlayerScreenTV
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
	case "android":
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
	}
	return PlatIPhone
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// InetAtoN conver ip addr to uint32.
func InetAtoN(s string) (sum uint32) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
	ip = ip.To4()
	if ip == nil {
		return
	}
	sum += uint32(ip[0]) << 24
	sum += uint32(ip[1]) << 16
	sum += uint32(ip[2]) << 8
	sum += uint32(ip[3])
	return sum
}

// InetNtoA conver uint32 to ip addr.
func InetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}
