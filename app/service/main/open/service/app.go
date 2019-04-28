package service

import (
	"context"
	"time"

	"go-common/app/admin/main/open/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Secret .
func (s *Service) Secret(c context.Context, sappKey string) (res string, err error) {
	ok := false
	if res, ok = s.appsecrets[sappKey]; !ok {
		log.Error("appkey(%s) not found in cache", sappKey)
		err = ecode.NothingFound
	}
	return
}

// AppID .
func (s *Service) AppID(c context.Context, appKey string) (appID int64, err error) {
	if appKey == "" {
		log.Error("appkey is null")
		err = ecode.RequestErr
		return
	}
	var ok bool
	if appID, ok = s.appIDs[appKey]; !ok {
		log.Error("failed to get appid of appkey(%v) in local cache")
		err = ecode.NothingFound
	}
	return
}

// loadApp .
func (s *Service) loadAppSecrets() {
	var (
		apps       []*model.App
		appsecrets = map[string]string{}
	)
	if err := s.dao.DB.Table("dm_apps").Find(&apps).Error; err != nil {
		log.Error("s.dao.DB error (%v)", err)
		if len(s.appsecrets) == 0 {
			if s.appsecrets, err = s.dao.AppkeyCache(context.Background()); err != nil {
				log.Error("s.dao.AppkeyCache error (%v)", err)
			}
		}
		return
	}
	for _, v := range apps {
		appsecrets[v.AppKey] = v.AppSecret
		s.appIDs[v.AppKey] = v.AppID
	}
	if length := len(appsecrets); length != 0 {
		s.appsecrets = appsecrets
		log.Info("loadAppSecrets refresh success! lines:%d", length)
		//update memcache
		if err := s.dao.SetAppkeyCache(context.Background(), appsecrets); err != nil {
			log.Error("Refresh data failed in memcache! error(%v)", err)
		}
	}
}

// loadAppSecretsproc .
func (s *Service) loadAppSecretsproc() {
	var duration time.Duration
	if duration = time.Duration(s.c.AppTicker); duration == 0 {
		//default value
		duration = time.Duration(5 * time.Minute)
	}
	ticker := time.NewTicker(duration)
	for range ticker.C {
		s.loadAppSecrets()
	}
	ticker.Stop()
}
