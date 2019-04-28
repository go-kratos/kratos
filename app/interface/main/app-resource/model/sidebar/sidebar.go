package sidebar

import (
	resource "go-common/app/service/main/resource/model"
	"go-common/library/time"
)

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
	NeedLogin  int8      `json:"-"`
	WhiteURL   string    `json:"-"`
	Language   string    `json:"-"`
}

type Limit struct {
	ID        int64  `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}

func (s *SideBar) Change(rsb *resource.SideBar) {
	s.ID = rsb.ID
	s.Tip = rsb.Tip
	s.Rank = rsb.Rank
	s.Logo = rsb.Logo
	s.LogoWhite = rsb.LogoWhite
	s.Name = rsb.Name
	s.Param = rsb.Param
	s.Module = rsb.Module
	s.Plat = rsb.Plat
	s.Build = rsb.Build
	s.Conditions = rsb.Conditions
	s.OnlineTime = rsb.OnlineTime
	s.NeedLogin = rsb.NeedLogin
	s.WhiteURL = rsb.WhiteURL
	s.Language = rsb.Language
}
