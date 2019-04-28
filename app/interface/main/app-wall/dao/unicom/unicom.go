package unicom

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	"go-common/library/database/elastic"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	//unicom
	_inOrderSyncSQL = `INSERT IGNORE INTO unicom_order (usermob,cpid,spid,type,ordertime,canceltime,endtime,channelcode,province,area,ordertype,videoid,ctime,mtime) 
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE cpid=?,spid=?,type=?,ordertime=?,canceltime=?,endtime=?,channelcode=?,province=?,area=?,ordertype=?,videoid=?,mtime=?`
	_inAdvanceSyncSQL = `INSERT IGNORE INTO unicom_order_advance (usermob,userphone,cpid,spid,ordertime,channelcode,province,area,ctime,mtime) 
	VALUES(?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE cpid=?,spid=?,ordertime=?,channelcode=?,province=?,area=?,mtime=?`
	_upOrderFlowSQL   = `UPDATE unicom_order SET time=?,flowbyte=?,mtime=? WHERE usermob=?`
	_orderUserSyncSQL = `SELECT usermob,cpid,spid,type,ordertime,canceltime,endtime,channelcode,province,area,ordertype,videoid,time,flowbyte FROM unicom_order WHERE usermob=? 
	ORDER BY type DESC`
	_inIPSyncSQL = `INSERT IGNORE INTO unicom_ip (ipbegion,ipend,provinces,isopen,opertime,sign,ctime,mtime) VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		ipbegion=?,ipend=?,provinces=?,isopen=?,opertime=?,sign=?,mtime=?`
	_ipSyncSQL = `SELECT ipbegion,ipend FROM unicom_ip WHERE isopen=1`
	//pack
	_inPackSQL = `INSERT IGNORE INTO unicom_pack (usermob,mid) VALUES(?,?)`
	_packSQL   = `SELECT usermob,mid FROM unicom_pack WHERE usermob=?`
	// update unicom ip
	_inUnicomIPSyncSQL = `INSERT IGNORE INTO unicom_ip (ipbegion,ipend,isopen,ctime,mtime) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		ipbegion=?,ipend=?,isopen=?,mtime=?`
	_upUnicomIPSQL = `UPDATE unicom_ip SET isopen=?,mtime=? WHERE ipbegion=? AND ipend=?`
	// unicom integral change
	_inUserBindSQL = `INSERT IGNORE INTO unicom_user_bind (usermob,phone,mid,state,integral,flow) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		phone=?,state=?`
	_userBindSQL         = `SELECT usermob,phone,mid,state,integral,flow,monthlytime FROM unicom_user_bind WHERE state=1 AND mid=?`
	_userBindPhoneMidSQL = `SELECT mid FROM unicom_user_bind WHERE phone=? AND state=1`
	_upUserIntegralSQL   = `UPDATE unicom_user_bind SET integral=?,flow=? WHERE phone=? AND mid=?`
	_userPacksSQL        = `SELECT id,ptype,pdesc,amount,capped,integral,param FROM unicom_user_packs WHERE state=1`
	_userPacksByIDSQL    = `SELECT id,ptype,pdesc,amount,capped,integral,param,state FROM unicom_user_packs WHERE id=?`
	_upUserPacksSQL      = `UPDATE unicom_user_packs SET amount=?,state=? WHERE id=?`
	_inUserPackLogSQL    = `INSERT INTO unicom_user_packs_log (phone,usermob,mid,request_no,ptype,integral,pdesc) VALUES (?,?,?,?,?,?,?)`
	_userBindOldSQL      = `SELECT usermob,phone,mid,state,integral,flow FROM unicom_user_bind WHERE phone=? ORDER BY mtime DESC limit 1`
	_userPacksLogSQL     = `SELECT phone,integral,ptype,pdesc FROM unicom_user_packs_log WHERE mtime>=? AND mtime<?`
)

type Dao struct {
	db      *xsql.DB
	client  *httpx.Client
	uclient *httpx.Client
	//unicom
	inOrderSyncSQL   *xsql.Stmt
	inAdvanceSyncSQL *xsql.Stmt
	upOrderFlowSQL   *xsql.Stmt
	orderUserSyncSQL *xsql.Stmt
	inIPSyncSQL      *xsql.Stmt
	ipSyncSQL        *xsql.Stmt
	// unicom integral change
	inUserBindSQL       *xsql.Stmt
	userBindSQL         *xsql.Stmt
	userBindPhoneMidSQL *xsql.Stmt
	upUserIntegralSQL   *xsql.Stmt
	userPacksSQL        *xsql.Stmt
	userPacksByIDSQL    *xsql.Stmt
	upUserPacksSQL      *xsql.Stmt
	inUserPackLogSQL    *xsql.Stmt
	userBindOldSQL      *xsql.Stmt
	//pack
	inPackSQL *xsql.Stmt
	packSQL   *xsql.Stmt
	// memcache
	mc             *memcache.Pool
	expire         int32
	flowKeyExpired int32
	flowWait       int32
	// unicom url
	unicomIPURL           string
	unicomFlowExchangeURL string
	// order url
	orderURL       string
	ordercancelURL string
	sendsmscodeURL string
	smsNumberURL   string
	// elastic
	es *elastic.Elastic
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:      xsql.NewMySQL(c.MySQL.Show),
		client:  httpx.NewClient(conf.Conf.HTTPBroadband),
		uclient: httpx.NewClient(conf.Conf.HTTPUnicom),
		// unicom url
		unicomIPURL:           c.Host.Unicom + _unicomIPURL,
		unicomFlowExchangeURL: c.Host.UnicomFlow + _unicomFlowExchangeURL,
		// memcache
		mc:             memcache.NewPool(c.Memcache.Operator.Config),
		expire:         int32(time.Duration(c.Memcache.Operator.Expire) / time.Second),
		flowKeyExpired: int32(time.Duration(c.Unicom.KeyExpired) / time.Second),
		flowWait:       int32(time.Duration(c.Unicom.FlowWait) / time.Second),
		// order url
		orderURL:       c.Host.Broadband + _orderURL,
		ordercancelURL: c.Host.Broadband + _ordercancelURL,
		sendsmscodeURL: c.Host.Broadband + _sendsmscodeURL,
		smsNumberURL:   c.Host.Broadband + _smsNumberURL,
		// elastic
		es: elastic.NewElastic(nil),
	}
	d.inOrderSyncSQL = d.db.Prepared(_inOrderSyncSQL)
	d.inAdvanceSyncSQL = d.db.Prepared(_inAdvanceSyncSQL)
	d.upOrderFlowSQL = d.db.Prepared(_upOrderFlowSQL)
	d.orderUserSyncSQL = d.db.Prepared(_orderUserSyncSQL)
	d.inIPSyncSQL = d.db.Prepared(_inIPSyncSQL)
	d.ipSyncSQL = d.db.Prepared(_ipSyncSQL)
	// unicom integral change
	d.inUserBindSQL = d.db.Prepared(_inUserBindSQL)
	d.userBindSQL = d.db.Prepared(_userBindSQL)
	d.userBindPhoneMidSQL = d.db.Prepared(_userBindPhoneMidSQL)
	d.upUserIntegralSQL = d.db.Prepared(_upUserIntegralSQL)
	d.userPacksSQL = d.db.Prepared(_userPacksSQL)
	d.upUserPacksSQL = d.db.Prepared(_upUserPacksSQL)
	d.userPacksByIDSQL = d.db.Prepared(_userPacksByIDSQL)
	d.inUserPackLogSQL = d.db.Prepared(_inUserPackLogSQL)
	d.userBindOldSQL = d.db.Prepared(_userBindOldSQL)
	//pack
	d.inPackSQL = d.db.Prepared(_inPackSQL)
	d.packSQL = d.db.Prepared(_packSQL)
	return
}

// InOrdersSync insert OrdersSync
func (d *Dao) InOrdersSync(ctx context.Context, usermob string, u *unicom.UnicomJson, now time.Time) (row int64, err error) {
	res, err := d.inOrderSyncSQL.Exec(ctx, usermob,
		u.Cpid, u.Spid, u.TypeInt, u.Ordertime, u.Canceltime, u.Endtime,
		u.Channelcode, u.Province, u.Area, u.Ordertypes, u.Videoid, now, now,
		u.Cpid, u.Spid, u.TypeInt, u.Ordertime, u.Canceltime, u.Endtime,
		u.Channelcode, u.Province, u.Area, u.Ordertypes, u.Videoid, now)
	if err != nil {
		log.Error("d.inOrderSyncSQL.Exec error(%v)", err)
		return
	}
	utmp := &unicom.Unicom{}
	utmp.UnicomJSONTOUincom(usermob, u)
	if err = d.UpdateUnicomCache(ctx, usermob, utmp); err != nil {
		log.Error("d.UpdateUnicomCache usermob(%v) error(%v)", usermob, err)
		return 0, err
	}
	return res.RowsAffected()
}

// InAdvanceSync insert AdvanceSync
func (d *Dao) InAdvanceSync(ctx context.Context, usermob string, u *unicom.UnicomJson, now time.Time) (row int64, err error) {
	res, err := d.inAdvanceSyncSQL.Exec(ctx, usermob, u.Userphone,
		u.Cpid, u.Spid, u.Ordertime, u.Channelcode, u.Province, u.Area, now, now,
		u.Cpid, u.Spid, u.Ordertime, u.Channelcode, u.Province, u.Area, now)
	if err != nil {
		log.Error("d.inAdvanceSyncSQL.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// FlowSync update OrdersSync
func (d *Dao) FlowSync(ctx context.Context, flowbyte int, usermob, time string, now time.Time) (row int64, err error) {
	res, err := d.upOrderFlowSQL.Exec(ctx, time, flowbyte, now, usermob)
	if err != nil {
		log.Error("d.upOrderFlowSQL.Exec error(%v)", err)
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
		if err = rows.Scan(&u.Usermob, &u.Cpid, &u.Spid, &u.TypeInt, &u.Ordertime, &u.Canceltime, &u.Endtime, &u.Channelcode, &u.Province,
			&u.Area, &u.Ordertypes, &u.Videoid, &u.Time, &u.Flowbyte); err != nil {
			log.Error("OrdersUserFlow row.Scan err (%v)", err)
			return
		}
		u.UnicomChange()
		res = append(res, u)
	}
	return
}

// InIPSync insert IpSync
func (d *Dao) InIPSync(ctx context.Context, u *unicom.UnicomIpJson, now time.Time) (row int64, err error) {
	res, err := d.inIPSyncSQL.Exec(ctx, u.Ipbegin, u.Ipend, u.Provinces, u.Isopen, u.Opertime, u.Sign, now, now,
		u.Ipbegin, u.Ipend, u.Provinces, u.Isopen, u.Opertime, u.Sign, now)
	if err != nil {
		log.Error("d.inIPSyncSQL.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// InPack insert Privilege pack
func (d *Dao) InPack(ctx context.Context, usermob string, mid int64) (row int64, err error) {
	res, err := d.inPackSQL.Exec(ctx, usermob, mid)
	if err != nil {
		log.Error("d.inPackSQL.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// Pack select Privilege pack
func (d *Dao) Pack(ctx context.Context, usermobStr string) (res map[string]map[int64]struct{}, err error) {
	row := d.packSQL.QueryRow(ctx, usermobStr)
	var (
		usermob string
		mid     int64
	)
	res = map[string]map[int64]struct{}{}
	if err = row.Scan(&usermob, &mid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("OrdersUserFlow rows.Scan err (%v)", err)
		}
	}
	if user, ok := res[usermob]; !ok {
		res[usermob] = map[int64]struct{}{
			mid: struct{}{},
		}
	} else {
		user[mid] = struct{}{}
	}
	return
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

// InUserBind unicom insert user bind
func (d *Dao) InUserBind(ctx context.Context, ub *unicom.UserBind) (row int64, err error) {
	res, err := d.inUserBindSQL.Exec(ctx, ub.Usermob, ub.Phone, ub.Mid, ub.State, ub.Integral, ub.Flow, ub.Phone, ub.State)
	if err != nil {
		log.Error("insert user bind error(%v)", err)
		return
	}
	return res.RowsAffected()
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

// UserPacks user pack list
func (d *Dao) UserPacks(ctx context.Context) (res []*unicom.UserPack, err error) {
	rows, err := d.userPacksSQL.Query(ctx)
	if err != nil {
		log.Error("user pack sql error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &unicom.UserPack{}
		if err = rows.Scan(&u.ID, &u.Type, &u.Desc, &u.Amount, &u.Capped, &u.Integral, &u.Param); err != nil {
			log.Error("user pack sql error(%v)", err)
			return
		}
		res = append(res, u)
	}
	return
}

// UserPackByID user pack by id
func (d *Dao) UserPackByID(ctx context.Context, id int64) (res map[int64]*unicom.UserPack, err error) {
	res = map[int64]*unicom.UserPack{}
	row := d.userPacksByIDSQL.QueryRow(ctx, id)
	if row == nil {
		log.Error("user pack sql is null")
		return
	}
	u := &unicom.UserPack{}
	if err = row.Scan(&u.ID, &u.Type, &u.Desc, &u.Amount, &u.Capped, &u.Integral, &u.Param, &u.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("userPacksByIDSQL row.Scan error(%v)", err)
		}
		return
	}
	res[id] = u
	return
}

// UpUserPacks update user packs
func (d *Dao) UpUserPacks(ctx context.Context, u *unicom.UserPack, id int64) (row int64, err error) {
	res, err := d.upUserPacksSQL.Exec(ctx, u.Amount, u.State, id)
	if err != nil {
		log.Error("update user pack sql error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// UpUserIntegral update unicom user integral
func (d *Dao) UpUserIntegral(ctx context.Context, ub *unicom.UserBind) (row int64, err error) {
	res, err := d.upUserIntegralSQL.Exec(ctx, ub.Integral, ub.Flow, ub.Phone, ub.Mid)
	if err != nil {
		log.Error("update user integral sql error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// UserBindPhoneMid mid by phone
func (d *Dao) UserBindPhoneMid(ctx context.Context, phone string) (mid int64, err error) {
	row := d.userBindPhoneMidSQL.QueryRow(ctx, phone)
	if row == nil {
		log.Error("user pack sql is null")
		return
	}
	if err = row.Scan(&mid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("userPacksByIDSQL row.Scan error(%v)", err)
		}
		return
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

// UserBindOld user by phone
func (d *Dao) UserBindOld(ctx context.Context, phone string) (res *unicom.UserBind, err error) {
	row := d.userBindOldSQL.QueryRow(ctx, phone)
	if row == nil {
		log.Error("user bind old sql is null")
		return
	}
	res = &unicom.UserBind{}
	if err = row.Scan(&res.Usermob, &res.Phone, &res.Mid, &res.State, &res.Integral, &res.Flow); err != nil {
		log.Error("userBindSQL row.Scan error(%v)", err)
		res = nil
		return
	}
	return
}

// UserPacksLog user pack logs
func (d *Dao) UserPacksLog(ctx context.Context, start, end time.Time) (res []*unicom.UserPackLog, err error) {
	rows, err := d.db.Query(ctx, _userPacksLogSQL, start, end)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &unicom.UserPackLog{}
		if err = rows.Scan(&u.Phone, &u.Integral, &u.Type, &u.UserDesc); err != nil {
			log.Error("user packs log sql error(%v)", err)
			return
		}
		res = append(res, u)
	}
	err = rows.Err()
	return
}

// BeginTran begin a transacition
func (d *Dao) BeginTran(ctx context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(ctx)
}
