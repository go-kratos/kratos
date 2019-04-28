package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
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
			s.setting = &model.Setting{}
		}
		for name, value := range settings {
			switch name {
			case "recheck_view":
				var recheckView = &model.Recheck{}
				if err = json.Unmarshal([]byte(value), recheckView); err != nil {
					log.Error("setting.Unmarshal(%s) error(%+v)", value, err)
					dao.PromError("service:配置项无效")
				} else {
					s.setting.Recheck = recheckView
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
