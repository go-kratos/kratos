package dao

import (
	"context"
	"strconv"

	"go-common/library/log"
)

// PubExtraRet pub extra msg into databus.
func (d *Dao) PubExtraRet(c context.Context, mid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(mid, 10)
	if err = d.dbExtraAnswerRet.Send(c, key, msg); err != nil {
		log.Error("d.dbExtraAnswerRet.Send(%s, %v) error (%v)", key, msg, err)
	}
	return
}
