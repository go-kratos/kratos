package mall

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const _favCount = "/mall-ugc/ugc/vote/user/wishcount"

type Dao struct {
	client   *bm.Client
	favCount string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:   bm.NewClient(c.HTTPClient),
		favCount: c.Host.Mall + _favCount,
	}
	return
}

func (d *Dao) FavCount(c context.Context, mid int64) (count int32, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int   `json:"code"`
		Data int32 `json:"data"`
	}
	if err = d.client.Get(c, d.favCount, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favCount+"?"+params.Encode())
		return
	}
	count = res.Data
	return
}
