package dao

import (
	"context"
	"time"

	antispam "go-common/app/service/main/antispam/rpc/client"
	"go-common/app/service/main/filter/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"

	bm "go-common/library/net/http/blademaster"
)

// Dao struct .
type Dao struct {
	conf        *conf.Config
	mysql       *sql.DB
	hbase       *hbase.Client
	mc          *memcache.Pool
	antispamRPC *antispam.Client
	httpClient  *bm.Client

	// expire
	mcKeyFilterExp int32
	mcFilterExp    int32
	// AIscore
	aiScoreURL       string
	mngReplyDelURL   string
	mngReplyLabelURL string
}

// New .
func New(c *conf.Config) *Dao {
	d := &Dao{
		conf:             c,
		mcKeyFilterExp:   int32(time.Duration(c.Memcache.Expire.FilterKeyExpire) / time.Second),
		mcFilterExp:      int32(time.Duration(c.Memcache.Expire.FilterExpire) / time.Second),
		antispamRPC:      antispam.NewClient(nil),
		httpClient:       bm.NewClient(c.HTTPClient),
		aiScoreURL:       c.Property.AIHost.AI + _getScore,
		mngReplyDelURL:   c.Property.AIHost.Manager + _replyDel,
		mngReplyLabelURL: c.Property.AIHost.Manager + _replyLabel,
	}
	// mysql
	d.mysql = sql.NewMySQL(c.MySQL)
	// mc
	d.mc = memcache.NewPool(c.Memcache.Mc)
	// hbase
	if c.HBase != nil {
		d.hbase = hbase.NewClient(c.HBase.Config)
	}
	return d
}

// Ping .
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.mysql.Ping(c); err != nil {
		return
	}
	return d.PingMc(c)
}

// PingMc .
func (d *Dao) PingMc(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{
		Key:        "ping",
		Object:     1,
		Flags:      memcache.FlagJSON,
		Expiration: 60,
	}
	err = conn.Set(item)
	conn.Close()
	return
}

// Close .
func (d *Dao) Close() {
	if d.mysql != nil {
		d.mysql.Close()
	}
}

// BeginTran .
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.mysql.Begin(c)
}
