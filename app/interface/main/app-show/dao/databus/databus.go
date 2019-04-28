package databus

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-show/conf"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Dao is show dao.
type Dao struct {
	// databus
	dataBus *databus.Databus
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// databus
		dataBus: databus.New(c.DislikeDataBus),
	}
	return
}

func (d *Dao) Pub(ctx context.Context, buvid, gt string, id, mid int64) (err error) {
	key := strconv.FormatInt(mid, 10)
	msg := struct {
		Buvid string `json:"buvid"`
		Goto  string `json:"goto"`
		ID    int64  `json:"id"`
		Mid   int64  `json:"mid"`
	}{Buvid: buvid, Goto: gt, ID: id, Mid: mid}
	if err = d.dataBus.Send(ctx, key, msg); err != nil {
		log.Error("d.dataBus.Pub(%s,%v) error (%v)", key, msg, err)
	}
	return
}
