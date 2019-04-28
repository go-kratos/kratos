package bloomfilter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/Dai0522/go-hash/murmur3"
)

// BloomFilter .
type BloomFilter struct {
	strategy Strategy
	bits     Bitmap
	numHash  int
}

// New BloomFilter
func New(expect uint64, fpp float64) (*BloomFilter, error) {
	m := optimalNumOfBits(expect, fpp)
	b, err := NewLockFreeBitmap(m)
	if err != nil {
		return nil, err
	}
	return &BloomFilter{
		strategy: &Murur3_128Strategy{
			hashFunc: murmur3.New(),
		},
		bits:    b,
		numHash: optimalNumOfHash(expect, m),
	}, nil
}

// Put little endian byte array
func (bf *BloomFilter) Put(b []byte) bool {
	return bf.strategy.Put(b, bf.numHash, bf.bits)
}

// PutUint16 .
func (bf *BloomFilter) PutUint16(i uint16) bool {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, i)
	return bf.strategy.Put(b, bf.numHash, bf.bits)
}

// PutUint32 .
func (bf *BloomFilter) PutUint32(i uint32) bool {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return bf.strategy.Put(b, bf.numHash, bf.bits)
}

// PutUint64 .
func (bf *BloomFilter) PutUint64(i uint64) bool {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return bf.strategy.Put(b, bf.numHash, bf.bits)
}

// MightContain little endian byte array
func (bf *BloomFilter) MightContain(b []byte) bool {
	return bf.strategy.MightContain(b, bf.numHash, bf.bits)
}

// MightContainUint16 .
func (bf *BloomFilter) MightContainUint16(i uint16) bool {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, i)
	return bf.strategy.MightContain(b, bf.numHash, bf.bits)
}

// MightContainUint32 .
func (bf *BloomFilter) MightContainUint32(i uint32) bool {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return bf.strategy.MightContain(b, bf.numHash, bf.bits)
}

// MightContainUint64 .
func (bf *BloomFilter) MightContainUint64(i uint64) bool {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return bf.strategy.MightContain(b, bf.numHash, bf.bits)
}

// ExpectedFpp returns the probability that mightContain will erroneously
// return true for an object that has not actually been put in
func (bf *BloomFilter) ExpectedFpp() float64 {
	return math.Pow(float64(bf.bits.BitCount())/float64(bf.bits.BitSize()), float64(bf.numHash))
}

// ApproximateElementCount returns an estimate for the total number of
// distinct elements that have been added to this Bloom filter. This
// approximation is reasonably accurate if it does not exceed the value
// of that was used when constructing the filter
func (bf *BloomFilter) ApproximateElementCount() int {
	size := float64(bf.bits.BitSize())
	count := float64(bf.bits.BitCount())
	res := math.Log1p(-1*count/size) * size / float64(bf.numHash)
	return int(math.Ceil(res + 0.5))
}

// Serialized serialized bloom filter
func (bf *BloomFilter) Serialized() *[]byte {
	// Serial form:
	// 1 signed byte for the strategy
	// 1 unsigned byte for the number of hash functions
	// 1 big endian int, the number of longs in our bitset
	// N big endian longs of our bitset
	var buf bytes.Buffer
	buf.WriteByte(byte(1))
	buf.WriteByte(byte(bf.numHash))

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, bf.bits.Size())
	buf.Write(size)

	dataBuf := make([]byte, 8)
	data := *bf.bits.Data()
	for _, v := range data {
		binary.BigEndian.PutUint64(dataBuf, v)
		buf.Write(dataBuf)
	}

	res := buf.Bytes()
	return &res
}

// Load load serialized bloom filter into memory
func Load(b *[]byte) (*BloomFilter, error) {
	if len(*b) < 10 {
		return nil, errors.New("invaled data")
	}
	numHash := int((*b)[1])
	length := binary.BigEndian.Uint32((*b)[2:6])
	data := make([]uint64, length)
	for i := 0; i < int(length); i++ {
		j := (i * 8) + 6
		data[i] = binary.BigEndian.Uint64((*b)[j : j+8])
	}
	bits := LoadLockFreeBitmap(&data)

	bf := &BloomFilter{
		strategy: &Murur3_128Strategy{
			hashFunc: murmur3.New(),
		},
		bits:    bits,
		numHash: numHash,
	}
	return bf, nil
}

// Merge return dst bloom filter ptr
func Merge(src *BloomFilter, dst *BloomFilter) *BloomFilter {
	if src == nil || dst == nil {
		return dst
	}

	dst.bits.Merge(src.bits.Data())
	return dst
}

func optimalNumOfHash(n, m uint64) int {
	return int(math.Max(1, math.Floor((float64(m/n)*math.Log(2))+0.5)))
}

func optimalNumOfBits(n uint64, p float64) uint64 {
	if p == 0.0 {
		p = math.SmallestNonzeroFloat64
	}
	return uint64(-1 * float64(n) * math.Log(p) / (math.Log(2) * math.Log(2)))
}
