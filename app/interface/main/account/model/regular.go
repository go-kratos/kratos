package model

import (
	"bytes"
	"regexp"
)

// var .
var (
	EmojiPattern = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
	NamePattern  = regexp.MustCompile("^[A-Za-z0-9\uAC00-\uD788\u3041-\u309E\u30A1-\u30FE\u3131-\u3163\u4E00-\u9FA5\uF92C-\uFA29_\\-]+$")
)

// HasEmoji is used to check string is contain emoji
func HasEmoji(s string) bool {
	return EmojiPattern.MatchString(s)
}

// ValidName check string is contain special characters.
func ValidName(s string) bool {
	h := []byte(s)
	if bytes.Contains(h, []byte("\xF0\x9F")) || bytes.Contains(h, []byte("\xC2\xA0")) {
		return false
	}
	return NamePattern.MatchString(s) && !EmojiPattern.MatchString(s)
}
