package ut

import (
	"go-common/library/time"
	"sync"
)

// RankResp resp result of rank list
type RankResp struct {
	UserName   string    `gorm:"column:username" json:"username"`
	Score      float64   `gorm:"-" json:"score"`
	Newton     float64   `gorm:"-" json:"newton"`
	Coverage   float64   `gorm:"-" json:"coverage"`
	PassRate   float64   `gorm:"-" json:"pass_rate"`
	Assertions int       `gorm:"-" json:"assertions"`
	Passed     int       `gorm:"-" json:"passed"`
	AvatarURL  string    `gorm:"-" json:"avatar_url"`
	Mtime      time.Time `gorm:"column:mtime" json:"mtime"`
	Rank       int       `gorm:"-" json:"rank"`
	Total      int       `gorm:"-" json:"total"`
	Change     int       `gorm:"-" json:"change"`
}

// RanksCache ranks cache.
type RanksCache struct {
	Slice []*RankResp
	Map   map[string]*RankResp
	sync.Mutex
}

// Image image of gitlab
type Image struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}
