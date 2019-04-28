package model

const (
	// IsUp is up
	IsUp = 1
	// NotUp not up
	NotUp = 0
)

//Up for db.
type Up struct {
	ID        int64 `json:"id"`
	MID       int64 `json:"mid"`
	Attribute int   `json:"attribute"`
}

//Info for auth by all platform
type Info struct {
	Archive     int `json:"archive"`
	ArchiveFake int `json:"archive_fake"`
}

// IdentifyAll for all type of uper identify.
type IdentifyAll struct {
	Archive int `json:"archive"`
	Article int `json:"article"`
	Pic     int `json:"pic"`
	Blink   int `json:"blink"`
}

// AttrSet set attribute.
func (u *Up) AttrSet(v int, bit uint8) {
	u.Attribute = u.Attribute&(^(1 << bit)) | (v << bit)
}

// AttrVal get attribute.
func (u *Up) AttrVal(bit uint8) int {
	return (u.Attribute >> bit) & int(1)
}

// Const State
const (
	// AttrNo attribute yes and no
	AttrNo  = int(0)
	AttrYes = int(1)
	// AttrArchiveUp attribute bit
	AttrArchiveUp    = uint8(0)
	AttrArchiveNewUp = uint8(1)
	AttrLiveUp       = uint8(2)
	AttrLiveWhiteUp  = uint8(3)
)

var (
	_attr = map[int]int{
		AttrNo:  AttrNo,
		AttrYes: AttrYes,
	}
	_bits = map[uint8]string{
		AttrArchiveUp:    "稿件作者-有过审稿",
		AttrArchiveNewUp: "稿件作者-有投过稿",
		AttrLiveUp:       "直播作者",
		AttrLiveWhiteUp:  "直播白名单",
	}
)

// BitDesc return bit desc.
func BitDesc(bit uint8) (desc string) {
	return _bits[bit]
}

// InAttr in correct attrs.
func InAttr(attr int) (ok bool) {
	_, ok = _attr[attr]
	return
}

// ListUpBaseArg arg
type ListUpBaseArg struct {
	LastID   int64   `form:"last_id"`
	Size     int     `form:"size"`
	Activity []int64 `form:"activity,split"`
}

// Validate ListUpBaseArg
func (arg *ListUpBaseArg) Validate() bool {
	if arg.Size < 0 || arg.Size > 1000 || arg.LastID < 0 {
		return false
	}
	if len(arg.Activity) > 0 {
		for _, v := range arg.Activity {
			if v < 0 || v > 4 {
				return false
			}
		}
	}
	return true
}
