package up

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-common/library/xstr"

	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
)

const (
	// _getUpInfoActiveSQL .
	_getUpInfoActiveSQL = "SELECT id, mid, active_tid FROM up_base_info WHERE mid = ?"
	// _getUpsInfoActiveSQL .
	_getUpsInfoActiveSQL = "SELECT id, mid, active_tid FROM up_base_info WHERE mid IN (%s)"
)

// RawUpInfoActive get up info active
func (d *Dao) RawUpInfoActive(ctx context.Context, mid int64) (upInfoActive *model.UpInfoActiveReply, err error) {
	row := global.GetUpCrmDB().QueryRow(ctx, _getUpInfoActiveSQL, mid)
	upInfoActive = &model.UpInfoActiveReply{}
	if err = row.Scan(&upInfoActive.ID, &upInfoActive.MID, &upInfoActive.ActiveTid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("RawUpInfoActive row.Scan error(%v)", err)
			err = errors.New("RawUpInfoActive get data failed")
		}
	}
	return
}

// RawUpsInfoActive get ups info active
func (d *Dao) RawUpsInfoActive(ctx context.Context, mids []int64) (res map[int64]*model.UpInfoActiveReply, err error) {
	res = make(map[int64]*model.UpInfoActiveReply)
	sql := fmt.Sprintf(_getUpsInfoActiveSQL, xstr.JoinInts(mids))
	log.Info("SQL: %s", sql)
	rows, err := global.GetUpCrmDB().Query(ctx, sql)
	if err != nil {
		log.Error("RawUpsInfoActive UpCrmDB.Query error(%v)", err)
		return
	}

	for rows.Next() {
		upInfoActive := model.UpInfoActiveReply{}
		if err = rows.Scan(&upInfoActive.ID, &upInfoActive.MID, &upInfoActive.ActiveTid); err != nil {
			log.Error("RawUpsInfoActive rows.Scan error(%v)", err)
			return
		}
		res[upInfoActive.MID] = &upInfoActive
	}

	return
}
