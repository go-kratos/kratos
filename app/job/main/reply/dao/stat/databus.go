package stat

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/reply/conf"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

type statMsg struct {
	ID        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Count     int    `json:"count"`
	Type      string `json:"type"`
}

// Dao stat dao.
type Dao struct {
	// new databus stats
	types   map[int8]string
	databus *databus.Databus
}

// New new a stat dao and return.
func New(c *conf.Config) *Dao {
	d := new(Dao)
	// new databus stats
	d.types = make(map[int8]string)
	for name, typ := range c.StatTypes {
		d.types[typ] = name
	}
	d.databus = databus.New(c.Databus.Stats)
	return d
}

// Send update stat.
func (d *Dao) Send(c context.Context, typ int8, oid int64, cnt int) (err error) {
	// new databus stats
	if name, ok := d.types[typ]; ok {
		m := &statMsg{
			ID:        oid,
			Type:      name,
			Count:     cnt,
			Timestamp: time.Now().Unix(),
		}
		if err = d.databus.Send(c, strconv.FormatInt(oid, 10), m); err != nil {
			log.Error("d.databus.Send(%d,%d,%d) error(%v)", typ, oid, cnt, err)
		}
	}
	return
}
