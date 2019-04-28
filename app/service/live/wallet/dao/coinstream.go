package dao

import (
	"context"
	"crypto/md5"
	"fmt"
	"go-common/app/service/live/wallet/model"
	"go-common/library/log"
	"time"
)

const (
	_newCoinStremaRecord = "INSERT INTO %s (uid, transaction_id, extend_tid, coin_type, delta_coin_num, org_coin_num, op_result, op_reason, op_type, op_time, reserved2, area, source, reserved3, reserved4,platform,reserved1, reserved5) values(?, ?, ?,?, ?, ?,?, ?, ?,?, ?, ?,?, ?, ?, ?, ?, ?)"
	_getCoinStream       = "SELECT uid,transaction_id,extend_tid, coin_type,delta_coin_num,org_coin_num,op_result,op_reason,op_type, reserved3,area,reserved2,source,reserved4,platform,reserved1,reserved5 FROM %s WHERE transaction_id = ? order by id desc limit %d,1"
	_getCoinStreamByUid  = "SELECT uid,transaction_id,extend_tid, coin_type,delta_coin_num,org_coin_num,op_result,op_reason,op_type, reserved3,area,reserved2,source,reserved4, platform , reserved1,reserved5 FROM %s WHERE transaction_id = ? and uid = ? order by id desc limit %d,1"
)

func hexdec(s string) uint64 {
	d := uint64(0)
	for i := 0; i < len(s); i++ {
		x := uint64(s[i])
		if x >= 'a' {
			x -= 'a' - 'A'
		}
		d1 := x - '0'
		if d1 > 9 {
			d1 = 10 + d1 - ('A' - '0')
		}
		if 0 > d1 || d1 > 15 {
			return 0
		}
		d = (16 * d) + d1
	}
	return d
}

func getCoinStreamTable(transactionId string) string {
	tlen := len(transactionId)
	year := hexdec(transactionId[tlen-3 : tlen-1])
	if year == 0 {
		log.Error("illegal tid : %s", transactionId)
	}
	year = year + 2000

	t := fmt.Sprintf("t_coin_stream_%d%02d", year, hexdec(transactionId[tlen-1:tlen]))
	return t
}

func GetTid(serviceType model.ServiceType, v interface{}) string {
	s := fmt.Sprintf("%v%s", v, randomString(5))
	bizTid := fmt.Sprintf("%x", md5.Sum([]byte(s)))
	now := time.Now()
	year := now.Year() - 2000
	month := int(now.Month())

	return fmt.Sprintf("%s%x%02x%x", bizTid, serviceType, year, month)
}

func (d *Dao) NewCoinStreamRecord(c context.Context, record *model.CoinStreamRecord) (int64, error) {
	s := fmt.Sprintf(_newCoinStremaRecord, getCoinStreamTable(record.TransactionId))
	date := model.GetWalletFormatTime(record.OpTime)
	return execSqlWithBindParams(d, c, &s, record.Uid, record.TransactionId, record.ExtendTid, record.CoinType,
		record.DeltaCoinNum, record.OrgCoinNum, record.OpResult, record.OpReason, record.OpType, date,
		record.BizCode, record.Area, record.Source, record.MetaData, record.BizSource,
		record.Platform, record.Reserved1, record.Version,
	)
}

func (d *Dao) GetCoinStreamByTid(c context.Context, tid string) (record *model.CoinStreamRecord, err error) {
	return d.GetCoinStreamByTidAndOffset(c, tid, 0)
}

func (d *Dao) GetCoinStreamByTidAndOffset(c context.Context, tid string, offset int) (record *model.CoinStreamRecord, err error) {
	s := fmt.Sprintf(_getCoinStream, getCoinStreamTable(tid), offset)
	row := d.db.QueryRow(c, s, tid)
	record = &model.CoinStreamRecord{}
	err = row.Scan(&record.Uid, &record.TransactionId, &record.ExtendTid, &record.CoinType, &record.DeltaCoinNum,
		&record.OrgCoinNum, &record.OpResult, &record.OpReason, &record.OpType, &record.MetaData, &record.Area,
		&record.BizCode, &record.Source, &record.BizSource, &record.Platform, &record.Reserved1, &record.Version)
	return
}

func (d *Dao) GetCoinStreamByUidTid(c context.Context, uid int64, tid string) (record *model.CoinStreamRecord, err error) {
	s := fmt.Sprintf(_getCoinStreamByUid, getCoinStreamTable(tid), 0)
	row := d.db.QueryRow(c, s, tid, uid)
	record = &model.CoinStreamRecord{}
	err = row.Scan(&record.Uid, &record.TransactionId, &record.ExtendTid, &record.CoinType, &record.DeltaCoinNum,
		&record.OrgCoinNum, &record.OpResult, &record.OpReason, &record.OpType, &record.MetaData, &record.Area,
		&record.BizCode, &record.Source, &record.BizSource, &record.Platform, &record.Reserved1, &record.Version)
	return
}
