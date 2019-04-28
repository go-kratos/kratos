package model

// associate_platform
const (
	AssociatePlatformNone int8 = iota
	AssociatePlatformAndroidPink
	AssociatePlatformIphonePink
	AssociatePlatformIpadPink
)

// AssociatePlatform get platfrom.
func AssociatePlatform(platfrom, device, mobiApp string) int8 {
	switch {
	case platfrom == "ios" && device == "phone" && mobiApp == "iphone":
		return AssociatePlatformIphonePink
	case platfrom == "ios" && device == "pad" && mobiApp == "iphone":
		return AssociatePlatformIpadPink
	case platfrom == "android" && mobiApp == "android":
		return AssociatePlatformAndroidPink
	default:
		return AssociatePlatformNone
	}
}

// AssociateVipResp  associate vip resp
type AssociateVipResp struct {
	Title             string `json:"title"`
	Remark            string `json:"remark"`
	Link              string `json:"link"`
	AssociatePlatform int8   `json:"associate_platform"`
}
