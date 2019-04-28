package dao

import (
	"context"
	"fmt"

	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"

	"github.com/bluele/gcache"
)

// Dao dao
type Dao struct {
	mc         *memcache.Pool
	db         *xsql.DB
	accountAPI account.AccountClient
	// local cache
	elecAVRankLC gcache.Cache
	elecUPRankLC gcache.Cache
}

// New init mysql db
func New() (dao *Dao) {
	dao = &Dao{
		mc:           memcache.NewPool(conf.Conf.Memcache),
		db:           xsql.NewMySQL(conf.Conf.MySQL),
		elecAVRankLC: gcache.New(conf.Conf.LocalCache.ElecAVRankSize).LFU().Build(),
		elecUPRankLC: gcache.New(conf.Conf.LocalCache.ElecUPRankSize).LFU().Build(),
	}
	var err error
	if dao.accountAPI, err = account.NewClient(conf.Conf.AccountGRPC); err != nil {
		panic(err)
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	return nil
}

func elecUPRankKey(upMID int64, ver int64) string {
	return fmt.Sprintf("ur_eur_%d_%d", ver, upMID)
}

func elecPrepUPRankKey(upMID int64, ver int64) string {
	return fmt.Sprintf("ur_epur_%d_%d", ver, upMID)
}

func elecAVRankKey(avID int64, ver int64) string {
	return fmt.Sprintf("ur_ear_%d_%d", ver, avID)
}

func elecPrepAVRankKey(avID int64, ver int64) string {
	return fmt.Sprintf("ur_epar_%d_%d", ver, avID)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	//mc: -key=elecUPRankKey -type=get
	CacheElecUPRank(c context.Context, mid int64) (*model.RankElecUPProto, error)
	//mc: -key=elecUPRankKey -expire=conf.Conf.CacheTTL.ElecUPRankTTL -encode=json
	AddCacheElecUPRank(c context.Context, mid int64, value *model.RankElecUPProto) error
	//mc: -key=elecUPRankKey
	DelCacheElecUPRank(c context.Context, mid int64) error

	//mc: -key=elecAVRankKey -type=get
	CacheElecAVRank(c context.Context, avID int64) (*model.RankElecAVProto, error)
	//mc: -key=elecAVRankKey -expire=conf.Conf.CacheTTL.ElecAVRankTTL -encode=json
	AddCacheElecAVRank(c context.Context, avID int64, value *model.RankElecAVProto) error
	//mc: -key=elecAVRankKey
	DelCacheElecAVRank(c context.Context, avID int64) error

	//mc: -key=elecPrepUPRankKey -type=get
	CacheElecPrepUPRank(c context.Context, mid int64) (*model.RankElecPrepUPProto, error)
	//mc: -key=elecPrepUPRankKey -expire=conf.Conf.CacheTTL.ElecPrepUPRankTTL -encode=json
	AddCacheElecPrepUPRank(c context.Context, mid int64, value *model.RankElecPrepUPProto) error
	//mc: -key=elecPrepUPRankKey
	DelCacheElecPrepUPRank(c context.Context, mid int64) error

	//mc: -key=elecPrepAVRankKey -type=get
	CacheElecPrepAVRank(c context.Context, avID int64) (*model.RankElecPrepAVProto, error)
	//mc: -key=elecPrepAVRankKey -expire=conf.Conf.CacheTTL.ElecPrepAVRankTTL -encode=json
	AddCacheElecPrepAVRank(c context.Context, avID int64, value *model.RankElecPrepAVProto) error
	//mc: -key=elecPrepAVRankKey
	DelCacheElecPrepAVRank(c context.Context, avID int64) error
}
