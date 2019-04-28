package dao

import (
	"context"
	"fmt"
	"go-common/app/service/live/wallet/model"
)

const (
	_newCoinExchange = "insert into %s (uid, transaction_id, src_type,src_num,dest_type,dest_num,status,exchange_time) values(?,?,?,?,?,?,?,?)"
)

func getCoinExchangeTableIndex(uid int64) string {
	return fmt.Sprintf("%02d", uid%10)
}

func getCoinExchangeTable(uid int64) string {
	return fmt.Sprintf("t_coin_exchange_%s", getCoinExchangeTableIndex(uid))
}

func (d *Dao) NewCoinExchangeRecord(c context.Context, record *model.CoinExchangeRecord) (int64, error) {
	s := fmt.Sprintf(_newCoinExchange, getCoinExchangeTable(record.Uid))
	date := model.GetWalletFormatTime(record.ExchangeTime)
	return execSqlWithBindParams(d, c, &s, record.Uid, record.TransactionId, record.SrcType, record.SrcNum, record.DestType, record.DestNum, record.Status, date)
}
