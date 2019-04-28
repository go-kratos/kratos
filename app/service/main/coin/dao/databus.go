package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"
)

// PubBigData pub msg into databus.
func (d *Dao) PubBigData(c context.Context, aid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(aid, 10)
	if err = d.dbBigData.Send(c, key, msg); err != nil {
		log.Error("dbBigData.Pub(%s, %v) error (%v)", key, msg, err)
		PromError("dbus:PubBigData")
	}
	return
}

// PubCoinJob pub job msg into databus.
func (d *Dao) PubCoinJob(c context.Context, aid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(aid, 10)
	for i := 0; i < 3; i++ {
		if err = d.dbCoinJob.Send(c, key, msg); err != nil {
			log.Error("d.dbCoinJob.Pub(%s, %v) error (%v) times: %v", key, msg, err, i+1)
			PromError("dbus:PubCoinJob")
			time.Sleep(time.Millisecond * 50)
			continue
		}
		break
	}
	return
}

// PubStat pub stat msg into databus.
func (d *Dao) PubStat(c context.Context, aid, tp, count int64) (err error) {
	var s = &struct {
		Type      string `json:"type"`
		ID        int64  `json:"id"`
		Count     int64  `json:"count"`
		Timestamp int64  `json:"timestamp"`
	}{
		ID:        aid,
		Count:     count,
		Timestamp: time.Now().Unix(),
	}
	// double write new databus.
	if b, ok := d.Businesses[tp]; ok {
		s.Type = b.Name
		if err = d.stat.Send(c, strconv.FormatInt(aid, 10), s); err != nil {
			log.Error("d.stat.Pub(%+v) error(%v)", s, err)
			PromError("dbus:stat")
		}
	}
	return
}
