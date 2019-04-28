package history

import (
	"context"
	"strconv"

	"go-common/library/log"
)

// PlayPro send history to databus.
func (d *Dao) PlayPro(c context.Context, key string, msg interface{}) (err error) {
	if err = d.playPro.Send(c, key, msg); err != nil {
		log.Error("d.pub.Pub(%s,%v) error(%v)", key, msg, err)
	}
	return
}

// Merge send history to databus.
func (d *Dao) Merge(c context.Context, mid int64, msg interface{}) (err error) {
	key := strconv.FormatInt(mid, 10)
	if err = d.merge.Send(c, key, msg); err != nil {
		log.Error("d.pub.Pub(%s,%v) error(%v)", key, msg, err)
	}
	return
}

// experiencePub send history to databus.
func (d *Dao) experiencePub(c context.Context, key string, msg interface{}) (err error) {
	if err = d.experience.Send(c, key, msg); err != nil {
		log.Error("d.pub.Pub(%s,%v) error(%v)", key, msg, err)
	}
	return
}

// ProPub send history to databus.
func (d *Dao) ProPub(c context.Context, key string, msg interface{}) (err error) {
	if err = d.proPub.Send(c, key, msg); err != nil {
		log.Error("d.ProPub.Pub(%s,%v) error(%v)", key, msg, err)
	}
	return
}
