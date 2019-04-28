package dao

import (
	"context"
	"fmt"
	"go-common/app/service/live/wallet/model"
	mc "go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_walletMcKey = "wu:%d" // 钱包数据的mc缓存
)

func mcKey(uid int64) string {
	return fmt.Sprintf(_walletMcKey, uid)
}

func (d *Dao) CacheVersion(c context.Context) int32 {
	return 1
}

func (d *Dao) IsNewVersion(c context.Context, detail *model.McDetail) bool {
	return detail.Version == d.CacheVersion(c)
}

// WalletCache 获取钱包缓存
func (d *Dao) WalletCache(c context.Context, uid int64) (detail *model.McDetail, err error) {
	key := mcKey(uid)
	conn := d.mc.Get(c)
	defer conn.Close()

	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			return
		}
		log.Error("[dao.mc_wallet|WalletCache] conn.Get(%s) error(%v)", key, err)
		err = ecode.ServerErr
		return
	}
	detail = &model.McDetail{}
	if err = conn.Scan(r, detail); err != nil {
		log.Error("[dao.mc_wallet|WalletCache] conn.Scan(%s) error(%v)", string(r.Value), err)
	}
	return
}

// SetWalletCache 设置钱包缓存
func (d *Dao) SetWalletCache(c context.Context, detail *model.McDetail, expire int32) (err error) {
	key := mcKey(detail.Detail.Uid)
	conn := d.mc.Get(c)
	defer conn.Close()

	if err = conn.Set(&mc.Item{
		Key:        key,
		Object:     detail,
		Flags:      mc.FlagProtobuf,
		Expiration: expire,
	}); err != nil {
		log.Error("[dao.mc_wallet|SetWalletCache] conn.Set(%s, %v) error(%v)", key, detail, err)
	}
	return
}

// DelWalletCache 删除等级缓存
func (d *Dao) DelWalletCache(c context.Context, uid int64) (err error) {
	key := mcKey(uid)
	conn := d.mc.Get(c)
	defer conn.Close()

	if err = conn.Delete(key); err == mc.ErrNotFound {
		return
	}
	if err != nil {
		log.Error("[dao.mc_wallet|DelWalletCache] conn.Delete(%s) error(%v)", key, err)
	}
	return
}
