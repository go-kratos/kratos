package manager

import (
	"context"
	gosql "database/sql"
	"net/url"
	"strings"
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao/global"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
	"go-common/library/xstr"
)

const (
	//URLUNames url for names
	URLUNames = "/x/admin/manager/users/unames"
	//URLUids url for uids
	URLUids = "/x/admin/manager/users/uids"
)

// Dao is redis dao.
type Dao struct {
	c          *conf.Config
	managerDB  *sql.DB
	HTTPClient *bm.Client
	// cache tool
	cache *fanout.Fanout
	// mc
	mc *memcache.Pool
	// upSpecial expiration
	upSpecialExpire int32
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		managerDB: sql.NewMySQL(c.DB.Manager),
		// http client
		HTTPClient:      bm.NewClient(c.HTTPClient.Normal),
		mc:              memcache.NewPool(c.Memcache.Up),
		upSpecialExpire: int32(time.Duration(c.Memcache.UpSpecialExpire) / time.Second),
		cache:           global.GetWorker(),
	}
	return d
}

// Close fn
func (d *Dao) Close() {
	if d.managerDB != nil {
		d.managerDB.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.managerDB.Ping(c)
}

func prepareAndExec(c context.Context, db *sql.DB, sqlstr string, args ...interface{}) (res gosql.Result, err error) {
	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlstr)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, sqlstr)
		return
	}
	defer stmt.Close()

	res, err = stmt.Exec(c, args...)
	if err != nil {
		log.Error("data base fail, err=%v", err)
		return
	}
	return
}

func prepareAndQuery(c context.Context, db *sql.DB, sqlstr string, args ...interface{}) (rows *sql.Rows, err error) {
	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlstr)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, sqlstr)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(c, args...)
	if err != nil {
		log.Error("data base fail, err=%v", err)
		return
	}
	return
}

//GetUNamesByUids get uname by uid
func (d *Dao) GetUNamesByUids(c context.Context, uids []int64) (res map[int64]string, err error) {
	var param = url.Values{}
	var uidStr = xstr.JoinInts(uids)
	param.Set("uids", uidStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[int64]string `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+URLUNames, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+URLUNames+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+URLUNames+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

//GetUIDByNames get uid by uname
func (d *Dao) GetUIDByNames(c context.Context, names []string) (res map[string]int64, err error) {
	var param = url.Values{}
	var namesStr = strings.Join(names, ",")
	param.Set("unames", namesStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[string]int64 `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+URLUids, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+URLUids+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+URLUids+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}
