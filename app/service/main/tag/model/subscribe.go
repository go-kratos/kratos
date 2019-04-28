package model

import (
	"fmt"
	"sort"
)

// SubSort Subscription sorted by create time.
type SubSort struct {
	Subs  []*Sub
	Order int
}

func (s *SubSort) Len() int      { return len(s.Subs) }
func (s *SubSort) Swap(i, j int) { s.Subs[i], s.Subs[j] = s.Subs[j], s.Subs[i] }
func (s *SubSort) Less(i, j int) bool {
	if s.Order == -1 {
		return s.Subs[i].MTime.Time().Unix() > s.Subs[j].MTime.Time().Unix() // DESC
	}
	return s.Subs[i].MTime.Time().Unix() < s.Subs[j].MTime.Time().Unix() // ASC
}

// Sort sort by ctime.
func (s *SubSort) Sort() (res []int64) {
	sort.Sort(s)
	for _, r := range s.Subs {
		res = append(res, r.Tid)
	}
	return
}

// Index return index for fids.
func Index(tids []int64) []byte {
	var (
		i int
		v int64
		n = len(tids) * 8
		b = make([]byte, n)
	)
	for i = 0; i < n; i += 8 {
		v = tids[i/8]
		b[i] = byte(v >> 56)
		b[i+1] = byte(v >> 48)
		b[i+2] = byte(v >> 40)
		b[i+3] = byte(v >> 32)
		b[i+4] = byte(v >> 24)
		b[i+5] = byte(v >> 16)
		b[i+6] = byte(v >> 8)
		b[i+7] = byte(v)
	}
	return b
}

// SetIndex set sort fids.
func SetIndex(b []byte) (tids []int64, err error) {
	var (
		i  int
		id int64
		n  = len(b)
	)
	if len(b)%8 != 0 {
		err = fmt.Errorf("invalid sort index:%v", b)
		return
	}
	tids = make([]int64, n/8)
	for i = 0; i < n; i += 8 {
		id = int64(b[i+7]) |
			int64(b[i+6])<<8 |
			int64(b[i+5])<<16 |
			int64(b[i+4])<<24 |
			int64(b[i+3])<<32 |
			int64(b[i+2])<<40 |
			int64(b[i+1])<<48 |
			int64(b[i])<<56
		tids[i/8] = id
	}
	return
}
