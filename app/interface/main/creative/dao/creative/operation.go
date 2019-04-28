package creative

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/creative/model/operation"
	"go-common/library/log"
)

const (
	// select
	_getToolByTypeSQL    = "SELECT id,type,rank,pic,link,content,remark,note,stime,etime,platform FROM operations WHERE type=? AND stime <= ? AND etime >=? AND dtime = '0000-00-00 00:00:00' ORDER BY rank"
	_getOperByTypesSQL   = "SELECT id,type,rank,pic,link,content,remark,stime,etime,app_pic,platform FROM operations WHERE type IN (%s) AND stime <= ? AND etime >=? AND dtime = '0000-00-00 00:00:00' ORDER BY rank"
	_getAllOperByTypeSQL = "SELECT id,type,rank,pic,link,content,stime,etime,app_pic,platform FROM operations WHERE type IN (%s) and stime <= ? AND dtime = '0000-00-00 00:00:00' ORDER BY ctime desc"
)

// Tool get all tool.
func (d *Dao) Tool(c context.Context, ty string) (ops []*operation.Operation, err error) {
	now := time.Now()
	rows, err := d.getToolStmt.Query(c, ty, now, now)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stime, etime time.Time
		op := &operation.Operation{}
		if err = rows.Scan(&op.ID, &op.Ty, &op.Rank, &op.Pic, &op.Link, &op.Content, &op.Remark, &op.Note, &stime, &etime, &op.Platform); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		op.Stime = stime.Format("2006-01-02 15:04:05")
		op.Etime = etime.Format("2006-01-02 15:04:05")
		ops = append(ops, op)
	}
	return
}

// Operations get all operations.
func (d *Dao) Operations(c context.Context, tys []string) (ops []*operation.Operation, err error) {
	now := time.Now()
	rows, err := d.creativeDb.Query(c, fmt.Sprintf(_getOperByTypesSQL, strings.Join(tys, ",")), now, now)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stime, etime time.Time
		op := &operation.Operation{}
		if err = rows.Scan(&op.ID, &op.Ty, &op.Rank, &op.Pic, &op.Link, &op.Content, &op.Remark, &stime, &etime, &op.AppPic, &op.Platform); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		op.Stime = stime.Format("2006-01-02 15:04:05")
		op.Etime = etime.Format("2006-01-02 15:04:05")
		ops = append(ops, op)
	}
	return
}

// AllOperByTypeSQL fn
func (d *Dao) AllOperByTypeSQL(c context.Context, tys []string) (ops []*operation.Operation, err error) {
	now := time.Now()
	rows, err := d.creativeDb.Query(c, fmt.Sprintf(_getAllOperByTypeSQL, strings.Join(tys, ",")), now)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stime, etime time.Time
		op := &operation.Operation{}
		if err = rows.Scan(&op.ID, &op.Ty, &op.Rank, &op.Pic, &op.Link, &op.Content, &stime, &etime, &op.AppPic, &op.Platform); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		op.Stime = stime.Format("2006-01-02 15:04:05")
		op.Etime = etime.Format("2006-01-02 15:04:05")
		ops = append(ops, op)
	}
	return
}
