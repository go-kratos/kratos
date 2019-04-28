package dao

import (
	"context"
	"database/sql"

	"go-common/app/job/main/passport-encrypt/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getOriginAccountSQL = "SELECT mid,userid,uname,pwd,salt,email,tel,mobile_verified,isleak,country_id,modify_time FROM aso_account where mid >= ? and mid < ?"
	_addAsoAccountSQL    = "INSERT INTO aso_account (mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak) VALUES(?,?,?,?,?,?,?,?,?,?)"
	_updateAsoAccountSQL = "UPDATE aso_account SET userid=?,uname=?,pwd=?,salt=?,email=?,tel=?,country_id=?,mobile_verified=?,isleak=? WHERE mid=?"
	_deleteAsoAccountSQL = "DELETE FROM aso_account WHERE mid=?"
)

// AsoAccounts get tokens by mid.
func (d *Dao) AsoAccounts(c context.Context, start, end int64) (res []*model.OriginAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, _getOriginAccountSQL, start, end); err != nil {
		log.Error("failed to get AsoAccounts, dao.originDB.AsoAccounts(%s) error(%v)", _getOriginAccountSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var telPtr, emailPtr *string
		r := new(model.OriginAccount)
		if err = rows.Scan(&r.Mid, &r.UserID, &r.Uname, &r.Pwd, &r.Salt, &emailPtr, &telPtr, &r.MobileVerified, &r.Isleak, &r.CountryID, &r.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		if telPtr != nil {
			r.Tel = *telPtr
		}
		if emailPtr != nil {
			r.Email = *emailPtr
		}
		res = append(res, r)
	}
	return
}

// AddAsoAccount add aso account.
func (d *Dao) AddAsoAccount(c context.Context, a *model.EncryptAccount) (affected int64, err error) {
	var res sql.Result
	var emailPtr *string

	if a.Email != "" {
		emailPtr = &a.Email
	}
	if res, err = d.encryptDB.Exec(c, _addAsoAccountSQL, a.Mid, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, a.Tel, a.CountryID, a.MobileVerified, a.Isleak); err != nil {
		log.Error("failed to add aso account, dao.encryptDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateAsoAccount update aso account.
func (d *Dao) UpdateAsoAccount(c context.Context, a *model.EncryptAccount) (affected int64, err error) {
	var res sql.Result
	var emailPtr *string
	if a.Email != "" {
		emailPtr = &a.Email
	}
	if res, err = d.encryptDB.Exec(c, _updateAsoAccountSQL, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, a.Tel, a.CountryID, a.MobileVerified, a.Isleak, a.Mid); err != nil {
		log.Error("failed to update aso account, dao.encryptDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelAsoAccount delete aso account.
func (d *Dao) DelAsoAccount(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.encryptDB.Exec(c, _deleteAsoAccountSQL, mid); err != nil {
		log.Error("failed to delete aso account, dao.encryptDB.Exec(%s) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}
