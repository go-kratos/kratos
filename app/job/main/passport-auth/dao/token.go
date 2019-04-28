package dao

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/job/main/passport-auth/model"
	"go-common/library/log"
)

const (
	_addTokenSQL = "INSERT IGNORE INTO user_token_%s (mid,appid,token,expires,type) VALUES (?,?,?,?,?)"
	_delTokenSQL = "DELETE FROM user_token_%s where token = ?"

	_addTokenDeletedSQL = "INSERT IGNORE INTO user_token_deleted_%s (mid,appid,token,expires,type,ctime) VALUES (?,?,?,?,?,?)"
)

// AddToken save token
func (d *Dao) AddToken(c context.Context, t *model.Token, token []byte, ct time.Time) (affected int64, err error) {
	var row sql.Result
	if row, err = d.db.Exec(c, fmt.Sprintf(_addTokenSQL, formatSuffix(ct)), t.Mid, t.AppID, token, t.Expires, t.Type); err != nil {
		log.Error("d.AddToken(%v) err(%v)", t, err)
		return
	}
	return row.RowsAffected()
}

// DelToken del token
func (d *Dao) DelToken(c context.Context, token []byte, ct time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_delTokenSQL, formatSuffix(ct)), token); err != nil {
		log.Error("del token failed, dao.db.Exec(%s) error(%v)", hex.EncodeToString(token), err)
		return
	}
	return res.RowsAffected()
}

// AddTokenDeleted save token deleted
func (d *Dao) AddTokenDeleted(c context.Context, t *model.Token, token []byte, ct time.Time) (affected int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_addTokenDeletedSQL, formatSuffix(ct)), t.Mid, t.AppID, token, t.Expires, t.Type, t.Ctime)
	if err != nil {
		log.Error("fail to add token deleted, token(%+v), tx.Exec() error(%+v)", t, err)
		return
	}
	return row.RowsAffected()
}

func formatSuffix(t time.Time) string {
	return t.Format("200601")
}
