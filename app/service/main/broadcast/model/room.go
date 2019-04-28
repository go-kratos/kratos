package model

import (
	"fmt"
	"net/url"
)

const (
	// NoRoom default no room key
	NoRoom = "noroom"
)

// EncodeRoomKey encode a room key.
func EncodeRoomKey(business string, room string) string {
	return fmt.Sprintf("%s://%s", business, room)
}

// DecodeRoomKey decode room key.
func DecodeRoomKey(key string) (string, string, error) {
	u, err := url.Parse(key)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}
