package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_treeURI         = "/v1/node/tree"
	_treeMachines    = "/v1/instance/hostnames"
	_treeSon         = "/v1/node/extree"
	_treeRole        = "/v1/node/role"
	_treeAppInstance = "/v1/instance/app"
	_authURI         = "/v1/auth"
	_authPlatformURI = "/v1/token"
	_treeOkCode      = 90000
	_sessIDKey       = "_AJSESSIONID"
	_treeRootName    = "bilibili."
	_questionMark    = "?"
)

// UserTree get user tree node.
func (d *Dao) UserTree(c context.Context, sessionID string) (tree *model.UserTree, err error) {
	var (
		req *http.Request
		res = &model.TreeResponse{}
	)
	if req, err = d.newTreeRequest(c, http.MethodGet, _treeURI, sessionID, nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeURI, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", _treeURI, res)
		return
	}
	tree = &res.Data
	return
}

// TreeSon get user tree node son node.
func (d *Dao) TreeSon(c context.Context, sessionID, treePath string) (treeSon map[string]interface{}, err error) {
	var (
		req *http.Request
		res = &model.TreeSonResponse{}
	)
	if req, err = d.newTreeRequest(c, http.MethodGet, _treeSon+"/"+d.getTreeFullPath(treePath), sessionID, nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeSon, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", _treeSon, res)
		return
	}
	return res.Data, nil
}

// TreeRoles get tree roles.
func (d *Dao) TreeRoles(c context.Context, sessionID, treePath string) (treeRoles []*model.TreeRole, err error) {
	var (
		req *http.Request
		res = &model.TreeRoleResponse{}
	)
	if req, err = d.newTreeRequest(c, http.MethodGet, _treeRole+"/"+d.getTreeFullPath(treePath), sessionID, nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeRole, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", _treeRole, res)
		return
	}
	return res.Data, nil
}

// TreeRolesAsPlatform get tree roles.
func (d *Dao) TreeRolesAsPlatform(c context.Context, treePath string) (treeRoles []*model.TreeRole, err error) {
	var (
		req *http.Request
		res = &model.TreeRoleResponse{}
	)
	if req, err = d.newPlatformTreeRequest(c, http.MethodGet, _treeRole+"/"+d.getTreeFullPath(treePath), nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeRole, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", _treeRole, res)
		return
	}
	return res.Data, nil
}

//QueryTreeInstances query tree instances
func (d *Dao) QueryTreeInstances(c context.Context, sessionID string, tir *model.TreeInstanceRequest) (tid map[string]*model.TreeInstance, err error) {
	var (
		req *http.Request
		res = &model.TreeInstancesResponse{}
	)
	if req, err = d.newTreeRequest(c, http.MethodGet, _treeMachines+_questionMark+tir.ToQueryURI(), sessionID, nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeRole, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		log.Error("Status url(%s) res(%v)", _treeRole, res)
		err = ecode.MerlinTreeResponseErr
		return
	}
	tid = res.Data
	return
}

// TreeAppInstance get user tree node app instance.
func (d *Dao) TreeAppInstance(c context.Context, treePaths []string) (treeAppInstances map[string][]*model.TreeAppInstance, err error) {
	var (
		req           *http.Request
		res           = &model.TreeAppInstanceResponse{}
		fullTreePaths []string
	)

	for _, treePath := range treePaths {
		fullTreePaths = append(fullTreePaths, d.getTreeFullPath(treePath))
	}

	treeAppInstanceRequest := &model.TreeAppInstanceRequest{
		Paths: fullTreePaths,
	}

	if req, err = d.newPlatformTreeRequest(c, http.MethodPost, _treeAppInstance+"?type=container", treeAppInstanceRequest); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) err(%v)", _treeRole, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeOkCode {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", _treeRole, res)
		return
	}
	return res.Data, nil
}

func (d *Dao) treeToken(c context.Context, sessionID string) (authToken string, err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
	)

	defer conn.Close()

	if item, err = conn.Get(sessionID); err == nil {
		if err = json.Unmarshal(item.Value, &authToken); err != nil {
			log.Error("treePlatformToken json parse error(%v)", err)
		}
		return
	}
	if authToken, err = d.authTree(c, sessionID); err != nil {
		return
	}

	item = &memcache.Item{Key: sessionID, Object: authToken, Flags: memcache.FlagJSON, Expiration: d.expire}
	d.tokenCacheSave(c, item)
	return

}

func (d *Dao) authTree(c context.Context, sessionID string) (authToken string, err error) {
	var (
		req      *http.Request
		tokenURL = d.c.ServiceTree.Host + _authURI
		res      struct {
			Code    int                    `json:"code"`
			Data    map[string]interface{} `json:"data"`
			Message string                 `json:"message"`
			Status  int                    `json:"status"`
		}
	)
	if req, err = d.newRequest(http.MethodGet, tokenURL, nil); err != nil {
		return
	}
	req.Header.Set("Cookie", _sessIDKey+"="+sessionID)
	if err = d.httpClient.Do(c, req, &res); err != nil {
		err = ecode.MerlinTreeRequestErr
		log.Error("d.Token url(%s) res($s) err(%v)", tokenURL, res, err)
		return
	}
	if res.Code != _treeOkCode {
		log.Error("Status url(%s) res(%v)", tokenURL, res)
		return
	}
	authToken = res.Data["token"].(string)
	return
}

func (d *Dao) newTreeRequest(c context.Context, method, uri, sessionID string, v interface{}) (req *http.Request, err error) {
	var authToken string
	if authToken, err = d.treeToken(c, sessionID); err != nil {
		return
	}
	if req, err = d.newRequest(method, d.c.ServiceTree.Host+uri, v); err != nil {
		return
	}
	req.Header.Set(_authHeader, authToken)
	return
}

func (d *Dao) treePlatformToken(c context.Context) (authToken string, err error) {
	var (
		conn     = d.mc.Get(c)
		item     *memcache.Item
		username = d.c.ServiceTree.Key
		secret   = d.c.ServiceTree.Secret
	)
	defer conn.Close()
	if item, err = conn.Get(secret); err == nil {
		if err = json.Unmarshal(item.Value, &authToken); err != nil {
			log.Error("treePlatformToken json parse error(%v)", err)
		}
		return
	}
	if authToken, err = d.authPlatformTree(c, username, secret); err != nil {
		return
	}
	item = &memcache.Item{Key: secret, Object: authToken, Flags: memcache.FlagJSON, Expiration: d.expire}
	d.tokenCacheSave(c, item)
	return
}

func (d *Dao) authPlatformTree(c context.Context, username, platformID string) (authToken string, err error) {
	var (
		req      *http.Request
		tokenURL = d.c.ServiceTree.Host + _authPlatformURI
		res      struct {
			Code    int                    `json:"code"`
			Data    map[string]interface{} `json:"data"`
			Message string                 `json:"message"`
			Status  int                    `json:"status"`
		}
	)

	treePlatformTokenRequest := &model.TreePlatformTokenRequest{
		UserName:   username,
		PlatformID: platformID,
	}

	if req, err = d.newRequest(http.MethodPost, tokenURL, treePlatformTokenRequest); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		err = ecode.MerlinTreeRequestErr
		log.Error("d.Token url(%s) res($s) err(%v)", tokenURL, res, err)
		return
	}
	if res.Code != _treeOkCode {
		log.Error("Status url(%s) res(%v)", tokenURL, res)
		return
	}
	authToken = res.Data["token"].(string)
	return
}

func (d *Dao) newPlatformTreeRequest(c context.Context, method, uri string, v interface{}) (req *http.Request, err error) {
	var authToken string
	if authToken, err = d.treePlatformToken(c); err != nil {
		return
	}
	if req, err = d.newRequest(method, d.c.ServiceTree.Host+uri, v); err != nil {
		return
	}
	req.Header.Set(_authHeader, authToken)
	return
}

func (d *Dao) getTreeFullPath(path string) string {
	return _treeRootName + path
}
