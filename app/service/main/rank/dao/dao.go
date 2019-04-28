package dao

import (
	"context"

	"go-common/app/service/main/rank/conf"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	dbArchive *xsql.DB
	dbStat    *xsql.DB
	dbTV      *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:         c,
		dbArchive: xsql.NewMySQL(c.MySQL.BilibiliArchive),
		dbStat:    xsql.NewMySQL(c.MySQL.ArchiveStat),
		dbTV:      xsql.NewMySQL(c.MySQL.BilibiliTV),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.dbArchive.Close()
	d.dbStat.Close()
	d.dbTV.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	var err error
	if err = d.dbArchive.Ping(c); err != nil {
		return err
	}
	if err = d.dbStat.Ping(c); err != nil {
		return err
	}
	return d.dbTV.Ping(c)
}
