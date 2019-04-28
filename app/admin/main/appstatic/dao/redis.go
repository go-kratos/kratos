package dao

import (
	"context"
	"time"

	"go-common/library/log"
)

const _pushKey = "appstatic-admin-topush"

// ZAddPush adds one to push data into the redis sorted set
func (d *Dao) ZAddPush(c context.Context, resID int) (err error) {
	var (
		conn  = d.redis.Get(c)
		ctime = time.Now().Unix()
	)
	defer conn.Close()
	if err = conn.Send("ZADD", _pushKey, ctime, resID); err != nil {
		log.Error("conn.Send(ZADD %s - %v) error(%v)", _pushKey, resID, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}
