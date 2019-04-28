package bloomfilter

import (
	"encoding/binary"
	"math"

	"github.com/Dai0522/go-hash/murmur3"
)

// Strategy .
type Strategy interface {
	Put([]byte, int, Bitmap) bool
	MightContain([]byte, int, Bitmap) bool
}

// Murur3_128Strategy .
type Murur3_128Strategy struct {
	hashFunc *murmur3.Murmur3
}

// Put .
func (s *Murur3_128Strategy) Put(b []byte, n int, bits Bitmap) bool {
	bitSize := bits.BitSize()
	hashCode := s.hashFunc.Murmur3_128(b)
	h1 := binary.LittleEndian.Uint64(hashCode[:8])
	h2 := binary.LittleEndian.Uint64(hashCode[8:])

	bitsChanged := false
	combine := h1
	for i := 0; i < n; i++ {
		res := bits.Set((combine & math.MaxInt64) % bitSize)
		bitsChanged = bitsChanged || res
		combine += h2
	}
	return bitsChanged
}

// MightContain .
func (s *Murur3_128Strategy) MightContain(b []byte, n int, bits Bitmap) bool {
	bitSize := bits.BitSize()
	hashCode := s.hashFunc.Murmur3_128(b)
	h1 := binary.LittleEndian.Uint64(hashCode[:8])
	h2 := binary.LittleEndian.Uint64(hashCode[8:])

	combine := h1
	for i := 0; i < n; i++ {
		if !bits.Get((combine & math.MaxInt64) % bitSize) {
			return false
		}
		combine += h2
	}
	return true
}
