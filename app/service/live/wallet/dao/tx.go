package dao

import (
	"context"
	"fmt"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"
)

const (
	_selWallet                  = "SELECT uid,gold,iap_gold,silver,cost_base,gold_recharge_cnt,gold_pay_cnt,silver_pay_cnt,snapshot_time,snapshot_gold,snapshot_iap_gold,snapshot_silver,reserved1,reserved2 FROM user_wallet_%d WHERE uid=? FOR UPDATE"
	_recharge                   = "UPDATE user_wallet_%d set %s = %s + %d,%s=%s+%d where uid = ?"
	_rechargeWihoutCnt          = "UPDATE user_wallet_%d set %s=%s + %d where uid = ?"
	_rechargeWithSnap           = "UPDATE user_wallet_%d set %s=%s + %d,%s = %s+ %d, snapshot_time = ? , snapshot_gold = ? , snapshot_iap_gold = ? , snapshot_silver = ? where uid = ?"
	_rechargeWithSnapWithoutCnt = "UPDATE user_wallet_%d set %s = %s + %d, snapshot_time = ? , snapshot_gold = ? , snapshot_iap_gold = ? , snapshot_silver = ? where uid = ?"
	_exchange                   = "UPDATE user_wallet_%d set %s = %s - %d , %s = %s + %d, %s = %s + %d, %s = %s + %d where uid = ?"
	_exhcangeWithSnap           = "UPDATE user_wallet_%d set %s = %s - %d , %s = %s + %d, %s = %s + %d, %s = %s + %d, snapshot_time = ? , snapshot_gold = ? , snapshot_iap_gold = ? , snapshot_silver = ? where uid = ? "
	_modifyCnt                  = "UPDATE user_wallet_%d set gold_pay_cnt = gold_pay_cnt + %d, gold_recharge_cnt = gold_recharge_cnt + %d, silver_pay_cnt = silver_pay_cnt + %d where uid = ?"
)

// 开启事务
func (d *Dao) BeginTx(c context.Context) (conn *sql.Tx, err error) {
	return d.db.Begin(c)
}

func (d *Dao) DoTx(c context.Context, doFunc func(conn *sql.Tx) (v interface{}, err error)) (v interface{}, err error) {
	conn, err := d.BeginTx(c)
	if err != nil {
		err = ecode.ServerErr
		return
	}
	v, err = doFunc(conn)
	var txErr error
	if err != nil {
		conn.Rollback()
		err = ecode.ServerErr
	} else {
		txErr = conn.Commit()
		if txErr != nil {
			err = ecode.ServerErr
			v = nil
		}
	}
	return v, err
}

// 为了后续的更新获取数据
func (d *Dao) WalletForUpdate(conn *sql.Tx, uid int64) (wallet *model.DetailWithSnapShot, err error) {
	row := conn.QueryRow(fmt.Sprintf(_selWallet, tableIndex(uid)), uid)
	wallet = &model.DetailWithSnapShot{}
	var snapShotTime time.Time
	if err = row.Scan(&wallet.Uid, &wallet.Gold, &wallet.IapGold, &wallet.Silver, &wallet.CostBase, &wallet.GoldRechargeCnt,
		&wallet.GoldPayCnt, &wallet.SilverPayCnt, &snapShotTime, &wallet.SnapShotGold, &wallet.SnapShotIapGold,
		&wallet.SnapShotSilver, &wallet.Reserved1, &wallet.Reserved2); err == sql.ErrNoRows {
		// 查询结果为空时，初始化数据
		_, err = d.InitWalletInTx(conn, uid, 0, 0, 0)
		wallet.SnapShotTime = snapShotTime.Format("2006-01-02 15:04:05")
		return
	}
	if err != nil {
		log.Error("[tx.wallet|Melonseed] row.Scan err: %s", err.Error())
		return
	}
	wallet.SnapShotTime = snapShotTime.Format("2006-01-02 15:04:05")
	return
}

// InitExp 初始化用户钱包,用于首次查询
func (d *Dao) InitWalletInTx(conn *sql.Tx, uid int64, gold int64, iap_gold int64, silver int64) (row int64, err error) {
	res, err := conn.Exec(fmt.Sprintf(_insWallet, tableIndex(uid)), uid, gold, iap_gold, silver)
	if err != nil {
		log.Error("[tx.wallet|InitWallet] Exec err: %v", err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) execSqlInTx(conn *sql.Tx, sql *string, params ...interface{}) (affect int64, err error) {
	res, err := conn.Exec(*sql, params...)
	if err != nil {
		log.Error("[tx.wallet|execSqlInTx] Exec err: %v sql:%s", err, *sql)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) changeCoinInTx(conn *sql.Tx, uid int64, sysCoinTypeNo int32, num int64, originWallet *model.DetailWithSnapShot, cntField string) (affect int64, err error) {
	// 判断

	coinType := model.GetSysCoinTypeByNo(sysCoinTypeNo)
	var s string
	absNum := num
	if absNum < 0 {
		absNum = -absNum
	}
	if model.TodayNeedSnapShot(originWallet) {
		if cntField == "" {
			s = getRechargeWithoutCntWithSnapShotSQL(uid, coinType, num)
		} else {
			s = getRechargeWithSnapShotSQL(uid, coinType, num, cntField, absNum)
		}
		date := time.Now().Format("2006-01-02 15:04:05")
		return d.execSqlInTx(conn, &s, date, originWallet.Gold, originWallet.IapGold, originWallet.Silver, uid)
	} else {
		if cntField == "" {
			s = getRechargeWithoutCntSQL(uid, coinType, num)
		} else {
			s = getRechargeSQL(uid, coinType, num, cntField, absNum)
		}
		return d.execSqlInTx(conn, &s, uid)
	}
}

// RechargeGold 充值IOS金瓜子 记入充值总值
func (d *Dao) RechargeCoinInTx(conn *sql.Tx, uid int64, sysCoinTypeNo int32, num int64, originWallet *model.DetailWithSnapShot) (affect int64, err error) {
	rechargerCntField := model.GetRechargeCnt(sysCoinTypeNo)
	return d.changeCoinInTx(conn, uid, sysCoinTypeNo, num, originWallet, rechargerCntField)
}

func (d *Dao) PayCoinInTx(conn *sql.Tx, uid int64, sysCoinTypeNo int32, num int64, originWallet *model.DetailWithSnapShot) (affect int64, err error) {
	cntField := model.GetPayCnt(sysCoinTypeNo)
	return d.changeCoinInTx(conn, uid, sysCoinTypeNo, -num, originWallet, cntField)
}

func (d *Dao) ModifyCoinInTx(conn *sql.Tx, uid int64, sysCoinTypeNo int32, num int64, originWallet *model.DetailWithSnapShot) (affect int64, err error) {
	return d.changeCoinInTx(conn, uid, sysCoinTypeNo, num, originWallet, "")
}

func (d *Dao) ExchangeCoinInTx(conn *sql.Tx, uid int64, srcSysCoinTypeNo int32, srcNum int64, destSysCoinTypeNo int32, destNum int64, originWallet *model.DetailWithSnapShot) (affect int64, err error) {
	rechargeCntNum := destNum
	var rechargerCntField string
	if destSysCoinTypeNo == model.SysCoinTypeSilver { // 如果为银瓜子则不存在充值统计总数字段 使用 消费统计字段代替　适配sql
		rechargerCntField = "silver_pay_cnt"
		rechargeCntNum = 0
	} else {
		rechargerCntField = model.GetRechargeCnt(destSysCoinTypeNo)
	}
	payCntField := model.GetPayCnt(srcSysCoinTypeNo)
	srcCoinType := model.GetSysCoinTypeByNo(srcSysCoinTypeNo)
	destCoinType := model.GetSysCoinTypeByNo(destSysCoinTypeNo)

	var s string
	if model.TodayNeedSnapShot(originWallet) {
		s = fmt.Sprintf(_exhcangeWithSnap, tableIndex(uid), srcCoinType, srcCoinType, srcNum, destCoinType, destCoinType, destNum, rechargerCntField, rechargerCntField, rechargeCntNum, payCntField, payCntField, srcNum)
		date := time.Now().Format("2006-01-02 15:04:05")
		return d.execSqlInTx(conn, &s, date, originWallet.Gold, originWallet.IapGold, originWallet.Silver, uid)
	} else {
		s = fmt.Sprintf(_exchange, tableIndex(uid), srcCoinType, srcCoinType, srcNum, destCoinType, destCoinType, destNum, rechargerCntField, rechargerCntField, rechargeCntNum, payCntField, payCntField, srcNum)
		return d.execSqlInTx(conn, &s, uid)
	}
}

func getRechargeWithSnapShotSQL(uid int64, coinType string, num int64, cntField string, cntNum int64) string {
	return fmt.Sprintf(_rechargeWithSnap, tableIndex(uid), coinType, coinType, num, cntField, cntField, cntNum)
}

func getRechargeWithoutCntWithSnapShotSQL(uid int64, coinType string, num int64) string {
	return fmt.Sprintf(_rechargeWithSnapWithoutCnt, tableIndex(uid), coinType, coinType, num)
}

func getRechargeSQL(uid int64, coinType string, num int64, cntField string, cntNum int64) string {
	return fmt.Sprintf(_recharge, tableIndex(uid), coinType, coinType, num, cntField, cntField, cntNum)
}

func getRechargeWithoutCntSQL(uid int64, coinType string, num int64) string {
	return fmt.Sprintf(_rechargeWihoutCnt, tableIndex(uid), coinType, coinType, num)
}

func (d *Dao) NewCoinStreamRecordInTx(conn *sql.Tx, record *model.CoinStreamRecord) (int64, error) {
	s := fmt.Sprintf(_newCoinStremaRecord, getCoinStreamTable(record.TransactionId))
	date := model.GetWalletFormatTime(record.OpTime)
	return d.execSqlInTx(conn, &s, record.Uid, record.TransactionId, record.ExtendTid, record.CoinType,
		record.DeltaCoinNum, record.OrgCoinNum, record.OpResult, record.OpReason, record.OpType, date,
		record.BizCode, record.Area, record.Source, record.MetaData, record.BizSource,
		record.Platform, record.Reserved1, record.Version)
}

func (d *Dao) NewCoinExchangeRecordInTx(conn *sql.Tx, record *model.CoinExchangeRecord) (int64, error) {
	s := fmt.Sprintf(_newCoinExchange, getCoinExchangeTable(record.Uid))
	date := model.GetWalletFormatTime(record.ExchangeTime)
	return d.execSqlInTx(conn, &s, record.Uid, record.TransactionId, record.SrcType, record.SrcNum, record.DestType, record.DestNum, record.Status, date)
}

func (d *Dao) ModifyCntInTx(conn *sql.Tx, uid int64, goldPay int, goldRecharge int, silverPay int) (int64, error) {
	s := fmt.Sprintf(_modifyCnt, tableIndex(uid), goldPay, goldRecharge, silverPay)
	return d.execSqlInTx(conn, &s, uid)
}
