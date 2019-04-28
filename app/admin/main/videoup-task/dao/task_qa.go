package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
)

const (
	_inTaskQA = "INSERT INTO task_qa(state,type,detail_id,uid,ctime,mtime) VALUES(?,?,?,?,?,?)"
	_upTaskQA = "UPDATE task_qa SET state=?, ftime=?, mtime=? WHERE id=?"
)

//InTaskQA insert a qa task
func (d *Dao) InTaskQA(tx *sql.Tx, uid int64, detailID int64, taskType int8) (id int64, err error) {
	now := time.Now()
	res, err := tx.Exec(_inTaskQA, model.QAStateWait, taskType, detailID, uid, now, now)
	if err != nil {
		PromeErr("arcdb: insert", "InTaskQA tx.Exe error(%v) uid(%d) detailid(%d)", err, uid, detailID)
		return
	}

	id, err = res.LastInsertId()
	return
}

//UpTask update qa task
func (d *Dao) UpTask(ctx context.Context, id int64, state int16, ftime time.Time) (rows int64, err error) {
	res, err := d.arcDB.Exec(ctx, _upTaskQA, state, ftime, ftime, id)
	if err != nil {
		PromeErr("arcdb: update", "UpTask d.arcDB.Exec error(%v) id(%d) state(%d)", err, id, state)
		return
	}

	rows, err = res.RowsAffected()
	return
}
