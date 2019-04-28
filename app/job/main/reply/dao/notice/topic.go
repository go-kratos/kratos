package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// Topic return topic link.
func (d *Dao) Topic(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(d.urlTopic, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title  string `json:"name"`
			PCLink string `json:"pc_url"`
			H5Link string `json:"h5_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", uri, res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.PCLink
	if link == "" {
		link = res.Data.H5Link
	}
	return
}

// Activity return topic link.
func (d *Dao) Activity(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(d.urlActivity, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title  string `json:"name"`
			PCLink string `json:"pc_url"`
			H5Link string `json:"h5_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", uri, res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.PCLink
	if link == "" {
		link = res.Data.H5Link
	}
	return
}

// ActivitySub return topic link.
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
	if err = d.httpClient.Get(c, d.urlActivitySub, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlActivitySub, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", d.urlActivitySub, res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.Link
	return
}
