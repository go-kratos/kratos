package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/vip/model"
	xtime "go-common/library/time"
)

// Tips config tips.
func (s *Service) Tips(c context.Context, arg *model.ArgTips) (res []*model.TipsResp, err error) {
	var (
		level    int8
		platform int
		ctime    = xtime.Time(time.Now().AddDate(-1, 0, 0).Unix())
		ok       bool
		key      string
		tip      *model.Tips
	)
	if len(s.tips) == 0 {
		return
	}
	if platform, ok = model.PlatformByName[arg.Platform]; !ok {
		return
	}
	key = s.tipsKey(int64(platform), arg.Position)
	if len(s.tips[key]) == 0 {
		return
	}
	tmp := []*model.Tips{}
	for _, v := range s.tips[key] {
		switch v.JudgeType {
		case model.VersionMoreThan:
			if arg.Version >= v.Version {
				tmp = append(tmp, v)
			}
		case model.VersionEqual:
			if arg.Version == v.Version {
				tmp = append(tmp, v)
			}
		case model.VersionLessThan:
			if arg.Version <= v.Version {
				tmp = append(tmp, v)
			}
		}
	}
	if platform == model.DevicePC {
		tmp = s.tips[key]
	}
	if len(tmp) == 0 {
		return
	}
	switch arg.Position {
	case model.PanelPosition:
		if len(tmp) == 1 {
			res = append(res, convertTip(tmp[0]))
			return
		}
		for _, v := range tmp {
			if v.Level > level {
				level = v.Level
			}
		}
		leveltmp := []*model.Tips{}
		for _, v := range tmp {
			if v.Level == level {
				leveltmp = append(leveltmp, v)
			}
		}
		if len(leveltmp) == 1 {
			res = append(res, convertTip(leveltmp[0]))
			return
		}
		for _, v := range leveltmp {
			if v.Ctime.Time().After(ctime.Time()) {
				ctime = v.Ctime
				tip = v
			}
		}
		res = append(res, convertTip(tip))
	case model.PgcPosition:
		sort.Slice(tmp, func(i int, j int) bool {
			if tmp[i].Level == tmp[j].Level {
				return tmp[i].Ctime.Time().After(tmp[j].Ctime.Time())
			}
			return tmp[i].Level > tmp[j].Level
		})
		if len(tmp) > _tiplimit {
			tmp = tmp[:_tiplimit]
		}
		for _, v := range tmp {
			res = append(res, &model.TipsResp{
				Version:    v.Version,
				Tip:        v.Tip,
				Link:       v.Link,
				ID:         v.ID,
				ButtonName: s.c.Property.TipButtonName,
				ButtonLink: s.c.Property.TipButtonLink,
			})
		}
	}
	return
}

func convertTip(t *model.Tips) (r *model.TipsResp) {
	return &model.TipsResp{
		Version: t.Version,
		Tip:     t.Tip,
		Link:    t.Link,
		ID:      t.ID,
	}
}
