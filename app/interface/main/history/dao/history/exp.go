package history

import (
	"context"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_keyFirst  = "his_f_"
	_timeMonth = "0102"
	_view      = "view"
)

type experienceMsg struct {
	Event string `json:"event"`
	Mid   int64  `json:"mid"`
	IP    string `json:"ip"`
	TS    int64  `json:"ts"`
}

// keyFirst return first key
func keyFirst(mid int64, t string) string {
	return _keyFirst + strconv.FormatInt(mid%1000, 10) + "_" + t
}

// PushFirstQueue push first view record every day into kafka.
func (d *Dao) PushFirstQueue(c context.Context, mid, aid, now int64) (err error) {
	var (
		today = time.Unix(now, 0)
		md    = today.Format(_timeMonth)
		key   = keyFirst(mid, md)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	rp, err := redis.Int(conn.Do("SISMEMBER", key, mid))
	if err != nil {
		log.Error("conn.Do(SISMEMBER, %s, %d) error(%v)", key, mid, err)
		err = nil
	}
	// if key exist , donot push to kafka
	if rp > 0 {
		return
	}
	midStr := strconv.FormatInt(mid, 10)
	ex := &experienceMsg{
		Event: _view,
		Mid:   mid,
		IP:    metadata.String(c, metadata.RemoteIP),
		TS:    now,
	}
	err = d.experiencePub(c, midStr, ex)
	conn.Send("SADD", key, mid)
	conn.Send("EXPIRE", key, 24*60*60)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive1() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive2() error(%v)", err)
		return
	}
	return
}
