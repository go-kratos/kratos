package dao

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_taskJobKey = "task_job"
)

// SetnxTaskJob setnx task_job value
func (d *Dao) SetnxTaskJob(c context.Context, value string) (ok bool, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("SETNX", _taskJobKey, value)); err != nil {
		log.Error("d.SetnxMask(value:%s),error(%v)", value, err)
		return
	}
	return
}

// GetTaskJob .
func (d *Dao) GetTaskJob(c context.Context) (value string, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if value, err = redis.String(conn.Do("GET", _taskJobKey)); err != nil {
		log.Error("d.GetMaskJob,error(%v)", err)
		return
	}
	return
}

// GetSetTaskJob .
func (d *Dao) GetSetTaskJob(c context.Context, value string) (old string, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if old, err = redis.String(conn.Do("GETSET", _taskJobKey, value)); err != nil {
		log.Error("d.GetSetTaskJob(value:%s),error(%v)", value, err)
		return
	}
	return
}

// DelTaskJob .
func (d *Dao) DelTaskJob(c context.Context) (err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", _taskJobKey); err != nil {
		log.Error("d.DelTaskJob,error(%v)", err)
		return
	}
	return
}
