package dao

/*
import (
	"context"

	xsql "go-common/database/sql"
)

type TxImpl struct {
	*xsql.Tx
}

func NewTx(ctx context.Context) (Tx, error) {
	t, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &TxImpl{t}, nil
}

func (tx *TxImpl) UpdateKeyword(k *Keyword) error {
	return updateKeyword(tx.Ctx, tx, k)
}

func (tx *TxImpl) InsertKeyword(k *Keyword) error {
	return insertKeyword(tx.Ctx, tx, k)
}

func (tx *TxImpl) InsertRule(r *Rule) error {
	return insertRule(tx.Ctx, tx, r)
}

func (tx *TxImpl) UpdateRegexp(r *Regexp) error {
	return updateRegexp(tx.Ctx, tx, r)
}

func (tx *TxImpl) InsertRegexp(r *Regexp) error {
	return insertRegexp(tx.Ctx, tx, r)
}

func (tx *TxImpl) QueryRow(_ context.Context, sql string, args ...interface{}) *xsql.Row {
	return tx.Tx.QueryRow(sql, args)
}

func (tx *TxImpl) Query(_ context.Context, sql string, args ...interface{}) (*xsql.Rows, error) {
	return tx.Tx.Query(sql, args)
}*/
