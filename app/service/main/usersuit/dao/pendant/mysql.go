package pendant

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/service/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_getGroupOnlineSQL        = "SELECT id,name,rank,status,image,image_model,frequency_limit,time_limit FROM pendant_group WHERE status = 1 or id in (30,31) ORDER BY rank"
	_getGroupByIDSQL          = "SELECT id,name,rank,status,image,image_model,frequency_limit,time_limit FROM pendant_group where id = ?"
	_selGIDRefPIDSQL          = "SELECT gid,pid FROM pendant_group_ref"
	_getPendantInfoSQL        = "SELECT id,name,image,image_model,status FROM pendant_info"
	_getPendantInfoByIDSQL    = "SELECT id,name,image,image_model,status FROM pendant_info where id = ? limit 1"
	_getPendantInfosSQL       = "SELECT id,name,image,image_model,status,rank FROM pendant_info ORDER BY rank"
	_getPendantPriceSQL       = "SELECT pid,type,price FROM pendant_price where pid = ? "
	_getOrderHistorySQL       = "SELECT mid,order_id,pay_id,appid,status,pid,time_length,cost,buy_time,is_callback,callback_time,pay_type FROM user_pendant_order WHERE %s"
	_getUserPackageSQL        = "SELECT mid,pid,expires,type,status,is_vip FROM user_pendant_pkg WHERE mid = ? AND pid = ? "
	_getUserPackageByMidSQL   = "SELECT mid,pid,expires,type,status,is_vip FROM user_pendant_pkg WHERE mid = ? AND expires >= ? AND status > 0 ORDER BY mtime DESC"
	_countOrderHistorySQL     = "SELECT count(1) FROM user_pendant_order WHERE %s"
	_getPendantEquipByMidSQL  = "SELECT mid,pid,expires FROM user_pendant_equip WHERE mid = ? and expires >= ?"
	_getPendantEquipByMidsSQL = "SELECT mid,pid,expires FROM user_pendant_equip WHERE mid IN (%s) and expires >= ?"
	_insertPendantPackageSQL  = "INSERT INTO user_pendant_pkg(mid,pid,expires,type,status,is_vip) VALUES (?,?,?,?,?,?)"
	_insertOrderHistory       = "INSERT INTO user_pendant_order(mid,order_id,pay_id,appid,status,pid,time_length,cost,buy_time,is_callback,callback_time,pay_type) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
	_insertOperationSQL       = "INSERT INTO pendant_grant_history(mid,pid,source_type,operator_name,operator_action) VALUES (?,?,?,?,?)"
	_insertEquipSQL           = "INSERT INTO user_pendant_equip(mid,pid,expires) VALUES (?,?,?) ON DUPLICATE KEY UPDATE pid=VALUES(pid),expires=VALUES(expires)"
	_updatePackageSQL         = "UPDATE user_pendant_pkg %s WHERE mid=? AND pid=?"
	_updatePackageExpireSQL   = "UPDATE user_pendant_pkg SET status=0 WHERE mid=? AND expires<?"
	_updateOrderInfoSQL       = "UPDATE user_pendant_order SET status=?,pay_id=?,is_callback=?,callback_time=? WHERE order_id=?"
	_updateEquipMIDSQL        = "UPDATE user_pendant_equip SET pid=0,expires=0 WHERE mid=?"
)

//PendantGroupInfo return all group info
func (d *Dao) PendantGroupInfo(c context.Context) (res []*model.PendantGroupInfo, err error) {
	var row *xsql.Rows
	res = make([]*model.PendantGroupInfo, 0)

	if row, err = d.db.Query(c, _getGroupOnlineSQL); err != nil {
		log.Error("PendantGroupInfo query error %v", err)
		return
	}
	defer row.Close()
	for row.Next() {
		info := new(model.PendantGroupInfo)
		if err = row.Scan(&info.ID, &info.Name, &info.Rank, &info.Status, &info.Image, &info.ImageModel, &info.FrequencyLimit, &info.TimeLimit); err != nil {
			log.Error("PendantGroupInfo scan error %v", err)
			return
		}
		res = append(res, info)
	}
	return
}

// GroupByID return group info by id
func (d *Dao) GroupByID(c context.Context, gid int64) (res *model.PendantGroupInfo, err error) {
	var row *xsql.Row
	res = new(model.PendantGroupInfo)
	row = d.db.QueryRow(c, _getGroupByIDSQL, gid)

	if err = row.Scan(&res.ID, &res.Name, &res.Rank, &res.Status, &res.Image, &res.ImageModel, &res.FrequencyLimit, &res.TimeLimit); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("PendantGroupInfo scan error %v", err)
		return
	}
	return
}

// GIDRefPID gid relation of pid .
func (d *Dao) GIDRefPID(c context.Context) (gidMap map[int64][]int64, pidMap map[int64]int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selGIDRefPIDSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	gidMap = make(map[int64][]int64)
	pidMap = make(map[int64]int64)
	for rows.Next() {
		var gid, pid int64
		if err = rows.Scan(&gid, &pid); err != nil {
			if err == xsql.ErrNoRows {
				gidMap = nil
				pidMap = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		pidMap[pid] = gid
		gidMap[gid] = append(gidMap[gid], pid)
	}
	return
}

// PendantList return pendant info
func (d *Dao) PendantList(c context.Context) (res []*model.Pendant, err error) {
	var (
		row *xsql.Rows
	)
	res = make([]*model.Pendant, 0)

	if row, err = d.db.Query(c, _getPendantInfosSQL); err != nil {
		log.Error("PendantInfo query error %v", err)
		return
	}
	defer row.Close()

	for row.Next() {
		info := new(model.Pendant)
		if err = row.Scan(&info.ID, &info.Name, &info.Image, &info.ImageModel, &info.Status, &info.Rank); err != nil {
			log.Error("PendantInfo scan error %v", err)
			return
		}
		res = append(res, info)
	}
	return
}

// Pendants return pendant info by ids
func (d *Dao) Pendants(c context.Context, pids []int64) (res []*model.Pendant, err error) {
	var (
		row *xsql.Rows
		bf  bytes.Buffer
	)
	res = make([]*model.Pendant, 0)
	bf.WriteString(_getPendantInfoSQL)
	bf.WriteString(" where id in(")
	bf.WriteString(xstr.JoinInts(pids))
	bf.WriteString(") and status = 1 ORDER BY rank")

	if row, err = d.db.Query(c, bf.String()); err != nil {
		log.Error("Pendants query error %v", err)
		return
	}

	defer row.Close()

	for row.Next() {
		info := new(model.Pendant)
		if err = row.Scan(&info.ID, &info.Name, &info.Image, &info.ImageModel, &info.Status); err != nil {
			log.Error("Pendants scan error %v", err)
			return
		}
		res = append(res, info)
	}

	return
}

// PendantInfo return pendant info by id
func (d *Dao) PendantInfo(c context.Context, pid int64) (res *model.Pendant, err error) {
	var (
		row *xsql.Row
	)
	res = new(model.Pendant)

	row = d.db.QueryRow(c, _getPendantInfoByIDSQL, pid)
	if err = row.Scan(&res.ID, &res.Name, &res.Image, &res.ImageModel, &res.Status); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("Pendant scan error %v", err)
		return
	}
	return
}

// PendantPrice return pendant price
func (d *Dao) PendantPrice(c context.Context, pid int64) (res map[int64]*model.PendantPrice, err error) {
	var row *xsql.Rows
	res = make(map[int64]*model.PendantPrice)

	if row, err = d.db.Query(c, _getPendantPriceSQL, pid); err != nil {
		log.Error("PendantPrice query error %v", err)
		return
	}
	defer row.Close()

	for row.Next() {
		info := new(model.PendantPrice)
		if err = row.Scan(&info.Pid, &info.Type, &info.Price); err != nil {
			log.Error("PendantPrice scan error %v", err)
			return
		}
		res[info.Type] = info
	}
	return
}

// getOrderInfoSQL return a sql string
func (d *Dao) getOrderInfoSQL(c context.Context, arg *model.ArgOrderHistory, tp string) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 5)
	var cond bytes.Buffer
	cond.WriteString("mid = ?")
	values = append(values, arg.Mid)

	if arg.OrderID != "" {
		cond.WriteString(" AND order_id = ?")
		values = append(values, arg.OrderID)
	}
	if arg.Pid != 0 {
		cond.WriteString(" AND pid = ?")
		values = append(values, arg.Pid)
	}
	if arg.Status != 0 {
		cond.WriteString(" AND status = ?")
		values = append(values, arg.Status)
	}
	if arg.PayType != 0 {
		cond.WriteString(" AND pay_type = ?")
		values = append(values, arg.PayType)
	}
	if arg.PayID != "" {
		cond.WriteString(" AND pay_id = ?")
		values = append(values, arg.PayID)
	}

	if arg.StartTime != 0 {
		cond.WriteString(" AND buy_time >= ?")
		values = append(values, arg.StartTime)
	}
	if arg.EndTime != 0 {
		cond.WriteString(" AND buy_time <= ?")
		values = append(values, arg.EndTime)
	}
	if tp == "info" {
		cond.WriteString(" order by buy_time DESC LIMIT ?,20")
		values = append(values, (arg.Page-1)*20)
		sql = fmt.Sprintf(_getOrderHistorySQL, cond.String())
	} else if tp == "count" {
		sql = fmt.Sprintf(_countOrderHistorySQL, cond.String())
	}
	return
}

// OrderInfo return order info
func (d *Dao) OrderInfo(c context.Context, arg *model.ArgOrderHistory) (res []*model.PendantOrderInfo, count int64, err error) {
	sqlstr, values := d.getOrderInfoSQL(c, arg, "info")
	var (
		row *xsql.Rows
		r   *xsql.Row
	)

	res = make([]*model.PendantOrderInfo, 0)

	if row, err = d.db.Query(c, sqlstr, values...); err != nil {
		log.Error("PendantOrderInfo query error %v", err)
		return
	}
	defer row.Close()
	cstr, values2 := d.getOrderInfoSQL(c, arg, "count")
	r = d.db.QueryRow(c, cstr, values2...)

	for row.Next() {
		info := new(model.PendantOrderInfo)
		if err = row.Scan(&info.Mid, &info.OrderID, &info.PayID, &info.AppID, &info.Stauts, &info.Pid, &info.TimeLength, &info.Cost, &info.BuyTime, &info.IsCallback, &info.CallbackTime, &info.PayType); err != nil {
			log.Error("PendantOrderInfo scan error %v", err)
			return
		}
		if info.PayType == 3 {
			info.PayPrice = info.PayPrice / 100
		}
		res = append(res, info)
	}
	err = r.Scan(&count)
	if err == xsql.ErrNoRows {
		res = nil
		err = nil
		return
	}
	return
}

// OrderInfoByID return order info by order id
func (d *Dao) OrderInfoByID(c context.Context, orderID string) (res *model.PendantOrderInfo, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getOrderHistorySQL, "order_id=?"), orderID)
	res = new(model.PendantOrderInfo)
	if err = row.Scan(&res.Mid, &res.OrderID, &res.PayID, &res.AppID, &res.Stauts, &res.Pid, &res.TimeLength, &res.Cost, &res.BuyTime, &res.IsCallback, &res.CallbackTime, &res.PayType); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("OrderInfoByID scan error %v", err)
		return
	}
	return
}

// AddOrderInfo add order log
func (d *Dao) AddOrderInfo(c context.Context, arg *model.PendantOrderInfo) (id int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertOrderHistory, arg.Mid, arg.OrderID, arg.PayID, arg.AppID, arg.Stauts, arg.Pid, arg.TimeLength, arg.Cost, arg.BuyTime, arg.IsCallback, arg.CallbackTime, arg.PayType); err != nil {
		log.Error("AddOrderInfo insert error %v", err)
		return
	}
	return res.LastInsertId()
}

// TxAddOrderInfo add order log
func (d *Dao) TxAddOrderInfo(c context.Context, arg *model.PendantOrderInfo, tx *xsql.Tx) (id int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertOrderHistory, arg.Mid, arg.OrderID, arg.PayID, arg.AppID, arg.Stauts, arg.Pid, arg.TimeLength, arg.Cost, arg.BuyTime, arg.IsCallback, arg.CallbackTime, arg.PayType); err != nil {
		log.Error("TxAddOrderInfo insert error %v", err)
		return
	}
	return res.LastInsertId()
}

// UpdateOrderInfo update order info
func (d *Dao) UpdateOrderInfo(c context.Context, arg *model.PendantOrderInfo) (id int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateOrderInfoSQL, arg.Stauts, arg.PayID, arg.IsCallback, arg.CallbackTime, arg.OrderID); err != nil {
		log.Error("UpdateOrderInfo update error %v", err)
		return
	}
	return res.LastInsertId()
}

// TxUpdateOrderInfo update order info
func (d *Dao) TxUpdateOrderInfo(c context.Context, arg *model.PendantOrderInfo, tx *xsql.Tx) (id int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateOrderInfoSQL, arg.Stauts, arg.PayID, arg.IsCallback, arg.CallbackTime, arg.OrderID); err != nil {
		log.Error("UpdateOrderInfo update error %v", err)
		return
	}
	return res.LastInsertId()
}

// PackageByMid get pendant in user's package
func (d *Dao) PackageByMid(c context.Context, mid int64) (res []*model.PendantPackage, err error) {
	var (
		row *xsql.Rows
		t   = time.Now().Unix()
	)
	res = make([]*model.PendantPackage, 0)

	if row, err = d.db.Query(c, _getUserPackageByMidSQL, mid, t); err != nil {
		log.Error("Package query error %v", err)
		return
	}
	defer row.Close()
	for row.Next() {
		info := new(model.PendantPackage)
		if err = row.Scan(&info.Mid, &info.Pid, &info.Expires, &info.Type, &info.Status, &info.IsVIP); err != nil {
			log.Error("Package scan error %v", err)
			return
		}
		res = append(res, info)
	}
	return
}

// PackageByID get pendant in user's package
func (d *Dao) PackageByID(c context.Context, mid, pid int64) (res *model.PendantPackage, err error) {
	var row *xsql.Row
	res = new(model.PendantPackage)
	row = d.db.QueryRow(c, _getUserPackageSQL, mid, pid)
	if err = row.Scan(&res.Mid, &res.Pid, &res.Expires, &res.Type, &res.Status, &res.IsVIP); err != nil {
		if err == xsql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("Package scan error %v", err)
		return
	}
	return
}

// EquipByMid obtain pendant equiped
func (d *Dao) EquipByMid(c context.Context, mid, t int64) (res *model.PendantEquip, noRow bool, err error) {
	var row *xsql.Row
	res = new(model.PendantEquip)
	row = d.db.QueryRow(c, _getPendantEquipByMidSQL, mid, t)
	if err = row.Scan(&res.Mid, &res.Pid, &res.Expires); err != nil {
		if err == xsql.ErrNoRows {
			noRow = true
			res = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

// EquipByMids obtain equipss by mids .
func (d *Dao) EquipByMids(c context.Context, mids []int64, t int64) (res map[int64]*model.PendantEquip, err error) {
	res = make(map[int64]*model.PendantEquip)
	rows, err := d.db.Query(c, fmt.Sprintf(_getPendantEquipByMidsSQL, xstr.JoinInts(mids)), t)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pe := &model.PendantEquip{}
		if err = rows.Scan(&pe.Mid, &pe.Pid, &pe.Expires); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			err = errors.WithStack(err)
		}
		if _, ok := res[pe.Mid]; !ok {
			res[pe.Mid] = pe
		}
	}
	err = rows.Err()
	return
}

// AddEquip add equip
func (d *Dao) AddEquip(c context.Context, arg *model.PendantEquip) (n int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertEquipSQL, arg.Mid, arg.Pid, arg.Expires); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpEquipMID uninstall user pid by mid.
func (d *Dao) UpEquipMID(c context.Context, mid int64) (n int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateEquipMIDSQL, mid); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// TxUpdatePackageInfo update package info
func (d *Dao) TxUpdatePackageInfo(c context.Context, arg *model.PendantPackage, tx *xsql.Tx) (n int64, err error) {
	var (
		bf     bytes.Buffer
		values = make([]interface{}, 0, 4)
		res    sql.Result
	)
	if arg.Status != 0 && arg.Expires != 0 {
		bf.WriteString("SET status=?,expires=?,type=?")
		values = append(values, arg.Status)
		values = append(values, arg.Expires)
		values = append(values, arg.Type)
	} else if arg.Status != 0 {
		bf.WriteString("SET status=?,type=?")
		values = append(values, arg.Status)
		values = append(values, arg.Type)
	} else if arg.Expires != 0 {
		bf.WriteString("SET expires=?,type=?")
		values = append(values, arg.Expires)
		values = append(values, arg.Type)
	}
	values = append(values, arg.Mid)
	values = append(values, arg.Pid)
	if res, err = tx.Exec(fmt.Sprintf(_updatePackageSQL, bf.String()), values...); err != nil {
		log.Error("TxUpdatePackageInfo update error %v", err)
		return
	}
	return res.RowsAffected()
}

// CheckPackageExpire check expire items and update
func (d *Dao) CheckPackageExpire(c context.Context, mid, expires int64) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updatePackageExpireSQL, mid, expires); err != nil {
		log.Error("CheckPackageExpire error %v", err)
		return
	}
	return res.RowsAffected()
}

// BeginTran begin a tx.
func (d *Dao) BeginTran(c context.Context) (res *xsql.Tx, err error) {
	if res, err = d.db.Begin(c); err != nil || res == nil {
		log.Error("BeginTran  error %v", err)
		return
	}
	return
}

// TxAddPackage add a pendant in package
func (d *Dao) TxAddPackage(c context.Context, arg *model.PendantPackage, tx *xsql.Tx) (id int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertPendantPackageSQL, arg.Mid, arg.Pid, arg.Expires, arg.Type, arg.Status, arg.IsVIP); err != nil {
		log.Error("TxAddPackage insert error %v", err)
		return
	}
	return res.LastInsertId()
}

// TxAddHistory add a history of operation
func (d *Dao) TxAddHistory(c context.Context, arg *model.PendantHistory, tx *xsql.Tx) (id int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertOperationSQL, arg.Mid, arg.Pid, arg.SourceType, arg.OperatorName, arg.OperatorAction); err != nil {
		log.Error("TxAddHistory insert error %v", err)
		return
	}
	return res.LastInsertId()
}
