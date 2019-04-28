package murmur3

import (
	"unsafe"
)

const (
	c1_128  = 0x87c37b91114253d5
	c2_128  = 0x4cf5ad432745937f
	size128 = 16
)

// murmur3_128
func murmur3_128(seed uint32, b []byte) (uint64, uint64) {
	h1 := uint64(seed)
	h2 := uint64(seed)

	length := len(b)
	h1, h2 = bmix64(b, h1, h2)

	// tail
	tail := b[length-(length&15):]
	var k1, k2 uint64
	switch length & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1_128
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2_128
		h1 ^= k1
	}

	h1 ^= uint64(length)
	h2 ^= uint64(length)

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1
	return h1, h2
}

// bmix64 .
func bmix64(b []byte, h1, h2 uint64) (uint64, uint64) {
	nblocks := len(b) / 16

	// body
	for i := 0; i < nblocks; i++ {
		t := (*[2]uint64)(unsafe.Pointer(&b[i*16]))
		k1, k2 := t[0], t[1]

		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2_128
		h1 ^= k1

		h1 = (h1 << 27) | (h1 >> 37) // rotl64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1_128
		h2 ^= k2

		h2 = (h2 << 31) | (h2 >> 33) // rotl64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	return h1, h2
}

// fmix64 .
func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33

	return k
}
