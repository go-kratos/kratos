package dao

import (
	"context"
	"database/sql"

	"go-common/app/service/main/passport-game/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_memberShard       = 30
	_getMemberInfoSQL  = "SELECT mid,face FROM member_%02d WHERE mid=?"
	_getAppsSQL        = "SELECT appid,appkey,app_secret FROM app"
	_addTokenSQL       = "INSERT INTO app_perm (mid,appid,app_subid,access_token,create_at,expires) VALUES (?,?,?,?,?,?)"
	_updateTokenSQL    = "UPDATE app_perm SET expires=? WHERE access_token=? AND expires<?"
	_getTokenSQL       = "SELECT mid,appid,app_subid,access_token,create_at,expires FROM app_perm WHERE access_token=?"
	_getAsoAccountSQL  = "SELECT mid,userid,uname,pwd,salt,isleak FROM aso_account WHERE userid=? OR email=? OR tel=?"
	_getAccountInfoSQL = "SELECT mid,userid,uname,email,tel FROM aso_account WHERE mid=?"
)

func hit(mid int64) int64 {
	return mid % _memberShard
}

// MemberInfo get member info.
func (d *Dao) MemberInfo(c context.Context, mid int64) (res *model.Info, err error) {
	var row = d.getMemberStmt[hit(mid)].QueryRow(c, mid)
	res = new(model.Info)
	if err = row.Scan(&res.Mid, &res.Face); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// Apps get all apps.
func (d *Dao) Apps(c context.Context) (res []*model.App, err error) {
	var rows *xsql.Rows
	if rows, err = d.cloudDB.Query(c, _getAppsSQL); err != nil {
		log.Error("get apps, dao.cloudDB.Query(%s) error(%v)", _getAppsSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		app := new(model.App)
		if err = rows.Scan(&app.AppID, &app.AppKey, &app.AppSecret); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, app)
	}
	err = rows.Err()
	return
}

// AddToken add token.
func (d *Dao) AddToken(c context.Context, t *model.Perm) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _addTokenSQL, t.Mid, t.AppID, t.AppSubID, t.AccessToken, t.CreateAt, t.Expires); err != nil {
		log.Error("add token, dao.cloudDB.Exec(%d, %d, %d, %s, %d, %d) error(%v)", t.Mid, t.AppID, t.AppSubID, t.AccessToken, t.CreateAt, t.Expires, err)
		return
	}
	return res.RowsAffected()
}

// UpdateToken update token.
func (d *Dao) UpdateToken(c context.Context, t *model.Perm) (affected int64, err error) {
	var res sql.Result
	if res, err = d.cloudDB.Exec(c, _updateTokenSQL, t.Expires, t.AccessToken, t.Expires); err != nil {
		log.Error("update token, dao.cloudDB.Exec(%d, %s, %d) error(%v)", t.Expires, t.AccessToken, t.Expires, err)
		return
	}
	return res.RowsAffected()
}

// Token get token.
func (d *Dao) Token(c context.Context, accessToken string) (res *model.Perm, err error) {
	row := d.cloudDB.QueryRow(c, _getTokenSQL, accessToken)
	res = new(model.Perm)
	if err = row.Scan(&res.Mid, &res.AppID, &res.AppSubID, &res.AccessToken, &res.CreateAt, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// AsoAccount get aso account.
func (d *Dao) AsoAccount(c context.Context, identify, identifyHash string) (res []*model.AsoAccount, err error) {
	var rows *xsql.Rows
	if rows, err = d.cloudDB.Query(c, _getAsoAccountSQL, identify, identifyHash, identifyHash); err != nil {
		log.Error("get apps, dao.cloudDB.Query(%s) error(%v)", _getAsoAccountSQL, err)
		return
	}
	for rows.Next() {
		aso := new(model.AsoAccount)
		if err = rows.Scan(&aso.Mid, &aso.UserID, &aso.Uname, &aso.Pwd, &aso.Salt, &aso.Isleak); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, aso)
	}
	err = rows.Err()
	return
}

// AccountInfo get account info.
func (d *Dao) AccountInfo(c context.Context, mid int64) (res *model.AsoAccount, err error) {
	row := d.cloudDB.QueryRow(c, _getAccountInfoSQL, mid)
	res = new(model.AsoAccount)
	if err = row.Scan(&res.Mid, &res.UserID, &res.Uname, &res.Email, &res.Tel); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// TokenFromOtherRegion get token from otherRegion.
func (d *Dao) TokenFromOtherRegion(c context.Context, accessToken string) (res *model.Perm, err error) {
	row := d.otherRegion.QueryRow(c, _getTokenSQL, accessToken)
	res = new(model.Perm)
	if err = row.Scan(&res.Mid, &res.AppID, &res.AppSubID, &res.AccessToken, &res.CreateAt, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
		return
	}
	return
}
