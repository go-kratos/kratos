package dao

import (
	"context"

	"go-common/app/admin/main/cache/model"
	"go-common/library/log"
)

var (
	opsMcURI    = "http://ops-cache.bilibili.co/manager/redisapp/get_all_memcache_json"
	opsRedisURI = "http://ops-cache.bilibili.co/manager/redisapp/get_all_redis_json"
)

// OpsMemcaches get all ops mc.
func (d *Dao) OpsMemcaches(c context.Context) (mcs []*model.OpsCacheMemcache, err error) {
	var res struct {
		Data []*model.OpsCacheMemcache `json:"data"`
	}
	if err = d.client.Get(c, opsMcURI, "", nil, &res); err != nil {
		log.Error("ops memcache url(%s) error(%v)", opsMcURI, err)
		return
	}
	mcs = res.Data
	return
}

// OpsRediss get all ops redis.
func (d *Dao) OpsRediss(c context.Context) (mcs []*model.OpsCacheRedis, err error) {
	var res struct {
		Data []*model.OpsCacheRedis `json:"data"`
	}
	if err = d.client.Get(c, opsRedisURI, "", nil, &res); err != nil {
		log.Error("ops redis url(%s) error(%v)", opsRedisURI, err)
		return
	}
	mcs = res.Data
	return
}
