package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/credit/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

const (
	_sendMsgURI      = "/api/notify/send.user.notify.do"
	_getQSURI        = "/laogai/question"
	_replyCountURI   = "/x/internal/v2/reply/mcount"
	_sendMedalURI    = "/x/internal/medal/grant"
	_managersURI     = "/x/admin/manager/users"
	_managerTotalURI = "/x/admin/manager/users/total"
	_addAppealURI    = "/x/internal/workflow/appeal/add"
	_appealDetailURI = "/x/internal/workflow/appeal/info"
	_appealListURI   = "/x/internal/workflow/appeal/list"
)

// Dao struct info of Dao.
type Dao struct {
	// mysql
	db *sql.DB
	// memcache
	mc              *memcache.Pool
	userExpire      int32
	minCommonExpire int32
	commonExpire    int32
	// redis
	redis       *redis.Pool
	redisExpire int64
	// databus stat
	dbusLabour *databus.Databus
	// http
	client *bm.Client
	// conf
	c *conf.Config
	// message api
	sendMsgURL string
	// big data api
	getQSURL string
	// account.co api
	sendMedalURL    string
	replyCountURL   string
	managersURL     string
	managerTotalURL string
	// appeal
	addAppealURL    string
	appealDetailURL string
	appealListURL   string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// db
		db: sql.NewMySQL(c.Mysql),
		// memcache
		mc:              memcache.NewPool(c.Memcache.Config),
		userExpire:      int32(time.Duration(c.Memcache.UserExpire) / time.Second),
		minCommonExpire: int32(time.Duration(c.Memcache.MinCommonExpire) / time.Second),
		commonExpire:    int32(time.Duration(c.Memcache.CommonExpire) / time.Second),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int64(time.Duration(c.Redis.Expire) / time.Second),
		// databus
		dbusLabour: databus.New(c.DataBus),
		// http client
		client: bm.NewClient(c.HTTPClient),
	}
	d.sendMsgURL = c.Host.MessageURI + _sendMsgURI
	d.getQSURL = c.Host.BigDataURI + _getQSURI
	d.sendMedalURL = c.Host.APICoURI + _sendMedalURI
	d.replyCountURL = c.Host.APICoURI + _replyCountURI
	d.managersURL = c.Host.ManagersURI + _managersURI
	d.managerTotalURL = c.Host.ManagersURI + _managerTotalURI
	d.addAppealURL = c.Host.APICoURI + _addAppealURI
	d.appealDetailURL = c.Host.APICoURI + _appealDetailURI
	d.appealListURL = c.Host.APICoURI + _appealListURI
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingMC(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin mysql transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}
