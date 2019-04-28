package dao

import (
	"context"
	"database/sql"

	"go-common/app/service/main/tv/internal/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_getUserInfoByMid = "SELECT `id`, `mid`, `ver`, `vip_type`, `pay_type`, `pay_channel_id`, `status`, `overdue_time`, `recent_pay_time`, `ctime`, `mtime` FROM `tv_user_info` WHERE `mid`=?"

	_insertUserInfo = "INSERT INTO tv_user_info (`mid`, `ver`, `vip_type`, `pay_type`, `pay_channel_id`, `status`, `overdue_time`, `recent_pay_time`) VALUES (?,?,?,?,?,?,?,?)"

	_updateUserInfo    = "UPDATE tv_user_info SET `status` = ?, `vip_type` = ?,  `overdue_time`=?, `recent_pay_time`=?, `ver` = `ver` + 1 WHERE `mid` = ? AND `ver` = ?"
	_updateUserStatus  = "UPDATE tv_user_info SET `status` = ?, `ver` = `ver` + 1  WHERE `mid` = ? AND `ver` = ?"
	_updateUserPayType = "UPDATE tv_user_info SET `pay_type` = ?, `ver` = `ver` + 1  WHERE `mid` = ? AND `ver` = ?"
)

// UserInfoByMid quires one row from tv_user_info.
func (d *Dao) RawUserInfoByMid(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	row := d.db.QueryRow(c, _getUserInfoByMid, mid)
	ui = &model.UserInfo{}
	err = row.Scan(&ui.ID, &ui.Mid, &ui.Ver, &ui.VipType, &ui.PayType, &ui.PayChannelId, &ui.Status, &ui.OverdueTime, &ui.RecentPayTime, &ui.Ctime, &ui.Mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getUserInfoByMid, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return ui, nil
}

// TxInsertUserInfo insert one row into tv_user_info.
func (d *Dao) TxInsertUserInfo(ctx context.Context, tx *xsql.Tx, ui *model.UserInfo) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertUserInfo, ui.Mid, ui.Ver, ui.VipType, ui.PayType, ui.PayChannelId, ui.Status, ui.OverdueTime, ui.RecentPayTime); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _insertUserInfo, err)
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdateUserInfo updates user info.
func (d *Dao) TxUpdateUserInfo(ctx context.Context, tx *xsql.Tx, ui *model.UserInfo) (err error) {
	if _, err = tx.Exec(_updateUserInfo, ui.Status, ui.VipType, ui.OverdueTime, ui.RecentPayTime, ui.Mid, ui.Ver); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _updateUserInfo, err)
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdateUserInfo updates vip status of user.
func (d *Dao) TxUpdateUserStatus(ctx context.Context, tx *xsql.Tx, ui *model.UserInfo) (err error) {
	if _, err = tx.Exec(_updateUserStatus, ui.Status, ui.Mid, ui.Ver); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _updateUserStatus, err)
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdateUserPayType updates pay type of user.
func (d *Dao) TxUpdateUserPayType(ctx context.Context, tx *xsql.Tx, ui *model.UserInfo) (err error) {
	if _, err = tx.Exec(_updateUserPayType, ui.PayType, ui.Mid, ui.Ver); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _updateUserPayType, err)
		err = errors.WithStack(err)
		return
	}
	return
}
