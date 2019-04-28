package migrate

import (
	"context"
	"go-common/app/job/live/push-search/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c                 *conf.Config
	SearchHBase       *hbase.Client
	RoomDb			  *sql.DB
}

// New init mysql db
func NewMigrate(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:           c,
		SearchHBase: hbase.NewClient(&c.SearchHBase.Config),
		RoomDb:      sql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.RoomDb.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return nil
}
