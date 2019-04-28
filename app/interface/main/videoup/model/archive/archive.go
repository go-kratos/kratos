package archive

import (
	"go-common/library/time"
)

// State + Attr + Copyright + Upfrom
const (
	// open state
	StateOpen   = int8(0)
	StateOrange = int8(1)
	// forbit state
	StateForbidWait     = int8(-1)
	StateForbidRecicle  = int8(-2)
	StateForbidPolice   = int8(-3)
	StateForbidLock     = int8(-4)
	StateForbidFackLock = int8(-5)
	StateForbidFixed    = int8(-6)
	StateForbidLater    = int8(-7)
	// StateForbidPatched   = int8(-8)
	StateForbidWaitXcode  = int8(-9)
	StateForbidAdminDelay = int8(-10)
	StateForbidFixing     = int8(-11)
	// StateForbidStorageFail = int8(-12)
	StateForbidOnlyComment = int8(-13)
	// StateForbidTmpRecicle  = int8(-14)
	StateForbidDispatch  = int8(-15)
	StateForbidXcodeFail = int8(-16)
	StateForbidSubmit    = int8(-30)
	StateForbidUserDelay = int8(-40)
	StateForbidUpDelete  = int8(-100)
	// attribute yes and no
	AttrYes = int32(1)
	AttrNo  = int32(0)
	// attribute bit
	AttrBitNoRank      = uint(0)
	AttrBitNoIndex     = uint(1)
	AttrBitNoWeb       = uint(2)
	AttrBitNoMobile    = uint(3)
	AttrBitNoSearch    = uint(4)
	AttrBitOverseaLock = uint(5)
	AttrBitNoRecommend = uint(6)
	// AttrBitHideCoins     = uint(7)
	AttrBitHasHD5 = uint(8)
	// AttrBitVisitorDm     = uint(9)
	AttrBitAllowBp   = uint(10)
	AttrBitIsBangumi = uint(11)
	// AttrBitAllowDownload = uint(12)
	AttrBitHideClick    = uint(13)
	AttrBitAllowTag     = uint(14)
	AttrBitIsFromArcAPI = uint(15)
	AttrBitJumpURL      = uint(16)
	AttrBitIsMovie      = uint(17)
	AttrBitBadgepay     = uint(18)
	AttrBitStaff        = uint(24) //联合投稿
	// copyright state
	CopyrightUnknow   = int8(0)
	CopyrightOriginal = int8(1)
	CopyrightCopy     = int8(2)
	//up_from
	UpFromWeb         = int8(0)
	UpFromPGC         = int8(1)
	UpFromWindows     = int8(2)
	UpFromAPP         = int8(3)
	UpFromMAC         = int8(4)
	UpFromSecretPGC   = int8(5)
	UpFromCoopera     = int8(6)
	UpFromCreator     = int8(7)  // 创作姬
	UpFromAPPAndroid  = int8(8)  // 安卓主APP
	UpFromAPPiOS      = int8(9)  // iOS主APP
	UpFromCM          = int8(10) // Web商单用户投稿
	UpFromIpad        = int8(11) // ipad投稿的用户
	AdvertisingTypeID = 166      // 广告分区的typeid
)

var (
	_copyright = map[int8]int8{
		CopyrightUnknow:   CopyrightUnknow,
		CopyrightOriginal: CopyrightOriginal,
		CopyrightCopy:     CopyrightCopy,
	}
)

// InCopyrights check copyright in all copyrights.
func InCopyrights(cp int8) (ok bool) {
	_, ok = _copyright[cp]
	return
}

// Archive is archive model.
type Archive struct {
	Aid    int64 `json:"aid"`
	Mid    int64 `json:"mid"`
	TypeID int16 `json:"tid"`
	// HumanRank int       `json:"-"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Cover     string `json:"cover"`
	Tag       string `json:"tag"`
	Duration  int64  `json:"duration"`
	Copyright int8   `json:"copyright"`
	Source    string `json:"source"`
	NoReprint int8   `json:"no_reprint"`
	UgcPay    int8   `json:"ugcpay"`
	OrderID   int64  `json:"order_id"`
	Desc      string `json:"desc"`
	MissionID int    `json:"mission_id"`
	// Round     int8      `json:"-"`
	// Forward   int64     `json:"-"`
	Attribute int32 `json:"attribute"`
	// Access    int16     `json:"-"`
	// desc_format
	DescFormatID int    `json:"desc_format_id,omitempty"`
	State        int8   `json:"state"`
	StateDesc    string `json:"state_desc"`
	// dynamic
	Dynamic string  `json:"dynamic"`
	Porder  *Porder `json:"porder"`
	// time
	DTime time.Time `json:"dtime"`
	PTime time.Time `json:"ptime"`
	CTime time.Time `json:"ctime"`
	// MTime     time.Time `json:"-"`
}

// NotAllowUp check archive is or not allow update state.
func (a *Archive) NotAllowUp() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidLock || a.State == StateForbidPolice
}

// AttrVal get attribute.
func (a *Archive) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// Type type from archive
type Type struct {
	ID          int16  `json:"id"`
	PID         int16  `json:"pid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DescFormat str
type DescFormat struct {
	ID        int   `json:"id"`
	TypeID    int16 `json:"typeid"`
	Copyright int8  `json:"copyright"`
	Lang      int8  `json:"lang"`
}

// FilterData filter-service data
type FilterData struct {
	Level  int64    `json:"level"`
	Limit  int64    `json:"limit"`
	Msg    string   `json:"msg"`
	TypeID []int64  `json:"typeid"`
	Hit    []string `json:"hit"`
}

// PayAsset  str
type PayAsset struct {
	Price         int            `json:"price"`
	PlatformPrice map[string]int `json:"platform_price"`
}
