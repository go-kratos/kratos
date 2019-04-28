package dao

import (
	"context"
	"net/http"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_treeTokenURI   = "/v1/auth"
	_treeNodeURI    = "/v1/node/apptree"
	_ajSessionID    = "_AJSESSIONID"
	_treeAdminURI   = "/v1/node/role/"
	_treeRoleAppURI = "/v1/node/role/app"
	_treeRspCode    = 90000
)

// QueryServiceTreeToken query service tree token by sessionID
func (d *Dao) QueryServiceTreeToken(c context.Context, sessionID string) (token string, err error) {
	var (
		req      *http.Request
		tokenURL = d.c.ServiceTree.Host + _treeTokenURI
		res      struct {
			Code    int                  `json:"code"`
			Data    *model.TokenResponse `json:"data"`
			Message string               `json:"message"`
			Status  int                  `json:"status"`
		}
	)

	if req, err = d.newRequest(http.MethodGet, tokenURL, nil); err != nil {
		return
	}
	req.Header.Set("Cookie", _ajSessionID+"="+sessionID)

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) res($s) error(%v)", tokenURL, res, err)
		return
	}
	if res.Code != _treeRspCode {
		err = ecode.MelloiTreeRequestErr
		log.Error("d.Tree.Response url(%s) resCode(%s) error(%v)", tokenURL, res.Code, err)
		return
	}
	token = res.Data.Token
	return
}

// QueryUserTree query user tree by user token
func (d *Dao) QueryUserTree(c context.Context, token string) (tree *model.UserTree, err error) {
	var (
		url = d.c.ServiceTree.Host + _treeNodeURI
		req *http.Request
		res = &model.TreeResponse{}
	)

	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) error(%v)", url, err)
		err = ecode.MerlinTreeRequestErr
		return
	}
	if res.Code != _treeRspCode {
		err = ecode.MelloiTreeRequestErr
		log.Error("Get tree error(%v)", err)
		return
	}
	tree = &res.Data
	return
}

// QueryUserRoleApp query User role app
func (d *Dao) QueryUserRoleApp(c context.Context, token string) (ra []*model.RoleApp, err error) {
	var (
		url = d.c.ServiceTree.Host + _treeRoleAppURI
		req *http.Request
		res = &model.TreeRoleApp{}
	)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	req.Header.Set("Context-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) error(%v)", url, err)
		return
	}

	if res.Code != _treeRspCode {
		err = ecode.MelloiTreeRequestErr
		log.Error("get tree admin error(%v)", err)
		return
	}
	ra = res.Data
	return
}

//QueryTreeAdmin query tree admin
func (d *Dao) QueryTreeAdmin(c context.Context, path string, token string) (ta *model.TreeAdminResponse, err error) {
	var (
		url = d.c.ServiceTree.Host + _treeAdminURI + path
		req *http.Request
	)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	req.Header.Set("Context-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)

	if err = d.httpClient.Do(c, req, &ta); err != nil {
		log.Error("d.Token url(%s) error(%v)", url, err)
		return
	}
	if ta.Code != _treeRspCode {
		err = ecode.MelloiTreeRequestErr
		log.Error("get tree admin error(%v)", err)
		return
	}
	return
}
