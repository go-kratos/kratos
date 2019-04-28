package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_userName    = "bilibilinik"
	_password    = "0989D4F0-AF9D-949F-C950-E22A3F891NIK"
	_successCode = 1

	_releaseURL       = "https://api.tapd.cn/releases?workspace_id=%s&id=%s"
	_bugPreURL        = "https://api.tapd.cn/bugs?workspace_id=%s&id=%s"
	_createBugURL     = "https://api.tapd.cn/bugs"
	_workspaceUserURL = "https://api.tapd.cn/workspaces/users?workspace_id=%s&fields=user"
)

//Bug Bug
func (d *Dao) Bug(bugURL string) (res *model.BugResponse, err error) {
	var req *http.Request
	res = &model.BugResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, bugURL, nil); err != nil {
		log.Error("Bug url(%s) res(%v) err(%v)", bugURL, res, err)
		return nil, ecode.MartheTapdErr
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		log.Error("WorkspaceUser url(%s) res(%v) err(%v)", bugURL, res, err)
		err = ecode.MartheTapdResDataErr
		return
	}

	if res.Status != _successCode {
		log.Error("WorkspaceUser url(%s) res(%v) err(%v)", bugURL, res, err)
		err = ecode.MartheTapdResDataErr
	}
	return
}

//BugPre Bug Pre
func (d *Dao) BugPre(workSpaceID, bugID string) (bug *model.Bug, err error) {
	var (
		req *http.Request
		res = &model.BugSingleResponse{}
	)

	url := fmt.Sprintf(_bugPreURL, workSpaceID, bugID)
	if req, err = d.newTapdRequest(http.MethodGet, url, nil); err != nil {
		log.Error("BugPre url(%s) res(%v) err(%v)", url, res, err)
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil || res.Data == nil {
		log.Error("BugPre url(%s) res(%v) err(%v)", url, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode || res.Data == nil {
		log.Error("BugPre url(%s) res(%v) err(%v)", url, res, err)
		err = ecode.MartheTapdResDataErr
		return
	}

	bug = res.Data.Bug
	return
}

//Release Release
func (d *Dao) Release(workSpaceID, releaseID string) (release *model.Release, err error) {
	var (
		req *http.Request
		res = &model.ReleaseResponse{}
	)

	releaseURL := fmt.Sprintf(_releaseURL, workSpaceID, releaseID)
	if req, err = d.newTapdRequest(http.MethodGet, releaseURL, nil); err != nil {
		log.Error("Release url(%s) res(%v) err(%v)", _releaseURL, res, err)
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil || res.Data == nil {
		log.Error("Release url(%s) res(%v) err(%v)", _releaseURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode || res.Data == nil {
		log.Error("Release url(%s) res(%v) err(%v)", _releaseURL, res, err)
		err = ecode.MartheTapdResDataErr
		return
	}

	release = res.Data.Release
	return
}

// ReleaseName ReleaseName
func (d *Dao) ReleaseName(workspaceID, releaseID string) (releaseName string, err error) {
	var (
		item    *memcache.Item
		conn    = d.mc.Get(context.Background())
		release *model.Release
	)
	defer conn.Close()

	if item, err = conn.Get(workspaceID + releaseID); err == nil {
		if err = json.Unmarshal(item.Value, &releaseName); err != nil {
			log.Error("Json unmarshal err(%v)", err)
		}
		return
	}

	if release, err = d.Release(workspaceID, releaseID); err != nil {
		return
	}
	releaseName = release.Name
	item = &memcache.Item{Key: workspaceID + releaseID, Object: release.Name, Flags: memcache.FlagJSON, Expiration: d.expire}
	d.cacheSave(context.Background(), item)
	return
}

//newTapdRequest new tapd request
func (d *Dao) newTapdRequest(method, url string, v interface{}) (req *http.Request, err error) {
	req, err = d.newRequest(method, url, v)
	req.SetBasicAuth(_userName, _password)
	return
}

// CreateBug Create Bug.
func (d *Dao) CreateBug(bug *model.Bug) (bugID string, err error) {
	var (
		req *http.Request
		res = &model.BugSingleResponse{}
	)

	if req, err = d.newTapdRequest(http.MethodPost, _createBugURL, bug); err != nil {
		log.Error("CreateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		log.Error("CreateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode {
		log.Error("CreateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdResDataErr
	}

	bugID = res.Data.Bug.ID
	return
}

// UpdateBug Update Bug.
func (d *Dao) UpdateBug(bug *model.UpdateBug) (err error) {
	var (
		req *http.Request
		res = &model.BugSingleResponse{}
	)

	if req, err = d.newTapdRequest(http.MethodPost, _createBugURL, bug); err != nil {
		log.Error("UpdateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		log.Error("UpdateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode {
		log.Error("UpdateBug url(%s) res(%v) err(%v)", _createBugURL, res, err)
		err = ecode.MartheTapdResDataErr
	}

	return
}

// WorkspaceUser Workspace User
func (d *Dao) WorkspaceUser(workspaceID string) (usernames []string, err error) {
	var (
		req *http.Request
		res = &model.WorkspaceUserResponse{}
	)

	workspaceUserURL := fmt.Sprintf(_workspaceUserURL, workspaceID)
	if req, err = d.newTapdRequest(http.MethodGet, workspaceUserURL, nil); err != nil {
		log.Error("WorkspaceUser url(%s) res(%v) err(%v)", workspaceUserURL, res, err)
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		log.Error("WorkspaceUser url(%s) res(%v) err(%v)", workspaceUserURL, res, err)
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode || res.Data == nil {
		log.Error("WorkspaceUser url(%s) res(%v) err(%v)", workspaceUserURL, res, err)
		err = ecode.MartheTapdResDataErr
		return
	}

	for _, workspaceUser := range res.Data {
		usernames = append(usernames, workspaceUser.UserWrapper.User)
	}
	return
}
