package murmur3

// Murmur3 .
type Murmur3 struct {
	seed uint32
}

// New .
func New() *Murmur3 {
	return NewWithSeed(0)
}

// NewWithSeed .
func NewWithSeed(s uint32) *Murmur3 {
	return &Murmur3{
		seed: s,
	}
}

// Murmur3_32 .
func (h *Murmur3) Murmur3_32(b []byte) []byte {
	return murmur3_32(h.seed, b)
}

// Murmur3_64 .
func (h *Murmur3) Murmur3_64(b []byte) []byte {
	return murmur3_64(h.seed, b)
}

// Murmur3_128 little endian []byte.
func (h *Murmur3) Murmur3_128(b []byte) []byte {
	h1, h2 := murmur3_128(h.seed, b)
	return []byte{
		byte(h1), byte(h1 >> 8), byte(h1 >> 16), byte(h1 >> 24),
		byte(h1 >> 32), byte(h1 >> 40), byte(h1 >> 48), byte(h1 >> 56),
		byte(h2), byte(h2 >> 8), byte(h2 >> 16), byte(h2 >> 24),
		byte(h2 >> 32), byte(h2 >> 40), byte(h2 >> 48), byte(h2 >> 56),
	}
}
