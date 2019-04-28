package dao

import (
	"context"
	"strings"

	"go-common/app/admin/main/upload/conf"

	"github.com/tsuna/gohbase"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	hbase gohbase.AdminClient
	Bfs   *Bfs
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		hbase: gohbase.NewAdminClient(strings.Join(c.Hbase.Zookeeper.Addrs, ",")),
		Bfs:   NewBfs(c),
	}
	return dao
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
