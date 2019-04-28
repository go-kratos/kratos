package dao

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/push/model"
	"go-common/library/log"
)

const (
	_addTaskSQL  = "insert into push_tasks (job,type,app_id,business_id,title,summary,link_type,link_value,sound,vibration,push_time,expire_time,status,`group`,image_url,extra) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_taskInfoSQL = "select id,job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url from push_tasks where id=?"
)

// TaskInfo .
func (d *Dao) TaskInfo(ctx context.Context, id int64) (t *model.Task, err error) {
	t = &model.Task{}
	if err = d.db.QueryRow(ctx, _taskInfoSQL, id).Scan(&t.ID, &t.Job, &t.Type, &t.AppID, &t.BusinessID, &t.Platform, &t.Title, &t.Summary, &t.LinkType, &t.LinkValue, &t.Build,
		&t.Sound, &t.Vibration, &t.PassThrough, &t.MidFile, &t.Progress, &t.PushTime, &t.ExpireTime, &t.Status, &t.Group, &t.ImageURL); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
			return
		}
		log.Error("d.TaskInfo(%s) error(%v)", id, err)
		return
	}
	t.PushTimeUnix = t.PushTime.Unix()
	t.ExpireTimeUnix = t.ExpireTime.Unix()
	return
}

// AddTask add data platform task
func (d *Dao) AddTask(ctx context.Context, t *model.Task) (id int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(ctx, _addTaskSQL, t.Job, t.Type, t.AppID, t.BusinessID, t.Title, t.Summary, t.LinkType, t.LinkValue, t.Sound, t.Vibration, t.PushTime, t.ExpireTime, t.Status, t.Group, t.ImageURL, t.Extra); err != nil {
		log.Error("d.AddTask(%+v) error(%v)", t, err)
		return
	}
	id, err = res.LastInsertId()
	return
}
