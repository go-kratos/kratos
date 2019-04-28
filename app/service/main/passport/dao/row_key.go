package dao

import (
	"strconv"
	"strings"

	"go-common/library/log"
)

const (
	_int64Max  = 0x7fffffffffffffff
	_uint32Max = 0xffffffff
)

// reverseID reverse a digital number represented in string,
// if len(id) < len, fill 0 on the right of reverse id to make reverse id len 10,
// if len(id) > len, will return empty string.
func reverseID(id string, l int) string {
	if len(id) > l {
		log.Error("len(%s) is %d, greater than the given l %d", id, len(id), l)
		return ""
	}
	// reverse id string
	runes := []rune(id)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	rid := string(runes)

	if len(id) == l {
		return rid
	}

	// fill with 0 on rid's right
	rid += strings.Repeat("0", l-len(id))
	return rid
}

func checkIDLen(id string) bool {
	return len(id) <= _maxIDLen
}

// diffTs return the last 10 digit of (int64_max - ts).
func diffTs(ts int64) string {
	i := _int64Max - ts
	s := strconv.FormatInt(i, 10)
	// during ts 0 - (int64 - now), cut the [9,19) part of s as result
	return s[9:19]
}

// diffID return the (unsigned_int32_max - id) convert to string in base 10.
// if len of the string < 10, fill 0 on the left to make len(res) equal to 10.
func diffID(id int64) string {
	i := _uint32Max - id
	s := strconv.FormatInt(i, 10)
	if len(s) == 10 {
		return s
	}
	return strings.Repeat("0", 10-len(s)) + s
}
