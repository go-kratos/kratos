package dao

import (
	"context"

	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"math/rand"
	"time"
)

type Dao struct {
	c             *conf.Config
	mc            *memcache.Pool
	db            *xsql.DB
	redis         *redis.Pool
	cacheExpire   int32
	httpClient    *httpx.Client
	changeDataBus *databus.Databus
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:             c,
		mc:            memcache.NewPool(c.Memcache.Wallet),
		db:            xsql.NewMySQL(c.DB.Wallet),
		redis:         redis.NewPool(c.Redis.Wallet),
		cacheExpire:   c.WalletExpire,
		httpClient:    httpx.NewClient(c.HTTPClient),
		changeDataBus: databus.New(c.DataBus.Change),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
	d.redis.Close()
	d.changeDataBus.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("PingDb error(%v)", err)
		return
	}
	if err = d.pingMC(c); err != nil {
		return err
	}
	return d.PingRedis(c)

}

// pingMc ping
func (d *Dao) pingMC(c context.Context) (err error) {
	item := &memcache.Item{
		Key:        "ping",
		Value:      []byte{1},
		Expiration: d.cacheExpire,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		log.Error("PingMemcache conn.Set(%v) error(%v)", item, err)
	}
	return
}

func (d *Dao) PingRedis(c context.Context) (err error) {
	var conn = d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

//　通过提供的sql和bind来更新，没有实际业务意义，只是为了少写重复代码
func execSqlWithBindParams(d *Dao, c context.Context, sql *string, bindParams ...interface{}) (affect int64, err error) {
	res, err := d.db.Exec(c, *sql, bindParams...)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", *sql, err)
		return
	}
	return res.RowsAffected()
}

func randomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (d *Dao) GetDetailByCache(c context.Context, uid int64) (wallet *model.Detail, err error) {
	mcDetail, err := d.WalletCache(c, uid)
	if err == ecode.ServerErr {
		return
	}
	if err == nil {
		if d.IsNewVersion(c, mcDetail) {
			wallet = mcDetail.Detail
			return
		} else {
			log.Info("user wallet hit but version old uid: %d", uid)
		}
	}
	if wallet, err = d.Detail(c, uid); err == nil {
		d.SetWalletCache(c, &model.McDetail{Detail: wallet, Exist: true, Version: d.CacheVersion(c)}, 86400)
	}
	return
}

// Melonseed 获取瓜子数
func (d *Dao) GetMelonseedByCache(c context.Context, uid int64) (wallet *model.Melonseed, err error) {
	detail, err := d.GetDetailByCache(c, uid)
	wallet = new(model.Melonseed)
	if err == nil {
		wallet = &model.Melonseed{
			Uid:     detail.Uid,
			Gold:    detail.Gold,
			IapGold: detail.IapGold,
			Silver:  detail.Silver,
		}
	}
	return
}

func (d *Dao) ModifyCoin(c context.Context, coinNum int, uid int64, coinTypeNo int32, delCache bool) (bool, error) {

	var (
		affect int64
		err    error
		res    bool
	)
	switch coinTypeNo {
	case model.SysCoinTypeIosGold:
		affect, err = d.AddIapGold(c, uid, coinNum)
	case model.SysCoinTypeSilver:
		affect, err = d.AddSilver(c, uid, coinNum)
	case model.SysCoinTypeGold:
		affect, err = d.AddGold(c, uid, coinNum)
	default:
		// do nothing
	}
	if affect > 0 {
		res = true
		if delCache {
			d.DelWalletCache(c, uid)
		}
	}
	return res, err
}

func (d *Dao) GetCoin(c context.Context, coinTypeNo int32, uid int64) (interface{}, error) {
	userCoin, err := d.Melonseed(c, uid)
	switch coinTypeNo {
	case model.SysCoinTypeIosGold:
		return userCoin.IapGold, err
	case model.SysCoinTypeGold:
		return userCoin.Gold, err
	case model.SysCoinTypeSilver:
		return userCoin.Silver, err
	case model.SysCoinTypeMetal:
		metal, err := d.GetMetal(c, uid)
		return metal, err
	default:
		return nil, nil
	}
}

func (d *Dao) RechargeCoin(c context.Context, coinNum int, uid int64, coinTypeNo int32, delCache bool) (bool, error) {
	var (
		affect int64
		err    error
		res    bool
	)
	switch coinTypeNo {
	case model.SysCoinTypeIosGold:
		affect, err = d.RechargeIapGold(c, uid, coinNum)
	case model.SysCoinTypeSilver:
		affect, err = d.AddSilver(c, uid, coinNum)
	case model.SysCoinTypeGold:
		affect, err = d.RechargeGold(c, uid, coinNum)
	default:
		// do nothing
	}
	if affect > 0 {
		res = true
		if delCache {
			d.DelWalletCache(c, uid)
		}
	}
	return res, err
}

func (d *Dao) ConsumeCoin(c context.Context, coinNum int, uid int64, coinTypeNo int32, seeds int64, delCache bool, reason interface{}) (success bool, err error) {
	var affect int64

	switch coinTypeNo {
	case model.SysCoinTypeGold:
		affect, err = d.ConsumeGold(c, uid, coinNum)
	case model.SysCoinTypeIosGold:
		affect, err = d.ConsumeIapGold(c, uid, coinNum)
	case model.SysCoinTypeSilver:
		affect, err = d.ConsumeSilver(c, uid, coinNum)
	case model.SysCoinTypeMetal:
		success, _, err = d.ModifyMetal(c, uid, int64(-1*coinNum), seeds, reason)
	default:
	}

	if model.IsLocalCoin(coinTypeNo) {
		if affect > 0 {
			success = true
			if delCache {
				d.DelWalletCache(c, uid)
			}
		}
	}

	return
}
