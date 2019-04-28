package dao

import (
	"context"
	"time"

	"go-common/app/job/main/passport-user/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	originDB  *sql.DB
	userDB    *sql.DB
	encryptDB *sql.DB
	mc        *memcache.Pool
	mcExpire  int32
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		originDB:  sql.NewMySQL(c.DB.OriginDB),
		userDB:    sql.NewMySQL(c.DB.UserDB),
		encryptDB: sql.NewMySQL(c.DB.EncryptDB),
		mc:        memcache.NewPool(c.Memcache.Config),
		mcExpire:  int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	return
}

// Ping check dao ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.originDB.Ping(c); err != nil {
		log.Info("dao.originDB.Ping() error(%v)", err)
	}
	if err = d.userDB.Ping(c); err != nil {
		log.Info("dao.userDB.Ping() error(%v)", err)
	}
	if err = d.encryptDB.Ping(c); err != nil {
		log.Info("dao.encryptDB.Ping() error(%v)", err)
	}
	if err = d.pingMC(c); err != nil {
		log.Info("d.pingMC() error(%v)", err)
	}
	return
}

// Close close connections of mc, cloudDB.
func (d *Dao) Close() (err error) {
	if d.originDB != nil {
		d.originDB.Close()
	}
	if d.userDB != nil {
		d.userDB.Close()
	}
	if d.encryptDB != nil {
		d.encryptDB.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	return
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.userDB.Begin(c)
}
