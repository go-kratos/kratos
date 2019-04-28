package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/passport-user-compare/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

var (
	selectAccountSQL           = "SELECT mid, userid, uname, pwd, salt, email, tel, mobile_verified, isleak, country_id, modify_time FROM aso_account WHERE mid > ? limit ? "
	selectAccountByMidSQL      = "SELECT mid, userid, uname, pwd, salt, email, tel, mobile_verified, isleak, country_id, modify_time FROM aso_account WHERE mid = ?"
	selectAccountByTelSQL      = "SELECT mid, userid, uname, pwd, salt, email, tel, mobile_verified, isleak, country_id, modify_time FROM aso_account WHERE tel = ?"
	selectAccountByMailSQL     = "SELECT mid, userid, uname, pwd, salt, email, tel, mobile_verified, isleak, country_id, modify_time FROM aso_account WHERE email = ?"
	selectAccountByTimeSQL     = "SELECT mid, userid, uname, pwd, salt, email, tel, mobile_verified, isleak, country_id, modify_time FROM aso_account WHERE modify_time >= ? AND modify_time < ?"
	selectAccountInfoSQL       = "SELECT id, mid, spacesta, safe_question, safe_answer, join_time, join_ip, active_time, modify_time FROM aso_account_info%d WHERE id > ? limit ?"
	selectAccountInfoByMidSQL  = "SELECT id, mid, spacesta, safe_question, safe_answer, join_time, join_ip, active_time, modify_time FROM aso_account_info%d WHERE mid = ? "
	selectAccountInfoByTimeSQL = "SELECT id, mid, spacesta, safe_question, safe_answer, join_time, join_ip, active_time, modify_time FROM aso_account_info%d WHERE modify_time >= ? AND modify_time < ? "
	selectAccountSnsSQL        = "SELECT mid, sina_uid, sina_access_token, sina_access_expires, qq_openid, qq_access_token, qq_access_expires FROM aso_account_sns WHERE mid > ? limit ?"
	selectAccountSnsByMidSQL   = "SELECT mid, sina_uid, sina_access_token, sina_access_expires, qq_openid, qq_access_token, qq_access_expires FROM aso_account_sns WHERE mid = ? "
	selectOriginTelBindLogSQL  = "SELECT timestamp FROM aso_telephone_bind_log WHERE mid = ? order by timestamp limit 1"
	selectAccountRegByMidSQL   = "SELECT id, mid, origintype, regtype, appid, ctime, mtime FROM aso_account_reg_origin_%d WHERE mid = ? "
	selectAccountRegByTimeSQL  = "SELECT id, mid, origintype, regtype, appid, ctime, mtime FROM aso_account_reg_origin_%d WHERE mtime >= ? AND mtime < ? "
)

// BatchQueryAccount batch query account
func (d *Dao) BatchQueryAccount(c context.Context, start, limit int64) (res []*model.OriginAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, selectAccountSQL, start, limit); err != nil {
		log.Error("fail to get BatchQueryAccount, dao.originDB.Query(%s) error(%v)", selectAccountSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var telPtr, emailPtr *string
		r := new(model.OriginAccount)
		if err = rows.Scan(&r.Mid, &r.UserID, &r.Uname, &r.Pwd, &r.Salt, &emailPtr, &telPtr, &r.MobileVerified, &r.Isleak, &r.CountryID, &r.MTime); err != nil {
			log.Error("BatchQueryAccount row.Scan() error(%v)", err)
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

// QueryAccountByMid query account by mid
func (d *Dao) QueryAccountByMid(c context.Context, mid int64) (res *model.OriginAccount, err error) {
	row := d.originDB.QueryRow(c, selectAccountByMidSQL, mid)
	res = new(model.OriginAccount)
	var telPtr, emailPtr *string
	if err = row.Scan(&res.Mid, &res.UserID, &res.Uname, &res.Pwd, &res.Salt, &emailPtr, &telPtr, &res.MobileVerified, &res.Isleak, &res.CountryID, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountByMid row.Scan() error(%v)", err)
		}
		return
	}
	if telPtr != nil {
		res.Tel = *telPtr
	}
	if emailPtr != nil {
		res.Email = *emailPtr
	}
	return
}

// QueryAccountByTel query account by mid
func (d *Dao) QueryAccountByTel(c context.Context, tel string) (res *model.OriginAccount, err error) {
	row := d.originDB.QueryRow(c, selectAccountByTelSQL, tel)
	res = new(model.OriginAccount)
	var telPtr, emailPtr *string
	if err = row.Scan(&res.Mid, &res.UserID, &res.Uname, &res.Pwd, &res.Salt, &emailPtr, &telPtr, &res.MobileVerified, &res.Isleak, &res.CountryID, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountByTel row.Scan() error(%v)", err)
		}
		return
	}
	if telPtr != nil {
		res.Tel = *telPtr
	}
	if emailPtr != nil {
		res.Email = *emailPtr
	}
	return
}

// QueryAccountByMail query account by mid
func (d *Dao) QueryAccountByMail(c context.Context, mail string) (res *model.OriginAccount, err error) {
	row := d.originDB.QueryRow(c, selectAccountByMailSQL, mail)
	res = new(model.OriginAccount)
	var telPtr, emailPtr *string
	if err = row.Scan(&res.Mid, &res.UserID, &res.Uname, &res.Pwd, &res.Salt, &emailPtr, &telPtr, &res.MobileVerified, &res.Isleak, &res.CountryID, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountByMail row.Scan() error(%v)", err)
		}
		return
	}
	if telPtr != nil {
		res.Tel = *telPtr
	}
	if emailPtr != nil {
		res.Email = *emailPtr
	}
	return
}

// BatchQueryAccountByTime batch query by time
func (d *Dao) BatchQueryAccountByTime(c context.Context, start, end time.Time) (res []*model.OriginAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, selectAccountByTimeSQL, start, end); err != nil {
		log.Error("fail to get BatchQueryAccountByTime, dao.originDB.Query(%s) error(%v)", selectAccountSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var telPtr, emailPtr *string
		r := new(model.OriginAccount)
		if err = rows.Scan(&r.Mid, &r.UserID, &r.Uname, &r.Pwd, &r.Salt, &emailPtr, &telPtr, &r.MobileVerified, &r.Isleak, &r.CountryID, &r.MTime); err != nil {
			log.Error("BatchQueryAccountByTime row.Scan() error(%v)", err)
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

// BatchQueryAccountInfo batch query account info
func (d *Dao) BatchQueryAccountInfo(c context.Context, start, limit int64, suffix int) (res []*model.OriginAccountInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, fmt.Sprintf(selectAccountInfoSQL, suffix), start, limit); err != nil {
		log.Error("fail to get BatchQueryAccountInfo, dao.originDB.Query(%s) error(%v)", selectAccountInfoSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Spacesta, &r.SafeQuestion, &r.SafeAnswer, &r.JoinTime, &r.JoinIP, &r.ActiveTime, &r.MTime); err != nil {
			log.Error("BatchQueryAccountInfo row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// QueryAccountInfoByMid query account info by mid
func (d *Dao) QueryAccountInfoByMid(c context.Context, mid int64) (res *model.OriginAccountInfo, err error) {
	row := d.originDB.QueryRow(c, fmt.Sprintf(selectAccountInfoByMidSQL, mid%30), mid)
	res = new(model.OriginAccountInfo)
	if err = row.Scan(&res.ID, &res.Mid, &res.Spacesta, &res.SafeQuestion, &res.SafeAnswer, &res.JoinTime, &res.JoinIP, &res.ActiveTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountInfoByMid row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// BatchQueryAccountInfoByTime batch query account info
func (d *Dao) BatchQueryAccountInfoByTime(c context.Context, start, end time.Time, suffix int) (res []*model.OriginAccountInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, fmt.Sprintf(selectAccountInfoByTimeSQL, suffix), start, end); err != nil {
		log.Error("fail to get BatchQueryAccountInfoByTime, dao.originDB.Query(%s) error(%v)", selectAccountInfoSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Spacesta, &r.SafeQuestion, &r.SafeAnswer, &r.JoinTime, &r.JoinIP, &r.ActiveTime, &r.MTime); err != nil {
			log.Error("BatchQueryAccountInfoByTime row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// BatchQueryAccountSns batch query account sns
func (d *Dao) BatchQueryAccountSns(c context.Context, start, limit int64) (res []*model.OriginAccountSns, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, selectAccountSnsSQL, start, limit); err != nil {
		log.Error("fail to get BatchQueryAccountSns, dao.originDB.Query(%s) error(%v)", selectAccountSnsSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountSns)
		if err = rows.Scan(&r.Mid, &r.SinaUID, &r.SinaAccessToken, &r.SinaAccessExpires, &r.QQOpenid, &r.QQAccessToken, &r.QQAccessExpires); err != nil {
			log.Error("BatchQueryAccountSns row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// QueryAccountSnsByMid query account sns by mid
func (d *Dao) QueryAccountSnsByMid(c context.Context, mid int64) (res *model.OriginAccountSns, err error) {
	row := d.originDB.QueryRow(c, selectAccountSnsByMidSQL, mid)
	res = new(model.OriginAccountSns)
	if err = row.Scan(&res.Mid, &res.SinaUID, &res.SinaAccessToken, &res.SinaAccessExpires, &res.QQOpenid, &res.QQAccessToken, &res.QQAccessExpires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountSnsByMid row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// QueryTelBindLog get aso tel bind log.
func (d *Dao) QueryTelBindLog(c context.Context, mid int64) (res int64, err error) {
	if err = d.originDB.QueryRow(c, selectOriginTelBindLogSQL, mid).Scan(&res); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("fail to get AsoTelBindLog, dao.originDB.QueryRow(%s) error(%v)", selectOriginTelBindLogSQL, err)
		}
		return
	}
	return
}

// QueryAccountRegByMid query account reg by mid
func (d *Dao) QueryAccountRegByMid(c context.Context, mid int64) (res *model.OriginAccountReg, err error) {
	row := d.originDB.QueryRow(c, fmt.Sprintf(selectAccountRegByMidSQL, mid%20), mid)
	res = new(model.OriginAccountReg)
	if err = row.Scan(&res.ID, &res.Mid, &res.OriginType, &res.RegType, &res.AppID, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("QueryAccountRegByMid row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// BatchQueryAccountRegByTime batch query account info
func (d *Dao) BatchQueryAccountRegByTime(c context.Context, start, end time.Time, suffix int) (res []*model.OriginAccountReg, err error) {
	var rows *xsql.Rows
	if rows, err = d.originDB.Query(c, fmt.Sprintf(selectAccountRegByTimeSQL, suffix), start, end); err != nil {
		log.Error("fail to get BatchQueryAccountRegByTime, dao.originDB.Query(%s) error(%v)", selectAccountInfoSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OriginAccountReg)
		if err = rows.Scan(&r.ID, &r.Mid, &r.OriginType, &r.RegType, &r.AppID, &r.CTime, &r.MTime); err != nil {
			log.Error("BatchQueryAccountInfoByTime row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}
