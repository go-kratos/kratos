package uuid

import (
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mrand "math/rand"
	"regexp"
	"strings"
	"time"
)

// seeded indicates if math/rand has been seeded
var seeded bool = false

// uuidRegex matches the UUID string
var uuidRegex *regexp.Regexp = regexp.MustCompile(`^\{?([a-fA-F0-9]{8})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{12})\}?$`)

// UUID type.
type UUID [16]byte

// Hex returns a hex string representation of the UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format.
func (this UUID) Hex() string {
	x := [16]byte(this)
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		x[0], x[1], x[2], x[3], x[4],
		x[5], x[6],
		x[7], x[8],
		x[9], x[10], x[11], x[12], x[13], x[14], x[15])

}

// Rand generates a new version 4 UUID.
func Rand() UUID {
	var x [16]byte
	randBytes(x[:])
	x[6] = (x[6] & 0x0F) | 0x40
	x[8] = (x[8] & 0x3F) | 0x80
	return x
}

// FromStr returns a UUID based on a string.
// The string could be in the following format:
//
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
//
// If the string is not in one of these formats, it'll return an error.
func FromStr(s string) (id UUID, err error) {
	if s == "" {
		err = errors.New("Empty string")
		return
	}

	parts := uuidRegex.FindStringSubmatch(s)
	if parts == nil {
		err = errors.New("Invalid string format")
		return
	}

	var array [16]byte
	slice, _ := hex.DecodeString(strings.Join(parts[1:], ""))
	copy(array[:], slice)
	id = array
	return
}

// MustFromStr behaves similarly to FromStr except that it'll panic instead of
// returning an error.
func MustFromStr(s string) UUID {
	id, err := FromStr(s)
	if err != nil {
		panic(err)
	}
	return id
}

// randBytes uses crypto random to get random numbers. If fails then it uses math random.
func randBytes(x []byte) {

	length := len(x)
	n, err := crand.Read(x)

	if n != length || err != nil {
		if !seeded {
			mrand.Seed(time.Now().UnixNano())
		}

		for length > 0 {
			length--
			x[length] = byte(mrand.Int31n(256))
		}
	}
}
