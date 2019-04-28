package bloomfilter

import (
	"errors"
	"math"
	"sync/atomic"
)

const (
	// AddrBits .
	AddrBits = 6
)

// Bitmap .
type Bitmap interface {
	Set(uint64) bool
	Get(uint64) bool
	BitSize() uint64
	BitCount() uint64
	Size() uint32
	Data() *[]uint64
	Merge(*[]uint64) bool
}

// LockFreeBitmap .
type LockFreeBitmap struct {
	data     []uint64
	bitCount uint64
}

// NewLockFreeBitmap .
func NewLockFreeBitmap(bits uint64) (*LockFreeBitmap, error) {
	size := int(math.Ceil(float64(bits) / 64.0))
	if size <= 0 {
		err := errors.New("data length is zero")
		return nil, err
	}
	bm := &LockFreeBitmap{
		data:     make([]uint64, size),
		bitCount: 0,
	}
	return bm, nil
}

// LoadLockFreeBitmap .
func LoadLockFreeBitmap(d *[]uint64) *LockFreeBitmap {
	count := uint64(0)
	for _, v := range *d {
		count += bitCount(v)
	}
	return &LockFreeBitmap{
		data:     *d,
		bitCount: count,
	}
}

// Set .
func (bits *LockFreeBitmap) Set(bitIndex uint64) bool {
	if bits.Get(bitIndex) {
		return false
	}

	longIndex := bitIndex >> AddrBits
	mask := uint64(1 << (bitIndex & 63))

	for {
		old := bits.data[longIndex]
		new := old | mask
		if old == new {
			return false
		}
		if atomic.CompareAndSwapUint64(&bits.data[longIndex], old, new) {
			break
		}
	}

	atomic.AddUint64(&bits.bitCount, 1)
	return true
}

// Get .
func (bits *LockFreeBitmap) Get(bitIndex uint64) bool {
	return bits.data[bitIndex>>AddrBits]&(uint64(1<<(bitIndex&63))) != 0
}

// BitSize .
func (bits *LockFreeBitmap) BitSize() uint64 {
	return uint64(len(bits.data) * 64)
}

// BitCount .
func (bits *LockFreeBitmap) BitCount() uint64 {
	return bits.bitCount
}

// Size .
func (bits *LockFreeBitmap) Size() uint32 {
	return uint32(len(bits.data))
}

// Data .
func (bits *LockFreeBitmap) Data() *[]uint64 {
	return &bits.data
}

// Merge .
func (bits *LockFreeBitmap) Merge(data *[]uint64) bool {
	if len(bits.data) != len(*data) {
		return false
	}
	for i := 0; i < len(bits.data); i++ {
		for {
			old := bits.data[i]
			new := bits.data[i] | (*data)[i]
			if old == new {
				break
			}
			if atomic.CompareAndSwapUint64(&bits.data[i], old, new) {
				break
			}
		}
	}
	return true
}

func bitCount(i uint64) uint64 {
	var c uint64
	for ; i != 0; i = i >> 1 {
		if i&1 == 1 {
			c++
		}
	}
	return c
}
