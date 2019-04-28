package up

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"go-common/app/admin/main/mcn/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"
)

const (
	_inMCNRenewalSQL = `INSERT INTO mcn_sign(mcn_mid,company_name,company_license_id,company_license_link,
		contract_link,contact_name,contact_title,contact_idcard,contact_phone,begin_date,end_date,state,permission) VALUES (?,?,?,?,?,?,?,?,?,?,?,15,?)`
	_inMCNPaySQL = "INSERT INTO mcn_sign_pay(mcn_mid,sign_id,due_date,pay_value) VALUES %s"
	_inMCNUPsSQL = "INSERT INTO mcn_up(sign_id,mcn_mid,up_mid,begin_date,end_date,contract_link,up_auth_link,state,state_change_time,up_type,site_link,confirm_time,permission,publication_price) VALUES %s"

	_upMCNStateSQL    = "UPDATE mcn_sign SET state=? WHERE id=? AND mcn_mid=?"
	_upMCNPaySQL      = "UPDATE mcn_sign_pay SET due_date=?,pay_value=? WHERE id=? AND mcn_mid=? AND sign_id=?"
	_upMCNPayStateSQL = "UPDATE mcn_sign_pay SET state=? WHERE id=? AND mcn_mid=? AND sign_id=?"
	// _upMCNSignRenewaIDSQL = "UPDATE mcn_sign SET renewal_id=? WHERE id=?"
	_upMCNUPStateSQL    = "UPDATE mcn_up SET state=?,state_change_time=? WHERE id =? AND sign_id=? AND mcn_mid=? AND up_mid=?"
	_upMCNImportUPSQL   = "UPDATE mcn_data_import_up SET is_reward=1 WHERE sign_id=? AND up_mid=?"
	_upMCNPermissionSQL = "UPDATE mcn_sign SET permission = ? WHERE id=?"
	_selMCNListSQL      = `SELECT s.id, s.mcn_mid,s.state,s.begin_date,s.end_date,s.permission,ifnull(ds.up_count,0),ifnull(ds.fans_count_accumulate,0),
	ifnull(ds.fans_count_online_accumulate,0),ifnull(ds.fans_count_real_accumulate,0),ifnull(ds.fans_count_cheat_accumulate,0),ifnull(ds.generate_date,0) 
	FROM (SELECT * FROM mcn_data_summary WHERE generate_date=(SELECT MAX(generate_date) FROM mcn_data_summary a WHERE a.sign_id=mcn_data_summary.sign_id) AND active_tid=? AND data_type=1) AS ds 
	RIGHT JOIN mcn_sign s ON s.id=ds.sign_id WHERE `
	_countMCNListSQL = `SELECT COUNT(*) FROM (SELECT * FROM mcn_data_summary WHERE generate_date=(SELECT MAX(generate_date) FROM mcn_data_summary a 
	WHERE a.sign_id=mcn_data_summary.sign_id) AND active_tid=? AND data_type=1) AS ds RIGHT JOIN mcn_sign s ON s.id=ds.sign_id WHERE %s`
	_selMCNPayInfosByID         = "SELECT id,mcn_mid,sign_id,due_date,pay_value,state FROM mcn_sign_pay WHERE id=?"
	_selMCNPayInfosBySignIDsSQL = "SELECT id,mcn_mid,sign_id,due_date,pay_value,state FROM mcn_sign_pay WHERE sign_id in(%s) ORDER BY due_date ASC"
	_selMCNRenewalUPsSQL        = "SELECT sign_id,mcn_mid,up_mid,begin_date,end_date,contract_link,up_auth_link,reject_reason,reject_time,state,state_change_time,up_type,site_link,confirm_time,permission,publication_price FROM mcn_up WHERE sign_id=? AND mcn_mid=? AND state IN(10,11,15)"
	_selMCNInfoSQL              = `SELECT s.id,s.mcn_mid,s.company_name,s.company_license_id,s.company_license_link,s.contract_link,s.contact_name,s.contact_title,s.contact_idcard,s.contact_phone,
	s.begin_date,s.end_date,ifnull(s.state,0),ifnull(ds.up_count,0),ifnull(ds.fans_count_accumulate,0),ifnull(ds.archive_count_accumulate,0),ifnull(ds.play_count_accumulate,0),
	ifnull(ds.fans_count_cheat_accumulate,0),ifnull(ds.fans_count_real_accumulate,0),ifnull(ds.fans_count_online_accumulate,0) 
	FROM mcn_sign s LEFT JOIN mcn_data_summary ds ON s.id=ds.sign_id AND ds.active_tid=? AND ds.data_type=1 WHERE s.id = ? ORDER BY ds.generate_date DESC LIMIT 1`
	_selMCNUPListSQL = `SELECT u.id,u.sign_id,u.mcn_mid,u.up_mid,u.publication_price,u.permission,ifnull(du.active_tid,0),ifnull(du.fans_count,0),ifnull(du.fans_count_active,0),u.begin_date,u.end_date,u.state,
	ifnull(du.fans_increase_accumulate,0),ifnull(du.archive_count,0),ifnull(du.play_count,0),u.contract_link,u.up_auth_link,u.up_type, ifnull(u.site_link, "") 
	FROM (SELECT * FROM mcn_data_up WHERE generate_date=(SELECT MAX(generate_date) FROM mcn_data_up a WHERE a.up_mid=mcn_data_up.up_mid)) du 
	RIGHT JOIN mcn_up u ON u.sign_id=du.sign_id AND u.up_mid=du.up_mid AND du.data_type=? WHERE u.state IN(10,11,12,14,15) AND `
	_countMCNUPListSQL = `SELECT COUNT(*) FROM (SELECT * FROM mcn_data_up WHERE generate_date=(SELECT MAX(generate_date) FROM mcn_data_up a WHERE a.up_mid=mcn_data_up.up_mid)) du 
	RIGHT JOIN mcn_up u ON u.sign_id=du.sign_id AND u.up_mid=du.up_mid AND du.data_type=? WHERE u.state IN(10,11,12,14,15) AND %s`
	_selMCNByMCNMIDSQL = `SELECT id,mcn_mid,company_name,company_license_id,company_license_link,contract_link,contact_name,
		contact_title,contact_idcard,contact_phone,begin_date,end_date,reject_reason,reject_time,state,ctime,mtime FROM mcn_sign WHERE state IN(10,11,15) AND mcn_mid = ? ORDER BY begin_date DESC LIMIT 1`
	_selMCNCheatListSQL       = `SELECT fans_count_accumulate,fans_count_cheat_accumulate,mcn_mid,sign_id,up_mid,fans_count_cheat_increase_day,fans_count_cheat_cleaned_accumulate FROM mcn_data_up_cheat WHERE `
	_countMCNCheatListSQL     = "SELECT COUNT(*) FROM mcn_data_up_cheat WHERE %s"
	_selMCNCheatUPListSQL     = `SELECT sign_id,generate_date,fans_count_accumulate,fans_count_cheat_accumulate,mcn_mid,fans_count_cheat_increase_day,fans_count_cheat_cleaned_accumulate FROM mcn_data_up_cheat WHERE up_mid=? ORDER BY generate_date DESC LIMIT ?,?`
	_countMCNCheatUPListSQL   = "SELECT COUNT(*) FROM mcn_data_up_cheat WHERE up_mid=?"
	_selMCNImportUPInfoSQL    = "SELECT id,mcn_mid,sign_id,up_mid,standard_fans_date,standard_archive_count,standard_fans_count,is_reward FROM mcn_data_import_up WHERE sign_id=? AND up_mid=? AND standard_fans_type=1"
	_selMCNDataSummaryListSQL = `SELECT id,sign_id,data_type,active_tid,generate_date,up_count,fans_count_online_accumulate,fans_count_real_accumulate,fans_count_cheat_accumulate,fans_count_increase_day,archive_count_accumulate,archive_count_day,
		play_count_accumulate,play_count_increase_day,fans_count_accumulate FROM mcn_data_summary WHERE `
	_countMCNDataSummaryListSQL = "SELECT COUNT(*) FROM mcn_data_summary WHERE %s"
)

// TxAddMCNRenewal .
func (d *Dao) TxAddMCNRenewal(tx *xsql.Tx, arg *model.MCNSign) (lastID int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inMCNRenewalSQL, arg.MCNMID, arg.CompanyName, arg.CompanyLicenseID, arg.CompanyLicenseLink, arg.ContractLink, arg.ContactName, arg.ContactTitle, arg.ContactIdcard, arg.ContactPhone, arg.BeginDate, arg.EndDate, arg.Permission); err != nil {
		return lastID, err
	}
	if lastID, err = res.LastInsertId(); err != nil {
		return lastID, errors.Errorf("TxAddMCNRenewal res.LastInsertId(%+v) error(%+v)", arg, err)
	}
	return lastID, nil
}

// TxAddMCNPays .
func (d *Dao) TxAddMCNPays(tx *xsql.Tx, lastID, mcnMID int64, arg []*model.SignPayReq) (err error) {
	l := len(arg)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*5)
	for _, v := range arg {
		valueStrings = append(valueStrings, "(?,?,?,?)")
		valueArgs = append(valueArgs, strconv.FormatInt(mcnMID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(lastID, 10))
		valueArgs = append(valueArgs, v.DueDate)
		valueArgs = append(valueArgs, strconv.FormatInt(v.PayValue, 10))
	}
	stmt := fmt.Sprintf(_inMCNPaySQL, strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	return
}

// TxAddMCNUPs .
func (d *Dao) TxAddMCNUPs(tx *xsql.Tx, signID, mcnMID int64, arg []*model.MCNUP) (err error) {
	l := len(arg)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*15)
	for _, v := range arg {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, strconv.FormatInt(signID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(mcnMID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.UPMID, 10))
		valueArgs = append(valueArgs, v.BeginDate.Time().Format(model.TimeFormatDay))
		valueArgs = append(valueArgs, v.EndDate.Time().Format(model.TimeFormatDay))
		valueArgs = append(valueArgs, v.ContractLink)
		valueArgs = append(valueArgs, v.UPAuthLink)
		valueArgs = append(valueArgs, strconv.FormatInt(int64(v.State), 10))
		valueArgs = append(valueArgs, v.StateChangeTime.Time().Format(model.TimeFormatSec))
		valueArgs = append(valueArgs, strconv.FormatInt(int64(v.UpType), 10))
		valueArgs = append(valueArgs, v.SiteLink)
		valueArgs = append(valueArgs, v.ConfirmTime.Time().Format(model.TimeFormatSec))
		valueArgs = append(valueArgs, strconv.FormatInt(int64(v.Permission), 10))
		valueArgs = append(valueArgs, strconv.FormatInt(int64(v.PublicationPrice), 10))
	}
	stmt := fmt.Sprintf(_inMCNUPsSQL, strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	return
}

// // TxUpMCNSignRenewaID .
// func (d *Dao) TxUpMCNSignRenewaID(tx *xsql.Tx, signID, renewalID int64) (rows int64, err error) {
// 	var res sql.Result
// 	if res, err = tx.Exec(_upMCNSignRenewaIDSQL, renewalID, signID); err != nil {
// 		return rows, err
// 	}
// 	return res.RowsAffected()
// }

// UpMCNState .
func (d *Dao) UpMCNState(c context.Context, arg *model.MCNStateEditReq) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNStateSQL, arg.State, arg.ID, arg.MCNMID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMCNPay update mcn_sign_pay.
func (d *Dao) UpMCNPay(c context.Context, arg *model.MCNPayEditReq) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNPaySQL, arg.DueDate, arg.PayValue, arg.ID, arg.MCNMID, arg.SignID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMCNPayState update mcn_sign_pay state.
func (d *Dao) UpMCNPayState(c context.Context, arg *model.MCNPayStateEditReq) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNPayStateSQL, arg.State, arg.ID, arg.MCNMID, arg.SignID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMCNImportUPRewardSign .
func (d *Dao) UpMCNImportUPRewardSign(c context.Context, arg *model.MCNImportUPRewardSignReq) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNImportUPSQL, arg.SignID, arg.UPMID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMCNPermission .
func (d *Dao) UpMCNPermission(c context.Context, signID int64, permission uint32) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNPermissionSQL, permission, signID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// MCNList .
func (d *Dao) MCNList(c context.Context, arg *model.MCNListReq) (res []*model.MCNListOne, ids, mids []int64, err error) {
	sql, values := d.buildMCNListSQL("list", arg)
	rows, err := d.db.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MCNListOne{}
		err = rows.Scan(&m.ID, &m.MCNMID, &m.State, &m.BeginDate, &m.EndDate, &m.Permission, &m.UPCount, &m.FansCountAccumulate, &m.FansCountOnlineAccumulate, &m.FansCountRealAccumulate, &m.FansCountCheatAccumulate, &m.GenerateDate)
		if err != nil {
			return
		}
		res = append(res, m)
		ids = append(ids, m.ID)
		mids = append(mids, m.MCNMID)
	}
	return
}

// buildMCNListSQL build a MCNList sql string.
func (d *Dao) buildMCNListSQL(SQLType string, arg *model.MCNListReq) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 11)
	var (
		cond      []string
		condStr   string
		curTime   = time.Now()
		date      = time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.Local)
		orderSign = false
	)
	// defualt where
	cond = append(cond, "s.state IN(10,11,12,13,14,15)")
	values = append(values, model.AllActiveTid)
	if arg.ExpirePay {
		cond = append(cond, "s.pay_expire_state = ?")
		values = append(values, model.MCNStatePayed)
	}
	if arg.ExpireSign {
		cond = append(cond, "s.end_date <= ? AND s.state = 10")
		values = append(values, date.AddDate(0, 0, 30))
	}
	if arg.FansNumMin != 0 {
		cond = append(cond, "ds.fans_count_accumulate >= ?")
		values = append(values, arg.FansNumMin)
	}
	if arg.FansNumMax != 0 && arg.FansNumMax >= arg.FansNumMin {
		cond = append(cond, "ds.fans_count_accumulate <= ?")
		values = append(values, arg.FansNumMax)
	}
	if arg.MCNMID != 0 {
		cond = append(cond, "s.mcn_mid=?")
		values = append(values, arg.MCNMID)
	}
	if arg.State != -1 {
		cond = append(cond, "s.state=?")
		values = append(values, arg.State)
	}
	var permission = arg.GetAttrPermitVal()
	if permission != 0 {
		cond = append(cond, "((permission & ?) = ?)")
		values = append(values, permission, permission)
	}
	if checkSort(arg.SortUP, orderSign) {
		arg.Order = "ds.up_count"
		arg.Sort = arg.SortUP
		orderSign = true
	}
	if checkSort(arg.SortAllFans, orderSign) {
		arg.Order = "ds.fans_count_accumulate"
		arg.Sort = arg.SortAllFans
		orderSign = true
	}
	if checkSort(arg.SortRiseFans, orderSign) {
		arg.Order = "ds.fans_count_online_accumulate"
		arg.Sort = arg.SortRiseFans
		orderSign = true
	}
	if checkSort(arg.SortTrueRiseFans, orderSign) {
		arg.Order = "ds.fans_count_real_accumulate"
		arg.Sort = arg.SortTrueRiseFans
		orderSign = true
	}
	if checkSort(arg.SortCheatFans, orderSign) {
		arg.Order = "ds.fans_count_cheat_accumulate"
		arg.Sort = arg.SortCheatFans
	}
	condStr = d.joinStringSQL(cond)
	switch SQLType {
	case "list":
		if arg.Export == model.ResponeModelCSV {
			sql = fmt.Sprintf(_selMCNListSQL+_orderByConditionNotLimitSQL, condStr, arg.Order, arg.Sort)
			return
		}
		sql = fmt.Sprintf(_selMCNListSQL+_orderByConditionSQL, condStr, arg.Order, arg.Sort)
		limit, offset := arg.PageArg.CheckPageValidation()
		values = append(values, offset, limit)
	case "count":
		sql = fmt.Sprintf(_countMCNListSQL, condStr)
	}
	return
}

func checkSort(arg string, orderSign bool) bool {
	return arg != "" && (arg == "asc" || arg == "desc") && !orderSign
}

// MCNListTotal .
func (d *Dao) MCNListTotal(c context.Context, arg *model.MCNListReq) (count int64, err error) {
	sql, values := d.buildMCNListSQL("count", arg)
	row := d.db.QueryRow(c, sql, values...)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// MCNPayInfo .
func (d *Dao) MCNPayInfo(c context.Context, arg *model.MCNPayStateEditReq) (m *model.SignPayInfoReply, err error) {
	row := d.db.QueryRow(c, _selMCNPayInfosByID, arg.ID)
	m = new(model.SignPayInfoReply)
	// id,mcn_mid,sign_id,due_date,pay_value,state
	if err = row.Scan(&m.SignPayID, &m.McnMid, &m.SignID, &m.DueDate, &m.PayValue, &m.State); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			m = nil
			return
		}
	}
	return
}

// MCNPayInfos .
func (d *Dao) MCNPayInfos(c context.Context, ids []int64) (res map[int64][]*model.SignPayInfoReply, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selMCNPayInfosBySignIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64][]*model.SignPayInfoReply, len(ids))
	for rows.Next() {
		m := &model.SignPayInfoReply{}
		err = rows.Scan(&m.SignPayID, &m.McnMid, &m.SignID, &m.DueDate, &m.PayValue, &m.State)
		if err != nil {
			return
		}
		res[m.SignID] = append(res[m.SignID], m)
	}
	return
}

// TxMCNRenewalUPs .
func (d *Dao) TxMCNRenewalUPs(tx *xsql.Tx, signID, mcnID int64) (ups []*model.MCNUP, err error) {
	var rows *xsql.Rows
	if rows, err = tx.Query(_selMCNRenewalUPsSQL, signID, mcnID); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := new(model.MCNUP)
		if err = rows.Scan(&up.SignID, &up.MCNMID, &up.UPMID, &up.BeginDate, &up.EndDate, &up.ContractLink,
			&up.UPAuthLink, &up.RejectReason, &up.RejectTime, &up.State, &up.StateChangeTime, &up.UpType,
			&up.SiteLink, &up.ConfirmTime, &up.Permission, &up.PublicationPrice); err != nil {
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// MCNInfo .
func (d *Dao) MCNInfo(c context.Context, arg *model.MCNInfoReq) (m *model.MCNInfoReply, err error) {
	row := d.db.QueryRow(c, _selMCNInfoSQL, model.AllActiveTid, arg.ID)
	m = new(model.MCNInfoReply)
	if err = row.Scan(&m.ID, &m.MCNMID, &m.CompanyName, &m.CompanyLicenseID, &m.CompanyLicenseLink, &m.ContractLink, &m.ContactName,
		&m.ContactTitle, &m.ContactIdcard, &m.ContactPhone, &m.BeginDate, &m.EndDate, &m.State, &m.UPCount, &m.FansCountAccumulate,
		&m.ArchiveCountAccumulate, &m.PlayCountAccumulate, &m.FansCountOnline, &m.FansCountRealAccumulate, &m.FansCountOnlineAccumulate); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			m = nil
			return
		}
	}
	return
}

// MCNUPList .
func (d *Dao) MCNUPList(c context.Context, arg *model.MCNUPListReq) (res []*model.MCNUPInfoReply, err error) {
	sql, values := d.buildMCNUPListSQL("list", arg)
	rows, err := d.db.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MCNUPInfoReply{}
		err = rows.Scan(&m.SignUpID, &m.SignID, &m.McnMid, &m.UpMid, &m.PubPrice, &m.Permission, &m.ActiveTid, &m.FansCount, &m.FansCountActive,
			&m.BeginDate, &m.EndDate, &m.State, &m.FansIncreaseAccumulate, &m.ArchiveCount, &m.PlayCount, &m.ContractLink, &m.UpAuthLink, &m.UPType, &m.SiteLink)
		if err != nil {
			return
		}
		res = append(res, m)
	}
	return
}

// MCNUPListTotal .
func (d *Dao) MCNUPListTotal(c context.Context, arg *model.MCNUPListReq) (count int64, err error) {
	sql, values := d.buildMCNUPListSQL("count", arg)
	row := d.db.QueryRow(c, sql, values...)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// buildMCNUPListSQL build a MCNUPList sql string.
func (d *Dao) buildMCNUPListSQL(SQLType string, arg *model.MCNUPListReq) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 8)
	var (
		cond      []string
		condStr   string
		orderSign = false
	)
	values = append(values, arg.DataType)
	cond = append(cond, "u.sign_id=?")
	values = append(values, arg.SignID)
	if arg.FansNumMin != 0 {
		cond = append(cond, "du.fans_count>=?")
		values = append(values, arg.FansNumMin)
	}
	if arg.FansNumMax != 0 && arg.FansNumMax >= arg.FansNumMin {
		cond = append(cond, "du.fans_count<=?")
		values = append(values, arg.FansNumMax)
	}
	if arg.ActiveTID != 0 {
		cond = append(cond, "du.active_tid=?")
		values = append(values, arg.ActiveTID)
	}
	if arg.State != -1 {
		cond = append(cond, "u.state=?")
		values = append(values, arg.State)
	}
	if arg.UPMID != 0 {
		cond = append(cond, "u.up_mid=?")
		values = append(values, arg.UPMID)
	}
	if arg.UpType != -1 {
		cond = append(cond, "u.up_type=?")
		values = append(values, arg.UpType)
	}

	var permission = arg.GetAttrPermitVal()
	if permission != 0 {
		cond = append(cond, "((permission & ?) = ?)")
		values = append(values, permission, permission)
	}
	if checkSort(arg.SortFansCount, orderSign) {
		arg.Order = "du.fans_count"
		arg.Sort = arg.SortFansCount
		orderSign = true
	}
	if checkSort(arg.SortFansCountActive, orderSign) {
		arg.Order = "du.fans_count_active"
		arg.Sort = arg.SortFansCountActive
		orderSign = true
	}
	if checkSort(arg.SortFansIncreaseAccumulate, orderSign) {
		arg.Order = "du.fans_increase_accumulate"
		arg.Sort = arg.SortFansIncreaseAccumulate
		orderSign = true
	}
	if checkSort(arg.SortArchiveCount, orderSign) {
		arg.Order = "du.archive_count"
		arg.Sort = arg.SortArchiveCount
		orderSign = true
	}
	if checkSort(arg.SortPlayCount, orderSign) {
		arg.Order = "du.play_count"
		arg.Sort = arg.SortPlayCount
		orderSign = true
	}
	if checkSort(arg.SortPubPrice, orderSign) {
		arg.Order = "u.publication_price"
		arg.Sort = arg.SortPubPrice
	}
	condStr = d.joinStringSQL(cond)
	switch SQLType {
	case "list":
		if arg.Export == model.ResponeModelCSV {
			sql = fmt.Sprintf(_selMCNUPListSQL+_orderByConditionNotLimitSQL, condStr, arg.Order, arg.Sort)
			return
		}
		sql = fmt.Sprintf(_selMCNUPListSQL+_orderByConditionSQL, condStr, arg.Order, arg.Sort)
		limit, offset := arg.PageArg.CheckPageValidation()
		values = append(values, offset, limit)
	case "count":
		sql = fmt.Sprintf(_countMCNUPListSQL, condStr)
		fmt.Println(sql)
	}
	return
}

// UpMCNUPState .
func (d *Dao) UpMCNUPState(c context.Context, arg *model.MCNUPStateEditReq) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMCNUPStateSQL, arg.State, time.Now(), arg.ID, arg.SignID, arg.MCNMID, arg.UPMID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// McnSignByMCNMID .
func (d *Dao) McnSignByMCNMID(c context.Context, MCNID int64) (m *model.MCNSignInfoReply, err error) {
	row := d.db.QueryRow(c, _selMCNByMCNMIDSQL, MCNID)
	m = new(model.MCNSignInfoReply)
	if err = row.Scan(&m.SignID, &m.McnMid, &m.CompanyName, &m.CompanyLicenseID, &m.CompanyLicenseLink, &m.ContractLink, &m.ContactName,
		&m.ContactTitle, &m.ContactIdcard, &m.ContactPhone, &m.BeginDate, &m.EndDate, &m.RejectReason, &m.RejectTime, &m.State, &m.Ctime, &m.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			m = nil
			return
		}
	}
	return
}

// MCNCheatList .
func (d *Dao) MCNCheatList(c context.Context, arg *model.MCNCheatListReq) (res []*model.MCNCheatReply, mids []int64, err error) {
	sql, values := d.buildMCNCheatListSQL("list", arg)
	rows, err := d.db.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MCNCheatReply{}
		err = rows.Scan(&m.FansCountReal, &m.FansCountCheatAccumulate, &m.MCNMID, &m.SignID, &m.UpMID, &m.FansCountCheatIncreaseDay, &m.FansCountCheatCleanedAccumulate)
		if err != nil {
			return
		}

		res = append(res, m)
		mids = append(mids, m.UpMID, m.MCNMID)
	}
	mids = SliceUnique(mids)
	return
}

// MCNCheatListTotal .
func (d *Dao) MCNCheatListTotal(c context.Context, arg *model.MCNCheatListReq) (count int64, err error) {
	sql, values := d.buildMCNCheatListSQL("count", arg)
	row := d.db.QueryRow(c, sql, values...)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

//SliceUnique unique
func SliceUnique(s []int64) []int64 {
	result := make([]int64, 0, len(s))
	temp := map[int64]struct{}{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// buildMCNCheatListSQL build a MCNCheatList sql string.
func (d *Dao) buildMCNCheatListSQL(SQLType string, arg *model.MCNCheatListReq) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 3)
	var (
		cond    []string
		condStr string
	)
	cond = append(cond, "generate_date=(SELECT generate_date FROM mcn_data_up_cheat ORDER BY generate_date DESC LIMIT 1)")
	if arg.MCNMID > 0 {
		cond = append(cond, "mcn_mid=?")
		values = append(values, arg.MCNMID)
	}
	if arg.UPMID > 0 {
		cond = append(cond, "up_mid=?")
		values = append(values, arg.UPMID)
	}
	condStr = d.joinStringSQL(cond)
	switch SQLType {
	case "list":
		sql = fmt.Sprintf(_selMCNCheatListSQL+_orderByConditionSQL, condStr, "generate_date DESC,fans_count_cheat_accumulate", "DESC")
		limit, offset := arg.PageArg.CheckPageValidation()
		values = append(values, offset, limit)
	case "count":
		sql = fmt.Sprintf(_countMCNCheatListSQL, condStr)
	}
	return
}

// MCNCheatUPList .
func (d *Dao) MCNCheatUPList(c context.Context, arg *model.MCNCheatUPListReq) (res []*model.MCNCheatUPReply, err error) {
	limit, offset := arg.PageArg.CheckPageValidation()
	rows, err := d.db.Query(c, _selMCNCheatUPListSQL, arg.UPMID, offset, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MCNCheatUPReply{}
		err = rows.Scan(&m.SignID, &m.GenerateDate, &m.FansCountReal, &m.FansCountCheatAccumulate, &m.MCNMID, &m.FansCountCheatIncreaseDay, &m.FansCountCheatCleanedAccumulate)
		if err != nil {
			return
		}
		res = append(res, m)
	}
	return
}

// MCNCheatUPListTotal .
func (d *Dao) MCNCheatUPListTotal(c context.Context, arg *model.MCNCheatUPListReq) (count int64, err error) {
	row := d.db.QueryRow(c, _countMCNCheatUPListSQL, arg.UPMID)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// MCNImportUPInfo .
func (d *Dao) MCNImportUPInfo(c context.Context, arg *model.MCNImportUPInfoReq) (m *model.MCNImportUPInfoReply, err error) {
	row := d.db.QueryRow(c, _selMCNImportUPInfoSQL, arg.SignID, arg.UPMID)
	m = new(model.MCNImportUPInfoReply)
	if err = row.Scan(&m.ID, &m.MCNMID, &m.SignID, &m.UpMID, &m.StandardFansDate, &m.StandardArchiveCount, &m.StandardFansCount, &m.IsReward); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// MCNIncreaseList .
func (d *Dao) MCNIncreaseList(c context.Context, arg *model.MCNIncreaseListReq) (res []*model.MCNIncreaseReply, err error) {
	sql, values := d.buildMCNIncreaseListSQL("list", arg)
	rows, err := d.db.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MCNIncreaseReply{}
		err = rows.Scan(&m.ID, &m.SignID, &m.DataType, &m.ActiveTID, &m.GenerateDate, &m.UPCount, &m.FansCountOnlineAccumulate, &m.FansCountRealAccumulate, &m.FansCountCheatAccumulate, &m.FansCountIncreaseDay,
			&m.ArchiveCountAccumulate, &m.ArchiveCountDay, &m.PlayCountAccumulate, &m.PlayCountIncreaseDay, &m.FansCountAccumulate)
		if err != nil {
			return
		}
		res = append(res, m)
	}
	return
}

// MCNIncreaseListTotal .
func (d *Dao) MCNIncreaseListTotal(c context.Context, arg *model.MCNIncreaseListReq) (count int64, err error) {
	sql, values := d.buildMCNIncreaseListSQL("count", arg)
	row := d.db.QueryRow(c, sql, values...)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// buildMCNIncreaseListSQL build a MCNIncreaseList sql string.
func (d *Dao) buildMCNIncreaseListSQL(SQLType string, arg *model.MCNIncreaseListReq) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 3)
	var (
		cond    []string
		condStr string
	)
	cond = append(cond, "sign_id=?")
	values = append(values, arg.SignID)
	if arg.DataType > 0 {
		cond = append(cond, "data_type=?")
		values = append(values, arg.DataType)
	}
	if arg.ActiveTID == 0 {
		arg.ActiveTID = model.AllActiveTid
	}
	if arg.ActiveTID > 0 {
		cond = append(cond, "active_tid=?")
		values = append(values, arg.ActiveTID)
	}
	condStr = d.joinStringSQL(cond)
	switch SQLType {
	case "list":
		sql = fmt.Sprintf(_selMCNDataSummaryListSQL+_orderByConditionSQL, condStr, "generate_date", "DESC")
		limit, offset := arg.PageArg.CheckPageValidation()
		values = append(values, offset, limit)
	case "count":
		sql = fmt.Sprintf(_countMCNDataSummaryListSQL, condStr)
	}
	return
}
