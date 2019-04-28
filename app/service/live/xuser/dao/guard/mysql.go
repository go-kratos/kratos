package guard

import (
	"context"
	"fmt"
	"go-common/library/xstr"
	"time"

	confm "go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/model"
	dhhm "go-common/app/service/live/xuser/model/dhh"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_guardTable = "ap_user_privilege"
)

var (
	// add guard info
	_addGuardInfo = "REPLACE INTO `%s` (`uid`,`target_id`,`privilege_type`,`start_time`,`expired_time`) VALUES(?,?,?,?,?);"
	// get guard info
	_getGuardInfo = "SELECT `id`,`uid`,`target_id`,`privilege_type`,`start_time`,`expired_time` FROM `%s` WHERE `uid`=? AND `expired_time`>=? ORDER BY `privilege_type` ASC;"
	// get guard info
	_getGuardInfo2 = "SELECT `id`,`uid`,`target_id`,`privilege_type`,`start_time`,`expired_time` FROM `%s` WHERE `uid`=? AND `target_id`=? AND `expired_time`>=? ORDER BY `privilege_type` ASC;"
	// update guard info
	_updGuardInfo = "UPDATE `%s` SET `expired_time`=date_add(expired_time, interval ? day) WHERE `uid`=? AND `target_id`=? AND `expired_time`>=? AND `privilege_type`%s?"
	// upsert guard info
	_upsertGuardInfo = "INSERT INTO `%s` (`uid`,`target_id`,`privilege_type`,`start_time`,`expired_time`) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE `start_time`=?,`expired_time`=?;"
	// 查询大航海信息
	_selUID           = "SELECT id,uid,target_id,privilege_type,start_time,expired_time,ctime,utime FROM ap_user_privilege where uid IN (%s) AND expired_time >= '%s' "
	_selAnchorUID     = "SELECT id,uid,target_id,privilege_type,start_time,expired_time,ctime,utime FROM ap_user_privilege where target_id IN (%s) AND expired_time >= '%s' "
	_errorDBLogPrefix = "xuser.dahanghai.dao.mysql|"
)

// GetByUIDs 批量查询
func (d *GuardDao) GetByUIDs(c context.Context, uids []int64) (dhhs []*dhhm.DHHDB, err error) {
	reqStartTime := confm.RecordTimeCost()
	dhhs = make([]*dhhm.DHHDB, 0)
	tm := time.Now()
	timeNow := tm.Format("2006-1-2 15:04:05")
	rows, err1 := d.db.Query(c, fmt.Sprintf(_selUID, xstr.JoinInts(uids), timeNow))
	if err1 != nil {
		reqAfterTime := confm.RecordTimeCost()
		err = err1
		log.Error(_errorDBLogPrefix+confm.GetFromDHHDBError+"|GetByUIDs err: %v|cost:%dms", err, reqAfterTime-reqStartTime)
		return
	}

	for rows.Next() {
		ele := &dhhm.DHHDB{}
		if err = rows.Scan(&ele.ID, &ele.UID, &ele.TargetId, &ele.PrivilegeType, &ele.StartTime, &ele.ExpiredTime, &ele.Ctime, &ele.Utime); err != nil {
			log.Error(_errorDBLogPrefix+confm.ScanFromDHHDBError+"|GetByUIDs rows.Scan err: %v", err)
			return
		}
		dhhs = append(dhhs, ele)
	}
	return
}

// GetByUIDsWithMap 批量查询
func (d *GuardDao) GetByUIDsWithMap(c context.Context, uids []int64) (dhhs map[int64][]*dhhm.DHHDB, err error) {
	reqStartTime := confm.RecordTimeCost()
	dhhs = make(map[int64][]*dhhm.DHHDB)
	tm := time.Now()
	timeNow := tm.Format("2006-1-2 15:04:05")
	rows, err1 := d.db.Query(c, fmt.Sprintf(_selUID, xstr.JoinInts(uids), timeNow))
	if err1 != nil {
		reqAfterTime := confm.RecordTimeCost()
		err = err1
		log.Error(_errorDBLogPrefix+confm.GetFromDHHDBError+"|GetByUIDs err: %v|cost:%dms", err, reqAfterTime-reqStartTime)
		return
	}

	for rows.Next() {
		ele := &dhhm.DHHDB{}
		if err = rows.Scan(&ele.ID, &ele.UID, &ele.TargetId, &ele.PrivilegeType, &ele.StartTime, &ele.ExpiredTime, &ele.Ctime, &ele.Utime); err != nil {
			log.Error(_errorDBLogPrefix+confm.ScanFromDHHDBError+"|GetByUIDs rows.Scan err: %v", err)
			return
		}
		if _, exist := dhhs[ele.UID]; !exist {
			dhhs[ele.UID] = make([]*dhhm.DHHDB, 0)
		}
		dhhs[ele.UID] = append(dhhs[ele.UID], ele)
	}
	return
}

// GetByAnchorUIDs 批量查询
func (d *GuardDao) GetByAnchorUIDs(c context.Context, uids []int64) (dhhs []*dhhm.DHHDB, err error) {
	reqStartTime := confm.RecordTimeCost()
	dhhs = make([]*dhhm.DHHDB, 0)
	tm := time.Now()
	timeNow := tm.Format("2006-1-2 15:04:05")
	rows, err1 := d.db.Query(c, fmt.Sprintf(_selAnchorUID, xstr.JoinInts(uids), timeNow))
	if err1 != nil {
		reqAfterTime := confm.RecordTimeCost()
		err = err1
		log.Error(_errorDBLogPrefix+confm.GetFromDHHDBError+"|GetByUIDs err: %v|cost:%dms", err, reqAfterTime-reqStartTime)
		return
	}

	for rows.Next() {
		ele := &dhhm.DHHDB{}
		if err = rows.Scan(&ele.ID, &ele.UID, &ele.TargetId, &ele.PrivilegeType, &ele.StartTime, &ele.ExpiredTime, &ele.Ctime, &ele.Utime); err != nil {
			log.Error(_errorDBLogPrefix+confm.ScanFromDHHDBError+"|GetByUIDs rows.Scan err: %v", err)
			return
		}
		dhhs = append(dhhs, ele)
	}
	return
}

// GetGuardByUID get guard info by uid
func (d *GuardDao) GetGuardByUID(ctx context.Context, uid int64) (info []*model.GuardInfo, err error) {
	sql := fmt.Sprintf(_getGuardInfo, _guardTable)
	rows, err := d.db.Query(ctx, sql, uid, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error("[dao.guard.mysql|GetGuardByUID] get user guard record error(%v), uid(%d)", err, uid)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var inf model.GuardInfo
		err = rows.Scan(&inf.Id, &inf.Uid, &inf.TargetId, &inf.PrivilegeType, &inf.StartTime, &inf.ExpiredTime)
		if err != nil {
			log.Error("[dao.guard.mysql|GetGuardByUID] scan user guard record error(%v), uid(%d)", err, uid)
			return nil, err
		}
		info = append(info, &inf)
	}
	return
}

// GetGuardByUIDRuid get guard info by uid and ruid
func (d *GuardDao) GetGuardByUIDRuid(ctx context.Context, uid int64, ruid int64) (info []*model.GuardInfo, err error) {
	sql := fmt.Sprintf(_getGuardInfo2, _guardTable)
	rows, err := d.db.Query(ctx, sql, uid, ruid, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error("[dao.guard.mysql|GetGuardByUIDRuid] get user guard record error(%v), uid(%d), ruid(%d)", err, uid, ruid)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var inf model.GuardInfo
		err = rows.Scan(&inf.Id, &inf.Uid, &inf.TargetId, &inf.PrivilegeType, &inf.StartTime, &inf.ExpiredTime)
		if err != nil {
			log.Error("[dao.guard.mysql|GetGuardByUIDRuid] scan user guard record error(%v), uid(%d), ruid(%d)", err, uid, ruid)
			return nil, err
		}
		info = append(info, &inf)
	}
	return
}

// AddGuard insert guard
func (d *GuardDao) AddGuard(ctx context.Context, req *model.GuardBuy) (err error) {
	sql := fmt.Sprintf(_addGuardInfo, _guardTable)
	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), int(23), int(59), int(59), int(999), time.Local).AddDate(0, 0, req.Num*30)
	res, err := d.db.Exec(ctx, sql, req.Uid, req.Ruid, req.GuardLevel, now.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		// unique key exists error
		log.Error("[dao.guard.mysql|AddGuard] add user guard record error(%v), req(%v)", err, req)
		return
	}
	if _, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		log.Error("[dao.guard.mysql|AddGuard] get last insert id error(%v), req(%v)", err, req)
	}
	return
}

// UpdateGuard update guard info
func (d *GuardDao) UpdateGuard(ctx context.Context, req *model.GuardBuy, cond string) (err error) {
	sql := fmt.Sprintf(_updGuardInfo, _guardTable, cond)
	_, err = d.db.Exec(ctx, sql, req.Num*30, req.Uid, req.Ruid, time.Now().Format("2006-01-02 15:04:05"), req.GuardLevel)
	if err != nil {
		log.Error("[dao.guard.mysql|UpdateGuard] update user guard record error(%v), req(%v)", err, req)
		return
	}
	return
}

// UpsertGuard upsert guard info
func (d *GuardDao) UpsertGuard(ctx context.Context, req *model.GuardBuy, expiredTime string) (err error) {
	sql := fmt.Sprintf(_upsertGuardInfo, _guardTable)
	now := time.Now()
	_, err = d.db.Exec(ctx, sql, req.Uid, req.Ruid, req.GuardLevel, now.Format("2006-01-02 15:04:05"), expiredTime, now.Format("2006-01-02 15:04:05"), expiredTime)
	if err != nil {
		log.Error("[dao.guard.mysql|UpsertGuard] upsert user guard record error(%v), req(%v)", err, req)
		return
	}
	return
}
