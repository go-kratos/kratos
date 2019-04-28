package model

// is
const (
	ValueForLocked = int8(1)
)

// Forbid is
type Forbid struct {
	AID    int64
	Plat   int8
	Lv     int8
	Locked int8
}
