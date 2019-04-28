package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/passport-game-cloud/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_memberShard         = 30
	_addMemberInfoSQL    = "INSERT INTO member_%02d (mid,face) VALUES(?,?) ON DUPLICATE KEY UPDATE face=?"
	_deleteMemberInfoSQL = "DELETE FROM member_%02d WHERE mid=?"
	_getMemberInfoSQL    = "SELECT mid,face FROM member_%02d WHERE mid=?"
	_addTokenSQL         = "INSERT INTO app_perm (mid,appid,app_subid,access_token,create_at,expires) VALUES (?,?,?,?,?,?)"
	_updateTokenSQL      = "UPDATE app_perm SET expires=? WHERE access_token=? AND expires<?"
	_deleteTokenSQL      = "DELETE FROM app_perm WHERE access_token=?"
	_getTokensSQL        = "SELECT access_token FROM app_perm WHERE mid=?"
	_addAsoAccountSQL    = "INSERT INTO aso_account (mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak) VALUES(?,?,?,?,?,?,?,?,?,?)"
	_updateAsoAccountSQL = "UPDATE aso_account SET userid=?,uname=?,pwd=?,salt=?,email=?,tel=?,country_id=?,mobile_verified=?,isleak=? WHERE mid=?"
	_deleteAsoAccountSQL = "DELETE FROM aso_account WHERE mid=?"
)

func hit(mid int64) int64 {
	return mid % _memberShard
}

// AddMemberInfo add member info.
func (d *Dao) AddMemberInfo(c context.Context, info *model.Info) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, fmt.Sprintf(_addMemberInfoSQL, hit(info.Mid)), info.Mid, info.Face, info.Face); err != nil {
		log.Error("failed to add member info, dao.cloudDB.Exec(%d, %s, %s) error(%v)", info.Mid, info.Face, info.Face, err)
		return
	}
	return res.RowsAffected()
}

// DelMemberInfo delete member info.
func (d *Dao) DelMemberInfo(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, fmt.Sprintf(_deleteMemberInfoSQL, hit(mid)), mid); err != nil {
		log.Error("failed to delete member info, dao.cloudDB.Exec(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// AddToken add token.
func (d *Dao) AddToken(c context.Context, t *model.Perm) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _addTokenSQL, t.Mid, t.AppID, t.AppSubID, t.AccessToken, t.CreateAt, t.Expires); err != nil {
		log.Error("failed to add token, dao.cloudDB.Exec(%d, %d, %d, %s, %d, %d) error(%v)", t.Mid, t.AppID, t.AppSubID, t.AccessToken, t.CreateAt, t.Expires, err)
		return
	}
	return res.RowsAffected()
}

// UpdateToken update token.
func (d *Dao) UpdateToken(c context.Context, t *model.Perm) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _updateTokenSQL, t.Expires, t.AccessToken, t.Expires); err != nil {
		log.Error("failed to update token, dao.cloudDB.Exec(%d, %s, %d) error(%v)", t.Expires, t.AccessToken, t.Expires, err)
		return
	}
	return res.RowsAffected()
}

// DelToken delete token.
func (d *Dao) DelToken(c context.Context, accessToken string) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _deleteTokenSQL, accessToken); err != nil {
		log.Error("failed to delete token, dao.cloudDB.Exec(%s) error(%v)", accessToken, err)
		return
	}
	return res.RowsAffected()
}

// Tokens get tokens by mid.
func (d *Dao) Tokens(c context.Context, mid int64) (res []string, err error) {
	var rows *xsql.Rows
	if rows, err = d.cloudDB.Query(c, _getTokensSQL, mid); err != nil {
		log.Error("failed to get tokens, dao.cloudDB.Tokens(%s) error(%v)", _getTokensSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var at string
		if err = rows.Scan(&at); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, at)
	}
	return
}

// AddAsoAccount add aso account.
func (d *Dao) AddAsoAccount(c context.Context, a *model.AsoAccount) (affected int64, err error) {
	var res sql.Result
	var telPtr, emailPtr *string
	if a.Tel != "" {
		telPtr = &a.Tel
	}
	if a.Email != "" {
		emailPtr = &a.Email
	}
	if res, err = d.cloudDB.Exec(c, _addAsoAccountSQL, a.Mid, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, telPtr, a.CountryID, a.MobileVerified, a.Isleak); err != nil {
		log.Error("failed to add aso account, dao.cloudDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateAsoAccount update aso account.
func (d *Dao) UpdateAsoAccount(c context.Context, a *model.AsoAccount) (affected int64, err error) {
	var res sql.Result
	var telPtr, emailPtr *string
	if a.Tel != "" {
		telPtr = &a.Tel
	}
	if a.Email != "" {
		emailPtr = &a.Email
	}
	if res, err = d.cloudDB.Exec(c, _updateAsoAccountSQL, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, telPtr, a.CountryID, a.MobileVerified, a.Isleak, a.Mid); err != nil {
		log.Error("failed to add aso account, dao.cloudDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelAsoAccount delete aso account.
func (d *Dao) DelAsoAccount(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _deleteAsoAccountSQL, mid); err != nil {
		log.Error("failed to delete aso account, dao.cloudDB.Exec(%s) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}
