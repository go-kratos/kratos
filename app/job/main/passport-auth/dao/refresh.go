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
	_addRefreshSQL = "INSERT IGNORE INTO user_refresh_%s (mid,appid,refresh,token,expires) VALUES (?,?,?,?,?)"
	_delRefreshSQL = "DELETE FROM user_refresh_%s WHERE refresh = ?"
)

// AddRefresh save token
func (d *Dao) AddRefresh(c context.Context, t *model.Refresh, refresh, token []byte, ct time.Time) (affected int64, err error) {
	var row sql.Result
	if row, err = d.db.Exec(c, fmt.Sprintf(_addRefreshSQL, formatRefreshSuffix(ct)), t.Mid, t.AppID, refresh, token, t.Expires); err != nil {
		log.Error("d.AddToken(%v) err(%v)", t, err)
		return
	}
	return row.RowsAffected()
}

// DelRefresh del token
func (d *Dao) DelRefresh(c context.Context, refresh []byte, ct time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_delRefreshSQL, formatRefreshSuffix(ct)), refresh); err != nil {
		log.Error("del token failed, dao.db.Exec(%s) error(%v)", hex.EncodeToString(refresh), err)
		return
	}
	return res.RowsAffected()
}

func formatRefreshSuffix(t time.Time) string {
	return formatByDate(t.Year(), int(t.Month()))
}

func formatByDate(year, month int) string {
	if month%2 == 1 {
		return fmt.Sprintf("%4d%02d", year, month)
	}
	return fmt.Sprintf("%4d%02d", year, month-1)
}
