package dao

import (
	"time"

	"go-common/app/service/bbq/push/conf"
	"go-common/app/service/bbq/push/dao/jpush"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	JPush *jpush.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	client := jpush.New(c.JPush.AppKey, c.JPush.SecretKey, time.Duration(c.JPush.Timeout))
	dao = &Dao{
		c:     c,
		JPush: client,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
}
