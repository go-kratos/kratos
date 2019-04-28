package dao

import (
	"context"

	"go-common/app/interface/main/passport-login/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	getUserBaseByMidSQL    = "SELECT mid,userid,pwd,salt,status FROM user_base WHERE deleted = 0 AND mid = ?"
	getUserBaseByUserIDSQL = "SELECT mid from user_base where userid= ?"
	getUserTelByTelSQL     = "SELECT mid from user_tel where tel= ?"
	getUserMailByMailSQL   = "SELECT mid FROM user_email where email = ?"
)

// GetUserByMid get user by userID
func (d *Dao) GetUserByMid(c context.Context, mid int64) (res *model.User, err error) {
	res = new(model.User)
	if err = d.userDB.QueryRow(c, getUserBaseByMidSQL, mid).Scan(&res.Mid, &res.UserID, &res.Pwd, &res.Salt, &res.Status); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("fail to get user by userid, dao.userDB.QueryRow(%s) error(%v)", getUserBaseByMidSQL, err)
		return
	}
	return
}

// GetUserBaseByUserID get user base by userID
func (d *Dao) GetUserBaseByUserID(c context.Context, userID string) (res int64, err error) {
	if err = d.userDB.QueryRow(c, getUserBaseByUserIDSQL, userID).Scan(&res); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = 0
			return
		}
		log.Error("fail to get userbase by userid, dao.userDB.QueryRow(%s) error(%v)", getUserBaseByMidSQL, err)
		return
	}
	return
}

// GetUserTelByTel get user tel by tel
func (d *Dao) GetUserTelByTel(c context.Context, tel []byte) (res int64, err error) {
	if err = d.userDB.QueryRow(c, getUserTelByTelSQL, tel).Scan(&res); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = 0
			return
		}
		log.Error("fail to get user tel by tel, dao.userDB.QueryRow(%s) error(%v)", getUserBaseByMidSQL, err)
		return
	}
	return
}

// GetUserMailByMail get user mail by mail
func (d *Dao) GetUserMailByMail(c context.Context, mail []byte) (res int64, err error) {
	if err = d.userDB.QueryRow(c, getUserMailByMailSQL, mail).Scan(&res); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = 0
			return
		}
		log.Error("fail to get user mail by mail, dao.userDB.QueryRow(%s) error(%v)", getUserBaseByMidSQL, err)
		return
	}
	return
}
