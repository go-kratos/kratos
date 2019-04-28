package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup-task/conf"
	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

//Dao dao.
type Dao struct {
	c          *conf.Config
	arcDB      *sql.DB
	arcReadDB  *sql.DB
	mngDB      *sql.DB
	upGroupURL string
	hbase      *hbase.Client
	redis      *redis.Pool
	hclient    *bm.Client
	//grpc
	acc account.AccountClient
}

var errCount = prom.BusinessErrCount

//New new dao
func New(config *conf.Config) (d *Dao) {
	d = &Dao{
		c:          config,
		arcDB:      sql.NewMySQL(config.DB.Archive),
		arcReadDB:  sql.NewMySQL(config.DB.ArchiveRead),
		mngDB:      sql.NewMySQL(config.DB.Manager),
		upGroupURL: config.Host.API + "/x/internal/uper/special/get",
		hbase:      hbase.NewClient(config.HBase.Config),
		redis:      redis.NewPool(config.Redis.Weight.Config),
		hclient:    bm.NewClient(config.HTTPClient),
	}
	var err error
	if d.acc, err = account.NewClient(config.GRPC.AccRPC); err != nil {
		panic(err)
	}
	return
}

//BeginTran begin transaction
func (d *Dao) BeginTran(ctx context.Context) (tx *sql.Tx, err error) {
	if tx, err = d.arcDB.Begin(ctx); err != nil {
		PromeErr("arcdb: begintran", "BeginTran d.arcDB.Begin error(%v)", err)
	}

	return
}

//Close close
func (d *Dao) Close() {
	if d.arcDB != nil {
		d.arcDB.Close()
	}
	if d.arcReadDB != nil {
		d.arcReadDB.Close()
	}
	if d.mngDB != nil {
		d.mngDB.Close()
	}
}

//Ping ping
func (d *Dao) Ping(ctx context.Context) (err error) {
	if d.arcDB != nil {
		if err = d.arcDB.Ping(ctx); err != nil {
			PromeErr("arcdb: ping", "d.arcDB.Ping error(%v)", err)
			return
		}
	}
	if d.arcReadDB != nil {
		if err = d.arcReadDB.Ping(ctx); err != nil {
			PromeErr("arcReaddb: ping", "d.arcReadDB.Ping error(%v)", err)
			return
		}
	}
	if d.mngDB != nil {
		if err = d.mngDB.Ping(ctx); err != nil {
			PromeErr("mngdb: ping", "d.mngDB.Ping error(%v)", err)
			return
		}
	}
	return
}

//AccountInfos get multi mids' accountinfo
func (d *Dao) AccountInfos(ctx context.Context, mids []int64) (info map[int64]*model.Info, err error) {
	tctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	var infosreply *account.InfosReply
	if infosreply, err = d.acc.Infos3(tctx, &account.MidsReq{Mids: mids}); err != nil {
		PromeErr("account_infos", "AccountInfos d.acc.Infos3 error(%v) mid(%d)", err, mids)
	}
	if infosreply != nil {
		info = infosreply.Infos
	}
	return
}

//PromeErr prome & log err
func PromeErr(name string, format string, args ...interface{}) {
	errCount.Incr(name)
	log.Error(format, args...)
}
