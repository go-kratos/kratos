package model

const (
	// Open is open switch
	Open = 1
	// Close is close switch
	Close = 0
)

//UpSwitch for db.
type UpSwitch struct {
	ID        int64 `json:"id"`
	MID       int64 `json:"mid"`
	Attribute int   `json:"attribute"`
}

// AttrSet set attribute.
func (u *UpSwitch) AttrSet(v int, bit uint8) {
	u.Attribute = u.Attribute&(^(1 << bit)) | (v << bit)
}

// AttrVal get attribute.
func (u *UpSwitch) AttrVal(bit uint8) int {
	return (u.Attribute >> bit) & int(1)
}

// Const State
const (
	// AttrPlayer flow up window 's switch of attribute bit
	AttrPlayer = uint8(0)
	// AttrHonorWeekly honor weekly subscription 's switch of attribute bit
	AttrHonorWeekly = uint8(1)
)
