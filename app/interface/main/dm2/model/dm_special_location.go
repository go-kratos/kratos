package model

import (
	"strings"
)

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
