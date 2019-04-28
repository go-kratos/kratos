package dao

import (
	"context"
	"time"

	"go-common/app/job/main/credit/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const (
	_delReplyURI       = "/x/internal/v2/reply/del"
	_delTagURI         = "/x/internal/tag/archive/del"
	_regReplyURI       = "/x/internal/v2/reply/subject/regist"
	_blockAccountURI   = "/x/internal/block/block"
	_unBlockAccountURI = "/x/internal/block/remove"
	_sendPendantURI    = "/x/internal/pendant/multiGrantByPid"
	_sendMsgURI        = "/api/notify/send.user.notify.do"
	_delDMURI          = "/x/internal/dmadmin/report/judge/result"
	_sendMedalURI      = "/x/internal/medal/grant"
	_addMoralURI       = "/api/moral/add"
	_upReplyStateURI   = "/x/internal/v2/reply/report/state"
	_modifyCoinsURI    = "/x/internal/v1/coin/user/modify"
	_filterURI         = "/x/internal/filter"
	_upAppealStateURI  = "/x/internal/workflow/appeal/v3/public/referee"
)

// Dao struct info of Dao.
type Dao struct {
	db     *sql.DB
	c      *conf.Config
	client *bm.Client
	// del path
	delReplyURL       string
	delTagURL         string
	delDMURL          string
	blockAccountURL   string
	unBlockAccountURL string
	regReplyURL       string
	sendPendantURL    string
	sendMsgURL        string
	sendMedalURL      string
	addMoralURL       string
	upReplyStateURL   string
	modifyCoinsURL    string
	filterURL         string
	upAppealStateURL  string
	// redis 	// redis
	redis       *redis.Pool
	redisExpire int64
	// memcache
	mc *memcache.Pool
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		db:                sql.NewMySQL(c.Mysql),
		client:            bm.NewClient(c.HTTPClient),
		delReplyURL:       c.Host.APICoURI + _delReplyURI,
		delTagURL:         c.Host.APICoURI + _delTagURI,
		regReplyURL:       c.Host.APICoURI + _regReplyURI,
		blockAccountURL:   c.Host.APICoURI + _blockAccountURI,
		unBlockAccountURL: c.Host.APICoURI + _unBlockAccountURI,
		sendPendantURL:    c.Host.APICoURI + _sendPendantURI,
		sendMsgURL:        c.Host.MsgCoURI + _sendMsgURI,
		delDMURL:          c.Host.APICoURI + _delDMURI,
		sendMedalURL:      c.Host.APICoURI + _sendMedalURI,
		addMoralURL:       c.Host.AccountCoURI + _addMoralURI,
		upReplyStateURL:   c.Host.APICoURI + _upReplyStateURI,
		modifyCoinsURL:    c.Host.APICoURI + _modifyCoinsURI,
		filterURL:         c.Host.APICoURI + _filterURI,
		upAppealStateURL:  c.Host.APICoURI + _upAppealStateURI,
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int64(time.Duration(c.Redis.Expire) / time.Second),
		// memcache
		mc: memcache.NewPool(c.Memcache.Config),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
