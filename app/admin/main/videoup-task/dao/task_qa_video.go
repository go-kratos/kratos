package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	xsql "database/sql"
	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_inQAVideo   = "INSERT INTO task_qa_video(cid,aid,task_id,task_utime,attribute,mid,fans,up_groups,arc_title,arc_typeid,audit_status,audit_tagid,audit_submit,audit_details,ctime,mtime) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_QATaskVideo = `SELECT qa.id, qa.state, qa.type, qa.detail_id, qa.uid, qa.ftime, qa.ctime,
qav.cid, qav.aid, coalesce(qav.task_id, 0) task_id, qav.task_utime, qav.attribute, qav.audit_tagid, qav.arc_title, qav.arc_typeid, qav.audit_status, qav.audit_submit, qav.audit_details, qav.mid, qav.fans, qav.up_groups
FROM task_qa qa LEFT JOIN task_qa_video qav ON qa.detail_id = qav.id WHERE qa.id = ? LIMIT 1`
	_QATaskVideoSimple = `SELECT qa.id, qa.state, qa.type, qa.detail_id, qa.uid, qa.ftime, qa.ctime,
qav.cid, qav.aid, coalesce(qav.task_id, 0) task_id,qav.task_utime, qav.attribute, qav.audit_tagid, qav.arc_title, qav.arc_typeid, qav.audit_status, qav.audit_submit, qav.mid, qav.fans, qav.up_groups
FROM task_qa qa LEFT JOIN task_qa_video qav ON qa.detail_id = qav.id WHERE qa.id = ? LIMIT 1`
	_QAVideoDetail   = "SELECT qa.id, qav.id, qav.mid, qav.task_utime FROM task_qa qa LEFT JOIN task_qa_video qav ON qa.detail_id = qav.id WHERE qa.id IN (%s)"
	_QAVideoByTASKID = `SELECT id FROM task_qa_video WHERE aid=? AND cid=? AND task_id=?`
	_upQAVideoUTime  = "UPDATE task_qa_video SET task_utime=? WHERE aid=? AND cid=? AND task_id=?"
	_delQAVideo      = "DELETE FROM task_qa_video WHERE mtime<? LIMIT ?"
	_delQATask       = "DELETE FROM task_qa WHERE mtime<? LIMIT ?"
)

//InsertQAVideo insert qa video detail
func (d *Dao) InsertQAVideo(tx *sql.Tx, dt *model.VideoDetail) (id int64, err error) {
	var (
		groups string
	)

	if len(dt.UPGroups) > 0 {
		var b []byte
		b, err = json.Marshal(dt.UPGroups)
		if err != nil {
			log.Error("InsertQAVideo json.Marshal(%v) error(%v) aid(%d) cid(%d)", dt.UPGroups, err, dt.AID, dt.CID)
			return
		}
		groups = string(b)
	}

	now := time.Now()
	res, err := tx.Exec(_inQAVideo, dt.CID, dt.AID, dt.TaskID, dt.TaskUTime, dt.Attribute, dt.MID, dt.Fans, groups,
		dt.ArcTitle, dt.ArcTypeID, dt.AuditStatus, dt.TagID, dt.AuditSubmit, dt.AuditDetails, now, now)
	if err != nil {
		PromeErr("arcdb: exec", "InsertQAVideo tx.Exec error(%v) aid(%d) cid(%d)", err, dt.AID, dt.CID)
		return
	}

	id, err = res.LastInsertId()
	return
}

//QATaskVideoByID get by id
func (d *Dao) QATaskVideoByID(ctx context.Context, id int64) (q *model.QATaskVideo, err error) {
	var (
		groups string
	)
	q = new(model.QATaskVideo)
	if err = d.arcReadDB.QueryRow(ctx, _QATaskVideo, id).Scan(&q.ID, &q.State, &q.Type, &q.DetailID, &q.UID, &q.FTime, &q.CTime,
		&q.CID, &q.AID, &q.TaskID, &q.TaskUTime, &q.Attribute, &q.TagID, &q.ArcTitle, &q.ArcTypeID, &q.AuditStatus, &q.AuditSubmit, &q.AuditDetails,
		&q.MID, &q.Fans, &groups); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			q = nil
		} else {
			PromeErr("arcReaddb: scan", "QATaskVideoByID row.Scan error(%v) id(%d)", err, id)
		}
		return
	}

	q.UPGroups = []int64{}
	if groups != "" {
		if err = json.Unmarshal([]byte(groups), &q.UPGroups); err != nil {
			log.Error("QATaskVideoByID json.Unmarshal(%s) error(%v) id(%d)", groups, err, id)
			return
		}
	}
	return
}

// QATaskVideoSimpleByID get without audit_details by id
func (d *Dao) QATaskVideoSimpleByID(ctx context.Context, id int64) (q *model.QATaskVideo, err error) {
	var (
		groups string
	)
	q = new(model.QATaskVideo)
	if err = d.arcReadDB.QueryRow(ctx, _QATaskVideoSimple, id).Scan(&q.ID, &q.State, &q.Type, &q.DetailID, &q.UID, &q.FTime, &q.CTime,
		&q.CID, &q.AID, &q.TaskID, &q.TaskUTime, &q.Attribute, &q.TagID, &q.ArcTitle, &q.ArcTypeID, &q.AuditStatus, &q.AuditSubmit,
		&q.MID, &q.Fans, &groups); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			q = nil
		} else {
			PromeErr("arcReaddb: scan", "QATaskVideoSimpleByID row.Scan error(%v) id(%d)", err, id)
		}
		return
	}

	q.UPGroups = []int64{}
	if groups != "" {
		if err = json.Unmarshal([]byte(groups), &q.UPGroups); err != nil {
			log.Error("QATaskVideoByID json.Unmarshal(%s) error(%v) id(%d)", groups, err, id)
			return
		}
	}
	return
}

//QAVideoDetail get detail id & task_utime
func (d *Dao) QAVideoDetail(ctx context.Context, ids []int64) (list map[int64]map[string]int64, arr []int64, err error) {
	var (
		idStr string
		rows  *sql.Rows
	)

	list = map[int64]map[string]int64{}
	arr = make([]int64, 0)
	idStr = xstr.JoinInts(ids)
	if rows, err = d.arcReadDB.Query(ctx, fmt.Sprintf(_QAVideoDetail, idStr)); err != nil {
		PromeErr("arcReaddb: query", "QAVideoDetail d.arcReadDB.Query error(%v) ids(%s)", err, idStr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, detailID, mid, taskUTime int64
		)
		if err = rows.Scan(&id, &detailID, &mid, &taskUTime); err != nil {
			PromeErr("arcReaddb: scan", "QAVideoDetail rows.Scan error(%v) ids(%s)", err, idStr)
			return
		}

		list[id] = map[string]int64{
			"detail_id":  detailID,
			"mid":        mid,
			"task_utime": taskUTime,
		}
		arr = append(arr, mid)
	}

	return
}

//GetQAVideoID get id by aid & cid & taskid
func (d *Dao) GetQAVideoID(ctx context.Context, aid int64, cid int64, taskID int64) (id int64, err error) {
	if err = d.arcDB.QueryRow(ctx, _QAVideoByTASKID, aid, cid, taskID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			id = 0
			err = nil
		} else {
			log.Error("GetQAVideoID scan error(%v) aid(%d) cid(%d) taskid(%d)", err, aid, cid, taskID)
		}
	}
	return
}

//UpdateQAVideoUTime update task_utime
func (d *Dao) UpdateQAVideoUTime(ctx context.Context, aid int64, cid int64, taskID, utime int64) (err error) {
	if _, err = d.arcDB.Exec(ctx, _upQAVideoUTime, utime, aid, cid, taskID); err != nil {
		PromeErr("arcdb: exec", "UpdateQAVideoUTime error(%v) aid(%d) cid(%d) taskid(%d)", err, aid, cid, taskID)
	}
	return
}

//DelQAVideo 删除数据
func (d *Dao) DelQAVideo(ctx context.Context, mtime time.Time, limit int) (rows int64, err error) {
	var (
		result xsql.Result
	)
	if result, err = d.arcDB.Exec(ctx, _delQAVideo, mtime, limit); err != nil {
		PromeErr("arcdb: exec", "DelQAVideo error(%v) mtime(%v)", err, mtime)
		return
	}

	rows, err = result.RowsAffected()
	return
}

//DelQATask 删除数据
func (d *Dao) DelQATask(ctx context.Context, mtime time.Time, limit int) (rows int64, err error) {
	var (
		result xsql.Result
	)
	if result, err = d.arcDB.Exec(ctx, _delQATask, mtime, limit); err != nil {
		PromeErr("arcdb: exec", "DelQATask error(%v) mtime(%v)", err, mtime)
		return
	}

	rows, err = result.RowsAffected()
	return
}
