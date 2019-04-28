package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_typeArticleFlow = "3"
)

// FlowSync 流量管理同步过审文章
func (d *Dao) FlowSync(c context.Context, mid, aid int64) (err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", _typeArticleFlow)
	resp := struct {
		Code int
		Data interface{}
	}{}
	if err = d.httpClient.Post(c, d.c.Job.FlowURL, "", params, &resp); err != nil {
		log.Error("flow: d.FlowSync.Post(%s) error(%+v)", d.c.Job.FlowURL+params.Encode(), err)
		PromError("flow:文章过审")
		return
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		log.Error("flow: d.FlowSync.Post(%s) error(%+v)", d.c.Job.FlowURL+"?"+params.Encode(), resp)
		PromError("flow:文章过审")
		return
	}
	log.Info("flow: dao.FlowSync success aid: %v mid: %v ", aid, mid)
	return
}
