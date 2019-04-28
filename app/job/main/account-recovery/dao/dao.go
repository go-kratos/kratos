package dao

import (
	"go-common/app/job/main/account-recovery/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c *conf.Config
	// httpClient
	httpClient *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
		// httpClient
		httpClient: bm.NewClient(c.HTTPClientConfig),
	}
	return
}
