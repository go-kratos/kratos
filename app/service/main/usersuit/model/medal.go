package model

import (
	"fmt"
	"math/rand"

	xtime "go-common/library/time"
)

const (
	// OwnerInstall is_activated=1.
	OwnerInstall = 1
	// OwnerUninstall is_activated=0.
	OwnerUninstall = 0
	// Level1 medal_info.level 普通.
	Level1 = int32(1)
	// Level2  medal_info.level 高级.
	Level2 = int32(2)
	// Level3 medal_info.level 稀有.
	Level3 = int32(3)
	// IsGet medal has get.
	IsGet = int32(1)
	// NotGet medal not get.
	NotGet = int32(0)
)

var medalLevel = map[int32]string{Level1: "普通勋章", Level2: "高级勋章", Level3: "稀有勋章"}

// MedalOwner struct.
type MedalOwner struct {
	ID          int64      `json:"id"`
	MID         int64      `json:"mid"`
	NID         int64      `json:"nid"`
	IsActivated int8       `json:"is_activated"`
	CTime       xtime.Time `json:"ctime"`
	MTime       xtime.Time `json:"mtime"`
}

// MedalGroup struct.
type MedalGroup struct {
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	PID      int64      `json:"pid"`
	Rank     int8       `json:"rank"`
	IsOnline int8       `josn:"is_online"`
	CTime    xtime.Time `json:"ctime"`
	MTime    xtime.Time `json:"mtime"`
}

// MedalMsg struct.
type MedalMsg struct {
	ID    int64      `json:"id"`
	MID   int64      `json:"mid"`
	NID   int64      `json:"nid"`
	CTime xtime.Time `json:"ctime"`
	MTime xtime.Time `json:"mtime"`
}

// MedalCheck struct.
type MedalCheck struct {
	Has  int32       `json:"has"`
	Info interface{} `json:"info"`
}

// Build build image and level info.
func (mi *MedalInfo) Build() {
	mi.Image = getImageURL(mi.Image)
	mi.ImageSmall = getImageURL(mi.ImageSmall)
	mi.LevelDesc = medalLevel[mi.Level]
}

// getImageUrl get image from BFS.
func getImageURL(imgSrc string) (imgURL string) {
	return fmt.Sprintf("http://i%d.hdslb.com%s", rand.Int63n(3), imgSrc)
}
