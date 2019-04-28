package dao

import (
	"context"
	rsql "database/sql"
	"fmt"
	"strings"

	"go-common/app/service/main/account-recovery/dao/sqlbuilder"
	"go-common/app/service/main/account-recovery/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_selectCountRecoveryInfo = "select count(rid) from account_recovery_info"
	_selectRecoveryInfoLimit = "select rid,mid,user_type,status,login_addrs,unames,reg_time,reg_type,reg_addr,pwds,phones,emails,safe_question,safe_answer,card_type,card_id," +
		"sys_login_addrs,sys_reg,sys_unames,sys_pwds,sys_phones,sys_emails,sys_safe,sys_card," +
		"link_email,operator,opt_time,remark,ctime,business from account_recovery_info %s"
	_getSuccessCount         = "SELECT count FROM account_recovery_success WHERE mid=?"
	_batchGetRecoverySuccess = "SELECT mid,count,ctime,mtime FROM account_recovery_success WHERE mid in (%s)"
	_updateSuccessCount      = "INSERT INTO account_recovery_success (mid, count) VALUES (?, 1) ON DUPLICATE KEY UPDATE count = count + 1"
	_batchUpdateSuccessCount = "INSERT INTO account_recovery_success (mid, count) VALUES %s ON DUPLICATE KEY UPDATE count = count + 1"
	_updateStatus            = "UPDATE account_recovery_info SET status=?,operator=?,opt_time=?,remark=? WHERE rid = ? AND `status`=0"
	_getNoDeal               = "SELECT COUNT(1) FROM account_recovery_info WHERE mid=? AND `status`=0"
	_updateUserType          = "UPDATE account_recovery_info SET user_type=? WHERE rid = ?"
	_insertRecoveryInfo      = "INSERT INTO account_recovery_info(login_addrs,unames,reg_time,reg_type,reg_addr,pwds,phones,emails,safe_question,safe_answer,card_type,card_id,link_email,mid,business,last_suc_count,last_suc_ctime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updateSysInfo           = "UPDATE account_recovery_info SET sys_login_addrs=?,sys_reg=?,sys_unames=?,sys_pwds=?,sys_phones=?,sys_emails=?,sys_safe=?,sys_card=?,user_type=? WHERE rid=?"
	_getUinfoByRid           = "SELECT mid,link_email,ctime FROM account_recovery_info WHERE rid=? LIMIT 1"
	_getUinfoByRidMore       = "SELECT rid,mid,link_email,ctime FROM account_recovery_info WHERE rid in (%s)"
	_selectUnCheckInfo       = "SELECT mid,login_addrs,unames,reg_time,reg_type,reg_addr,pwds,phones,emails,safe_question,safe_answer,card_type,card_id FROM account_recovery_info WHERE rid=? AND `status`=0 AND sys_card=''"
	_getStatusByRid          = "SELECT `status` FROM account_recovery_info WHERE rid=?"
	_getMailStatus           = "SELECT mail_status FROM account_recovery_info WHERE rid=?"
	_updateMailStatus        = "UPDATE account_recovery_info SET mail_status=1 WHERE rid=?"
	_insertRecoveryAddit     = "INSERT INTO account_recovery_addit(`rid`, `files`, `extra`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE files=VALUES(files),extra=VALUES(extra)"
	_updateRecoveryAddit     = "UPDATE account_recovery_addit SET `files` = ?,`extra` = ? WHERE rid = ?"
	_getRecoveryAddit        = "SELECT rid, `files`,`extra`, ctime, mtime FROM account_recovery_addit WHERE rid= ?"
	_batchRecoveryAAdit      = "SELECT rid, `files`, `extra`, ctime, mtime FROM account_recovery_addit WHERE rid in (%s)"
	_batchGetLastSuccess     = "SELECT mid,max(ctime) FROM account_recovery_info WHERE mid in (%s)  AND `status`=1 GROUP BY mid"
	_getLastSuccess          = "SELECT mid,max(ctime) FROM account_recovery_info WHERE mid = ? AND `status`=1"
)

// GetStatusByRid get status by rid
func (dao *Dao) GetStatusByRid(c context.Context, rid int64) (status int64, err error) {
	res := dao.db.Prepared(_getStatusByRid).QueryRow(c, rid)
	if err = res.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			status = -1
			err = nil
		} else {
			log.Error("GetStatusByRid row.Scan error(%v)", err)
		}
	}
	return
}

// GetSuccessCount get success count
func (dao *Dao) GetSuccessCount(c context.Context, mid int64) (count int64, err error) {
	res := dao.db.Prepared(_getSuccessCount).QueryRow(c, mid)
	if err = res.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			count = 0
			err = nil
		} else {
			log.Error("GetSuccessCount row.Scan error(%v)", err)
		}
	}
	return
}

// BatchGetRecoverySuccess batch get recovery success info
func (dao *Dao) BatchGetRecoverySuccess(c context.Context, mids []int64) (successMap map[int64]*model.RecoverySuccess, err error) {
	rows, err := dao.db.Query(c, fmt.Sprintf(_batchGetRecoverySuccess, xstr.JoinInts(mids)))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("BatchGetRecoverySuccess d.db.Query error(%v)", err)
		return
	}
	successMap = make(map[int64]*model.RecoverySuccess)
	for rows.Next() {
		r := new(model.RecoverySuccess)
		if err = rows.Scan(&r.SuccessMID, &r.SuccessCount, &r.FirstSuccessTime, &r.LastSuccessTime); err != nil {
			log.Error("BatchGetRecoverySuccess rows.Scan error(%v)", err)
			continue
		}
		successMap[r.SuccessMID] = r
	}
	return
}

// UpdateSuccessCount insert or update success count
func (dao *Dao) UpdateSuccessCount(c context.Context, mid int64) (err error) {
	_, err = dao.db.Exec(c, _updateSuccessCount, mid)
	return
}

// BatchUpdateSuccessCount batch insert or update success count
func (dao *Dao) BatchUpdateSuccessCount(c context.Context, mids string) (err error) {
	var s string
	midArr := strings.Split(mids, ",")
	for _, mid := range midArr {
		s = s + fmt.Sprintf(",(%s, 1)", mid)
	}
	_, err = dao.db.Exec(c, fmt.Sprintf(_batchUpdateSuccessCount, s[1:]))
	return
}

// GetNoDeal get no deal record
func (dao *Dao) GetNoDeal(c context.Context, mid int64) (count int64, err error) {
	res := dao.db.Prepared(_getNoDeal).QueryRow(c, mid)
	if err = res.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetNoDeal row.Scan error(%v)", err)
		return
	}
	return
}

// UpdateStatus update field status.
func (dao *Dao) UpdateStatus(c context.Context, status int64, rid int64, operator string, optTime xtime.Time, remark string) (err error) {
	_, err = dao.db.Exec(c, _updateStatus, status, operator, optTime, remark, rid)
	return
}

// UpdateUserType update field user_type.
func (dao *Dao) UpdateUserType(c context.Context, status int64, rid int64) (err error) {
	if _, err = dao.db.Exec(c, _updateUserType, status, rid); err != nil {
		log.Error("dao.db.Exec(%s, %d, %d) error(%v)", _updateUserType, status, rid, err)
	}
	return
}

// InsertRecoveryInfo insert data
func (dao *Dao) InsertRecoveryInfo(c context.Context, uinfo *model.UserInfoReq) (lastID int64, err error) {
	var res rsql.Result
	if res, err = dao.db.Exec(c, _insertRecoveryInfo, uinfo.LoginAddrs, uinfo.Unames, uinfo.RegTime, uinfo.RegType, uinfo.RegAddr,
		uinfo.Pwds, uinfo.Phones, uinfo.Emails, uinfo.SafeQuestion, uinfo.SafeAnswer, uinfo.CardType, uinfo.CardID, uinfo.LinkMail, uinfo.Mid, uinfo.Business, uinfo.LastSucCount, uinfo.LastSucCTime); err != nil {
		log.Error("dao.db.Exec(%s, %v) error(%v)", _insertRecoveryInfo, uinfo, err)
		return
	}
	return res.LastInsertId()
}

// UpdateSysInfo update sysinfo and user_type
func (dao *Dao) UpdateSysInfo(c context.Context, sys *model.SysInfo, userType int64, rid int64) (err error) {
	if _, err = dao.db.Exec(c, _updateSysInfo, &sys.SysLoginAddrs, &sys.SysReg, &sys.SysUNames, &sys.SysPwds, &sys.SysPhones,
		&sys.SysEmails, &sys.SysSafe, &sys.SysCard, userType, rid); err != nil {
		log.Error("dao.db.Exec(%s, %v) error(%v)", _updateSysInfo, sys, err)
	}
	return
}

// GetAllByCon get a pageData by more condition
func (dao *Dao) GetAllByCon(c context.Context, aq *model.QueryRecoveryInfoReq) ([]*model.AccountRecoveryInfo, int64, error) {
	query := sqlbuilder.NewSelectBuilder().Select("rid,mid,user_type,status,login_addrs,unames,reg_time,reg_type,reg_addr,pwds,phones,emails,safe_question,safe_answer,card_type,card_id,sys_login_addrs,sys_reg,sys_unames,sys_pwds,sys_phones,sys_emails,sys_safe,sys_card,link_email,operator,opt_time,remark,ctime,business,last_suc_count,last_suc_ctime").From("account_recovery_info")

	if aq.Bussiness != "" {
		query = query.Where(query.Equal("business", aq.Bussiness))
	}

	if aq.Status != nil {
		query = query.Where(fmt.Sprintf("status=%d", *aq.Status))
	}
	if aq.Game != nil {
		query = query.Where(fmt.Sprintf("user_type=%d", *aq.Game))
	}
	if aq.UID != 0 {
		query = query.Where(fmt.Sprintf("mid=%d", aq.UID))
	}
	if aq.RID != 0 {
		query = query.Where(fmt.Sprintf("rid=%d", aq.RID))
	}
	if aq.StartTime != 0 {
		query = query.Where(query.GE("ctime", aq.StartTime.Time()))
	}
	if aq.EndTime != 0 {
		query = query.Where(query.LE("ctime", aq.EndTime.Time()))
	}
	totalSQL, totalArg := query.Copy().Select("count(1)").Build()
	log.Info("Build GetAllByCon total count SQL: %s", totalSQL)
	page := aq.Page
	if page == 0 {
		page = 1
	}
	size := aq.Size
	if size == 0 {
		size = 50
	}
	query = query.Limit(int(size)).Offset(int(size * (page - 1))).OrderBy("rid DESC")
	rawSQL, rawArg := query.Build()
	log.Info("Build GetAllByCon SQL: %s", rawSQL)

	total := int64(0)
	row := dao.db.QueryRow(c, totalSQL, totalArg...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := dao.db.Query(c, rawSQL, rawArg...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	resultData := make([]*model.AccountRecoveryInfo, 0)
	for rows.Next() {
		r := new(model.AccountRecoveryInfo)
		if err = rows.Scan(&r.Rid, &r.Mid, &r.UserType, &r.Status, &r.LoginAddr, &r.UNames, &r.RegTime, &r.RegType, &r.RegAddr,
			&r.Pwd, &r.Phones, &r.Emails, &r.SafeQuestion, &r.SafeAnswer, &r.CardType, &r.CardID,
			&r.SysLoginAddr, &r.SysReg, &r.SysUNames, &r.SysPwds, &r.SysPhones, &r.SysEmails, &r.SysSafe, &r.SysCard,
			&r.LinkEmail, &r.Operator, &r.OptTime, &r.Remark, &r.CTime, &r.Bussiness, &r.LastSucCount, &r.LastSucCTime); err != nil {
			log.Error("GetAllByCon error (%+v)", err)
			continue
		}
		resultData = append(resultData, r)
	}
	return resultData, total, err
}

// QueryByID query by rid
func (dao *Dao) QueryByID(c context.Context, rid int64, fromTime, endTime xtime.Time) (res *model.AccountRecoveryInfo, err error) {
	sql1 := "select rid,mid,user_type,status,login_addrs,unames,reg_time,reg_type,reg_addr,pwds,phones,emails,safe_question,safe_answer,card_type,card_id," +
		"sys_login_addrs,sys_reg,sys_unames,sys_pwds,sys_phones,sys_emails,sys_safe,sys_card," +
		"link_email,operator,opt_time,remark,ctime,business from account_recovery_info where ctime between ? and ? and rid = ?"
	res = new(model.AccountRecoveryInfo)
	row := dao.db.QueryRow(c, sql1, fromTime, endTime, rid)
	if err = row.Scan(&res.Rid, &res.Mid, &res.UserType, &res.Status, &res.LoginAddr, &res.UNames, &res.RegTime, &res.RegType, &res.RegAddr,
		&res.Pwd, &res.Phones, &res.Emails, &res.SafeQuestion, &res.SafeAnswer, &res.CardType, &res.CardID,
		&res.SysLoginAddr, &res.SysReg, &res.SysUNames, &res.SysPwds, &res.SysPhones, &res.SysEmails, &res.SysSafe, &res.SysCard,
		&res.LinkEmail, &res.Operator, &res.OptTime, &res.Remark, &res.CTime, &res.Bussiness); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("QueryByID(%d) error(%v)", rid, err)
		return
	}
	return
}

//QueryInfoByLimit  page query through limit m,n
func (dao *Dao) QueryInfoByLimit(c context.Context, req *model.DBRecoveryInfoParams) (res []*model.AccountRecoveryInfo, total int64, err error) {
	p := make([]interface{}, 0)
	s := " where ctime between ? and ?"
	p = append(p, req.StartTime)
	p = append(p, req.EndTime)
	if req.ExistGame {
		s = s + " and user_type = ?"
		p = append(p, req.Game)
	}
	if req.ExistStatus {
		s = s + " and status = ?"
		p = append(p, req.Status)
	}
	if req.ExistMid {
		s = s + " and mid = ?"
		p = append(p, req.Mid)
	}

	var s2 = s + " order by rid desc limit ?,?"
	p2 := p
	p2 = append(p2, (req.CurrPage-1)*req.Size, req.Size)
	var rows *sql.Rows
	rows, err = dao.db.Query(c, fmt.Sprintf(_selectRecoveryInfoLimit, s2), p2...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("QueryInfo err: d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.AccountRecoveryInfo, 0, req.Size)
	for rows.Next() {
		r := new(model.AccountRecoveryInfo)
		if err = rows.Scan(&r.Rid, &r.Mid, &r.UserType, &r.Status, &r.LoginAddr, &r.UNames, &r.RegTime, &r.RegType, &r.RegAddr,
			&r.Pwd, &r.Phones, &r.Emails, &r.SafeQuestion, &r.SafeAnswer, &r.CardType, &r.CardID,
			&r.SysLoginAddr, &r.SysReg, &r.SysUNames, &r.SysPwds, &r.SysPhones, &r.SysEmails, &r.SysSafe, &r.SysCard,
			&r.LinkEmail, &r.Operator, &r.OptTime, &r.Remark, &r.CTime, &r.Bussiness); err != nil {
			log.Error("QueryInfo (%v) error (%v)", req, err)
			return
		}
		res = append(res, r)
	}
	row := dao.db.QueryRow(c, _selectCountRecoveryInfo+s, p...)
	if err = row.Scan(&total); err != nil {
		log.Error("QueryInfo total error (%v)", err)
		return
	}
	return
}

// GetUinfoByRid get mid,linkMail by rid
func (dao *Dao) GetUinfoByRid(c context.Context, rid int64) (mid int64, linkMail string, ctime string, err error) {
	res := dao.db.Prepared(_getUinfoByRid).QueryRow(c, rid)
	req := new(struct {
		Mid      int64
		LinKMail string
		Ctime    xtime.Time
	})

	if err = res.Scan(&req.Mid, &req.LinKMail, &req.Ctime); err != nil {
		if err == sql.ErrNoRows {
			req.Mid = 0
			err = nil
		} else {
			log.Error("GetUinfoByRid row.Scan error(%v)", err)
		}
	}
	mid = req.Mid
	linkMail = req.LinKMail
	ctime = req.Ctime.Time().Format("2006-01-02 15:04:05")
	return
}

// GetUinfoByRidMore get list of BatchAppeal by rid
func (dao *Dao) GetUinfoByRidMore(c context.Context, ridsStr string) (bathRes []*model.BatchAppeal, err error) {

	rows, err := dao.db.Prepared(fmt.Sprintf(_getUinfoByRidMore, ridsStr)).Query(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bathRes = make([]*model.BatchAppeal, 0, len(strings.Split(ridsStr, ",")))
	for rows.Next() {
		req := &model.BatchAppeal{}
		if err = rows.Scan(&req.Rid, &req.Mid, &req.LinkMail, &req.Ctime); err != nil {
			return
		}
		bathRes = append(bathRes, req)
	}
	return
}

// GetUnCheckInfo get uncheck info
func (dao *Dao) GetUnCheckInfo(c context.Context, rid int64) (r *model.UserInfoReq, err error) {
	row := dao.db.QueryRow(c, _selectUnCheckInfo, rid)
	r = new(model.UserInfoReq)
	if err = row.Scan(&r.Mid, &r.LoginAddrs, &r.Unames, &r.RegTime, &r.RegType, &r.RegAddr,
		&r.Pwds, &r.Phones, &r.Emails, &r.SafeQuestion, &r.SafeAnswer, &r.CardType, &r.CardID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetUnCheckInfo (%v) error (%v)", rid, err)
	}
	return
}

//BeginTran begin transaction
func (dao *Dao) BeginTran(ctx context.Context) (tx *sql.Tx, err error) {
	if tx, err = dao.db.Begin(ctx); err != nil {
		log.Error("db: begintran BeginTran d.db.Begin error(%v)", err)
	}

	return
}

// GetMailStatus get mail_status by rid
func (dao *Dao) GetMailStatus(c context.Context, rid int64) (mailStatus int64, err error) {
	res := dao.db.Prepared(_getMailStatus).QueryRow(c, rid)
	if err = res.Scan(&mailStatus); err != nil {
		if err == sql.ErrNoRows {
			mailStatus = -1
			err = nil
		} else {
			log.Error("GetStatusByRid row.Scan error(%v)", err)
		}
	}
	return
}

// UpdateMailStatus update mail_status.
func (dao *Dao) UpdateMailStatus(c context.Context, rid int64) (err error) {
	_, err = dao.db.Exec(c, _updateMailStatus, rid)
	return
}

// UpdateRecoveryAddit is
func (dao *Dao) UpdateRecoveryAddit(c context.Context, rid int64, files []string, extra string) (err error) {
	_, err = dao.db.Exec(c, _updateRecoveryAddit, strings.Join(files, ","), extra, rid)
	return
}

// GetRecoveryAddit is
func (dao *Dao) GetRecoveryAddit(c context.Context, rid int64) (addit *model.DBAccountRecoveryAddit, err error) {
	row := dao.db.QueryRow(c, _getRecoveryAddit, rid)
	addit = new(model.DBAccountRecoveryAddit)
	if err = row.Scan(&addit.Rid, &addit.Files, &addit.Extra, &addit.Ctime, &addit.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetRecoveryAddit (%v) error (%v)", rid, err)
	}
	return
}

// InsertRecoveryAddit is
func (dao *Dao) InsertRecoveryAddit(c context.Context, rid int64, files, extra string) (err error) {
	_, err = dao.db.Exec(c, _insertRecoveryAddit, rid, files, extra)
	return
}

//BatchGetRecoveryAddit is
func (dao *Dao) BatchGetRecoveryAddit(c context.Context, rids []int64) (addits map[int64]*model.DBAccountRecoveryAddit, err error) {

	rows, err := dao.db.Query(c, fmt.Sprintf(_batchRecoveryAAdit, xstr.JoinInts(rids)))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("BatchGetRecoveryAddit d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	addits = make(map[int64]*model.DBAccountRecoveryAddit)
	for rows.Next() {
		var addit = new(model.DBAccountRecoveryAddit)
		if err = rows.Scan(&addit.Rid, &addit.Files, &addit.Extra, &addit.Ctime, &addit.Mtime); err != nil {
			log.Error("BatchGetRecoveryAddit rows.Scan error(%v)", err)
			continue
		}
		addits[addit.Rid] = addit
	}
	return
}

// BatchGetLastSuccess batch get last find success info
func (dao *Dao) BatchGetLastSuccess(c context.Context, mids []int64) (lastSuccessMap map[int64]*model.LastSuccessData, err error) {
	rows, err := dao.db.Query(c, fmt.Sprintf(_batchGetLastSuccess, xstr.JoinInts(mids)))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("BatchGetLastSuccess d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	lastSuccessMap = make(map[int64]*model.LastSuccessData)
	for rows.Next() {
		r := new(model.LastSuccessData)
		if err = rows.Scan(&r.LastApplyMID, &r.LastApplyTime); err != nil {
			log.Error("BatchGetLastSuccess rows.Scan error(%v)", err)
			continue
		}
		lastSuccessMap[r.LastApplyMID] = r
	}
	return
}

// GetLastSuccess get last find success info
func (dao *Dao) GetLastSuccess(c context.Context, mid int64) (lastSuc *model.LastSuccessData, err error) {
	row := dao.db.QueryRow(c, _getLastSuccess, mid)
	lastSuc = new(model.LastSuccessData)
	if err = row.Scan(&lastSuc.LastApplyMID, &lastSuc.LastApplyTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetRecoveryAddit (%v) error (%v)", mid, err)
	}
	return
}
