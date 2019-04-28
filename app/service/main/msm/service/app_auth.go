package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/main/msm/model"
	"go-common/library/log"
)

const (
	_hmacTreeID = "%d%d"
	_msmTreeID  = int64(2888)
)

func signature(appID, serviceID int64, appAuth string) string {
	hmacTreeIDStr := fmt.Sprintf(_hmacTreeID, appID, serviceID)
	mac := hmac.New(sha1.New, []byte(appAuth))
	mac.Write([]byte(hmacTreeIDStr))
	return hex.EncodeToString(mac.Sum(nil))
}

// ServiceScopes ServiceScopes.
func (s *Service) ServiceScopes(c context.Context, appTreeID int64) (res map[int64]*model.Scope, err error) {
	res, ok := s.scopeMap[appTreeID]
	if !ok {
		res = make(map[int64]*model.Scope)
	}
	return
}

func (s *Service) updateScope() (err error) {
	var (
		c           = context.TODO()
		appAuthMap  map[int64]map[int64]*model.AppAuth
		appTokenMap map[int64]*model.AppToken
	)
	scopeMap := make(map[int64]map[int64]*model.Scope)
	if appAuthMap, err = s.dao.AllAppsAuth(c); err != nil {
		log.Error("The update scope process was abnormal, was blocked in DB!")
		return
	}
	if appTokenMap, err = s.dao.TreeAppInfo(c); err != nil {
		log.Error("The update scope process was abnormal, was blocked in service-tree service!")
		return
	}
	for sTreeID, sAuthMap := range appAuthMap {
		_, ok := appTokenMap[sTreeID]
		if !ok {
			log.Warn("This app(%d) has no app_auth records in the service-tree service.", sTreeID)
			continue
		}
		si := make(map[int64]*model.Scope)
		for appID, appAuth := range sAuthMap {
			appToken, b := appTokenMap[appID]
			if !b {
				log.Warn("This app(%d) has no app_auth records in the service-tree service.", appID)
				continue
			}
			scope := &model.Scope{
				AppTreeID:   appID,
				RPCMethods:  strings.Split(appAuth.RPCMethod, ","),
				HTTPMethods: strings.Split(appAuth.HTTPMethod, ","),
				Quota:       appAuth.Quota,
				Sign:        signature(appID, sTreeID, appToken.AppAuth),
			}
			si[appID] = scope
		}
		scopeMap[sTreeID] = si
	}
	s.scopeMap = scopeMap
	return
}

// updateScopeproc update scope info proc.
func (s *Service) updateScopeproc() {
	for {
		time.Sleep(time.Minute)
		s.updateScope()
	}
}

// updateMsmScopeproc update scope info proc.
func (s *Service) updateMsmScopeproc() {
	for {
		time.Sleep(time.Minute)
		s.updateMsmScope()
	}
}

func (s *Service) updateMsmScope() (err error) {
	var (
		c           = context.TODO()
		appTokenMap map[int64]*model.AppToken
	)
	scopes := make(map[int64]*model.Scope)
	if appTokenMap, err = s.dao.TreeAppInfo(c); err != nil {
		log.Error("The update scope process was abnormal, was blocked in service-tree service!")
		return
	}
	for treeID, appToken := range appTokenMap {
		scope := &model.Scope{
			AppTreeID: treeID,
			Sign:      signature(treeID, _msmTreeID, appToken.AppAuth),
		}
		scopes[treeID] = scope
	}
	s.msmScope = scopes
	return
}

// CheckSign CheckSign.
func (s *Service) CheckSign(appID int64, sign string) bool {
	scope, ok := s.msmScope[appID]
	if !ok || scope == nil {
		log.Error("CheckSign(%d) error(no msmTree info.)", appID)
		return false
	}
	return sign == scope.Sign
}
