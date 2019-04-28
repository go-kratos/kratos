package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selAppInfo = "SELECT id,name,app_key,purge_url from vip_app_info WHERE `type` = 1"

	//bcoin sql
	_selBcoinSalarySQL         = "SELECT id,mid,status,give_now_status,payday,amount,memo FROM vip_user_bcoin_salary WHERE 1=1"
	_selBcoinSalaryDataSQL     = "SELECT id,mid,status,give_now_status,payday,amount,memo FROM vip_user_bcoin_salary WHERE id>? AND id<=?"
	_selBcoinSalaryByMidSQL    = "SELECT id,mid,status,give_now_status,payday,amount,memo FROM vip_user_bcoin_salary WHERE mid IN (%v)"
	_selOldBcoinSalaryByMidSQL = "SELECT id,mid,IFNULL(status,0),IFNULL(give_now_status,0),month,IFNULL(amount,0),IFNULL(memo,'') FROM vip_bcoin_salary WHERE mid IN (%v)"
	_selOldBcoinSalarySQL      = "SELECT id,mid,IFNULL(status,0),IFNULL(give_now_status,0),month,IFNULL(amount,0),IFNULL(memo,'') FROM vip_bcoin_salary WHERE id>? AND id<=?"
	_addBcoinSalarySQL         = "INSERT INTO vip_user_bcoin_salary(mid,status,give_now_status,payday,amount,memo) VALUES(?,?,?,?,?,?)"
	_updateBcoinSalarySQL      = "UPDATE vip_user_bcoin_salary set status = ? WHERE mid = ? AND payday = ?"
	_updateBcoinSalaryBatchSQL = "UPDATE vip_user_bcoin_salary set status = ? WHERE id in (?)"
	_batchAddBcoinSalarySQL    = "INSERT INTO vip_user_bcoin_salary(mid,status,give_now_status,payday,amount,memo) VALUES "
	_selBcoinMaxIDSQL          = "SELECT IFNULL(MAX(id),0) FROM vip_user_bcoin_salary"
	_selOldBcoinMaxIDSQL       = "SELECT IFNULL(MAX(id),0) FROM vip_bcoin_salary"
	_delBcoinSalarySQL         = "DELETE FROM vip_user_bcoin_salary WHERE mid = ? AND payday = ?"

	_getAbleCodeSQL  = "SELECT code FROM vip_resource_code WHERE batch_code_id = ? AND status = 1 AND relation_id = '' LIMIT 1"
	_selBatchCodeSQL = "SELECT id,business_id,pool_id,status,type,batch_name,reason,unit,count,surplus_count,price,start_time,end_time FROM vip_resource_batch_code WHERE  id = ?"

	_updateCodeRelationIDSQL = "UPDATE vip_resource_code SET relation_id=?,bmid=? WHERE code=?"

	_selEffectiveVipList = "SELECT id,mid,vip_type,vip_status,vip_overdue_time,annual_vip_overdue_time FROM vip_user_info WHERE id>? AND id <=? "

	//push
	_selPushDataSQL    = "SELECT id,disable_type,group_name,title,content,push_total_count,pushed_count,progress_status,`status`,platform,link_type,link_url,error_code,expired_day_start,expired_day_end,effect_start_date,effect_end_date,push_start_time,push_end_time FROM vip_push_data WHERE effect_start_date <= ? AND effect_end_date >= ? "
	_updatePushDataSQL = "UPDATE vip_push_data SET progress_status=?,status=?, pushed_count=?,error_code=?,task_id=? WHERE id=?"
)

//SelOldBcoinMaxID sel oldbcoin maxID
func (d *Dao) SelOldBcoinMaxID(c context.Context) (maxID int64, err error) {
	row := d.oldDb.QueryRow(c, _selOldBcoinMaxIDSQL)
	if err = row.Scan(&maxID); err != nil {
		if err == sql.ErrStmtNil {
			err = nil
			maxID = 0
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
	}
	return
}

//SelBcoinMaxID sel bcoin maxID
func (d *Dao) SelBcoinMaxID(c context.Context) (maxID int64, err error) {
	row := d.db.QueryRow(c, _selBcoinMaxIDSQL)
	if err = row.Scan(&maxID); err != nil {

		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelBcoinSalaryData sel bcoinSalary data
func (d *Dao) SelBcoinSalaryData(c context.Context, startID int64, endID int64) (res []*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selBcoinSalaryDataSQL, startID, endID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Payday, &r.Amount, &r.Memo); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelBcoinSalaryDataMaps sel bcoin salary data convert map
func (d *Dao) SelBcoinSalaryDataMaps(c context.Context, mids []int64) (res map[int64][]*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selBcoinSalaryByMidSQL, xstr.JoinInts(mids))); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	res = make(map[int64][]*model.VipBcoinSalary)
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Payday, &r.Amount, &r.Memo); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		salaries := res[r.Mid]
		salaries = append(salaries, r)
		res[r.Mid] = salaries
	}
	err = rows.Err()
	return
}

//SelOldBcoinSalaryDataMaps sel old bcoin salary data convert map
func (d *Dao) SelOldBcoinSalaryDataMaps(c context.Context, mids []int64) (res map[int64][]*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, fmt.Sprintf(_selOldBcoinSalaryByMidSQL, xstr.JoinInts(mids))); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	res = make(map[int64][]*model.VipBcoinSalary)
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Payday, &r.Amount, &r.Memo); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		salaries := res[r.Mid]
		salaries = append(salaries, r)
		res[r.Mid] = salaries
	}
	err = rows.Err()
	return
}

//SelBcoinSalary sel bcoin salary data by query
func (d *Dao) SelBcoinSalary(c context.Context, arg *model.QueryBcoinSalary) (res []*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	sqlStr := ""
	if arg.StartID >= 0 {
		sqlStr += fmt.Sprintf(" AND id > %v ", arg.StartID)
	}
	if arg.EndID >= 0 {
		sqlStr += fmt.Sprintf(" AND id <= %v ", arg.EndID)
	}
	if arg.StartMonth > 0 {
		sqlStr += fmt.Sprintf(" AND payday>= '%v' ", arg.StartMonth.Time().Format("2006-01-02"))
	}
	if arg.EndMonth > 0 {
		sqlStr += fmt.Sprintf(" AND payday <= '%v' ", arg.EndMonth.Time().Format("2006-01-02"))
	}
	if arg.GiveNowStatus > -1 {
		sqlStr += fmt.Sprintf(" AND give_now_status = %v ", arg.GiveNowStatus)
	}
	if arg.Status > -1 {
		sqlStr += fmt.Sprintf(" AND status = %v ", arg.Status)
	}
	if rows, err = d.db.Query(c, _selBcoinSalarySQL+sqlStr); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Payday, &r.Amount, &r.Memo); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelOldBcoinSalary sel old bcoin salary
func (d *Dao) SelOldBcoinSalary(c context.Context, startID, endID int64) (res []*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _selOldBcoinSalarySQL, startID, endID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Payday, &r.Amount, &r.Memo); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//AddBcoinSalary add bcoin salary
func (d *Dao) AddBcoinSalary(c context.Context, arg *model.VipBcoinSalaryMsg) (err error) {
	if _, err = d.db.Exec(c, _addBcoinSalarySQL, &arg.Mid, &arg.Status, &arg.GiveNowStatus, &arg.Payday, &arg.Amount, &arg.Memo); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_exec")
	}
	return
}

//UpdateBcoinSalary update bcoin salary
func (d *Dao) UpdateBcoinSalary(c context.Context, payday string, mid int64, status int8) (err error) {
	if _, err = d.db.Exec(c, _updateBcoinSalarySQL, status, mid, payday); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_exec")
	}
	return
}

//DelBcoinSalary del bcoin salary
func (d *Dao) DelBcoinSalary(c context.Context, payday string, mid int64) (err error) {
	if _, err = d.db.Exec(c, _delBcoinSalarySQL, mid, payday); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdateBcoinSalaryBatch update bcoin salary batch
func (d *Dao) UpdateBcoinSalaryBatch(c context.Context, ids []int64, status int8) (err error) {
	if len(ids) <= 0 {
		return
	}
	sqlStr := ""
	for _, v := range ids {
		sqlStr += "," + strconv.FormatInt(v, 10)
	}
	if _, err = d.db.Exec(c, _updateBcoinSalaryBatchSQL, status, sqlStr[1:]); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//BatchAddBcoinSalary batch add bcoin salary data
func (d *Dao) BatchAddBcoinSalary(bcoins []*model.VipBcoinSalary) (err error) {
	var values []string
	if len(bcoins) <= 0 {
		return
	}
	for _, v := range bcoins {
		str := fmt.Sprintf("('%v','%v','%v','%v','%v','%v')", v.Mid, v.Status, v.GiveNowStatus, v.Payday.Time().Format("2006-01-02"), v.Amount, v.Memo)
		values = append(values, str)
	}
	valueStr := strings.Join(values, ",")

	if _, err = d.db.Exec(context.TODO(), _batchAddBcoinSalarySQL+valueStr); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_exec")
		return
	}
	return
}

//SelAppInfo sel vip_app_info data
func (d *Dao) SelAppInfo(c context.Context) (res []*model.VipAppInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _selAppInfo); err != nil {
		log.Error("SelAppInfo db.query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipAppInfo)
		if err = rows.Scan(&r.ID, &r.Name, &r.AppKey, &r.PurgeURL); err != nil {
			log.Error("row.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)

	}
	err = rows.Err()
	return
}

//SelEffectiveVipList sel effective vip data
func (d *Dao) SelEffectiveVipList(c context.Context, id, endID int) (res []*model.VipUserInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _selEffectiveVipList, id, endID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Type, &r.Status, &r.OverdueTime, &r.AnnualVipOverdueTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
