package notice

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/notice"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getSQL = `SELECT id,plat,title,content,build,conditions,area,url,type,ef_time,ex_time FROM notice WHERE state=1 AND ef_time<? AND ex_time>? ORDER BY mtime DESC`
)

// Dao is notice dao.
type Dao struct {
	db  *sql.DB
	get *sql.Stmt
}

// New new a notice dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	d.get = d.db.Prepared(_getSQL)
	return
}

// GetAll get all notice data.
func (d *Dao) All(ctx context.Context, now time.Time) (res []*notice.Notice, err error) {
	rows, err := d.get.Query(ctx, now, now)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	res = []*notice.Notice{}
	for rows.Next() {
		b := &notice.Notice{}
		if err = rows.Scan(&b.ID, &b.Plat, &b.Title, &b.Content, &b.Build, &b.Condition, &b.Area, &b.URI, &b.Type, &b.Start, &b.End); err != nil {
			log.Error("rows.Scan err (%v)", err)
			return nil, err

		}
		res = append(res, b)
	}
	return
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}
