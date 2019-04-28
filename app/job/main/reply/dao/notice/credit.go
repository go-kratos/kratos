package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_urlBan        = "https://www.bilibili.com/blackroom/ban/%d"
	_urlNotice     = "https://www.bilibili.com/blackroom/notice/%d"
	_urlCreditLink = "https://www.bilibili.com/judgement/case/%d"
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

// Credit return link.
func (d *Dao) Credit(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*credit `json:"data"`
	}
	if err = d.httpClient.Get(c, d.urlCredit, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlCredit, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlCredit, params.Encode(), res.Code)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_urlCreditLink, oid)
	return
}

// Notice get blackromm notice info.
func (d *Dao) Notice(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*notice `json:"data"`
	}
	if err = d.httpClient.Get(c, d.urlNotice, "", params, &res); err != nil {
		log.Error("httpNotice(%s) error(%v)", d.urlNotice, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_urlNotice, oid)
	return
}

// Ban get blackroom ban info.
func (d *Dao) Ban(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int            `json:"code"`
		Data map[int64]*ban `json:"data"`
	}
	if err = d.httpClient.Get(c, d.urlBan, "", params, &res); err != nil {
		log.Error("httpNotice(%s) error(%v)", d.urlBan, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_urlBan, oid)
	return
}
