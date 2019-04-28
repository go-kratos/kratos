// Copyright 2013 Hui Chen
// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

/* Package murmur Murmur3 32bit hash function based on
http://en.wikipedia.org/wiki/MurmurHash
*/

package murmur

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	c3 = 0x85ebca6b
	c4 = 0xc2b2ae35
	r1 = 15
	r2 = 13
	m  = 5
	n  = 0xe6546b64
)

var (
	// defaultSeed default murmur seed
	defaultSeed = uint32(1)
)

// Sum32 returns a hash from the provided key.
func Sum32(key string, seed ...uint32) (hash uint32) {
	if len(seed) > 0 {
		return Murmur3([]byte(key), seed[0])
	}

	return Murmur3([]byte(key))
}

// Murmur3 murmur []byte Hash32
func Murmur3(key []byte, seed ...uint32) (hash uint32) {
	hash = defaultSeed
	if len(seed) > 0 {
		hash = seed[0]
	}

	iByte := 0
	for ; iByte+4 <= len(key); iByte += 4 {
		k := uint32(key[iByte]) | uint32(key[iByte+1])<<8 |
			uint32(key[iByte+2])<<16 | uint32(key[iByte+3])<<24

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	var remainingBytes uint32
	switch len(key) - iByte {
	case 3:
		remainingBytes += uint32(key[iByte+2]) << 16
		fallthrough
	case 2:
		remainingBytes += uint32(key[iByte+1]) << 8
		fallthrough
	case 1:
		remainingBytes += uint32(key[iByte])
		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes = remainingBytes * c2
		hash ^= remainingBytes
	}

	hash ^= uint32(len(key))
	hash ^= hash >> 16
	hash *= c3
	hash ^= hash >> 13
	hash *= c4
	hash ^= hash >> 16

	return
}
