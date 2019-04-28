package model

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	_regFmt = `.*/bfs/([\S]+)/%s.xml`
)

// DmSpecialContent .
type DmSpecialContent struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

// DmSpecial special dm bfs location
type DmSpecial struct {
	ID        int64
	Type      int32
	Oid       int64
	Locations string
}

// Split .
func (d *DmSpecial) Split() []string {
	return strings.Split(d.Locations, ",")
}

// Join .
func (d *DmSpecial) Join(s []string) {
	d.Locations = strings.Join(s, ",")
}

// Find find url if exist
func (d *DmSpecial) Find(sha1Sum string) string {
	locations := d.Split()
	reg := regexp.MustCompile(fmt.Sprintf(_regFmt, sha1Sum))
	for _, location := range locations {
		if reg.MatchString(location) {
			return location
		}
	}
	return ""
}
