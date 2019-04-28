package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/passport-game-data/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getAsoAccountRangeCloudSQL = "SELECT mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak,ctime,mtime FROM aso_account WHERE mtime>=? AND mtime<?"
	_getAsoAccountsCloudSQL     = "SELECT mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak,ctime,mtime FROM aso_account WHERE mid in(%s)"

	_getRangeAsoAccountsLocalSQL = "SELECT mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak,modify_time FROM aso_account WHERE modify_time>=? AND modify_time<?"
	_getAsoAccountsLocalSQL      = "SELECT mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak,modify_time FROM aso_account WHERE mid in(%s)"

	_updateAsoAccountCloudSQL = "UPDATE aso_account SET userid=?,uname=?,pwd=?,salt=?,email=?,tel=?,country_id=?,mobile_verified=?,isleak=? WHERE mid=? AND mtime=?"

	_addIgnoreAsoAccountCloudSQL  = "INSERT IGNORE INTO aso_account (mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak) VALUES(?,?,?,?,?,?,?,?,?,?)"
	_addIgnoreAsoAccountsCloudSQL = "INSERT IGNORE INTO aso_account (mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak) VALUES %s"
)

// AddAsoAccountsCloud batch add aso account to cloud.
func (d *Dao) AddAsoAccountsCloud(c context.Context, vs []*model.AsoAccount) (err error) {
	if len(vs) == 0 {
		return
	}

	var args = make([]string, 0)
	for _, v := range vs {
		args = append(args, getValues(v))
	}

	s := fmt.Sprintf(_addIgnoreAsoAccountsCloudSQL, strings.Join(args, ","))
	if _, err = d.cloudDB.Exec(c, s); err != nil {
		log.Error("d.cloudDB.Exec(%s) error(%v)", s, err)
	}
	return
}

func getValues(a *model.AsoAccount) string {
	email := "NULL"
	tel := "NULL"
	if len(a.Email) > 0 {
		email = "'" + a.Email + "'"
	}

	if len(a.Tel) > 0 {
		tel = "'" + a.Tel + "'"
	}

	return fmt.Sprintf(`(%d,'%s','%s','%s','%s',%s,%s,%d,%d,%d)`, a.Mid, a.UserID, a.Uname, a.Pwd, a.Salt, email, tel, a.CountryID, a.MobileVerified, a.Isleak)
}

// AsoAccountRangeCloud get aso account from cloud.
func (d *Dao) AsoAccountRangeCloud(c context.Context, st, ed time.Time) (res []*model.AsoAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.cloudDB.Query(c, _getAsoAccountRangeCloudSQL, st, ed); err != nil {
		log.Error("get aso account range cloud, dao.cloudDB.Query(%s) error(%v)", _getAsoAccountRangeCloudSQL, err)
		return
	}
	for rows.Next() {
		a := new(model.AsoAccount)
		var telPtr, emailPtr *string
		if err = rows.Scan(&a.Mid, &a.UserID, &a.Uname, &a.Pwd, &a.Salt, &emailPtr, &telPtr, &a.CountryID, &a.MobileVerified, &a.Isleak, &a.Ctime, &a.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				res = nil
				return
			}
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if telPtr != nil {
			a.Tel = *telPtr
		}
		if emailPtr != nil {
			a.Email = *emailPtr
		}
		res = append(res, a)
	}
	err = rows.Err()
	return
}

// AsoAccountsCloud get aso accounts from cloud.
func (d *Dao) AsoAccountsCloud(c context.Context, vs []int64) (res []*model.AsoAccount, err error) {
	if len(vs) == 0 {
		return
	}
	var args = make([]string, 0, len(vs))
	for _, v := range vs {
		args = append(args, fmt.Sprintf(`'%d'`, v))
	}
	if len(args) == 0 {
		return
	}
	s := fmt.Sprintf(_getAsoAccountsCloudSQL, strings.Join(args, ","))

	var rows *xsql.Rows
	if rows, err = d.cloudDB.Query(c, s); err != nil {
		log.Error("get aso accounts cloud, dao.cloudDB.Query(%s) error(%v)", s, err)
		return
	}
	for rows.Next() {
		a := new(model.AsoAccount)
		var telPtr, emailPtr *string
		if err = rows.Scan(&a.Mid, &a.UserID, &a.Uname, &a.Pwd, &a.Salt, &emailPtr, &telPtr, &a.CountryID, &a.MobileVerified, &a.Isleak, &a.Ctime, &a.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				res = nil
				return
			}
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if telPtr != nil {
			a.Tel = *telPtr
		}
		if emailPtr != nil {
			a.Email = *emailPtr
		}
		res = append(res, a)
	}
	err = rows.Err()
	return
}

// AsoAccountRangeLocal get aso account from local range start and end time.
func (d *Dao) AsoAccountRangeLocal(c context.Context, st, ed time.Time) (res []*model.OriginAsoAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.localDB.Query(c, _getRangeAsoAccountsLocalSQL, st, ed); err != nil {
		log.Error("get aso account range local, dao.localDB.Query(%s) error(%v)", _getRangeAsoAccountsLocalSQL, err)
		return
	}
	for rows.Next() {
		a := new(model.OriginAsoAccount)
		var telPtr, emailPtr *string
		if err = rows.Scan(&a.Mid, &a.UserID, &a.Uname, &a.Pwd, &a.Salt, &emailPtr, &telPtr, &a.CountryID, &a.MobileVerified, &a.Isleak, &a.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				res = nil
				return
			}
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if telPtr != nil {
			a.Tel = *telPtr
		}
		if emailPtr != nil {
			a.Email = *emailPtr
		}
		res = append(res, a)
	}
	err = rows.Err()
	return
}

// AsoAccountsLocal get aso accounts from origin.
func (d *Dao) AsoAccountsLocal(c context.Context, vs []int64) (res []*model.OriginAsoAccount, err error) {
	if len(vs) == 0 {
		return
	}
	var args = make([]string, 0, len(vs))
	for _, v := range vs {
		args = append(args, fmt.Sprintf(`'%d'`, v))
	}
	if len(args) == 0 {
		return
	}
	s := fmt.Sprintf(_getAsoAccountsLocalSQL, strings.Join(args, ","))

	var rows *xsql.Rows
	if rows, err = d.localDB.Query(c, s); err != nil {
		log.Error("get aso accounts local, dao.localDB.Query(%s) error(%v)", s, err)
		return
	}
	for rows.Next() {
		a := new(model.OriginAsoAccount)
		var telPtr, emailPtr *string
		if err = rows.Scan(&a.Mid, &a.UserID, &a.Uname, &a.Pwd, &a.Salt, &emailPtr, &telPtr, &a.CountryID, &a.MobileVerified, &a.Isleak, &a.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				res = nil
				return
			}
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if telPtr != nil {
			a.Tel = *telPtr
		}
		if emailPtr != nil {
			a.Email = *emailPtr
		}
		res = append(res, a)
	}
	err = rows.Err()
	return
}

// UpdateAsoAccountCloud update aso account.
func (d *Dao) UpdateAsoAccountCloud(c context.Context, a *model.AsoAccount, mt time.Time) (affected int64, err error) {
	var telPtr, emailPtr *string
	if a.Tel != "" {
		telPtr = &a.Tel
	}
	if a.Email != "" {
		emailPtr = &a.Email
	}
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _updateAsoAccountCloudSQL, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, telPtr, a.CountryID, a.MobileVerified, a.Isleak, a.Mid, mt); err != nil {
		log.Error("failed to update aso account, dao.cloudDB.Exec(%s) email(%s) tel(%s) error(%v)", _updateAsoAccountCloudSQL, emailPtr, telPtr, err)
		return
	}
	return res.RowsAffected()
}

// AddIgnoreAsoAccount add ignore aso account.
func (d *Dao) AddIgnoreAsoAccount(c context.Context, a *model.AsoAccount) (affected int64, err error) {
	var res sql.Result
	var telPtr, emailPtr *string
	if a.Tel != "" {
		telPtr = &a.Tel
	}
	if a.Email != "" {
		emailPtr = &a.Email
	}
	if res, err = d.cloudDB.Exec(c, _addIgnoreAsoAccountCloudSQL, a.Mid, a.UserID, a.Uname, a.Pwd, a.Salt, emailPtr, telPtr, a.CountryID, a.MobileVerified, a.Isleak); err != nil {
		log.Error("failed to add ignore aso account, dao.cloudDB.Exec(%s) email(%s) tel(%s) error(%s)", _addIgnoreAsoAccountCloudSQL, a.Email, a.Tel, err)
		return
	}
	return res.RowsAffected()
}
