package dao

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_pendantPKG   = "pkg_"
	_pendantEquip = "pe_"
)

func keyPendantPKG(uid int64) string {
	return _pendantPKG + strconv.FormatInt(uid, 10)
}

func keyEquip(uid int64) string {
	return _pendantEquip + strconv.FormatInt(uid, 10)
}

// DelPKGCache del package cache
func (d *Dao) DelPKGCache(c context.Context, uids []int64) (err error) {
	var (
		args = redis.Args{}
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, v := range uids {
		args = args.Add(keyPendantPKG(v))
	}
	if err = conn.Send("DEL", args...); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", args, err)
	}
	return
}

// DelEquipsCache del batch equip cache .
func (d *Dao) DelEquipsCache(c context.Context, uids []int64) (err error) {
	var (
		args = redis.Args{}
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, v := range uids {
		args = args.Add(keyEquip(v))
	}
	if err = conn.Send("DEL", args...); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", args, err)
	}
	return
}
