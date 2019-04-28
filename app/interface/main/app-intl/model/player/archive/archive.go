package archive

// is
const (
	StateOpen       = int32(0)
	AttrNo          = int32(0)
	AttrYes         = int32(1)
	AttrBitBadgepay = uint(18)
	AttrBitUGCPay   = uint(22)
	AttrBitIsPGC    = uint(9)
)

// IsNormal check archive is normal.
func (info *Info) IsNormal() bool {
	return info.State >= StateOpen
}

// IsPGC is.
func (info *Info) IsPGC() bool {
	return info.AttrVal(AttrBitIsPGC) == AttrYes
}

// AttrVal get attr val by bit.
func (info *Info) AttrVal(bit uint) int32 {
	return (info.Attribute >> bit) & int32(1)
}

// HasCid check cid is in info.Cids.
func (info *Info) HasCid(cid int64) (ok bool) {
	for _, id := range info.Cids {
		if cid == id {
			ok = true
			break
		}
	}
	return
}
