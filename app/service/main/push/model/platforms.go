package model

import "strings"

// PushSDK* for parameter 'push_sdk' in http report API.
const (
	// PushSDKApns apns sdk.
	PushSDKApns = 1
	// PushSDKXiaomi mipush sdk.
	PushSDKXiaomi = 2
	// PushSDKHuawei huawei sdk.
	PushSDKHuawei = 3
	// PushSDKOppo oppo sdk.
	PushSDKOppo = 5
	// PushSDKJpush jpush sdk.
	PushSDKJpush = 6
	// PushSDKFCM fcm sdk
	PushSDKFCM = 7
)

const (
	// PlatformUnknown unknown.
	PlatformUnknown = 0
	// PlatformAndroid Android.
	PlatformAndroid = 1
	// PlatformIPhone iPhone.
	PlatformIPhone = 2
	// PlatformIPad iPad.
	PlatformIPad = 3
	// PlatformXiaomi mipush.
	PlatformXiaomi = 4
	// PlatformHuawei huawei.
	PlatformHuawei = 5
	// PlatformOppo oppo.
	PlatformOppo = 8
	// PlatformJpush jpush.
	PlatformJpush = 9
	// PlatformFCM fcm
	PlatformFCM = 10
)

// Platforms all platform
var Platforms = []int{
	PlatformIPhone,
	PlatformIPad,
	PlatformXiaomi,
	PlatformHuawei,
	PlatformOppo,
	PlatformJpush,
	PlatformFCM,
}

// Platform gets real platform.
func Platform(platform string, pushSDK int) int {
	switch pushSDK {
	case PushSDKApns:
		platform = strings.ToLower(platform)
		if strings.HasPrefix(platform, "iphone") {
			return PlatformIPhone
		} else if strings.HasPrefix(platform, "ipad") {
			return PlatformIPad
		}
	case PushSDKXiaomi:
		return PlatformXiaomi
	case PushSDKHuawei:
		return PlatformHuawei
	case PushSDKOppo:
		return PlatformOppo
	case PushSDKJpush:
		return PlatformJpush
	case PushSDKFCM:
		return PlatformFCM
	}
	// TODO add more brands
	return PlatformUnknown
}
