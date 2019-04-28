package model

import (
	"go-common/app/service/main/filter/conf"
	xtime "go-common/library/time"
)

// Area struct .
type Area struct {
	ID         int        `json:"id"`
	GroupID    int        `json:"group_id"`
	Name       string     `json:"name"`
	ShowName   string     `json:"show_name"`
	CommonFlag bool       `json:"common_flag"`
	Ctime      xtime.Time `json:"ctime"`
	Mtime      xtime.Time `json:"mtime"`
}

// RubbishName .
func (a *Area) RubbishName() string {
	switch a.Name {
	case "live_danmu":
		return "live_dm"
	case "bplus_xiaoxi", "bplus_xiaoxiliaotianshi":
		return "message"
	default:
		return a.Name
	}
}

// IsFullLevel 是否是全level过滤
func (a *Area) IsFullLevel() bool {
	for _, name := range conf.Conf.Property.FilterFullLevelList {
		if a.Name == name {
			return true
		}
	}
	return false
}

// IsFilter 是否过滤敏感词
func (a *Area) IsFilter() bool {
	switch a.Name {
	case "message":
		return false
	default:
		return true
	}
}

// IsFilterCommon 是否过滤基础库
func (a *Area) IsFilterCommon() bool {
	switch a.Name {
	case "common":
		return true
	default:
		return a.CommonFlag
	}
}

// IsFilterKey .
func (a *Area) IsFilterKey(keys []string) bool {
	switch a.Name {
	case "danmu":
		return true
	default:
		if len(keys) > 0 {
			return true
		}
		return false
	}
}

// IsFilterRubbish .
func (a *Area) IsFilterRubbish(oid int64) bool {
	switch a.Name {
	case "message":
		return true
	case "reply", "danmu", "live_danmu", "bplus_xiaoxi", "bplus_xiaoxiliaotianshi":
		if oid != 0 {
			return true
		}
		return false
	default:
		return false
	}
}

// IsAIFilter 是否过AI过滤.
func (a *Area) IsAIFilter() bool {
	switch a.Name {
	case "reply", "live_danmu", "danmu":
		return true
	default:
		return false
	}
}
