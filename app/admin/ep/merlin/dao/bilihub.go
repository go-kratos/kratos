package dao

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_projectURI    = "/api/projects"
	_repositoryURI = "/api/repositories"
	_authHubURI    = "/api/platformUsers?role=2&session="
)

// AuthHub Auth Hub.
func (d *Dao) AuthHub(c context.Context, sessionID string) (err error) {
	var (
		req *http.Request
		res interface{}
	)
	url := d.c.BiliHub.Host + _authHubURI + sessionID
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, &res); err != nil {
		err = ecode.MerlinNoHubAccount
		log.Error("d.AuthHub url(%s) err(%v)", url, err)
	}
	return
}

// HubProjects Get Hub Projects.
func (d *Dao) HubProjects(c context.Context, projectName string) (projects []*model.HubProject, err error) {
	var req *http.Request

	url := d.c.BiliHub.Host + _projectURI + "?page=1&page_size=15&name=" + projectName
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, &projects); err != nil {
		log.Error("d.HubProjects url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
		return
	}

	return
}

// ProjectRepoCount Get Project RepoCount.
func (d *Dao) ProjectRepoCount(c context.Context, projectID int) (total int, err error) {

	var (
		req *http.Request
		res = &model.GetHubProjectDetailResponse{}
	)

	url := d.c.BiliHub.Host + _projectURI + "/" + strconv.Itoa(projectID)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.ProjectRepoCount url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
		return
	}

	return res.RepoCount, err
}

// ProjectRepos Get Project Repos.
func (d *Dao) ProjectRepos(c context.Context, projectID, pageNum, pageSize int, keyWord string) (projectRepos []*model.ProjectRepository, err error) {
	var req *http.Request

	host := d.c.BiliHub.Host + _repositoryURI + "?project_id=%s&page=%s&page_size=%s&q=%s"
	url := fmt.Sprintf(host, strconv.Itoa(projectID), strconv.Itoa(pageNum), strconv.Itoa(pageSize), keyWord)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, &projectRepos); err != nil {
		log.Error("d.ProjectRepos url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
	}

	return
}

// RepoTags Get Repo Tags.
func (d *Dao) RepoTags(c context.Context, repoName string) (repoTags []*model.RepositoryTagResponse, err error) {
	var req *http.Request

	host := d.c.BiliHub.Host + _repositoryURI + "/%s/tags"
	url := fmt.Sprintf(host, repoName)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, &repoTags); err != nil {
		log.Error("d.RepoTags url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
	}

	return
}

// DeleteRepoTag Delete Repo Tag.
func (d *Dao) DeleteRepoTag(c context.Context, repoName, tagName string) (err error) {
	var req *http.Request

	host := d.c.BiliHub.Host + _repositoryURI + "/%s/tags/%s"
	url := fmt.Sprintf(host, repoName, tagName)

	if req, err = d.newRequest(http.MethodDelete, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, nil); err != nil {
		log.Error("d.DeleteRepoTag url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
	}

	return
}

// DeleteRepo Delete Repo.
func (d *Dao) DeleteRepo(c context.Context, repoName string) (err error) {
	var req *http.Request

	url := d.c.BiliHub.Host + _repositoryURI + "/" + repoName

	if req, err = d.newRequest(http.MethodDelete, url, nil); err != nil {
		return
	}

	req.SetBasicAuth(d.c.BiliHub.Username, d.c.BiliHub.Password)

	if err = d.httpClient.Do(c, req, nil); err != nil {
		log.Error("d.DeleteRepo url(%s) err(%v)", url, err)
		err = ecode.MerlinHubRequestErr
	}

	return
}

// ImageTask Image Task.
func (d *Dao) ImageTask(imageTask func()) {
	d.hubCache.Do(context.Background(), func(c context.Context) {
		imageTask()
	})
}
