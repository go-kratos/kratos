package dao

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_maskJobKey = "mask_job"
)

// SetnxMaskJob setnx mask_job value
func (d *Dao) SetnxMaskJob(c context.Context, value string) (ok bool, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("SETNX", _maskJobKey, value)); err != nil {
		log.Error("d.SetnxMask(value:%s),error(%v)", value, err)
		return
	}
	return
}

// GetMaskJob .
func (d *Dao) GetMaskJob(c context.Context) (value string, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if value, err = redis.String(conn.Do("GET", _maskJobKey)); err != nil {
		log.Error("d.GetMaskJob,error(%v)", err)
		return
	}
	return
}

// GetSetMaskJob .
func (d *Dao) GetSetMaskJob(c context.Context, value string) (old string, err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if old, err = redis.String(conn.Do("GETSET", _maskJobKey, value)); err != nil {
		log.Error("d.GetSetMaskJob(value:%s),error(%v)", value, err)
		return
	}
	return
}

// DelMaskJob .
func (d *Dao) DelMaskJob(c context.Context) (err error) {
	var (
		conn = d.dmRds.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", _maskJobKey); err != nil {
		log.Error("d.DelMaskJob,error(%v)", err)
		return
	}
	return
}
