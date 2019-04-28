package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/service/main/vipinfo/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_byMidSQL  = "SELECT id,mid,vip_type,vip_pay_type,vip_status,pay_channel_id,vip_start_time,vip_overdue_time,annual_vip_overdue_time,ctime,mtime,vip_recent_time,ios_overdue_time,ver FROM vip_user_info WHERE mid = ?;"
	_byMidsSQL = "SELECT id,mid,vip_type,vip_pay_type,vip_status,pay_channel_id,vip_start_time,vip_overdue_time,annual_vip_overdue_time,ctime,mtime,vip_recent_time,ios_overdue_time,ver FROM vip_user_info WHERE mid IN (%s);"
)

//RawInfo select user info by mid.
func (d *Dao) RawInfo(c context.Context, mid int64) (r *model.VipUserInfo, err error) {
	var row = d.db.QueryRow(c, _byMidSQL, mid)
	r = new(model.VipUserInfo)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipPayType, &r.VipStatus, &r.PayChannelID, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime,
		&r.Ctime, &r.Mtime, &r.VipRecentTime, &r.IosOverdueTime, &r.Ver); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao info bymid(%d)", mid)
	}
	return
}

// RawInfos get user infos.
func (d *Dao) RawInfos(c context.Context, mids []int64) (res map[int64]*model.VipUserInfo, err error) {
	var rows *xsql.Rows
	res = make(map[int64]*model.VipUserInfo, len(mids))
	midStr := xstr.JoinInts(mids)
	if rows, err = d.db.Query(c, fmt.Sprintf(_byMidsSQL, midStr)); err != nil {
		err = errors.Wrapf(err, "dao infos mids(%s)", midStr)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipPayType, &r.VipStatus, &r.PayChannelID, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime,
			&r.Ctime, &r.Mtime, &r.VipRecentTime, &r.IosOverdueTime, &r.Ver); err != nil {
			err = errors.Wrapf(err, "dao infos scan mids(%s)", midStr)
			res = nil
			return
		}
		res[r.Mid] = r
	}
	err = rows.Err()
	return
}
