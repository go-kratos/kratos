package bws

import (
	"context"
	"time"

	"go-common/app/interface/main/activity/conf"
	bwsmdl "go-common/app/interface/main/activity/model/bws"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_bindingSQL  = "UPDATE act_bws_users SET mid = ? WHERE `key`= ?"
	_usersMidSQL = "SELECT id,mid,`key`,ctime,mtime,bid FROM act_bws_users WHERE mid = ?"
	_usersKeySQL = "SELECT id,mid,`key`,ctime,mtime,bid FROM act_bws_users WHERE `key` = ?"
	_usersIDSQL  = "SELECT id,mid,`key`,ctime,mtime,bid FROM act_bws_users WHERE id = ?"
)

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Dao dao.
type Dao struct {
	// config
	c *conf.Config
	// db
	db              *xsql.DB
	mc              *memcache.Pool
	mcExpire        int32
	redis           *redis.Pool
	redisExpire     int32
	userAchExpire   int32
	userPointExpire int32
	achCntExpire    int32
	cacheCh         chan func()
}

// New dao new.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c:               c,
		db:              xsql.NewMySQL(c.MySQL.Like),
		mc:              memcache.NewPool(c.Memcache.Like),
		mcExpire:        int32(time.Duration(c.Memcache.BwsExpire) / time.Second),
		redis:           redis.NewPool(c.Redis.Config),
		cacheCh:         make(chan func(), 1024),
		redisExpire:     int32(time.Duration(c.Redis.Expire) / time.Second),
		userAchExpire:   int32(time.Duration(c.Redis.UserAchExpire) / time.Second),
		userPointExpire: int32(time.Duration(c.Redis.UserPointExpire) / time.Second),
		achCntExpire:    int32(time.Duration(c.Redis.AchCntExpire) / time.Second),
	}
	return
}

// RawUsersMid get users by mid
func (d *Dao) RawUsersMid(c context.Context, mid int64) (res *bwsmdl.Users, err error) {
	res = &bwsmdl.Users{}
	row := d.db.QueryRow(c, _usersMidSQL, mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Key, &res.Ctime, &res.Mtime, &res.Bid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("RawUsersMid:row.Scan error(%v)", err)
		}
	}
	return
}

// RawUsersKey get users by key
func (d *Dao) RawUsersKey(c context.Context, key string) (res *bwsmdl.Users, err error) {
	res = &bwsmdl.Users{}
	row := d.db.QueryRow(c, _usersKeySQL, key)
	if err = row.Scan(&res.ID, &res.Mid, &res.Key, &res.Ctime, &res.Mtime, &res.Bid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("RawUsersKey:row.Scan error(%v)", err)
		}
	}
	return
}

// Binding binding mid
func (d *Dao) Binding(c context.Context, loginMid int64, p *bwsmdl.ParamBinding) (err error) {
	if _, err = d.db.Exec(c, _bindingSQL, loginMid, p.Key); err != nil {
		log.Error("Binding: db.Exec(%d,%s) error(%v)", loginMid, p.Key, err)
	}
	return
}

// UserByID .
func (d *Dao) UserByID(c context.Context, keyID int64) (res *bwsmdl.Users, err error) {
	res = &bwsmdl.Users{}
	row := d.db.QueryRow(c, _usersIDSQL, keyID)
	if err = row.Scan(&res.ID, &res.Mid, &res.Key, &res.Ctime, &res.Mtime, &res.Bid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("UserByID:row.Scan error(%v)", err)
		}
	}
	return
}
