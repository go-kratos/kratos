package sidebar

import (
	"go-common/library/time"
)

// SideBar struct
type SideBar struct {
	ID         int64     `json:"id,omitempty"`
	Tip        int       `json:"tip,omitempty"`
	Rank       int       `json:"rank,omitempty"`
	Logo       string    `json:"logo,omitempty"`
	LogoWhite  string    `json:"logo_white,omitempty"`
	Name       string    `json:"name,omitempty"`
	Param      string    `json:"param,omitempty"`
	Module     int       `json:"module,omitempty"`
	Plat       int8      `json:"-"`
	Build      int       `json:"-"`
	Conditions string    `json:"-"`
	OnlineTime time.Time `json:"online_time"`
}

// Limit struct
type Limit struct {
	ID        int64  `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}
