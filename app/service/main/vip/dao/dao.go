package dao

import (
	"context"
	"time"

	"go-common/app/service/main/vip/conf"
	eleclient "go-common/app/service/main/vip/dao/ele-api-client"
	mailclient "go-common/app/service/main/vip/dao/mail-api-client"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Dao struct info of Dao.
type Dao struct {
	// mysql
	db    *sql.DB
	olddb *sql.DB
	// http
	client *bm.Client
	// ele api http client
	eleclient *eleclient.EleClient
	// mail api http client
	mailclient *mailclient.Client
	// conf
	c              *conf.Config
	msgURI         string
	payURI         string
	payCloseURL    string
	vipURI         string
	passportDetail string
	mc             *memcache.Pool
	mcExpire       int32
	errProm        *prom.Prom
	loginOutURL    string
	//redis pool
	redis *redis.Pool
	// cache async save
	cache *fanout.Fanout
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// db
		db:    sql.NewMySQL(c.Mysql),
		olddb: sql.NewMySQL(c.OldMysql),
		// http client.
		client:         bm.NewClient(c.HTTPClient),
		msgURI:         c.MsgURI + _SendUserNotify,
		payURI:         c.PayURI + _PayAPIAdd,
		payCloseURL:    c.Property.PayCoURL + _payClose,
		passportDetail: c.Property.PassportURL + _passportDetail,
		vipURI:         c.VipURI + _CleanCache,
		mc:             memcache.NewPool(c.Memcache.Config),
		mcExpire:       int32(time.Duration(c.Memcache.Expire) / time.Second),
		errProm:        prom.BusinessErrCount,
		loginOutURL:    c.Property.PassportURL + _loginout,
		redis:          redis.NewPool(c.Redis.Config),
		// cache chan
		cache: fanout.New("cache", fanout.Worker(10), fanout.Buffer(10240)),
	}
	// ele
	d.eleclient = eleclient.NewEleClient(c.ELEConf, d.client)
	d.mailclient = mailclient.NewClient(d.client)
	return
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		return
	}
	if err = d.olddb.Ping(c); err != nil {
		return
	}
	return d.db.Ping(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.olddb != nil {
		d.olddb.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
}

//StartTx start tx
func (d *Dao) StartTx(c context.Context) (tx *sql.Tx, err error) {
	if d.db != nil {
		tx, err = d.db.Begin(c)
	}
	return
}

//OldStartTx old start tx
func (d *Dao) OldStartTx(c context.Context) (tx *sql.Tx, err error) {
	if d.db != nil {
		tx, err = d.olddb.Begin(c)
	}
	return
}
