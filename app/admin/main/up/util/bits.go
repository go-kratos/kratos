package util

//IsBitSet bit is from 0 to 31
func IsBitSet(attr int, bit uint) bool {
	return IsBitSet64(int64(attr), bit)
}

// IsBitSet64 bit is from 0 to 63
func IsBitSet64(attr int64, bit uint) bool {
	if bit >= 64 {
		return false
	}

	return (attr & (1 << bit)) != 0
}

//SetBit64 set bit to 1
func SetBit64(attr int64, bit uint) int64 {
	return attr | (1 << bit)
}

//UnSetBit64 set bit to 0
func UnSetBit64(attr int64, bit uint) int64 {
	return attr & ^(1 << bit)
}
