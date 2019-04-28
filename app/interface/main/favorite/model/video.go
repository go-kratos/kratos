package model

import (
	"fmt"
	"sort"

	xtime "go-common/library/time"
)

const (
	bit1                 = int8(1)
	bit2                 = int8(1) << 1
	StateDefaultPublic   = int8(0)        // binary 00 / int 0
	StateDefaultNoPublic = int8(0) | bit1 // binary 01 / int 1
	StateNormalPublic    = bit2 | int8(0) // binary 10 / int 2
	StateNormalNoPublic  = bit2 | bit1    // binary 11 / int 3
	// DefaultFolderName default name of favorite folder
	DefaultFolderName = "默认收藏夹"
	// AllFidFlag all folder id flag
	AllFidFlag = -1
	// CDFlag cool down flag
	CDFlag = -1
	// search error code
	SearchErrWordIllegal = -110 // 非法搜索词错误
	// clean state
	StateAllowToClean = 0
	StateCleaning     = 1
	StateCleanCD      = 2
)

type VideoFolder struct {
	Fid        int64      `json:"fid"`
	Mid        int64      `json:"mid"`
	Name       string     `json:"name"`
	MaxCount   int        `json:"max_count"`
	CurCount   int        `json:"cur_count"`
	AttenCount int        `json:"atten_count"`
	Favoured   int8       `json:"favoured"`
	State      int8       `json:"state"`
	CTime      xtime.Time `json:"ctime"`
	MTime      xtime.Time `json:"mtime"`
	Cover      []*Cover   `json:"cover,omitempty"`
}

// IsPublic return true if folder is public.
func (f *VideoFolder) IsPublic() bool {
	return f.State&bit1 == int8(0)
}

// IsDefault return true if folder is default.
func (f *VideoFolder) IsDefault() bool {
	return f.State&bit2 == int8(0)
}

// StatePub return folder's public state.
func (f *VideoFolder) StatePub() int8 {
	return f.State & bit1
}

// StateDef return folder's default state.
func (f *VideoFolder) StateDef() int8 {
	return f.State & bit2
}

// IsDefault return true if state is default state.
func IsDefault(state int8) bool {
	return (state&(int8(1)<<1) == int8(0))
}

// CheckPublic check user update public value in [0, 1].
func CheckPublic(state int8) bool {
	return state == int8(0) || state == bit1
}

type VideoFolders []*VideoFolder

func (f VideoFolders) Len() int { return len(f) }

func (f VideoFolders) Less(i, j int) bool {
	if f[i].IsDefault() {
		return true
	}
	if f[j].IsDefault() {
		return false
	}
	if f[i].Fid > f[j].Fid {
		return true
	}
	return false
}

func (f VideoFolders) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// Cover image
type Cover struct {
	Aid int64  `json:"aid"`
	Pic string `json:"pic"`
}

// VideoFolderSort folder index.
type VideoFolderSort struct {
	ID    int64              `json:"id"`
	Mid   int64              `json:"mid"`
	Sort  []int64            `json:"sort"`
	Map   map[int64]struct{} `json:"map"`
	CTime xtime.Time         `json:"ctime"`
	MTime xtime.Time         `json:"mtime"`
}

// Index return index for fids.
func (f *VideoFolderSort) Index() []byte {
	var (
		i  int
		v  int64
		fs = f.Sort
		n  = len(fs) * 8
		b  = make([]byte, n)
	)
	for i = 0; i < n; i += 8 {
		v = fs[i/8]
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
func (f *VideoFolderSort) SetIndex(b []byte) (err error) {
	var (
		i   int
		id  int64
		n   = len(b)
		ids = make([]int64, n)
	)
	if len(b)%8 != 0 {
		err = fmt.Errorf("invalid sort index:%v", b)
		return
	}
	f.Map = make(map[int64]struct{}, n)
	for i = 0; i < n; i += 8 {
		id = int64(b[i+7]) |
			int64(b[i+6])<<8 |
			int64(b[i+5])<<16 |
			int64(b[i+4])<<24 |
			int64(b[i+3])<<32 |
			int64(b[i+2])<<40 |
			int64(b[i+1])<<48 |
			int64(b[i])<<56
		ids[i/8] = id
		f.Map[id] = struct{}{}
	}
	f.Sort = ids
	return
}

// SortFavs sort the favorites by index.
func (f *VideoFolderSort) SortFavs(fs map[int64]*VideoFolder, isSelf bool) (res []*VideoFolder, update bool) {
	var (
		ok     bool
		id     int64
		sorted []int64
		fav    *VideoFolder
		idx    = f.Sort
	)
	res = make([]*VideoFolder, 0, len(fs))
	if len(f.Sort) == 0 {
		for _, fav = range fs {
			if !isSelf && !fav.IsPublic() {
				continue
			}
			res = append(res, fav)
		}
		sort.Sort(VideoFolders(res))
		return
	}
	if len(idx) != len(fs) {
		sorted = append(sorted, idx[0])
		for id = range fs {
			if _, ok = f.Map[id]; !ok {
				sorted = append(sorted, id)
			}
		}
		for _, id := range idx[1:] {
			if _, ok = fs[id]; ok {
				sorted = append(sorted, id)
			}
		}
		update = true
		f.Sort = sorted
	}
	for _, id = range f.Sort {
		if fav, ok = fs[id]; ok {
			if !isSelf && !fav.IsPublic() {
				continue
			}
			res = append(res, fav)
		}
	}
	return
}
