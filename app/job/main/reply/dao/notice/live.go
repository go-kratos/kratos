package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_liveSmallVideoLink = "http://vc.bilibili.com/video/%d"
	_liveNoticeLink     = "http://link.bilibili.com/p/eden/news#/newsdetail?id=%d"
	_livePictureLink    = "http://h.bilibili.com/ywh/%d"
)

// LiveSmallVideo return link.
func (d *Dao) LiveSmallVideo(c context.Context, oid int64) (title, link string, err error) {
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
	if err = d.httpClient.Get(c, d.urlLiveSmallVideo, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlLiveSmallVideo, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlLiveSmallVideo, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Description
	link = fmt.Sprintf(_liveSmallVideoLink, oid)
	return
}

// LiveActivity return link.
func (d *Dao) LiveActivity(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, d.urlLiveActivity, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlLiveActivity, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlLiveActivity, params.Encode(), res.Code)
		return
	}
	title = res.Data.Name
	link = res.Data.URL
	return
}

// LiveNotice return link.
func (d *Dao) LiveNotice(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, d.urlLiveNotice, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlLiveNotice, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlLiveNotice, params.Encode(), res.Code)
		return
	}
	title = res.Data.Title
	link = fmt.Sprintf(_liveNoticeLink, oid)
	return
}

// LivePicture return link.
func (d *Dao) LivePicture(c context.Context, oid int64) (title, link string, err error) {
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
	if err = d.httpClient.Get(c, d.urlLivePicture, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlLivePicture, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlLivePicture, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Title
	if title == "" {
		title = res.Data.Item.Desc
	}
	link = fmt.Sprintf(_livePictureLink, oid)
	return
}
