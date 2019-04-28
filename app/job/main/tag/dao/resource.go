package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_upResTagSQL     = "UPDATE resource_tag_%s SET state=?,mid=?,role=?,mtime=? WHERE oid=? AND `type`=? AND tid in (%s)"
	_upTagResSQL     = "UPDATE tag_resource_%s SET state=?,mid=?,role=?,mtime=? WHERE oid=? AND `type`=? AND tid=?"
	_insertResTagSQL = "INSERT IGNORE INTO resource_tag_%s (oid,type,tid,mid,role,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=?,mtime=?"
	_insertTagResSQL = "INSERT IGNORE INTO tag_resource_%s (oid,type,tid,mid,role,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=?,mtime=?"
)

// UpdateResTags .
func (d *Dao) UpdateResTags(c context.Context, rt *model.ResTag) (affect int64, err error) {
	var (
		res sql.Result
	)
	tx, err := d.platform.Begin(c)
	if err != nil {
		log.Error("d.UpdateResTags(%v) BeginTx error(%v)", rt, err)
		return
	}
	if res, err = tx.Exec(fmt.Sprintf(_upResTagSQL, d.hit(rt.Oid), xstr.JoinInts(rt.Tids)), rt.State, rt.Mid, rt.Role, time.Time(rt.MTime), rt.Oid, rt.Type); err != nil {
		log.Error("d.UpdateResTags(_upResTagSQL): %v error(%v)", rt, err)
		tx.Rollback()
		return
	}
	for _, tid := range rt.Tids {
		if res, err = tx.Exec(fmt.Sprintf(_upTagResSQL, d.hit(tid)), rt.State, rt.Mid, rt.Role, time.Time(rt.MTime), rt.Oid, rt.Type, tid); err != nil {
			log.Error("d.UpdateResTags(_upTagResSQL):%v error(%v)", rt, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("d.UpdateResTags(%v) Commit() error(%v)", rt, err)
		tx.Rollback()
		return
	}
	return res.RowsAffected()
}

// InsertResTags .
func (d *Dao) InsertResTags(c context.Context, rt *model.ResTag) (affect int64, err error) {
	var (
		res sql.Result
	)
	tx, err := d.platform.Begin(c)
	if err != nil {
		log.Error("d.InsertResTags(%v) BeginTx error(%v)", rt, err)
		return
	}
	for _, tid := range rt.Tids {
		if _, err = tx.Exec(fmt.Sprintf(_insertResTagSQL, d.hit(rt.Oid)), rt.Oid, rt.Type, tid, rt.Mid, rt.Role, rt.State, time.Time(rt.MTime), time.Time(rt.MTime), rt.Mid, rt.Role, rt.State, time.Time(rt.MTime)); err != nil {
			log.Error("d.InsertResTags(_insertResTagSQL):%v error(%v)", rt, err)
			tx.Rollback()
			return
		}
		if res, err = tx.Exec(fmt.Sprintf(_insertTagResSQL, d.hit(tid)), rt.Oid, rt.Type, tid, rt.Mid, rt.Role, rt.State, time.Time(rt.MTime), time.Time(rt.MTime), rt.Mid, rt.Role, rt.State, time.Time(rt.MTime)); err != nil {
			log.Error("d.InsertResTags(_insertTagResSQL):%v error(%v)", rt, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("d.InsertResTags(%v) Commit() error(%v)", rt, err)
		tx.Rollback()
		return
	}
	return res.RowsAffected()
}
