package dao

import (
	"context"
	"go-common/library/log"
	"strconv"
)

func (d *Dao) Pub(c context.Context, uid int64, msg interface{}) error {
	key := strconv.FormatInt(uid, 10)
	err := d.PushSearchDataBus.Send(c, key, msg)
	if err != nil {
		log.Error("pub wallet change failed uid:%d, msg:%+v", uid, msg)
	}
	return err
}
