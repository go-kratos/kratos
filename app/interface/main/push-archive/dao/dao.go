package dao

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	xredis "go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"

	"go-common/library/database/hbase.v2"
)

// Dao .
type Dao struct {
	c                         *conf.Config
	db                        *xsql.DB
	redis                     *xredis.Pool
	relationHBase             *hbase.Client
	relationHBaseReadTimeout  time.Duration
	relationHBaseWriteTimeout time.Duration
	fanHBase                  *hbase.Client
	fanHBaseReadTimeout       time.Duration
	httpClient                *xhttp.Client
	settingStmt               *xsql.Stmt
	setSettingStmt            *xsql.Stmt
	settingsMaxIDStmt         *xsql.Stmt
	setStatisticsStmt         *xsql.Stmt
	UpperLimitExpire          int32
	FanGroups                 map[string]*FanGroup
	GroupOrder                []string
	Proportions               []Proportion
	ActiveDefaultTime         map[int]int
	PushBusinessID            string
	PushAuth                  string
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// New creates a push-service DAO instance.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:                         c,
		db:                        xsql.NewMySQL(c.MySQL),
		relationHBase:             hbase.NewClient(&c.HBase.Config),
		relationHBaseReadTimeout:  time.Duration(c.HBase.ReadTimeout),
		relationHBaseWriteTimeout: time.Duration(c.HBase.WriteTimeout),
		fanHBase:                  hbase.NewClient(&c.FansHBase.Config),
		fanHBaseReadTimeout:       time.Duration(c.FansHBase.ReadTimeout),
		redis:                     xredis.NewPool(c.Redis),
		httpClient:                xhttp.NewClient(c.HTTPClient),
		UpperLimitExpire:          int32(time.Duration(c.ArcPush.UpperLimitExpire) / time.Second),
		FanGroups:                 NewFanGroups(c),
		Proportions:               NewProportion(c.ArcPush.Proportions),
	}
	d.settingStmt = d.db.Prepared(_settingSQL)
	d.setSettingStmt = d.db.Prepared(_setSettingSQL)
	d.settingsMaxIDStmt = d.db.Prepared(_settingsMaxIDSQL)
	d.setStatisticsStmt = d.db.Prepared(_inStatisticsSQL)
	for _, gp := range c.ArcPush.Order {
		if _, exist := d.FanGroups[gp]; !exist {
			log.Error("order config error, group %s not exist", gp)
			fmt.Printf("order config error, group %s not exist\r\n\r\n", gp)
			os.Exit(1)
		}
	}
	d.GroupOrder = c.ArcPush.Order
	// default active time
	d.ActiveDefaultTime = map[int]int{}
	for _, one := range c.ArcPush.ActiveTime {
		d.ActiveDefaultTime[one] = 1
	}
	return d
}

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

// PromInfoAdd add prom info by value
func PromInfoAdd(name string, value int64) {
	infosCount.Add(name, value)
}

// PromChanLen channel length
func PromChanLen(name string, length int64) {
	infosCount.State(name, length)
}

// BeginTx begin transaction.
func (d *Dao) BeginTx(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// Close dao.
func (d *Dao) Close() (err error) {
	if err = d.relationHBase.Close(); err != nil {
		log.Error("d.relationHBase.Close() error(%v)", err)
		PromError("hbase:close")
	}
	if err = d.fanHBase.Close(); err != nil {
		log.Error("d.fanHBase.Close() error(%v)", err)
		PromError("fanHBase:close")
	}
	if err = d.redis.Close(); err != nil {
		log.Error("d.redis.Close() error(%v)", err)
		PromError("redis:close")
	}
	if err = d.db.Close(); err != nil {
		log.Error("d.db.Close() error(%v)", err)
		PromError("db:close")
	}
	return
}

// Ping check connection status.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		PromError("mysql:Ping")
		log.Error("d.db.Ping error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		PromError("redis:Ping")
		log.Error("d.redis.Ping error(%v)", err)
	}
	return
}

// Batch 批量处理
func Batch(list *[]int64, batchSize int, retry int, params *model.BatchParam, f func(fans *[]int64, params map[string]interface{}) error) {
	if params == nil {
		log.Warn("Batch params(%+v) nil", params)
		return
	}
	for {
		var (
			mids []int64
			err  error
		)
		l := len(*list)
		if l == 0 {
			break
		} else if l <= batchSize {
			mids = (*list)[:l]
		} else {
			mids = (*list)[:batchSize]
			l = batchSize
		}
		*list = (*list)[l:]

		params.Handler(&params.Params, mids)
		for i := 0; i < retry; i++ {
			if err = f(&mids, params.Params); err == nil {
				break
			}
		}
		if err != nil {
			log.Error("Batch error(%v), params(%+v)", err, params)
		}
	}
}
