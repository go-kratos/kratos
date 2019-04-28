package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// consts
const (
	AttrFollowing = uint32(1) << 1
	AttrFriend    = uint32(1) << 2
)

// DelFollowerCache del follower cache
func (d *Dao) DelFollowerCache(fid int64) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(fid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(context.TODO(), d.clearFollowerPath, "", params, &res); err != nil {
		log.Error("d.client.Post error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.clearFollowerPath+"?"+params.Encode(), res.Code)
	}
	return
}
