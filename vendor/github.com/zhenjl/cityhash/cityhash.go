/*
 * Go implementation of Google city hash (MIT license)
 * https://code.google.com/p/cityhash/
 *
 * MIT License http://www.opensource.org/licenses/mit-license.php
 *
 * I don't even want to pretend to understand the details of city hash.
 * I am only reproducing the logic in Go as faithfully as I can.
 *
 */

package cityhash

import (
	"encoding/binary"
	"unsafe"
)

/*
var (
	little bool
)

func init() {
	if IsLittleEndian() {
		little = true
	} else {
		little = false
	}
}
*/

func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

func unalignedLoad64(p []byte) (result uint64) {
	return binary.LittleEndian.Uint64(p)
	/*
		if little {
			result = binary.LittleEndian.Uint64(p)
		} else {
			result = binary.BigEndian.Uint64(p)
		}

		return result
	*/
}

func unalignedLoad32(p []byte) (result uint32) {
	return binary.LittleEndian.Uint32(p)
	/*
		if little {
			result = binary.LittleEndian.Uint32(p)
		} else {
			result = binary.BigEndian.Uint32(p)
		}

		return result
	*/
}

func bswap64(x uint64) uint64 {
	// Copied from netbsd's bswap64.c
	return ((x << 56) & 0xff00000000000000) |
		((x << 40) & 0x00ff000000000000) |
		((x << 24) & 0x0000ff0000000000) |
		((x << 8) & 0x000000ff00000000) |
		((x >> 8) & 0x00000000ff000000) |
		((x >> 24) & 0x0000000000ff0000) |
		((x >> 40) & 0x000000000000ff00) |
		((x >> 56) & 0x00000000000000ff)
}

func bswap32(x uint32) uint32 {
	// Copied from netbsd's bswap32.c
	return ((x << 24) & 0xff000000) |
		((x << 8) & 0x00ff0000) |
		((x >> 8) & 0x0000ff00) |
		((x >> 24) & 0x000000ff)
}

func uint32InExpectedOrder(x uint32) uint32 {
	/*
		if !little {
			return bswap32(x)
		}
	*/

	return x
}

func uint64InExpectedOrder(x uint64) uint64 {
	/*
		if !little {
			return bswap64(x)
		}
	*/

	return x
}

// If I understand the original code correctly, it's expecting to load either 8 or 4
// byes in little endian order. so let's just simplify it a bit since we will do that
// anyway..
// https://code.google.com/p/cityhash/source/browse/trunk/src/city.cc#112
func fetch64(p []byte) uint64 {
	return binary.LittleEndian.Uint64(p)
	//return uint64InExpectedOrder(unalignedLoad64(p))
}

func fetch32(p []byte) uint32 {
	return binary.LittleEndian.Uint32(p)
	//return uint32InExpectedOrder(unalignedLoad32(p))
}

const (
	k0 uint64 = 0xc3a5c85c97cb3127
	k1 uint64 = 0xb492b66fbe98f273
	k2 uint64 = 0x9ae16a3b2f90404f
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
)

func fmix(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func rotate64(val uint64, shift uint32) uint64 {
	// Avoid shifting by 64: doing so yields an undefined result.
	if shift != 0 {
		return ((val >> shift) | (val << (64 - shift)))
	}

	return val
}

func rotate32(val uint32, shift uint32) uint32 {
	// Avoid shifting by 32: doing so yields an undefined result.
	if shift != 0 {
		return ((val >> shift) | (val << (32 - shift)))
	}

	return val
}

func swap64(a, b *uint64) {
	*a, *b = *b, *a
}

func swap32(a, b *uint32) {
	*a, *b = *b, *a
}

func permute3(a, b, c *uint32) {
	swap32(a, b)
	swap32(a, c)
}

func mur(a, h uint32) uint32 {
	a *= c1
	a = rotate32(a, 17)
	a *= c2
	h ^= a
	h = rotate32(h, 19)

	//return h * 5 + 0xe6546b64
	z := h*5 + 0xe6546b64
	return z
}

func hash32Len13to24(s []byte, length uint32) uint32 {
	var a uint32 = fetch32(s[(length>>1)-4:])
	var b uint32 = fetch32(s[4:])
	var c uint32 = fetch32(s[length-8:])
	var d uint32 = fetch32(s[(length >> 1):])
	var e uint32 = fetch32(s)
	var f uint32 = fetch32(s[length-4:])
	var h uint32 = length

	return fmix(mur(f, mur(e, mur(d, mur(c, mur(b, mur(a, h)))))))
}

func hash32Len0to4(s []byte, length uint32) uint32 {
	var b, c uint32 = 0, 9

	tmp := s[:length]
	for _, v := range tmp {
		b = uint32(int64(b)*int64(c1) + int64(int8(v)))
		c ^= b
	}

	return fmix(mur(b, mur(length, c)))
}

func hash32Len5to12(s []byte, length uint32) uint32 {
	var a, b, c uint32 = length, length * 5, 9
	var d uint32 = b

	a += fetch32(s)
	b += fetch32(s[length-4:])
	c += fetch32(s[((length >> 1) & 4):])

	return fmix(mur(c, mur(b, mur(a, d))))
}

func CityHash32(s []byte, length uint32) uint32 {
	if length <= 4 {
		return hash32Len0to4(s, length)
	} else if length <= 12 {
		return hash32Len5to12(s, length)
	} else if length <= 24 {
		return hash32Len13to24(s, length)
	}

	// length > 24
	var h uint32 = length
	var g uint32 = c1 * length
	var f uint32 = g
	var a0 uint32 = rotate32(fetch32(s[length-4:])*c1, 17) * c2
	var a1 uint32 = rotate32(fetch32(s[length-8:])*c1, 17) * c2
	var a2 uint32 = rotate32(fetch32(s[length-16:])*c1, 17) * c2
	var a3 uint32 = rotate32(fetch32(s[length-12:])*c1, 17) * c2
	var a4 uint32 = rotate32(fetch32(s[length-20:])*c1, 17) * c2
	h ^= a0
	h = rotate32(h, 19)
	h = h*5 + 0xe6546b64
	h ^= a2
	h = rotate32(h, 19)
	h = h*5 + 0xe6546b64
	g ^= a1
	g = rotate32(g, 19)
	g = g*5 + 0xe6546b64
	g ^= a3
	g = rotate32(g, 19)
	g = g*5 + 0xe6546b64
	f += a4
	f = rotate32(f, 19)
	f = f*5 + 0xe6546b64

	var iters uint32 = (length - 1) / 20
	for {
		var a0 uint32 = rotate32(fetch32(s)*c1, 17) * c2
		var a1 uint32 = fetch32(s[4:])
		var a2 uint32 = rotate32(fetch32(s[8:])*c1, 17) * c2
		var a3 uint32 = rotate32(fetch32(s[12:])*c1, 17) * c2
		var a4 uint32 = fetch32(s[16:])
		h ^= a0
		h = rotate32(h, 18)
		h = h*5 + 0xe6546b64
		f += a1
		f = rotate32(f, 19)
		f = f * c1
		g += a2
		g = rotate32(g, 18)
		g = g*5 + 0xe6546b64
		h ^= a3 + a1
		h = rotate32(h, 19)
		h = h*5 + 0xe6546b64
		g ^= a4
		g = bswap32(g) * 5
		h += a4 * 5
		h = bswap32(h)
		f += a0
		permute3(&f, &h, &g)
		s = s[20:]

		iters--
		if iters == 0 {
			break
		}
	}

	g = rotate32(g, 11) * c1
	g = rotate32(g, 17) * c1
	f = rotate32(f, 11) * c1
	f = rotate32(f, 17) * c1
	h = rotate32(h+g, 19)
	h = h*5 + 0xe6546b64
	h = rotate32(h, 17) * c1
	h = rotate32(h+f, 19)
	h = h*5 + 0xe6546b64
	h = rotate32(h, 17) * c1
	return h
}

func shiftMix(val uint64) uint64 {
	return val ^ (val >> 47)
}

type Uint128 [2]uint64

func (this *Uint128) setLower64(l uint64) {
	this[0] = l
}

func (this *Uint128) setHigher64(h uint64) {
	this[1] = h
}

func (this Uint128) Lower64() uint64 {
	return this[0]
}

func (this Uint128) Higher64() uint64 {
	return this[1]
}

func (this Uint128) Bytes() []byte {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b, this[0])
	binary.LittleEndian.PutUint64(b[8:], this[1])
	return b
}

func hash128to64(x Uint128) uint64 {
	// Murmur-inspired hashing.
	const kMul uint64 = 0x9ddfea08eb382d69
	var a uint64 = (x.Lower64() ^ x.Higher64()) * kMul
	a ^= (a >> 47)
	var b uint64 = (x.Higher64() ^ a) * kMul
	b ^= (b >> 47)
	b *= kMul
	return b
}

func hashLen16(u, v uint64) uint64 {
	return hash128to64(Uint128{u, v})
}

func hashLen16_3(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	var a uint64 = (u ^ v) * mul
	a ^= (a >> 47)
	var b uint64 = (v ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func hashLen0to16(s []byte, length uint32) uint64 {
	if length >= 8 {
		var mul uint64 = k2 + uint64(length)*2
		var a uint64 = fetch64(s) + k2
		var b uint64 = fetch64(s[length-8:])
		var c uint64 = rotate64(b, 37)*mul + a
		var d uint64 = (rotate64(a, 25) + b) * mul
		return hashLen16_3(c, d, mul)
	}

	if length >= 4 {
		var mul uint64 = k2 + uint64(length)*2
		var a uint64 = uint64(fetch32(s))
		return hashLen16_3(uint64(length)+(a<<3), uint64(fetch32(s[length-4:])), mul)
	}

	if length > 0 {
		var a uint8 = uint8(s[0])
		var b uint8 = uint8(s[length>>1])
		var c uint8 = uint8(s[length-1])
		var y uint32 = uint32(a) + (uint32(b) << 8)
		var z uint32 = length + (uint32(c) << 2)
		return shiftMix(uint64(y)*k2^uint64(z)*k0) * k2
	}

	return k2
}

func hashLen17to32(s []byte, length uint32) uint64 {
	var mul uint64 = k2 + uint64(length)*2
	var a uint64 = fetch64(s) * k1
	var b uint64 = fetch64(s[8:])
	var c uint64 = fetch64(s[length-8:]) * mul
	var d uint64 = fetch64(s[length-16:]) * k2
	return hashLen16_3(rotate64(a+b, 43)+rotate64(c, 30)+d, a+rotate64(b+k2, 18)+c, mul)
}

func weakHashLen32WithSeeds(w, x, y, z, a, b uint64) Uint128 {
	a += w
	b = rotate64(b+a+z, 21)
	var c uint64 = a
	a += x
	a += y
	b += rotate64(a, 44)
	return Uint128{a + z, b + c}
}

func weakHashLen32WithSeeds_3(s []byte, a, b uint64) Uint128 {
	return weakHashLen32WithSeeds(fetch64(s), fetch64(s[8:]), fetch64(s[16:]), fetch64(s[24:]), a, b)
}

func hashLen33to64(s []byte, length uint32) uint64 {
	var mul uint64 = k2 + uint64(length)*2
	var a uint64 = fetch64(s) * k2
	var b uint64 = fetch64(s[8:])
	var c uint64 = fetch64(s[length-24:])
	var d uint64 = fetch64(s[length-32:])
	var e uint64 = fetch64(s[16:]) * k2
	var f uint64 = fetch64(s[24:]) * 9
	var g uint64 = fetch64(s[length-8:])
	var h uint64 = fetch64(s[length-16:]) * mul
	var u uint64 = rotate64(a+g, 43) + (rotate64(b, 30)+c)*9
	var v uint64 = ((a + g) ^ d) + f + 1
	var w uint64 = bswap64((u+v)*mul) + h
	var x uint64 = rotate64(e+f, 42) + c
	var y uint64 = (bswap64((v+w)*mul) + g) * mul
	var z uint64 = e + f + c
	a = bswap64((x+z)*mul+y) + b
	b = shiftMix((z+a)*mul+d+h) * mul
	return b + x
}

func CityHash64(s []byte, length uint32) uint64 {
	if length <= 32 {
		if length <= 16 {
			return hashLen0to16(s, length)
		} else {
			return hashLen17to32(s, length)
		}
	} else if length <= 64 {
		return hashLen33to64(s, length)
	}

	var x uint64 = fetch64(s[length-40:])
	var y uint64 = fetch64(s[length-16:]) + fetch64(s[length-56:])
	var z uint64 = hashLen16(fetch64(s[length-48:])+uint64(length), fetch64(s[length-24:]))
	var v Uint128 = weakHashLen32WithSeeds_3(s[length-64:], uint64(length), z)
	var w Uint128 = weakHashLen32WithSeeds_3(s[length-32:], y+k1, x)
	x = x*k1 + fetch64(s)

	length = (length - 1) & ^uint32(63)
	for {
		x = rotate64(x+y+v.Lower64()+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Higher64()+fetch64(s[48:]), 42) * k1
		x ^= w.Higher64()
		y += v.Lower64() + fetch64(s[40:])
		z = rotate64(z+w.Lower64(), 33) * k1
		v = weakHashLen32WithSeeds_3(s, v.Higher64()*k1, x+w.Lower64())
		w = weakHashLen32WithSeeds_3(s[32:], z+w.Higher64(), y+fetch64(s[16:]))
		swap64(&z, &x)
		s = s[64:]
		length -= 64

		if length == 0 {
			break
		}
	}

	return hashLen16(hashLen16(v.Lower64(), w.Lower64())+shiftMix(y)*k1+z, hashLen16(v.Higher64(), w.Higher64())+x)
}

func CityHash64WithSeed(s []byte, length uint32, seed uint64) uint64 {
	return CityHash64WithSeeds(s, length, k2, seed)
}

func CityHash64WithSeeds(s []byte, length uint32, seed0, seed1 uint64) uint64 {
	return hashLen16(CityHash64(s, length)-seed0, seed1)
}

func cityMurmur(s []byte, length uint32, seed Uint128) Uint128 {
	var a uint64 = seed.Lower64()
	var b uint64 = seed.Higher64()
	var c uint64 = 0
	var d uint64 = 0
	var l int32 = int32(length) - 16

	if l <= 0 { // len <= 16
		a = shiftMix(a*k1) * k1
		c = b*k1 + hashLen0to16(s, length)

		if length >= 8 {
			d = shiftMix(a + fetch64(s))
		} else {
			d = shiftMix(a + c)
		}
	} else { // len > 16
		c = hashLen16(fetch64(s[length-8:])+k1, a)
		d = hashLen16(b+uint64(length), c+fetch64(s[length-16:]))
		a += d

		for {
			a ^= shiftMix(fetch64(s)*k1) * k1
			a *= k1
			b ^= a
			c ^= shiftMix(fetch64(s[8:])*k1) * k1
			c *= k1
			d ^= c
			s = s[16:]
			l -= 16

			if l <= 0 {
				break
			}
		}
	}

	a = hashLen16(a, c)
	b = hashLen16(d, b)
	return Uint128{a ^ b, hashLen16(b, a)}
}

func CityHash128WithSeed(s []byte, length uint32, seed Uint128) Uint128 {
	if length < 128 {
		return cityMurmur(s, length, seed)
	}

	orig_length := length
	var t []byte = s

	// We expect length >= 128 to be the common case.  Keep 56 bytes of state:
	// v, w, x, y, and z.
	var v, w Uint128
	var x uint64 = seed.Lower64()
	var y uint64 = seed.Higher64()
	var z uint64 = uint64(length) * k1

	v.setLower64(rotate64(y^k1, 49)*k1 + fetch64(s))
	v.setHigher64(rotate64(v.Lower64(), 42)*k1 + fetch64(s[8:]))
	w.setLower64(rotate64(y+z, 35)*k1 + x)
	w.setHigher64(rotate64(x+fetch64(s[88:]), 53) * k1)

	// This is the same inner loop as CityHash64(), manually unrolled.
	for {
		x = rotate64(x+y+v.Lower64()+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Higher64()+fetch64(s[48:]), 42) * k1
		x ^= w.Higher64()
		y += v.Lower64() + fetch64(s[40:])
		z = rotate64(z+w.Lower64(), 33) * k1
		v = weakHashLen32WithSeeds_3(s, v.Higher64()*k1, x+w.Lower64())
		w = weakHashLen32WithSeeds_3(s[32:], z+w.Higher64(), y+fetch64(s[16:]))
		swap64(&z, &x)
		s = s[64:]
		x = rotate64(x+y+v.Lower64()+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Higher64()+fetch64(s[48:]), 42) * k1
		x ^= w.Higher64()
		y += v.Lower64() + fetch64(s[40:])
		z = rotate64(z+w.Lower64(), 33) * k1
		v = weakHashLen32WithSeeds_3(s, v.Higher64()*k1, x+w.Lower64())
		w = weakHashLen32WithSeeds_3(s[32:], z+w.Higher64(), y+fetch64(s[16:]))
		swap64(&z, &x)
		s = s[64:]
		length -= 128

		if length < 128 {
			break
		}
	}

	x += rotate64(v.Lower64()+z, 49) * k0
	y = y*k0 + rotate64(w.Higher64(), 37)
	z = z*k0 + rotate64(w.Lower64(), 27)
	w.setLower64(w.Lower64() * 9)
	v.setLower64(v.Lower64() * k0)

	// If 0 < length < 128, hash up to 4 chunks of 32 bytes each from the end of s.
	var tail_done uint32
	for tail_done = 0; tail_done < length; {
		tail_done += 32
		y = rotate64(x+y, 42)*k0 + v.Higher64()
		w.setLower64(w.Lower64() + fetch64(t[orig_length-tail_done+16:]))
		x = x*k0 + w.Lower64()
		z += w.Higher64() + fetch64(t[orig_length-tail_done:])
		w.setHigher64(w.Higher64() + v.Lower64())
		v = weakHashLen32WithSeeds_3(t[orig_length-tail_done:], v.Lower64()+z, v.Higher64())
		v.setLower64(v.Lower64() * k0)
	}

	// At this point our 56 bytes of state should contain more than
	// enough information for a strong 128-bit hash.  We use two
	// different 56-byte-to-8-byte hashes to get a 16-byte final result.
	x = hashLen16(x, v.Lower64())
	y = hashLen16(y+z, w.Lower64())

	return Uint128{hashLen16(x+v.Higher64(), w.Higher64()) + y,
		hashLen16(x+w.Higher64(), y+v.Higher64())}
}

func CityHash128(s []byte, length uint32) (result Uint128) {
	if length >= 16 {
		result = CityHash128WithSeed(s[16:length], length-16, Uint128{fetch64(s), fetch64(s[8:length]) + k0})
	} else {
		result = CityHash128WithSeed(s, length, Uint128{k0, k1})
	}

	return
}
