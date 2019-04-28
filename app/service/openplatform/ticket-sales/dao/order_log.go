package dao

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	time2 "time"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/time"
)

const (
	_addOrderLog = "insert into order_log (uid,ip,order_id,op_data,remark,op_object,op_name,ctime) values(?,?,?,?,?,?,?,?)"

	_getOrderLogList = "select id,uid, order_id, ip,op_data,remark,op_object,op_name,ctime,mtime from  order_log where order_id =  %v order by %v  desc limit %v,%v"
	_getOrderCnt     = "select count(1)  from  order_log where order_id =  ? "
)

//GetOrderLogList 获取订单日志列表
func (d *Dao) GetOrderLogList(c context.Context, OID int64, index int64, size int64, orderBy string) (res []*v1.OrderLog, err error) {
	res = make([]*v1.OrderLog, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_getOrderLogList, OID, orderBy, index, size))
	if err != nil {
		log.Warn("d.GetOrderLogList(%v) d.db.Query() error(%v)", OID, err)
		return
	}

	defer rows.Close()
	var ip net.IP
	for rows.Next() {
		r := &v1.OrderLog{}
		if err = rows.Scan(&r.ID, &r.UID, &r.OID, &ip, &r.OpData, &r.Remark, &r.OpObject, &r.OpName, &r.CTime, &r.MTime); err != nil {
			return
		}
		r.IP = ip.String()
		res = append(res, r)
	}
	return
}

//GetOrderLogCnt 获取订单日志条数
func (d *Dao) GetOrderLogCnt(c context.Context, OID int64) (cnt int64, err error) {
	err = d.db.QueryRow(c, _getOrderCnt, OID).Scan(&cnt)
	if err != nil {
		log.Warn("d.GetOrderLogCnt error(%v)", err)
		return
	}
	return
}

//AddOrderLog 添加订单日志
func (d *Dao) AddOrderLog(c context.Context, oi *v1.OrderLog) (cnt int64, err error) {
	var res sql.Result
	if oi.CTime == 0 {
		oi.CTime = time.Time(time2.Now().Unix())
	}
	var ip []byte
	if ip = net.ParseIP(metadata.String(c, metadata.RemoteIP)); ip == nil {
		ip = []byte{}
	}
	if res, err = d.db.Exec(c, _addOrderLog, oi.UID, ip, oi.OID, oi.OpData, oi.Remark, oi.OpObject, oi.OpName, oi.CTime); err != nil {
		log.Warn("创建订单日志%d失败", oi.OID)
		return
	}
	return res.LastInsertId()
}
