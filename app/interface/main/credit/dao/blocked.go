package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_addBlockedInfoSQL = `INSERT INTO blocked_info(uid,origin_title,blocked_remark,origin_url,origin_content,origin_content_modify,origin_type,
		punish_time,punish_type,blocked_days,publish_status,blocked_type,blocked_forever,reason_type,oper_id,moral_num,operator_name) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_addBatchBlockedInfoSQL = `INSERT INTO blocked_info(uid,origin_title,blocked_remark,origin_url,origin_content,origin_content_modify,origin_type,
		punish_time,punish_type,blocked_days,publish_status,blocked_type,blocked_forever,reason_type,oper_id,moral_num,operator_name) VALUES %s`
	_blockedCountSQL      = `SELECT COUNT(*) FROM blocked_info WHERE uid=? AND punish_type IN (2,3) AND status = 0`
	_blockedNumUserSQL    = `SELECT COUNT(*) FROM blocked_info WHERE uid = ? AND status = 0`
	_blkHistoryCountSQL   = "SELECT COUNT(*) FROM blocked_info WHERE uid = ? AND ctime >= ? AND status = 0"
	_blockedTotalSQL      = "SELECT COUNT(*) AS num FROM blocked_info WHERE uid=? AND ctime >? AND status = 0"
	_blockedInfosByMidSQL = `SELECT id,case_id,uid,origin_title,origin_url,origin_content,origin_content_modify,origin_type,punish_time,punish_type,blocked_days,publish_status,blocked_type,
	reason_type,blocked_remark,ctime from blocked_info WHERE uid=? AND status = 0`
	_blockedListSQL = `SELECT id,origin_type,blocked_type,publish_time FROM blocked_info WHERE publish_status = 1 %s %s AND status = 0 ORDER BY publish_time desc`
	_blkHistorysSQL = `SELECT id,uid,blocked_days,blocked_forever,blocked_remark,moral_num,origin_content_modify,origin_title,origin_type,origin_url,punish_time,
	punish_type,reason_type FROM blocked_info WHERE uid = ? AND ctime >= ? AND status = 0 ORDER BY id LIMIT ?,?`
	_blockedInfoIDSQL = `SELECT id,uid,uname,origin_content,origin_content_modify,origin_type,punish_time,punish_type,moral_num,blocked_days,reason_type,blocked_forever,origin_title,
	origin_url,blocked_type,blocked_remark,case_id,ctime,publish_status from blocked_info WHERE id=? AND status = 0`
	_blockedInfoIDsSQL = `SELECT id,uid,blocked_days,blocked_forever,blocked_remark,moral_num,origin_content_modify,origin_title,origin_type,origin_url,punish_time,
	punish_type,reason_type FROM blocked_info WHERE id IN (%s) AND status = 0`
	_blockedInfosSQL = `SELECT id,uid,uname,origin_content_modify,origin_type,punish_time,punish_type,moral_num,blocked_days,reason_type,blocked_forever,origin_title,
	origin_url,blocked_type,blocked_remark,case_id,ctime,publish_status FROM blocked_info WHERE id IN (%s) AND publish_status = 1 AND status = 0 ORDER BY publish_time desc`
)

// AddBlockedInfo add blocked info
func (d *Dao) AddBlockedInfo(c context.Context, r *model.BlockedInfo) (err error) {
	if _, err = d.db.Exec(c, _addBlockedInfoSQL, r.UID, r.OriginTitle, r.BlockedRemark, r.OriginURL, r.OriginContent, r.OriginContent,
		r.OriginType, r.PunishTime.Time(), r.PunishType, r.BlockedDays, r.PublishStatus, r.BlockedType, r.BlockedForever,
		r.ReasonType, r.OID, r.MoralNum, r.OperatorName); err != nil {
		err = errors.Wrap(err, "AddBlockedInfo")
	}
	return
}

// TxAddBlockedInfo add blocked info
func (d *Dao) TxAddBlockedInfo(tx *sql.Tx, rs []*model.BlockedInfo) (err error) {
	l := len(rs)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*17)
	for _, v := range rs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, strconv.FormatInt(v.UID, 10))
		valueArgs = append(valueArgs, v.OriginTitle)
		valueArgs = append(valueArgs, v.BlockedRemark)
		valueArgs = append(valueArgs, v.OriginURL)
		valueArgs = append(valueArgs, v.OriginContent)
		valueArgs = append(valueArgs, v.OriginContent)
		valueArgs = append(valueArgs, strconv.FormatInt(v.OriginType, 10))
		valueArgs = append(valueArgs, v.PunishTime.Time())
		valueArgs = append(valueArgs, strconv.FormatInt(v.PunishType, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.BlockedDays, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.PublishStatus, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.BlockedType, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.BlockedForever, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.ReasonType, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.OID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(v.MoralNum, 10))
		valueArgs = append(valueArgs, v.OperatorName)
	}
	stmt := fmt.Sprintf(_addBatchBlockedInfoSQL, strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		err = errors.Wrapf(err, "TxAddBlockedInfo tx.Exec() error(%+v)", err)
	}
	return
}

// BlockedCount get user blocked count.
func (d *Dao) BlockedCount(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _blockedCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "BlockedCount scan fail")
	}
	return
}

// BlockedNumUser get blocked user number.
func (d *Dao) BlockedNumUser(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _blockedNumUserSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrap(err, "BlockedNumUser")
	}
	return
}

// BLKHistoryCount get blocked historys count.
func (d *Dao) BLKHistoryCount(c context.Context, ArgHis *model.ArgHistory) (count int64, err error) {
	row := d.db.QueryRow(c, _blkHistoryCountSQL, ArgHis.MID, xtime.Time(ArgHis.STime))
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrap(err, "BLKHistoryCount")
	}
	return
}

// BlockTotalTime get block total by time.
func (d *Dao) BlockTotalTime(c context.Context, mid int64, ts time.Time) (total int64, err error) {
	row := d.db.QueryRow(c, _blockedTotalSQL, mid, ts)
	if err = row.Scan(&total); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		err = nil
		total = 0
	}
	return
}

// BlockedUserList get user blocked list.
func (d *Dao) BlockedUserList(c context.Context, mid int64) (res []*model.BlockedInfo, err error) {
	rows, err := d.db.Query(c, _blockedInfosByMidSQL, mid)
	if err != nil {
		log.Error("d.getBlockedInfosByMidStmt.Query(mid %d) error(%v)", mid, err)
		err = errors.Wrap(err, "BlockedUserList")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.BlockedInfo{}
		if err = rows.Scan(&r.ID, &r.CaseID, &r.UID, &r.OriginTitle, &r.OriginURL, &r.OriginContent, &r.OriginContentModify, &r.OriginType, &r.PunishTime, &r.PunishType,
			&r.BlockedDays, &r.PublishStatus, &r.BlockedType, &r.ReasonType, &r.BlockedRemark, &r.CTime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BlockedList get blocked list.
func (d *Dao) BlockedList(c context.Context, otype, btype int8) (res []*model.BlockedInfo, err error) {
	var ostr, bstr string
	if otype != 0 {
		ostr = fmt.Sprintf("AND origin_type=%d ", otype)
	}
	if btype >= 0 {
		bstr = fmt.Sprintf("AND blocked_type=%d ", btype)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_blockedListSQL, ostr, bstr))
	if err != nil {
		err = errors.Wrap(err, "BlockedInfos")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.BlockedInfo{}
		if err = rows.Scan(&r.ID, &r.OriginType, &r.BlockedType, &r.PublishTime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BLKHistorys get blocked historys list.
func (d *Dao) BLKHistorys(c context.Context, ah *model.ArgHistory) (res []*model.BlockedInfo, err error) {
	rows, err := d.db.Query(c, _blkHistorysSQL, ah.MID, xtime.Time(ah.STime), (ah.PN-1)*ah.PS, ah.PS)
	if err != nil {
		err = errors.Wrap(err, "BLKHistorys")
		return
	}
	defer rows.Close()
	for rows.Next() {
		bi := new(model.BlockedInfo)
		if err = rows.Scan(&bi.ID, &bi.UID, &bi.BlockedDays, &bi.BlockedForever, &bi.BlockedRemark, &bi.MoralNum, &bi.OriginContentModify, &bi.OriginTitle,
			&bi.OriginType, &bi.OriginURL, &bi.PunishTime, &bi.PunishType, &bi.ReasonType); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.Wrap(err, "BLKHistorys")
			return
		}
		res = append(res, bi)
	}
	err = rows.Err()
	return
}

// BlockedInfoByID get blocked info by id.
func (d *Dao) BlockedInfoByID(c context.Context, id int64) (r *model.BlockedInfo, err error) {
	row := d.db.QueryRow(c, _blockedInfoIDSQL, id)
	r = new(model.BlockedInfo)
	if err = row.Scan(&r.ID, &r.UID, &r.Uname, &r.OriginContent, &r.OriginContentModify, &r.OriginType, &r.PunishTime, &r.PunishType, &r.MoralNum,
		&r.BlockedDays, &r.ReasonType, &r.BlockedForever, &r.OriginTitle, &r.OriginURL, &r.BlockedType, &r.BlockedRemark, &r.CaseID, &r.CTime, &r.PublishStatus); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.Wrap(err, "BlockedInfoByID")
	}
	return
}

// BlockedInfoIDs get blocked info by ids
func (d *Dao) BlockedInfoIDs(c context.Context, ids []int64) (res map[int64]*model.BlockedInfo, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_blockedInfoIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.Wrap(err, "BlockedInfoIDs")
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.BlockedInfo, len(ids))
	for rows.Next() {
		bi := new(model.BlockedInfo)
		if err = rows.Scan(&bi.ID, &bi.UID, &bi.BlockedDays, &bi.BlockedForever, &bi.BlockedRemark, &bi.MoralNum, &bi.OriginContentModify, &bi.OriginTitle,
			&bi.OriginType, &bi.OriginURL, &bi.PunishTime, &bi.PunishType, &bi.ReasonType); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.Wrap(err, "BlockedInfoIDs")
			return
		}
		res[bi.ID] = bi
	}
	err = rows.Err()
	return
}

// BlockedInfos get blocked infos. Queryed without mid or id, public default.
func (d *Dao) BlockedInfos(c context.Context, ids []int64) (res []*model.BlockedInfo, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_blockedInfosSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.Wrap(err, "BlockedInfos")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.BlockedInfo{}
		if err = rows.Scan(&r.ID, &r.UID, &r.Uname, &r.OriginContentModify, &r.OriginType, &r.PunishTime, &r.PunishType, &r.MoralNum,
			&r.BlockedDays, &r.ReasonType, &r.BlockedForever, &r.OriginTitle, &r.OriginURL, &r.BlockedType, &r.BlockedRemark, &r.CaseID, &r.CTime, &r.PublishStatus); err != nil {
			log.Error("BlockedInfos err %v", err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
