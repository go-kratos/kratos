package service

import (
	"encoding/hex"
)

func hexEncode(b []byte) string {
	return hex.EncodeToString(b)
}

func hexDecode(s string) (res []byte, err error) {
	return hex.DecodeString(s)
}
