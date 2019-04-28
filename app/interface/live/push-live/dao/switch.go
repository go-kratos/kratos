package dao

import (
	"context"
	"fmt"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_shard             = 10 //分表十张
	_getMidsByTargetID = "SELECT uid FROM app_switch_config_%s WHERE target_id=? AND type=? AND switch=?"
)

// tableIndex return index by target_id
func tableIndex(targetID int64) string {
	return fmt.Sprintf("%02d", targetID%_shard)
}

// GetFansBySwitch 获取直播开关数据
func (d *Dao) GetFansBySwitch(c context.Context, targetID int64) (fans map[int64]bool, err error) {
	var mid int64
	fans = make(map[int64]bool)
	sql := fmt.Sprintf(_getMidsByTargetID, tableIndex(targetID))
	rows, err := d.db.Query(c, sql, targetID, model.LivePushType, model.LivePushSwitchOn)
	if err != nil {
		err = errors.WithStack(err)
		fmt.Printf("%v", err)
		log.Error("[dao.switch|GetSwitchMids] db.Query() error(%v)", err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&mid); err != nil {
			err = errors.WithStack(err)
			fmt.Printf("%v", err)
			log.Error("[dao.switch|GetSwitchMids] rows.Scan() error(%v)", err)
			return
		}
		fans[mid] = true
	}
	return
}
