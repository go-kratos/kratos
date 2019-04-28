package vip

import (
	"context"
	"fmt"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/model"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"strconv"
)

// Dao vip dao
type Dao struct {
	c     *conf.Config
	db    *xsql.DB
	redis *redis.Pool
}

// New new vip dao
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		db:    xsql.NewMySQL(c.LiveUserMysql),
		redis: redis.NewPool(c.VipRedis),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
	d.redis.Close()
}

// initInfo init info struct
func (d *Dao) initInfo(info *model.VipInfo) *model.VipInfo {
	if info == nil {
		info = &model.VipInfo{Vip: 0, VipTime: model.TimeEmpty, Svip: 0, SvipTime: model.TimeEmpty}
	} else {
		if info.VipTime == "" {
			info.VipTime = model.TimeEmpty
		}
		if info.SvipTime == "" {
			info.SvipTime = model.TimeEmpty
		}
	}
	return info
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return nil
}

// toInt try trans interface input to int
func toInt(in interface{}) (int, error) {
	switch in.(type) {
	case int:
		return in.(int), nil
	case int32:
		return int(in.(int32)), nil
	case int64:
		return int(in.(int64)), nil
	case float32:
		return int(in.(float32)), nil
	case float64:
		return int(in.(float64)), nil
	case string:
		i, err := strconv.Atoi(in.(string))
		if err != nil {
			return 0, err
		}
		return i, nil
	case []byte:
		i, err := strconv.Atoi(string(in.([]byte)))
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("invalid input(%v)", in)
}
