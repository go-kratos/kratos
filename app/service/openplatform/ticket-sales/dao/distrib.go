package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/log"
)

const (
	_hasDisOrder = "SELECT order_id,status FROM dist_order where order_id=? and status=?"
	_getOrder    = "SELECT order_id,cm_amount,cm_method,cm_price,dist_user,status,pid,count,sid,type,payment_amount,serial_num,ctime,mtime FROM dist_order where order_id = ?"
	_addDisOrder = "INSERT INTO dist_order(order_id,cm_amount,cm_method,cm_price,dist_user,status,pid,count,sid,type,payment_amount,serial_num) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
)

// InsertOrder 同步订单到分销表方法
func (d *Dao) InsertOrder(c context.Context, oi *model.DistOrderArg) (lastID int64, err error) {
	res, err := d.db.Exec(c, _addDisOrder, &oi.Oid, &oi.CmAmount, &oi.CmMethod, &oi.CmPrice, &oi.Duid, &oi.Stat, &oi.Pid, &oi.Count, &oi.Sid, &oi.Type, &oi.PayAmount, &oi.Serial)
	if err != nil {
		errors.Wrap(err, fmt.Sprintf("db.Exec(%s) err ", _addDisOrder))
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// HasOrder 检查订单是否存在
func (d *Dao) HasOrder(c context.Context, oi *model.DistOrderArg) (has bool, err error) {
	has = false
	row := d.db.QueryRow(c, _hasDisOrder, &oi.Oid, &oi.Stat)
	err = row.Scan(&oi.Oid, &oi.Stat)
	if err == nil {
		has = true
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errors.Wrap(err, "dao Hasorder err")
		return
	}
	return
}

// GetOrder 检查订单是否存在
func (d *Dao) GetOrder(c context.Context, oid uint64) (res []*model.OrderInfo, err error) {
	rows, err := d.db.Query(c, _getOrder, oid)
	if err != nil {
		log.Error("[dao.distrib|GetOrder] d.db.Query err: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		oi := &model.OrderInfo{}

		if err = rows.Scan(&oi.Oid, &oi.CmAmount, &oi.CmMethod, &oi.CmPrice, &oi.Duid, &oi.Stat, &oi.Pid, &oi.Count, &oi.Sid, &oi.Type, &oi.PayAmount, &oi.Serial, &oi.Ctime, &oi.Mtime); err != nil {
			log.Error("[dao.distrib|GetOrder] rows.Scan err: %v", err)
			return
		}
		res = append(res, oi)
	}
	return
}
