package dao

import (
	"context"
	"time"

	"go-common/app/service/main/history/conf"
	"go-common/app/service/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/database/tidb"
	"go-common/library/queue/databus"
)

// Dao dao
type Dao struct {
	c                     *conf.Config
	tidb                  *tidb.DB
	redis                 *redis.Pool
	redisExpire           int32
	mergeDbus             *databus.Databus
	businessesStmt        *tidb.Stmts
	historiesStmt         *tidb.Stmts
	historyStmt           *tidb.Stmts
	insertStmt            *tidb.Stmts
	deleteHistoriesStmt   *tidb.Stmts
	clearAllHistoriesStmt *tidb.Stmts
	userHideStmt          *tidb.Stmts
	updateUserHideStmt    *tidb.Stmts
	Businesses            map[int64]*model.Business
	BusinessNames         map[string]*model.Business
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:           c,
		redis:       redis.NewPool(c.Redis.Config),
		tidb:        tidb.NewTiDB(c.TiDB),
		mergeDbus:   databus.New(c.DataBus.Merge),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	dao.businessesStmt = dao.tidb.Prepared(_businessesSQL)
	dao.historiesStmt = dao.tidb.Prepared(_historiesSQL)
	dao.deleteHistoriesStmt = dao.tidb.Prepared(_deleteHistoriesSQL)
	dao.clearAllHistoriesStmt = dao.tidb.Prepared(_clearAllHistoriesSQL)
	dao.historyStmt = dao.tidb.Prepared(_historySQL)
	dao.userHideStmt = dao.tidb.Prepared(_userHide)
	dao.updateUserHideStmt = dao.tidb.Prepared(_updateUserHide)
	dao.insertStmt = dao.tidb.Prepared(_addHistorySQL)

	dao.loadBusiness()
	go dao.loadBusinessproc()
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.tidb.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.tidb.Ping(c); err != nil {
		return
	}
	return d.pingRedis(c)
}

// LoadBusiness .
func (d *Dao) loadBusiness() {
	var business []*model.Business
	var err error
	businessMap := make(map[string]*model.Business)
	businessIDMap := make(map[int64]*model.Business)
	for {
		if business, err = d.QueryBusinesses(context.TODO()); err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, b := range business {
			businessMap[b.Name] = b
			businessIDMap[b.ID] = b
		}
		d.BusinessNames = businessMap
		d.Businesses = businessIDMap
		return
	}
}

func (d *Dao) loadBusinessproc() {
	for {
		time.Sleep(time.Minute * 5)
		d.loadBusiness()
	}
}
