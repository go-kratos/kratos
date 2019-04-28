package dao

import (
	"context"

	"go-common/app/admin/main/push/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_dpConditionSQL             = `select id,job,task,conditions,sql_stmt,status,status_url,file from push_dataplatform_conditions where job=?`
	_addDpConditionSQL          = `insert into push_dataplatform_conditions (job,task,conditions,sql_stmt,status,status_url,file) values (?,?,?,?,?,?,?) on duplicate key update task=?,conditions=?,sql_stmt=?,status=?,status_url=?,file=?`
	_updateDpConditionStatusSQL = `update push_dataplatform_conditions set status=? where job=?`
)

// AddDPCondition add data platform task
func (d *Dao) AddDPCondition(ctx context.Context, cond *model.DPCondition) (id int64, err error) {
	res, err := d.db.Exec(ctx, _addDpConditionSQL, cond.Job, cond.Task, cond.Condition, cond.SQL, cond.Status, cond.StatusURL, cond.File,
		cond.Task, cond.Condition, cond.SQL, cond.Status, cond.StatusURL, cond.File)
	if err != nil {
		log.Error("d.AddDPCondition(%+v) error(%v)", cond, err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// DPCondition .
func (d *Dao) DPCondition(ctx context.Context, job string) (c *model.DPCondition, err error) {
	c = new(model.DPCondition)
	if err = d.db.QueryRow(ctx, _dpConditionSQL, job).Scan(&c.ID, &c.Job, &c.Task, &c.Condition, &c.SQL, &c.Status, &c.StatusURL, &c.File); err != nil {
		if err == sql.ErrNoRows {
			c = nil
			err = nil
		}
		return
	}
	return
}

// UpdateDpCondtionStatus .
func (d *Dao) UpdateDpCondtionStatus(ctx context.Context, job string, status int) (err error) {
	_, err = d.db.Exec(ctx, _updateDpConditionStatusSQL, status, job)
	return
}

// Partitions 获取一级分区数据
func (d *Dao) Partitions(ctx context.Context) (m map[int]string, err error) {
	var res = struct {
		Code int `json:"code"`
		Data map[int]struct {
			ID   int    `json:"id"`
			Pid  int    `json:"pid"`
			Name string `json:"name"`
		} `json:"data"`
	}{}
	if err = d.httpClient.Get(ctx, d.c.Cfg.PartitionsURL, "", nil, &res); err != nil {
		return
	}
	if !ecode.Int(res.Code).Equal(ecode.OK) {
		err = ecode.Int(res.Code)
		return
	}
	m = make(map[int]string)
	for _, v := range res.Data {
		if v.Pid != 0 {
			continue
		}
		m[v.ID] = v.Name
	}
	return
}
