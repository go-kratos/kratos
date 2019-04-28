package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_taskReleaseSQL = "update task SET admin_id=0,state=0,uid=0,gtime=0 where state=? AND mtime<=?"
	_taskClearSQL   = "DELETE FROM task WHERE mtime<=? AND state>=3 LIMIT ?"
)

// TaskActiveConfigs list config
func (d *Dao) TaskActiveConfigs(c context.Context) (configs []*model.TaskConfig, err error) {
	db := d.orm.Model(&model.TaskConfig{}).Where("state=?", model.ConfigStateOn)

	if err = db.Find(&configs).Error; err != nil {
		log.Error("query error(%v)", err)
		return
	}
	return
}

// TaskActiveConsumer list consumer
func (d *Dao) TaskActiveConsumer(c context.Context) (consumerCache map[string]map[int64]struct{}, err error) {
	rows, err := d.orm.Table("task_consumer").Select("business_id,flow_id,uid").Where("state=?", model.ConsumerStateOn).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	consumerCache = make(map[string]map[int64]struct{})
	for rows.Next() {
		var bizID, flowID, UID int64
		if err = rows.Scan(&bizID, &flowID, &UID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}

		key := fmt.Sprintf("%d-%d", bizID, flowID)

		if _, ok := consumerCache[key]; ok {
			consumerCache[key][UID] = struct{}{}
		} else {
			consumerCache[key] = map[int64]struct{}{UID: {}}
		}
	}
	return
}

// KickOutConsumer 踢出用户
func (d *Dao) KickOutConsumer(c context.Context, bizid, flowid, uid int64) (err error) {
	return d.orm.Table("task_consumer").Where("business_id=? AND flow_id=? AND uid=?", bizid, flowid, uid).
		Update("state", model.ConsumerStateOff).Error
}

// Resource .
func (d *Dao) Resource(c context.Context, rid int64) (res *model.Resource, err error) {
	res = &model.Resource{}
	if err = d.orm.Where("id = ?", rid).First(res).Error; err == gorm.ErrRecordNotFound {
		res = nil
		err = nil
	}
	return
}

//RscState resource state
func (d *Dao) RscState(c context.Context, rid int64) (state int64, err error) {
	err = d.orm.Table("resource_result").Select("state").Where("rid=?", rid).Row().Scan(&state)
	return
}

// TaskRelease .
func (d *Dao) TaskRelease(c context.Context, mtime time.Time) (err error) {
	return d.orm.Exec(_taskReleaseSQL, model.TaskStateDispatch, mtime).Error
}

// ReleaseByConsumer .
func (d *Dao) ReleaseByConsumer(c context.Context, bizid, flowid, uid int64) (err error) {
	return d.orm.Table("task").Where("business_id=? AND flow_id=? AND uid=? AND (state=1 or (admin_id>0 AND state=0))", bizid, flowid, uid).Update(
		map[string]interface{}{
			"uid":      0,
			"state":    0,
			"gtime":    0,
			"admin_id": 0,
		}).Error
}

//Report .
func (d *Dao) Report(c context.Context, rt *model.Report) (err error) {
	return d.orm.Create(rt).Error
}

//TaskClear 已完成任务最多保留3天
func (d *Dao) TaskClear(c context.Context, mtime time.Time, limit int64) (rows int64, err error) {
	db := d.orm.Exec(_taskClearSQL, mtime, limit)
	rows, err = db.RowsAffected, db.Error
	return
}

//CheckFlow 检查资源是否在对应流程上
func (d *Dao) CheckFlow(c context.Context, rid, flowid int64) (ok bool, err error) {
	var id int64
	err = d.orm.Table("net_flow_resource").Select("id").
		Where("rid=? AND flow_id=? AND state!=-1", rid, flowid).Row().Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("CheckFlow(%d,%d) error(%v)", rid, flowid)
		}
		return
	}
	if id > 0 {
		ok = true
	}
	return
}

// CreateTask .
func (d *Dao) CreateTask(c context.Context, task *model.Task) error {
	return d.orm.Table("task").Where("rid=? AND flow_id=? AND state<?", task.RID, task.FlowID, model.TaskStateSubmit).
		Assign(map[string]interface{}{
			"mid":   task.MID,
			"fans":  task.Fans,
			"group": task.Group,
		}).FirstOrCreate(task).Error
}
