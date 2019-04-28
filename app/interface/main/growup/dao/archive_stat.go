package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/interface/main/growup/model"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ArticleStat article stat
func (d *Dao) ArticleStat(c context.Context, mid int64, ip string) (res article.UpStat, err error) {
	arg := &article.ArgMid{Mid: mid, RealIP: ip}
	if res, err = d.art.CreationUpStat(c, arg); err != nil {
		log.Error("d.art.CreationUpStat(%+v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

func (d *Dao) getUpBaseStatCache(c context.Context, mid int64, date string) (data *model.UpBaseStat, err error) {
	key := fmt.Sprintf("growup-up-status:%d-%s", mid, date)
	res, err := d.getCacheVal(c, key)
	if err != nil {
		log.Error("d.getCacheVal error(%v)", err)
		return
	}
	if res == nil {
		return
	}
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", res, err)
	}
	return
}

// setUpBaseStatCache add stat cache.
func (d *Dao) setUpBaseStatCache(c context.Context, mid int64, date string, st *model.UpBaseStat) (err error) {
	key := fmt.Sprintf("growup-up-status:%d-%s", mid, date)
	v, err := json.Marshal(st)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	return d.setCacheKV(c, key, v, d.redisExpire)
}

// UpStat get up stat from hbase
func (d *Dao) UpStat(c context.Context, mid int64, dt string) (st *model.UpBaseStat, err error) {
	// try cache
	st, err = d.getUpBaseStatCache(c, mid, dt)
	if err != nil {
		log.Error("d.getUpBaseStatCache(%d) error(%v)", mid, err)
		return
	}
	if st != nil {
		return
	}
	// from hbase
	if st, err = d.BaseUpStat(c, mid, dt); st != nil {
		d.AddCache(func() {
			d.setUpBaseStatCache(context.TODO(), mid, dt, st)
		})
	}
	return
}
