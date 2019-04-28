package offer

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	//offer
	_inClkSQL     = "INSERT IGNORE INTO wall_click (ip,cid,mac,idfa,cb,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
	_inActSQL     = "INSERT IGNORE INTO wall_active (ip,mid,rmac,mac,idfa,device,ctime,mtime) VALUES(?,?,?,?,?,?,?,?)"
	_upActIdfaSQL = "UPDATE wall_click SET idfa_active=?,mtime=? WHERE idfa=?"
	_cbSQL        = "SELECT cid,cb FROM wall_click WHERE (idfa=? OR idfa=?) and ctime>? and ctime<? ORDER BY id DESC LIMIT 1"
	_rMacCntSQL   = "SELECT COUNT(*) FROM wall_active WHERE rmac=?"
	_idfaSQL      = "SELECT id FROM wall_active WHERE idfa=?"
	_inANClkSQL   = "INSERT IGNORE INTO `wall_click_an` (`channel`,`imei`,`androidid`,`mac`,`ip`,`cb`,`ctime`,`mtime`) VALUES (?,?,?,?,?,?,?,?)"
	_anActiveSQL  = "SELECT COUNT(*) FROM `wall_click_an` WHERE `type`=1 AND ((`androidid`=? AND `androidid`<>'') OR ((`imei`=? OR `imei`=?) AND `imei`<>''))"
	_anCbSQL      = "SELECT `id`,`channel`,`cb`,`type` FROM `wall_click_an` WHERE (`androidid`=? AND `androidid`<>'') OR (`imei`=? AND `imei`<>'') ORDER BY `id` DESC LIMIT 1"
	_anGdtCbSQL   = "SELECT `id`,`channel`,`cb`,`type` FROM `wall_click_an` WHERE `imei`=? ORDER BY `id` DESC LIMIT 1"
	_upClkActSQL  = "UPDATE `wall_click_an` SET `type`=? WHERE `id`=?"
)

// Dao is wall dao.
type Dao struct {
	db           *xsql.DB
	inClkSQL     *xsql.Stmt
	inActSQL     *xsql.Stmt
	upActIdfaSQL *xsql.Stmt
	cbSQL        *xsql.Stmt
	rMacCntSQL   *xsql.Stmt
	idfaSQL      *xsql.Stmt
	// android
	inANClkSQL  *xsql.Stmt
	anActiveSQL *xsql.Stmt
	anCbSQL     *xsql.Stmt
	anGdtCbSQL  *xsql.Stmt
	upClkActSQL *xsql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: xsql.NewMySQL(c.MySQL.Show),
	}
	d.inClkSQL = d.db.Prepared(_inClkSQL)
	d.inActSQL = d.db.Prepared(_inActSQL)
	d.upActIdfaSQL = d.db.Prepared(_upActIdfaSQL)
	d.cbSQL = d.db.Prepared(_cbSQL)
	d.rMacCntSQL = d.db.Prepared(_rMacCntSQL)
	d.idfaSQL = d.db.Prepared(_idfaSQL)
	// android
	d.inANClkSQL = d.db.Prepared(_inANClkSQL)
	d.anActiveSQL = d.db.Prepared(_anActiveSQL)
	d.anCbSQL = d.db.Prepared(_anCbSQL)
	d.anGdtCbSQL = d.db.Prepared(_anGdtCbSQL)
	d.upClkActSQL = d.db.Prepared(_upClkActSQL)
	return
}

func (d *Dao) InClick(ctx context.Context, ip uint32, cid, mac, idfa, cb string, now time.Time) (row int64, err error) {
	res, err := d.inClkSQL.Exec(ctx, ip, cid, mac, idfa, cb, now, now)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) InActive(ctx context.Context, ip uint32, mid, rmac, mac, idfa, device string, now time.Time) (row int64, err error) {
	res, err := d.inActSQL.Exec(ctx, ip, mid, rmac, mac, idfa, device, now, now)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) UpIdfaActive(ctx context.Context, idfaAct, idfa string, now time.Time) (row int64, err error) {
	res, err := d.upActIdfaSQL.Exec(ctx, idfaAct, now, idfa)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) Callback(ctx context.Context, idfa, gdtIdfa string, now time.Time) (cid, cb string, err error) {
	lo := now.AddDate(0, 0, -1)
	row := d.cbSQL.QueryRow(ctx, idfa, gdtIdfa, lo, now)
	if row == nil {
		log.Error("d.cbSQL.QueryRow is null")
		return
	}
	if err = row.Scan(&cid, &cb); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Callback row.Scan error(%v)", err)
		}
		return
	}
	return
}

func (d *Dao) RMacCount(ctx context.Context, rmac string) (count int, err error) {
	row := d.rMacCntSQL.QueryRow(ctx, rmac)
	if row == nil {
		log.Error("d.rMacCntSQL.QueryRow is null")
		return
	}
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("RMacCount row.Scan error(%v)", err)
		}
		return
	}
	return
}

func (d *Dao) Exists(ctx context.Context, idfa string) (exist bool, err error) {
	row := d.idfaSQL.QueryRow(ctx, idfa)
	var id int
	if row == nil {
		log.Error("d.idfaSQL.QueryRow is null")
		return
	}
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("RMacCount row.Scan error(%v)", err)
		}
		return
	}
	exist = id > 0
	return
}

func (d *Dao) InANClick(c context.Context, channel, imei, androidid, mac, cb string, ip uint32, now time.Time) (row int64, err error) {
	res, err := d.inANClkSQL.Exec(c, channel, imei, androidid, mac, ip, cb, now, now)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (d *Dao) ANActive(c context.Context, androidid, imei, gdtImei string) (count int, err error) {
	row := d.anActiveSQL.QueryRow(c, androidid, imei, gdtImei)
	if row == nil {
		return
	}
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

func (d *Dao) ANCallback(c context.Context, androidid, imei, gdtImei string) (id int64, channel, cb string, typ int, err error) {
	row := d.anCbSQL.QueryRow(c, androidid, imei)
	if err = row.Scan(&id, &channel, &cb, &typ); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			return
		}
	}
	if gdtImei != "" {
		var (
			gdtID             int64
			gdtChannel, gdtCb string
			gdtTyp            int
		)
		row := d.anGdtCbSQL.QueryRow(c, gdtImei)
		if err = row.Scan(&gdtID, &gdtChannel, &gdtCb, &gdtTyp); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			return
		}
		if gdtID > id {
			id = gdtID
			channel = gdtChannel
			cb = gdtCb
			typ = gdtTyp
		}
	}
	return
}

func (d *Dao) ANClickAct(c context.Context, id int64, typ int) (row int64, err error) {
	res, err := d.upClkActSQL.Exec(c, typ, id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
