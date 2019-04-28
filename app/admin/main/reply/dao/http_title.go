package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// api
	_apiNotice       = "http://api.bilibili.co/x/internal/credit/publish/infos"
	_apiBan          = "http://api.bilibili.co/x/internal/credit/blocked/infos"
	_apiCredit       = "http://api.bilibili.co/x/internal/credit/blocked/cases"
	_apiLiveVideo    = "http://api.vc.bilibili.co/clip/v1/video/detail"
	_apiLiveActivity = "http://api.live.bilibili.co/comment/v1/relation/get_by_id"
	_apiLiveNotice   = "http://api.vc.bilibili.co/news/v1/notice/info"
	_apiLivePicture  = "http://api.vc.bilibili.co/link_draw/v1/doc/detail"
	_apiActivitySub  = "http://matsuri.bilibili.co//activity/subject/url"
	_apiTopic        = "http://matsuri.bilibili.co/activity/page/one/%d"
	_apiTopics       = "http://matsuri.bilibili.co/activity/pages"
	_apiDynamic      = "http://api.vc.bilibili.co/dynamic_repost/v0/dynamic_repost/ftch_rp_cont?dynamic_ids[]=%d"
	// link
	_linkBan         = "https://www.bilibili.com/blackroom/ban/%d"
	_linkNotice      = "https://www.bilibili.com/blackroom/notice/%d"
	_linkCredit      = "https://www.bilibili.com/judgement/case/%d"
	_linkLiveVideo   = "http://vc.bilibili.com/video/%d"
	_linkLiveNotice  = "http://link.bilibili.com/p/eden/news#/newsdetail?id=%d"
	_linkLivePicture = "http://h.bilibili.com/ywh/%d"
	_linkDynamic     = "http://t.bilibili.com/%d"
)

type notice struct {
	Title string `json:"title"`
}

type ban struct {
	Title string `json:"punishTitle"`
}
type credit struct {
	Title string `json:"punishTitle"`
}

// NoticeTitle get blackromm notice info.
func (d *Dao) NoticeTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*notice `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiNotice, "", params, &res); err != nil {
		log.Error("httpNotice(%s) error(%v)", _apiNotice, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkNotice, oid)
	return
}

// BanTitle get ban info.
func (d *Dao) BanTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int            `json:"code"`
		Data map[int64]*ban `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiBan, "", params, &res); err != nil {
		log.Error("httpBan(%s) error(%v)", _apiBan, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkBan, oid)
	return
}

// CreditTitle return link.
func (d *Dao) CreditTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*credit `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiCredit, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", _apiCredit, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiCredit, params.Encode(), res.Code)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkCredit, oid)
	return
}

// LiveVideoTitle get live video title.
func (d *Dao) LiveVideoTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("video_id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item *struct {
				Description string `json:"description"`
			} `json:"item"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveVideo, "", params, &res); err != nil {
		log.Error("httpLiveVideoTitle(%s?%s) error(%v)", _apiLiveVideo, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveVideo, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Description
	link = fmt.Sprintf(_linkLiveVideo, oid)
	return
}

// LiveActivityTitle get live activity info.
func (d *Dao) LiveActivityTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveActivity, "", params, &res); err != nil {
		log.Error("httpLiveActivityTitle(%s?%s) error(%v)", _apiLiveActivity, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveActivity, params.Encode(), res.Code)
		return
	}
	title = res.Data.Name
	link = res.Data.URL
	return
}

// LiveNoticeTitle get live notice info.
func (d *Dao) LiveNoticeTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveNotice, "", params, &res); err != nil {
		log.Error("LiveNoticeTitle(%s?%s) error(%v)", _apiLiveNotice, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveNotice, params.Encode(), res.Code)
		return
	}
	title = res.Data.Title
	link = fmt.Sprintf(_linkLiveNotice, oid)
	return
}

// LivePictureTitle get live picture info.
func (d *Dao) LivePictureTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("doc_id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item *struct {
				Title string `json:"title"`
				Desc  string `json:"description"`
			} `json:"item"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLivePicture, "", params, &res); err != nil {
		log.Error("LivePictureTitle(%s?%s) error(%v)", _apiLivePicture, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLivePicture, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Title
	if title == "" {
		title = res.Data.Item.Desc
	}
	link = fmt.Sprintf(_linkLivePicture, oid)
	return
}

// TopicTitle get topic info.
func (d *Dao) TopicTitle(c context.Context, oid int64) (title, link string, err error) {
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title  string `json:"name"`
			PCLink string `json:"pc_url"`
			H5Link string `json:"h5_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, fmt.Sprintf(_apiTopic, oid), "", nil, &res); err != nil {
		log.Error("TopicTitle(%s) error(%v)", fmt.Sprintf(_apiTopic, oid), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", fmt.Sprintf(_apiTopic, oid), res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.PCLink
	if link == "" {
		link = res.Data.H5Link
	}
	return
}

// TopicsLink get topic info.
func (d *Dao) TopicsLink(c context.Context, links map[int64]string, isTopic bool) (err error) {
	if len(links) == 0 {
		return
	}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			List []struct {
				ID    int64  `json:"id"`
				PCURL string `json:"pc_url"`
				H5URL string `json:"h5_url"`
			} `json:"list"`
		} `json:"data"`
	}
	var ids []int64
	for oid := range links {
		ids = append(ids, oid)
	}
	params := url.Values{}
	params.Set("pids", xstr.JoinInts(ids))
	params.Set("all", "isOne")
	if isTopic {
		params.Set("mold", "1")
	}
	if err = d.httpClient.Get(c, _apiTopics, "", params, &res); err != nil {
		log.Error("TopicsTitle(%s) error(%v)", _apiTopics, err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", _apiTopics, res.Code)
		return
	}
	for _, data := range res.Data.List {
		if data.PCURL != "" {
			links[data.ID] = data.PCURL
		} else {
			links[data.ID] = data.H5URL
		}
	}
	return
}

// ActivitySub return activity sub link.
func (d *Dao) ActivitySub(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"name"`
			Link  string `json:"act_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiActivitySub, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", _apiActivitySub, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", _apiActivitySub, res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.Link
	return
}

// DynamicTitle return link and title.
func (d *Dao) DynamicTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(_apiDynamic, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Pairs []struct {
				DynamicID int64  `json:"dynamic_id"`
				Content   string `json:"rp_cont"`
				Type      int32  `json:"type"`
			} `json:"pairs"`
			TotalCount int64 `json:"total_count"`
		} `json:"data,omitempty"`
		Message string `json:"message"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}

	if res.Code != 0 || res.Data == nil || len(res.Data.Pairs) == 0 {
		err = fmt.Errorf("get dynamic failed!url:%s?%s code:%d message:%s pairs:%v", uri, params.Encode(), res.Code, res.Message, res.Data.Pairs)
		return
	}
	title = res.Data.Pairs[0].Content
	link = fmt.Sprintf(_linkDynamic, oid)
	return
}
