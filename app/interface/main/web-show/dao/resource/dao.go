package resource

import (
	"context"

	"go-common/app/interface/main/web-show/conf"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

//Dao struct
type Dao struct {
	db      *xsql.DB
	videodb *xsql.DB
	// ad_active
	selAllVdoActStmt   *xsql.Stmt
	selVdoActMTCntStmt *xsql.Stmt
	delAllVdoActStmt   *xsql.Stmt
	// ad
	selAdVdoActStmt   *xsql.Stmt
	selAdMtCntVdoStmt *xsql.Stmt
	// res
	selAllResStmt    *xsql.Stmt
	selAllAssignStmt *xsql.Stmt
	selDefBannerStmt *xsql.Stmt
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		db:      xsql.NewMySQL(c.MySQL.Res),
		videodb: xsql.NewMySQL(c.MySQL.Ads),
	}
	dao.initActive()
	dao.initRes()
	dao.initAd()
	return
}

// Close close the resource.
func (dao *Dao) Close() {
	dao.db.Close()
}

// PromError err
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Ping Dao
func (dao *Dao) Ping(c context.Context) (err error) {
	if err = dao.db.Ping(c); err != nil {
		log.Error("dao.db.Ping error(%v)", err)
		return
	}
	if err = dao.videodb.Ping(c); err != nil {
		log.Error("dao.videodb.Ping error(%v)", err)
	}
	return
}

//BeginTran Dao
func (dao *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	tx, err = dao.videodb.Begin(c)
	return
}
