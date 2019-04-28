package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/passport-user-compare/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

var (
	insertUserBaseSQL = "INSERT INTO user_base (mid,userid,pwd,salt,status,deleted,mtime) VALUES(?,?,?,?,?,?,?)"
	selectUserBaseSQL = "SELECT mid, userid, pwd, salt, status FROM user_base WHERE mid = ?"
	updateUserBaseSQL = "UPDATE user_base SET userid = ?, pwd = ? ,salt = ? ,status = ? WHERE mid = ?"

	insertUserTelSQL = "INSERT  INTO user_tel (mid,tel,cid,tel_bind_time,mtime) VALUES(?,?,?,?,?)"
	selectUserTelSQL = "SELECT mid, tel, cid, tel_bind_time FROM user_tel WHERE mid = ?"
	updateUserTelSQL = "UPDATE user_tel SET tel = ? ,cid = ? WHERE mid = ?"

	insertUserMailSQL         = "INSERT INTO user_email (mid,email,verified,email_bind_time,mtime) VALUES(?,?,?,?,?)"
	selectUserMailSQL         = "SELECT mid,email,email_bind_time FROM user_email WHERE mid = ?"
	updateUserMailSQL         = "UPDATE user_email SET email = ?  WHERE mid = ?"
	updateUserMailVerifiedSQL = "UPDATE user_email SET verified = ?  WHERE mid = ?"

	insertUserSafeQuestionSQL = "INSERT INTO user_safe_question%02d (mid,safe_question,safe_answer,safe_bind_time) VALUES(?,?,?,?)"
	selectUserSafeQuestionSQL = "SELECT mid, safe_question, safe_answer FROM user_safe_question%02d WHERE mid = ?"
	updateUserSafeQuestionSQL = "UPDATE user_safe_question%02d SET safe_question = ? ,safe_answer = ? WHERE mid = ? "

	insertThirdBindSQL = "INSERT INTO user_third_bind (mid,openid,platform,token,expires) VALUES(?,?,?,?,?)"
	selectThirdBindSQL = "SELECT mid, openid, platform, token FROM user_third_bind WHERE mid = ? AND platform = ?"
	updateThirdBindSQL = "UPDATE user_third_bind SET openid = ? ,token = ? WHERE mid = ? AND platform = ? "

	queryCountryCodeSQL = "SELECT id,code FROM country_code"

	queryMidByTelSQL   = "SELECT mid FROM user_tel WHERE tel = ? and cid = ?"
	queryMidByEmailSQL = "SELECT mid FROM user_email WHERE email = ?"

	getUnverifiedEmail = "SELECT mid FROM user_email where mid > ? and verified = 0 limit 20000"

	_getUserRegOriginByMidSQL         = "SELECT mid,join_ip,join_time,origin,reg_type,appid from user_reg_origin%02d where mid = ?"
	_insertUpdateUserRegOriginTypeSQL = "INSERT INTO user_reg_origin%02d (mid,join_ip,join_time,origin,reg_type,appid,ctime,mtime) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE join_ip=?,join_time=?,origin=?,reg_type=?,appid=?,ctime=?,mtime=?"
)

// QueryUserBase query user basic info by mid
func (d *Dao) QueryUserBase(c context.Context, mid int64) (res *model.UserBase, err error) {
	row := d.userDB.QueryRow(c, selectUserBaseSQL, mid)
	res = new(model.UserBase)
	if err = row.Scan(&res.Mid, &res.UserID, &res.Pwd, &res.Salt, &res.Status); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// UpdateUserBase update user base
func (d *Dao) UpdateUserBase(c context.Context, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, updateUserBaseSQL, a.UserID, a.Pwd, a.Salt, a.Status, a.Mid); err != nil {
		log.Error("failed to update user base, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUserBase add user base.
func (d *Dao) InsertUserBase(c context.Context, a *model.UserBase) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, insertUserBaseSQL, a.Mid, a.UserID, a.Pwd, a.Salt, a.Status, a.Deleted, a.MTime); err != nil {
		log.Error("fail to add user base, userBase(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// QueryUserTel query user tel info by mid
func (d *Dao) QueryUserTel(c context.Context, mid int64) (res *model.UserTel, err error) {
	row := d.userDB.QueryRow(c, selectUserTelSQL, mid)
	res = new(model.UserTel)
	var cidPtr *string
	if err = row.Scan(&res.Mid, &res.Tel, &cidPtr, &res.TelBindTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	if cidPtr != nil {
		res.Cid = *cidPtr
	}
	return
}

// UpdateUserTel update user tel
func (d *Dao) UpdateUserTel(c context.Context, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = a.Cid
	}
	if res, err = d.userDB.Exec(c, updateUserTelSQL, telPtr, cidPtr, a.Mid); err != nil {
		log.Error("failed to update user tel, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUserTel insert user tel
func (d *Dao) InsertUserTel(c context.Context, a *model.UserTel) (affected int64, err error) {
	var (
		res    sql.Result
		telPtr *[]byte
		cidPtr string
	)
	if len(a.Tel) != 0 {
		telPtr = &a.Tel
		cidPtr = a.Cid
	}
	if res, err = d.userDB.Exec(c, insertUserTelSQL, a.Mid, telPtr, cidPtr, a.TelBindTime, a.MTime); err != nil {
		log.Error("fail to add user tel, userTel(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// QueryUserMail query user mail info by mid
func (d *Dao) QueryUserMail(c context.Context, mid int64) (res *model.UserEmail, err error) {
	row := d.userDB.QueryRow(c, selectUserMailSQL, mid)
	res = new(model.UserEmail)
	if err = row.Scan(&res.Mid, &res.Email, &res.EmailBindTime); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// UpdateUserMail update user tel
func (d *Dao) UpdateUserMail(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, updateUserMailSQL, a.Email, a.Mid); err != nil {
		log.Error("failed to update user mail, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserMailVerified update user email verified
func (d *Dao) UpdateUserMailVerified(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, updateUserMailVerifiedSQL, a.Verified, a.Mid); err != nil {
		log.Error("failed to update user mail verified, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUserEmail add user email.
func (d *Dao) InsertUserEmail(c context.Context, a *model.UserEmail) (affected int64, err error) {
	var (
		res      sql.Result
		emailPtr *[]byte
	)
	if len(a.Email) != 0 {
		emailPtr = &a.Email
	}
	if res, err = d.userDB.Exec(c, insertUserMailSQL, a.Mid, emailPtr, a.Verified, a.EmailBindTime, a.MTime); err != nil {
		log.Error("fail to add user email, userEmail(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// QueryUserSafeQuestion query user safe question by mid
func (d *Dao) QueryUserSafeQuestion(c context.Context, mid int64) (res *model.UserSafeQuestion, err error) {
	row := d.userDB.QueryRow(c, fmt.Sprintf(selectUserSafeQuestionSQL, mid%safeQuestionSegment), mid)
	res = new(model.UserSafeQuestion)
	if err = row.Scan(&res.Mid, &res.SafeQuestion, &res.SafeAnswer); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// UpdateUserSafeQuestion update user tel
func (d *Dao) UpdateUserSafeQuestion(c context.Context, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(updateUserSafeQuestionSQL, a.Mid%safeQuestionSegment), a.SafeQuestion, a.SafeAnswer, a.Mid); err != nil {
		log.Error("failed to update user safe question, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUserSafeQuestion insert user safe question
func (d *Dao) InsertUserSafeQuestion(c context.Context, a *model.UserSafeQuestion) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(insertUserSafeQuestionSQL, a.Mid%safeQuestionSegment), a.Mid, a.SafeQuestion, a.SafeAnswer, a.SafeBindTime); err != nil {
		log.Error("fail to add user safe question, userSafeQuestion(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// QueryUserThirdBind query user third bind by mid and platform
func (d *Dao) QueryUserThirdBind(c context.Context, mid, platform int64) (res *model.UserThirdBind, err error) {
	row := d.userDB.QueryRow(c, selectThirdBindSQL, mid, platform)
	res = new(model.UserThirdBind)
	if err = row.Scan(&res.Mid, &res.OpenID, &res.PlatForm, &res.Token); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// UpdateUserThirdBind update user third bind
func (d *Dao) UpdateUserThirdBind(c context.Context, a *model.UserThirdBind) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, updateThirdBindSQL, a.OpenID, a.Token, a.Mid, a.PlatForm); err != nil {
		log.Error("failed to update user third bind sql, dao.userDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUserThirdBind insert user third bind.
func (d *Dao) InsertUserThirdBind(c context.Context, a *model.UserThirdBind) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, insertThirdBindSQL, a.Mid, a.OpenID, a.PlatForm, a.Token, a.Expires); err != nil {
		log.Error("fail to add user third bind, userThirdBind(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// QueryCountryCode query country code
func (d *Dao) QueryCountryCode(c context.Context) (res map[int64]string, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, queryCountryCodeSQL); err != nil {
		log.Error("fail to get CountryCodeMap, dao.originDB.Query(%s) error(%v)", queryCountryCodeSQL, err)
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
	if err = d.userDB.QueryRow(c, queryMidByTelSQL, a.Tel, a.Cid).Scan(&mid); err != nil {
		log.Error("fail to get mid by tel, dao.userDB.QueryRow(%s) error(%+v)", queryMidByTelSQL, err)
		return
	}
	return
}

// GetMidByEmail get mid by email.
func (d *Dao) GetMidByEmail(c context.Context, a *model.UserEmail) (mid int64, err error) {
	if err = d.userDB.QueryRow(c, queryMidByEmailSQL, a.Email).Scan(&mid); err != nil {
		log.Error("fail to get mid by email, dao.userDB.QueryRow(%s) error(%+v)", queryMidByEmailSQL, err)
		return
	}
	return
}

// GetUnverifiedEmail get unverified email.
func (d *Dao) GetUnverifiedEmail(c context.Context, start int64) (res []*model.UserEmail, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, getUnverifiedEmail, start); err != nil {
		log.Error("fail to get UnverifiedEmail, dao.userDB.Query(%s) error(%v)", getUnverifiedEmail, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserEmail)
		if err = rows.Scan(&r.Mid); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// GetUserRegOriginByMid get user reg origin by mid.
func (d *Dao) GetUserRegOriginByMid(c context.Context, mid int64) (res *model.UserRegOrigin, err error) {
	res = &model.UserRegOrigin{}
	if err = d.userDB.QueryRow(c, fmt.Sprintf(_getUserRegOriginByMidSQL, tableIndex(mid)), mid).Scan(&res.Mid, &res.JoinIP, &res.JoinTime, &res.Origin, &res.RegType, &res.AppID); err != nil {
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

// InsertUpdateUserRegOriginType insert update user reg origin type.
func (d *Dao) InsertUpdateUserRegOriginType(c context.Context, a *model.UserRegOrigin) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, fmt.Sprintf(_insertUpdateUserRegOriginTypeSQL, tableIndex(a.Mid)), a.Mid, a.JoinIP, a.JoinTime, a.Origin, a.RegType, a.AppID, a.CTime, a.MTime,
		a.JoinIP, a.JoinTime, a.Origin, a.RegType, a.AppID, a.CTime, a.MTime); err != nil {
		log.Error("fail to insert update user reg origin type, userRegOrigin(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

func tableIndex(mid int64) int64 {
	return mid % 100
}
