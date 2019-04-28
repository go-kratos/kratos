package manager

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/admin/main/videoup/model/manager"
	"go-common/library/log"
	"net/url"
)

const (
	_upGroupURI = "/x/internal/uper/group/get"
)

// UpGroups get all up groups
func (d *Dao) UpGroups(c context.Context) (groups map[int64]*manager.UpGroup, err error) {
	groups = make(map[int64]*manager.UpGroup)
	params := url.Values{}
	params.Set("state", "1")
	var res *manager.UpGroupData
	if err = d.httpClient.Get(c, d.c.Host.API+_upGroupURI, "", params, &res); err != nil {
		log.Error("d.UpGroups() error(%v)", err)
		return
	}
	if res == nil {
		err = errors.New("error when get up groups")
		return
	}
	for _, v := range res.Data {
		groups[v.ID] = v
	}
	return
}
