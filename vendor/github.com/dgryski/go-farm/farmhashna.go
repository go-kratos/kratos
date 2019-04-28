package farm

func shiftMix(val uint64) uint64 {
	return val ^ (val >> 47)
}

func hashLen16(u, v uint64) uint64 {
	return hash128to64(uint128{u, v})
}

func hashLen16Mul(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	a := (u ^ v) * mul
	a ^= (a >> 47)
	b := (v ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func hashLen0to16(s []byte) uint64 {
	slen := uint64(len(s))
	if slen >= 8 {
		mul := k2 + slen*2
		a := fetch64(s, 0) + k2
		b := fetch64(s, int(slen-8))
		c := rotate64(b, 37)*mul + a
		d := (rotate64(a, 25) + b) * mul
		return hashLen16Mul(c, d, mul)
	}

	if slen >= 4 {
		mul := k2 + slen*2
		a := fetch32(s, 0)
		return hashLen16Mul(slen+(uint64(a)<<3), uint64(fetch32(s, int(slen-4))), mul)
	}
	if slen > 0 {
		a := s[0]
		b := s[slen>>1]
		c := s[slen-1]
		y := uint32(a) + (uint32(b) << 8)
		z := uint32(slen) + (uint32(c) << 2)
		return shiftMix(uint64(y)*k2^uint64(z)*k0) * k2
	}
	return k2
}

// This probably works well for 16-byte strings as well, but it may be overkill
// in that case.
func hashLen17to32(s []byte) uint64 {
	slen := len(s)
	mul := k2 + uint64(slen*2)
	a := fetch64(s, 0) * k1
	b := fetch64(s, 8)
	c := fetch64(s, slen-8) * mul
	d := fetch64(s, slen-16) * k2
	return hashLen16Mul(rotate64(a+b, 43)+rotate64(c, 30)+d, a+rotate64(b+k2, 18)+c, mul)
}

// Return a 16-byte hash for 48 bytes.  Quick and dirty.
// Callers do best to use "random-looking" values for a and b.
func weakHashLen32WithSeedsWords(w, x, y, z, a, b uint64) (uint64, uint64) {
	a += w
	b = rotate64(b+a+z, 21)
	c := a
	a += x
	a += y
	b += rotate64(a, 44)
	return a + z, b + c
}

// Return a 16-byte hash for s[0] ... s[31], a, and b.  Quick and dirty.
func weakHashLen32WithSeeds(s []byte, a, b uint64) (uint64, uint64) {
	return weakHashLen32WithSeedsWords(fetch64(s, 0),
		fetch64(s, 8),
		fetch64(s, 16),
		fetch64(s, 24),
		a,
		b)
}

// Return an 8-byte hash for 33 to 64 bytes.
func hashLen33to64(s []byte) uint64 {
	slen := len(s)
	mul := k2 + uint64(slen)*2
	a := fetch64(s, 0) * k2
	b := fetch64(s, 8)
	c := fetch64(s, slen-8) * mul
	d := fetch64(s, slen-16) * k2
	y := rotate64(a+b, 43) + rotate64(c, 30) + d
	z := hashLen16Mul(y, a+rotate64(b+k2, 18)+c, mul)
	e := fetch64(s, 16) * mul
	f := fetch64(s, 24)
	g := (y + fetch64(s, slen-32)) * mul
	h := (z + fetch64(s, slen-24)) * mul
	return hashLen16Mul(rotate64(e+f, 43)+rotate64(g, 30)+h, e+rotate64(f+a, 18)+g, mul)
}

// NOTE: s must be <= 64 bytes
func naHash64(s []byte) uint64 {
	slen := len(s)
	const seed uint64 = 81
	if slen <= 16 {
		return hashLen0to16(s)
	}
	if slen <= 32 {
		return hashLen17to32(s)
	}
	return hashLen33to64(s)
}

func naHash64WithSeed(s []byte, seed uint64) uint64 {
	return naHash64WithSeeds(s, k2, seed)
}

func naHash64WithSeeds(s []byte, seed0, seed1 uint64) uint64 {
	return hashLen16(naHash64(s)-seed0, seed1)
}
