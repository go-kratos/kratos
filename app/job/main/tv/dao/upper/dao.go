package upper

import (
	"context"
	"time"

	"go-common/app/job/main/tv/conf"
	acccli "go-common/app/service/main/account/api"
	accwar "go-common/app/service/main/account/api"
	account "go-common/app/service/main/account/model"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/stat/prom"

	"github.com/pkg/errors"
)

var (
	missedCount = prom.CacheMiss
	cachedCount = prom.CacheHit
)

// Dao is account dao.
type Dao struct {
	accClient accwar.AccountClient
	mc        *memcache.Pool
	mcExpire  int32
	DB        *sql.DB
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		DB:       sql.NewMySQL(c.Mysql),
	}
	var err error
	if d.accClient, err = acccli.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}

// Card3 get card info by mid
func (d *Dao) Card3(c context.Context, mid int64) (res *account.Card, err error) {
	var (
		arg  = &accwar.MidReq{Mid: mid}
		resp *accwar.CardReply
	)
	if resp, err = d.accClient.Card3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	res = resp.Card
	return
}
