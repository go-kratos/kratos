package jpush

const (
	// PlatformIOS .
	PlatformIOS = "ios"
	// PlatformAndroid .
	PlatformAndroid = "android"
	// PlatformWinphone .
	PlatformWinphone = "winphone"
	// PlatformAll .
	PlatformAll = "all"
)

// Platform .
type Platform struct {
	OS      interface{}
	osArray []string
}

// NewPlatform .
func NewPlatform(os ...string) *Platform {
	p := new(Platform)
	for _, v := range os {
		switch v {
		case PlatformIOS, PlatformAndroid, PlatformWinphone:
			p.osArray = append(p.osArray, v)
		case PlatformAll:
			p.OS = PlatformAll
			return p
		}
	}
	p.OS = p.osArray
	return p
}
