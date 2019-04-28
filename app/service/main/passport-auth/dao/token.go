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
	_getTokenSQL = "SELECT mid,appid,token,expires,type FROM user_token_%s where token = ? limit 1"
)

// Token get token by access_token
func (d *Dao) Token(c context.Context, token []byte, ct time.Time) (res *model.Token, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getTokenSQL, formatSuffix(ct)), token)
	res = new(model.Token)
	var tokenByte []byte
	if err = row.Scan(&res.Mid, &res.AppID, &tokenByte, &res.Expires, &res.Type); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	res.Token = hex.EncodeToString(tokenByte)
	return
}

func formatSuffix(t time.Time) string {
	return t.Format("200601")
}
