package dao

import (
	"context"
	"strconv"

	"go-common/library/log"
)

// PubScoreChange pub spy score change msg into databus.
func (d *Dao) PubScoreChange(c context.Context, mid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(mid, 10)
	if err = d.dbScoreChange.Send(c, key, msg); err != nil {
		log.Error("d.dbScoreChange.Send(%s, %v) error (%v)", key, msg, err)
	}
	return
}
