package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_typeArticle = "12"
)

// LikeSync like sync
func (d *Dao) LikeSync(c context.Context, aid, likes int64) (err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("likes", strconv.FormatInt(likes, 10))
	params.Set("type", _typeArticle)
	resp := struct {
		Code int
		Data interface{}
	}{}
	if err = d.httpClient.Get(c, d.c.Job.ActLikeURL, "", params, &resp); err != nil {
		log.Error("activity: d.LikeSync.Get(%s) error(%+v)", d.c.Job.ActLikeURL+params.Encode(), err)
		PromError("activity:同步点赞数")
		return
	}
	if resp.Code != 0 {
		// 未参与活动
		if resp.Code == -403 {
			return
		}
		err = ecode.Int(resp.Code)
		log.Error("activity: d.LikeSync.Get(%s) error(%+v)", d.c.Job.ActLikeURL+"?"+params.Encode(), resp)
		PromError("activity:同步点赞数")
		return
	}
	log.Info("activity: dao.LikeSync success aid: %v count: %v ", aid, likes)
	return
}
