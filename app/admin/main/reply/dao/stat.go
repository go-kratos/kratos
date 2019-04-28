package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"
)

type statMsg struct {
	ID        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Count     int32  `json:"count"`
	Type      string `json:"type"`
}

// SendStats update stat.
func (d *Dao) SendStats(c context.Context, typ int32, oid int64, cnt int32) (err error) {
	// new databus stats
	if name, ok := d.statsTypes[typ]; ok {
		m := &statMsg{
			ID:        oid,
			Type:      name,
			Count:     cnt,
			Timestamp: time.Now().Unix(),
		}
		if err = d.statsBus.Send(c, strconv.FormatInt(oid, 10), m); err != nil {
			log.Error("d.databus.Send(%d,%d,%d) error(%v)", typ, oid, cnt, err)
		}
	}
	return
}
