package telecom

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/telecom"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_inOrderSyncSQL = `INSERT IGNORE INTO telecom_order (request_no,result_type,flowpackageid,flowpackagesize,flowpackagetype,trafficattribution,begintime,endtime,
		ismultiplyorder,settlementtype,operator,order_status,remainedrebindnum,maxbindnum,orderid,sign_no,accesstoken,phoneid,isrepeatorder,paystatus,
		paytime,paychannel,sign_status,refund_status) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE
		request_no=?,result_type=?,flowpackagesize=?,flowpackagetype=?,trafficattribution=?,begintime=?,endtime=?,
		ismultiplyorder=?,settlementtype=?,operator=?,order_status=?,remainedrebindnum=?,maxbindnum=?,orderid=?,sign_no=?,accesstoken=?,
		isrepeatorder=?,paystatus=?,paytime=?,paychannel=?,sign_status=?,refund_status=?`
	_inRechargeSyncSQL = `INSERT INTO telecom_recharge (request_no,fcrecharge_no,recharge_status,ordertotalsize,flowbalance) VALUES(?,?,?,?,?)`
	_orderByPhoneSQL   = `SELECT phoneid,orderid,order_status,sign_no,isrepeatorder,begintime,endtime FROM telecom_order WHERE phoneid=?`
	_orderByOrderIDSQL = `SELECT phoneid,orderid,order_status,sign_no,isrepeatorder,begintime,endtime FROM telecom_order WHERE orderid=?`
)

type Dao struct {
	c                    *conf.Config
	client               *httpx.Client
	payInfoURL           string
	cancelRepeatOrderURL string
	sucOrderListURL      string
	telecomReturnURL     string
	telecomCancelPayURL  string
	phoneAreaURL         string
	orderStateURL        string
	smsSendURL           string
	phoneKeyExpired      int32
	payKeyExpired        int32
	db                   *xsql.DB
	inOrderSyncSQL       *xsql.Stmt
	inRechargeSyncSQL    *xsql.Stmt
	orderByPhoneSQL      *xsql.Stmt
	orderByOrderIDSQL    *xsql.Stmt
	phoneRds             *redis.Pool
	// memcache
	mc     *memcache.Pool
	expire int32
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                    c,
		client:               httpx.NewClient(conf.Conf.HTTPTelecom),
		payInfoURL:           c.Host.Telecom + _payInfo,
		cancelRepeatOrderURL: c.Host.Telecom + _cancelRepeatOrder,
		sucOrderListURL:      c.Host.Telecom + _sucOrderList,
		phoneAreaURL:         c.Host.Telecom + _phoneArea,
		orderStateURL:        c.Host.Telecom + _orderState,
		telecomReturnURL:     c.Host.TelecomReturnURL,
		telecomCancelPayURL:  c.Host.TelecomCancelPayURL,
		smsSendURL:           c.Host.Sms + _smsSendURL,
		db:                   xsql.NewMySQL(c.MySQL.Show),
		phoneRds:             redis.NewPool(c.Redis.Recommend.Config),
		//reids
		phoneKeyExpired: int32(time.Duration(c.Telecom.KeyExpired) / time.Second),
		payKeyExpired:   int32(time.Duration(c.Telecom.PayKeyExpired) / time.Second),
		// memcache
		mc:     memcache.NewPool(c.Memcache.Operator.Config),
		expire: int32(time.Duration(c.Memcache.Operator.Expire) / time.Second),
	}
	d.inOrderSyncSQL = d.db.Prepared(_inOrderSyncSQL)
	d.inRechargeSyncSQL = d.db.Prepared(_inRechargeSyncSQL)
	d.orderByPhoneSQL = d.db.Prepared(_orderByPhoneSQL)
	d.orderByOrderIDSQL = d.db.Prepared(_orderByOrderIDSQL)
	return
}

// InOrderSync
func (d *Dao) InOrderSync(ctx context.Context, requestNo, resultType int, phone string, t *telecom.TelecomJSON) (row int64, err error) {
	res, err := d.inOrderSyncSQL.Exec(ctx, requestNo, resultType, t.FlowpackageID, t.FlowPackageSize, t.FlowPackageType, t.TrafficAttribution, t.BeginTime, t.EndTime,
		t.IsMultiplyOrder, t.SettlementType, t.Operator, t.OrderStatus, t.RemainedRebindNum, t.MaxbindNum, t.OrderID, t.SignNo, t.AccessToken,
		phone, t.IsRepeatOrder, t.PayStatus, t.PayTime, t.PayChannel, t.SignStatus, t.RefundStatus,
		requestNo, resultType, t.FlowPackageSize, t.FlowPackageType, t.TrafficAttribution, t.BeginTime, t.EndTime,
		t.IsMultiplyOrder, t.SettlementType, t.Operator, t.OrderStatus, t.RemainedRebindNum, t.MaxbindNum, t.OrderID,
		t.SignNo, t.AccessToken, t.IsRepeatOrder, t.PayStatus, t.PayTime, t.PayChannel, t.SignStatus, t.RefundStatus)
	if err != nil {
		log.Error("d.inOrderSyncSQL.Exec error(%v)", err)
		return
	}
	tmp := &telecom.OrderInfo{}
	tmp.OrderInfoJSONChange(t)
	phoneInt, _ := strconv.Atoi(t.PhoneID)
	if err = d.AddTelecomCache(ctx, phoneInt, tmp); err != nil {
		log.Error("s.AddTelecomCache error(%v)", err)
	}
	orderID, _ := strconv.ParseInt(t.OrderID, 10, 64)
	if err = d.AddTelecomOrderIDCache(ctx, orderID, tmp); err != nil {
		log.Error("s.AddTelecomOrderIDCache error(%v)", err)
	}
	return res.RowsAffected()
}

// InRechargeSync
func (d *Dao) InRechargeSync(ctx context.Context, r *telecom.RechargeJSON) (row int64, err error) {
	res, err := d.inRechargeSyncSQL.Exec(ctx, r.RequestNo, r.FcRechargeNo, r.RechargeStatus, r.OrderTotalSize, r.FlowBalance)
	if err != nil {
		log.Error("d.inRechargeSyncSQL.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) OrdersUserFlow(ctx context.Context, phoneID int) (res map[int]*telecom.OrderInfo, err error) {
	res = map[int]*telecom.OrderInfo{}
	var (
		PhoneIDStr string
		OrderIDStr string
	)
	t := &telecom.OrderInfo{}
	row := d.orderByPhoneSQL.QueryRow(ctx, phoneID)
	if err = row.Scan(&PhoneIDStr, &OrderIDStr, &t.OrderState, &t.SignNo, &t.IsRepeatorder, &t.Begintime, &t.Endtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("OrdersUserFlow row.Scan err (%v)", err)
		}
		return
	}
	t.TelecomChange()
	t.PhoneID, _ = strconv.Atoi(PhoneIDStr)
	t.OrderID, _ = strconv.ParseInt(OrderIDStr, 10, 64)
	if t.PhoneID > 0 {
		res[t.PhoneID] = t
	}
	return
}

func (d *Dao) OrdersUserByOrderID(ctx context.Context, orderID int64) (res map[int64]*telecom.OrderInfo, err error) {
	res = map[int64]*telecom.OrderInfo{}
	var (
		PhoneIDStr string
		OrderIDStr string
	)
	t := &telecom.OrderInfo{}
	row := d.orderByOrderIDSQL.QueryRow(ctx, orderID)
	if err = row.Scan(&PhoneIDStr, &OrderIDStr, &t.OrderState, &t.SignNo, &t.IsRepeatorder, &t.Begintime, &t.Endtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("OrdersUserFlow row.Scan err (%v)", err)
		}
		return
	}
	t.TelecomChange()
	t.PhoneID, _ = strconv.Atoi(PhoneIDStr)
	t.OrderID, _ = strconv.ParseInt(OrderIDStr, 10, 64)
	if t.OrderID > 0 {
		res[t.OrderID] = t
	}
	return
}
