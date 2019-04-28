package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/job/main/passport/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_insertLoginLog    = "INSERT INTO aso_login_log%d(`mid`, `timestamp`, `loginip`, `type`, `server`) VALUES %s"
	_queryTelBindLog   = "SELECT id, mid, tel, timestamp FROM aso_telephone_bind_log where id = ?"
	_queryEmailBindLog = "SELECT id, mid, email, timestamp FROM aso_email_bind_log where id = ?"
	_batchGetPwdLog    = "select id, timestamp, mid, ip, old_pwd, old_salt, new_pwd, new_salt from aso_pwd_log where id < ? order by id desc limit 1000"
	_getPwdLog         = "select id, timestamp, mid, ip, old_pwd, old_salt, new_pwd, new_salt from aso_pwd_log where id = ?"
)

// AddLoginLog insert service to db.
func (d *Dao) AddLoginLog(vs []*model.LoginLog) (err error) {
	if len(vs) == 0 {
		return
	}
	var args = make([]string, 0, len(vs))
	for _, v := range vs {
		args = append(args, fmt.Sprintf(`(%d,%d,%d,%d,'%s')`, v.Mid, v.Timestamp, v.LoginIP, v.Type, v.Server))
	}
	if len(args) == 0 {
		return
	}
	s := fmt.Sprintf(_insertLoginLog, vs[0].Mid%10, strings.Join(args, ","))
	if _, err = d.logDB.Exec(context.Background(), s); err != nil {
		log.Error("d.logDB.Exec(%s) error(%v)", s, err)
	}
	return
}

// QueryTelBindLog query from id
func (d *Dao) QueryTelBindLog(id int64) (res *model.TelBindLog, err error) {
	if id <= 0 {
		return
	}
	res = new(model.TelBindLog)
	row := d.asoDB.QueryRow(context.Background(), _queryTelBindLog, id)
	if err = row.Scan(&res.ID, &res.Mid, &res.Tel, &res.Timestamp); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("QueryTelBindLog err(%+v)", err)
		return
	}
	return
}

// QueryEmailBindLog query from id
func (d *Dao) QueryEmailBindLog(id int64) (res *model.EmailBindLog, err error) {
	if id <= 0 {
		return
	}
	res = new(model.EmailBindLog)
	row := d.asoDB.QueryRow(context.Background(), _queryEmailBindLog, id)
	if err = row.Scan(&res.ID, &res.Mid, &res.Email, &res.Timestamp); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("QueryEmailBindLog err(%+v)", err)
		return
	}
	return
}

// BatchGetPwdLog batch get pwd log
func (d *Dao) BatchGetPwdLog(c context.Context, id int64) (res []*model.PwdLog, err error) {
	var rows *sql.Rows
	if rows, err = d.asoDB.Query(c, _batchGetPwdLog, id); err != nil {
		log.Error("batch get pwd log, dao.db.Query(%s) error(%v)", _batchGetPwdLog, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pwd := new(model.PwdLog)
		if err = rows.Scan(&pwd.ID, &pwd.Timestamp, &pwd.Mid, &pwd.IP, &pwd.OldPwd, &pwd.OldSalt, &pwd.NewPwd, &pwd.NewSalt); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, pwd)
	}
	return
}

// GetPwdLog  get pwd log
func (d *Dao) GetPwdLog(c context.Context, id int64) (res *model.PwdLog, err error) {
	res = new(model.PwdLog)
	row := d.asoDB.QueryRow(c, _getPwdLog, id)
	if err = row.Scan(&res.ID, &res.Timestamp, &res.Mid, &res.IP, &res.OldPwd, &res.OldSalt, &res.NewPwd, &res.NewSalt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
		return
	}
	return
}
