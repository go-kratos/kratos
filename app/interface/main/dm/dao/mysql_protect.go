package dao

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-common/app/interface/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// dm protect apply
	_addProtect      = "INSERT INTO dm_protect_apply (cid, uid, apply_uid, aid, playtime, dmid, msg, status, ctime, mtime) VALUES %s"
	_selProtect      = "SELECT ctime FROM dm_protect_apply WHERE dmid=? ORDER BY id DESC"
	_getPrttApply    = "SELECT id,cid,apply_uid,aid,playtime,msg,ctime from dm_protect_apply where uid=? %s and status=-1"
	_getProtectAids  = "SELECT aid FROM dm_protect_apply WHERE uid=? GROUP BY aid"
	_paStatus        = "UPDATE dm_protect_apply SET status=%d WHERE uid=%d AND id IN (%s)"
	_chgPaSwitch     = "INSERT INTO dm_protect_up (uid, notice_switch) VALUES(?, ?) ON DUPLICATE KEY UPDATE notice_switch=?"
	_getProtectByIDs = "SELECT cid,dmid FROM dm_protect_apply WHERE uid=%d AND id IN (%s)"
	_paUsrStat       = "SELECT aid,apply_uid,status,ctime from dm_protect_apply where ctime>?"
	_paStatistics    = "SELECT uid from dm_protect_apply where status=-1 GROUP BY uid"
	_paNoticeSwitch  = "SELECT uid FROM dm_protect_up WHERE uid IN (%s) AND notice_switch=0"
)

// AddProtectApply 添加保护弹幕申请
func (d *Dao) AddProtectApply(c context.Context, pas []*model.Pa) (affect int64, err error) {
	var (
		values, s string
	)
	for _, pa := range pas {
		pa.Msg = strings.Replace(pa.Msg, `\`, `\\`, -1)
		pa.Msg = strings.Replace(pa.Msg, `'`, `\'`, -1)
		values += fmt.Sprintf(`(%d, %d, %d, %d, %f, %d, '%s', %d, '%s', '%s'),`, pa.CID, pa.UID, pa.ApplyUID, pa.AID, pa.Playtime, pa.DMID, pa.Msg, pa.Status, pa.Ctime.Format("2006-01-02 15:04:05"), pa.Mtime.Format("2006-01-02 15:04:05"))
	}
	if len(values) > 0 {
		values = values[0:(len(values) - 1)]
	}
	s = fmt.Sprintf(_addProtect, values)
	res, err := d.biliDM.Exec(c, s)
	if err != nil {
		log.Error("d.biliDM.Exec(%v) error(%v)", s, err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// ProtectApplyTime 根据dmid获取保护弹幕申请
func (d *Dao) ProtectApplyTime(c context.Context, dmid int64) (t time.Time, err error) {
	row := d.biliDM.QueryRow(c, _selProtect, dmid)
	if err = row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ProtectApplies 保护弹幕申请列表
func (d *Dao) ProtectApplies(c context.Context, uid, aid int64, order string) (res []*model.Apply, err error) {
	var (
		query string
		rows  *sql.Rows
	)
	if aid > 0 {
		query = fmt.Sprintf(_getPrttApply, "AND aid=?")
		rows, err = d.biliDM.Query(c, query, uid, aid)
	} else {
		query = fmt.Sprintf(_getPrttApply, "")
		rows, err = d.biliDM.Query(c, query, uid)
	}
	if err != nil {
		log.Error("d.biliDM.Query(%s,%v,%v) error(%v)", query, uid, aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Apply{}
		t := time.Time{}
		if err = rows.Scan(&r.ID, &r.CID, &r.ApplyUID, &r.AID, &r.Playtime, &r.Msg, &t); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		r.Ctime = t.Format("2006-01-02 15:04:05")
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
		return
	}
	if order == "playtime" {
		sort.Sort(model.ApplySortPlaytime(res))
	} else {
		sort.Sort(model.ApplySortID(res))
	}
	return
}

// ProtectAids 被申请保护弹幕稿件列表
func (d *Dao) ProtectAids(c context.Context, uid int64) (res []int64, err error) {
	rows, err := d.biliDM.Query(c, _getProtectAids, uid)
	if err != nil {
		log.Error("d.biliDM.Query(%s,%d) error(%v)", _getProtectAids, uid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, aid)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UptPaStatus 修改保护弹幕状态
func (d *Dao) UptPaStatus(c context.Context, uid int64, ids string, status int) (affect int64, err error) {
	s := fmt.Sprintf(_paStatus, status, uid, ids)
	res, err := d.biliDM.Exec(c, s)
	if err != nil {
		log.Error("d.biliDM.Exec(%v) error(%v)", s, err)
	}
	affect, err = res.RowsAffected()
	return
}

// ProtectApplyByIDs get protect apply by dmid
func (d *Dao) ProtectApplyByIDs(c context.Context, uid int64, ids string) (res map[int64][]int64, err error) {
	res = make(map[int64][]int64, 10)
	s := fmt.Sprintf(_getProtectByIDs, uid, ids)
	rows, err := d.biliDM.Query(c, s)
	if err != nil {
		log.Error("d.biliDM.Query(%s) error(%v)", s, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			dmid int64
			cid  int64
		)
		if err = rows.Scan(&cid, &dmid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		_, ok := res[cid]
		if ok {
			res[cid] = append(res[cid], dmid)
		} else {
			res[cid] = []int64{dmid}
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UptPaNoticeSwitch 设置申请保护弹幕站内通知开关
func (d *Dao) UptPaNoticeSwitch(c context.Context, uid int64, status int) (affect int64, err error) {
	res, err := d.biliDM.Exec(c, _chgPaSwitch, uid, status, status)
	if err != nil {
		log.Error("d.biliDM.Exec(%s,%d,%d) error(%v)", _chgPaSwitch, status, uid, err)
	}
	affect, err = res.RowsAffected()
	return
}

// PaNoticeClose 获取关闭申请保护弹幕站内通知
func (d *Dao) PaNoticeClose(c context.Context, uids []int64) (res map[int64]bool, err error) {
	if len(uids) < 1 {
		return
	}
	res = make(map[int64]bool, len(uids))
	s := fmt.Sprintf(_paNoticeSwitch, xstr.JoinInts(uids))
	rows, err := d.biliDM.Query(c, s)
	if err != nil {
		log.Error("d.biliDM.Query(%s) error(%v)", s, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid int64
		if err = rows.Scan(&uid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[uid] = true
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return

}

// ProtectApplyStatistics 保护弹幕申请up统计
func (d *Dao) ProtectApplyStatistics(c context.Context) (res []int64, err error) {
	rows, err := d.biliDM.Query(c, _paStatistics)
	if err != nil {
		log.Error("d.biliDM.Query(%s) error(%v)", _paStatistics, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid int64
		if err = rows.Scan(&uid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, uid)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// twoDayAgo22 两天前的22点
func twoDayAgo22() string {
	yesterday := time.Now().Add(-48 * time.Hour)
	year, month, day := yesterday.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Format("2006-01-02") + " 22:00:00"
}

// PaUsrStat 保护弹幕申请用户统计
func (d *Dao) PaUsrStat(c context.Context) (res []*model.ApplyUserStat, err error) {
	rows, err := d.biliDM.Query(c, _paUsrStat, twoDayAgo22())
	if err != nil {
		log.Error("d.biliDM.Query(%s) error(%v)", _paUsrStat, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ApplyUserStat{}
		if err = rows.Scan(&r.Aid, &r.UID, &r.Status, &r.Ctime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}
