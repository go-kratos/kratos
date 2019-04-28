package pub

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"
)

type statMessage struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// PubStats update object's fav count
func (d *Dao) PubStats(c context.Context, typ int8, oid int64, cnt int64) (err error) {
	if name, ok := d.consumersMap[typ]; ok {
		msg := &statMessage{
			Type:      name,
			ID:        oid,
			Count:     cnt,
			TimeStamp: time.Now().Unix(),
		}
		if err = d.databus2.Send(c, strconv.FormatInt(oid, 10), msg); err != nil {
			log.Error("d.databus2.Send(%d,%d,%v) error(%v)", typ, oid, msg, err)
		}
	}
	return
}
