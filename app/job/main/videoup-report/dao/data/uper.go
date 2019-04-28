package data

import (
	"context"
	"go-common/app/job/main/videoup-report/model/manager"
	"go-common/library/log"
	"net/url"
	"strconv"
)

const (
	_midGroupsURI = "/x/internal/uper/special/get_by_mid"
)

// UpGroups get all up groups
func (d *Dao) MidGroups(c context.Context, mid int64) (groups map[int64]*manager.UpGroup, err error) {
	groups = make(map[int64]*manager.UpGroup)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                `json:"code"`
		Msg  string             `json:"message"`
		Data []*manager.UpGroup `json:"data"`
	}
	if err = d.client.Get(c, d.c.Host.API+_midGroupsURI, "", params, &res); err != nil {
		log.Error("d.UpGroups() error(%v)", err)
		return
	}
	if res.Data == nil {
		log.Warn("MidGroups(%d) error when get up groups", mid)
		return
	}
	for _, v := range res.Data {
		groups[v.GroupID] = v
	}
	return
}
