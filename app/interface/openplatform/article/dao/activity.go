package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// HandleActivity add or delete activity
func (d *Dao) HandleActivity(c context.Context, mid, aid, actID int64, state int, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("state", strconv.Itoa(state)) //-1 0-待审 1
	params.Set("type", strconv.Itoa(12))
	var res struct {
		Code int `json:"code"`
	}
	log.Info("HandleActivity url(%s)", d.c.Article.ActAddURI+"?"+params.Encode())
	if err = d.httpClient.RESTfulPost(c, d.c.Article.ActAddURI, ip, params, &res, actID); err != nil {
		log.Error("activity: HandleActivity url(%s) response(%s) error(%+v)", d.c.Article.ActAddURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeActivityErr
		PromError("activity:活动绑定")
		return
	}
	if res.Code != 0 {
		log.Error("activity: HandleActivity url(%s) res(%v)", d.c.Article.ActAddURI+"?"+params.Encode(), res)
		err = ecode.CreativeActivityErr
		PromError("activity:活动绑定")
	}
	return
}

// DelActivity delete activity
func (d *Dao) DelActivity(c context.Context, aid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("otype", strconv.Itoa(12))
	params.Set("state", strconv.Itoa(-1))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, d.c.Article.ActDelURI, ip, params, &res); err != nil {
		log.Error("DelActivity url(%s) response(%s) error(%+v)", d.c.Article.ActDelURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeActivityErr
		PromError("activity:活动取消绑定")
		return
	}
	if res.Code != 0 {
		log.Error("DelActivity url(%s) res(%v)", d.c.Article.ActDelURI+"?"+params.Encode(), res)
		err = ecode.CreativeActivityErr
		PromError("activity:活动取消绑定")
	}
	return
}

// Activity .
func (d *Dao) Activity(c context.Context) (resp map[int64]*model.Activity, err error) {
	var res struct {
		Code int               `json:"errno"`
		Msg  string            `json:"msg"`
		Data []*model.Activity `json:"data"`
	}
	err = d.httpClient.Get(c, d.c.Article.ActURI, "", nil, &res)
	if err != nil {
		PromError("activity:在线活动")
		log.Error("activity: d.client.Get(%s) error(%+v)", d.c.Article.ActURI+"?", err)
		return
	}
	if res.Code != 0 {
		PromError("activity:在线活动")
		log.Error("activity: url(%s) res code(%d) msg: %s", d.c.Article.ActURI+"?", res.Code, res.Msg)
		err = ecode.Int(res.Code)
		return
	}
	for _, act := range res.Data {
		if resp == nil {
			resp = make(map[int64]*model.Activity)
		}
		resp[act.ID] = act
	}
	return
}
