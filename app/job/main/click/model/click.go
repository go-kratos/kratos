package model

//web+h5+outside+ios+android
const (
	TypeForWeb                   = "web"
	TypeForH5                    = "h5"
	TypeForOutside               = "outside"
	TypeForIOS                   = "ios"
	TypeForAndroid               = "android"
	TypeForAndroidTv             = "android_tv"
	PlatForWeb                   = int8(0)
	PlatForH5                    = int8(1)
	PlatForOuter                 = int8(2)
	PlatForIos                   = int8(3)
	PlatForAndroid               = int8(4)
	PlatForAndroidTV             = int8(5)
	PlatForAutoPlayIOS           = int8(6)
	PlafForAutoPlayInlineIOS     = int8(7)
	PlatForAutoPlayAndroid       = int8(8)
	PlatForAutoPlayInlineAndroid = int8(9)
	_maxDBTimes                  = 6
)

// ClickInfo is
type ClickInfo struct {
	Aid            int64
	Web            int64
	H5             int64
	Outer          int64
	Ios            int64
	Android        int64
	AndroidTV      int64
	Sum            int64
	DBTimes        int
	LastChangeTime int64
}

// NeedRelease is
func (c *ClickInfo) NeedRelease() bool {
	if c.DBTimes > _maxDBTimes {
		return true
	}
	return false
}

// ArcDuration is
type ArcDuration struct {
	Duration int64
	GotTime  int64
}

// Ready is
func (c *ClickInfo) Ready(ts int64) {
	c.Sum = c.Sum + c.GetSum()
	c.LastChangeTime = ts
	c.Web, c.H5, c.Outer, c.Ios, c.Android, c.AndroidTV = 0, 0, 0, 0, 0, 0
	c.DBTimes++
}

// GetSum is
func (c *ClickInfo) GetSum() (sum int64) {
	sum = c.Web + c.H5 + c.Outer + c.Ios + c.Android + c.AndroidTV
	return
}
