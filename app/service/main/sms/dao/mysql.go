package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/sms/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_templateByStatusSQL = "SELECT id,code,stype,template,status,submitter,approver FROM sms_template_new WHERE status=?"
	_taskSQL             = "SELECT id,type,business_id,template_code,`desc`,file_path,file_name,send_time,status,ctime,mtime FROM sms_tasks WHERE status=? AND send_time<=? LIMIT 1 FOR UPDATE"
	_upadteTaskStatusSQL = "UPDATE sms_tasks SET status=? WHERE id=?"
)

// TemplateByStatus select template by status
func (d *Dao) TemplateByStatus(ctx context.Context, status int) (res []*model.ModelTemplate, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(ctx, _templateByStatusSQL, status); err != nil {
		log.Error("d.TemplateByStatus.Query(%d) error(%v)", status, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ModelTemplate)
		if err = rows.Scan(&r.ID, &r.Code, &r.Stype, &r.Template, &r.Status, &r.Submitter, &r.Approver); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BeginTx begin transaction.
func (d *Dao) BeginTx(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// TxTask gets prepared task.
func (d *Dao) TxTask(tx *xsql.Tx) (t *model.ModelTask, err error) {
	t = new(model.ModelTask)
	if err = tx.QueryRow(_taskSQL, model.TaskStatusPrepared, time.Now()).Scan(&t.ID, &t.Type, &t.BusinessID,
		&t.TemplateCode, &t.Desc, &t.FilePath, &t.FileName, &t.SendTime, &t.Status, &t.Ctime, &t.Mtime); err != nil {
		t = nil
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.TxTask() QueryRow() error(%v)", err)
	}
	return
}

// TxUpdateTaskStatus updates task status by tx.
func (d *Dao) TxUpdateTaskStatus(tx *xsql.Tx, taskID int64, status int32) (err error) {
	if _, err = tx.Exec(_upadteTaskStatusSQL, status, taskID); err != nil {
		log.Error("d.TxUpdateTaskStatus() Exec(%s,%d) error(%v)", taskID, status, err)
	}
	return
}

// UpdateTaskStatus updates task status.
func (d *Dao) UpdateTaskStatus(ctx context.Context, taskID int64, status int32) (err error) {
	if _, err = d.db.Exec(ctx, _upadteTaskStatusSQL, status, taskID); err != nil {
		log.Error("d.UpdateTaskStatus() Exec(%s,%d) error(%v)", taskID, status, err)
	}
	return
}
