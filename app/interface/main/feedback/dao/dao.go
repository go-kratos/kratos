package dao

import (
	"context"

	"go-common/app/interface/main/feedback/conf"
	"go-common/library/database/sql"
	"net/http"
)

// Dao is feedback dao.
type Dao struct {
	// conf
	c *conf.Config
	// db
	dbMs *sql.DB
	//db stmt
	// session
	selSsn         *sql.Stmt
	selSsnByMid    *sql.Stmt
	selTagID       *sql.Stmt
	inSsn          *sql.Stmt
	inSsnTag       *sql.Stmt
	upSsn          *sql.Stmt
	upSsnMtime     *sql.Stmt
	upSsnSta       *sql.Stmt
	selSSnID       *sql.Stmt
	selSSnCntByMid *sql.Stmt
	// reply
	selReply      *sql.Stmt
	selReplyByMid *sql.Stmt
	selReplyBySid *sql.Stmt
	inReply       *sql.Stmt
	// tag
	selTagBySsnID *sql.Stmt
	// bfs
	bfsClient *http.Client
}

// New dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:    c,
		dbMs: sql.NewMySQL(c.MySQL.Master),
	}
	// session
	d.selSsn = d.dbMs.Prepared(_selSsn)
	d.selSsnByMid = d.dbMs.Prepared(_selSsnByMid)
	d.inSsn = d.dbMs.Prepared(_inSsn)
	d.inSsnTag = d.dbMs.Prepared(_inSsnTag)
	d.upSsn = d.dbMs.Prepared(_upSsn)
	d.upSsnMtime = d.dbMs.Prepared(_upSsnMtime)
	d.upSsnSta = d.dbMs.Prepared(_upSsnState)
	d.selSSnID = d.dbMs.Prepared(_selSSnID)
	d.selSSnCntByMid = d.dbMs.Prepared(_selSSnCntByMid)
	// reply
	d.selReply = d.dbMs.Prepared(_selReply)
	d.selReplyByMid = d.dbMs.Prepared(_selReplyByMid)
	d.inReply = d.dbMs.Prepared(_inReply)
	d.selTagID = d.dbMs.Prepared(_selTagID)
	d.selReplyBySid = d.dbMs.Prepared(_selReplyBySid)
	d.selTagBySsnID = d.dbMs.Prepared(_selTagBySsnID)
	// init bfs http client
	d.bfsClient = http.DefaultClient
	return
}

// BeginTran begin tran.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	tx, err = d.dbMs.Begin(c)
	return
}
