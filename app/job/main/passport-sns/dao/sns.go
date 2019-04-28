package dao

import (
	"context"
	"database/sql"
	"fmt"
	"hash/crc32"

	"go-common/app/job/main/passport-sns/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getSnsUserByMidSQL      = "SELECT mid,unionid,platform,expires FROM sns_user WHERE mid = ? and platform = ?"
	_getSnsUserByUnionIDSQL  = "SELECT mid,unionid,platform FROM sns_user WHERE unionid = ? and platform = ?"
	_addSnsUserSQL           = "INSERT INTO sns_user (mid,unionid,platform,expires) VALUES(?,?,?,?)"
	_addSnsOpenIDSQL         = "INSERT IGNORE INTO sns_openid_%02d (openid,unionid,appid,platform) VALUES(?,?,?,?)"
	_addSnsTokenSQL          = "INSERT INTO sns_token (mid,openid,unionid,platform,token,expires,appid) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?, openid =?, token =?, expires =?, appid =?"
	_delSnsUserSQL           = "DELETE FROM sns_user WHERE mid = ? and platform = ?"
	_updateSnsUserExpiresSQL = "UPDATE sns_user SET expires = ? where mid = ? and platform = ?"
	_updateSnsUserSQL        = "UPDATE sns_user SET unionid = ?, expires = ? where mid = ? and platform = ?"
	_updateSnsTokenSQL       = "UPDATE sns_token SET token =?, expires =? WHERE mid =? and platform = ?"
)

// SnsUserByMid get sns user by mid and platform
func (d *Dao) SnsUserByMid(c context.Context, mid int64, platform int) (res *model.SnsUser, err error) {
	res = new(model.SnsUser)
	row := d.snsDB.QueryRow(c, _getSnsUserByMidSQL, mid, platform)
	if err = row.Scan(&res.Mid, &res.UnionID, &res.Platform, &res.Expires); err != nil {
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
	row := d.snsDB.QueryRow(c, _getSnsUserByUnionIDSQL, unionID, platform)
	if err = row.Scan(&res.Mid, &res.UnionID, &res.Platform); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("SnsUserByUnionID unionID(%d) platform(%d) row.Scan() error(%+v)", unionID, platform, err)
		return
	}
	return
}

// AddSnsUser add sns user.
func (d *Dao) AddSnsUser(c context.Context, mid, expires int64, unionID string, platform int) (affected int64, err error) {
	var res sql.Result
	if res, err = d.snsDB.Exec(c, _addSnsUserSQL, mid, unionID, platform, expires); err != nil {
		log.Error("AddSnsUser mid(%d) platform(%d) unionID(%s) expires(%d) d.snsDB.Exec() error(%+v)", mid, platform, unionID, platform, err)
		return
	}
	return res.RowsAffected()
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

// DelSnsUser del sns user.
func (d *Dao) DelSnsUser(c context.Context, mid int64, platform int) (affected int64, err error) {
	var res sql.Result
	if res, err = d.snsDB.Exec(c, _delSnsUserSQL, mid, platform); err != nil {
		log.Error("DelSnsUser mid(%d) platform(%d) d.snsDB.Exec() error(%+v)", mid, platform, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateSnsUserExpires update sns user expires.
func (d *Dao) TxUpdateSnsUserExpires(tx *xsql.Tx, a *model.SnsUser) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateSnsUserExpiresSQL, a.Expires, a.Mid, a.Platform); err != nil {
		log.Error("TxUpdateSnsUser(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateSnsUser update sns user.
func (d *Dao) TxUpdateSnsUser(tx *xsql.Tx, mid, expires int64, unionID string, platform int) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateSnsUserSQL, unionID, expires, mid, platform); err != nil {
		log.Error("TxUpdateSnsUser mid(%d) platform(%d) unionID(%s) expires(%d) d.snsDB.Exec() error(%+v)", mid, platform, unionID, platform, err)
		return
	}
	return res.RowsAffected()
}

// UpdateSnsUser update sns user.
func (d *Dao) UpdateSnsUser(c context.Context, mid, expires int64, unionID string, platform int) (affected int64, err error) {
	var res sql.Result
	if res, err = d.snsDB.Exec(c, _updateSnsUserSQL, unionID, expires, mid, platform); err != nil {
		log.Error("UpdateSnsUser mid(%d) platform(%d) unionID(%s) expires(%d) d.snsDB.Exec() error(%+v)", mid, platform, unionID, platform, err)
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

// UpdateSnsToken update sns token.
func (d *Dao) UpdateSnsToken(c context.Context, a *model.SnsToken) (affected int64, err error) {
	var res sql.Result
	if res, err = d.snsDB.Exec(c, _updateSnsTokenSQL, a.Token, a.Expires, a.Mid, a.Platform); err != nil {
		log.Error("UpdateSnsToken(%+v) tx.Exec() error(%+v)", a, err)
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
