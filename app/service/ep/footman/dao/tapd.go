package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-common/app/service/ep/footman/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_userName    = "bilibilinik"
	_password    = "0989D4F0-AF9D-949F-C950-E22A3F891NIK"
	_startPN     = 1
	_successCode = 1
)

//Iteration fetch iterations
func (d *Dao) Iteration(iterationURL string) (res *model.IterationResponse, err error) {
	var req *http.Request
	res = &model.IterationResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, iterationURL, nil); err != nil {
		return
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//Story fetch stories
func (d *Dao) Story(storyURL string) (res *model.StoryResponse, err error) {
	var req *http.Request
	res = &model.StoryResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, storyURL, nil); err != nil {
		return nil, ecode.MartheTapdErr
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//Bug Bug
func (d *Dao) Bug(bugURL string) (res *model.BugResponse, err error) {
	var req *http.Request
	res = &model.BugResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, bugURL, nil); err != nil {
		return nil, ecode.MartheTapdErr
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//BugPre Bug Pre
func (d *Dao) BugPre(workSpaceID, bugID string) (bug *model.Bug, err error) {
	var (
		req *http.Request
		res = &model.BugSingleResponse{}
	)

	url := fmt.Sprintf(model.BugPreURL, workSpaceID, bugID)
	if req, err = d.newTapdRequest(http.MethodGet, url, nil); err != nil {
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil || res.Data == nil {
		err = ecode.MartheTapdErr
		return
	}

	bug = res.Data.Bug
	return
}

//SpecStory fetch specific story
func (d *Dao) SpecStory(storyURL string) (res *model.SpecStoryResponse, err error) {
	var req *http.Request
	res = &model.SpecStoryResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, storyURL, nil); err != nil {
		return
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//StoryChange fetch story changes
func (d *Dao) StoryChange(storyChangeURL string) (res *model.StoryChangeResponse, err error) {
	var req *http.Request
	res = &model.StoryChangeResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, storyChangeURL, nil); err != nil {
		return
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//NameMap fetch story status name mapping
func (d *Dao) NameMap(nameMapURL string) (res *model.NameMapResponse, err error) {
	var req *http.Request
	res = &model.NameMapResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, nameMapURL, nil); err != nil {
		return
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//Category fetch project category
func (d *Dao) Category(categoryURL string) (res *model.CategoryResponse, err error) {
	var req *http.Request
	res = &model.CategoryResponse{}
	if req, err = d.newTapdRequest(http.MethodGet, categoryURL, nil); err != nil {
		return
	}
	err = d.httpClient.Do(context.TODO(), req, res)
	return
}

//AllIterations get all iterations by query url
func (d *Dao) AllIterations(ps int, iterationURL string) (res *model.IterationResponse, err error) {
	var (
		tempRes       *model.IterationResponse
		iterationPage = _startPN
		tempPS        = ps
		currentTime   int
	)
	url := fmt.Sprintf(iterationURL, iterationPage)
	for tempPS == ps && currentTime < d.c.Tapd.RetryTime {
		if tempRes, err = d.Iteration(url); err != nil || len(tempRes.Data) == 0 {
			currentTime = currentTime + 1
			err = nil
			continue
		}
		if res == nil {
			res = tempRes
		} else if len(tempRes.Data) > 0 {
			res.Data = append(res.Data, tempRes.Data...)
		}
		tempPS = len(tempRes.Data)
		iterationPage++
		url = fmt.Sprintf(iterationURL, iterationPage)
	}
	return
}

//AllStories get all stories by query url
func (d *Dao) AllStories(ps int, storyURL string) (res *model.StoryResponse, err error) {
	var (
		tempRes     *model.StoryResponse
		storyPage   = _startPN
		tempPS      = ps
		currentTime int
	)
	url := fmt.Sprintf(storyURL, storyPage)
	for tempPS == ps && currentTime < d.c.Tapd.RetryTime {
		if tempRes, err = d.Story(url); err != nil || len(tempRes.Data) == 0 {
			currentTime = currentTime + 1
			err = nil
			continue
		}
		if res == nil {
			res = tempRes
		} else if len(tempRes.Data) > 0 {
			res.Data = append(res.Data, tempRes.Data...)
		}
		tempPS = len(tempRes.Data)
		storyPage++
		url = fmt.Sprintf(storyURL, storyPage)
	}
	return
}

//AllStoryChanges get all story changes by query url
func (d *Dao) AllStoryChanges(ps int, storyChangeURL string) (res *model.StoryChangeResponse, err error) {
	var (
		tempRes         *model.StoryChangeResponse
		storyChangePage = _startPN
		tempPS          = ps
		currentTime     int
	)
	url := fmt.Sprintf(storyChangeURL, storyChangePage)
	for tempPS == ps && currentTime < d.c.Tapd.RetryTime {
		if tempRes, err = d.StoryChange(url); err != nil || len(tempRes.Data) == 0 {
			currentTime = currentTime + 1
			err = nil
			continue
		}
		if res == nil {
			res = tempRes
		} else if len(tempRes.Data) > 0 {
			res.Data = append(res.Data, tempRes.Data...)
		}
		tempPS = len(tempRes.Data)
		storyChangePage++
		url = fmt.Sprintf(storyChangeURL, storyChangePage)
	}
	return
}

//AllCategories get all categories of project
func (d *Dao) AllCategories(ps int, categoryURL string) (categoryMap map[string]string, err error) {
	var (
		tempRes, res *model.CategoryResponse
		categoryPage = _startPN
		tempPS       = ps
		currentTime  int
	)
	url := fmt.Sprintf(categoryURL, categoryPage)
	for tempPS == ps && currentTime < d.c.Tapd.RetryTime {
		if tempRes, err = d.Category(url); err != nil {
			currentTime = currentTime + 1
			err = nil
			continue
		}
		if res == nil {
			res = tempRes
		} else if len(tempRes.Data) > 0 {
			res.Data = append(res.Data, tempRes.Data...)
		}
		tempPS = len(tempRes.Data)
		categoryPage++
		url = fmt.Sprintf(categoryURL, categoryPage)
	}

	categoryMap = make(map[string]string)
	for _, category := range res.Data {
		categoryMap[category.Category.ID] = category.Category.Name
	}
	return
}

//AllBugs All Bugs
func (d *Dao) AllBugs(ps int, bugURL string) (res *model.BugResponse, err error) {
	var (
		tempRes     *model.BugResponse
		bugPage     = _startPN
		tempPS      = ps
		currentTime int
	)
	url := fmt.Sprintf(bugURL, bugPage)
	for tempPS == ps && currentTime < d.c.Tapd.RetryTime {
		if tempRes, err = d.Bug(url); err != nil || len(tempRes.Data) == 0 {
			currentTime = currentTime + 1
			err = nil
			continue
		}
		if res == nil {
			res = tempRes
		} else if len(tempRes.Data) > 0 {
			res.Data = append(res.Data, tempRes.Data...)
		}
		tempPS = len(tempRes.Data)
		bugPage++
		url = fmt.Sprintf(bugURL, bugPage)
	}
	return
}

//CategoryPre Category Pre
func (d *Dao) CategoryPre(workSpaceID, releaseID string) (category *model.Category, err error) {
	var (
		req *http.Request
		res = &model.CategoryPreResponse{}
	)

	releaseURL := fmt.Sprintf(model.CategoryPreURL, workSpaceID, releaseID)
	if req, err = d.newTapdRequest(http.MethodGet, releaseURL, nil); err != nil {
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil || res.Data == nil {
		err = ecode.MartheTapdErr
		return
	}

	category = res.Data.Category
	return
}

// CategoryPreName Category PreName
func (d *Dao) CategoryPreName(workspaceID, categoryID string) (categoryPreName string, err error) {
	var (
		item     *memcache.Item
		conn     = d.mc.Get(context.Background())
		category *model.Category
	)
	defer conn.Close()

	if item, err = conn.Get(workspaceID + categoryID); err == nil {
		if err = json.Unmarshal(item.Value, &categoryPreName); err != nil {
			log.Error("Json unmarshal err(%v)", err)
		}
		return
	}

	if category, err = d.CategoryPre(workspaceID, categoryID); err != nil {
		return
	}
	categoryPreName = category.Name
	item = &memcache.Item{Key: workspaceID + categoryID, Object: category.Name, Flags: memcache.FlagJSON, Expiration: d.expire}
	d.cacheSave(context.Background(), item)
	return
}

//Release Release
func (d *Dao) Release(workSpaceID, releaseID string) (release *model.Release, err error) {
	var (
		req *http.Request
		res = &model.ReleaseResponse{}
	)

	releaseURL := fmt.Sprintf(model.ReleaseURL, workSpaceID, releaseID)
	if req, err = d.newTapdRequest(http.MethodGet, releaseURL, nil); err != nil {
		err = ecode.MartheTapdErr
		return
	}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil || res.Data == nil {
		err = ecode.MartheTapdErr
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

	if req, err = d.newTapdRequest(http.MethodPost, model.CreateBugURL, bug); err != nil {
		err = ecode.MartheTapdErr
		return
	}

	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode {
		err = ecode.MartheTapdErr
		return
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

	if req, err = d.newTapdRequest(http.MethodPost, model.CreateBugURL, bug); err != nil {
		err = ecode.MartheTapdErr
		return
	}

	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		err = ecode.MartheTapdErr
		return
	}

	if res.Status != _successCode {
		err = ecode.MartheTapdErr
	}

	return
}
