package dao

import (
	"context"
	"math"

	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
)

const (
	_oldElecOrderList   = `SELECT id,mid,pay_mid,order_no,elec_num,status,ctime,mtime FROM elec_pay_order WHERE id>? ORDER BY ID ASC LIMIT ?`
	_oldElecMessageList = `SELECT id,mid,ref_mid,IFNULL(ref_id,0),message,av_no,date_version,type,state,ctime,mtime FROM elec_message WHERE id>? ORDER BY ID ASC LIMIT ?`

	_upsertElecMessage = "INSERT INTO elec_message (id,`ver`,avid,up_mid,pay_mid,message,replied,hidden,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `ver`=?,avid=?,up_mid=?,pay_mid=?,message=?,replied=?,hidden=?"
	_upsertElecReply   = "INSERT INTO elec_reply (id,message_id,reply,hidden,ctime,mtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE message_id=?,reply=?,hidden=?"
	_upsertElecAVRank  = "INSERT INTO rank_elec_av (ver,avid,up_mid,pay_mid,pay_amount,hidden) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE pay_amount = pay_amount + ?"
	_upsertElecUPRank  = "INSERT INTO rank_elec_up (ver,up_mid,pay_mid,pay_amount,hidden) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE pay_amount = pay_amount + ?"

	_selectOldElecTradeInfo = "SELECT id,order_no,IFNULL(av_no,0) FROM elec_pay_trade_info WHERE order_no=? LIMIT 1"
)

// UpsertElecMessage .
func (d *Dao) UpsertElecMessage(ctx context.Context, data *model.DBElecMessage) (err error) {
	_, err = d.dbrank.Exec(ctx, _upsertElecMessage, data.ID, data.Ver, data.AVID, data.UPMID, data.PayMID, data.Message, data.Replied, data.Hidden, data.CTime, data.MTime, data.Ver, data.AVID, data.UPMID, data.PayMID, data.Message, data.Replied, data.Hidden)
	return
}

// UpsertElecReply .
func (d *Dao) UpsertElecReply(ctx context.Context, data *model.DBElecReply) (err error) {
	_, err = d.dbrank.Exec(ctx, _upsertElecReply, data.ID, data.MSGID, data.Reply, data.Hidden, data.CTime, data.MTime, data.MSGID, data.Reply, data.Hidden)
	return
}

// TXUpsertElecAVRank .
func (d *Dao) TXUpsertElecAVRank(ctx context.Context, tx *xsql.Tx, ver int64, avID int64, upMID, payMID int64, deltaPayAmount int64, hidden bool) (err error) {
	_, err = tx.Exec(_upsertElecAVRank, ver, avID, upMID, payMID, deltaPayAmount, hidden, deltaPayAmount)
	return
}

// TXUpsertElecUPRank .
func (d *Dao) TXUpsertElecUPRank(ctx context.Context, tx *xsql.Tx, ver int64, upMID, payMID int64, deltaPayAmount int64, hidden bool) (err error) {
	_, err = tx.Exec(_upsertElecUPRank, ver, upMID, payMID, deltaPayAmount, hidden, deltaPayAmount)
	return
}

// RawOldElecTradeInfo .
func (d *Dao) RawOldElecTradeInfo(ctx context.Context, orderID string) (data *model.DBOldElecPayTradeInfo, err error) {
	row := d.dbrankold.QueryRow(ctx, _selectOldElecTradeInfo, orderID)
	data = &model.DBOldElecPayTradeInfo{}
	if err = row.Scan(&data.ID, &data.OrderID, &data.AVID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// OldElecMessageList .
func (d *Dao) OldElecMessageList(ctx context.Context, startID int64, limit int) (maxID int64, list []*model.DBOldElecMessage, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.dbrankold.Query(ctx, _oldElecMessageList, startID, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p = &model.DBOldElecMessage{}
		)
		if err = rows.Scan(&p.ID, &p.MID, &p.RefMID, &p.RefID, &p.Message, &p.AVID, &p.DateVer, &p.Type, &p.State, &p.CTime, &p.MTime); err != nil {
			return
		}
		if maxID < p.ID {
			maxID = p.ID
		}
		list = append(list, p)
	}

	if err = rows.Err(); err != nil {
		return
	}
	return
}

// OldElecOrderList .
func (d *Dao) OldElecOrderList(ctx context.Context, startID int64, limit int) (maxID int64, list []*model.DBOldElecPayOrder, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.dbrankold.Query(ctx, _oldElecOrderList, startID, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p = &model.DBOldElecPayOrder{}
		)
		if err = rows.Scan(&p.ID, &p.UPMID, &p.PayMID, &p.OrderID, &p.ElecNum, &p.Status, &p.CTime, &p.MTime); err != nil {
			return
		}
		if maxID < p.ID {
			maxID = p.ID
		}
		list = append(list, p)
	}

	if err = rows.Err(); err != nil {
		return
	}
	return
}

const (
	_elecAddSetting    = `INSERT INTO elec_user_setting (mid,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=value|?`
	_elecDeleteSetting = `INSERT INTO elec_user_setting (mid,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=value&?`
)

// ElecAddSetting .
func (d *Dao) ElecAddSetting(ctx context.Context, defaultValue int32, mid int64, bitValue int32) (err error) {
	defaultValue |= bitValue
	_, err = d.dbrank.Exec(ctx, _elecAddSetting, mid, defaultValue, bitValue)
	return
}

// ElecDeleteSetting .
func (d *Dao) ElecDeleteSetting(ctx context.Context, defaultValue int32, mid int64, bitValue int32) (err error) {
	andValue := math.MaxInt32 ^ bitValue
	defaultValue &= andValue
	_, err = d.dbrank.Exec(ctx, _elecDeleteSetting, mid, defaultValue, andValue)
	return
}
