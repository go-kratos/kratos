package dao

import (
	"context"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"time"
)

// SetAppCovCache set apps and depts coverage into redis.
func (d *Dao) SetAppCovCache(c context.Context) (err error) {
	var (
		apps       []*ut.App
		deptApps   []*ut.App
		conn       = d.Redis.Get(c)
		expireTime = int64(time.Duration(d.c.Redis.ExpireTime) / time.Second)
	)
	defer conn.Close()
	if err = d.DB.Find(&apps).Error; err != nil {
		log.Error("d.AddAppCovCache DB.Find Error(%v)", err)
		return
	}
	if err = d.DB.Raw(`select substring_index(substring_index(path,"/",4),"/",-1) as path,avg(coverage) as coverage,max(mtime) as mtime from ut_app where has_ut=1 group by substring_index(substring_index(path,"/",4),"/",-1)`).
		Find(&deptApps).Error; err != nil {
		log.Error("d.AddAppCovCache DB.Find Error(%v)", err)
		return
	}
	apps = append(apps, deptApps...)
	for _, app := range apps {
		if err = conn.Send("SETEX", app.Path, expireTime, app.Coverage); err != nil {
			log.Error("d.AddAppCovCache Redis.SETEX Error(%v)", err)
			return
		}
	}
	return
}

// GetAppCovCache get apps and depts coverage from redis.
func (d *Dao) GetAppCovCache(c context.Context, path string) (coverage float64, err error) {
	var (
		conn  = d.Redis.Get(c)
		exist bool
		key   = path
	)
	defer conn.Close()
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("d.GetAppCov Error(%v)", err)
		return
	}
	if exist {
		if coverage, err = redis.Float64(conn.Do("GET", key)); err != nil {
			log.Error("d.GetAppCov Error(%v)", err)
			return
		}
	}
	return
}
