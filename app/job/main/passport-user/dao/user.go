package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/passport-user/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addCountryCodeSQL                   = "INSERT INTO country_code (id,code,cname,rank,type,ename) VALUES(?,?,?,?,?,?)"
	_getAesKeySQL                        = "SELECT `key` FROM user_secret WHERE key_type = 2"
	_getSaltSQL                          = "SELECT `key` FROM user_secret WHERE key_type = 3"
	_getCountryCodeMapSQL                = "SELECT id,code FROM country_code"
	_getUserTelSQL                       = "SELECT mid FROM user_tel WHERE mid > ? limit ?"
	_getUserEmailByMidSQL                = "SELECT mid,email,verified,email_bind_time,ctime,mtime FROM user_email WHERE mid = ?"
	_getUserTelByMidSQL                  = "SELECT mid,tel,cid,tel_bind_time,ctime,mtime FROM user_tel WHERE mid = ?"
	_addUserBaseSQL                      = "INSERT INTO user_base (mid,userid,pwd,salt,status,deleted,mtime) VALUES(?,?,?,?,?,?,?)"
	_addUserEmailSQL                     = "INSERT INTO user_email (mid,email,verified,email_bind_time,mtime) VALUES(?,?,?,?,?)"
	_addUserTelSQL                       = "INSERT INTO user_tel (mid,tel,cid,tel_bind_time,mtime) VALUES(?,?,?,?,?)"
	_addUserSafeQuestionSQL              = "INSERT INTO user_safe_question%02d (mid,safe_question,safe_answer,safe_bind_time) VALUES(?,?,?,?)"
	_addUserThirdBindSQL                 = "INSERT INTO user_third_bind (mid,openid,platform,token,expires) VALUES(?,?,?,?,?)"
	_updateUserBaseSQL                   = "UPDATE user_base SET userid=?,pwd=?,salt=?,status=? WHERE mid =?"
	_updateUserEmailSQL                  = "UPDATE user_email SET email=? WHERE mid =?"
	_updateUserEmailAndBindTimeSQL       = "UPDATE user_email SET email=?,email_bind_time=? WHERE mid =?"
	_updateUserTelSQL                    = "UPDATE user_tel SET tel=?,cid=? WHERE mid =?"
	_updateUserTelAndBindTimeSQL         = "UPDATE user_tel SET tel=?,cid=?,tel_bind_time=? WHERE mid =?"
	_updateUserSafeQuestionSQL           = "UPDATE user_safe_question%02d SET safe_question=?,safe_answer=? WHERE mid =?"
	_updateUserThirdBindSQL              = "UPDATE user_third_bind SET openid=?,token=?,expires=? WHERE mid =? and platform=?"
	_updateUserEmailVerifiedSQL          = "UPDATE user_email SET verified=? WHERE mid =?"
	_updateUserEmailBindTimeSQL          = "UPDATE user_email SET verified=?,email_bind_time=? WHERE mid =?"
	_updateUserTelBindTimeSQL            = "UPDATE user_tel SET tel_bind_time=? WHERE mid =?"
	_insertUpdateUserRegOriginSQL        = "INSERT INTO user_reg_origin%02d (mid,join_ip,join_ip_v6,port,join_time) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE join_ip=?,join_ip_v6=?,port=?,join_time=?"
	_insertUpdateUserRegOriginTypeSQL    = "INSERT INTO user_reg_origin%02d (mid,origin,reg_type,appid,ctime,mtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE origin=?,reg_type=?,appid=?,ctime=?,mtime=?"
	_delUserBase                         = "UPDATE user_base SET deleted=1 WHERE mid =?"
	_delUserTel                          = "UPDATE user_tel SET tel=null,cid=null WHERE mid =?"
	_delUserEmail                        = "UPDATE user_email SET email=null WHERE mid =?"
	_getMidByTelSQL                      = "SELECT mid FROM user_tel WHERE tel = ? and cid = ?"
	_getMidByEmailSQL                    = "SELECT mid FROM user_email WHERE email = ?"
	_getUserBaseByMidSQL                 = "SELECT mid,userid,pwd,salt,status,ctime,mtime FROM user_base WHERE mid = ?"
	_getUserSafeQuestionByMidSQL         = "SELECT mid,safe_question,safe_answer,safe_bind_time,ctime,mtime FROM user_safe_question%02d where mid = ?"
	_getUserThirdBindByMidSQL            = "SELECT id,mid,openid,platform,token,expires,ctime,mtime FROM user_third_bind where mid = ?"
	_getUserThirdBindByMidAndPlatformSQL = "SELECT id,mid,openid,platform,token,expires,ctime,mtime FROM user_third_bind where mid = ? and platform = ? limit 1"
	_getUserRegOriginByMidSQL            = "SELECT mid,join_ip,join_ip_v6,port,join_time,origin,reg_type,appid from user_reg_origin%02d where mid = ?"

	_insertUpdateUserBaseSQL         = "INSERT INTO user_base (mid,userid,pwd,salt,status,deleted,mtime) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE userid=?,pwd=?,salt=?,status=?"
	_insertIgnoreUserSafeQuestionSQL = "INSERT IGNORE INTO user_safe_question%02d (mid,safe_question,safe_answer,safe_bind_time) VALUES(?,?,?,?)"
	_insertIgnoreUserRegOriginSQL    = "INSERT IGNORE INTO user_reg_origin%02d (mid,join_ip,join_ip_v6,port,join_time) VALUES (?,?,?,?,?)"
	_delUserThirdBindSQL             = "UPDATE user_third_bind SET openid='',token='',expires=0 WHERE mid =?"
)

// AddCountryCode add country code.
func (d *Dao) AddCountryCode(c context.Context, a *model.CountryCode) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _addCountryCodeSQL, a.ID, a.Code, a.Cname, a.Rank, a.Type, a.Ename); err != nil {
		log.Error("fail to add country code, countryCode(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AesKey get aes key.
func (d *Dao) AesKey(c context.Context) (res string, err error) {
	if err = d.encryptDB.QueryRow(c, _getAesKeySQL).Scan(&res); err != nil {
		log.Error("fail to get AesKey, dao.encryptDB.QueryRow(%s) error(%v)", _getAesKeySQL, err)
		return
	}
	return
}

// Salt get salt.
func (d *Dao) Salt(c context.Context) (res string, err error) {
	if err = d.encryptDB.QueryRow(c, _getSaltSQL).Scan(&res); err != nil {
		log.Error("fail to get Salt, dao.encryptDB.QueryRow(%s) error(%v)", _getSaltSQL, err)
		return
	}
	return
}

// UserTel get user tel.
func (d *Dao) UserTel(c context.Context, start, count int64) (res []*model.UserTel, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, _getUserTelSQL, start, count); err != nil {
		log.Error("fail to get UserTel, dao.userDB.Query(%s) error(%v)", _getUserTelSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserTel)
		if err = rows.Scan(&r.Mid); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// GetUserEmailByMid get user email by mid.
func (d *Dao) GetUserEmailByMid(c context.Context, mid int64) (res *model.UserEmail, err error) {
	res = &model.UserEmail{}
	if err = d.userDB.QueryRow(c, _getUserEmailByMidSQL, mid).Scan(&res.Mid, &res.Email, &res.Verified, &res.EmailBindTime, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserEmail by mid(%d), dao.encryptDB.QueryRow(%s) error(%v)", mid, _getUserEmailByMidSQL, err)
		}
		return
	}
	return
}

// GetUserTelByMid get user email by mid.
func (d *Dao) GetUserTelByMid(c context.Context, mid int64) (res *model.UserTel, err error) {
	var cidPtr *string
	res = &model.UserTel{}
	if err = d.userDB.QueryRow(c, _getUserTelByMidSQL, mid).Scan(&res.Mid, &res.Tel, &cidPtr, &res.TelBindTime, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserTel by mid(%d), dao.encryptDB.QueryRow(%s) error(%v)", mid, _getUserTelByMidSQL, err)
		}
		return
	}
	if cidPtr != nil {
		res.Cid = *cidPtr
	}
	return
}

// AddUserBase add user base.
func (d *Dao) AddUserBase(c context.Context, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _addUserBaseSQL, a.Mid, a.UserID, a.Pwd, a.Salt, a.Status, a.Deleted, a.MTime); err != nil {
		log.Error("fail to add user base, userBase(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddUserBase add user base.
func (d *Dao) TxAddUserBase(tx *xsql.Tx, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_addUserBaseSQL, a.Mid, a.UserID, a.Pwd, a.Salt, a.Status, a.Deleted, a.MTime); err != nil {
		log.Error("fail to add user base, userBase(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AddUserEmail add user email.
func (d *Dao) AddUserEmail(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = d.userDB.Exec(c, _addUserEmailSQL, a.Mid, emailPtr, a.Verified, a.EmailBindTime, a.MTime); err != nil {
		log.Error("fail to add user email, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddUserEmail add user email.
func (d *Dao) TxAddUserEmail(tx *xsql.Tx, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = tx.Exec(_addUserEmailSQL, a.Mid, emailPtr, a.Verified, a.EmailBindTime, a.MTime); err != nil {
		log.Error("fail to add user email, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AddUserTel add user tel.
func (d *Dao) AddUserTel(c context.Context, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr *string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = &a.Cid
	}
	if res, err = d.userDB.Exec(c, _addUserTelSQL, a.Mid, telPtr, cidPtr, a.TelBindTime, a.MTime); err != nil {
		log.Error("fail to add user tel, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddUserTel add user tel.
func (d *Dao) TxAddUserTel(tx *xsql.Tx, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr *string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = &a.Cid
	}
	if res, err = tx.Exec(_addUserTelSQL, a.Mid, telPtr, cidPtr, a.TelBindTime, a.MTime); err != nil {
		log.Error("fail to add user tel, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AddUserSafeQuestion add user safe question.
func (d *Dao) AddUserSafeQuestion(c context.Context, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(_addUserSafeQuestionSQL, tableIndex(a.Mid)), a.Mid, a.SafeQuestion, a.SafeAnswer, a.SafeBindTime); err != nil {
		log.Error("fail to add user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddUserSafeQuestion add user safe question.
func (d *Dao) TxAddUserSafeQuestion(tx *xsql.Tx, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addUserSafeQuestionSQL, tableIndex(a.Mid)), a.Mid, a.SafeQuestion, a.SafeAnswer, a.SafeBindTime); err != nil {
		log.Error("fail to add user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AddUserThirdBind add user third bind.
func (d *Dao) AddUserThirdBind(c context.Context, a *model.UserThirdBind) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _addUserThirdBindSQL, a.Mid, a.OpenID, a.PlatForm, a.Token, a.Expires); err != nil {
		log.Error("fail to add user third bind, userThirdBind(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxAddUserThirdBind add user third bind.
func (d *Dao) TxAddUserThirdBind(tx *xsql.Tx, a *model.UserThirdBind) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_addUserThirdBindSQL, a.Mid, a.OpenID, a.PlatForm, a.Token, a.Expires); err != nil {
		log.Error("fail to add user third bind, userThirdBind(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserBase update user base.
func (d *Dao) UpdateUserBase(c context.Context, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateUserBaseSQL, a.UserID, a.Pwd, a.Salt, a.Status, a.Mid); err != nil {
		log.Error("fail to update user base, userBase(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserEmail update user email.
func (d *Dao) UpdateUserEmail(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = d.userDB.Exec(c, _updateUserEmailSQL, emailPtr, a.Mid); err != nil {
		log.Error("fail to update user email, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserEmailAndBindTime update user email and bind time.
func (d *Dao) UpdateUserEmailAndBindTime(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = d.userDB.Exec(c, _updateUserEmailAndBindTimeSQL, emailPtr, a.EmailBindTime, a.Mid); err != nil {
		log.Error("fail to update user email and bind time, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUserEmail update user email.
func (d *Dao) TxUpdateUserEmail(tx *xsql.Tx, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = tx.Exec(_updateUserEmailSQL, emailPtr, a.Mid); err != nil {
		log.Error("fail to update user email, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserTel update user tel.
func (d *Dao) UpdateUserTel(c context.Context, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr *string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = &a.Cid
	}
	if res, err = d.userDB.Exec(c, _updateUserTelSQL, telPtr, cidPtr, a.Mid); err != nil {
		log.Error("fail to update user tel, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserTelAndBindTime update user tel and bind time.
func (d *Dao) UpdateUserTelAndBindTime(c context.Context, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr *string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = &a.Cid
	}
	if res, err = d.userDB.Exec(c, _updateUserTelAndBindTimeSQL, telPtr, cidPtr, a.TelBindTime, a.Mid); err != nil {
		log.Error("fail to update user tel and bind time, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserSafeQuesion update user safe question.
func (d *Dao) UpdateUserSafeQuesion(c context.Context, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(_updateUserSafeQuestionSQL, tableIndex(a.Mid)), a.SafeQuestion, a.SafeAnswer, a.Mid); err != nil {
		log.Error("fail to update user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUserSafeQuesion update user safe question.
func (d *Dao) TxUpdateUserSafeQuesion(tx *xsql.Tx, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateUserSafeQuestionSQL, tableIndex(a.Mid)), a.SafeQuestion, a.SafeAnswer, a.Mid); err != nil {
		log.Error("fail to update user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserThirdBind update user third bind.
func (d *Dao) UpdateUserThirdBind(c context.Context, a *model.UserThirdBind) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateUserThirdBindSQL, a.OpenID, a.Token, a.Expires, a.Mid, a.PlatForm); err != nil {
		log.Error("fail to update user third bind, userThirdBind(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserEmailVerified update user email verified.
func (d *Dao) UpdateUserEmailVerified(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateUserEmailVerifiedSQL, a.Verified, a.Mid); err != nil {
		log.Error("fail to update user email verified, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUserEmailVerified update user email verified.
func (d *Dao) TxUpdateUserEmailVerified(tx *xsql.Tx, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateUserEmailVerifiedSQL, a.Verified, a.Mid); err != nil {
		log.Error("fail to update user email verified, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserEmailBindTime update user email bind time.
func (d *Dao) UpdateUserEmailBindTime(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateUserEmailBindTimeSQL, a.Verified, a.EmailBindTime, a.Mid); err != nil {
		log.Error("fail to update user email bind time, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUserEmailBindTime update user email bind time.
func (d *Dao) TxUpdateUserEmailBindTime(tx *xsql.Tx, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateUserEmailBindTimeSQL, a.Verified, a.EmailBindTime, a.Mid); err != nil {
		log.Error("fail to update user email bind time, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserTelBindTime update user tel bind time.
func (d *Dao) UpdateUserTelBindTime(c context.Context, a *model.UserTel) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateUserTelBindTimeSQL, a.TelBindTime, a.Mid); err != nil {
		log.Error("fail to update user tel bind time, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxDelUserBase update user base deleted = 1.
func (d *Dao) TxDelUserBase(tx *xsql.Tx, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_delUserBase, mid); err != nil {
		log.Error("fail to del user base, mid(%d) dao.userDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// TxDelUserTel update user tel deleted = 1.
func (d *Dao) TxDelUserTel(tx *xsql.Tx, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_delUserTel, mid); err != nil {
		log.Error("fail to del user tel, mid(%d) dao.userDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// TxDelUserEmail update user email deleted = 1.
func (d *Dao) TxDelUserEmail(tx *xsql.Tx, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_delUserEmail, mid); err != nil {
		log.Error("fail to del user email, mid(%d) dao.userDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertUpdateUserRegOrigin insert update user reg origin.
func (d *Dao) TxInsertUpdateUserRegOrigin(tx *xsql.Tx, a *model.UserRegOrigin) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_insertUpdateUserRegOriginSQL, tableIndex(a.Mid)), a.Mid, a.JoinIP, a.JoinIPV6, a.Port, a.JoinTime, a.JoinIP, a.JoinIPV6, a.Port, a.JoinTime); err != nil {
		log.Error("fail to insert update user reg origin, userRegOrigin(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// InsertUpdateUserRegOriginType insert update user reg origin type.
func (d *Dao) InsertUpdateUserRegOriginType(c context.Context, a *model.UserRegOrigin) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(_insertUpdateUserRegOriginTypeSQL, tableIndex(a.Mid)), a.Mid, a.Origin, a.RegType, a.AppID, a.CTime, a.MTime,
		a.Origin, a.RegType, a.AppID, a.CTime, a.MTime); err != nil {
		log.Error("fail to insert update user reg origin type, userRegOrigin(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// CountryCodeMap get country code map.
func (d *Dao) CountryCodeMap(c context.Context) (res map[int64]string, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, _getCountryCodeMapSQL); err != nil {
		log.Error("fail to get CountryCodeMap, dao.userDB.Query(%s) error(%+v)", _getCountryCodeMapSQL, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]string)
	for rows.Next() {
		var (
			id   int64
			code string
		)
		if err = rows.Scan(&id, &code); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res[id] = code
	}
	return
}

// GetMidByTel get mid by tel.
func (d *Dao) GetMidByTel(c context.Context, a *model.UserTel) (mid int64, err error) {
	if err = d.userDB.QueryRow(c, _getMidByTelSQL, a.Tel, a.Cid).Scan(&mid); err != nil {
		log.Error("fail to get mid by tel, dao.userDB.QueryRow(%s) error(%+v)", _getMidByTelSQL, err)
		return
	}
	return
}

// GetMidByEmail get mid by email.
func (d *Dao) GetMidByEmail(c context.Context, a *model.UserEmail) (mid int64, err error) {
	if err = d.userDB.QueryRow(c, _getMidByEmailSQL, a.Email).Scan(&mid); err != nil {
		log.Error("fail to get mid by email, dao.userDB.QueryRow(%s) error(%+v)", _getMidByEmailSQL, err)
		return
	}
	return
}

// GetUserBaseByMid get user base by mid.
func (d *Dao) GetUserBaseByMid(c context.Context, mid int64) (res *model.UserBase, err error) {
	res = &model.UserBase{}
	if err = d.userDB.QueryRow(c, _getUserBaseByMidSQL, mid).Scan(&res.Mid, &res.UserID, &res.Pwd, &res.Salt, &res.Status, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserBase by mid(%d), dao.userDB.QueryRow(%s) error(%+v)", mid, _getUserBaseByMidSQL, err)
		}
		return
	}
	return
}

// GetUserSafeQuestionByMid get user safe question by mid.
func (d *Dao) GetUserSafeQuestionByMid(c context.Context, mid int64) (res *model.UserSafeQuestion, err error) {
	res = &model.UserSafeQuestion{}
	if err = d.userDB.QueryRow(c, fmt.Sprintf(_getUserSafeQuestionByMidSQL, tableIndex(mid)), mid).Scan(&res.Mid, &res.SafeQuestion, &res.SafeAnswer, &res.SafeBindTime, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserSafeQuestion by mid(%d), dao.userDB.QueryRow(%s) error(%+v)", mid, _getUserSafeQuestionByMidSQL, err)
		}
		return
	}
	return
}

// GetUserThirdBindByMid get user third bind by mid.
func (d *Dao) GetUserThirdBindByMid(c context.Context, mid int64) (res []*model.UserThirdBind, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, _getUserThirdBindByMidSQL, mid); err != nil {
		log.Error("fail to get UserThirdBind, dao.userDB.Query(%s) error(%+v)", _getUserThirdBindByMidSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserThirdBind)
		if err = rows.Scan(&r.ID, &r.Mid, &r.OpenID, &r.PlatForm, &r.Token, &r.Expires, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// GetUserThirdBindByMidAndPlatform get user third bind by mid and platform.
func (d *Dao) GetUserThirdBindByMidAndPlatform(c context.Context, mid, platform int64) (res *model.UserThirdBind, err error) {
	res = &model.UserThirdBind{}
	if err = d.userDB.QueryRow(c, _getUserThirdBindByMidAndPlatformSQL, mid, platform).Scan(&res.ID, &res.Mid, &res.OpenID, &res.PlatForm, &res.Token, &res.Expires, &res.CTime, &res.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserSafeQuestion by mid(%d) platform(%d), dao.userDB.QueryRow(%s) error(%+v)", mid, platform, _getUserThirdBindByMidAndPlatformSQL, err)
		}
		return
	}
	return
}

// GetUserRegOriginByMid get user reg origin by mid.
func (d *Dao) GetUserRegOriginByMid(c context.Context, mid int64) (res *model.UserRegOrigin, err error) {
	res = &model.UserRegOrigin{}
	if err = d.userDB.QueryRow(c, fmt.Sprintf(_getUserRegOriginByMidSQL, tableIndex(mid)), mid).Scan(&res.Mid, &res.JoinIP, &res.JoinIPV6, &res.Port, &res.JoinTime, &res.Origin, &res.RegType, &res.AppID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("fail to get UserRegOrigin by mid(%d), dao.userDB.QueryRow(%s) error(%+v)", mid, _getUserRegOriginByMidSQL, err)
		}
		return
	}
	return
}

// InsertUpdateUserBase insert update user base.
func (d *Dao) InsertUpdateUserBase(c context.Context, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _insertUpdateUserBaseSQL, a.Mid, a.UserID, a.Pwd, a.Salt, a.Status, a.Deleted, a.MTime, a.UserID, a.Pwd, a.Salt, a.Status); err != nil {
		log.Error("fail to insert update user base, userBase(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertIgnoreUserSafeQuestion insert ignore user safe question.
func (d *Dao) TxInsertIgnoreUserSafeQuestion(tx *xsql.Tx, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_insertIgnoreUserSafeQuestionSQL, tableIndex(a.Mid)), a.Mid, a.SafeQuestion, a.SafeAnswer, a.SafeBindTime); err != nil {
		log.Error("fail to insert ignore user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertIgnoreUserRegOrigin insert ignore user reg origin.
func (d *Dao) TxInsertIgnoreUserRegOrigin(tx *xsql.Tx, a *model.UserRegOrigin) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(fmt.Sprintf(_insertIgnoreUserRegOriginSQL, tableIndex(a.Mid)), a.Mid, a.JoinIP, a.JoinIPV6, a.Port, a.JoinTime); err != nil {
		log.Error("fail to insert ignore user reg origin, userRegOrigin(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// DelUserThirdBind del user third bind.
func (d *Dao) DelUserThirdBind(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _delUserThirdBindSQL, mid); err != nil {
		log.Error("fail to del user third bind, mid(%d) dao.userDB.Exec() error(%+v)", mid, err)
		return
	}
	return res.RowsAffected()
}

func tableIndex(mid int64) int64 {
	return mid % 100
}
