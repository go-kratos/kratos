package feedback

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/feedback"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_addURL    = "/x/internal/feedback/ugc/add"
	_listURL   = "/x/internal/feedback/ugc/session"
	_detailURL = "/x/internal/feedback/ugc/reply"
	_tagsURL   = "/x/internal/feedback/ugc/tag"
	_closeURL  = "/x/internal/feedback/ugc/session/close"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	addURI    string
	listURI   string
	detailURI string
	tagsURI   string
	closeURI  string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		client:    httpx.NewClient(c.HTTPClient.Slow),
		addURI:    c.Host.API + _addURL,
		listURI:   c.Host.API + _listURL,
		detailURI: c.Host.API + _detailURL,
		tagsURI:   c.Host.API + _tagsURL,
		closeURI:  c.Host.API + _closeURL,
	}
	return
}

// AddFeedback add feedback
func (d *Dao) AddFeedback(c context.Context, mid, tagID, sessionID int64, qq, content, aid, browser, imgURL, platform, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("session_id", strconv.FormatInt(sessionID, 10))
	params.Set("content", content)
	params.Set("browser", browser)
	params.Set("img_url", imgURL)
	params.Set("qq", qq)
	params.Set("tag_id", strconv.FormatInt(tagID, 10))
	params.Set("aid", aid)
	params.Set("platform", platform)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.addURI, ip, params, &res); err != nil {
		log.Error("AddFeedback url(%s) response(%+v) error(%v)", d.addURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFeedbackErr
		return
	}
	if res.Code != 0 {
		log.Error("AddFeedback url(%s) res(%v)", d.addURI+"?"+params.Encode(), res)
		err = ecode.CreativeFeedbackErr
		return
	}
	return
}

// Tags get feedback tags/types
func (d *Dao) Tags(c context.Context, mid int64, ip string) (tls *feedback.TagList, err error) {
	params := url.Values{}
	params.Set("type", "0")
	params.Set("platform", "ugc,article")
	params.Set("mid", strconv.FormatInt(mid, 10))

	var res struct {
		Code int `json:"code"`
		Data struct {
			Tags  map[string][]*feedback.Tag `json:"tags"`
			Limit int                        `json:"limit"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.tagsURI, ip, params, &res); err != nil {
		log.Error("Feedback Tags url(%s) response(%+v) error(%v)", d.tagsURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFeedbackErr
		return
	}
	if res.Code != 0 {
		log.Error("Feedback Tags url(%s) res(%v)", d.tagsURI+"?"+params.Encode(), res)
		err = ecode.CreativeFeedbackErr
		return
	}
	tls = &feedback.TagList{}
	if res.Data.Tags == nil {
		return
	}
	platforms := []string{
		"article",
		"ugc",
	}
	pls := make([]*feedback.Platform, 0, len(platforms))
	for _, key := range platforms {
		if list, ok := res.Data.Tags[key]; ok {
			pl := &feedback.Platform{}
			tgs := make([]*feedback.Tag, 0, len(res.Data.Tags[key]))
			for _, v := range list {
				tg := &feedback.Tag{}
				tg.ID = v.ID
				tg.Name = v.Name
				tgs = append(tgs, tg)
			}
			pl.Tags = tgs
			if key == "ugc" {
				pl.ZH = "视频稿件"
			}
			if key == "article" {
				pl.ZH = "专栏稿件"
			}
			pl.EN = key
			pls = append(pls, pl)
		}
	}
	tls.Platforms = pls
	tls.Limit = res.Data.Limit
	return
}

// Detail get feedback details
func (d *Dao) Detail(c context.Context, mid, sessionID int64, ip string) (data []*feedback.Reply, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("session_id", strconv.FormatInt(sessionID, 10))
	var res struct {
		Code int               `json:"code"`
		Data []*feedback.Reply `json:"data"`
	}
	if err = d.client.Get(c, d.detailURI, ip, params, &res); err != nil {
		log.Error("Feedback Detail url(%s) response(%+v) error(%v)", d.detailURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFeedbackErr
		return
	}
	if res.Code != 0 {
		log.Error("Feedback Detail url(%s) res(%v)", d.detailURI+"?"+params.Encode(), res)
		err = ecode.CreativeFeedbackErr
		return
	}
	data = res.Data
	return
}

// Feedbacks get feedback list
func (d *Dao) Feedbacks(c context.Context, mid, ps, pn, tagID int64, state, start, end, platform, ip string) (data []*feedback.Feedback, count int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("platform", platform)
	if start != "" {
		params.Set("start", start)
		params.Set("end", end)
	}
	if state != "" {
		params.Set("state", state)
	}
	if tagID != 0 {
		params.Set("tag_id", strconv.FormatInt(tagID, 10))
	}
	var res struct {
		Code  int                  `json:"code"`
		Data  []*feedback.Feedback `json:"data"`
		Count int64                `json:"total"`
	}
	if err = d.client.Get(c, d.listURI, ip, params, &res); err != nil {
		log.Error("Feedbacks url(%s) response(%+v) error(%v)", d.listURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFeedbackErr
		return
	}
	if res.Code != 0 {
		log.Error("Feedbacks url(%s) res(%v)", d.listURI+"?"+params.Encode(), res)
		err = ecode.CreativeFeedbackErr
		return
	}
	data = res.Data
	count = res.Count
	return
}

// CloseSession close feedback by user
func (d *Dao) CloseSession(c context.Context, sessionID int64, ip string) (err error) {
	params := url.Values{}
	params.Set("session_id", strconv.FormatInt(sessionID, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.closeURI, ip, params, &res); err != nil {
		log.Error("Feedback CloseSession url(%s) response(%+v) error(%v)", d.closeURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFeedbackErr
		return
	}
	if res.Code != 0 {
		log.Error("Feedback CloseSession url(%s) res(%v)", d.closeURI+"?"+params.Encode(), res)
		err = ecode.CreativeFeedbackErr
		return
	}
	return
}
