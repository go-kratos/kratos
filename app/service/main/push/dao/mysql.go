package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// task
	_addTaskSQL            = "INSERT INTO push_tasks (job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url,extra) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_upadteTaskStatusSQL   = "UPDATE push_tasks SET status=? WHERE id=?"
	_upadteTaskProgressSQL = "UPDATE push_tasks SET progress=? WHERE id=?"
	_taskByIDSQL           = "SELECT id,job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url FROM push_tasks WHERE id=?"
	_taskByPlatformSQL     = "SELECT id,job,type,app_id,business_id,platform_id,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url FROM push_tasks WHERE status=? AND push_time<=? AND dtime=0 and platform_id=? and mtime>? LIMIT 1 FOR UPDATE"

	// business
	_businessesSQL = "SELECT id,app_id,name,description,token,sound,vibration,receive_switch,push_switch,silent_time,push_limit_user,whitelist FROM push_business WHERE dtime=0"

	// auth
	_authsSQL = "SELECT app_id,platform_id,name,`key`,value,bundle_id FROM push_auths WHERE dtime=0"

	// callback
	_addCallbackSQL = `INSERT INTO push_callbacks (task,app,platform,mid,buvid,token,pid,click,brand,extra) VALUES(?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE app=?,platform=?,mid=?,buvid=?,pid=?,click=?`
)

// AddTask adds task.
func (d *Dao) AddTask(c context.Context, t *model.Task) (id int64, err error) {
	var (
		res         sql.Result
		platform    = model.JoinInts(t.Platform)
		build, _    = json.Marshal(t.Build)
		progress, _ = json.Marshal(t.Progress)
		extra, _    = json.Marshal(t.Extra)
	)
	if res, err = d.addTaskStmt.Exec(c, t.Job, t.Type, t.APPID, t.BusinessID, platform, t.Title, t.Summary, t.LinkType, t.LinkValue,
		build, t.Sound, t.Vibration, t.PassThrough, t.MidFile, progress, t.PushTime, t.ExpireTime, t.Status, t.Group, t.ImageURL, extra); err != nil {
		log.Error("d.AddTask(%+v) error(%v)", t, err)
		PromError("mysql:添加推送任务")
		return
	}
	id, err = res.LastInsertId()
	return
}

// TxTaskByPlatform gets prepared task by platform.
func (d *Dao) TxTaskByPlatform(tx *xsql.Tx, platform int) (t *model.Task, err error) {
	var (
		id       int64
		build    string
		progress string
		now      = time.Now()
		since    = now.Add(-7 * 24 * time.Hour)
	)
	t = &model.Task{Progress: &model.Progress{}}
	if err = tx.QueryRow(_taskByPlatformSQL, model.TaskStatusPrepared, now, platform, since).Scan(&id, &t.Job, &t.Type, &t.APPID, &t.BusinessID, &t.PlatformID, &t.Title, &t.Summary, &t.LinkType, &t.LinkValue, &build,
		&t.Sound, &t.Vibration, &t.PassThrough, &t.MidFile, &progress, &t.PushTime, &t.ExpireTime, &t.Status, &t.Group, &t.ImageURL); err != nil {
		t = nil
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.TxPreparedTask() QueryRow(%d,%v) error(%v)", platform, now, err)
		PromError("mysql:按状态查询任务")
		return
	}
	t.ID = strconv.FormatInt(id, 10)
	t.Build = model.ParseBuild(build)
	if progress != "" {
		if err = json.Unmarshal([]byte(progress), t.Progress); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", progress, err)
			PromError("mysql:unmarshal进度")
		}
	}
	return
}

// Task loads task by task id.
func (d *Dao) Task(c context.Context, taskID string) (t *model.Task, err error) {
	var (
		platform string
		build    string
		progress string
		id, _    = strconv.ParseInt(taskID, 10, 64)
	)
	t = &model.Task{Progress: &model.Progress{}}
	if err = d.taskStmt.QueryRow(c, id).Scan(&id, &t.Job, &t.Type, &t.APPID, &t.BusinessID, &platform, &t.Title, &t.Summary, &t.LinkType, &t.LinkValue, &build,
		&t.Sound, &t.Vibration, &t.PassThrough, &t.MidFile, &progress, &t.PushTime, &t.ExpireTime, &t.Status, &t.Group, &t.ImageURL); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
			return
		}
		log.Error("d.taskStmt.QueryRow(%s) error(%v)", taskID, err)
		PromError("mysql:按ID查询任务")
		return
	}
	t.ID = strconv.FormatInt(id, 10)
	t.Platform = model.SplitInts(platform)
	t.Build = model.ParseBuild(build)
	if progress != "" {
		if err = json.Unmarshal([]byte(progress), t.Progress); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", progress, err)
			PromError("mysql:unmarshal进度")
		}
	}
	if t.Progress == nil {
		t.Progress = &model.Progress{}
	}
	return
}

// UpdateTaskStatus update task status.
func (d *Dao) UpdateTaskStatus(c context.Context, taskID string, status int8) (err error) {
	id, _ := strconv.ParseInt(taskID, 10, 64)
	if _, err = d.updateTaskStatusStmt.Exec(c, status, id); err != nil {
		log.Error("d.updateTaskStatusStmt.Exec(%s,%d) error(%v)", taskID, status, err)
		PromError("mysql:更新推送任务状态")
	}
	return
}

// TxUpdateTaskStatus updates task status by tx.
func (d *Dao) TxUpdateTaskStatus(tx *xsql.Tx, taskID string, status int8) (err error) {
	id, _ := strconv.ParseInt(taskID, 10, 64)
	if _, err = tx.Exec(_upadteTaskStatusSQL, status, id); err != nil {
		log.Error("d.TxUpdateTaskStatus() Exec(%s,%d) error(%v)", taskID, status, err)
		PromError("mysql:更新推送任务状态")
	}
	return
}

// UpdateTaskProgress updates task's progress.
func (d *Dao) UpdateTaskProgress(c context.Context, taskID string, p *model.Progress) (err error) {
	var b []byte
	if b, err = json.Marshal(p); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", p, err)
		return
	}
	id, _ := strconv.ParseInt(taskID, 10, 64)
	if _, err = d.updateTaskProgressStmt.Exec(c, string(b), id); err != nil {
		log.Error("d.updateTaskProgress.Exec(%s,%+v) error(%v)", taskID, p, err)
		PromError("mysql:更新推送任务进度")
	}
	return
}

// Businesses gets all business info.
func (d *Dao) Businesses(c context.Context) (res map[int64]*model.Business, err error) {
	rows, err := d.businessesStmt.Query(c)
	if err != nil {
		log.Error("d.businessesStmt.Query() error(%v)", err)
		PromError("mysql:查询业务方")
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Business)
	for rows.Next() {
		var (
			silentTime string
			b          = &model.Business{}
		)
		if err = rows.Scan(&b.ID, &b.APPID, &b.Name, &b.Desc, &b.Token,
			&b.Sound, &b.Vibration, &b.ReceiveSwitch, &b.PushSwitch, &silentTime, &b.PushLimitUser, &b.Whitelist); err != nil {
			PromError("mysql:查询业务方Scan")
			log.Error("d.Business() Scan() error(%v)", err)
			return
		}
		b.SilentTime = model.ParseSilentTime(silentTime)
		res[b.ID] = b
	}
	return
}

func (d *Dao) auths(c context.Context) (res []*model.Auth, err error) {
	var rows *xsql.Rows
	if rows, err = d.authsStmt.Query(c); err != nil {
		log.Error("d.authsStmt.Query() error(%v)", err)
		PromError("mysql:获取 auths")
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Auth{}
		if err = rows.Scan(&a.APPID, &a.PlatformID, &a.Name, &a.Key, &a.Value, &a.BundleID); err != nil {
			PromError("mysql:获取 auths Scan")
			log.Error("d.auths() Scan() error(%v)", err)
			return
		}
		res = append(res, a)
	}
	return
}

// AddCallback adds callback.
func (d *Dao) AddCallback(c context.Context, cb *model.Callback) (err error) {
	var extra []byte
	if cb.Extra != nil {
		extra, _ = json.Marshal(cb.Extra)
	}
	if _, err = d.addCallbackStmt.Exec(c, cb.Task, cb.APP, cb.Platform, cb.Mid, cb.Buvid, cb.Token, cb.Pid, cb.Click, cb.Brand, string(extra),
		cb.APP, cb.Platform, cb.Mid, cb.Buvid, cb.Pid, cb.Click); err != nil {
		log.Error("d.AddCallback(%+v) error(%v)", cb, err)
		PromError("mysql:新增callback")
		return
	}
	PromInfo("mysql:新增callback")
	return
}
