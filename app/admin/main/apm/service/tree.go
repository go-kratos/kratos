package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/tree"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_discoveryIDKey = "%s_discoveryIDKey"
)

// Appids get appid by username.
func (s *Service) Appids(c context.Context, username, cookie string) (appids []string, err error) {
	s.treeLock.RLock()
	nodem, ok := s.treeCache[username]
	tmp := make(map[string]struct{})
	if !ok {
		s.treeLock.RUnlock()
		s.TreeSync(c, username, cookie)
		s.treeLock.RLock()
		nodem = s.treeCache[username]
	}
	for _, v := range nodem {
		nameArr := strings.Split(v.Path, ".")
		newName := nameArr[1] + "." + nameArr[2] + "." + nameArr[3]
		if _, ok := tmp[newName]; !ok {
			appids = append(appids, newName)
			tmp[newName] = struct{}{}
		}
	}
	s.treeLock.RUnlock()
	return
}

// Projects get projects by username.
func (s *Service) Projects(c context.Context, username, cookie string) (projects []string, err error) {
	s.treeLock.RLock()
	nodem, ok := s.treeCache[username]
	mm := make(map[string]struct{})
	if !ok {
		s.treeLock.RUnlock()
		s.TreeSync(c, username, cookie)
		s.treeLock.RLock()
		nodem = s.treeCache[username]
	}
	for _, v := range nodem {
		nameArr := strings.Split(v.Path, ".")
		newName := nameArr[1] + "." + nameArr[2]
		if _, ok := mm[newName]; ok {
			continue
		}
		mm[newName] = struct{}{}
		projects = append(projects, newName)
	}
	s.treeLock.RUnlock()
	return
}

// TreeSync sync tree cache by username.
func (s *Service) TreeSync(c context.Context, username, cookie string) {
	nodem, err := s.roleTrees(c, username, cookie)
	if err != nil {
		log.Error("TreeSync(%s) error(%v)", username, err)
		return
	}
	s.treeLock.Lock()
	s.treeCache[username] = nodem
	s.treeLock.Unlock()
}

// trees get service trees by username.
func (s *Service) roleTrees(c context.Context, username, cookie string) (trees []*tree.Node, err error) {
	url := "http://easyst.bilibili.co/v1/auth"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", cookie)
	var result = &tree.TokenResult{}
	if err = s.client.Do(c, req, result); err != nil {
		log.Error("TreeSync(%s) error(%v)", username, err)
		err = ecode.RequestErr
		return
	}

	url = "http://easyst.bilibili.co/v1/node/role/app"
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Error("TreeSync(%s) get token error(%v)", username, err)
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", result.Data.Token)
	var dat = &tree.Resp{}
	if err = s.client.Do(c, req, dat); err != nil {
		log.Error("TreeSync(%s) token(%s) error(%v)", username, result.Data.Token, err)
		err = ecode.RequestErr
		return
	}
	if len(dat.Data) == 0 {
		log.Error("TreeSync(%s) no data", username)
		return
	}
	trees = dat.Data
	log.Info("TreeSync(%s) data(%v)", username, trees)
	return
}

// Trees tree list
func (s *Service) Trees(c context.Context, username, cookie string) (nodem []*tree.Node, err error) {
	s.treeLock.RLock()
	nodem, ok := s.treeCache[username]
	if !ok {
		s.treeLock.RUnlock()
		s.TreeSync(c, username, cookie)
		s.treeLock.RLock()
		nodem = s.treeCache[username]
	}
	s.treeLock.RUnlock()
	return
}

// AllTrees AllTrees
// func (s *Service) AllTrees(c context.Context, username, cookie string) (trees []*tree.Info, err error) {
// 	var (
// 		jsonBytes []byte
// 		url       = "http://easyst.bilibili.co/v1/token"
// 	)
// 	body := &struct {
// 		Username   string `json:"user_name"`
// 		PlatformID string `json:"platform_id"`
// 	}{
// 		Username:   "main",
// 		PlatformID: conf.Conf.Tree.PlatformID,
// 	}
// 	if jsonBytes, err = json.Marshal(body); err != nil {
// 		log.Error("json.Marshal(body) error(%v)", err)
// 		return
// 	}
// 	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
// 	if err != nil {
// 		log.Error("http.NewRequest failed", err)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	result := &struct {
// 		Code    int                     `json:"code"`
// 		Data    *model.ServiceTreeToken `json:"data"`
// 		Message string                  `json:"message"`
// 		Status  int                     `json:"status"`
// 	}{}
// 	if err = s.client.Do(c, req, result); err != nil {
// 		log.Error("TreesAll(%s) get token error(%v)", username, err)
// 		err = ecode.RequestErr
// 		return
// 	}
// 	url = "http://easyst.bilibili.co/v1/node/app/extendinfo"
// 	if req, err = http.NewRequest("GET", url, nil); err != nil {
// 		log.Error("TreesAll(%s) error(%v)", username, err)
// 		err = ecode.RequestErr
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("X-Authorization-Token", result.Data.Token)
// 	var dat = &tree.Rest{}
// 	if err = s.client.Do(c, req, dat); err != nil {
// 		log.Error("TreesAll(%s) token(%s) error(%v)", username, result.Data.Token, err)
// 		err = ecode.RequestErr
// 		return
// 	}
// 	if len(dat.Data) == 0 {
// 		log.Error("TreesAll(%s) no data", username)
// 		return
// 	}
// 	trees = dat.Data
// 	log.Info("TreesAll(%s) data(%v)", username, trees)
// 	return
// }

// DiscoveryID get appid by username.
func (s *Service) DiscoveryID(c context.Context, username, cookie string) (appids []string, err error) {
	keyName := discoveryIDKey(username)
	s.discoveryIDLock.RLock()
	nodem, ok := s.discoveryIDCache[keyName]
	if !ok || (time.Since(nodem.CTime) > 60*time.Second) {
		s.discoveryIDLock.RUnlock()
		s.DiscoveryTreeSync(c, username, cookie)
		s.discoveryIDLock.RLock()
		nodem = s.discoveryIDCache[keyName]
	}
	s.discoveryIDLock.RUnlock()
	tmp := make(map[string]struct{})
	for _, v := range nodem.Data {
		var newName string
		if v.DiscoveryID != "" {
			newName = v.DiscoveryID
		} else {
			nameArr := strings.Split(v.AppID, ".")
			newName = nameArr[0] + "." + nameArr[1] + "." + nameArr[2]
		}
		if _, ok := tmp[newName]; !ok {
			appids = append(appids, newName)
			tmp[newName] = struct{}{}
		}
	}
	return
}

// discoveryAllTrees ...
func (s *Service) discoveryAllTrees(c context.Context, username, cookie string) (dat *tree.Resd, err error) {
	url := "http://easyst.bilibili.co/v1/token"
	var jsonBytes []byte
	body := &struct {
		Username   string `json:"user_name"`
		PlatformID string `json:"platform_id"`
	}{
		Username:   "msm",
		PlatformID: conf.Conf.Tree.MsmPlatformID,
	}
	if jsonBytes, err = json.Marshal(body); err != nil {
		log.Error("json.Marshal(body) error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", cookie)
	var result = &tree.TokenResult{}
	if err = s.client.Do(c, req, result); err != nil {
		log.Error("TreeSync(%s) error(%v)", username, err)
		err = ecode.RequestErr
		return
	}
	dat = &tree.Resd{}
	url = "http://easyst.bilibili.co/v1/node/app/secretinfo/prod"
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Error("TreeSync(%s) get token error(%v)", username, err)
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", result.Data.Token)
	if err = s.client.Do(c, req, dat); err != nil {
		log.Error("TreeSync(%s) token(%s) error(%v)", username, result.Data.Token, err)
		err = ecode.RequestErr
		return
	}
	if len(dat.Data) == 0 {
		log.Error("TreeSync(%s) no data", username)
		return
	}
	dat.CTime = time.Now()
	log.Info("TreeSync(%s) data(%v)", username, dat)
	return
}

func discoveryIDKey(username string) string {
	return fmt.Sprintf(_discoveryIDKey, username)
}

// DiscoveryTreeSync sync tree cache by username discoverykey.
func (s *Service) DiscoveryTreeSync(c context.Context, username, cookie string) {
	keyName := discoveryIDKey(username)
	nodem, err := s.discoveryAllTrees(c, username, cookie)
	if err != nil {
		log.Error("DiscoveryTreeSync(%s) error(%v)", username, err)
		return
	}

	s.discoveryIDLock.Lock()
	s.discoveryIDCache[keyName] = nodem
	s.discoveryIDLock.Unlock()
}
