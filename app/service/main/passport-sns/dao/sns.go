package dao

import (
	"context"
	"database/sql"
	"fmt"
	"hash/crc32"

	"go-common/app/service/main/passport-sns/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getSnsAppsSQL          = "SELECT appid,appsecret,platform,business FROM sns_apps"
	_getSnsUsersSQL         = "SELECT mid,unionid,platform,expires FROM sns_user WHERE mid = ?"
	_getSnsTokensSQL        = "SELECT mid,openid,unionid,platform,token,expires FROM sns_token WHERE mid = ?"
	_getSnsUserByMidSQL     = "SELECT mid,unionid,platform FROM sns_user WHERE mid = ? and platform = ?"
	_getSnsUserByUnionIDSQL = "SELECT mid,unionid,platform FROM sns_user WHERE unionid = ? and platform = ?"
	_addSnsUserSQL          = "INSERT INTO sns_user (mid,unionid,platform,expires) VALUES(?,?,?,?)"
	_addSnsOpenIDSQL        = "INSERT IGNORE INTO sns_openid_%02d (openid,unionid,appid,platform) VALUES(?,?,?,?)"
	_addSnsTokenSQL         = "INSERT INTO sns_token (mid,openid,unionid,platform,token,expires,appid) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?, openid =?, token =?, expires =?, appid =?"
	_updateSnsTokenSQL      = "UPDATE sns_token SET token =?, expires =? WHERE mid =? and platform = ?"
	_updateSnsUserSQL       = "UPDATE sns_user SET expires =? WHERE mid =? and platform =?"
	_delSnsUserSQL          = "DELETE FROM sns_user WHERE mid = ? and platform = ?"
	_delSnsUsersSQL         = "DELETE FROM sns_user WHERE mid = ?"
)

// SnsApps get sns apps
func (d *Dao) SnsApps(c context.Context) (res []*model.SnsApps, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getSnsAppsSQL); err != nil {
		log.Error("SnsApps dao.db.Query error(%+v)", err)
		return
	}
	res = make([]*model.SnsApps, 0)
	defer rows.Close()
	for rows.Next() {
		r := new(model.SnsApps)
		if err = rows.Scan(&r.AppID, &r.AppSecret, &r.Platform, &r.Business); err != nil {
			log.Error("SnsApps row.Scan() error(%+v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// SnsUsers get sns users
func (d *Dao) SnsUsers(c context.Context, mid int64) (res []*model.SnsUser, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getSnsUsersSQL, mid); err != nil {
		log.Error("SnsUsers dao.db.Query error(%+v)", err)
		return
	}
	res = make([]*model.SnsUser, 0)
	defer rows.Close()
	for rows.Next() {
		r := new(model.SnsUser)
		if err = rows.Scan(&r.Mid, &r.UnionID, &r.Platform, &r.Expires); err != nil {
			log.Error("SnsUsers row.Scan() error(%+v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// SnsTokens get sns tokens
func (d *Dao) SnsTokens(c context.Context, mid int64) (res []*model.SnsToken, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getSnsTokensSQL, mid); err != nil {
		log.Error("SnsTokens dao.db.Query error(%+v)", err)
		return
	}
	res = make([]*model.SnsToken, 0)
	defer rows.Close()
	for rows.Next() {
		r := new(model.SnsToken)
		if err = rows.Scan(&r.Mid, &r.OpenID, &r.UnionID, &r.Platform, &r.Token, &r.Expires); err != nil {
			log.Error("SnsTokens row.Scan() error(%+v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// SnsUserByMid get sns user by mid and platform
func (d *Dao) SnsUserByMid(c context.Context, mid int64, platform int) (res *model.SnsUser, err error) {
	res = new(model.SnsUser)
	row := d.db.QueryRow(c, _getSnsUserByMidSQL, mid, platform)
	if err = row.Scan(&res.Mid, &res.UnionID, &res.Platform); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("SnsUserByMid mid(%d) platform(%d) row.Scan() error(%+v)", mid, platform, err)
		return
	}
	return
}

// SnsUserByUnionID get sns user by unionID and platform
func (d *Dao) SnsUserByUnionID(c context.Context, unionID string, platform int) (res *model.SnsUser, err error) {
	res = new(model.SnsUser)
	row := d.db.QueryRow(c, _getSnsUserByUnionIDSQL, unionID, platform)
	if err = row.Scan(&res.Mid, &res.UnionID, &res.Platform); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("SnsUserByUnionID unionID(%s) platform(%d) row.Scan() error(%+v)", unionID, platform, err)
		return
	}
	return
}

// TxAddSnsUser add sns user.
func (d *Dao) TxAddSnsUser(tx *xsql.Tx, a *model.SnsUser) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_addSnsUserSQL, a.Mid, a.UnionID, a.Platform, a.Expires); err != nil {
		log.Error("TxAddSnsUser(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddSnsOpenID add sns openid.
func (d *Dao) TxAddSnsOpenID(tx *xsql.Tx, a *model.SnsOpenID) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addSnsOpenIDSQL, openIDSuffix(a.OpenID)), a.OpenID, a.UnionID, a.AppID, a.Platform); err != nil {
		log.Error("TxAddSnsOpenID(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddSnsToken add sns token.
func (d *Dao) TxAddSnsToken(tx *xsql.Tx, a *model.SnsToken) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_addSnsTokenSQL, a.Mid, a.OpenID, a.UnionID, a.Platform, a.Token, a.Expires, a.AppID, a.Mid, a.OpenID, a.Token, a.Expires, a.AppID); err != nil {
		log.Error("TxAddSnsToken(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateSnsUser update sns user expires.
func (d *Dao) TxUpdateSnsUser(tx *xsql.Tx, a *model.SnsUser) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateSnsUserSQL, a.Expires, a.Mid, a.Platform); err != nil {
		log.Error("TxUpdateSnsUser(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateSnsToken update sns token.
func (d *Dao) TxUpdateSnsToken(tx *xsql.Tx, a *model.SnsToken) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateSnsTokenSQL, a.Token, a.Expires, a.Mid, a.Platform); err != nil {
		log.Error("TxUpdateSnsToken(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// DelSnsUser del sns user.
func (d *Dao) DelSnsUser(c context.Context, mid int64, platform int) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _delSnsUserSQL, mid, platform); err != nil {
		log.Error("DelSnsUser mid(%d) platform(%d) d.db.Exec() error(%+v)", mid, platform, err)
		return
	}
	return res.RowsAffected()
}

// DelSnsUsers del sns user by mid.
func (d *Dao) DelSnsUsers(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _delSnsUsersSQL, mid); err != nil {
		log.Error("DelAllSnsUser mid(%d) d.db.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

func openIDSuffix(openID string) int {
	v := int(crc32.ChecksumIEEE([]byte(openID)))
	if v < 0 {
		v = -v
	}
	return v % 100
}
