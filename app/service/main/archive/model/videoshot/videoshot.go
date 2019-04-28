package videoshot

import (
	"time"

	xtime "go-common/library/time"
)

var (
	_verDate = time.Date(2015, 6, 1, 0, 0, 0, 0, time.Local)
)

// Videoshot is struct.
type Videoshot struct {
	Cid     int64
	Count   int
	version int
	CTime   xtime.Time
	MTime   xtime.Time
}

// Version get version.
func (v *Videoshot) Version() int {
	if v.version > 0 {
		return v.version
	}
	return int(v.MTime.Time().Sub(_verDate) / time.Second)
}

// SetVersion set version from cache.
func (v *Videoshot) SetVersion(version int) {
	v.version = version
}
