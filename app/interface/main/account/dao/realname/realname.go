package realname

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/account/conf"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

var (
	telInfoURI = "/intranet/acc/telInfo/mid"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	client *bm.Client
	mc     *memcache.Pool
}

// New new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(c.HTTPClient.Normal),
		mc:     memcache.NewPool(c.AccMemcache),
	}
	return
}

// TelInfo tel info.
func (d *Dao) TelInfo(c context.Context, mid int64) (tel string, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Mid      int64  `json:"mid"`
			Tel      string `json:"tel"`
			JoinIP   string `json:"join_ip"`
			JoinTime int64  `json:"join_time"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.c.Host.Passport+telInfoURI, "", params, &resp); err != nil {
		err = errors.Errorf("realname TelInfo d.httpClient.Do() error(%+v)", err)
		return
	}
	if resp.Code != 0 {
		err = errors.Errorf("realname TelInfo url(%s) res(%+v) err(%+v)", telInfoURI+"?"+params.Encode(), resp, ecode.Int(resp.Code))
		return
	}
	tel = resp.Data.Tel
	return
}
