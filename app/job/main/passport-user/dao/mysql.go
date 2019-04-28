package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/job/main/passport-user/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getOriginAccountSQL          = "SELECT mid,userid,uname,pwd,salt,email,tel,mobile_verified,isleak,country_id,modify_time FROM aso_account WHERE mid > ? limit ?"
	_getOriginAccountInfoSQL      = "SELECT id,mid,spacesta,safe_question,safe_answer,join_time,join_ip,join_ip_v6,port,active_time,modify_time FROM aso_account_info%d WHERE id > ? limit ?"
	_getOriginAccountRegSQL       = "SELECT id,mid,origintype,regtype,appid,ctime,mtime FROM aso_account_reg_origin_%d WHERE id > ? limit ?"
	_getOriginAccountSnsSQL       = "SELECT mid,sina_uid,sina_access_token,sina_access_expires,qq_openid,qq_access_token,qq_access_expires FROM aso_account_sns WHERE mid > ? limit ?"
	_getOriginCountryCodeSQL      = "SELECT id,code,cname,rank,type,ename FROM aso_country_code"
	_getOriginTelBindLogSQL       = "SELECT timestamp FROM aso_telephone_bind_log WHERE mid = ? order by timestamp limit 1"
	_getOriginAccountByMidSQL     = "SELECT mid,userid,uname,pwd,salt,email,tel,mobile_verified,isleak,country_id,modify_time FROM aso_account WHERE mid = ?"
	_getOriginAccountSnsByMidSQL  = "SELECT mid,sina_uid,sina_access_token,sina_access_expires,qq_openid,qq_access_token,qq_access_expires FROM aso_account_sns WHERE mid = ?"
	_getOriginAccountInfoByMidSQL = "SELECT id,mid,spacesta,safe_question,safe_answer,join_time,join_ip,join_ip_v6,port,active_time,modify_time FROM aso_account_info%d WHERE mid = ?"
	_getOriginAccountRegByMidSQL  = "SELECT id,mid,origintype,regtype,appid,ctime,mtime FROM aso_account_reg_origin_%d WHERE mid = ?"
)

// AsoAccount get account by mid.
func (d *Dao) AsoAccount(c context.Context, start, count int64) (res []*model.OriginAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, _getOriginAccountSQL, start, count); err != nil {
		log.Error("fail to get AsoAccount, dao.originDB.Query(%s) error(%v)", _getOriginAccountSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var telPtr, emailPtr *string
		r := new(model.OriginAccount)
		if err = rows.Scan(&r.Mid, &r.UserID, &r.Uname, &r.Pwd, &r.Salt, &emailPtr, &telPtr, &r.MobileVerified, &r.Isleak, &r.CountryID, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		if telPtr != nil {
			r.Tel = *telPtr
		}
		if emailPtr != nil {
			r.Email = strings.ToLower(*emailPtr)
		}
		res = append(res, r)
	}
	return
}

// AsoAccountInfo get account info by id.
func (d *Dao) AsoAccountInfo(c context.Context, start, count, suffix int64) (res []*model.OriginAccountInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, fmt.Sprintf(_getOriginAccountInfoSQL, suffix), start, count); err != nil {
		log.Error("fail to get AsoAccountInfo, dao.originDB.Query(%s) error(%v)", _getOriginAccountInfoSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Spacesta, &r.SafeQuestion, &r.SafeAnswer, &r.JoinTime, &r.JoinIP, &r.JoinIPV6, &r.Port, &r.ActiveTime, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// AsoAccountReg get account reg by id.
func (d *Dao) AsoAccountReg(c context.Context, start, count, suffix int64) (res []*model.OriginAccountReg, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, fmt.Sprintf(_getOriginAccountRegSQL, suffix), start, count); err != nil {
		log.Error("fail to get AsoAccountReg, dao.originDB.Query(%s) error(%v)", _getOriginAccountRegSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountReg)
		if err = rows.Scan(&r.ID, &r.Mid, &r.OriginType, &r.RegType, &r.AppID, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// AsoAccountSns get account sns by id.
func (d *Dao) AsoAccountSns(c context.Context, start, count int64) (res []*model.OriginAccountSns, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, _getOriginAccountSnsSQL, start, count); err != nil {
		log.Error("fail to get AsoAccountSns, dao.originDB.Query(%s) error(%v)", _getOriginAccountSnsSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountSns)
		if err = rows.Scan(&r.Mid, &r.SinaUID, &r.SinaAccessToken, &r.SinaAccessExpires, &r.QQOpenid, &r.QQAccessToken, &r.QQAccessExpires); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// AsoCountryCode get aso country code.
func (d *Dao) AsoCountryCode(c context.Context) (res []*model.CountryCode, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, _getOriginCountryCodeSQL); err != nil {
		log.Error("fail to get AsoCountryCode, dao.originDB.Query(%s) error(%v)", _getOriginCountryCodeSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CountryCode)
		if err = rows.Scan(&r.ID, &r.Code, &r.Cname, &r.Rank, &r.Type, &r.Ename); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// AsoTelBindLog get aso tel bind log.
func (d *Dao) AsoTelBindLog(c context.Context, mid int64) (res int64, err error) {
	if err = d.originDB.QueryRow(c, _getOriginTelBindLogSQL, mid).Scan(&res); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("fail to get AsoTelBindLog, dao.originDB.QueryRow(%s) error(%v)", _getOriginTelBindLogSQL, err)
		}
		return
	}
	return
}

// GetAsoAccountByMid get aso account by mid.
func (d *Dao) GetAsoAccountByMid(c context.Context, mid int64) (res *model.OriginAccount, err error) {
	row := d.originDB.QueryRow(c, _getOriginAccountByMidSQL, mid)
	var telPtr, emailPtr *string
	res = &model.OriginAccount{}
	if err = row.Scan(&res.Mid, &res.UserID, &res.Uname, &res.Pwd, &res.Salt, &emailPtr, &telPtr, &res.MobileVerified, &res.Isleak, &res.CountryID, &res.MTime); err != nil {
		log.Error("fail to get AsoAccount, dao.originDB.QueryRow(%s) error(%v)", _getOriginAccountByMidSQL, err)
		return
	}
	if telPtr != nil {
		res.Tel = *telPtr
	}
	if emailPtr != nil {
		res.Email = strings.ToLower(*emailPtr)
	}
	return
}

// GetAsoAccountSnsByMid get aso account sns by mid.
func (d *Dao) GetAsoAccountSnsByMid(c context.Context, mid int64) (res *model.OriginAccountSns, err error) {
	row := d.originDB.QueryRow(c, _getOriginAccountSnsByMidSQL, mid)
	res = &model.OriginAccountSns{}
	if err = row.Scan(&res.Mid, &res.SinaUID, &res.SinaAccessToken, &res.SinaAccessExpires, &res.QQOpenid, &res.QQAccessToken, &res.QQAccessExpires); err != nil {
		log.Error("fail to get AsoAccountSns, dao.originDB.QueryRow(%s) error(%v)", _getOriginAccountSnsByMidSQL, err)
		return
	}
	return
}

// GetAsoAccountInfoByMid get aso account info by mid.
func (d *Dao) GetAsoAccountInfoByMid(c context.Context, mid int64) (res *model.OriginAccountInfo, err error) {
	row := d.originDB.QueryRow(c, fmt.Sprintf(_getOriginAccountInfoByMidSQL, accountInfoIndex(mid)), mid)
	res = &model.OriginAccountInfo{}
	if err = row.Scan(&res.ID, &res.Mid, &res.Spacesta, &res.SafeQuestion, &res.SafeAnswer, &res.JoinTime, &res.JoinIP, &res.JoinIPV6, &res.Port, &res.ActiveTime, &res.MTime); err != nil {
		log.Error("fail to get AsoAccountInfo, dao.originDB.QueryRow(%s) error(%v)", _getOriginAccountInfoByMidSQL, err)
		return
	}
	return
}

// GetAsoAccountRegByMid get aso account reg by mid.
func (d *Dao) GetAsoAccountRegByMid(c context.Context, mid int64) (res *model.OriginAccountReg, err error) {
	row := d.originDB.QueryRow(c, fmt.Sprintf(_getOriginAccountRegByMidSQL, accountRegIndex(mid)), mid)
	res = &model.OriginAccountReg{}
	if err = row.Scan(&res.ID, &res.Mid, &res.OriginType, &res.RegType, &res.AppID, &res.CTime, &res.MTime); err != nil {
		log.Error("fail to get AsoAccountReg, dao.originDB.QueryRow(%s) error(%v)", _getOriginAccountRegByMidSQL, err)
		return
	}
	return
}

func accountInfoIndex(mid int64) int64 {
	return mid % 30
}

func accountRegIndex(mid int64) int64 {
	return mid % 20
}
