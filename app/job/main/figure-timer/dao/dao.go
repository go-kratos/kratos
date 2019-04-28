package dao

import (
	"context"
	"time"

	"go-common/app/job/main/figure-timer/conf"
	"go-common/app/job/main/figure-timer/model"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
)

// Int is dao interface
type Int interface {
	HBase
	Mysql
	Redis
	Close()
	Ping(c context.Context) error
}

// HBase dao hbase interface
type HBase interface {
	UserInfo(c context.Context, mid int64, weekVer int64) (userInfo *model.UserInfo, err error)
	PutCalcRecord(c context.Context, record *model.FigureRecord, weekTS int64) (err error)
	CalcRecords(c context.Context, mid int64, weekTSFrom, weekTSTo int64) (figureRecords []*model.FigureRecord, err error)
	ActionCounter(c context.Context, mid int64, ts int64) (counter *model.ActionCounter, err error)
}

// Mysql dao mysql interface
type Mysql interface {
	Figure(c context.Context, mid int64) (figure *model.Figure, err error)
	Figures(c context.Context, fromMid int64, limit int) (figures []*model.Figure, end bool, err error)
	UpsertFigure(c context.Context, figure *model.Figure) (id int64, err error)
	InsertRankHistory(c context.Context, rank *model.Rank) (id int64, err error)
	UpsertRank(c context.Context, rank *model.Rank) (id int64, err error)
}

// Redis dao redis interface
type Redis interface {
	FigureCache(c context.Context, mid int64) (figure *model.Figure, err error)
	SetFigureCache(c context.Context, figure *model.Figure) (err error)
	PendingMidsCache(c context.Context, version int64, shard int64) (mids []int64, err error)
	RemoveCache(c context.Context, mid int64) (err error)
}

// Dao struct info of Dao.
type Dao struct {
	c           *conf.Config
	mysql       *sql.DB
	hbase       *hbase.Client
	redis       *redis.Pool
	redisExpire int32
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		mysql:       sql.NewMySQL(c.Mysql),
		hbase:       hbase.NewClient(c.Hbase.Config),
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.mysql != nil {
		d.mysql.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.hbase != nil {
		d.hbase.Close()
	}
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.mysql.Ping(c); err != nil {
		return
	}
	// if err = d.hbase.Ping(c); err != nil {
	// 	return
	// }
	if err = d.PingRedis(c); err != nil {
		return
	}
	return
}
