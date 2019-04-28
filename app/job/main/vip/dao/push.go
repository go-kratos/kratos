package dao

import (
	"context"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

//PushDatas get push datas
func (d *Dao) PushDatas(c context.Context, curtime string) (res []*model.VipPushData, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selPushDataSQL, curtime, curtime); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipPushData)
		if err = rows.Scan(&r.ID, &r.DisableType, &r.GroupName, &r.Title, &r.Content, &r.PushTotalCount, &r.PushedCount, &r.ProgressStatus, &r.Status, &r.Platform, &r.LinkType, &r.LinkURL, &r.ErrorCode, &r.ExpiredDayStart, &r.ExpiredDayEnd, &r.EffectStartDate, &r.EffectEndDate, &r.PushStartTime, &r.PushEndTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//UpdatePushData update push data
func (d *Dao) UpdatePushData(c context.Context, status, progressStatus int8, pushedCount int32, errcode, data, id int64) (err error) {
	if _, err = d.db.Exec(c, _updatePushDataSQL, progressStatus, status, pushedCount, errcode, data, id); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_err")
	}
	return
}
