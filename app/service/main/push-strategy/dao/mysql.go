package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	pushmdl "go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// task
	_addTaskSQL = "INSERT INTO push_tasks (job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url,extra) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_taskSQL    = "SELECT id,job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url,extra FROM push_tasks WHERE id=?"
	// app && business
	_appsSQL       = "SELECT id,name,push_limit_user FROM push_apps WHERE dtime=0"
	_businessesSQL = "SELECT id,app_id,name,description,token,sound,vibration,receive_switch,push_switch,silent_time,push_limit_user,whitelist FROM push_business WHERE dtime=0"
	// setting
	_settingsByRangeSQL = "SELECT mid,value FROM push_user_settings WHERE id>? AND id<=? AND dtime=0"
	_maxSettingIDSQL    = "SELECT MAX(id) FROM push_user_settings"
)

// Apps get all app info
func (d *Dao) Apps(ctx context.Context) (res map[int64]*pushmdl.APP, err error) {
	rows, err := d.appsStmt.Query(ctx)
	if err != nil {
		log.Error("d.appsStmt.Query() error(%v)", err)
		PromError("mysql:查询应用")
		return
	}
	defer rows.Close()
	res = make(map[int64]*pushmdl.APP)
	for rows.Next() {
		app := new(pushmdl.APP)
		if err = rows.Scan(&app.ID, &app.Name, &app.PushLimitUser); err != nil {
			log.Error("d.Apps() Scan() error(%v)", err)
			PromError("mysql:查询应用Scan")
			return
		}
		res[app.ID] = app
	}
	err = rows.Err()
	return
}

// Businesses gets all business info.
func (d *Dao) Businesses(ctx context.Context) (res map[int64]*pushmdl.Business, err error) {
	rows, err := d.businessesStmt.Query(ctx)
	if err != nil {
		log.Error("d.businessesStmt.Query() error(%v)", err)
		PromError("mysql:查询业务方")
		return
	}
	defer rows.Close()
	res = make(map[int64]*pushmdl.Business)
	for rows.Next() {
		var (
			silentTime string
			b          = &pushmdl.Business{}
		)
		if err = rows.Scan(&b.ID, &b.APPID, &b.Name, &b.Desc, &b.Token,
			&b.Sound, &b.Vibration, &b.ReceiveSwitch, &b.PushSwitch, &silentTime, &b.PushLimitUser, &b.Whitelist); err != nil {
			PromError("mysql:查询业务方Scan")
			log.Error("d.Business() Scan() error(%v)", err)
			return
		}
		b.SilentTime = pushmdl.ParseSilentTime(silentTime)
		res[b.ID] = b
	}
	err = rows.Err()
	return
}

// AddTask adds task.
func (d *Dao) AddTask(ctx context.Context, t *pushmdl.Task) (id int64, err error) {
	var (
		res         sql.Result
		platform    = pushmdl.JoinInts(t.Platform)
		build, _    = json.Marshal(t.Build)
		progress, _ = json.Marshal(t.Progress)
		extra, _    = json.Marshal(t.Extra)
	)
	if res, err = d.addTaskStmt.Exec(ctx, t.Job, t.Type, t.APPID, t.BusinessID, platform, t.Title, t.Summary, t.LinkType, t.LinkValue,
		build, t.Sound, t.Vibration, t.PassThrough, t.MidFile, progress, t.PushTime, t.ExpireTime, t.Status, t.Group, t.ImageURL, extra); err != nil {
		log.Error("d.AddTask(%+v) error(%v)", t, err)
		PromError("mysql:添加推送任务")
		return
	}
	id, err = res.LastInsertId()
	return
}

// Task loads task by id.
func (d *Dao) Task(ctx context.Context, id int64) (t *pushmdl.Task, err error) {
	var (
		platform string
		build    string
		progress string
		extra    string
		now      = time.Now()
	)
	t = &pushmdl.Task{Progress: &pushmdl.Progress{}, Extra: &pushmdl.TaskExtra{}}
	if err = d.taskStmt.QueryRow(ctx, id).Scan(&id, &t.Job, &t.Type, &t.APPID, &t.BusinessID, &platform, &t.Title, &t.Summary, &t.LinkType, &t.LinkValue, &build,
		&t.Sound, &t.Vibration, &t.PassThrough, &t.MidFile, &progress, &t.PushTime, &t.ExpireTime, &t.Status, &t.Group, &t.ImageURL, &extra); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
			return
		}
		log.Error("d.taskStmt.QueryRow(%d) error(%v)", id, now, err)
		PromError("mysql:按ID查询任务")
		return
	}
	t.ID = strconv.FormatInt(id, 10)
	t.Platform = pushmdl.SplitInts(platform)
	t.Build = pushmdl.ParseBuild(build)
	if progress != "" {
		if err = json.Unmarshal([]byte(progress), t.Progress); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", progress, err)
			PromError("mysql:unmarshal progress")
		}
	}
	if extra != "" {
		if err = json.Unmarshal([]byte(extra), t.Extra); err != nil {
			log.Error("json.Unmarshal(%s) extra error(%v)", extra, err)
			PromError("mysql:unmarshal extra")
		}
	}
	return
}

// MaxSettingID gets max setting id in DB.
func (d *Dao) MaxSettingID(ctx context.Context) (id int64, err error) {
	if err = d.maxSettingIDStmt.QueryRow(ctx).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.maxSettingIDStmt.QueryRow.Scan error(%v)", err)
		PromError("db:max setting id")
	}
	return
}

// SettingsByRange gets user setting by range.
func (d *Dao) SettingsByRange(ctx context.Context, start, end int64) (res map[int64]map[int]int, err error) {
	var rows *xsql.Rows
	if rows, err = d.settingsByRangeStmt.Query(ctx, start, end); err != nil {
		log.Error("d.settingsStmt.Query(%d,%d) error(%v)", start, end, err)
		PromError("mysql:Settings")
		return
	}
	defer rows.Close()
	res = make(map[int64]map[int]int)
	for rows.Next() {
		var (
			mid int64
			v   string
			st  = make(map[int]int)
		)
		if err = rows.Scan(&mid, &v); err != nil {
			PromError("mysql:Settings scan")
			log.Error("d.Settings() Scan() error(%v)", err)
			return
		}
		if e := json.Unmarshal([]byte(v), &st); e != nil {
			log.Error("d.Settings() json unmarshal(%s) error(%v)", v, st)
			continue
		}
		res[mid] = st
	}
	err = rows.Err()
	return
}
