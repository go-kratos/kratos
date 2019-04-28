package vip

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/live/xuser/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
	xsql "go-common/library/database/sql"
)

const (
	_userVipRecordPrefix = "user_vip_record_%d"
	_userVipRecordCount  = 10
	_userLevelPrefix     = "user_"
)

var (
	errUpdateVipTimeInvalid = errors.New("update vip but vip_time invalid")
)

var (
	// get vip info from user_x table
	_getVipInfo = "SELECT `vip`,`vip_time`,`svip`,`svip_time` FROM `%s` WHERE uid=?;"

	// insert into user_vip_record_n
	_insertUserVipRecord = "INSERT INTO `%s` (`uid`,`vip_type`,`vip_num`,`order_id`,`platform`,`source`) VALUES (?,?,?,?,?,?);"

	// update user_vip_record_n before_time & status after update vip success
	_updateUserVipRecordStatus = "UPDATE `%s` SET `before_vip`=?,`before_vip_time`=?,`before_svip`=?,`before_svip_time`=?,`status`=? WHERE `id`=?;"

	// update user_x vip info
	_updateAllVip = "UPDATE `%s` SET `vip`=?,`vip_time`=?,`svip`=?,`svip_time`=? WHERE uid=?;"
	_updateVip    = "UPDATE `%s` SET `vip`=?,`vip_time`=? WHERE uid=?;"

	// insert vip
	_insertVip = "INSERT INTO `%s` (`uid`,`vip`,`vip_time`,`svip`,`svip_time`) VALUES (?,?,?,?,?);"

	// delete vip
	_deleteVip = "DELETE FROM `%s` WHERE `uid`=? LIMIT 1;"

	// insert ap_vip_record
	_insertApVipRecord = "INSERT INTO `ap_vip_record`(`uid`,`type`,`vip_time`,`platform`) VALUES (?,?,?,?);"
)

// GetVipFromDB get vip info by uid
// return error maybe sql no row or just scan error, how to handle decided by upper business
func (d *Dao) GetVipFromDB(ctx context.Context, uid int64) (info *model.VipInfo, err error) {
	var (
		vipTime, sVipTime xtime.Time
		currentTime       = xtime.Time(time.Now().Unix())
	)
	row := d.db.QueryRow(ctx, fmt.Sprintf(_getVipInfo, getUserLevelTable(uid)), uid)
	info = &model.VipInfo{}
	if err = row.Scan(&info.Vip, &vipTime, &info.Svip, &sVipTime); err != nil {
		log.Error("[dao.vip.mysql|GetVipFromDB] row scan error(%v), uid(%d)", err, uid)
		// no rows in user_x table, async insert one, don't return error
		if err == xsql.ErrNoRows {
			go d.createVip(context.TODO(), uid, info)
			err = nil
			return
		}
		return
	}

	// format info vip time
	if vipTime <= 0 {
		info.VipTime = model.TimeEmpty
	} else {
		info.VipTime = vipTime.Time().Format(model.TimeNano)
	}
	if sVipTime <= 0 {
		info.SvipTime = model.TimeEmpty
	} else {
		info.SvipTime = sVipTime.Time().Format(model.TimeNano)
	}

	// format vip & svip
	// 注意!!! db里的数据不一定正确，可能含有已经过期的time
	if vipTime <= currentTime {
		info.Vip = 0
	}
	if sVipTime <= currentTime {
		info.Svip = 0
	}

	return
}

// AddVip update user_n vip fields, add vip/svip time
// weather add vip or svip, vipTime should not be empty
func (d *Dao) AddVip(ctx context.Context, uid int64, vipTime, sVipTime xtime.Time) (row int64, err error) {
	var (
		vt, st      string
		res         sql.Result
		updateType  string
		currentTime = xtime.Time(time.Now().Unix())
	)
	if vipTime <= currentTime {
		return 0, errUpdateVipTimeInvalid
	}

	vt = vipTime.Time().Format(model.TimeNano)
	if sVipTime > currentTime {
		// update vip and svip
		st = sVipTime.Time().Format(model.TimeNano)
		updateType = "all"
		res, err = d.db.Exec(ctx, fmt.Sprintf(_updateAllVip, getUserLevelTable(uid)), 1, vt, 1, st, uid)
	} else {
		// update vip only
		updateType = "vip"
		res, err = d.db.Exec(ctx, fmt.Sprintf(_updateVip, getUserLevelTable(uid)), 1, vt, uid)
	}
	if err != nil {
		log.Error("[dao.vip.mysql|AddVip] update vip error(%v), type(%s), uid(%d), vip(%s), svip(%s)",
			err, updateType, uid, vt, st)
		return
	}
	row, _ = res.RowsAffected()
	return
}

// createVip create user_n vip row, for internal usage only. Do not use in business!
func (d *Dao) createVip(ctx context.Context, uid int64, info *model.VipInfo) (err error) {
	info = d.initInfo(info)
	log.Info("[dao.vip.mysql|createVip] create user_n row, uid(%d), info(%v)", uid, info)
	_, err = d.db.Exec(ctx, fmt.Sprintf(_insertVip, getUserLevelTable(uid)),
		uid, info.Vip, info.VipTime, info.Svip, info.SvipTime)
	if err != nil {
		log.Error("[dao.vip.mysql|createVip] create error(%v), uid(%d), info(%v)", err, uid, info)
	}
	return
}

// deleteVip delete user_n vip row, for internal usage only. Do not use in business!
func (d *Dao) deleteVip(ctx context.Context, uid int64) (err error) {
	log.Info("[dao.vip.mysql|deleteVip] delete user_n row, uid(%d)", uid)
	_, err = d.db.Exec(ctx, fmt.Sprintf(_deleteVip, getUserLevelTable(uid)), uid)
	if err != nil {
		log.Error("[dao.vip.mysql|deleteVip] delete error(%v), uid(%d)", err, uid)
	}
	return
}

// CreateVipRecord create user vip record if not exists
// return error maybe unique key exists err or other db error, upper business should notice
// unique key is (uid,order_id)
func (d *Dao) CreateVipRecord(ctx context.Context, req *model.VipBuy) (recordID int64, err error) {
	res, err := d.db.Exec(ctx, fmt.Sprintf(_insertUserVipRecord, getUserVipRecordTable(req.Uid)),
		req.Uid, req.GoodID, req.GoodNum, req.OrderID, req.Platform, req.Source)
	if err != nil {
		log.Error("[dao.vip.mysql|CreateUserVipRecord] create user vip record error(%v), req(%v)", err, req)
		return
	}
	if recordID, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		log.Error("[dao.vip.mysql|CreateUserVipRecord] get last insert id error(%v), req(%v)", err, req)
	}
	return
}

// UpdateVipRecord update user vip record after buy success
func (d *Dao) UpdateVipRecord(ctx context.Context, recordID, uid int64, info *model.VipInfo) (err error) {
	execSql := fmt.Sprintf(_updateUserVipRecordStatus, getUserVipRecordTable(uid))
	_, err = d.db.Exec(ctx, execSql, info.Vip, info.VipTime, info.Svip, info.SvipTime, model.BuyStatusSuccess, recordID)
	if err != nil {
		log.Error("[dao.vip.mysql|UpdateVipRecordLater] update error(%v), record id(%d), uid(%d), info(%v)",
			err, recordID, uid, info)
	}
	return
}

// CreateApVipRecord create ap_vip_record
func (d *Dao) CreateApVipRecord(ctx context.Context, record *model.VipRecord) (row int64, err error) {
	res, err := d.db.Exec(ctx, _insertApVipRecord, record.Uid, record.VipType, record.AfterVipTime, record.Platform)
	if err != nil {
		log.Error("[dao.vip.mysql|CreateApVipRecord] insert ap_vip_record error(%v), record(%v)", err, record)
		return
	}
	row, _ = res.RowsAffected()
	return
}

// getUserLevelTable get user_x table by uid
func getUserLevelTable(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(uidStr))
	cipher := md5Ctx.Sum(nil)
	return _userLevelPrefix + hex.EncodeToString(cipher)[0:1]
}

// getUserVipRecordTable get user_vip_record_x table by uid
func getUserVipRecordTable(uid int64) string {
	return fmt.Sprintf(_userVipRecordPrefix, uid%_userVipRecordCount)
}
