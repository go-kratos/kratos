package model

import (
	"errors"
	"go-common/library/time"
)

var (
	// ErrFavVideoExist error video has been favoured.
	ErrFavVideoExist = errors.New("error video has been favoured")
	// ErrFavVideoAlreadyDel error video has been unfavoured.
	ErrFavVideoAlreadyDel = errors.New("error video has been unfavoured")
)

const (
	bit1 = int8(1)
	bit2 = int8(1) << 1

	// StateDefaultPublic default public folder.
	StateDefaultPublic = int8(0) // binary 00 / int 0
	// StateDefaultNoPublic default private folder.
	StateDefaultNoPublic = int8(0) | bit1 // binary 01 / int 1
	// StateNormalPublic nomal public folder.
	StateNormalPublic = bit2 | int8(0) // binary 10 / int 2
	// StateNormalNoPublic nomal private folder.
	StateNormalNoPublic = bit2 | bit1 // binary 11 / int 3

	// DefaultFolderName name of favorite folder.
	DefaultFolderName = "默认收藏夹"
)

// Favorite .
type Favorite struct {
	Fid        int64     `json:"fid"`
	Mid        int64     `json:"mid"`
	Name       string    `json:"name"`
	MaxCount   int       `json:"max_count"`
	CurCount   int       `json:"cur_count"`
	AttenCount int       `json:"atten_count"`
	State      int8      `json:"state"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"-"`
	Cover      []*Cover  `json:"cover,omitempty"`
}

// Archive .
type Archive struct {
	ID    int64     `json:"id"`
	Mid   int64     `json:"mid"`
	Fid   int64     `json:"fid"`
	Aid   int64     `json:"aid"`
	CTime time.Time `json:"-"`
	MTime time.Time `json:"-"`
}

// IsPublic return true if folder is public.
func (f *Favorite) IsPublic() bool {
	return f.State&bit1 == int8(0)
}

// IsDefault return true if folder is default.
func (f *Favorite) IsDefault() bool {
	return f.State&bit2 == int8(0)
}

// StatePub return folder's public state.
func (f *Favorite) StatePub() int8 {
	return f.State & bit1
}

// StateDef return folder's default state.
func (f *Favorite) StateDef() int8 {
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

// Favorites .
type Favorites []*Favorite

func (f Favorites) Len() int { return len(f) }

func (f Favorites) Less(i, j int) bool {
	if f[i].State < f[j].State {
		return true
	}
	if f[i].State == f[j].State && f[i].MaxCount > f[j].MaxCount {
		return true
	}
	if f[i].State == f[j].State && f[i].MaxCount <= f[j].MaxCount && f[i].CTime < f[j].CTime {
		return true
	}
	return false
}

func (f Favorites) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// Cover image
type Cover struct {
	Aid int64  `json:"aid"`
	Pic string `json:"pic"`
}
