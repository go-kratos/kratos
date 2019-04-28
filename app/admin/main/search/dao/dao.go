package dao

import (
	"context"

	"go-common/app/admin/main/search/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
	"go-common/library/sync/errgroup"

	"gopkg.in/olivere/elastic.v5"
)

const (
	_managerDep    = "/x/admin/manager/users/udepts"
	_managerUnames = "/x/admin/manager/users/unames"
	_managerIP     = "/x/location/infos"
)

// Dao .
type Dao struct {
	c             *conf.Config
	esPool        map[string]*elastic.Client
	db            *sql.DB
	client        *bm.Client
	managerDep    string
	managerUnames string
	managerIP     string
	queryConfStmt *sql.Stmt
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		db:            sql.NewMySQL(c.DB.Search),
		client:        bm.NewClient(c.HTTPClient),
		managerDep:    c.Prop.Manager + _managerDep,
		managerUnames: c.Prop.Manager + _managerUnames,
		managerIP:     c.Prop.API + _managerIP,
	}
	d.esPool = newEsPool(c, d)
	d.NewLog()
	go d.NewLogProcess()
	d.queryConfStmt = d.db.Prepared(_queryConfSQL)
	return
}

// BulkItem .
type BulkItem interface {
	IndexName() string
	IndexType() string
	IndexID() string
}

// BulkMapItem .
type BulkMapItem interface {
	IndexName() string
	IndexType() string
	IndexID() string
	PField() map[string]interface{}
}

// newEsCluster cluster action
func newEsPool(c *conf.Config, d *Dao) (esCluster map[string]*elastic.Client) {
	esCluster = make(map[string]*elastic.Client)
	for esName, e := range c.Es {
		cof := []elastic.ClientOptionFunc{}
		cof = append(cof, elastic.SetURL(e.Addr...))
		if esName == "ops_log" {
			cof = append(cof, elastic.SetSniff(false))
		}
		client, err := elastic.NewClient(cof...)
		if err != nil {
			PromError("es:集群连接失败", "cluster: %s, %v", esName, err)
			continue
		}
		esCluster[esName] = client
	}
	return
}

// PromError prometheus error count.
func PromError(name, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Ping health
func (d *Dao) Ping(c context.Context) (err error) {
	group := errgroup.Group{}
	group.Go(func() (err error) {
		err = d.db.Ping(context.Background())
		if err != nil {
			PromError("DB:Ping", "DB:Ping error(%v)", err)
		}
		return
	})
	for name, client := range d.esPool {
		group.Go(func() (err error) {
			_, _, err = client.Ping(d.c.Es[name].Addr[0]).Do(context.Background())
			if err != nil {
				PromError("Es:Ping", "%s:Ping error(%v)", name, err)
			}
			return
		})
	}
	return group.Wait()
}
