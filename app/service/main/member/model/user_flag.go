package model

const (
	// NickUpdated 首次改昵称
	NickUpdated = uint(1)
)

// HasAttr get attr.
func HasAttr(flag uint, bit uint) bool {
	return flag&bit == bit
}

// SetAttr set attr.
func SetAttr(flag uint, bit uint) uint {
	return flag | bit
}
