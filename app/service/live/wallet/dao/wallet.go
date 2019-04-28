package dao

import (
	"context"
	"database/sql"
	"fmt"
	"go-common/app/service/live/wallet/model"
	"go-common/library/log"
)

const (
	_shard           = 10
	_insWallet       = "INSERT IGNORE INTO user_wallet_%d (uid,gold,iap_gold,silver) VALUES(?,?,?,?)"
	_selMelonseed    = "SELECT uid,gold,iap_gold,silver FROM user_wallet_%d WHERE uid=?"
	_selDetail       = "SELECT uid,gold,iap_gold,silver,gold_recharge_cnt,gold_pay_cnt,silver_pay_cnt, cost_base FROM user_wallet_%d WHERE uid=?"
	_addGold         = "INSERT INTO user_wallet_%d(uid,gold) VALUES(?,?) ON DUPLICATE KEY UPDATE gold=gold+%d"
	_rechargeGold    = "INSERT INTO user_wallet_%d(uid,gold,gold_recharge_cnt) VALUES(?,?,?) ON DUPLICATE KEY UPDATE gold=gold+%d,gold_recharge_cnt=gold_recharge_cnt+%d"
	_addIapGold      = "INSERT INTO user_wallet_%d(uid,iap_gold) VALUES(?,?) ON DUPLICATE KEY UPDATE iap_gold=iap_gold+%d"
	_rechargeIapGold = "INSERT INTO user_wallet_%d(uid,iap_gold,gold_recharge_cnt) VALUES(?,?,?) ON DUPLICATE KEY UPDATE iap_gold=iap_gold+%d,gold_recharge_cnt=gold_recharge_cnt+%d"
	_addSilver       = "INSERT INTO user_wallet_%d(uid,silver) VALUES(?,?) ON DUPLICATE KEY UPDATE silver=silver+%d"
	_consumeGold     = "UPDATE user_wallet_%d SET gold=gold-%d,gold_pay_cnt=gold_pay_cnt+%d WHERE uid=?"
	_consumeIapGold  = "UPDATE user_wallet_%d SET iap_gold=iap_gold-%d,gold_pay_cnt=gold_pay_cnt+%d WHERE uid=?"
	_consumeSilver   = "UPDATE user_wallet_%d SET silver=silver-%d,silver_pay_cnt=silver_pay_cnt+%d WHERE uid=?"
	_changeCosetBase = "UPDATE user_wallet_%d SET cost_base=? WHERE uid=?"
)

func tableIndex(uid int64) int64 {
	return uid % _shard
}

// InitExp 初始化用户钱包,用于首次查询
func (d *Dao) InitWallet(c context.Context, uid int64, gold int64, iap_gold int64, silver int64) (row int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insWallet, tableIndex(uid)), uid, gold, iap_gold, silver)
	if err != nil {
		log.Error("[dao.wallet|InitWallet] d.db.Exec err: %v", err)
		return
	}
	return res.RowsAffected()
}

// Melonseed 获取瓜子数
func (d *Dao) Melonseed(c context.Context, uid int64) (wallet *model.Melonseed, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selMelonseed, tableIndex(uid)), uid)
	wallet = &model.Melonseed{}
	if err = row.Scan(&wallet.Uid, &wallet.Gold, &wallet.IapGold, &wallet.Silver); err == sql.ErrNoRows {
		// 查询结果为空时，初始化数据
		_, err = d.InitWallet(c, uid, 0, 0, 0)
	}
	if err != nil {
		log.Error("[dao.wallet|Melonseed] row.Scan err: %v", err)
		return
	}
	return
}

// Detail 详细数据
func (d *Dao) Detail(c context.Context, uid int64) (detail *model.Detail, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selDetail, tableIndex(uid)), uid)
	detail = &model.Detail{}
	err = row.Scan(&detail.Uid, &detail.Gold, &detail.IapGold, &detail.Silver, &detail.GoldRechargeCnt,
		&detail.GoldPayCnt, &detail.SilverPayCnt, &detail.CostBase)
	if err == sql.ErrNoRows {
		detail.Gold = 0
		detail.IapGold = 0
		detail.Silver = 0
		detail.GoldPayCnt = 0
		detail.GoldRechargeCnt = 0
		detail.SilverPayCnt = 0
		detail.CostBase = 0
		detail.Uid = uid
		err = nil
		return
	}
	if err != nil {
		log.Error("[dao.wallet|Detail] row.Scan err: %v", err)
		return
	}
	return
}

func (d *Dao) DetailWithoutDefault(c context.Context, uid int64) (detail *model.Detail, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selDetail, tableIndex(uid)), uid)
	detail = &model.Detail{}
	err = row.Scan(&detail.Uid, &detail.Gold, &detail.IapGold, &detail.Silver, &detail.GoldRechargeCnt,
		&detail.GoldPayCnt, &detail.SilverPayCnt, &detail.CostBase)
	return
}

// AddGold 添加金瓜子
func (d *Dao) AddGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_addGold, tableIndex(uid), num)
	return execSqlWithBindParams(d, c, &s, uid, num)
}

// RechargeGold 充值金瓜子 记入充值总值
func (d *Dao) RechargeGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_rechargeGold, tableIndex(uid), num, num)
	return execSqlWithBindParams(d, c, &s, uid, num, num)
}

// RechargeGold 添加IOS金瓜子
func (d *Dao) AddIapGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_addIapGold, tableIndex(uid), num)
	return execSqlWithBindParams(d, c, &s, uid, num)
}

// RechargeGold 充值IOS金瓜子 记入充值总值
func (d *Dao) RechargeIapGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_rechargeIapGold, tableIndex(uid), num, num)
	return execSqlWithBindParams(d, c, &s, uid, num, num)
}

// AddSilver 添加银瓜子
func (d *Dao) AddSilver(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_addSilver, tableIndex(uid), num)
	return execSqlWithBindParams(d, c, &s, uid, num)
}

// AddSilver 消费金瓜子 记入消费总值
func (d *Dao) ConsumeGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_consumeGold, tableIndex(uid), num, num)
	return execSqlWithBindParams(d, c, &s, uid)
}

// AddSilver 消费IOS金瓜子 记入消费总值
func (d *Dao) ConsumeIapGold(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_consumeIapGold, tableIndex(uid), num, num)
	return execSqlWithBindParams(d, c, &s, uid)
}

// ConsumeSilver 消费银瓜子 记入消费总值
func (d *Dao) ConsumeSilver(c context.Context, uid int64, num int) (affect int64, err error) {
	s := fmt.Sprintf(_consumeSilver, tableIndex(uid), num, num)
	return execSqlWithBindParams(d, c, &s, uid)
}

func (d *Dao) ChangeCostBase(c context.Context, uid int64, num int64) (affect int64, err error) {
	s := fmt.Sprintf(_changeCosetBase, tableIndex(uid))
	return execSqlWithBindParams(d, c, &s, num, uid)
}

// 更新镜像时间　方便测试
func (d *Dao) UpdateSnapShotTime(context context.Context, uid int64, time string) (int64, error) {
	s := fmt.Sprintf("update user_wallet_%d set snapshot_time = ?", tableIndex(uid))
	return execSqlWithBindParams(d, context, &s, time)
}
