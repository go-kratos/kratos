package dao

import (
	"bytes"
	"sync"

	"go-common/app/interface/main/report-click/conf"
	"go-common/library/queue/databus"
)

// Dao report-click dao
type Dao struct {
	c       *conf.Config
	merge   *databus.Databus
	msgs    chan []byte
	spliter []byte
	bfp     sync.Pool
}

// New dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:       c,
		merge:   databus.New(c.DataBus.Merge),
		msgs:    make(chan []byte, 1024),
		spliter: []byte("\001"),
		bfp: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer([]byte{})
			},
		},
	}
	go d.pubproc()
	return
}

// Close close kafka connection.
func (d *Dao) Close() {
	d.msgs <- d.spliter
	if d.merge != nil {
		d.merge.Close()
	}
}
