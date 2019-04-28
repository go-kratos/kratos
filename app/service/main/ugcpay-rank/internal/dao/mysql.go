package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/ugcpay-rank/internal/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_countElecUPRank      = "SELECT count(1) FROM rank_elec_up WHERE up_mid=? AND ver=? AND hidden=0"
	_countUPTotalElec     = "SELECT count(1) FROM rank_elec_up WHERE up_mid=? AND ver<>0 AND hidden=0"
	_selectElecUPRankList = "SELECT id,ver,up_mid,pay_mid,pay_amount,hidden,ctime,mtime FROM rank_elec_up WHERE up_mid=? AND ver=? AND hidden=0 ORDER BY pay_amount DESC,mtime ASC LIMIT ?"
	_selectElecUPRank     = "SELECT id,ver,up_mid,pay_mid,pay_amount,hidden,ctime,mtime FROM rank_elec_up WHERE up_mid=? AND ver=? AND pay_mid=? LIMIT 1"

	_countElecAVRank      = "SELECT count(1) FROM rank_elec_av WHERE avid=? AND ver=? AND hidden=0"
	_selectElecAVRankList = "SELECT id,ver,avid,up_mid,pay_mid,pay_amount,hidden,ctime,mtime FROM rank_elec_av WHERE avid=? AND ver=? AND hidden=0 ORDER BY pay_amount DESC,mtime ASC LIMIT ?"
	_selectElecAVRank     = "SELECT id,ver,avid,up_mid,pay_mid,pay_amount,hidden,ctime,mtime FROM rank_elec_av WHERE avid=? AND ver=? AND pay_mid=? LIMIT 1"

	_selectElecUPMessages      = "SELECT id,ver,avid,up_mid,pay_mid,message,replied,hidden FROM elec_message WHERE pay_mid in (%s) AND up_mid=? AND ver=? ORDER BY ID ASC"
	_selectElecAVMessagesByVer = "SELECT id,ver,avid,up_mid,pay_mid,message,replied,hidden FROM elec_message WHERE pay_mid in (%s) AND avid=? AND ver=? ORDER BY ID ASC"
	_selectElecAVMessages      = "SELECT id,ver,avid,up_mid,pay_mid,message,replied,hidden FROM elec_message WHERE pay_mid in (%s) AND avid=? ORDER BY ID ASC"

	_selectElecUPUserRank = "SELECT count(1) FROM rank_elec_up WHERE up_mid=? AND ver=? AND pay_amount>? AND mtime<? LIMIT 1"
	_selectElecAVUserRank = "SELECT count(1) FROM rank_elec_av WHERE avid=? AND ver=? AND pay_amount>? AND mtime<? LIMIT 1"
)

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// RawElecUPRankList get elec up rank
func (d *Dao) RawElecUPRankList(ctx context.Context, upMID int64, ver int64, limit int) (list []*model.DBElecUPRank, err error) {
	rows, err := d.db.Query(ctx, _selectElecUPRankList, upMID, ver, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.DBElecUPRank{}
		if err = rows.Scan(&data.ID, &data.Ver, &data.UPMID, &data.PayMID, &data.PayAmount, &data.Hidden, &data.CTime, &data.MTime); err != nil {
			return
		}
		list = append(list, data)
	}
	err = rows.Err()
	return
}

// RawElecUPRank get elec up rank
func (d *Dao) RawElecUPRank(ctx context.Context, upMID int64, ver int64, payMID int64) (data *model.DBElecUPRank, err error) {
	row := d.db.Master().QueryRow(ctx, _selectElecUPRank, upMID, ver, payMID)
	if err != nil {
		return
	}
	data = &model.DBElecUPRank{}
	if err = row.Scan(&data.ID, &data.Ver, &data.UPMID, &data.PayMID, &data.PayAmount, &data.Hidden, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// RawCountElecUPRank .
func (d *Dao) RawCountElecUPRank(ctx context.Context, upMID int64, ver int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countElecUPRank, upMID, ver)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// RawCountUPTotalElec stupid bug from prod.
func (d *Dao) RawCountUPTotalElec(ctx context.Context, upMID int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countUPTotalElec, upMID)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// RawElecAVRankList .
func (d *Dao) RawElecAVRankList(ctx context.Context, avID int64, ver int64, limit int) (list []*model.DBElecAVRank, err error) {
	rows, err := d.db.Query(ctx, _selectElecAVRankList, avID, ver, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.DBElecAVRank{}
		if err = rows.Scan(&data.ID, &data.Ver, &data.AVID, &data.UPMID, &data.PayMID, &data.PayAmount, &data.Hidden, &data.CTime, &data.MTime); err != nil {
			return
		}
		list = append(list, data)
	}
	err = rows.Err()
	return
}

// RawElecAVRank .
func (d *Dao) RawElecAVRank(ctx context.Context, avID int64, ver int64, payMID int64) (data *model.DBElecAVRank, err error) {
	rows := d.db.Master().QueryRow(ctx, _selectElecAVRank, avID, ver, payMID)
	if err != nil {
		return
	}
	data = &model.DBElecAVRank{}
	if err = rows.Scan(&data.ID, &data.Ver, &data.AVID, &data.UPMID, &data.PayMID, &data.PayAmount, &data.Hidden, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// RawCountElecAVRank .
func (d *Dao) RawCountElecAVRank(ctx context.Context, avID int64, ver int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countElecAVRank, avID, ver)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// RawElecUPMessages .
func (d *Dao) RawElecUPMessages(ctx context.Context, payMIDs []int64, upMID int64, ver int64) (dataMap map[int64]*model.DBElecMessage, err error) {
	dataMap = make(map[int64]*model.DBElecMessage)
	if len(payMIDs) <= 0 {
		return
	}
	sql := fmt.Sprintf(_selectElecUPMessages, xstr.JoinInts(payMIDs))
	rows, err := d.db.Query(ctx, sql, upMID, ver)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.DBElecMessage{}
		if err = rows.Scan(&data.ID, &data.Ver, &data.AVID, &data.UPMID, &data.PayMID, &data.Message, &data.Replied, &data.Hidden); err != nil {
			return
		}
		dataMap[data.PayMID] = data
	}
	return
}

// RawElecAVMessagesByVer  .
func (d *Dao) RawElecAVMessagesByVer(ctx context.Context, payMIDs []int64, avID int64, ver int64) (dataMap map[int64]*model.DBElecMessage, err error) {
	dataMap = make(map[int64]*model.DBElecMessage)
	if len(payMIDs) <= 0 {
		return
	}
	sql := fmt.Sprintf(_selectElecAVMessagesByVer, xstr.JoinInts(payMIDs))
	rows, err := d.db.Query(ctx, sql, avID, ver)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.DBElecMessage{}
		if err = rows.Scan(&data.ID, &data.Ver, &data.AVID, &data.UPMID, &data.PayMID, &data.Message, &data.Replied, &data.Hidden); err != nil {
			return
		}
		dataMap[data.PayMID] = data
	}
	return
}

// RawElecAVMessages .
func (d *Dao) RawElecAVMessages(ctx context.Context, payMIDs []int64, avID int64) (dataMap map[int64]*model.DBElecMessage, err error) {
	dataMap = make(map[int64]*model.DBElecMessage)
	if len(payMIDs) <= 0 {
		return
	}
	sql := fmt.Sprintf(_selectElecAVMessages, xstr.JoinInts(payMIDs))
	rows, err := d.db.Query(ctx, sql, avID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.DBElecMessage{}
		if err = rows.Scan(&data.ID, &data.Ver, &data.AVID, &data.UPMID, &data.PayMID, &data.Message, &data.Replied, &data.Hidden); err != nil {
			return
		}
		dataMap[data.PayMID] = data
	}
	return
}

// RawElecUPUserRank .
func (d *Dao) RawElecUPUserRank(ctx context.Context, upMID int64, ver int64, payAmount int64, mtime time.Time) (rank int, err error) {
	row := d.db.QueryRow(ctx, _selectElecUPUserRank, upMID, ver, payAmount, mtime)
	if err = row.Scan(&rank); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			rank = 0
		}
		return
	}
	return
}

// RawElecAVUserRank .
func (d *Dao) RawElecAVUserRank(ctx context.Context, avID int64, ver int64, payAmount int64, mtime time.Time) (rank int, err error) {
	row := d.db.QueryRow(ctx, _selectElecAVUserRank, avID, ver, payAmount, mtime)
	if err = row.Scan(&rank); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			rank = 0
		}
		return
	}
	return
}

const (
	_elecUserSettingList = "SELECT id,mid,value FROM elec_user_setting WHERE id>? ORDER BY id ASC LIMIT ?"
)

// RawElecUserSettings .
func (d *Dao) RawElecUserSettings(ctx context.Context, id int, limit int) (res map[int64]model.ElecUserSetting, maxID int, err error) {
	rows, err := d.db.Master().Query(ctx, _elecUserSettingList, id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64]model.ElecUserSetting)
	for rows.Next() {
		var (
			data model.ElecUserSetting
			mid  int64
			id   int
		)
		if err = rows.Scan(&id, &mid, &data); err != nil {
			err = errors.WithStack(err)
			return
		}
		res[mid] = data
		if maxID < id {
			maxID = id
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
