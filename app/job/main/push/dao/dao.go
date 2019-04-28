package dao

import (
	"context"

	"go-common/app/job/main/push/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	c                    *conf.Config
	db                   *xsql.DB
	mc                   *memcache.Pool
	httpClient           *bm.Client
	dpClient             *bm.Client
	delCallbacksStmt     *xsql.Stmt
	delTasksStmt         *xsql.Stmt
	reportLastIDStmt     *xsql.Stmt
	reportsByRangeStmt   *xsql.Stmt
	updateTaskStatusStmt *xsql.Stmt
	updateTaskStmt       *xsql.Stmt
	updateDpCondStmt     *xsql.Stmt
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// New creates a dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		db:         xsql.NewMySQL(c.MySQL),
		mc:         memcache.NewPool(c.Memcache.Config),
		httpClient: bm.NewClient(c.HTTPClient),
		dpClient:   bm.NewClient(c.DpClient),
	}
	d.delCallbacksStmt = d.db.Prepared(_delCallbacksSQL)
	d.delTasksStmt = d.db.Prepared(_delTasksSQL)
	d.reportLastIDStmt = d.db.Prepared(_reportLastIDSQL)
	d.reportsByRangeStmt = d.db.Prepared(_reportsByRangeSQL)
	d.updateTaskStatusStmt = d.db.Prepared(_upadteTaskStatusSQL)
	d.updateTaskStmt = d.db.Prepared(_upadteTaskSQL)
	d.updateDpCondStmt = d.db.Prepared(_updateDpCondSQL)
	return
}

// PromError prometheus error count.
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo prometheus info count.
func PromInfo(name string) {
	infosCount.Incr(name)
}

// Ping reports the health of the db/cache etc.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	err = d.pingMC(c)
	return
}

// Close .
func (d *Dao) Close() {
	d.db.Close()
	d.mc.Close()
}
