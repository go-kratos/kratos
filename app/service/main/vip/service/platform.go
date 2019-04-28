package service

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/log"
)

func (s *Service) loadplatformconf() (err error) {
	temp, err := s.dao.PlatformAll(context.Background())
	if err != nil || temp == nil {
		return
	}
	tempPlatConf := map[string]int64{}
	if len(temp) == 0 {
		log.Warn("s.dao.PlatformAll is empty!")
		return
	}
	for _, v := range temp {
		plat := v.ID
		key := ""
		if v.PanelType != model.PanelTypeDefault {
			key += v.PanelType
		}
		if _, ok := model.PlatformMap[v.Platform]; ok {
			key += ":" + v.Platform
		}
		if _, ok := model.MobiAPPIDMap[v.MobiApp]; ok {
			key += ":" + v.MobiApp
		}
		if _, ok := model.DeviceMap[v.Device]; ok {
			key += ":" + v.Device
		}
		tempPlatConf[key] = plat
	}
	s.pLock.Lock()
	s.platformConf = tempPlatConf
	s.pLock.Unlock()
	return
}

// GetPlatID panel_type>platform>mobi_app>device .
func (s *Service) GetPlatID(c context.Context, platform, panelType, mobiApp, device string, build int64) (platID int64) {
	key := ""
	// 兼容bug start
	if platform == "ipad" {
		log.Info("GetPlatID platform(%s)", platform)
		platform = "ios"
	}
	// 兼容bug end
	if platform == "" {
		platform = "pc"
	}
	if panelType != model.PanelTypeDefault {
		key += panelType
	}
	if _, ok := model.PlatformMap[platform]; ok {
		key += ":" + platform
	}
	// 兼容老蓝版
	if mobiApp == "iphone" && (build > 7000 && build < 8000) {
		mobiApp = "iphone_b"
	}
	if _, ok := model.MobiAPPIDMap[mobiApp]; ok {
		key += ":" + mobiApp
	}
	if _, ok := model.DeviceMap[device]; ok {
		key += ":" + device
	}
	s.pLock.RLock()
	if plat, ok := s.platformConf[key]; ok {
		platID = plat
	}
	s.pLock.RUnlock()
	log.Info("GetPlatID platform(%s),panelType(%s),mobiApp(%s),device(%s),build(%d),platID(%d)", platform, panelType, mobiApp, device, build, platID)
	return
}

// PlatformByID platform by id.
func (s *Service) PlatformByID(c context.Context, id int64) (r *model.ConfPlatform, err error) {
	return s.dao.PlatformByID(c, id)
}
