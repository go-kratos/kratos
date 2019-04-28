package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

func (s *Service) loadBannersproc() {
	for {
		now := time.Now()
		ts := now.Unix()
		if (s.bannersMap == nil) || (ts%s.dao.UpdateBannersInterval == 0) {
			banners, err := s.banners(context.TODO(), s.c.Article.BannerIDs)
			if err != nil {
				dao.PromError("service:更新banner数据")
				time.Sleep(time.Second)
				continue
			}
			s.bannersMap = banners
		}
		// 这里不是每秒钟一更新
		time.Sleep(time.Second)
	}
}

func (s *Service) loadActBannersproc() {
	for {
		now := time.Now()
		ts := now.Unix()
		if (s.actBannersMap == nil) || (ts%s.dao.UpdateBannersInterval == 0) {
			banners, err := s.banners(context.TODO(), s.c.Article.ActBannerIDs)
			if err != nil {
				dao.PromError("service:更新actBanner数据")
				time.Sleep(time.Second)
				continue
			}
			s.actBannersMap = banners
		}
		// 这里不是每秒钟一更新
		time.Sleep(time.Second)
	}
}

// Banners get banners
func (s *Service) Banners(c context.Context, plat int8, build int, t time.Time) (res []*model.Banner, err error) {
	tStr := strconv.FormatInt((t.UnixNano() / 1e6), 10)
	for _, banner := range s.bannersMap[int8(plat)] {
		if !invalidBuild(build, banner.Build, banner.Condition) {
			b := &model.Banner{}
			*b = *banner
			b.RequestID = tStr
			res = append(res, b)
		}
	}
	return
}

func (s *Service) actBanners(c context.Context, plat int8, t time.Time) (res []*model.Banner, err error) {
	tStr := strconv.FormatInt((t.UnixNano() / 1e6), 10)
	for _, banner := range s.actBannersMap[int8(plat)] {
		b := &model.Banner{}
		*b = *banner
		b.RequestID = tStr
		res = append(res, b)
	}
	return
}

func invalidBuild(srcBuild, cfgBuild int, cfgCond string) bool {
	if cfgBuild != 0 && cfgCond != "" {
		switch cfgCond {
		case "gt":
			if cfgBuild >= srcBuild {
				return true
			}
		case "lt":
			if cfgBuild <= srcBuild {
				return true
			}
		case "eq":
			if cfgBuild != srcBuild {
				return true
			}
		case "ne":
			if cfgBuild == srcBuild {
				return true
			}
		}
	}
	return false
}

func (s *Service) banners(c context.Context, resIDs []int) (res map[int8][]*model.Banner, err error) {
	arg := &resmdl.ArgRess{ResIDs: resIDs}
	bs, err := s.resRPC.Resources(c, arg)
	if err != nil {
		dao.PromError("banner:RPC")
		log.Error("s.resRPC.Resources(%+v) err: %+v", arg, err)
		return
	}
	res = make(map[int8][]*model.Banner)
	for _, r := range bs {
		for i, a := range r.Assignments {
			b := &model.Banner{
				ID:       a.ID,
				Title:    a.Name,
				URL:      a.URL,
				Image:    a.Pic,
				Position: i + 1,
				Plat:     int8(r.Platform),
				Rule:     string(a.Rule),
				ResID:    r.ID,
			}
			if b.Rule != "" {
				var tmp *model.BannerRule
				if json.Unmarshal([]byte(b.Rule), &tmp) == nil {
					b.Build = tmp.Build
					b.Condition = tmp.Condition
				}
			}
			b.Plat = model.ConvertPlat(b.Plat)
			res[b.Plat] = append(res[b.Plat], b)
		}
	}
	return
}
