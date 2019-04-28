package dao

import (
	"context"
	"database/sql"

	"go-common/app/job/main/passport-sns/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addSnsLogSQL = "INSERT INTO sns_log (mid,openid,unionid,appid,platform,operator,operate,description) VALUES(?,?,?,?,?,?,?,?)"

	_getAsoAccountSnsAllSQL = "SELECT mid,sina_uid,sina_access_token,sina_access_expires,qq_openid,qq_access_token,qq_access_expires FROM aso_account_sns WHERE mid > ? order by mid limit 20000"
	_getAsoAccountSnsSQL    = "SELECT mid,sina_uid,sina_access_token,sina_access_expires,qq_openid,qq_access_token,qq_access_expires FROM aso_account_sns WHERE (qq_openid != '' or sina_uid != 0) and mid > ? order by mid limit 20000"
)

// AddSnsLog add sns log.
func (d *Dao) AddSnsLog(c context.Context, a *model.SnsLog) (affected int64, err error) {
	var res sql.Result
	if res, err = d.snsDB.Exec(c, _addSnsLogSQL, a.Mid, a.OpenID, a.UnionID, a.AppID, a.Platform, a.Operator, a.Operate, a.Description); err != nil {
		log.Error("AddSnsLog(%+v) tx.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AsoAccountSnsAll get account sns
func (d *Dao) AsoAccountSnsAll(c context.Context, start int64) (res []*model.AsoAccountSns, err error) {
	var rows *xsql.Rows
	if rows, err = d.asoDB.Query(c, _getAsoAccountSnsAllSQL, start); err != nil {
		log.Error("fail to get AsoAccountSns, dao.asoDB.Query() error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AsoAccountSns)
		if err = rows.Scan(&r.Mid, &r.SinaUID, &r.SinaAccessToken, &r.SinaAccessExpires, &r.QQOpenid, &r.QQAccessToken, &r.QQAccessExpires); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// AsoAccountSns get account sns by id.
func (d *Dao) AsoAccountSns(c context.Context, start int64) (res []*model.AsoAccountSns, err error) {
	var rows *xsql.Rows
	if rows, err = d.asoDB.Query(c, _getAsoAccountSnsSQL, start); err != nil {
		log.Error("fail to get AsoAccountSns, dao.asoDB.Query() error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AsoAccountSns)
		if err = rows.Scan(&r.Mid, &r.SinaUID, &r.SinaAccessToken, &r.SinaAccessExpires, &r.QQOpenid, &r.QQAccessToken, &r.QQAccessExpires); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}
