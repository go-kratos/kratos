package dao

import (
	"context"
	"strconv"

	"go-common/library/log"
)

// PubLabour pub labour answer log msg into databus.
func (d *Dao) PubLabour(c context.Context, aid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(aid, 10)
	if err = d.dbusLabour.Send(c, key, msg); err != nil {
		log.Error("PubLabour.Pub(%s, %v) error (%v)", key, msg, err)
	}
	return
}
