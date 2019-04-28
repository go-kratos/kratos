package tag

import (
	tagClient "go-common/app/interface/main/tag/rpc/client"
	"go-common/app/job/main/videoup-report/conf"
	bm "go-common/library/net/http/blademaster"
)

//Dao tag dao
type Dao struct {
	c                       *conf.Config
	client                  *bm.Client
	upBindURL, adminBindURL string
	tagDisRPC               *tagClient.Service
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		client:       bm.NewClient(c.HTTPClient.Write),
		upBindURL:    c.Host.API + _upBindURI,
		adminBindURL: c.Host.API + _adminBindURI,
		tagDisRPC:    tagClient.New2(c.TagDisConf),
	}
	return d
}
