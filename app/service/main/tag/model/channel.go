package model

// const channel const value.
const (
	ChannelCategoryAttrINT = uint(0) //频道分类海外版位
)

// AttrVal get attr flag.
func (t *ChannelCategory) AttrVal(bit uint) int32 {
	return (t.Attr >> bit) & int32(1)
}

// AttrSet get attr flag.
func (t *ChannelCategory) AttrSet(bit uint, v int32) {
	t.Attr = t.Attr&(^(1 << bit)) | (v << bit)
}
