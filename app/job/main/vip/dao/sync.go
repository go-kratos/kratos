package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/job/main/vip/model"

	"github.com/pkg/errors"
)

const (
	//sync
	_syncAddUser    = "INSERT INTO vip_user_info(mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_recent_time,ios_overdue_time,pay_channel_id,ctime,mtime,ver) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_syncUpdateUser = "UPDATE vip_user_info SET vip_type=?,vip_pay_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,vip_recent_time=?,ios_overdue_time=?,pay_channel_id=?,ctime=?,mtime=?,ver=? WHERE mid=? AND ver=?"
)

//SyncAddUser insert vipUserInfo
func (d *Dao) SyncAddUser(c context.Context, r *model.VipUserInfo) (err error) {
	if _, err = d.db.Exec(c, _syncAddUser, r.Mid, r.Type, r.PayType, r.Status, r.StartTime, r.OverdueTime, r.AnnualVipOverdueTime, r.RecentTime, r.OverdueTime, r.PayChannelID, r.Ctime, r.Mtime, r.Ver); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SyncUpdateUser insert vipUserInfo
func (d *Dao) SyncUpdateUser(c context.Context, r *model.VipUserInfo, ver int64) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _syncUpdateUser, r.Type, r.PayType, r.Status, r.OverdueTime, r.AnnualVipOverdueTime, r.RecentTime, r.IosOverdueTime, r.PayChannelID, r.Ctime, r.Mtime, r.Ver, r.Mid, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
