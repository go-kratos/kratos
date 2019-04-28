package dao

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go-common/app/job/main/history/conf"
	"go-common/app/service/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/database/tidb"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/database/hbase.v2"
)

var errFlushRequest = errors.New("error flush history to store")

// Dao dao.
type Dao struct {
	conf           *conf.Config
	HTTPClient     *bm.Client
	URL            string
	info           *hbase.Client
	redis          *redis.Pool
	db             *tidb.DB
	longDB         *tidb.DB
	insertStmt     *tidb.Stmts
	businessesStmt *tidb.Stmts
	allHisStmt     *tidb.Stmts
	delUserStmt    *tidb.Stmts
	BusinessesMap  map[int64]*model.Business
	BusinessNames  map[string]*model.Business
}

// New new history dao and return.
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		conf:       c,
		redis:      redis.NewPool(c.Redis),
		HTTPClient: bm.NewClient(c.Job.Client),
		URL:        c.Job.URL,
		info:       hbase.NewClient(c.Info.Config),
		db:         tidb.NewTiDB(c.TiDB),
		longDB:     tidb.NewTiDB(c.LongTiDB),
	}
	dao.businessesStmt = dao.db.Prepared(_businessesSQL)
	dao.insertStmt = dao.db.Prepared(_addHistorySQL)
	dao.allHisStmt = dao.db.Prepared(_allHisSQL)
	dao.delUserStmt = dao.db.Prepared(_delUserHisSQL)
	dao.loadBusiness()
	go dao.loadBusinessproc()
	return
}

// Flush flush history to store by mids.
func (d *Dao) Flush(c context.Context, mids string, stime int64) (err error) {
	params := url.Values{}
	params.Set("mids", mids)
	params.Set("time", fmt.Sprintf("%d", stime))
	var res = &struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}{}
	if err = d.HTTPClient.Post(c, d.URL, "", params, res); err != nil {
		log.Error("d.HTTPClient.Post(%s?%s) error(%v)", d.URL, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.HTTPClient.Post(%s?%s) code:%d msg:%s", d.URL, params.Encode(), res.Code, res.Msg)
		err = errFlushRequest
		return
	}
	return
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// Close close the redis and kafka resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
	d.longDB.Close()
}

func (d *Dao) loadBusiness() {
	var business []*model.Business
	var err error
	businessMap := make(map[string]*model.Business)
	businessIDMap := make(map[int64]*model.Business)
	for {
		if business, err = d.Businesses(context.TODO()); err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, b := range business {
			businessMap[b.Name] = b
			businessIDMap[b.ID] = b
		}
		d.BusinessNames = businessMap
		d.BusinessesMap = businessIDMap
		return
	}
}

func (d *Dao) loadBusinessproc() {
	for {
		time.Sleep(time.Minute * 5)
		d.loadBusiness()
	}
}
