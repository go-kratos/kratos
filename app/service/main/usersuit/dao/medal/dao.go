package medal

import (
	"context"
	"time"

	"go-common/app/service/main/usersuit/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/bluele/gcache"
)

const (
	_sendMsgPath           = "/api/notify/send.user.notify.do"
	_getWaredFansMedalPath = "/fans_medal/v1/fans_medal/get_weared_medal"
)

// Dao struct info of Dao.
type Dao struct {
	db *sql.DB

	c      *conf.Config
	client *bm.Client
	// memcache
	mc          *memcache.Pool
	mcExpire    int32
	pointExpire int32
	// send message URI.
	sendMsgURI string
	// get weared fans medal URI.
	getWaredFansMedalURI string
	// medalStore
	medalStore gcache.Cache
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     sql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
		// memcache
		mc:                   memcache.NewPool(c.Memcache.Config),
		mcExpire:             int32(time.Duration(c.Memcache.MedalExpire) / time.Second),
		pointExpire:          int32(time.Duration(c.Memcache.PointExpire) / time.Second),
		sendMsgURI:           c.Host.MessageCo + _sendMsgPath,
		getWaredFansMedalURI: c.Host.LiveAPICo + _getWaredFansMedalPath,
		medalStore:           gcache.New(c.MedalCache.Size).LFU().Build(),
	}
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
