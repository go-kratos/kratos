package dao

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/service/main/passport-auth/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getRefreshSQL = "SELECT mid,appid,refresh,token,expires FROM user_refresh_%s WHERE refresh = ? limit 1"
)

// Refresh get token by access_token
func (d *Dao) Refresh(c context.Context, rk []byte, ct time.Time) (res *model.Refresh, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getRefreshSQL, formatRefreshSuffix(ct)), rk)
	res = new(model.Refresh)
	var refresh, token []byte
	if err = row.Scan(&res.Mid, &res.AppID, &refresh, &token, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	res.Refresh = hex.EncodeToString(refresh)
	res.Token = hex.EncodeToString(token)
	return
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
