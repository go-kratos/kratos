package dao

import (
	"context"
	"encoding/hex"
	"fmt"

	"go-common/app/job/main/identify/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_newSessionBinByteLen = 16

	_getCookieDeletedSQL = "SELECT id,mid,session FROM user_cookie_deleted_%s where id > ? limit ?"
	_getTokenDeletedSQL  = "SELECT id,mid,token FROM user_token_deleted_%s where id > ? limit ?"
)

// CookieDeleted get cookie deleted
func (d *Dao) CookieDeleted(c context.Context, start, count int64, suffix string) (res []*model.AuthCookie, err error) {
	var rows *xsql.Rows
	if rows, err = d.authDB.Query(c, fmt.Sprintf(_getCookieDeletedSQL, suffix), start, count); err != nil {
		log.Error("fail to get CookieDeleted, dao.authDB.Query(%s) error(%v)", _getCookieDeletedSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var session []byte
		a := new(model.AuthCookie)
		if err = rows.Scan(&a.ID, &a.Mid, &session); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		a.Session = encodeSession(session)
		res = append(res, a)
	}
	return
}

// TokenDeleted get token deleted
func (d *Dao) TokenDeleted(c context.Context, start, count int64, suffix string) (res []*model.AuthToken, err error) {
	var rows *xsql.Rows
	if rows, err = d.authDB.Query(c, fmt.Sprintf(_getTokenDeletedSQL, suffix), start, count); err != nil {
		log.Error("fail to get TokenDeleted, dao.authDB.Query(%s) error(%v)", _getTokenDeletedSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var token []byte
		a := new(model.AuthToken)
		if err = rows.Scan(&a.ID, &a.Mid, &token); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		a.Token = hex.EncodeToString(token)
		res = append(res, a)
	}
	return
}

func encodeSession(b []byte) (s string) {
	// format new
	if len(b) == _newSessionBinByteLen {
		return hex.EncodeToString(b)
	}
	// or format old
	return string(b)
}
