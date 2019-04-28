package dao

import (
	"context"
	"database/sql"
	"fmt"

	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_shard         = 50
	_multi         = 100
	_selCoin       = "SELECT coins FROM coins_%02d where mid=?"
	_txUpdateCoins = "INSERT INTO coins_%02d(mid,coins) VALUES(?,?) ON DUPLICATE KEY UPDATE coins=coins+?"
)

// UpdateCoin update user coins.
func (d *Dao) UpdateCoin(c context.Context, mid int64, coin float64) (err error) {
	count := int64(coin * _multi)
	if _, err = d.coin.Exec(c, fmt.Sprintf(_txUpdateCoins, mid%_shard), mid, count, count); err != nil {
		PromError("db:UpdateCoin")
		log.Errorv(c, log.KV("log", "UpdateCoin"), log.KV("err", err))
		return
	}
	return
}

// TxUpdateCoins update coins
func (d *Dao) TxUpdateCoins(c context.Context, tx *xsql.Tx, mid int64, coin float64) (err error) {
	count := int64(coin * _multi)
	_, err = tx.Exec(fmt.Sprintf(_txUpdateCoins, mid%_shard), mid, count, count)
	if err != nil {
		PromError("db:TxUpdateCoins")
		log.Errorv(c, log.KV("log", "TxUpdateCoins"), log.KV("err", err), log.KV("mid", mid))
		return
	}
	return
}

// TxUserCoin tx user coin
func (d *Dao) TxUserCoin(c context.Context, tx *xsql.Tx, mid int64) (count float64, err error) {
	var coin int64
	row := tx.QueryRow(fmt.Sprintf(_selCoin, mid%_shard), mid)
	if err = row.Scan(&coin); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:TxUserCoin")
		log.Errorv(c, log.KV("log", "TxUserCoin"), log.KV("err", err))
	}
	count = float64(coin) / _multi
	return
}

// RawUserCoin get user coins.
func (d *Dao) RawUserCoin(c context.Context, mid int64) (res float64, err error) {
	var count int64
	row := d.coin.QueryRow(c, fmt.Sprintf(_selCoin, mid%_shard), mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:Coins")
		log.Errorv(c, log.KV("log", "Coins"), log.KV("err", err))
	}
	res = float64(count) / _multi
	return
}
