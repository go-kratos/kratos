package model

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

// TODO move to model

const (
	shortUrlLength = 6
)

var (
	chars = [62]string{
		"a", "b", "c", "d", "e", "f", "g", "h",
		"i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
	}
)

// generate short url from long url
func Generate(long string) [4]string {
	var resUrl [4]string
	h := md5.New()
	h.Write([]byte(long))
	hexstr := hex.EncodeToString(h.Sum(nil))
	for i := 0; i < 4; i++ {
		start := i * 8
		end := start + 8
		s := hexstr[start:end]
		hexInt, _ := strconv.ParseInt(s, 16, 64)
		hexInt = 0x3FFFFFFF & hexInt
		var out string = ""
		for n := 0; n < shortUrlLength; n++ {
			index := 0x0000003D & hexInt
			out += chars[index]
			hexInt = hexInt >> 5
		}
		resUrl[i] = out
	}
	return resUrl
}
