package model

import (
	"go-common/library/time"

	arcmdl "go-common/app/service/main/archive/api"
	favmdl "go-common/app/service/main/favorite/model"
	xtime "go-common/library/time"
)

// PlDBusType databus type playlist
const PlDBusType = "playlist"

// ArcSort struct.
type ArcSort struct {
	Aid  int64  `json:"aid"`
	Sort int64  `json:"sort"`
	Desc string `json:"desc"`
}

// Videos  add video result.
type Videos struct {
	RightAids []int64 `json:"right_aids"`
	WrongAids []int64 `json:"wrong_aids"`
}

// Playlist struct.
type Playlist struct {
	Pid int64 `json:"pid"`
	*favmdl.Folder
	Stat         *Stat          `json:"stat,omitempty"`
	Author       *arcmdl.Author `json:"owner,omitempty"`
	FavoriteTime time.Time      `json:"favorite_time,omitempty"`
	IsFavorite   bool           `json:"is_favorite"`
}

// Stat playlist stat.
type Stat struct {
	Pid   int64 `json:"pid"`
	View  int64 `json:"view"`
	Fav   int64 `json:"favorite"`
	Reply int64 `json:"reply"`
	Share int64 `json:"share"`
}

// PlStat playlist stat
type PlStat struct {
	ID    int64      `json:"id"`
	Mid   int64      `json:"mid"`
	Fid   int64      `json:"fid"`
	View  int64      `json:"view"`
	Reply int64      `json:"reply"`
	Fav   int64      `json:"favorite"`
	Share int64      `json:"share"`
	MTime xtime.Time `json:"mtime"`
}

// View arc view.
type View struct {
	*arcmdl.Arc
	Pages []*arcmdl.Page `json:"pages"`
}

// PlView playlist view struct
type PlView struct {
	*View
	PlayDesc string `json:"play_desc"`
}

// ArcList playlist archive list.
type ArcList struct {
	List []*PlView `json:"list"`
}

// ToView to view page struct.
type ToView struct {
	*Playlist
	List     []*View `json:"list"`
	Favorite bool    `json:"favorite"`
}

// SearchArc search archive struct
type SearchArc struct {
	Aid         int64  `json:"aid"`
	Title       string `json:"title"`
	Pic         string `json:"pic"`
	Duration    string `json:"duration"`
	Mid         int64  `json:"mid"`
	Author      string `json:"author"`
	Play        int64  `json:"play"`
	Review      int64  `json:"review"`
	VideoReview int64  `json:"video_review"`
	Favorites   int64  `json:"favorites"`
}
