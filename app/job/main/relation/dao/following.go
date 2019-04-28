package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/relation/model"
	"go-common/library/log"
)

// UpdateFollowingCache update following cache.
func (d *Dao) UpdateFollowingCache(r *model.Relation) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(r.Mid, 10))
	params.Set("fid", strconv.FormatInt(r.Fid, 10))
	params.Set("attribute", strconv.FormatInt(int64(r.Attribute), 10))
	mt, err := time.Parse(time.RFC3339, r.MTime)
	if err != nil {
		mt = time.Now()
	}
	params.Set("mtime", strconv.FormatInt(mt.Unix(), 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(context.TODO(), d.clearFollowingPath, "", params, &res); err != nil {
		log.Error("d.client.Post error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.clearFollowingPath+"?"+params.Encode(), res.Code)
	}
	return
}
