package service

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
)

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// MD52IntStr converts MD5 checksum bytes to big int string in base 10.
func MD52IntStr(md5d []byte) (res string) {
	b := big.NewInt(0)
	b.SetBytes(md5d)
	res = b.Text(10)
	return
}

// IntStr2Md5 converts big int string in base 10 to MD5 checksum bytes.
func IntStr2Md5(intStr string) (res []byte) {
	b := big.NewInt(0)
	b.SetString(intStr, 10)
	return b.Bytes()
}

// hexEncode
func hexEncode(b []byte) string {
	return hex.EncodeToString(b)
}

// hexDecode
func hexDecode(s string) (res []byte, err error) {
	return hex.DecodeString(s)
}
