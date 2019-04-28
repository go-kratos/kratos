package unicom

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/job/main/app-wall/conf"
	"go-common/app/job/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	// unicom integral change
	_upUserIntegralSQL = `UPDATE unicom_user_bind SET integral=?,flow=?,monthlytime=? WHERE mid=? AND state=1`
	_orderUserSyncSQL  = `SELECT usermob,spid,type,ordertime,endtime FROM unicom_order WHERE usermob=? AND type=0 ORDER BY type DESC`
	_bindAllSQL        = `SELECT mid,usermob,monthlytime FROM unicom_user_bind WHERE state=1 LIMIT ?,?`
	_userBindSQL       = `SELECT usermob,phone,mid,state,integral,flow,monthlytime FROM unicom_user_bind WHERE state=1 AND mid=?`
	// update unicom ip
	_inUnicomIPSyncSQL = `INSERT IGNORE INTO unicom_ip (ipbegion,ipend,isopen,ctime,mtime) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		ipbegion=?,ipend=?,isopen=?,mtime=?`
	_upUnicomIPSQL        = `UPDATE unicom_ip SET isopen=?,mtime=? WHERE ipbegion=? AND ipend=?`
	_ipSyncSQL            = `SELECT ipbegion,ipend FROM unicom_ip WHERE isopen=1`
	_inUserPackLogSQL     = `INSERT INTO unicom_user_packs_log (phone,usermob,mid,request_no,ptype,integral,pdesc) VALUES (?,?,?,?,?,?,?)`
	_inUserIntegralLogSQL = `INSERT INTO unicom_user_integral_log (phone,mid,unicom_desc,ptype,integral,flow,pdesc) VALUES (?,?,?,?,?,?,?)`
)

type Dao struct {
	db      *xsql.DB
	uclient *httpx.Client
	// memcache
	mc             *memcache.Pool
	flowKeyExpired int32
	expire         int32
	// unicom integral change
	upUserIntegralSQL    *xsql.Stmt
	orderUserSyncSQL     *xsql.Stmt
	bindAllSQL           *xsql.Stmt
	userBindSQL          *xsql.Stmt
	ipSyncSQL            *xsql.Stmt
	inUserPackLogSQL     *xsql.Stmt
	inUserIntegralLogSQL *xsql.Stmt
	// unicom url
	unicomFlowExchangeURL string
	unicomIPURL           string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:      xsql.NewMySQL(c.MySQL.Show),
		uclient: httpx.NewClient(conf.Conf.HTTPUnicom),
		// memcache
		mc:             memcache.NewPool(c.Memcache.Operator.Config),
		expire:         int32(time.Duration(c.Unicom.PackKeyExpired) / time.Second),
		flowKeyExpired: int32(time.Duration(c.Unicom.KeyExpired) / time.Second),
		// unicom url
		unicomFlowExchangeURL: c.Host.UnicomFlow + _unicomFlowExchangeURL,
		unicomIPURL:           c.Host.Unicom + _unicomIPURL,
	}
	// unicom integral change
	d.upUserIntegralSQL = d.db.Prepared(_upUserIntegralSQL)
	d.orderUserSyncSQL = d.db.Prepared(_orderUserSyncSQL)
	d.bindAllSQL = d.db.Prepared(_bindAllSQL)
	d.userBindSQL = d.db.Prepared(_userBindSQL)
	d.ipSyncSQL = d.db.Prepared(_ipSyncSQL)
	d.inUserPackLogSQL = d.db.Prepared(_inUserPackLogSQL)
	d.inUserIntegralLogSQL = d.db.Prepared(_inUserIntegralLogSQL)
	return
}

// UpUserIntegral update unicom user integral
func (d *Dao) UpUserIntegral(ctx context.Context, ub *unicom.UserBind) (row int64, err error) {
	res, err := d.upUserIntegralSQL.Exec(ctx, ub.Integral, ub.Flow, ub.Monthly, ub.Mid)
	if err != nil {
		log.Error("update user integral sql error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// OrdersUserFlow select user OrdersSync
func (d *Dao) OrdersUserFlow(ctx context.Context, usermob string) (res []*unicom.Unicom, err error) {
	rows, err := d.orderUserSyncSQL.Query(ctx, usermob)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &unicom.Unicom{}
		if err = rows.Scan(&u.Usermob, &u.Spid, &u.TypeInt, &u.Ordertime, &u.Endtime); err != nil {
			log.Error("OrdersUserFlow row.Scan err (%v)", err)
			return
		}
		res = append(res, u)
	}
	return
}

//BindAll select bind all mid state 1
func (d *Dao) BindAll(ctx context.Context, start, end int) (res []*unicom.UserBind, err error) {
	rows, err := d.bindAllSQL.Query(ctx, start, end)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &unicom.UserBind{}
		if err = rows.Scan(&u.Mid, &u.Usermob, &u.Monthly); err != nil {
			log.Error("BindAll rows.Scan error(%v)", err)
			return
		}
		res = append(res, u)
	}
	return
}

// UserBind unicom select user bind
func (d *Dao) UserBind(ctx context.Context, mid int64) (res *unicom.UserBind, err error) {
	row := d.userBindSQL.QueryRow(ctx, mid)
	if row == nil {
		log.Error("userBindSQL is null")
		return
	}
	res = &unicom.UserBind{}
	if err = row.Scan(&res.Usermob, &res.Phone, &res.Mid, &res.State, &res.Integral, &res.Flow, &res.Monthly); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("userBindSQL row.Scan error(%v)", err)
		}
		res = nil
		return
	}
	return
}

// InUnicomIPSync insert or update unicom_ip
func (d *Dao) InUnicomIPSync(tx *xsql.Tx, u *unicom.UnicomIP, now time.Time) (row int64, err error) {
	res, err := tx.Exec(_inUnicomIPSyncSQL, u.Ipbegin, u.Ipend, 1, now, now,
		u.Ipbegin, u.Ipend, 1, now)
	if err != nil {
		log.Error("tx.inUnicomIPSyncSQL.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// UpUnicomIP update unicom_ip state
func (d *Dao) UpUnicomIP(tx *xsql.Tx, ipstart, ipend, state int, now time.Time) (row int64, err error) {
	res, err := tx.Exec(_upUnicomIPSQL, state, now, ipstart, ipend)
	if err != nil {
		log.Error("tx.upUnicomIPSQL.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// IPSync select all ipSync
func (d *Dao) IPSync(ctx context.Context) (res []*unicom.UnicomIP, err error) {
	rows, err := d.ipSyncSQL.Query(ctx)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	res = []*unicom.UnicomIP{}
	for rows.Next() {
		u := &unicom.UnicomIP{}
		if err = rows.Scan(&u.Ipbegin, &u.Ipend); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		u.UnicomIPChange()
		res = append(res, u)
	}
	return
}

// InUserPackLog insert unicom user pack log
func (d *Dao) InUserPackLog(ctx context.Context, u *unicom.UserPackLog) (row int64, err error) {
	res, err := d.inUserPackLogSQL.Exec(ctx, u.Phone, u.Usermob, u.Mid, u.RequestNo, u.Type, u.Integral, u.Desc)
	if err != nil {
		log.Error("insert user pack log integral sql error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// InUserIntegralLog insert unicom user add integral and flow log
func (d *Dao) InUserIntegralLog(ctx context.Context, u *unicom.UserIntegralLog) (row int64, err error) {
	res, err := d.inUserIntegralLogSQL.Exec(ctx, u.Phone, u.Mid, u.UnicomDesc, u.Type, u.Integral, u.Flow, u.Desc)
	if err != nil {
		log.Error("insert user add integral and flow sql error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// BeginTran begin a transacition
func (d *Dao) BeginTran(ctx context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(ctx)
}
