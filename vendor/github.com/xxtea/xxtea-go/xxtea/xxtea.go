/**********************************************************\
|                                                          |
| xxtea.go                                                 |
|                                                          |
| XXTEA encryption algorithm library for Golang.           |
|                                                          |
| Encryption Algorithm Authors:                            |
|      David J. Wheeler                                    |
|      Roger M. Needham                                    |
|                                                          |
| Code Author: Ma Bingyao <mabingyao@gmail.com>            |
| LastModified: Mar 10, 2015                               |
|                                                          |
\**********************************************************/

package xxtea

import (
	"encoding/base64"
	"strings"
)

const delta = 0x9E3779B9

func toBytes(v []uint32, includeLength bool) []byte {
	length := uint32(len(v))
	n := length << 2
	if includeLength {
		m := v[length-1]
		n -= 4
		if (m < n-3) || (m > n) {
			return nil
		}
		n = m
	}
	bytes := make([]byte, n)
	for i := uint32(0); i < n; i++ {
		bytes[i] = byte(v[i>>2] >> ((i & 3) << 3))
	}
	return bytes
}

func toUint32s(bytes []byte, includeLength bool) (v []uint32) {
	length := uint32(len(bytes))
	n := length >> 2
	if length&3 != 0 {
		n++
	}
	if includeLength {
		v = make([]uint32, n+1)
		v[n] = length
	} else {
		v = make([]uint32, n)
	}
	for i := uint32(0); i < length; i++ {
		v[i>>2] |= uint32(bytes[i]) << ((i & 3) << 3)
	}
	return v
}

func mx(sum uint32, y uint32, z uint32, p uint32, e uint32, k []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[p&3^e] ^ z))
}

func fixk(k []uint32) []uint32 {
	if len(k) < 4 {
		key := make([]uint32, 4)
		copy(key, k)
		return key
	}
	return k
}

func encrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = fixk(k)
	var y, z, sum, e, p, q uint32
	z = v[n]
	sum = 0
	for q = 6 + 52/length; q > 0; q-- {
		sum += delta
		e = sum >> 2 & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += mx(sum, y, z, p, e, k)
			z = v[p]
		}
		y = v[0]
		v[n] += mx(sum, y, z, p, e, k)
		z = v[n]
	}
	return v
}

func decrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = fixk(k)
	var y, z, sum, e, p, q uint32
	y = v[0]
	q = 6 + 52/length
	for sum = q * delta; sum != 0; sum -= delta {
		e = sum >> 2 & 3
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= mx(sum, y, z, p, e, k)
			y = v[p]
		}
		z = v[n]
		v[0] -= mx(sum, y, z, p, e, k)
		y = v[0]
	}
	return v
}

// Encrypt the data with key.
// data is the bytes to be encrypted.
// key is the encrypt key. It is the same as the decrypt key.
func Encrypt(data []byte, key []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	return toBytes(encrypt(toUint32s(data, true), toUint32s(key, false)), false)
}

// Decrypt the data with key.
// data is the bytes to be decrypted.
// key is the decrypted key. It is the same as the encrypt key.
func Decrypt(data []byte, key []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	return toBytes(decrypt(toUint32s(data, false), toUint32s(key, false)), true)
}

// Encrypt the data with key.
// data is the string to be encrypted.
// key is the string of encrypt key.
func EncryptString(str, key string) string {
	s := []byte(str)
	k := []byte(key)
	b64 := base64.StdEncoding
	return b64.EncodeToString(Encrypt(s, k))
}

// Decrypt the data with key.
// data is the string to be decrypted.
// key is the decrypted key. It is the same as the encrypt key.
func DecryptString(str, key string) (string, error) {
	k := []byte(key)
	b64 := base64.StdEncoding
	decodeStr, err := b64.DecodeString(str)
	if err != nil {
		return "", err
	}
	result := Decrypt([]byte(decodeStr), k)
	return string(result), nil
}

// Encrypt the string with key and convert the string to URL format
func EncryptStdToURLString(str, key string) string {
	return encryptBase64ToUrlFormat(EncryptString(str, key))
}

// Decrypt the URL string with key and convert the URL string to the origin string
func DecryptURLToStdString(str, key string) (string, error) {
	return DecryptString(decryptBase64ToStdFormat(str), key)
}

// Replace std character to URL character in base64 string
func encryptBase64ToUrlFormat(str string) string {
	str = strings.Replace(str, "+", "-", -1)
	str = strings.Replace(str, "/", "_", -1)
	str = strings.Replace(str, "=", "~", -1)
	return str
}

// Replace URL character to origin character in base64 string
func decryptBase64ToStdFormat(str string) string {
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	str = strings.Replace(str, "~", "=", -1)
	return str
}
