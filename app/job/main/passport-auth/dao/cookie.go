package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/passport-auth/model"
	"go-common/library/log"
)

const (
	_addCookieSQL          = "INSERT IGNORE INTO user_cookie_%s (mid,session,csrf,type,expires) VALUES (?,?,?,?,?)"
	_delCookieBySessionSQL = "DELETE FROM user_cookie_%s where session = ?"

	_addCookieDeletedSQL = "INSERT IGNORE INTO user_cookie_deleted_%s (mid,session,csrf,type,expires,ctime) VALUES (?,?,?,?,?,?)"
)

// AddCookie save cookie
func (d *Dao) AddCookie(c context.Context, cookie *model.Cookie, session, csrf []byte, ct time.Time) (affected int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_addCookieSQL, formatSuffix(ct)), cookie.Mid, session, csrf, cookie.Type, cookie.Expires)
	if err != nil {
		log.Error("dao.db.Exec(%v) err(%v)", cookie, err)
		return
	}
	return row.RowsAffected()
}

// DelCookie del cookie by session
func (d *Dao) DelCookie(c context.Context, session []byte, ct time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_delCookieBySessionSQL, formatSuffix(ct)), session); err != nil {
		log.Error("del cookie by session , dao.db.Exec(%s) error(%v)", session, err)
		return
	}
	return res.RowsAffected()
}

// AddCookieDeleted save cookie deleted
func (d *Dao) AddCookieDeleted(c context.Context, cookie *model.Cookie, session, csrf []byte, ct time.Time) (affected int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_addCookieDeletedSQL, formatSuffix(ct)), cookie.Mid, session, csrf, cookie.Type, cookie.Expires, cookie.Ctime)
	if err != nil {
		log.Error("fail to add cookie deleted, cookie(%+v), tx.Exec() error(%+v)", cookie, err)
		return
	}
	return row.RowsAffected()
}
