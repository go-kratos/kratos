package mobile

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/mobile"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inOrderSyncSQL = `INSERT IGNORE INTO mobile_order (orderid,userpseudocode,channelseqid,price,actiontime,actionid,effectivetime,expiretime,channelid,productid,ordertype,threshold)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE orderid=?,channelseqid=?,price=?,actiontime=?,actionid=?,effectivetime=?,expiretime=?,channelid=?,productid=?,ordertype=?,threshold=?`
	_upOrderFlowSQL     = `UPDATE mobile_order SET threshold=?,resulttime=? WHERE userpseudocode=? AND productid=?`
	_orderSyncByUserSQL = `SELECT orderid,userpseudocode,channelseqid,price,actionid,effectivetime,expiretime,channelid,productid,ordertype,threshold FROM mobile_order WHERE effectivetime<? AND userpseudocode=?`
)

type Dao struct {
	db               *xsql.DB
	inOrderSyncSQL   *xsql.Stmt
	upOrderFlowSQL   *xsql.Stmt
	orderUserSyncSQL *xsql.Stmt
	// memcache
	mc     *memcache.Pool
	expire int32
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: xsql.NewMySQL(c.MySQL.Show),
		// memcache
		mc:     memcache.NewPool(c.Memcache.Operator.Config),
		expire: int32(time.Duration(c.Memcache.Operator.Expire) / time.Second),
	}
	d.inOrderSyncSQL = d.db.Prepared(_inOrderSyncSQL)
	d.upOrderFlowSQL = d.db.Prepared(_upOrderFlowSQL)
	d.orderUserSyncSQL = d.db.Prepared(_orderSyncByUserSQL)
	return
}

func (d *Dao) InOrdersSync(ctx context.Context, u *mobile.MobileXML) (row int64, err error) {
	res, err := d.inOrderSyncSQL.Exec(ctx, u.Orderid, u.Userpseudocode, u.Channelseqid, u.Price, u.Actiontime, u.Actionid, u.Effectivetime, u.Expiretime,
		u.Channelid, u.Productid, u.Ordertype, u.Threshold,
		u.Orderid, u.Channelseqid, u.Price, u.Actiontime, u.Actionid, u.Effectivetime, u.Expiretime,
		u.Channelid, u.Productid, u.Ordertype, u.Threshold)
	if err != nil {
		log.Error("d.inOrderSyncSQL.Exec error(%v)", err)
		return
	}
	tmp := &mobile.Mobile{}
	tmp.MobileXMLMobile(u)
	if err = d.UpdateMobileCache(ctx, u.Userpseudocode, tmp); err != nil {
		log.Error("mobile d.UpdateMobileCache usermob(%v) error(%v)", u.Userpseudocode, err)
		return
	}
	return res.RowsAffected()
}

// FlowSync update OrdersSync
func (d *Dao) FlowSync(ctx context.Context, u *mobile.MobileXML) (row int64, err error) {
	res, err := d.upOrderFlowSQL.Exec(ctx, u.Threshold, u.Resulttime, u.Userpseudocode, u.Productid)
	if err != nil {
		log.Error("d.upOrderFlowSQL.Exec error(%v)", err)
		return
	}
	thresholdInt, _ := strconv.Atoi(u.Threshold)
	tmp := &mobile.Mobile{
		Threshold: thresholdInt,
		Productid: u.Productid,
	}
	if err = d.UpdateMobileFlowCache(ctx, u.Userpseudocode, tmp); err != nil {
		log.Error("mobile d.UpdateMobileFlowCache usermob(%v) error(%v)", u.Userpseudocode, err)
		return
	}
	return res.RowsAffected()
}

// OrdersUserFlow select user OrdersSync
func (d *Dao) OrdersUserFlow(ctx context.Context, usermob string, now time.Time) (res map[string][]*mobile.Mobile, err error) {
	rows, err := d.orderUserSyncSQL.Query(ctx, now, usermob)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	res = map[string][]*mobile.Mobile{}
	for rows.Next() {
		u := &mobile.Mobile{}
		if err = rows.Scan(&u.Orderid, &u.Userpseudocode, &u.Channelseqid, &u.Price, &u.Actionid, &u.Effectivetime, &u.Expiretime,
			&u.Channelid, &u.Productid, &u.Ordertype, &u.Threshold); err != nil {
			log.Error("Mobile OrdersUserFlow rows.Scan err (%v)", err)
			return
		}
		u.MobileChange()
		res[u.Userpseudocode] = append(res[u.Userpseudocode], u)
	}
	return
}
