package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/bfs/conf"
	bm "go-common/library/net/http/blademaster"

	"github.com/samuel/go-zookeeper/zk"
)

// Dao dao
type Dao struct {
	c       *conf.Config
	zkcs    map[string]*zk.Conn
	httpCli *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	zkcs := make(map[string]*zk.Conn)
	for name, zkCfg := range c.Zookeepers {
		zkc, _, err := zk.Connect(zkCfg.Addrs, time.Duration(zkCfg.Timeout))
		if err != nil {
			panic(err)
		}
		zkcs[name] = zkc
	}
	dao = &Dao{
		c:       c,
		zkcs:    zkcs,
		httpCli: bm.NewClient(c.HTTPClient),
	}
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}

// Close close the resource.
func (d *Dao) Close() {
	for _, zkc := range d.zkcs {
		zkc.Close()
	}
}
