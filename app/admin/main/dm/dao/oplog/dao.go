package oplog

import (
	"go-common/app/admin/main/dm/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao struct for querying infoc data storing in hbase
type Dao struct {
	httpCli                    *bm.Client
	key, secret, infocQueryURL string
}

// New new a Dao instance and init.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		key:           c.HTTPInfoc.ClientConfig.Key,
		secret:        c.HTTPInfoc.ClientConfig.Secret,
		httpCli:       bm.NewClient(c.HTTPInfoc.ClientConfig),
		infocQueryURL: c.HTTPInfoc.InfocQueryURL,
	}
	return
}
