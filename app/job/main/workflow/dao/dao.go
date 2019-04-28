package dao

import (
	"context"

	"go-common/app/job/main/workflow/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

const (
	_srhURL    = "/x/admin/search/workflow/common"
	_notifyURL = "/api/notify/send.user.notify.do"
)

// Dao struct info of Dao.
type Dao struct {
	c *conf.Config
	// orm
	WriteORM *gorm.DB
	ReadORM  *gorm.DB
	// redis
	redis *redis.Pool
	// search
	httpSearch *bm.Client
	// url
	searchURL  string
	messageURL string
	// es client
	es *elastic.Elastic
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		WriteORM:   orm.NewMySQL(c.ORM.Write),
		ReadORM:    orm.NewMySQL(c.ORM.Read),
		redis:      redis.NewPool(c.Redis),
		httpSearch: bm.NewClient(c.HTTPSearch),
		searchURL:  c.Host.SearchURI + _srhURL,
		messageURL: c.Host.MessageURI + _notifyURL,
		es:         elastic.NewElastic(nil),
	}
	d.initORM()
	return
}

// Close close connections of ORM.
func (d *Dao) Close() (err error) {
	if d.WriteORM != nil {
		d.WriteORM.Close()
	}
	if d.ReadORM != nil {
		d.ReadORM.Close()
	}
	return
}

// Ping ping health of ORM.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.WriteORM.DB().PingContext(c); err != nil {
		return
	}
	err = d.ReadORM.DB().PingContext(c)
	return
}

func (d *Dao) initORM() {
	d.WriteORM.LogMode(true)
	d.ReadORM.LogMode(true)
}
