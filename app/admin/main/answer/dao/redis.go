package dao

import (
	"context"

	"go-common/library/log"
)

const (
	_baseQsIdsKey = "v3_tc_bqs"
)

// SetQidCache set question id into question set
func (d *Dao) SetQidCache(c context.Context, id int64) (err error) {
	var (
		key  = _baseQsIdsKey
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SET", key, id); err != nil {
		log.Error("conn.Send(SET, %s, %d) error(%v)", key, id, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}
