package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
)

func (s *Service) loadSettings() {
	for {
		settings, err := s.dao.Settings(context.TODO())
		if err != nil || len(settings) == 0 {
			dao.PromError("service:获取配置")
			time.Sleep(time.Second)
			continue
		}
		if s.setting == nil {
			s.setting = &artmdl.Setting{}
		}
		for name, value := range settings {
			switch name {
			case "apply_open":
				var open bool
				if open, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ApplyOpen = open
				}
			case "apply_limit":
				var limit int64
				if limit, err = strconv.ParseInt(value, 10, 64); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseInt(%v:%v) err: %+v", name, value, err)
				} else {
					s.setting.ApplyLimit = limit
				}
			case "frozen_duration":
				var duration int64
				if duration, err = strconv.ParseInt(value, 10, 64); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseInt(%v:%v) err: %+v", name, value, err)
				} else {
					s.setting.ApplyFrozenDuration = duration
				}
			case "show_rec_new_arts":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowRecommendNewArticles = show
				}
			case "show_rank_note":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowRankNote = show
				}
			case "show_app_home_rank":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowAppHomeRank = show
				}
			case "show_later_watch":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowLaterWatch = show
				}
			case "show_small_window":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowSmallWindow = show
				}
			case "hotspot":
				var show bool
				if show, err = strconv.ParseBool(value); err != nil {
					dao.PromError("service:配置项无效")
					log.Error("strconv.ParseBool(%v: %v) err: %+v", name, value, err)
				} else {
					s.setting.ShowHotspot = show
				}
			}
		}
		return
	}
}

func (s *Service) loadSettingsproc() {
	for {
		time.Sleep(time.Minute)
		s.loadSettings()
	}
}
