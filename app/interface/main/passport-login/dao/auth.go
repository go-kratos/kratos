package dao

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/interface/main/passport-login/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addCookieSQL         = "INSERT INTO user_cookie_%s (mid,session,csrf,type,expires) VALUES (?,?,?,?,?)"
	_addTokenSQL          = "INSERT INTO user_token_%s (mid,appid,token,expires,type) VALUES (?,?,?,?,?)"
	_addRefreshSQL        = "INSERT INTO user_refresh_%s (mid,appid,refresh,token,expires) VALUES (?,?,?,?,?)"
	_addOldCookieSQL      = "INSERT INTO aso_cookie_token%d (mid,session_data,csrf_token,type,expire_time) VALUES (?,?,?,?,?)"
	_addOldTokenSQL       = "INSERT INTO aso_app_perm_%s (mid,appid,access_token,refresh_token,app_subid,create_at,expires,type) VALUES (?,?,?,?,?,?,?,?)"
	_getCookieSQL         = "SELECT mid,session,csrf,type,expires FROM user_cookie_%s WHERE session = ? limit 1"
	_getTokenSQL          = "SELECT mid,appid,token,expires,type FROM user_token_%s WHERE token = ? limit 1"
	_getRefreshSQL        = "SELECT mid,appid,refresh,token,expires FROM user_refresh_%s WHERE refresh = ? limit 1"
	_delCookieSQL         = "DELETE FROM user_cookie_%s WHERE session = ?"
	_delTokenSQL          = "DELETE FROM user_token_%s WHERE token = ?"
	_delRefreshSQL        = "DELETE FROM user_refresh_%s WHERE refresh = ?"
	_delOldCookieSQL      = "DELETE FROM aso_cookie_token%d WHERE session_data = ?"
	_delOldTokenSQL       = "DELETE FROM aso_app_perm_%s WHERE access_token = ?"
	_delTokenByMidSQL     = "DELETE FROM user_token_%s WHERE mid = ?"
	_delCookieByMidSQL    = "DELETE FROM user_cookie_%s WHERE mid = ?"
	_delRefreshByMidSQL   = "DELETE FROM user_refresh_%s WHERE mid = ?"
	_delOldCookieByMidSQL = "DELETE FROM aso_cookie_token%d WHERE mid = ?"
	_delOldTokenByMidSQL  = "DELETE FROM aso_app_perm_%s WHERE mid = ?"
)

// AddCookie add cookie
func (d *Dao) AddCookie(c context.Context, cookie *model.Cookie, now time.Time) (affected int64, err error) {
	row, err := d.authDB.Exec(c, fmt.Sprintf(_addCookieSQL, formatAuthSuffix(now)), cookie.Mid, cookie.Session, cookie.CSRF, cookie.Type, cookie.Expires)
	if err != nil {
		log.Error("fail to add cookie(%+v), sessionHex(%s), dao.authDB.Exec() error(%+v)", cookie, hex.EncodeToString(cookie.Session), err)
		return
	}
	return row.RowsAffected()
}

// AddToken add token
func (d *Dao) AddToken(c context.Context, token *model.Token, t time.Time) (affected int64, err error) {
	var row sql.Result
	row, err = d.authDB.Exec(c, fmt.Sprintf(_addTokenSQL, formatAuthSuffix(t)), token.Mid, token.AppID, token.Token, token.Expires, token.Type)
	if err != nil {
		log.Error("fail to add token(%+v), tokenHex(%s), dao.authDB.Exec() error(%+v)", token, hex.EncodeToString(token.Token), err)
		return
	}
	return row.RowsAffected()
}

// AddRefresh add refresh
func (d *Dao) AddRefresh(c context.Context, refresh *model.Refresh, t time.Time) (affected int64, err error) {
	var row sql.Result
	row, err = d.authDB.Exec(c, fmt.Sprintf(_addRefreshSQL, formatRefreshSuffix(t)), refresh.Mid, refresh.AppID, refresh.Refresh, refresh.Token, refresh.Expires)
	if err != nil {
		log.Error("fail to add refresh(%+v), refreshHex(%s), dao.authDB.Exec() error(%+v)", refresh, hex.EncodeToString(refresh.Refresh), err)
		return
	}
	return row.RowsAffected()
}

// AddOldCookie add old cookie
func (d *Dao) AddOldCookie(c context.Context, cookie *model.OldCookie) (affected int64, err error) {
	row, err := d.originDB.Exec(c, fmt.Sprintf(_addOldCookieSQL, oldCookieSuffix(cookie.Mid)), cookie.Mid, cookie.Session, cookie.CSRFToken, cookie.Type, cookie.Expires)
	if err != nil {
		log.Error("fail to add oldCookie(%+v), dao.originDB.Exec() error(%+v)", cookie, err)
		return
	}
	return row.RowsAffected()
}

// AddOldToken add old token
func (d *Dao) AddOldToken(c context.Context, token *model.OldToken, t time.Time) (affected int64, err error) {
	row, err := d.originDB.Exec(c, fmt.Sprintf(_addOldTokenSQL, formatAuthSuffix(t)), token.Mid, token.AppID, token.AccessToken, token.RefreshToken, token.AppSubID, token.CreateAt, token.Expires, token.Type)
	if err != nil {
		log.Error("fail to add oldToken(%+v), dao.originDB.Exec() error(%+v)", token, err)
		return
	}
	return row.RowsAffected()
}

// GetCookie get cookie by session
func (d *Dao) GetCookie(c context.Context, session []byte, t time.Time) (res *model.Cookie, err error) {
	row := d.authDB.QueryRow(c, fmt.Sprintf(_getCookieSQL, formatAuthSuffix(t)), session)
	res = new(model.Cookie)
	if err = row.Scan(&res.Mid, &res.Session, &res.CSRF, &res.Type, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("fail to get cookie(%s), dao.authDB.QueryRow() error(%+v)", hex.EncodeToString(session), err)
		}
	}
	return
}

// GetToken get token by access_token
func (d *Dao) GetToken(c context.Context, token []byte, t time.Time) (res *model.Token, err error) {
	row := d.authDB.QueryRow(c, fmt.Sprintf(_getTokenSQL, formatAuthSuffix(t)), token)
	res = new(model.Token)
	if err = row.Scan(&res.Mid, &res.AppID, &res.Token, &res.Expires, &res.Type); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("fail to get token(%s), dao.authDB.QueryRow() error(%+v)", hex.EncodeToString(token), err)
		}
	}
	return
}

// GetRefresh get refresh by refresh_token
func (d *Dao) GetRefresh(c context.Context, refresh []byte, t time.Time) (res *model.Refresh, err error) {
	row := d.authDB.QueryRow(c, fmt.Sprintf(_getRefreshSQL, formatRefreshSuffix(t)), refresh)
	res = new(model.Refresh)
	if err = row.Scan(&res.Mid, &res.AppID, &res.Refresh, &res.Token, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("fail to get refresh(%s), dao.authDB.QueryRow() error(%+v)", hex.EncodeToString(refresh), err)
		}
	}
	return
}

// DelCookie del cookie by session
func (d *Dao) DelCookie(c context.Context, session []byte, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delCookieSQL, formatAuthSuffix(t)), session); err != nil {
		log.Error("fail to del cookie by session(%s), dao.authDB.Exec() error(%+v)", hex.EncodeToString(session), err)
		return
	}
	return res.RowsAffected()
}

// DelToken del token by access_token
func (d *Dao) DelToken(c context.Context, token []byte, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delTokenSQL, formatAuthSuffix(t)), token); err != nil {
		log.Error("fail to del token(%s), dao.authDB.Exec() error(%+v)", hex.EncodeToString(token), err)
		return
	}
	return res.RowsAffected()
}

// DelRefresh del refresh
func (d *Dao) DelRefresh(c context.Context, refresh []byte, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delRefreshSQL, formatRefreshSuffix(t)), refresh); err != nil {
		log.Error("fail to del refresh(%s), dao.authDB.Exec() error(%+v)", hex.EncodeToString(refresh), err)
		return
	}
	return res.RowsAffected()
}

// DelOldCookie del old cookie by session_data
func (d *Dao) DelOldCookie(c context.Context, session string, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.originDB.Exec(c, fmt.Sprintf(_delOldCookieSQL, oldCookieSuffix(mid)), session); err != nil {
		log.Error("fail to del old cookie(%s), dao.originDB.Exec() error(%+v)", session, err)
		return
	}
	return res.RowsAffected()
}

// DelOldToken del old token by access_token
func (d *Dao) DelOldToken(c context.Context, token string, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.originDB.Exec(c, fmt.Sprintf(_delOldTokenSQL, formatAuthSuffix(t)), token); err != nil {
		log.Error("fail to del old token(%s), dao.originDB.Exec() error(%+v)", token, err)
		return
	}
	return res.RowsAffected()
}

// DelCookieByMid del cookie by mid
func (d *Dao) DelCookieByMid(c context.Context, mid int64, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delCookieByMidSQL, formatAuthSuffix(t)), mid); err != nil {
		log.Error("fail to del cookie by mid(%d), dao.authDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// DelTokenByMid del token by mid
func (d *Dao) DelTokenByMid(c context.Context, mid int64, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delTokenByMidSQL, formatAuthSuffix(t)), mid); err != nil {
		log.Error("fail to del token by mid(%d), dao.authDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// DelRefreshByMid del refresh by mid
func (d *Dao) DelRefreshByMid(c context.Context, mid int64, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.authDB.Exec(c, fmt.Sprintf(_delRefreshByMidSQL, formatRefreshSuffix(t)), mid); err != nil {
		log.Error("fail to del refresh by mid(%d), dao.authDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// DelOldCookieByMid del old cookie by mid
func (d *Dao) DelOldCookieByMid(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.originDB.Exec(c, fmt.Sprintf(_delOldCookieByMidSQL, oldCookieSuffix(mid)), mid); err != nil {
		log.Error("fail to del old token by mid(%d), dao.originDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// DelOldTokenByMid del old token by mid
func (d *Dao) DelOldTokenByMid(c context.Context, mid int64, t time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.originDB.Exec(c, fmt.Sprintf(_delOldTokenByMidSQL, formatAuthSuffix(t)), mid); err != nil {
		log.Error("fail to del old token by mid(%d), dao.originDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

func formatAuthSuffix(t time.Time) string {
	return t.Format("200601")
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

func oldCookieSuffix(mid int64) int64 {
	return mid % 30
}
