package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Dao struct .
type Dao struct {
	conf  *conf.Config
	mysql *sql.DB
	hbase *hbase.Client
	mc    *memcache.Pool
	// http
	client *bm.Client

	// key stmt
	ConKeyStmt        *sql.Stmt
	ConKeyIDStmt      *sql.Stmt
	insertConkeyStmt  *sql.Stmt
	inserKeyStmt      *sql.Stmt
	delKeyFidStmt     *sql.Stmt
	updateConKeyStmt  *sql.Stmt
	searchKeyAreaStmt *sql.Stmt
	insertFkLogStmt   *sql.Stmt
	fkLogsStmt        *sql.Stmt

	// expire
	mcExp int32

	// AIscore
	aiScoreURL string
}

// New conf .
func New(c *conf.Config) *Dao {
	d := &Dao{
		conf:  c,
		mcExp: int32(time.Duration(c.Memcache.Expire.Expire) / time.Second),
		// http client
		client:     bm.NewClient(c.HTTPClient.Normal),
		aiScoreURL: c.Host.AI + _getScore,
	}
	// mysql
	d.mysql = sql.NewMySQL(c.MySQL)

	// mc
	d.mc = memcache.NewPool(c.Memcache.Mc)

	// hbase
	if c.HBase != nil {
		d.hbase = hbase.NewClient(c.HBase.Config)
	}

	// key
	d.ConKeyStmt = d.mysql.Prepared(_conKey)
	d.ConKeyIDStmt = d.mysql.Prepared(_conKeyByID)
	d.insertConkeyStmt = d.mysql.Prepared(_insertConkey)
	d.inserKeyStmt = d.mysql.Prepared(_inserKey)
	d.delKeyFidStmt = d.mysql.Prepared(_delKeyFid)
	d.updateConKeyStmt = d.mysql.Prepared(_updateConkey)
	d.searchKeyAreaStmt = d.mysql.Prepared(_searchKeyArea)
	d.insertFkLogStmt = d.mysql.Prepared(_insertFkLog)
	d.fkLogsStmt = d.mysql.Prepared(_fkLogs)

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
	if d.hbase != nil {
		d.hbase.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
}

// BeginTran .
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	if tx, err = d.mysql.Begin(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}
