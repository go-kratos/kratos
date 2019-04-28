package show

import (
	"context"
	"time"

	"go-common/app/job/main/app/conf"
	"go-common/library/database/sql"
)

const (
	// tran pub
	_ptimeSQL = "SELECT plat FROM show_time WHERE state=1 AND ptime<?"
	_upStSQL  = "UPDATE show_time SET state=0 WHERE plat=?"
	_delHdSQL = "DELETE FROM show_head WHERE plat=?"
	_delItSQL = "DELETE FROM show_item WHERE plat=?"
	_cpHdSQL  = "INSERT INTO show_head(id,plat,title,type,style,param,rank,build,conditions,lang_id,ctime,mtime) SELECT id,plat,title,type,style,param,rank,build,conditions,lang_id,ctime,mtime FROM show_head_temp WHERE plat=?"
	_cpItSQL  = "INSERT INTO show_item(id,sid,plat,title,random,cover,param,ctime,mtime) SELECT id,sid,plat,title,random,cover,param,ctime,mtime FROM show_item_temp WHERE plat=?"
)

// Dao is show dao.
type Dao struct {
	db       *sql.DB
	getPTime *sql.Stmt
}

// New new a show dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// mysql
		db: sql.NewMySQL(c.MySQL.Show),
	}
	d.getPTime = d.db.Prepared(_ptimeSQL)
	return
}

// BeginTran begin a transacition
func (d *Dao) BeginTran(ctx context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(ctx)
}

// PTime get timing publis time.
func (d *Dao) PTime(ctx context.Context, now time.Time) (ps []int8, err error) {
	rows, err := d.getPTime.Query(ctx, now)
	if err != nil {
		return
	}
	defer rows.Close()
	var plat int8
	for rows.Next() {
		if err = rows.Scan(&plat); err != nil {
			return
		}
		ps = append(ps, plat)
	}
	return
}

// Pub check ptime and publish.
func (d *Dao) Pub(tx *sql.Tx, plat int8) (err error) {
	if _, err = tx.Exec(_delHdSQL, plat); err != nil {
		return
	}
	if _, err = tx.Exec(_delItSQL, plat); err != nil {
		return
	}
	if _, err = tx.Exec(_cpHdSQL, plat); err != nil {
		return
	}
	if _, err = tx.Exec(_cpItSQL, plat); err != nil {
		return
	}
	if _, err = tx.Exec(_upStSQL, plat); err != nil {
		return
	}
	return
}

func (dao *Dao) PingDB(c context.Context) (err error) {
	return dao.db.Ping(c)
}
