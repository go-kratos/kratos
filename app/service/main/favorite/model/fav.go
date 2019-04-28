package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"go-common/library/time"
)

const (
	// default name of folder
	InitFolderName = "默认收藏夹"
	// state
	StateNormal = int8(0)
	StateIsDel  = int8(1)
	// attr bit bit  from left
	AttrBitPublic      = uint(0)
	AttrBitDefault     = uint(1)
	AttrBitAudit       = uint(2)
	AttrBitAdminDelete = uint(3)
	AttrBitName        = uint(4)
	AttrBitDesc        = uint(5)
	AttrBitCover       = uint(6)
	AttrBitSensitive   = uint(7)

	AttrIsPublic  = int32(0) // 公开
	AttrIsDefault = int32(0) // 默认
	// foler attr
	AttrBitPrivate      = int32(1)
	AttrBitNoDefault    = int32(1) << AttrBitDefault
	AttrBitNeedAudit    = int32(1) << AttrBitAudit
	AttrBitHitSensitive = int32(1) << AttrBitSensitive

	AttrDefaultPublic   = 0                                 // binary 0 / int 0
	AttrDefaultNoPublic = AttrBitPrivate                    // binary 01 / int 1
	AttrNormalPublic    = AttrBitNoDefault                  // binary 10 / int 2
	AttrNormalNoPublic  = AttrBitNoDefault | AttrBitPrivate // binary 11 / int 3
	// limit
	DefaultFolderLimit = 50000
	NormalFolderLimit  = 999
	// cache
	CacheNotFound = -1
	// max type
	TypeMax = 20
	// sort field
	SortPubtime = "pubtime"
	SortMtime   = "mtime"
	SortView    = "view"
)

func (r *Resource) ResourceID() int64 {
	return r.Oid*100 + int64(r.Typ)
}

func IsMediaList(typ int32) bool {
	return typ == int32(TypeVideo) || typ == int32(TypeMusicNew)
}

type Favorite struct {
	ID       int64     `json:"id"`
	Oid      int64     `json:"oid"`
	Mid      int64     `json:"mid"`
	Fid      int64     `json:"fid"`
	Type     int8      `json:"type"`
	State    int8      `json:"state"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
	Sequence uint64    `json:"sequence"`
}

func (f *Favorite) ResourceID() int64 {
	return int64(f.Oid)*100 + int64(f.Type)
}

type Favorites struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Count int `json:"count"`
	} `json:"page"`
	List []*Favorite `json:"list"`
}

type User struct {
	ID    int64     `json:"id"`
	Oid   int64     `json:"oid"`
	Mid   int64     `json:"mid"`
	Type  int8      `json:"type"`
	State int8      `json:"state"`
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"mtime"`
}
type UserList struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
	List []*User `json:"list"`
}

// AttrVal get attr val by bit.
func (f *Folder) AttrVal(bit uint) int32 {
	return (f.Attr >> bit) & int32(1)
}

// AttrSet set attr value by bit.
func (f *Folder) AttrSet(v int32, bit uint) {
	f.Attr = f.Attr&(^(1 << bit)) | (v << bit)
}

// IsDefault return true if folder is default.
func (f *Folder) IsDefault() bool {
	return f.Attr&AttrBitNoDefault == int32(0)
}

// IsPublic return true if folder is public.
func (f *Folder) IsPublic() bool {
	return f.AttrVal(AttrBitPublic) == AttrIsPublic
}

// Access return true if the user has the access permission to the folder.
func (f *Folder) Access(mid int64) bool {
	return f.IsPublic() || f.Mid == mid
}

// IsLimited return true if folder count is eq or gt conf limit.
func (f *Folder) IsLimited(cnt int, defaultLimit int, normalLimit int) bool {
	switch f.IsDefault() {
	case true:
		return f.Count+cnt > defaultLimit
	case false:
		return f.Count+cnt > normalLimit
	}
	return true
}

func (f *Folder) MediaID() int64 {
	return f.ID*100 + f.Mid%100
}

func CheckArg(tp int8, oid int64) error {
	if tp <= 0 || oid <= 0 {
		return errors.New("negative number and zero not allowed")
	}
	return CheckType(tp)
}

func CheckType(typ int8) error {
	if typ < Article || typ > TypeMax {
		return errors.New("type code out of range")
	}
	return nil
}

// CompleteURL adds host on path.
func CompleteURL(path string) (url string) {
	if path == "" {
		// url = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	url = path
	if strings.Index(path, "//") == 0 || strings.Index(path, "http://") == 0 || strings.Index(path, "https://") == 0 {
		return
	}
	url = "https://i0.hdslb.com" + url
	return
}

// CleanURL cuts host.
func CleanURL(url string) (path string) {
	path = url
	if strings.Index(url, "//") == 0 {
		path = url[14:]
	} else if strings.Index(url, "http://") == 0 {
		path = url[19:]
	} else if strings.Index(url, "https://") == 0 {
		path = url[20:]
	}
	return
}

// Folders .
type Folders []*Folder

func (f Folders) Len() int { return len(f) }

func (f Folders) Less(i, j int) bool {
	if f[i].IsDefault() {
		return true
	}
	if f[j].IsDefault() {
		return false
	}
	if f[i].ID > f[j].ID {
		return true
	}
	return false
}

func (f Folders) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// FolderSort folder index.
type FolderSort struct {
	ID    int64              `json:"id"`
	Type  int8               `json:"type"`
	Mid   int64              `json:"mid"`
	Sort  []int64            `json:"sort"`
	Map   map[int64]struct{} `json:"-"`
	CTime time.Time          `json:"ctime"`
	MTime time.Time          `json:"mtime"`
}

// Index return index for fids.
func (f *FolderSort) Index() []byte {
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
func (f *FolderSort) SetIndex(b []byte) (err error) {
	var (
		i   int
		id  int64
		n   = len(b)
		ids = make([]int64, n/8)
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

// ToBytes return []byte for ids.
func ToBytes(ids []int64) []byte {
	var (
		i int
		v int64
		n = len(ids) * 8
		b = make([]byte, n)
	)
	for i = 0; i < n; i += 8 {
		v = ids[i/8]
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

// ToInt64s bytes to int64s.
func ToInt64s(b []byte) (ids []int64, err error) {
	var (
		i  int
		id int64
		n  = len(b)
	)
	ids = make([]int64, n/8)
	if len(b)%8 != 0 {
		err = fmt.Errorf("invalid bytes:%v", b)
		return
	}
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
	}
	return
}

// SortFolders sort the favorites by index.
func (f *FolderSort) SortFolders(fs map[int64]*Folder, isSelf bool) (res []*Folder, update bool) {
	var (
		ok     bool
		id     int64
		sorted []int64
		fav    *Folder
		idx    = f.Sort
	)
	res = make([]*Folder, 0, len(fs))
	if len(f.Sort) == 0 {
		for _, fav = range fs {
			if !isSelf && !fav.IsPublic() {
				continue
			}
			res = append(res, fav)
		}
		sort.Sort(Folders(res))
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
