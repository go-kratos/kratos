package dao

import (
	"context"
	"fmt"
	"go-common/app/job/live/wallet/model"
	"go-common/library/log"
)

const (
	_shard       = 10
	_insWallet   = "INSERT IGNORE INTO user_wallet_%d(uid,gold,iap_gold,silver) VALUES(?,?,?,?)"
	_mergeWallet = "INSERT INTO user_wallet_%d(uid,gold,iap_gold,silver) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE gold=gold+%d,iap_gold=iap_gold+%d,silver=silver+%d"
)

func tableIndex(uid int64) int64 {
	return uid % _shard
}

func (d *Dao) InitWallet(c context.Context, user *model.User) (row int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insWallet, tableIndex(user.Uid)), user.Uid, user.Gold, user.IapGold, user.Silver)
	if err != nil {
		log.Error("[dao.wallet|InitWallet] d.db.Exec err: %v {uid:%d gold:%d iap_gold:%d silver:%d}", err, user.Uid, user.Gold, user.IapGold, user.Silver)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) MergeWallet(c context.Context, uid int64, gold int64, iap_gold int64, silver int64) (row int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_mergeWallet, tableIndex(uid), gold, iap_gold, silver), uid, gold, iap_gold, silver)
	if err != nil {
		log.Error("[dao.wallet|MergeWallet] d.db.Exec err: %v {uid:%d gold:%d iap_gold:%d silver:%d}", err, uid, gold, iap_gold, silver)
		return
	}
	return res.RowsAffected()
}
